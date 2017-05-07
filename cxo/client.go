package cxo

import (
	"errors"
	"fmt"
	"github.com/evanlinjin/bbs/typ"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/skycoin/src/cipher"
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
		_, e := c.Subscribe(bc.PubKey)
		if e != nil {
			fmt.Println("[BOARD CONFIG] Subscription failed:", bc.PubKey, e)
		}
	}
	return e
}

// Shutdown shutdowns the Client.
func (c *Client) Shutdown() error {
	return c.cxo.Close()
}

// Subscribe subscribes to a board.
func (c *Client) Subscribe(pk cipher.PubKey) (*typ.BoardConfig, error) {
	// Check if we already have the config, if not make one.
	bc, e := c.bManager.GetConfig(pk)
	if e != nil {
		bc, _ = typ.NewBoardConfig(pk, "")
	}
	// Attempt subscribe in cxo.
	if c.cxo.Subscribe(pk) == false {
		return nil, errors.New("subscription failed")
	}
	// Add to BoardManager.
	c.bManager.AddConfig(bc)
	return bc, nil
}

// Unsubscribe unsubscribes from a board.
func (c *Client) Unsubscribe(pk cipher.PubKey) bool {
	c.bManager.RemoveConfig(pk)
	return c.cxo.Unsubscribe(pk)
}

// InjectBoard injects a new master board with given seed.
func (c *Client) InjectBoard(board *typ.Board, seed string) error {
	if c.config.Master == false {
		return errors.New("action not allowed; node is not master")
	}
	// Make BoardConfig.
	bc := typ.NewMasterBoardConfig(board, seed)
	// Check if BoardConfig exists.
	if c.bManager.HasConfig(bc.PubKey) == true {
		return errors.New("board already exists")
	}
	// Add to cxo.
	e := c.cxo.Execute(func(ct *node.Container) error {
		r := ct.NewRoot(bc.PubKey, bc.SecKey)
		boardRef := ct.Save(*board)
		boardContainer := typ.BoardContainer{Board: boardRef}
		r.Inject(boardContainer, bc.SecKey)
		return nil
	})
	if e != nil {
		return e
	}
	// Add config to BoardManager.
	return c.bManager.AddConfig(bc)
}

// ObtainBoard obtains the latest board from cxo.
//func (c *Client) ObtainBoard