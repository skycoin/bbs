package state

import (
	"fmt"
	"github.com/skycoin/bbs/src/misc/keys"
	"github.com/skycoin/bbs/src/store/cxo/setup"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state/pack"
	"github.com/skycoin/bbs/src/store/state/views"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"testing"
)

const (
	ListenAddress = "[::]:18998"
)

func prepareNode(t *testing.T) *node.Node {
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
	if e := in.Process([]cipher.PubKey{}); e != nil {
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

func prepareInstance(_ *testing.T, n *node.Node, pk cipher.PubKey) *BoardInstance {
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

func obtainBoardPubKey(t *testing.T, bi *BoardInstance) cipher.PubKey {
	var pk cipher.PubKey
	e := bi.ViewPack(func(p *skyobject.Pack, h *pack.Headers) error {
		if p == nil {
			log.Println("*skyobject.Pack is empty")
		}
		root := p.Root()
		if root == nil {
			log.Println("*skyobject.Root is empty")
		}
		pk = root.Pub
		return nil
	})
	if e != nil {
		t.Fatal("failed to view pack:", e)
	}
	return pk
}

func obtainThreadList(t *testing.T, bi *BoardInstance) []cipher.SHA256 {
	var threads []cipher.SHA256
	bi.ViewPack(func(p *skyobject.Pack, h *pack.Headers) error {
		pages, e := object.GetPages(p, &object.GetPagesIn{
			RootPage:  false,
			BoardPage: true,
			DiffPage:  false,
			UsersPage: false,
		})
		if e != nil {
			return e
		}
		threads = make([]cipher.SHA256, pages.BoardPage.GetThreadCount())
		return pages.BoardPage.RangeThreadPages(
			func(i int, tp *object.ThreadPage) error {
				thread, e := tp.GetThread()
				if e != nil {
					t.Fatal("failed to extract thread:", e)
				}
				threads[i] = thread.GetHeader().GetHash()
				return nil
			},
		)
	})
	return threads
}

func addThread(t *testing.T, bi *BoardInstance, threadIndex int, userSeed []byte) (cipher.SHA256, uint64) {
	_, csk := cipher.GenerateDeterministicKeyPair(userSeed)
	in := &object.NewThreadIO{
		BoardPubKeyStr:   obtainBoardPubKey(t, bi).Hex(),
		Name:             fmt.Sprintf("Thread %d", threadIndex),
		Body:             fmt.Sprintf("A test thread created of index %d.", threadIndex),
		CreatorSecKeyStr: csk.Hex(),
	}
	if e := in.Process(); e != nil {
		t.Fatal("failed to process new thread input:", e)
	}
	goal, e := bi.Submit(in.Transport)
	if e != nil {
		t.Fatal("failed to create new thread:", e)
	}
	hash, _ := keys.GetHash(in.Transport.Header.Hash)
	return hash, goal
}

func addPost(t *testing.T, bi *BoardInstance, threadHash cipher.SHA256, postIndex int, userSeed []byte) uint64 {
	_, csk := cipher.GenerateDeterministicKeyPair(userSeed)
	in := &object.NewPostIO{
		BoardPubKeyStr:   obtainBoardPubKey(t, bi).Hex(),
		ThreadRefStr:     threadHash.Hex(),
		Name:             fmt.Sprintf("Post %d", postIndex),
		Body:             fmt.Sprintf("A test post created of index %d.", postIndex),
		CreatorSecKeyStr: csk.Hex(),
	}
	if e := in.Process(); e != nil {
		t.Fatal("failed to process new post input:", e)
	}
	goal, e := bi.Submit(in.Transport)
	if e != nil {
		t.Fatal("failed to create new post:", e)
	}
	return goal
}

func TestBoardInstance_Init(t *testing.T) {
	const (
		bSeed = "a"
	)
	_, close := initInstance(t, bSeed)
	close()
}
