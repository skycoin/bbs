package state

import (
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/bbs/src/store/cxo/setup"
	"github.com/skycoin/cxo/skyobject"
	"testing"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/bbs/src/store/state/views"
	"github.com/skycoin/bbs/src/store/object"
	"fmt"
)

const (
	ListenAddress = "[::]:18998"

)

func prepareNode(t *testing.T) (*node.Node) {
	c := node.NewConfig()
	c.Skyobject.Registry = skyobject.NewRegistry(
		setup.PrepareRegistry)
	c.InMemoryDB = true
	c.EnableListener = true
	c.PublicServer = true
	c.Listen = ListenAddress
	c.EnableRPC = false
	c.RemoteClose = false

	out, e := node.NewNode(c)
	if e != nil {
		t.Fatal("failed to create cxo node:", e)
	}

	return out
}

func prepareBoard(t *testing.T, n *node.Node, seed string) (cipher.PubKey, cipher.SecKey, *skyobject.Root) {
	in := &object.NewBoardIO{
		Name: fmt.Sprintf("Board of seed '%s'", seed),
		Body: fmt.Sprintf("A test board created with seed '%s'.", seed),
	}
	if e := in.Process(nil); e != nil {
		t.Fatal("failed to process new board input:", e)
	}
	if e := n.AddFeed(in.BoardPubKey); e != nil {
		t.Fatal("failed to add feed:", e)
	}
	out, e := setup.NewBoard(n, in)
	if e != nil {
		t.Fatal("failed to create new board:", e)
	}
	return in.BoardPubKey, in.BoardSecKey, out
}

func prepareInstance(t *testing.T, n *node.Node, pk cipher.PubKey) *BoardInstance {
	return new(BoardInstance).Init(n, pk, views.AddContent(), views.AddFollow())
}

func initInstance(t *testing.T, seed string) (*BoardInstance, func()) {
	n := prepareNode(t)
	pk, sk, r := prepareBoard(t, n, seed)
	bi := prepareInstance(t, n, pk)

	if e := bi.UpdateWithReceived(r, sk); e != nil {
		t.Fatal("failed to update board instance:", e)
	}

	return bi, func() {
		bi.Close()
		n.Close()
	}
}

func TestBoardInstance_Init(t *testing.T) {
	const (
		bSeed = "a"
	)
	_, close := initInstance(t, bSeed)
	close()
}