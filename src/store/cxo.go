package store

import (
	"github.com/skycoin/bbs/src/boo"
	"github.com/skycoin/bbs/src/store/obj"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/node/gnet"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const (
	cxoConfigSubDir = "cxo"     // sub directory of configuration folder to store cxo db.
	cxoDBFileName   = "node.db" // file name of cxo db file.
)

// CXOConfig provides configuration options for CXO.
type CXOConfig struct {
	Master       *bool   // Whether node is master.
	TestMode     *bool   // Whether node is in test mode.
	MemoryMode   *bool   // Whether to use local storage in runtime.
	ConfigDir    *string // Configuration directory.
	CXOPort      *int    // CXO listening port.
	CXORPCEnable *bool   // Whether to enable CXO RPC.
	CXORPCPort   *int    // CXO RPC listening port.
}

// CXO updates events from cxo.
type CXO struct {
	sync.Mutex
	config *CXOConfig
	node   *node.Node
}

// NewCXO creates a new CXO.
func NewCXO(config *CXOConfig, updateFunc func(root *node.Root)) (*CXO, error) {
	updater := &CXO{
		config: config,
	}
	if e := updater.setupCXO(updateFunc); e != nil {
		return nil, e
	}
	return updater, nil
}

// Sets up CXO node.
func (c *CXO) setupCXO(updateFunc func(root *node.Root)) error {
	// Setup registry.
	r := skyobject.NewRegistry()
	r.Register("Board", obj.Board{})
	r.Register("Thread", obj.Thread{})
	r.Register("Post", obj.Post{})
	r.Register("ThreadPage", obj.ThreadPage{})
	r.Register("BoardPage", obj.BoardPage{})
	r.Register("ExternalRoot", obj.ExternalRoot{})
	r.Done()

	// Setup CXO Configurations.
	nc := node.NewConfig()
	nc.MaxMessageSize = 0 /* TODO: Adjust. */
	nc.InMemoryDB = *c.config.MemoryMode
	nc.DataDir = filepath.Join(*c.config.ConfigDir, cxoConfigSubDir)
	nc.DBPath = filepath.Join(nc.DataDir, cxoDBFileName)

	nc.EnableListener = true
	nc.Listen = "[::]:" + strconv.Itoa(*c.config.CXOPort)
	nc.EnableRPC = *c.config.CXORPCEnable
	nc.RPCAddress = "[::]:" + strconv.Itoa(*c.config.CXORPCPort)
	nc.GCInterval = 5 * time.Second

	// Setup CXO Callbacks.
	nc.OnRootFilled = updateFunc

	// Attempt to setup CXO Node.
	var e error
	c.node, e = node.NewNodeReg(nc, r)
	return e
}

// Close closes the BBS container.
func (c *CXO) Close() error {
	if e := c.node.Close(); e != nil {
		return e
	}
	select {
	case <-c.node.Quiting():
	}
	return nil
}

func (c *CXO) GetRoot(pk cipher.PubKey) *node.Root {
	return c.node.Container().LastFullRoot(pk)
}

// NewRoot creates a new root.
func (c *CXO) NewRoot(seed []byte, modifier RootModifier) (cipher.PubKey, cipher.SecKey, error) {
	pk, sk := cipher.GenerateDeterministicKeyPair(seed)
	c.Lock()
	defer c.Unlock()
	root, e := c.node.Container().NewRoot(pk, sk)
	if e != nil {
		return pk, sk, boo.New(boo.Internal, "failed to create root:", e.Error())
	}
	return pk, sk, modifier(root)
}

// RootModifier modifies the root.
type RootModifier func(r *node.Root) error

// GetMasterWalker obtains walker of root that node can edit.
func (c *CXO) ModifyRoot(pk cipher.PubKey, sk cipher.SecKey, modifier RootModifier) error {
	c.Lock()
	defer c.Unlock()
	return modifier(c.node.Container().LastRootSk(pk, sk))
}

// IsMaster returns whether master or not.
func (c *CXO) IsMaster() bool {
	return *c.config.Master
}

// Feeds returns a list of all feeds we are subscribed to.
func (c *CXO) Feeds() []cipher.PubKey {
	c.Lock()
	defer c.Unlock()
	return c.node.Feeds()
}

// GetAddress gets the address of the node.
func (c *CXO) GetAddress() string {
	return c.node.Pool().Address()
}

// Subscribe subscribes to a cxo feed.
func (c *CXO) Subscribe(addr string, pk cipher.PubKey) error {
	c.Lock()
	defer c.Unlock()
	if addr == "" {
		c.node.Subscribe(nil, pk)
		return nil
	}
	conn, e := c.node.Pool().Dial(addr)
	if e != nil {
		switch e {
		case gnet.ErrAlreadyListen:
		default:
			return e
		}
	}
	return c.node.SubscribeResponse(conn, pk)
}

// Unsubscribe unsubscribes from a cxo feed.
func (c *CXO) Unsubscribe(addr string, pk cipher.PubKey) error {
	c.Lock()
	defer c.Unlock()
	if addr == "" {
		c.node.Unsubscribe(nil, pk)
		return nil
	}
	conn := c.node.Pool().Connection(addr)
	if conn == nil {
		c.node.Unsubscribe(nil, pk)
		return nil
	}
	c.node.Unsubscribe(conn, pk)
	return nil
}

// GetConnections gets connections.
func (c *CXO) GetConnections() []string {
	c.Lock()
	defer c.Unlock()
	conns := c.node.Pool().Connections()
	addresses := make([]string, len(conns))
	for i, conn := range conns {
		addresses[i] = conn.Address()
	}
	return addresses
}

// Connect adds a connection.
func (c *CXO) Connect(addr string) error {
	c.Lock()
	defer c.Unlock()
	if _, e := c.node.Pool().Dial(addr); e != nil {
		return e
	}
	return nil
}

// Disconnect removes a connection.
func (c *CXO) Disconnect(addr string) error {
	c.Lock()
	defer c.Unlock()
	conn := c.node.Pool().Connection(addr)
	if conn == nil {
		return nil
	}
	return conn.Close()
}
