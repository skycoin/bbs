package cxo

import (
	"errors"
	"fmt"
	"github.com/evanlinjin/bbs/typ"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
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
			fmt.Println("[BOARD CONFIG]", bc.PubKey.Hex(), e)
			fmt.Println("[BOARD CONFIG] Removing", bc.PubKey.Hex())
			c.bManager.RemoveConfig(bc.PubKey)
		}
		//if bc.Master {
		//	if c.config.Master == false {
		//		continue
		//	}
		//	e := c.InjectBoard(bc)
		//	if e != nil {
		//		fmt.Println("[BOARD CONFIG]", bc.PubKey.Hex(), e)
		//		fmt.Println("[BOARD CONFIG] Removing", bc.PubKey.Hex())
		//		c.bManager.RemoveConfig(bc.PubKey)
		//	}
		//} else {
		//	_, e := c.Subscribe(bc.PubKey)
		//	if e != nil {
		//		fmt.Println("[BOARD CONFIG]", bc.PubKey.Hex(), e)
		//		fmt.Println("[BOARD CONFIG] Removing", bc.PubKey.Hex())
		//		c.bManager.RemoveConfig(bc.PubKey)
		//	}
		//}
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

// InjectBoard injects a new master board with given BoardConfig.
func (c *Client) InjectBoard(bc *typ.BoardConfig) error {
	if c.config.Master == false {
		return errors.New("action not allowed; node is not master")
	}
	// Make Board from BoardConfig.
	board := typ.NewBoardFromConfig(bc)
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
	// Subscribe to board.
	if c.cxo.Subscribe(bc.PubKey) == false {
		return errors.New("failed to subscribe to board in cxo")
	}
	// Add config to BoardManager.
	return c.bManager.AddConfig(bc)
}

// ObtainBoard obtains a board of specified PubKey from cxo.
// Also obtains the BoardContainer.
// Assumes public key provided is valid.
func (c *Client) ObtainBoard(bpk cipher.PubKey) (
	b *typ.Board, bCont *typ.BoardContainer, e error,
) {
	bCont = new(typ.BoardContainer)
	b = new(typ.Board)

	e = c.cxo.Execute(func(ct *node.Container) error {
		r := ct.Root(bpk)
		if r == nil {
			c.bManager.RemoveConfig(bpk)
			return errors.New("nil root, removed from config")
		}
		values, e := r.Values()
		if e != nil {
			return e
		}
		if len(values) != 1 {
			return errors.New("invalid root")
		}
		// Get BoardContainer.
		if e := encoder.DeserializeRaw(values[0].Data(), bCont); e != nil {
			return e
		}
		// Get Board.
		bValue, e := ct.GetObject("Board", bCont.Board)
		if e != nil {
			return e
		}
		if e := encoder.DeserializeRaw(bValue.Data(), b); e != nil {
			return e
		}
		return nil
	})
	return
}

// ObtainAllBoards obtains all the boards that we are subscribed to.
func (c *Client) ObtainAllBoards() (
	bList []*typ.Board, bContList []*typ.BoardContainer, e error,
) {
	// Get board configs.
	bConfList := c.bManager.GetList()
	// Prepare outputs.
	bList = make([]*typ.Board, len(bConfList))
	bContList = make([]*typ.BoardContainer, len(bConfList))
	// Loop.
	for i := 0; i < len(bConfList); i++ {
		b, bCont, err := c.ObtainBoard(bConfList[i].PubKey)
		if err != nil {
			e = err
			return
		}
		bList[i], bContList[i] = b, bCont
	}
	return
}
