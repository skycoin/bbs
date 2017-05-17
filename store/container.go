package store

import (
	"github.com/evanlinjin/bbs/cmd"
	"github.com/evanlinjin/bbs/store/typ"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/cipher"
	"strconv"
	"time"
)

type Container struct {
	*node.Container
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

	// Setup cxo client.
	c.client, e = node.NewClient(node.NewClientConfig(), skyobject.NewContainer(r))
	if e != nil {
		return
	}

	// Run cxo client.
	if e = c.client.Start("[::]:" + strconv.Itoa(c.config.CXOPort())); e != nil {
		return
	}

	// Set Container.
	c.Container = c.client.Container()

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
func (c *Container) ChangeBoardURL(bpk cipher.PubKey, url string) error {
	w := c.LastRoot(bpk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBcFinder()); e != nil {
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
	w := c.LastRoot(bpk).Walker()
	bc := typ.BoardContainer{}
	if e := w.AdvanceFromRoot(&bc, makeBcFinder()); e != nil {
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
		w := c.LastRoot(bpk).Walker()
		bc, b := typ.BoardContainer{}, typ.Board{}
		if e := w.AdvanceFromRoot(&bc, makeBcFinder()); e != nil {
			continue
		}
		w.AdvanceFromRefField("Board", &b)
		boards[i] = &b
	}
	return boards
}

// NewBoard attempts to create a new board from a given board and seed.
func (c *Container) NewBoard(board *typ.Board, seed string) (cipher.PubKey, cipher.SecKey, error) {
	bpk, bsk := board.TouchWithSeed([]byte(seed))
	r, e := c.NewRoot(bpk, bsk)
	if e != nil {
		return bpk, bsk, e
	}
	bRef := r.Save(*board)
	bCont := typ.BoardContainer{Board: bRef}
	_, _, e = r.Inject("BoardContainer", bCont)
	return bpk, bsk, e
}

// GetThreads attempts to obtain a list of threads from a board of public key.
func (c *Container) GetThreads(bpk cipher.PubKey) ([]*typ.Thread, error) {
	w := c.LastRoot(bpk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBcFinder()); e != nil {
		return nil, e
	}
	threads := make([]*typ.Thread, len(bc.Threads))
	for i, tRef := range bc.Threads {
		tData, has := c.Get(tRef)
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
func (c *Container) NewThread(bpk cipher.PubKey, thread *typ.Thread) (e error) {
	w := c.LastRoot(bpk).Walker()
	bc := &typ.BoardContainer{}
	if e = w.AdvanceFromRoot(bc, makeBcFinder()); e != nil {
		return e
	}
	var tRef skyobject.Reference
	if tRef, e = w.AppendToRefsField("Threads", *thread); e != nil {
		return e
	}
	_, e = w.AppendToRefsField("ThreadPages", typ.ThreadPage{Thread: tRef})
	thread.Ref = cipher.SHA256(tRef).Hex()
	return
}

// GetPosts attempts to obtain posts from a specified board and thread.
func (c *Container) GetPosts(bpk cipher.PubKey, tRef skyobject.Reference) ([]*typ.Post, error) {
	w := c.LastRoot(bpk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBcFinder()); e != nil {
		return nil, e
	}
	tp := &typ.ThreadPage{}
	if e := w.AdvanceFromRefsField("ThreadPages", tp, makeTpFinder(tRef)); e != nil {
		return nil, e
	}
	posts := make([]*typ.Post, len(tp.Posts))
	for i, pRef := range tp.Posts {
		pData, has := c.Get(pRef)
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
func (c *Container) NewPost(bpk cipher.PubKey, tRef skyobject.Reference, post *typ.Post) error {
	w := c.LastRoot(bpk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBcFinder()); e != nil {
		return e
	}
	tp := &typ.ThreadPage{}
	if e := w.AdvanceFromRefsField("ThreadPages", tp, makeTpFinder(tRef)); e != nil {
		return e
	}
	_, e := w.AppendToRefsField("Posts", *post)
	return e
}
