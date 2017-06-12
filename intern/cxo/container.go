package cxo

import (
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/cmd/bbsnode/args"
	"github.com/skycoin/bbs/intern/typ"
	"github.com/skycoin/bbs/misc"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/node/gnet"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type Container struct {
	misc.PrintMux
	c      *node.Container
	client *node.Client
	config *args.Config
	msgs   chan *Msg

	cbAddFeed     func(key cipher.PubKey) (bool, error)
	cbDelFeed     func(key cipher.PubKey) (bool, error)
	cbConnect     func(address string) error
	cbDisconnect  func(address string) error
	cbConnections func() ([]string, error)
	cbClose       func() error

	tempFiles []string
	quit      chan struct{}
}

func NewContainer(config *args.Config) (*Container, error) {
	c := &Container{
		config: config,
		msgs:   make(chan *Msg),
	}

	// Prepare waiter, error channel and timeout.
	var wg sync.WaitGroup
	eChan := make(chan error, 100)
	timeout := 20 * time.Second

	// run stuff.
	go c.setupCXOClient(wg, eChan, timeout)
	if config.CXOUseInternal() {
		go c.setupInternalCXODaemon(wg, eChan, timeout)
	} else {
		go c.setupCXORPCClient(wg, eChan, timeout)
	}

	// Wait.
	time.Sleep(time.Second)
	wg.Wait()

	// Check for errors.
	select {
	case e := <-eChan:
		if e != nil {
			return nil, e
		}
	default:
		break
	}

	log.Println("[CXOCONTAINER] Connection to cxo daemon established!")

	go c.service()
	return c, nil
}

// Helper function that sets up cxo client.
func (c *Container) setupCXOClient(wg sync.WaitGroup, eChan chan error, timeout time.Duration) {
	wg.Add(1)
	defer wg.Done()
	timer := time.NewTimer(timeout)

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
	cc := node.NewClientConfig()
	cc.MaxMessageSize = 0 /* TODO: Should really look into adjusting this in the future. */
	cc.InMemoryDB = c.config.CXOUseMemory()
	cc.DataDir = c.config.CXODir()
	cc.DBPath = filepath.Join(cc.DataDir, "client.db")

	// Setup CXO Callbacks.
	if c.config.Master() {
		cc.OnRootFilled = c.rootFilledCallBack
	}
	cc.OnAddFeed = c.feedAddedInternalCB
	cc.OnDelFeed = c.feedDeletedInternalCB

	// Change some configurations if test mode.
	if c.config.TestMode() {

		// Make temp file.
		tempDir, e := ioutil.TempDir("", "skybbs")
		if e != nil {
			eChan <- errors.Wrap(e, "unable to create temp dir")
			return
		}
		tempFile, e := ioutil.TempFile(tempDir, "")
		if e != nil {
			eChan <- errors.Wrap(e, "unable to create temp file")
			return
		}
		tempFileName := tempFile.Name()
		tempFile.Close()

		// Change stuff.
		c.tempFiles = append(c.tempFiles, tempFileName)
		cc.DataDir = tempDir
		cc.DBPath = tempFileName
		cc.InMemoryDB = false
	}

	// Attempt to set up CXO Client.
	var e error
	for {
		c.client, e = node.NewClient(cc, r)
		if e == nil {
			break
		}
		select {
		case <-timer.C:
			eChan <- errors.Wrap(e, "timeout before cxo client initiation")
			return
		default:
			time.Sleep(time.Second)
			continue
		}
	}
	addr := "[::]:" + strconv.Itoa(c.config.CXOPort())
	for {
		e = c.client.Start(addr)
		if e == nil {
			break
		}
		select {
		case <-timer.C:
			eChan <- errors.Wrap(e, "timeout before cxo client start")
			return
		default:
			time.Sleep(time.Second)
			continue
		}
	}
	for {
		if c.client.Conn().State() == gnet.ConnStateConnected {
			break
		}
		select {
		case <-timer.C:
			eChan <- errors.New("timeout before cxo client-daemon connection")
			return
		default:
			time.Sleep(time.Second)
			continue
		}
	}
	// Setup internal container.
	c.c = c.client.Container()
	return
}

// Helper function that sets up cxo rpc client.
func (c *Container) setupCXORPCClient(wg sync.WaitGroup, eChan chan error, timeout time.Duration) {
	wg.Add(1)
	defer wg.Done()
	timer := time.NewTimer(timeout)

	// Attempt to set up CXO RPC.
	addr := "[::]:" + strconv.Itoa(c.config.CXORPCPort())
	for {
		rpc, e := node.NewRPCClient(addr)
		if e == nil {
			c.cbAddFeed = rpc.AddFeed
			c.cbDelFeed = rpc.DelFeed
			c.cbClose = rpc.Close
			c.cbConnect = rpc.Connect
			c.cbDisconnect = rpc.Disconnect
			c.cbConnections = rpc.OutgoingConnections
			break
		}
		select {
		case <-timer.C:
			eChan <- errors.Wrap(e, "timeout before rpc connected")
		default:
			time.Sleep(time.Second)
			continue
		}
	}
	return
}

// Helper function that sets up cxo server (daemon).
func (c *Container) setupInternalCXODaemon(wg sync.WaitGroup, eChan chan error, timeout time.Duration) {
	wg.Add(1)
	defer wg.Done()
	timer := time.NewTimer(timeout)

	// Setup CXO Server Configuration.
	sc := node.NewServerConfig()
	sc.DataDir = c.config.CXODir()
	sc.DBPath = filepath.Join(sc.DataDir, "server.db")
	sc.InMemoryDB = c.config.CXOUseMemory()
	sc.MaxMessageSize = 0
	sc.Listen = "[::]:" + strconv.Itoa(c.config.CXOPort())
	sc.RPCAddress = "[::]:" + strconv.Itoa(c.config.CXORPCPort())
	sc.Config = gnet.NewConfig()

	// Change some configurations if test mode.
	if c.config.TestMode() {

		// Make temp file.
		tempDir, e := ioutil.TempDir("", "skybbs")
		if e != nil {
			eChan <- errors.Wrap(e, "unable to create temp dir")
			return
		}
		tempFile, e := ioutil.TempFile(tempDir, "")
		if e != nil {
			eChan <- errors.Wrap(e, "unable to create temp file")
			return
		}
		tempFileName := tempFile.Name()
		tempFile.Close()

		// Change stuff.
		c.tempFiles = append(c.tempFiles, tempFileName)
		sc.DataDir = tempDir
		sc.DBPath = tempFileName
		sc.InMemoryDB = false
	}

	// Attempt to run server.
	var cxoServer *node.Server
	var e error
	for {
		cxoServer, e = node.NewServer(sc)
		if e != nil {
			log.Println("[CONTAINER] New CXO Server Error:", e)
		} else {
			c.cbAddFeed = func(key cipher.PubKey) (bool, error) {
				return cxoServer.AddFeed(key), nil
			}
			c.cbDelFeed = func(key cipher.PubKey) (bool, error) {
				return cxoServer.DelFeed(key), nil
			}
			c.cbClose = cxoServer.Close
			c.cbConnect = cxoServer.Connect
			c.cbDisconnect = cxoServer.Disconnect
			c.cbConnections = func() ([]string, error) {
				connections := cxoServer.Connections()
				addresses := []string{}
				for _, conn := range connections {
					if !conn.IsIncoming() {
						addresses = append(addresses, conn.Address())
					}
				}
				return addresses, nil
			}
			break
		}
		select {
		case <-timer.C:
			eChan <- errors.Wrap(e, "timeout before cxo server created")
		default:
			time.Sleep(time.Second)
			continue
		}
	}
	for {
		e := cxoServer.Start()
		if e == nil {
			break
		}
		select {
		case <-timer.C:
			eChan <- errors.Wrap(e, "timeout before cxo server started")
		default:
			time.Sleep(time.Second)
			continue
		}
	}
	return
}

// Close closes the BBS container.
func (c *Container) Close() error {
	for {
		select {
		case c.quit <- struct{}{}:
		default:
			goto ServiceFinished
		}
	}
ServiceFinished:
	for _, fName := range c.tempFiles {
		os.Remove(fName)
	}
	defer c.cbClose()
	return c.client.Close()
}

// Helper function for running garbage collector service.
func (c *Container) service() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-c.quit:
			return
		case <-ticker.C:
			c.Lock(c.service)
			c.c.GC(false)
			c.Unlock()
		}
	}
}

// Feeds returns a list of all feeds we are subscribed to.
func (c *Container) Feeds() []cipher.PubKey {
	c.Lock(c.Feeds)
	defer c.Unlock()
	return c.client.Feeds()
}

// Subscribe subscribes to a cxo feed.
func (c *Container) Subscribe(pk cipher.PubKey) (bool, error) {
	c.Lock(c.Subscribe)
	defer c.Unlock()
	if _, e := c.cbAddFeed(pk); e != nil {
		return false, e
	}
	return c.client.Subscribe(pk), nil
}

// Unsubscribe unsubscribes from a cxo feed.
func (c *Container) Unsubscribe(pk cipher.PubKey) (bool, error) {
	c.Lock(c.Unsubscribe)
	defer c.Unlock()
	if _, e := c.cbDelFeed(pk); e != nil {
		return false, e
	}
	return c.client.Unsubscribe(pk, false), nil
}

// Connect connects to external BBS Node.
func (c *Container) Connect(address string) error {
	c.Lock(c.Connect)
	defer c.Unlock()
	return c.cbConnect(address)
}

// Disconnect disconnects from external BBS Node.
func (c *Container) Disconnect(address string) error {
	c.Lock(c.Disconnect)
	defer c.Unlock()
	return c.cbDisconnect(address)
}

// GetConnections gets a list of all external BBS Node addresses we are connected to.
func (c *Container) GetConnections() ([]string, error) {
	c.Lock(c.GetConnections)
	defer c.Unlock()
	return c.cbConnections()
}
