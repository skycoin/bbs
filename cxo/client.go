package cxo

import (
	"errors"
	"github.com/evanlinjin/bbs/types"
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
	*CXOConfig
	Client       *node.Client
	UserManager  *types.UserManager
	BoardManager *types.BoardManager
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
		CXOConfig:    conf,
		Client:       client,
		UserManager:  types.NewUserManager(),
		BoardManager: types.NewBoardManager(conf.Master),
	}
	c.initUsers()
	return &c, nil
}

func (c *Client) initUsers() {
	m := c.UserManager
	if m.Load() != nil {
		m.Clear()
		m.AddNewRandomMaster()
		m.Save()
	}
}

// Launch runs the Client.
func (c *Client) Launch() error {
	if e := c.Client.Start("[::]:" + strconv.Itoa(c.Port)); e != nil {
		return e
	}
	time.Sleep(5 * time.Second)
	c.Client.Execute(func(ct *node.Container) error {
		ct.Register("Board", types.Board{},
			"Thread", types.Thread{},
			"Post", types.Post{})
		return nil
	})
	return nil
}

// Shutdown shutdowns the Client.
func (c *Client) Shutdown() error {
	return c.Client.Close()
}

// SubscribeToBoard subscribes to a board.
func (c *Client) SubscribeToBoard(pk cipher.PubKey) (*types.BoardConfig, error) {
	bc := &types.BoardConfig{
		Master:       false,
		PublicKeyStr: pk.Hex(),
	}
	// See if board exists in manager, if not create it.
	if c.BoardManager.AddConfig(bc) != nil {
		var e error
		bc, e = c.BoardManager.GetConfig(pk)
		if e != nil {
			return nil, e
		}
	}
	// Subscribe to board in cxo.
	if c.Client.Subscribe(pk) == false {
		c.BoardManager.RemoveConfig(pk)
		return nil, errors.New("cxo failed to subscribe")
	}
	return bc, nil
}

// UnSubscribeFromBoard unsubscribes from a board.
func (c *Client) UnSubscribeFromBoard(pk cipher.PubKey) bool {
	c.BoardManager.RemoveConfig(pk)
	return c.Client.Unsubscribe(pk)
}

// InjectBoard injects a board with specified BoardConfig.
func (c *Client) InjectBoard(bc *types.BoardConfig) error {
	// Check if BoardConfig is master.
	if bc.Master == false {
		return errors.New("is not master")
	}
	// Check if BoardConfig already exists in BoardManager.
	if c.BoardManager.HasConfig(bc.PublicKey) == true {
		return errors.New("config already exists")
	}
	// Create Board for cxo.
	c.Client.Execute(func(ct *node.Container) error {
		root := ct.NewRoot(bc.PublicKey, bc.SecretKey)
		board := types.NewBoard(bc.Name, bc.URL)
		root.Inject(*board, bc.SecretKey)
		return nil
	})
	// Subscribe to board in cxo.
	if c.Client.Subscribe(bc.PublicKey) == false {
		c.BoardManager.RemoveConfig(bc.PublicKey)
		return errors.New("cxo failed to subscribe")
	}
	// Add to BoardManager.
	e := c.BoardManager.AddConfig(bc)
	return e
}

// ListThreads lists the threads of specified board.
func (c *Client) ListThreads(pk cipher.PubKey) (*types.BoardView, error) {
	bc, e := c.BoardManager.GetConfig(pk)
	if e != nil {
		return nil, e
	}
	return types.NewBoardView(bc, c.Client)
}

// NewThread creates a new thread under specified board.
func (c *Client) NewThread(pk cipher.PubKey, title, desc string) (*types.Thread, error) {
	bc, e := c.BoardManager.GetConfig(pk)
	if e != nil {
		return nil, e
	}
	// Check if master.
	if bc.Master == false {
		return nil, errors.New("not master")
	}
	// Obtain secret key from config.
	sk, e := cipher.SecKeyFromHex(bc.SecretKeyStr)
	if e != nil {
		return nil, e
	}
	// Create thread and add to cxo.
	thread := types.NewThread(title, desc)
	e = c.Client.Execute(func(ct *node.Container) (_ error) {
		root := ct.Root(pk)
		root.Inject(thread, sk)
		return
	})
	if e != nil {
		return nil, e
	}
	return thread, nil
}

// ListPosts lists all posts of a specified board and thread.
func (c *Client) ListPosts(id []byte) (*types.ThreadPage, error) {
	// Find thread of id.

	e := c.Client.Execute(func(ct *node.Container) error {
		//ct.

		return nil
	})
	if e != nil {
		return nil, e
	}

	return nil, nil
}
