package cxo

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/object/revisions/r0"
	"github.com/skycoin/bbs/src/store/state"
	"github.com/skycoin/bbs/src/store/state/views"
	"github.com/skycoin/bbs/src/store/state/views/content_view"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/node/gnet"
	"github.com/skycoin/cxo/node/log"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
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
	SubDir        = "cxo"
	FileName      = "bbs.json"
	ExportSubDir  = "exports"
	ExportFileExt = ".export"
	RetryDuration = time.Second * 5
)

// ManagerConfig represents the configuration for CXO Manager.
type ManagerConfig struct {
	Memory       *bool   // Whether to enable memory mode.
	Master       *bool   // Whether node is to host boards and submission address.
	Defaults     *bool   // Whether to have default connection / subscription.
	Config       *string // Configuration directory.
	CXOPort      *int    // CXO listening port.
	CXORPCEnable *bool   // Whether to enable CXO RPC.
	CXORPCPort   *int    // CXO RPC port.
}

// Manager manages interaction with CXO and storing/retrieving node configuration files.
type Manager struct {
	mux      sync.Mutex
	c        *ManagerConfig
	l        *log2.Logger
	file     *object.CXOFileManager
	node     *node.Node
	compiler *state.Compiler
	wg       sync.WaitGroup
	newRoots chan *skyobject.Root
	quit     chan struct{}
}

// NewManager creates a new CXO manager.
func NewManager(config *ManagerConfig, compilerConfig *state.CompilerConfig) *Manager {
	manager := &Manager{
		c: config,
		l: inform.NewLogger(true, os.Stdout, LogPrefix),
		file: object.NewCXOFileManager(&object.CXOFileManagerConfig{
			Memory:   config.Memory,
			Defaults: config.Defaults,
		}),
		newRoots: make(chan *skyobject.Root, 10),
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

	// Prepare CXO file.
	if e := manager.prepareFile(); e != nil {
		manager.l.Panicln("failed to start CXO manager:", e)
	}

	// Init directories.
	if !*config.Memory {
		e := os.MkdirAll(
			path.Join(*config.Config, SubDir, ExportSubDir),
			os.FileMode(0700))
		if e != nil {
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
	c.Skyobject.Log.Prefix = "[server cxo] "

	c.Skyobject.Registry = skyobject.NewRegistry(func(t *skyobject.Reg) {
		t.Register(r0.RootPageName, r0.RootPage{})
		t.Register(r0.BoardPageName, r0.BoardPage{})
		t.Register(r0.ThreadPageName, r0.ThreadPage{})
		t.Register(r0.DiffPageName, r0.DiffPage{})
		t.Register(r0.UsersPageName, r0.UsersPage{})
		t.Register(r0.UserActivityPageName, r0.UserActivityPage{})
		t.Register(r0.BoardName, r0.Board{})
		t.Register(r0.ThreadName, r0.Thread{})
		t.Register(r0.PostName, r0.Post{})
		t.Register(r0.VoteName, r0.Vote{})
		t.Register(r0.UserName, r0.User{})
	})

	//c.MaxMessageSize = 0 // TODO -> Adjust.
	c.InMemoryDB = *m.c.Memory
	c.DataDir = filepath.Join(*m.c.Config, SubDir)
	c.EnableListener = true
	c.Listen = "[::]:" + strconv.Itoa(*m.c.CXOPort)
	c.EnableRPC = *m.c.CXORPCEnable
	c.RemoteClose = false
	c.RPCAddress = "[::]:" + strconv.Itoa(*m.c.CXORPCPort)

	c.OnRootFilled = func(c *node.Conn, root *skyobject.Root) {
		m.l.Printf("Received board '%s'", root.Hash.Hex()[:5]+"...")
		select {
		case m.newRoots <- root:
		default:
		}
	}

	// Service discovery / auto root sync.
	c.OnSubscribeRemote = func(c *node.Conn, bpk cipher.PubKey) error {
		m.l.Printf("Found board '(%s) %s'", c.Address(), bpk.Hex()[:5]+"...")
		if e := c.Node().AddFeed(bpk); e != nil {
			m.l.Println(" - Failed to relay feed:", e)
		} else {
			m.l.Println(" - Feed relayed.")
		}
		return nil
	}

	c.OnCreateConnection = func(c *node.Conn) {
		for _, pk := range m.node.Feeds() {
			c.Subscribe(pk)
		}
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

	if e := m.file.RangeConnections(func(address string) {
		if _, e := m.connectNode(address); e != nil {
			m.l.Println("prepareFile() Connection error:", e)
		}
	}); e != nil {
		return e
	}

	if e := m.file.RangeMasterSubs(func(pk cipher.PubKey, sk cipher.SecKey) {
		if r, e := m.node.Container().LastRoot(pk); e != nil {
			m.l.Println("prepareFile() LastRoot failed with error:", e)
		} else {
			m.compiler.NewRootsChan() <- r
		}
	}); e != nil {
		return e
	}

	if e := m.file.RangeRemoteSubs(func(pk cipher.PubKey) {
		if r, e := m.node.Container().LastRoot(pk); e != nil {
			m.l.Println("prepareFile() LastRoot failed with error:", e)
		} else {
			m.compiler.NewRootsChan() <- r
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
	return path.Join(*m.c.Config, SubDir, ExportSubDir, name+ExportFileExt)
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
			m.file.RangeConnections(func(address string) {
				m.connectNode(address)
			})
			m.file.Save(m.filePath())
		}
	}
}

/*
	<<< CONNECTION >>>
*/

func (m *Manager) GetConnections() []r0.Connection {
	i, out := 0, make([]r0.Connection, m.file.ConnectionsLen())
	m.file.RangeConnections(func(address string) {
		conn, state := m.node.Pool().Connection(address), ""
		if conn == nil {
			state = "DISCONNECTED"
		} else {
			state = conn.State().String()
		}
		out[i] = r0.Connection{
			Address: address,
			State:   state,
		}
		i += 1
	})
	return out
}

func (m *Manager) Connect(address string) error {

	// Add to file first.
	if e := m.file.AddConnection(address); e != nil {
		return e
	}
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
	if e := m.file.RemoveConnection(address); e != nil {
		return e
	}
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
	//if e := m.compiler.InitBoard(false, bpk, bsk); e != nil {
	//	return e
	//}
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

/*
	<<< UNSUBSCRIBE >>>
*/

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

func (m *Manager) GetBoards() ([]interface{}, []interface{}, error) {

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

	if *m.c.Master == false {
		return boo.Newf(boo.NotMaster,
			"boards can only be created on master nodes")
	}

	if e := m.file.AddMasterSub(in.BoardPubKey, in.BoardSecKey); e != nil {
		return e
	}
	m.subscribeNode(in.BoardPubKey)

	if r, e := newBoard(m.node, in); e != nil {
		return e
	} else {
		m.compiler.NewRootsChan() <- r
		return nil
	}
}

func newBoard(node *node.Node, in *object.NewBoardIO) (*skyobject.Root, error) {
	pack, e := node.Container().NewRoot(
		in.BoardPubKey,
		in.BoardSecKey,
		skyobject.HashTableIndex|skyobject.EntireTree,
		node.Container().CoreRegistry().Types(),
	)
	if e != nil {
		return nil, e
	}
	defer pack.Close()

	pack.Append(
		&r0.RootPage{
			Typ: r0.RootTypeBoard,
			Rev: 0,
			Del: false,
		},
		&r0.BoardPage{
			Board: pack.Ref(in.Board),
		},
		&r0.DiffPage{},
		&r0.UsersPage{},
	)
	if e := pack.Save(); e != nil {
		return nil, e
	}
	node.Publish(pack.Root())
	return pack.Root(), nil
}

/*
	<<< IMPORT / EXPORT >>>
*/

//func (m *Manager) ExportBoard(pk cipher.PubKey, name string) (string, *r0.ExpRoot, error) {
//	if *m.c.Memory {
//		return "", nil, nil
//	}
//	bi, e := m.compiler.GetBoard(pk)
//	if e != nil {
//		return "", nil, e
//	}
//	out, e := bi.Export()
//	if e != nil {
//		return "", nil, e
//	}
//	path := m.exportPath(name)
//	if e := file.SaveJSON(path, out, os.FileMode(0600)); e != nil {
//		return "", nil, e
//	}
//	return path, out, nil
//}

//func (m *Manager) ExportAllMasters(prefix string) error {
//	if *m.c.Memory {
//		return nil
//	}
//	return m.file.RangeMasterSubs(func(pk cipher.PubKey, _ cipher.SecKey) {
//		m.l.Printf("Exporting master '%s' ...",
//			pk.Hex()[:5]+"...")
//		if path, _, e := m.ExportBoard(pk, prefix+"."+pk.Hex()); e != nil {
//			m.l.Println(" - FAILED, Error:", e)
//		} else {
//			m.l.Printf(" - SAVED TO '%s'.", path)
//		}
//	})
//}
//
//func (m *Manager) ImportBoard(pk cipher.PubKey, name string) (string, *r0.ExpRoot, error) {
//	if *m.c.Memory {
//		return "", nil, nil
//	}
//
//	path := m.exportPath(name)
//	out := new(r0.ExpRoot)
//	if e := file.LoadJSON(path, out); e != nil {
//		return "", nil, e
//	}
//
//	_, has := m.file.GetMasterSubSecKey(pk)
//	if !has {
//		return "", nil, boo.Newf(boo.NotAuthorised,
//			"this node is not the master of board of public key '%s'", pk.Hex())
//	}
//	bi, e := m.compiler.GetBoard(pk)
//	if e != nil {
//		return "", nil, e
//	}
//	if e := bi.Import(out); e != nil {
//		return "", nil, e
//	}
//	return path, out, nil
//}
