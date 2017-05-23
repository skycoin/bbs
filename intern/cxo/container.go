package cxo

import (
	"errors"
	"fmt"
	"github.com/evanlinjin/bbs/cmd"
	"github.com/evanlinjin/bbs/intern/typ"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"strconv"
	"time"
)

type Container struct {
	c      *node.Container
	client *node.Client
	config *cmd.Config
}

func NewContainer(config *cmd.Config) (c *Container, e error) {
	c = &Container{config: config}

	// Setup cxo registry.
	r := skyobject.NewRegistry()
	r.Register("Board", typ.Board{})
	r.Register("Thread", typ.Thread{})
	r.Register("Post", typ.Post{})
	r.Register("ThreadPage", typ.ThreadPage{})
	r.Register("BoardContainer", typ.BoardContainer{})
	r.Done()

	// Setup cxo config.
	cc := node.NewClientConfig()
	cc.InMemoryDB = config.CXOUseMemory()
	cc.DataDir = config.CXODir()

	// Setup cxo client.
	c.client, e = node.NewClient(cc, r)
	if e != nil {
		return
	}

	// Run cxo client.
	if e = c.client.Start("[::]:" + strconv.Itoa(c.config.CXOPort())); e != nil {
		return
	}

	// Set Container.
	c.c = c.client.Container()

	// Wait.
	time.Sleep(5 * time.Second)
	return
}

func (c *Container) Close() error                      { return c.client.Close() }
func (c *Container) Connected() bool                   { return c.client.IsConnected() }
func (c *Container) Feeds() []cipher.PubKey            { return c.client.Feeds() }
func (c *Container) Subscribe(pk cipher.PubKey) bool   { return c.client.Subscribe(pk) }
func (c *Container) Unsubscribe(pk cipher.PubKey) bool { return c.client.Unsubscribe(pk) }

// ChangeBoardURL changes the board's URL of given public key.
func (c *Container) ChangeBoardURL(bpk cipher.PubKey, bsk cipher.SecKey, url string) error {
	r := c.c.LastRootSk(bpk, bsk)
	w := r.Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder()); e != nil {
		return e
	}
	b := &typ.Board{}
	if e := w.AdvanceFromRefField("Board", b); e != nil {
		return e
	}
	b.URL = url
	w.Retreat()
	_, e := w.ReplaceInRefField("Board", *b)
	return e
}

// GetBoard attempts to obtain the board of a given public key.
func (c *Container) GetBoard(bpk cipher.PubKey) (*typ.Board, error) {
	w := c.c.LastRoot(bpk).Walker()
	bc := typ.BoardContainer{}
	if e := w.AdvanceFromRoot(&bc, makeBoardContainerFinder()); e != nil {
		return nil, e
	}
	b := &typ.Board{}
	e := w.AdvanceFromRefField("Board", b)
	return b, e
}

// GetBoards attempts to obtain a list of boards from the given public keys.
func (c *Container) GetBoards(bpks ...cipher.PubKey) []*typ.Board {
	boards := make([]*typ.Board, len(bpks))
	for i, bpk := range bpks {
		w := c.c.LastRoot(bpk).Walker()
		bc, b := typ.BoardContainer{}, typ.Board{}
		if e := w.AdvanceFromRoot(&bc, makeBoardContainerFinder()); e != nil {
			continue
		}
		w.AdvanceFromRefField("Board", &b)
		boards[i] = &b
	}
	return boards
}

// NewBoard attempts to create a new board from a given board and seed.
func (c *Container) NewBoard(board *typ.Board, pk cipher.PubKey, sk cipher.SecKey) error {
	r, e := c.c.NewRoot(pk, sk)
	if e != nil {
		fmt.Println("[Container.NewBoard] [113] Error:", e)
		return e
	}
	bRef := r.Save(*board)
	bCont := typ.BoardContainer{Board: bRef}
	_, _, e = r.Inject("BoardContainer", bCont)
	if e != nil {
		fmt.Println("[Container.NewBoard] [120] Error:", e)
	}
	return e
}

// RemoveBoard attempts to remove a board by a given public key.
func (c *Container) RemoveBoard(bpk cipher.PubKey, bsk cipher.SecKey) error {
	w := c.c.LastRootSk(bpk, bsk).Walker()
	fmt.Println("Removing board:", bpk.Hex())
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder()); e != nil {
		return e
	}

	return w.RemoveCurrent()
}

// GetThreads attempts to obtain a list of threads from a board of public key.
func (c *Container) GetThreads(bpk cipher.PubKey) ([]*typ.Thread, error) {
	w := c.c.LastRoot(bpk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder()); e != nil {
		return nil, e
	}
	threads := make([]*typ.Thread, len(bc.Threads))
	for i, tRef := range bc.Threads {
		tData, has := c.c.Get(tRef)
		if has == false {
			continue
		}
		threads[i] = new(typ.Thread)
		if e := encoder.DeserializeRaw(tData, threads[i]); e != nil {
			return nil, e
		}
		threads[i].Ref = cipher.SHA256(tRef).Hex()
	}
	return threads, nil
}

// NewThread attempts to create a new thread from a board of given public key.
func (c *Container) NewThread(bpk cipher.PubKey, bsk cipher.SecKey, thread *typ.Thread) (e error) {
	w := c.c.LastRootSk(bpk, bsk).Walker()
	bc := &typ.BoardContainer{}
	if e = w.AdvanceFromRoot(bc, makeBoardContainerFinder()); e != nil {
		return e
	}
	thread.MasterBoard = bpk.Hex()
	var tRef skyobject.Reference
	if tRef, e = w.AppendToRefsField("Threads", *thread); e != nil {
		return e
	}
	_, e = w.AppendToRefsField("ThreadPages", typ.ThreadPage{Thread: tRef})
	thread.Ref = cipher.SHA256(tRef).Hex()
	return
}

// RemoveThread attempts to remove a thread from a board of given public key.
func (c *Container) RemoveThread(bpk cipher.PubKey, bsk cipher.SecKey, tRef skyobject.Reference) (e error) {
	w := c.c.LastRootSk(bpk, bsk).Walker()
	bc := &typ.BoardContainer{}
	if e = w.AdvanceFromRoot(bc, makeBoardContainerFinder()); e != nil {
		return e
	}

	fmt.Println("Removing thread:", tRef.String())
	e = w.RemoveInRefsByRef("Threads", tRef)
	if e != nil {
		return errors.New("remove thread from threads failed: " + e.Error())
	}

	fmt.Println("Removing from thread pages")
	e = w.RemoveInRefsField("ThreadPages", makeThreadPageFinder(tRef))
	if e != nil {
		return errors.New("remove thread from threadpages failed: " + e.Error())
	}

	return nil
}

// GetThreadPage requests a page from a thread
func (c *Container) GetThreadPage(bpk cipher.PubKey, tRef skyobject.Reference) (*typ.Thread, []*typ.Post, error) {
	w := c.c.LastRoot(bpk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder()); e != nil {
		return nil, nil, e
	}
	// Get thread.
	tData, has := c.c.Get(tRef)
	if has == false {
		return nil, nil, errors.New("unable to obtain thread")
	}
	thread := new(typ.Thread)
	if e := encoder.DeserializeRaw(tData, thread); e != nil {
		return nil, nil, e
	}
	thread.Ref = cipher.SHA256(tRef).Hex()
	// Get posts.
	tp := &typ.ThreadPage{}
	if e := w.AdvanceFromRefsField("ThreadPages", tp, makeThreadPageFinder(tRef)); e != nil {
		return nil, nil, e
	}
	posts := make([]*typ.Post, len(tp.Posts))
	for i, pRef := range tp.Posts {
		pData, has := c.c.Get(pRef)
		if has == false {
			continue
		}
		posts[i] = new(typ.Post)
		if e := encoder.DeserializeRaw(pData, posts[i]); e != nil {
			return nil, nil, e
		}
	}
	return thread, posts, nil
}

// GetPosts attempts to obtain posts from a specified board and thread.
func (c *Container) GetPosts(bpk cipher.PubKey, tRef skyobject.Reference) ([]*typ.Post, error) {
	w := c.c.LastRoot(bpk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder()); e != nil {
		return nil, e
	}
	tp := &typ.ThreadPage{}
	if e := w.AdvanceFromRefsField("ThreadPages", tp, makeThreadPageFinder(tRef)); e != nil {
		return nil, e
	}
	posts := make([]*typ.Post, len(tp.Posts))
	for i, pRef := range tp.Posts {
		pData, has := c.c.Get(pRef)
		if has == false {
			continue
		}
		posts[i] = new(typ.Post)
		if e := encoder.DeserializeRaw(pData, posts[i]); e != nil {
			return nil, e
		}
	}
	return posts, nil
}

// NewPost attempts to create a new post in a given board and thread.
func (c *Container) NewPost(bpk cipher.PubKey, bsk cipher.SecKey, tRef skyobject.Reference, post *typ.Post) error {
	w := c.c.LastRootSk(bpk, bsk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder()); e != nil {
		return e
	}
	tp := &typ.ThreadPage{}
	if e := w.AdvanceFromRefsField("ThreadPages", tp, makeThreadPageFinder(tRef)); e != nil {
		return e
	}
	t := &typ.Thread{}
	if _, e := w.GetFromRefField("Thread", t); e != nil {
		return e
	}
	if t.MasterBoard != bpk.Hex() {
		return errors.New("this board is not master of this thread")
	}
	_, e := w.AppendToRefsField("Posts", *post)
	return e
}

// RemovePost attempts to remove a post in a given board and thread.
func (c *Container) RemovePost(bpk cipher.PubKey, bsk cipher.SecKey, tRef, pRef skyobject.Reference) error {
	w := c.c.LastRootSk(bpk, bsk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder()); e != nil {
		return e
	}
	tp := &typ.ThreadPage{}
	if e := w.AdvanceFromRefsField("ThreadPages", tp, makeThreadPageFinder(tRef)); e != nil {
		return e
	}
	t := &typ.Thread{}
	if _, e := w.GetFromRefField("Thread", t); e != nil {
		return e
	}
	if t.MasterBoard != bpk.Hex() {
		return errors.New("this board is not master of this thread")
	}
	if e := w.RemoveInRefsByRef("Posts", pRef); e != nil {
		return errors.New("remove post from posts failed: " + e.Error())
	}
	return nil
}

// ImportThread imports a thread from a board to another board (which this node owns). If already imported replaces it.
func (c *Container) ImportThread(fromBpk, toBpk cipher.PubKey, toBsk cipher.SecKey, tRef skyobject.Reference) error {
	// Get from 'from' Board.
	w := c.c.LastRoot(fromBpk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder()); e != nil {
		return e
	}

	// Obtain thread and thread page.
	tp := &typ.ThreadPage{}
	if e := w.AdvanceFromRefsField("ThreadPages", tp, makeThreadPageFinder(tRef)); e != nil {
		return e
	}
	t := &typ.Thread{}
	if _, e := w.GetFromRefField("Thread", t); e != nil {
		return e
	}

	// Get from 'to' Board.
	w = c.c.LastRootSk(toBpk, toBsk).Walker()
	bc = &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder()); e != nil {
		return e
	}
	if e := w.ReplaceInRefsField("ThreadPages", *tp, makeThreadPageFinder(tRef)); e != nil {
		/* THREAD DOES NOT EXIST */
		// Append thread and thread page.
		if _, e := w.AppendToRefsField("Threads", *t); e != nil {
			return e
		}
		if _, e := w.AppendToRefsField("ThreadPages", *tp); e != nil {
			return e
		}
		return nil
	}
	return nil
}
