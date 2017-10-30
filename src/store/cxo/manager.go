package cxo

import (
	"context"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/msgs"
	"github.com/skycoin/bbs/src/store/cxo/setup"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/object/revisions/r0"
	"github.com/skycoin/bbs/src/store/object/transfer"
	"github.com/skycoin/bbs/src/store/state"
	"github.com/skycoin/bbs/src/store/state/views"
	"github.com/skycoin/bbs/src/store/state/views/content_view"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/node/gnet"
	"github.com/skycoin/cxo/node/log"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/file"
	"io/ioutil"
	log2 "log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const (
	LogPrefix     = "CXO"
	SubDir        = "cxo_r1"
	FileName      = "bbs.json"
	ExportSubDir  = "exports"
	ExportFileExt = ".export"
	RetryDuration = time.Second * 5
)

// ManagerConfig represents the configuration for CXO Manager.
type ManagerConfig struct {
	Memory             *bool    // Whether to enable memory mode.
	Defaults           *bool    // Whether to have default connection / subscription.
	MessengerAddresses []string // Messenger addresses.
	Config             *string  // Configuration directory.
	CXOPort            *int     // CXO listening port.
	CXORPCEnable       *bool    // Whether to enable CXO RPC.
	CXORPCPort         *int     // CXO RPC port.
}

// Manager manages interaction with CXO and storing/retrieving node configuration files.
type Manager struct {
	mux      sync.Mutex
	c        *ManagerConfig
	l        *log2.Logger
	file     *object.CXOFileManager
	node     *node.Node
	compiler *state.Compiler
	relay    *msgs.Relay
	wg       sync.WaitGroup
	newRoots chan state.RootWrap
	quit     chan struct{}
}

// NewManager creates a new CXO manager.
func NewManager(config *ManagerConfig, compilerConfig *state.CompilerConfig, relayConfig *msgs.RelayConfig) *Manager {
	manager := &Manager{
		c: config,
		l: inform.NewLogger(true, os.Stdout, LogPrefix),
		file: object.NewCXOFileManager(&object.CXOFileManagerConfig{
			Memory:   config.Memory,
			Defaults: config.Defaults,
		}),
		relay:    msgs.NewRelay(relayConfig),
		newRoots: make(chan state.RootWrap, 10),
		quit:     make(chan struct{}),
	}

	// Prepare CXO node.
	if e := manager.prepareNode(); e != nil {
		manager.l.Panicln("failed to start CXO manager:", e)
	}

	// Prepare CXO compiler.
	manager.compiler = state.NewCompiler(
		compilerConfig, manager.file, manager.newRoots, manager.node,
		views.AddContent(),
		views.AddFollow(),
	)

	// Prepare messenger relay.
	if e := manager.relay.Open(manager.compiler); e != nil {
		manager.l.Panicln("failed to start CXO manager:", e)
	}

	// Prepare CXO file.
	if e := manager.prepareFile(); e != nil {
		manager.l.Panicln("failed to start CXO manager:", e)
	}

	// Init directories.
	if !*config.Memory {
		if e := os.MkdirAll(
			path.Join(*config.Config, SubDir),
			os.FileMode(0700),
		); e != nil {
			manager.l.Panicln(e)
		}
		if e := os.MkdirAll(
			path.Join(*config.Config, ExportSubDir),
			os.FileMode(0700),
		); e != nil {
			manager.l.Panicln(e)
		}
	}

	go manager.retryLoop()
	return manager
}

// Close quits the CXO manager.
func (m *Manager) Close() {
	for {
		select {
		case m.quit <- struct{}{}:
		default:
			m.relay.Close()
			m.compiler.Close()
			m.wg.Wait()
			if e := m.node.Close(); e != nil {
				m.l.Println("Error on close:", e.Error())
			}
			<-m.node.Quiting()
			return
		}
	}
}

// Relay obtains the messenger relay.
func (m *Manager) Relay() *msgs.Relay {
	return m.relay
}

/*
	<<< HELPER FUNCTIONS >>>
*/

// Sets up the CXO node.
func (m *Manager) prepareNode() error {
	c := node.NewConfig()

	c.Log.Prefix = "[CXO] "
	c.Log.Debug = false
	c.Log.Pins = log.All // all
	c.Config.Logger = log.NewLogger(log.Config{Output: ioutil.Discard})

	c.Skyobject.Log.Debug = true
	c.Skyobject.Log.Pins = skyobject.PackSavePin // all
	c.Skyobject.Log.Prefix = "[CXO] "

	// CXO schema are registered here.
	c.Skyobject.Registry = skyobject.NewRegistry(
		setup.PrepareRegistry)

	//c.MaxMessageSize = 0 // TODO -> Adjust.
	c.InMemoryDB = *m.c.Memory
	c.DataDir = filepath.Join(*m.c.Config, SubDir)
	c.EnableListener = true
	c.PublicServer = true
	c.DiscoveryAddresses = m.c.MessengerAddresses
	c.Listen = "[::]:" + strconv.Itoa(*m.c.CXOPort)
	c.EnableRPC = *m.c.CXORPCEnable
	c.RemoteClose = false
	c.RPCAddress = "[::]:" + strconv.Itoa(*m.c.CXORPCPort)

	c.OnRootReceived = func(c *node.Conn, root *skyobject.Root) {
		m.l.Printf("Receiving board '%s' (NOT FILLED)", root.Pub.Hex()[:5]+"...")
		root.IsFull = false
		m.newRoots <- state.RootWrap{Root: root}
	}

	c.OnRootFilled = func(c *node.Conn, root *skyobject.Root) {
		m.l.Printf("Receiving board '%s' (FILLED)", root.Pub.Hex()[:5]+"...")
		root.IsFull = true
		m.newRoots <- state.RootWrap{Root: root}
	}

	c.OnCreateConnection = func(c *node.Conn) {
		for _, pk := range m.node.Feeds() {
			c.Subscribe(pk)
		}
	}

	c.OnCloseConnection = func(c *node.Conn) {
		m.node.Connect(c.Address())
	}

	var e error
	if m.node, e = node.NewNode(c); e != nil {
		return e
	}
	return nil
}

// Sets up CXO file.
func (m *Manager) prepareFile() error {
	if e := m.file.Load(m.filePath()); e != nil {
		return e
	}
	if e := m.file.RangeMasterSubs(func(pk cipher.PubKey, sk cipher.SecKey) {
		if r, e := m.node.Container().LastRoot(pk); e != nil {
			m.l.Println("prepareFile() LastRoot failed with error:", e)
		} else {
			m.compiler.UpdateBoard(r)
		}
	}); e != nil {
		return e
	}
	if e := m.file.RangeRemoteSubs(func(pk cipher.PubKey) {
		if e := m.subscribeNode(pk); e != nil {
			m.l.Println("prepareFile() subscribeNode() failed with error:", e)
		}
		if r, e := m.node.Container().LastRoot(pk); e != nil {
			m.l.Println("prepareFile() LastRoot failed with error:", e)
		} else {
			m.compiler.UpdateBoard(r)
		}
	}); e != nil {
		return e
	}
	return nil
}

func (m *Manager) filePath() string {
	return path.Join(*m.c.Config, SubDir, FileName)
}

func (m *Manager) exportPath(name string) string {
	return path.Join(*m.c.Config, ExportSubDir, name+ExportFileExt)
}

func (m *Manager) retryLoop() {
	m.wg.Add(1)
	defer m.wg.Done()

	ticker := time.NewTicker(RetryDuration)
	defer ticker.Stop()

	for {
		select {
		case <-m.quit:
			m.file.Save(m.filePath())
			return
		case <-ticker.C:
			m.file.Save(m.filePath())
		}
	}
}

/*
	<<< CONNECTION >>>
*/

func (m *Manager) GetConnections() []r0.Connection {
	connections := m.node.Connections()
	out := make([]r0.Connection, len(connections))
	for i, conn := range m.node.Connections() {
		out[i] = r0.Connection{
			Address: conn.Address(),
			State:   conn.Gnet().State().String(),
		}
	}
	return out
}

func (m *Manager) Connect(address string) error {
	if _, e := m.connectNode(address); e != nil {
		switch boo.Type(e) {
		case boo.AlreadyExists:
			return nil
		default:
			return e
		}
	}
	return nil
}

func (m *Manager) connectNode(address string) (*gnet.Conn, error) {
	if connection, e := m.node.Pool().Dial(address); e != nil {
		switch e {
		case gnet.ErrClosed, gnet.ErrConnectionsLimit:
			return nil, boo.WrapType(e, boo.Internal)
		case gnet.ErrAlreadyListen:
			return nil, boo.WrapType(e, boo.AlreadyExists)
		default:
			return nil, boo.WrapType(e, boo.InvalidInput)
		}
	} else {
		return connection, nil
	}
}

func (m *Manager) Disconnect(address string) error {
	if e := m.disconnectNode(address); e != nil {
		return e
	}
	return nil
}

func (m *Manager) disconnectNode(address string) error {
	conn := m.node.Pool().Connection(address)
	if conn == nil {
		return nil
	}
	if e := conn.Close(); e != nil {
		return e
	}
	return nil
}

/*
	<<< SUBSCRIPTION >>>
*/

func (m *Manager) GetSubscriptions() []cipher.PubKey {
	out, _ := m.file.GetSubKeyList()
	return out
}

func (m *Manager) SubscribeRemote(bpk cipher.PubKey) error {
	if e := m.file.AddRemoteSub(bpk); e != nil {
		return e
	}
	m.subscribeNode(bpk)
	return nil
}

func (m *Manager) SubscribeMaster(bpk cipher.PubKey, bsk cipher.SecKey) error {
	if e := m.file.AddMasterSub(bpk, bsk); e != nil {
		return e
	}
	m.subscribeNode(bpk)
	return nil
}

func (m *Manager) subscribeNode(bpk cipher.PubKey) error {
	if e := m.node.AddFeed(bpk); e != nil {
		return boo.WrapType(e, boo.Internal,
			"failed to add feed")
	}
	m.l.Printf("Subscribing to feed '%s'...", bpk.Hex()[:5]+"...")
	connections, count := m.node.Connections(), 0
	for _, conn := range connections {
		if e := conn.Subscribe(bpk); e != nil {
			switch e {
			case node.ErrSubscriptionRejected:
				m.l.Printf(" - '%s' rejected feed '%s'",
					conn.Address(), bpk.Hex()[:5]+"...")
			default:
				return boo.WrapType(e, boo.Internal,
					"failed to subscribe")
			}
		} else {
			count += 1
			m.l.Printf(" - '%s' accepted feed '%s'",
				conn.Address(), bpk.Hex()[:5]+"...")
		}
	}
	m.l.Printf("Feed '%s' accepted on %d/%d connections",
		bpk.Hex()[:5]+"...", count, len(connections))
	return nil
}

func (m *Manager) UnsubscribeRemote(bpk cipher.PubKey) error {
	if m.file.HasRemoteSub(bpk) {
		if e := m.file.RemoveSub(bpk); e != nil {
			return e
		}
		m.unsubscribeNode(bpk)
		return nil
	}
	return boo.Newf(boo.NotFound,
		"remote board of public key '%s' not found in cxo file", bpk.Hex()[:5]+"...")
}

func (m *Manager) UnsubscribeMaster(bpk cipher.PubKey) error {
	if m.file.HasMasterSub(bpk) {
		if e := m.file.RemoveSub(bpk); e != nil {
			return e
		}
		m.unsubscribeNode(bpk)
		return nil
	}
	return boo.Newf(boo.NotFound,
		"master board of public key '%s' not found in cxo file", bpk.Hex()[:5]+"...")
}

func (m *Manager) unsubscribeNode(bpk cipher.PubKey) {
	m.compiler.DeleteBoard(bpk)
	m.node.DelFeed(bpk)
}

/*
	<<< CONTENT >>>
*/

func (m *Manager) GetBoardInstance(bpk cipher.PubKey) (*state.BoardInstance, error) {
	return m.compiler.GetBoard(bpk)
}

func (m *Manager) GetBoards(ctx context.Context) ([]interface{}, []interface{}, error) {

	var masterOut []interface{}
	m.file.RangeMasterSubs(func(pk cipher.PubKey, sk cipher.SecKey) {
		bi, e := m.compiler.GetBoard(pk)
		if e != nil {
			m.l.Println(e)
			return
		}
		bView, e := bi.Get(views.Content, content_view.Board)
		if e != nil {
			m.l.Println(e)
			return
		}
		masterOut = append(masterOut, bView)
	})

	var remoteOut []interface{}
	m.file.RangeRemoteSubs(func(pk cipher.PubKey) {
		bi, e := m.compiler.GetBoard(pk)
		if e != nil {
			m.l.Println(e)
			return
		}
		bView, e := bi.Get(views.Content, content_view.Board)
		if e != nil {
			m.l.Println(e)
			return
		}
		remoteOut = append(remoteOut, bView)
	})

	return masterOut, remoteOut, nil
}

func (m *Manager) NewBoard(in *object.NewBoardIO) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	if e := m.file.AddMasterSub(in.BoardPubKey, in.BoardSecKey); e != nil {
		return e
	}
	m.subscribeNode(in.BoardPubKey)

	if r, e := setup.NewBoard(m.node, in); e != nil {
		return e
	} else {
		m.compiler.UpdateBoardWithContext(context.Background(), r)
		return nil
	}
}

/*
	<<< ADMIN >>>
*/

//func (m *Manager) GetObject(hash cipher.SHA256)

/*
	<<< DISCOVERER >>>
*/

func (m *Manager) GetDiscoveredBoards() []string {
	return m.relay.GetBoards()
}

/*
	<<< IMPORT / EXPORT >>>
*/

func (m *Manager) ExportBoard(pk cipher.PubKey, name string) (string, *transfer.RootRep, error) {
	if *m.c.Memory {
		return "", nil, nil
	}
	bi, e := m.compiler.GetBoard(pk)
	if e != nil {
		return "", nil, e
	}
	out, e := bi.Export()
	if e != nil {
		return "", nil, e
	}
	path := m.exportPath(name)
	if e := file.SaveJSON(path, out, os.FileMode(0600)); e != nil {
		return "", nil, e
	}
	return path, out, nil
}

func (m *Manager) ImportBoard(pk cipher.PubKey, name string) (string, *transfer.RootRep, error) {
	if *m.c.Memory {
		return "", nil, nil
	}

	path := m.exportPath(name)
	out := new(transfer.RootRep)
	if e := file.LoadJSON(path, out); e != nil {
		return "", nil, e
	}

	_, has := m.file.GetMasterSubSecKey(pk)
	if !has {
		return "", nil, boo.Newf(boo.NotAuthorised,
			"this node is not the master of board of public key '%s'", pk.Hex())
	}
	bi, e := m.compiler.GetBoard(pk)
	if e != nil {
		return "", nil, e
	}
	if e := bi.Import(out); e != nil {
		return "", nil, e
	}
	return path, out, nil
}
