package extern

import (
	"github.com/evanlinjin/bbs/cmd"
	"github.com/evanlinjin/bbs/store"
	"github.com/evanlinjin/bbs/typ"
	"github.com/skycoin/skycoin/src/cipher"
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
func (g *Gateway) NewBoard(name, desc, seed string) (bi store.BoardInfo, e error) {
	b := &typ.Board{Name: name, Desc: desc}
	bpk, bsk, e := g.container.NewBoard(b, seed)
	if e != nil {
		return
	}
	if e = g.boardSaver.MasterAdd(bpk, bsk); e != nil {
		return
	}
	bi, _ = g.boardSaver.Get(bpk)
	return
}