package cxo

import (
	"errors"
	"fmt"
	"github.com/evanlinjin/bbs/typ"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"strconv"
	"time"
)

// CXOConfig represents a configuration for CXO Client.
type CXOConfig struct {
	Master bool
	Port   int
	RPCAddr string
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
	config *CXOConfig
	cxo    *node.Client
	c      *node.Container
	U      *typ.UserManager
	B      *typ.BoardManager
}

// NewClient creates a new Client.
func NewClient(conf *CXOConfig) (*Client, error) {

	// Setup cxo client.
	r := skyobject.NewRegistry()
	r.Register("Board", typ.Board{})
	r.Register("Thread", typ.Thread{})
	r.Register("Post", typ.Post{})
	r.Register("ThreadPage", typ.ThreadPage{})
	r.Register("BoardContainer", typ.BoardContainer{})
	r.Done()

	client, e := node.NewClient(node.NewClientConfig(), skyobject.NewContainer(r))
	if e != nil {
		return nil, e
	}
	c := Client{
		config: conf,
		cxo:    client,
		c:      client.Container(),
		U:      typ.NewUserManager(),
		B:      typ.NewBoardManager(conf.Master, conf.RPCAddr),
	}
	c.initUsers()
	c.initBoards()
	return &c, nil
}

func (c *Client) initUsers() {
	m := c.U
	if m.Load() != nil {
		m.Clear()
		m.AddNewRandomMaster()
		if e := m.Save(); e != nil {
			panic(e)
		}
	}
}

func (c *Client) initBoards() {
	m := c.B
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
	// Load boards from BoardManager.
	for _, bc := range c.B.Boards {

		// [USE THIS PART WHEN KONSTANTIN GIVES ME FIX]
		_, e := c.Subscribe(bc.PubKey)
		if e != nil {
			fmt.Println("[BOARD CONFIG]", bc.PubKey.Hex(), e)
			fmt.Println("[BOARD CONFIG] Removing", bc.PubKey.Hex())
			c.B.RemoveConfig(bc.PubKey)
		}
	}
	return nil
}

// Shutdown shutdowns the Client.
func (c *Client) Shutdown() error {
	return c.cxo.Close()
}

// Subscribe subscribes to a board.
func (c *Client) Subscribe(pk cipher.PubKey) (*typ.BoardConfig, error) {
	// Check if we already have the config, if not make one.
	bc, e := c.B.GetConfig(pk)
	if e != nil {
		bc, _ = typ.NewBoardConfig(pk, "")
	}
	// Attempt subscribe in cxo.
	c.cxo.Subscribe(pk)

	// Add to BoardManager.
	c.B.AddConfig(bc)
	return bc, nil
}

// Unsubscribe unsubscribes from a board.
func (c *Client) Unsubscribe(pk cipher.PubKey) bool {
	c.B.RemoveConfig(pk)
	return c.cxo.Unsubscribe(pk)
}

// InjectBoard injects a new master board with given BoardConfig.
func (c *Client) InjectBoard(bc *typ.BoardConfig) (
	b *typ.Board, bCont *typ.BoardContainer, e error,
) {
	if c.config.Master == false {
		e = errors.New("action not allowed; node is not master")
		return
	}
	// Make Board from BoardConfig.
	b = typ.NewBoardFromConfig(bc)
	// Check if BoardConfig exists.
	if c.B.HasConfig(bc.PubKey) == true {
		e = errors.New("board already exists")
		return
	}
	// Add to cxo.
	{
		var r *node.Root
		r, e = c.c.NewRoot(bc.PubKey, bc.SecKey)

		boardRef := r.Save(*b)
		bCont = typ.NewBoardContainer(boardRef)
		r.Inject("BoardContainer", *bCont)
	}

	// Subscribe to board.
	c.cxo.Subscribe(bc.PubKey)

	// Add config to BoardManager.
	e = c.B.AddConfig(bc)
	return
}

// ObtainBoard obtains a board of specified PubKey from cxo.
// Also obtains the BoardContainer.
// Assumes public key provided is valid.
func (c *Client) ObtainBoard(bpk cipher.PubKey) (
	b *typ.Board, bCont *typ.BoardContainer, e error,
) {
	bCont = new(typ.BoardContainer)
	b = new(typ.Board)

	// Obtain root.
	r := c.c.LastRoot(bpk)
	if r == nil {
		c.B.RemoveConfig(bpk)
		e = errors.New("nil root, removed from config")
		return
	}
	// Obtain root values.
	var values []*skyobject.Value
	values, e = r.Values()
	if e != nil {
		return
	}
	// Loop through and get latest BoardContainer.
	for _, v := range values {
		bContTemp := &typ.BoardContainer{}
		// Get BoardContainer.
		if e = encoder.DeserializeRaw(v.Data(), bContTemp); e != nil {
			return
		}
		if bContTemp.Seq >= bCont.Seq {
			bCont = bContTemp
		}
	}
	// Get Board.
	var bData []byte
	var has bool
	bData, has = r.Get(bCont.Board)
	if has == false {
		e = errors.New("unable to obtain board")
		return
	}
	if e = encoder.DeserializeRaw(bData, b); e != nil {
		return
	}
	b.PubKey = bpk.Hex()
	return
}

// ObtainAllBoards obtains all the boards that we are subscribed to.
func (c *Client) ObtainAllBoards() (
	bList []*typ.Board, bContList []*typ.BoardContainer, e error,
) {
	// Get board configs.
	bConfList := c.B.GetList()
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

// InjectThread injects a thread in specified board.
func (c *Client) InjectThread(bpk cipher.PubKey, thread *typ.Thread) (
	b *typ.Board, bCont *typ.BoardContainer, tList []*typ.Thread, e error,
) {
	if c.config.Master == false {
		e = errors.New("action not allowed; node is not master")
		return
	}
	// Check thread.
	if e = thread.CheckAndPrep(); e != nil {
		return
	}
	// Obtain board configuration.
	bc, e := c.B.GetConfig(bpk)
	if e != nil {
		return
	}
	// See if we are master of the board.
	if bc.Master == false {
		e = errors.New("not master")
		return
	}
	// Obtain board and board container.
	b, bCont, e = c.ObtainBoard(bpk)
	if e != nil {
		return
	}
	{
		var (
			tp    *typ.ThreadPage
			tpRef skyobject.Reference
		)
		// Obtain root.
		var r *node.Root
		r, e = c.c.NewRoot(bc.PubKey, bc.SecKey)
		if e != nil {
			return
		}
		// Save thread.
		tRef := r.Save(*thread)
		// Add thread to BoardContainer if not exist.
		if e = bCont.AddThread(tRef); e != nil {
			goto InjectThreadListThreads
		}
		// Add empty ThreadPage to BoardContainer.
		tp = typ.NewThreadPage(tRef)
		tpRef = r.Save(*tp)
		bCont.AddThreadPage(tpRef)
		// Increment sequence of BoardContainer.
		bCont.Touch()
		r.Inject("BoardContainer", *bCont)

	InjectThreadListThreads:
		// Obtain all threads in board.
		tList = make([]*typ.Thread, len(bCont.Threads))
		for i, tRef := range bCont.Threads {
			tList[i] = typ.InitThread(tRef)
			tData, has := r.Get(tRef)
			if has == false {
				e = errors.New("unable to retrieve thread")
				return
			}
			if e = encoder.DeserializeRaw(tData, tList[i]); e != nil {
				return
			}
		}
	}
	return
}

// ObtainThread obtains a single thread.
func (c *Client) ObtainThread(bpk cipher.PubKey, tRef skyobject.Reference) (
	b *typ.Board, bCont *typ.BoardContainer, t *typ.Thread, e error,
) {
	// Obtain board and board container.
	b, bCont, e = c.ObtainBoard(bpk)
	if e != nil {
		return
	}

	// Find thread.
	t = &typ.Thread{}
	tData, has := c.c.Get(tRef)
	if has == false {
		e = errors.New("unable to find thread")
		return
	}
	e = encoder.DeserializeRaw(tData, t)
	return
}

// ObtainThreads obtains threads of a specified board.
func (c *Client) ObtainThreads(bpk cipher.PubKey) (
	b *typ.Board, bCont *typ.BoardContainer, tList []*typ.Thread, e error,
) {
	// Obtain board and board container.
	b, bCont, e = c.ObtainBoard(bpk)
	if e != nil {
		return
	}

	// Obtain all threads in board.
	tList = make([]*typ.Thread, len(bCont.Threads))
	for i, tRef := range bCont.Threads {
		tList[i] = typ.InitThread(tRef)
		tData, has := c.c.Get(tRef)
		if has == false {
			e = errors.New("unable to find thread")
			return
		}
		if e = encoder.DeserializeRaw(tData, tList[i]); e != nil {
			return
		}
	}
	return
}

// ObtainPosts obtains the posts of specified board and thread.
func (c *Client) ObtainPosts(bpk cipher.PubKey, tRef skyobject.Reference) (
	b *typ.Board, bCont *typ.BoardContainer, t *typ.Thread, tPage *typ.ThreadPage, pList []*typ.Post,
	tpRef skyobject.Reference, e error,
) {
	// Obtain board, board container and thread.
	b, bCont, t, e = c.ObtainThread(bpk, tRef)
	if e != nil {
		return
	}

	// Find thread page.
	tpFound := false
	tPage = &typ.ThreadPage{}
	for _, tpRef = range bCont.ThreadPages {
		tpData, has := c.c.Get(tpRef)
		if has == false {
			e = errors.New("unable to find thread")
			return
		}
		if e = encoder.DeserializeRaw(tpData, tPage); e != nil {
			return
		}
		if tPage.Thread == tRef {
			tpFound = true
			break
		}
	}
	if tpFound == false {
		e = errors.New("thread page not found")
		return
	}
	// Obtain posts from thread page.
	pList = make([]*typ.Post, len(tPage.Posts))
	for i, pRef := range tPage.Posts {
		pData, has := c.c.Get(pRef)
		if has == false {
			e = errors.New("post not found")
			return
		}
		pList[i] = &typ.Post{}
		if e = encoder.DeserializeRaw(pData, pList[i]); e != nil {
			return
		}
		if e = pList[i].CheckCreator(); e != nil {
			return
		}
	}
	return
}

// InjectPost injects a post in specified Board and Thread.
// Note that post needs to be signed properly.
func (c *Client) InjectPost(bpk cipher.PubKey, tRef skyobject.Reference, post *typ.Post) (
	b *typ.Board, bCont *typ.BoardContainer,
	t *typ.Thread, tPage *typ.ThreadPage, pList []*typ.Post, e error,
) {
	// Obtain posts and whatnot.
	oldTpRef := skyobject.Reference{}
	b, bCont, t, tPage, pList, oldTpRef, e = c.ObtainPosts(bpk, tRef)
	if e != nil {
		return
	}
	// Obtain board configuration.
	bc, e := c.B.GetConfig(bpk)
	if e != nil {
		return
	}
	// See if we are master of the board.
	if bc.Master == false {
		e = errors.New("not master")
		return
	}
	// Check post to inject.
	if e = post.CheckContent(); e != nil {
		return
	}
	if e = post.CheckCreator(); e != nil {
		return
	}
	if e = post.CheckSig(); e != nil {
		return
	}
	// Touch post.
	post.Touch()

	// Obtain Root.
	var r *node.Root
	r, e = c.c.NewRoot(bc.PubKey, bc.SecKey)
	if e != nil {
		return
	}

	// Save post in cxo container and obtain it's reference.
	pRef := r.Save(*post)

	// Add post to thread page and save new thread page in cxo container.
	// Hence, obtain the new thread page's reference.
	tPage.AddPost(pRef)
	newTpRef := r.Save(tPage)

	// Replace thread page reference in board container.
	// Hence, increment sequence of ThreadContainer.
	if e = bCont.ReplaceThreadPage(oldTpRef, newTpRef); e != nil {
		return
	}
	bCont.Touch()

	// Inject new board container in root.
	_, _, e = r.Inject("BoardContainer", bCont)
	if e != nil {
		return
	}

	// Add post to output post list.
	pList = append(pList, post)
	return
}
