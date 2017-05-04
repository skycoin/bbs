package cxo

import (
	"errors"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/skycoin/src/cipher"
	"strconv"
	"time"
	"github.com/evanlinjin/bbs/types"
)

// CXOConfig represents a configuration for CXO Client.
type CXOConfig struct {
	Master          bool
	Port            int
	BoardConfigFile string
	ConfigDir       string
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
		BoardManager: types.NewBoardManager(conf.Master),
	}
	return &c, nil
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
		Master:    false,
		PublicKey: pk.Hex(),
	}
	// See if board exists in manager, if not create it.
	if c.BoardManager.AddConfig(bc) != nil {
		var e error
		bc, e = c.BoardManager.GetConfig(pk)
		if e != nil {
			return nil, e
		}
	}
	c.Client.Subscribe(pk)
	return bc, nil
}

// UnSubscribeFromBoard unsubscribes from a board.
func (c *Client) UnSubscribeFromBoard(pk cipher.PubKey) bool {
	c.BoardManager.RemoveConfig(pk)
	return c.Client.Unsubscribe(pk)
}

// ListBoards lists all boards we are subscribed to.
func (c *Client) ListBoards() []*types.BoardConfig {
	return c.BoardManager.GetList()
}

// NewBoard creates a new board with a seed.
func (c *Client) NewBoard(name, seed string) (*types.BoardConfig, error) {
	bc, pk, sk, e := c.BoardManager.NewMasterConfigFromSeed(seed, "")
	if e != nil {
		return nil, e
	}
	// Create Board for cxo.
	e = c.Client.Execute(func(ct *node.Container) error {
		root := ct.NewRoot(pk, sk)
		board := types.NewBoard(name, "")
		root.Inject(*board, sk)
		return nil
	})
	c.Client.Subscribe(pk)
	return bc, e
}

// ListThreads lists the threads of specified board.
func (c *Client) ListThreads(pk cipher.PubKey) (*types.BoardPage, error) {
	bc, e := c.BoardManager.GetConfig(pk)
	if e != nil {
		return nil, e
	}
	return bc.NewBoardPage(pk, c.Client)
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
	sk, e := cipher.SecKeyFromHex(bc.SecretKey)
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
