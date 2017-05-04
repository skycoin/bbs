package cxo

import (
	"errors"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/skycoin/src/cipher"
	"time"
)

// CXOConfig represents a configuration for Client.
type CXOConfig struct {
	Master          bool
	Address         string
	BoardConfigFile string
}

// NewCXOConfig creates a new CXOConfig.
func NewCXOConfig() *CXOConfig {
	c := CXOConfig{
		Master:  false,
		Address: "[::]:8998",
	}
	return &c
}

// Client contains all the boards, threads and posts.
type Client struct {
	*CXOConfig
	Client *node.Client

	// TODO: Should create identity manager.
	CurrentIdentity cipher.PubKey
	Identities      map[cipher.PubKey]*UserConfig

	BoardManager *BoardManager
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
		Identities:   make(map[cipher.PubKey]*UserConfig),
		BoardManager: NewBoardManager(),
	}
	return &c, nil
}

// Launch runs the Client.
func (c *Client) Launch() error {
	if e := c.Client.Start(c.Address); e != nil {
		return e
	}
	time.Sleep(5 * time.Second)
	c.InitRegister()
	return nil
}

// Shutdown shutdowns the Client.
func (c *Client) Shutdown() error {
	return c.Client.Close()
}

// InitRegister initiates schema registration.
func (c *Client) InitRegister() {
	c.Client.Execute(func(ct *node.Container) (_ error) {
		ct.Register("Board", Board{},
			"Thread", Thread{},
			"Post", Post{})
		return
	})
}

// AddIdentity Adds an identity.
func (c *Client) AddIdentity(uc *UserConfig) error {
	if uc == nil {
		return errors.New("nil UserConfig")
	}
	if e := uc.CheckKeys(); e != nil {
		return e
	}
	c.Identities[uc.PublicKey] = uc

	// If first identity, set as current.
	// TODO: Improve this.
	if len(c.Identities) == 1 {
		c.SetCurrentIdentity(uc.PublicKey)
	}
	return nil
}

// AddRandomIdentity adds a random identity.
func (c *Client) AddRandomIdentity() error {
	pk, sk := cipher.GenerateKeyPair()
	uc := &UserConfig{
		PublicKey: pk,
		SecretKey: sk,
	}
	return c.AddIdentity(uc)
}

// SetCurrentIdentity sets the current identity.
func (c *Client) SetCurrentIdentity(pk cipher.PubKey) error {
	if e := pk.Verify(); e != nil {
		return errors.New("invalid public key")
	}
	if _, has := c.Identities[pk]; has == false {
		return errors.New("identity unavaliable")
	}
	c.CurrentIdentity = pk
	return nil
}

// CheckIdentity checks whether current identity is set.
func (c *Client) CheckCurrentIdentity() bool {
	_, has := c.Identities[c.CurrentIdentity]
	return has
}

// SubscribeToBoard subscribes to a board.
func (c *Client) SubscribeToBoard(pk cipher.PubKey) (*BoardConfig, error) {
	bc := &BoardConfig{
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
func (c *Client) ListBoards() []*BoardConfig {
	return c.BoardManager.GetList()
}

// NewBoard creates a new board with a seed.
func (c *Client) NewBoard(name, seed string) (*BoardConfig, error) {
	bc, pk, sk, e := c.BoardManager.NewMasterConfigFromSeed(seed, "")
	if e != nil {
		return nil, e
	}
	// Create Board for cxo.
	e = c.Client.Execute(func(ct *node.Container) error {
		root := ct.NewRoot(pk, sk)
		board := NewBoard(name, "")
		root.Inject(*board, sk)
		return nil
	})
	c.Client.Subscribe(pk)
	return bc, e
}

// ListThreads lists the threads of specified board.
func (c *Client) ListThreads(pk cipher.PubKey) (*BoardPage, error) {
	bc, e := c.BoardManager.GetConfig(pk)
	if e != nil {
		return nil, e
	}
	return bc.NewBoardPage(pk, c.Client)
}

// NewThread creates a new thread under specified board.
func (c *Client) NewThread(pk cipher.PubKey, title, desc string) (*Thread, error) {
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
	thread := NewThread(title, desc)
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
func (c *Client) ListPosts(id []byte) (*ThreadPage, error) {
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
