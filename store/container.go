package store

import (
	"github.com/evanlinjin/bbs/cmd"
	"github.com/evanlinjin/bbs/typ"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
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

	// Wait.
	time.Sleep(5 * time.Second)
	return
}

func (c *Container) Stop() error                       { return c.client.Close() }
func (c *Container) Connected() bool                   { return c.client.IsConnected() }
func (c *Container) Feeds() []cipher.PubKey            { return c.client.Feeds() }
func (c *Container) Subscribe(pk cipher.PubKey) bool   { return c.client.Subscribe(pk) }
func (c *Container) Unsubscribe(pk cipher.PubKey) bool { return c.client.Unsubscribe(pk) }

func (c *Container) GetBoard(bpk cipher.PubKey) (*node.RootWalker, *typ.Board, error) {
	w := c.LastRoot(bpk).Walker()
	bc := typ.BoardContainer{}
	if e := w.AdvanceFromRoot(&bc, findBoardContainer); e != nil {
		return w, nil, e
	}
	b := &typ.Board{}
	e := w.AdvanceFromRefField("Board", b)
	return w, b, e
}

func (c *Container) GetBoards(bpks ...cipher.PubKey) []*typ.Board {
	boards := make([]*typ.Board, len(bpks))
	for i, bpk := range bpks {
		w := c.LastRoot(bpk).Walker()
		bc, b := typ.BoardContainer{}, typ.Board{}
		if e := w.AdvanceFromRoot(&bc, findBoardContainer); e != nil {
			continue
		}
		w.AdvanceFromRefField("Board", &b)
		boards[i] = &b
	}
	return boards
}

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
