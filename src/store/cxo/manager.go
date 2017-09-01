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
	SubDir        = "cxo"
	DBName        = "bbs.db"
	FileName      = "bbs.json"
	RetryDuration = time.Second * 5
)

// ManagerConfig represents the configuration for CXO Manager.
type ManagerConfig struct {
	Memory       *bool   // Whether to enable memory mode.
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
	file     *object.CXOFile
	node     *node.Node
	compiler *state.Compiler
	wg       sync.WaitGroup
	newRoots chan *skyobject.Root
	quit     chan struct{}
}

// NewManager creates a new CXO manager.
func NewManager(config *ManagerConfig, compilerConfig *state.CompilerConfig) *Manager {
	manager := &Manager{
		c:        config,
		l:        inform.NewLogger(true, os.Stdout, LogPrefix),
		file:     new(object.CXOFile),
		newRoots: make(chan *skyobject.Root, 10),
		quit:     make(chan struct{}),
	}
	if e := manager.setup(); e != nil {
		manager.l.Panicln("failed to start CXO manager:", e)
	}
	manager.compiler = state.NewCompiler(
		compilerConfig, manager.file, manager.newRoots, manager.node,
		views.AddContent(),
		views.AddFollow(),
	)

	manager.file.Lock()
	masterSubs := manager.file.MasterSubs
	remoteSubs := manager.file.RemoteSubs
	manager.file.Unlock()

	for _, sub := range masterSubs {
		if e := manager.compiler.InitBoard(false, sub.PK, sub.SK); e != nil {
			manager.l.Println("compiler.InitBoard() :", e)
		}
	}
	for _, sub := range remoteSubs {
		if e := manager.compiler.InitBoard(false, sub.PK); e != nil {
			manager.l.Println("compiler.InitBoard() :", e)
		}
	}
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

// Sets up the CXO manager.
func (m *Manager) setup() error {
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
	c.DBPath = filepath.Join(c.DataDir, DBName)
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
		go func() {
			m.wg.Add(1)
			defer m.wg.Done()

			for _, pk := range m.node.Feeds() {
				c.Subscribe(pk)
			}
		}()
	}
	var e error
	if m.node, e = node.NewNode(c); e != nil {
		return e
	}
	if e := m.load(); e != nil {
		return e
	}
	m.init()
	return nil
}

func (m *Manager) filePath() string {
	return path.Join(*m.c.Config, SubDir, FileName)
}

func (m *Manager) load() error {
	if *m.c.Memory {
		return nil
	}
	if e := file.LoadJSON(m.filePath(), m.file); e != nil {
		if !os.IsNotExist(e) {
			return boo.WrapType(e, boo.InvalidRead,
				"failed to read CXO file")
		}
	}
	return nil
}

func (m *Manager) save() error {
	if *m.c.Memory {
		return nil
	}
	if e := file.SaveJSON(m.filePath(), m.file, os.FileMode(0600)); e != nil {
		return boo.WrapType(e, boo.Internal,
			"failed to save CXO file")
	}
	return nil
}

func (m *Manager) init() {
	m.file.Lock()
	defer m.file.Unlock()

	for _, sub := range m.file.MasterSubs {
		m.subscribeNode(sub.PK)
	}

	for _, address := range m.file.Connections {
		m.connectNode(address)
	}

	for _, sub := range m.file.RemoteSubs {
		m.subscribeNode(sub.PK)
	}
}

/*
	<<< CONNECT >>>
*/

func (m *Manager) GetConnections() []r0.Connection {
	m.file.Lock()
	defer m.file.Unlock()

	out := make([]r0.Connection, len(m.file.Connections))
	for i, address := range m.file.Connections {
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
	}
	return out
}

func (m *Manager) Connect(address string) error {
	if e := m.connectFile(address); e != nil {
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

func (m *Manager) connectFile(address string) error {
	if !m.file.AddConnection(address) {
		return boo.Newf(boo.AlreadyExists,
			"connection to address '%s' already exists", address)
	}
	return m.save()
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

/*
	<<< DISCONNECT >>>
*/

func (m *Manager) Disconnect(address string) error {
	if e := m.disconnectFile(address); e != nil {
		return e
	}
	if e := m.disconnectNode(address); e != nil {
		return e
	}
	return nil
}

func (m *Manager) disconnectFile(address string) error {
	if !m.file.RemoveConnection(address) {
		return boo.Newf(boo.NotFound,
			"connection to address '%s' is not recorded in cxo file", address)
	}
	return m.save()
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
	<<< SUBSCRIBE >>>
*/

func (m *Manager) GetSubscriptions() []cipher.PubKey {
	m.file.Lock()
	defer m.file.Unlock()
	subs := append(m.file.MasterSubs, m.file.RemoteSubs...)
	out := make([]cipher.PubKey, len(subs))
	for i, sub := range subs {
		out[i] = sub.PK
	}
	return out
}

func (m *Manager) SubscribeRemote(bpk cipher.PubKey) error {
	if e := m.subscribeFileRemote(bpk); e != nil {
		return e
	}
	m.subscribeNode(bpk)
	return nil
}

func (m *Manager) SubscribeMaster(bpk cipher.PubKey, bsk cipher.SecKey) error {
	if e := m.subscribeFileMaster(bpk, bsk); e != nil {
		return e
	}
	m.subscribeNode(bpk)
	if e := m.compiler.InitBoard(false, bpk, bsk); e != nil {
		return e
	}
	return nil
}

func (m *Manager) subscribeFileRemote(bpk cipher.PubKey) error {
	if !m.file.AddRemoteSub(bpk) {
		return boo.Newf(boo.AlreadyExists,
			"already subscribed to remote board '%s'", bpk.Hex()[:5]+"...")
	}
	return m.save()
}

func (m *Manager) subscribeFileMaster(bpk cipher.PubKey, bsk cipher.SecKey) error {
	if !m.file.AddMasterSub(bpk, bsk) {
		return boo.Newf(boo.AlreadyExists,
			"already subscribed to master board '%s'", bpk.Hex()[:5]+"...")
	}
	return m.save()
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
	if e := m.unsubscribeFileRemote(bpk); e != nil {
		return e
	}
	m.unsubscribeNode(bpk)
	return nil
}

func (m *Manager) UnsubscribeMaster(bpk cipher.PubKey) error {
	if e := m.unsubscribeFileMaster(bpk); e != nil {
		return e
	}
	m.unsubscribeNode(bpk)
	return nil
}

func (m *Manager) unsubscribeFileRemote(bpk cipher.PubKey) error {
	if !m.file.RemoveRemoteSub(bpk) {
		return boo.Newf(boo.NotFound,
			"remote board of public key '%s' not found in cxo file", bpk.Hex()[:5]+"...")
	}
	return m.save()
}

func (m *Manager) unsubscribeFileMaster(bpk cipher.PubKey) error {
	if !m.file.RemoveMasterSub(bpk) {
		return boo.Newf(boo.NotFound,
			"master board of public key '%s' not found in cxo file", bpk.Hex()[:5]+"...")
	}
	return m.save()
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
	m.file.Lock()
	defer m.file.Unlock()

	masterOut := make([]interface{}, len(m.file.MasterSubs))
	for i, sub := range m.file.MasterSubs {
		bi, e := m.compiler.GetBoard(sub.PK)
		if e != nil {
			return nil, nil, e
		}
		bView, e := bi.Get(views.Content, content_view.Board)
		if e != nil {
			return nil, nil, e
		}
		masterOut[i] = bView
	}

	remoteOut := make([]interface{}, len(m.file.RemoteSubs))
	for i, sub := range m.file.RemoteSubs {
		bi, e := m.compiler.GetBoard(sub.PK)
		if e != nil {
			m.l.Println(e)
			continue
		}
		bView, e := bi.Get(views.Content, content_view.Board)
		if e != nil {
			m.l.Println(e)
			continue
		}
		remoteOut[i] = bView
	}

	return masterOut, remoteOut, nil
}

func (m *Manager) NewBoard(in *object.NewBoardIO) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	if e := m.subscribeFileMaster(in.BoardPubKey, in.BoardSecKey); e != nil {
		return e
	}
	m.subscribeNode(in.BoardPubKey)

	if e := newBoard(m.node, in); e != nil {
		return e
	}
	if e := m.compiler.InitBoard(false, in.BoardPubKey, in.BoardSecKey); e != nil {
		return e
	}
	return nil
}

func (m *Manager) NewThread(thread *r0.Thread) (uint64, error) {
	bi, e := m.compiler.GetBoard(thread.OfBoard)
	if e != nil {
		return 0, e
	}
	return bi.NewThread(thread)
}

func newBoard(node *node.Node, in *object.NewBoardIO) error {
	pack, e := node.Container().NewRoot(
		in.BoardPubKey,
		in.BoardSecKey,
		skyobject.HashTableIndex|skyobject.EntireTree,
		node.Container().CoreRegistry().Types(),
	)
	if e != nil {
		return e
	}
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
		return e
	}
	node.Publish(pack.Root())
	return nil
}
