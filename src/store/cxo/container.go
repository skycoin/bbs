package cxo

import (
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/cmd/bbsnode/args"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/bbs/src/store/typ"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/node/gnet"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Container struct {
	misc.PrintMux
	c    *node.Container
	node *node.Node

	config *args.Config
	msgs   chan *Msg

	tempFiles []string
}

func NewContainer(config *args.Config) (*Container, error) {
	c := &Container{
		config: config,
		msgs:   make(chan *Msg),
	}
	// Setup stuff.
	if e := c.setupCXONode(); e != nil {
		return nil, e
	}
	return c, nil
}

// Helper function that sets up cxo client.
func (c *Container) setupCXONode() error {

	// Setup CXO Registry.
	r := skyobject.NewRegistry()
	r.Register("Board", typ.Board{})
	r.Register("Thread", typ.Thread{})
	r.Register("Post", typ.Post{})
	r.Register("ThreadPage", typ.ThreadPage{})
	r.Register("BoardContainer", typ.BoardContainer{})

	r.Register("Vote", typ.Vote{})
	r.Register("ThreadVotes", typ.ThreadVotes{})
	r.Register("PostVotes", typ.PostVotes{})
	r.Register("ThreadVotesContainer", typ.ThreadVotesContainer{})
	r.Register("PostVotesContainer", typ.PostVotesContainer{})
	r.Done()

	// Setup CXO Configuration.
	nc := node.NewConfig()
	nc.MaxMessageSize = 0 /* TODO: Should really look into adjusting this in the future. */
	nc.InMemoryDB = c.config.CXOUseMemory()
	nc.DataDir = c.config.CXODir()
	nc.DBPath = filepath.Join(nc.DataDir, "node.db")

	nc.EnableListener = true
	nc.Listen = "[::]:" + strconv.Itoa(c.config.CXOPort())
	nc.RPCAddress = "[::]:" + strconv.Itoa(c.config.CXORPCPort())
	nc.GCInterval = 5 * time.Second

	// Setup CXO Callbacks.
	if c.config.Master() {
		nc.OnRootFilled = c.rootFilledInternalCB
		nc.PublicServer = true
	}
	nc.OnSubscriptionAccepted = c.subAcceptedInternalCB
	nc.OnSubscriptionRejected = c.subRejectedInternalCB
	nc.OnCreateConnection = c.connCreatedInternalCB
	nc.OnCloseConnection = c.connClosedInternalCB

	// Change some configurations if test mode.
	if c.config.TestMode() {

		// Make temp file.
		tempDir, e := ioutil.TempDir("", "skybbs")
		if e != nil {
			return errors.Wrap(e, "unable to create temp dir")
		}
		tempFile, e := ioutil.TempFile(tempDir, "")
		if e != nil {
			return errors.Wrap(e, "unable to create temp file")
		}
		tempFileName := tempFile.Name()
		tempFile.Close()

		// Change stuff.
		c.tempFiles = append(c.tempFiles, tempFileName)
		nc.DataDir = tempDir
		nc.DBPath = tempFileName
		nc.InMemoryDB = false
	}

	// Attempt to set up CXO Node.
	var e error
	if c.node, e = node.NewNodeReg(nc, r); e != nil {
		return e
	}
	// Setup internal container.
	c.c = c.node.Container()
	return nil
}

// Close closes the BBS container.
func (c *Container) Close() error {
	for _, fName := range c.tempFiles {
		os.Remove(fName)
	}
	if e := c.node.Close(); e != nil {
		return e
	}
	select {
	case <-c.node.Quiting():
	}
	return nil
}

// Lock for thread safety.
func (c *Container) Lock(function interface{}) {
	c.PrintMux.Lock(function)
	c.c.LockGC()
}

// Unlock for thread safety.
func (c *Container) Unlock() {
	c.PrintMux.Unlock()
	c.c.UnlockGC()
}

// Feeds returns a list of all feeds we are subscribed to.
func (c *Container) Feeds() []cipher.PubKey {
	c.Lock(c.Feeds)
	defer c.Unlock()
	return c.node.Feeds()
}

// Gets the address of the node.
func (c *Container) GetAddress() string {
	return c.node.Pool().Address()
}

// Subscribe subscribes to a cxo feed.
func (c *Container) Subscribe(addr string, pk cipher.PubKey) error {
	c.Lock(c.Subscribe)
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
func (c *Container) Unsubscribe(addr string, pk cipher.PubKey) error {
	c.Lock(c.Unsubscribe)
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
func (c *Container) GetConnections() []string {
	c.Lock(c.GetConnections)
	defer c.Unlock()
	conns := c.node.Pool().Connections()
	addresses := make([]string, len(conns))
	for i, conn := range conns {
		addresses[i] = conn.Address()
	}
	return addresses
}

// Connect adds a connection.
func (c *Container) Connect(addr string) error {
	c.Lock(c.Connect)
	defer c.Unlock()
	if _, e := c.node.Pool().Dial(addr); e != nil {
		return e
	}
	return nil
}

// Disconnect removes a connection.
func (c *Container) Disconnect(addr string) error {
	c.Lock(c.Disconnect)
	defer c.Unlock()
	conn := c.node.Pool().Connection(addr)
	if conn == nil {
		return nil
	}
	return conn.Close()
}
