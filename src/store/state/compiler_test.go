package state

import (
	"testing"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state/views"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/bbs/src/store/cxo/setup"
	"github.com/skycoin/skycoin/src/cipher"
	"context"
	"github.com/skycoin/net/skycoin-messenger/factory"
	"time"
	"github.com/skycoin/bbs/src/misc/boo"
)

func prepareMessengerServer(t *testing.T, address string) *factory.MessengerFactory {
	f := factory.NewMessengerFactory()
	if e := f.Listen(address); e != nil {
		t.Fatal(e)
	}
	time.Sleep(time.Second)
	return f
}

func prepareCompiler(t *testing.T, address string, discoveryAddresses []string) *Compiler {
	var (
		memMode = true
		haveDefaults = false
		updateInterval = 1
	)

	compilerConfig := &CompilerConfig{
		UpdateInterval: &updateInterval,
	}

	fileManager := object.NewCXOFileManager(&object.CXOFileManagerConfig{
		Memory:   &memMode,
		Defaults: &haveDefaults,
	})

	newRootsChan := make(chan RootWrap)

	c := node.NewConfig()
	{
		c.Skyobject.Registry = skyobject.NewRegistry(setup.PrepareRegistry)
		c.InMemoryDB = true
		c.EnableListener = true
		c.PublicServer = true
		c.Listen = address
		c.EnableRPC = false
		c.RemoteClose = false
		c.DiscoveryAddresses = discoveryAddresses

		c.OnRootReceived = func(c *node.Conn, root *skyobject.Root) {
			root.IsFull = false
			newRootsChan <- RootWrap{Root: root}
		}

		c.OnRootFilled = func(c *node.Conn, root *skyobject.Root) {
			root.IsFull = true
			newRootsChan <- RootWrap{Root: root}
		}
	}

	cxo, e := node.NewNode(c)
	if e != nil {
		t.Fatal(e)
	}

	compiler := NewCompiler(
		compilerConfig,
		fileManager,
		newRootsChan,
		cxo,
		views.AddContent(),
		views.AddFollow(),
	)

	return compiler
}

func closeCompiler(t *testing.T, c *Compiler) {
	c.Close()
	if e := c.node.Close(); e != nil {
		t.Fatal(e)
	}
}

func newBoard(c *Compiler, seed, name, body string) (*object.NewBoardIO, error) {
	in := &object.NewBoardIO{
		Name: name,
		Body: body,
		Seed: seed,
	}
	if e := in.Process([]cipher.PubKey{}); e != nil {
		return in, boo.Wrap(e, "in.Process")
	}

	if e := c.file.AddMasterSub(in.BoardPubKey, in.BoardSecKey); e != nil {
		return in, e
	}
	if e := c.node.AddFeed(in.BoardPubKey); e != nil {
		return in, e
	}

	r, e := setup.NewBoard(c.node, in)
	if e != nil {
		return in, boo.Wrap(e, "setup.NewBoard")
	}
	c.UpdateBoardWithContext(context.Background(), r)
	return in, nil
}

func subscribeMaster(c *Compiler, pk cipher.PubKey, sk cipher.SecKey) error {
	if e := c.file.AddMasterSub(pk, sk); e != nil {
		return e
	}
	return c.node.AddFeed(pk)
}

func TestNewCompiler(t *testing.T) {
	const (
		MessengerServerAddress = "[::]:11001"
		Node1Address           = "[::]:11002"
		Node2Address           = "[::]:11003"
		BoardSeed              = "a"
	)
	var (
		f         = prepareMessengerServer(t, MessengerServerAddress)
		compiler1 = prepareCompiler(t, Node1Address, []string{MessengerServerAddress})
		compiler2 = prepareCompiler(t, Node2Address, []string{MessengerServerAddress})
	)
	defer f.Close()
	defer closeCompiler(t, compiler1)
	defer closeCompiler(t, compiler2)

	in, e := newBoard(compiler1, BoardSeed, "Test Board V1", "A test board (v1).")
	if e != nil {
		t.Fatal(e)
	}

	time.Sleep(time.Second)

	if e := subscribeMaster(compiler2, in.BoardPubKey, in.BoardSecKey); e != nil {
		t.Fatal(e)
	}

	time.Sleep(10 * time.Second)
}