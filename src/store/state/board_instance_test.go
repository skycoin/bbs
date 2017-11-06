package state

import (
	"fmt"
	"github.com/skycoin/bbs/src/store/cxo/setup"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"testing"
	"time"
	"github.com/skycoin/bbs/src/misc/tag"
	"encoding/json"
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
	bpk, bsk := cipher.GenerateDeterministicKeyPair([]byte(seed))
	content := new(object.Content)
	content.SetBody(&object.Body{
		Type: object.V5BoardType,
		Name: fmt.Sprintf("Board of seed '%s'", seed),
		Body: fmt.Sprintf("A test board created with seed '%s'.", seed),
		TS:   time.Now().UnixNano(),
	})
	var (
		hash = cipher.SumSHA256(content.Body)
		sig  = cipher.SignHash(hash, bsk)
	)
	content.SetHeader(&object.ContentHeaderData{
		Hash: hash.Hex(),
		Sig:  sig.Hex(),
	})

	if e := n.AddFeed(bpk); e != nil {
		t.Fatal("failed to add feed:", e)
	}

	out, e := setup.NewBoard(n, content, bpk, bsk)
	if e != nil {
		t.Fatal("failed to create new board:", e)
	}
	return bpk, bsk, out
}

func prepareInstance(_ *testing.T, n *node.Node, pk cipher.PubKey) *BoardInstance {
	return new(BoardInstance).Init(n, pk)
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
	e := bi.ViewPack(func(p *skyobject.Pack, h *Headers) error {
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
	bi.ViewPack(func(p *skyobject.Pack, h *Headers) error {
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
	cpk, csk := cipher.GenerateDeterministicKeyPair(userSeed)
	body := &object.Body{
		Type: object.V5ThreadType,
		TS:   time.Now().UnixNano(),
		OfBoard: obtainBoardPubKey(t, bi).Hex(),
		Name: fmt.Sprintf("Thread %d", threadIndex),
		Body: fmt.Sprintf("A test thread created of index %d.", threadIndex),
		Creator: cpk.Hex(),
	}
	raw, _ := json.Marshal(body)
	sig := cipher.SignHash(cipher.SumSHA256(raw), csk)
	transport, e := object.NewTransport(raw, sig)
	if e != nil {
		t.Fatal("failed to generate transport:", e)
	}
	goal, e := bi.Submit(transport)
	if e != nil {
		t.Fatal("failed to create new thread:", e)
	}
	hash, e := tag.GetHash(transport.Header.Hash)
	if e != nil {
		t.Fatal("")
	}
	return hash, goal
}

func addPost(t *testing.T, bi *BoardInstance, threadHash cipher.SHA256, postIndex int, userSeed []byte) uint64 {
	cpk, csk := cipher.GenerateDeterministicKeyPair(userSeed)
	body := &object.Body{
		Type: object.V5PostType,
		TS: time.Now().UnixNano(),
		OfBoard: obtainBoardPubKey(t, bi).Hex(),
		OfThread: threadHash.Hex(),
		Name: fmt.Sprintf("Post %d", postIndex),
		Body: fmt.Sprintf("A test post created of index %d.", postIndex),
		Creator: cpk.Hex(),
	}
	raw, _ := json.Marshal(body)
	sig := cipher.SignHash(cipher.SumSHA256(raw), csk)
	transport, e := object.NewTransport(raw, sig)
	if e != nil {
		t.Fatal("failed to generate transport:", e)
	}
	goal, e := bi.Submit(transport)
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

func TestBoardInstance_UpdateWithReceived(t *testing.T) {
	const (
		MessengerServerAddress = "[::]:11001"
		Compiler1Address       = "[::]:11002"
		Compiler2Address       = "[::]:11003"
		DisruptorAddress       = "[::]:11004"
		BoardSeedA             = "board a"
		BoardSeedB             = "board b"
	)
	var (
		f                 = prepareMessengerServer(t, MessengerServerAddress)
		compiler1RootChan = make(chan *skyobject.Root)
		compiler1         = prepareCompiler(t, Compiler1Address, []string{MessengerServerAddress},
			func(c *node.Conn, root *skyobject.Root) {
				go func() {
					compiler1RootChan <- root
				}()
			})
		compiler2RootChan = make(chan *skyobject.Root)
		compiler2         = prepareCompiler(t, Compiler2Address, []string{MessengerServerAddress},
			func(c *node.Conn, root *skyobject.Root) {
				go func() {
					compiler2RootChan <- root
				}()
			})
		disruptorRootChan = make(chan *skyobject.Root)
		disruptor         = prepareDisruptor(t, DisruptorAddress, []string{MessengerServerAddress},
			func(c *node.Conn, root *skyobject.Root) {
				go func() {
					disruptorRootChan <- root
				}()
			})
	)
	defer f.Close()
	defer closeCompiler(t, compiler1)
	defer closeCompiler(t, compiler2)
	defer disruptor.Close()

	t.Run("node_deals_with_invalid_master_root", func(t *testing.T) {
		bpk, bsk, e := newBoard(compiler1, BoardSeedA, "Test Board A", "A test board (A).")
		if e != nil {
			t.Fatal(e)
		}

		if e := disruptor.AddFeed(bpk); e != nil {
			t.Fatal(e)
		}

		// Wait for valid root to be received by disruptor.
		root := <-disruptorRootChan
		if len(root.Refs) != object.RootChildrenCount {
			t.Fatalf("disruptor received invalid root: child_count(%d) expected(%d)",
				len(root.Refs), object.RootChildrenCount)
		}

		if e := performDisruption(t, disruptor, bpk, bsk); e != nil {
			t.Fatal(e)
		}

		// Wait for invalid root to be received by compiler1.
		root = <-compiler1RootChan
		if len(root.Refs) == object.RootChildrenCount {
			t.Fatal("compiler1 received valid root, when expecting something invalid")
		}

		// Wait for valid root to be received by disruptor.
		root = <-disruptorRootChan
		if len(root.Refs) != object.RootChildrenCount {
			t.Fatalf("disruptor received invalid root: child_count(%d) expected(%d)",
				len(root.Refs), object.RootChildrenCount)
		}
	})

	t.Run("node_deals_with_invalid_remote_root", func(t *testing.T) {
		bpk, bsk, e := newBoard(compiler1, BoardSeedB, "Test Board B", "A test board (B).")
		if e != nil {
			t.Fatal(e)
		}

		if e := disruptor.AddFeed(bpk); e != nil {
			t.Fatal(e)
		}
		<-disruptorRootChan

		if e := subscribeRemote(compiler2, bpk); e != nil {
			t.Fatal(e)
		}
		<-compiler2RootChan

		for i := 0; i < 10; i++ {
			if e := performDisruption(t, disruptor, bpk, bsk); e != nil {
				t.Fatal(e)
			}
			<-compiler1RootChan
		}
	})
}
