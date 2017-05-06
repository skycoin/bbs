package cxo

import (
	"github.com/evanlinjin/bbs/typ"
	"github.com/skycoin/cxo/node"
	"strconv"
	"time"
)

// CXOConfig represents a configuration for CXO Client.
type CXOConfig struct {
	Master bool
	Port   int
}

// NewCXOConfig creates a new CXOConfig.
func NewCXOConfig() *CXOConfig {
	c := CXOConfig{
		Master: false,
		Port:   8998,
	}
	return &c
}

// Client acts as a client for cxo.
type Client struct {
	config   *CXOConfig
	cxo      *node.Client
	uManager *typ.UserManager
	bManager *typ.BoardManager
}

// NewClient creates a new Client.
func NewClient(conf *CXOConfig) (*Client, error) {

	// Setup cxo client.
	clientConfig := node.NewClientConfig()
	client, e := node.NewClient(clientConfig)
	if e != nil {
		return nil, e
	}
	c := Client{
		config:   conf,
		cxo:      client,
		uManager: typ.NewUserManager(),
		bManager: typ.NewBoardManager(conf.Master),
	}
	c.initUsers()
	c.initBoards()
	return &c, nil
}

func (c *Client) initUsers() {
	m := c.uManager
	if m.Load() != nil {
		m.Clear()
		m.AddNewRandomMaster()
		if e := m.Save(); e != nil {
			panic(e)
		}
	}
}

func (c *Client) initBoards() {
	m := c.bManager
	if m.Load() != nil {
		m.Clear()
		if e := m.Save(); e != nil {
			panic(e)
		}
	}
}

// Launch runs the Client.
func (c *Client) Launch() error {
	// Connect to cxo daemon.
	if e := c.cxo.Start("[::]:" + strconv.Itoa(c.config.Port)); e != nil {
		return e
	}
	// Wait for safety.
	time.Sleep(5 * time.Second)
	// Register schemas for cxo.
	e := c.cxo.Execute(func(ct *node.Container) error {
		ct.Register("Board", typ.Board{})
		ct.Register("Thread", typ.Thread{})
		ct.Register("Post", typ.Post{})
		ct.Register("ThreadPage", typ.ThreadPage{})
		ct.Register("BoardContainer", typ.BoardContainer{})
		return nil
	})
	if e != nil {
		return e
	}
	// Load boards from BoardManager.
	for _, bc := range c.bManager.Boards {
		if bc.Master && c.config.Master {
			//c.InjectBoard(bc)
		} else {
			//c.SubscribeToBoard(bc.PublicKey)
		}
	}
	return e
}

// Shutdown shutdowns the Client.
func (c *Client) Shutdown() error {
	return c.cxo.Close()
}
