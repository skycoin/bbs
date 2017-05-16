package extern

import (
	"github.com/evanlinjin/bbs/cmd"
	"github.com/evanlinjin/bbs/store"
	"github.com/evanlinjin/bbs/typ"
	"github.com/skycoin/skycoin/src/cipher"
	"errors"
)

// Gateway represents the intermediate between External calls and internal processing.
// It can be seen as a security layer.
type Gateway struct {
	config     *cmd.Config
	container  *store.Container
	boardSaver *store.BoardSaver
}

// NewGateway creates a new Gateway.
func NewGateway(
	config *cmd.Config,
	container *store.Container,
	boardSaver *store.BoardSaver,
) *Gateway {
	return &Gateway{
		config:     config,
		container:  container,
		boardSaver: boardSaver,
	}
}

// GetSubscriptions lists all subscriptions.
func (g *Gateway) GetSubscriptions() []store.BoardInfo {
	return g.boardSaver.List()
}

// GetSubscription gets a subscription.
func (g *Gateway) GetSubscription(bpk cipher.PubKey) (store.BoardInfo, bool) {
	return g.boardSaver.Get(bpk)
}

// Subscribe subscribes to a board.
func (g *Gateway) Subscribe(bpk cipher.PubKey) {
	g.boardSaver.Add(bpk)
}

// Unsubscribe unsubscribes from a board.
func (g *Gateway) Unsubscribe(bpk cipher.PubKey) {
	g.boardSaver.Remove(bpk)
}

// GetBoards lists all boards.
func (g *Gateway) GetBoards() []*typ.Board {
	return g.container.GetBoards(g.boardSaver.ListKeys()...)
}

// NewBoard creates a new board.
func (g *Gateway) NewBoard(board *typ.Board, seed string) (bi store.BoardInfo, e error) {
	bpk, bsk, e := g.container.NewBoard(board, seed)
	if e != nil {
		return
	}
	if e = g.boardSaver.MasterAdd(bpk, bsk); e != nil {
		return
	}
	bi, _ = g.boardSaver.Get(bpk)
	return
}

// GetThreads obtains threads of boards we are subscribed to.
// Input `bpks` acts as a filter.
// If no `bpks` are specified, threads of all boards will be obtained.
// If one or more `bpks` are specified, only threads under those boards will be returned.
func (g *Gateway) GetThreads(bpks ...cipher.PubKey) []*typ.Thread {
	tMap := make(map[string]*typ.Thread)
	switch len(bpks) {
	case 0:
		for _, bpk := range g.boardSaver.ListKeys() {
			ts, e := g.container.GetThreads(bpk)
			if e != nil {
				continue
			}
			for _, t := range ts {
				tMap[t.Hash] = t
			}
		}
	default:
		for _, bpk := range bpks {
			if _, has := g.boardSaver.Get(bpk); has == false {
				return nil
			}
			ts, e := g.container.GetThreads(bpk)
			if e != nil {
				continue
			}
			for _, t := range ts {
				tMap[t.Hash] = t
			}
		}
	}
	threads, i := make([]*typ.Thread, len(tMap)), 0
	for _, t := range tMap {
		threads[i] = t
		i += 1
	}
	return threads
}

// NewThread creates a new thread and sets the board of public key as it's master.
func (g *Gateway) NewThread(bpk cipher.PubKey, thread *typ.Thread) error {
	bi, has := g.boardSaver.Get(bpk)
	if has == false {
		return errors.New("not subscribed to board")
	}
	if bi.BoardConfig.Master == true {
		// Via Container.
		if e := g.container.NewThread(bpk, thread); e != nil {
			return e
		}
	} else {
		// Via RPC Client.
		return errors.New("not implemented")
	}
	return nil
}