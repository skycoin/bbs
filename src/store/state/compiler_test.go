package state

import (
	"context"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/cxo/setup"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/net/skycoin-messenger/factory"
	"github.com/skycoin/skycoin/src/cipher"
	"testing"
	"time"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

func prepareMessengerServer(t *testing.T, address string) *factory.MessengerFactory {
	f := factory.NewMessengerFactory()
	if e := f.Listen(address); e != nil {
		t.Fatal(e)
	}
	time.Sleep(time.Second)
	return f
}

func prepareCompiler(t *testing.T, address string, discoveryAddresses []string, cb func(c *node.Conn, root *skyobject.Root)) *Compiler {
	var (
		memMode        = true
		updateInterval = 1
	)

	compilerConfig := &CompilerConfig{
		UpdateInterval: &updateInterval,
	}

	fileManager := object.NewCXOFileManager(&object.CXOFileManagerConfig{
		Memory: &memMode,
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
			cb(c, root)
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
	)

	return compiler
}

func closeCompiler(t *testing.T, c *Compiler) {
	c.Close()
	if e := c.node.Close(); e != nil {
		t.Fatal(e)
	}
}

func newBoard(c *Compiler, seed, name, body string) (cipher.PubKey, cipher.SecKey, error) {
	pk, sk := cipher.GenerateDeterministicKeyPair([]byte(seed))
	data := &object.Body{
		Type: object.V5BoardType,
		TS: time.Now().UnixNano(),
		Name: name,
		Body: body,
	}
	content := new(object.Content)
	content.SetBody(data)
	content.SetHeader(&object.ContentHeaderData{
		Hash: cipher.SumSHA256(encoder.Serialize(data)).Hex(),
	})

	if e := c.file.AddMasterSub(pk, sk); e != nil {
		return cipher.PubKey{}, cipher.SecKey{}, e
	}
	if e := c.node.AddFeed(pk); e != nil {
		return cipher.PubKey{}, cipher.SecKey{}, e
	}

	r, e := setup.NewBoard(c.node, content, pk, sk)
	if e != nil {
		return cipher.PubKey{}, cipher.SecKey{}, boo.Wrap(e, "setup.NewBoard")
	}
	c.UpdateBoardWithContext(context.Background(), r)
	return pk, sk, nil
}

func subscribeMaster(c *Compiler, pk cipher.PubKey, sk cipher.SecKey) error {
	if e := c.file.AddMasterSub(pk, sk); e != nil {
		return e
	}
	return c.node.AddFeed(pk)
}

func subscribeRemote(c *Compiler, pk cipher.PubKey) error {
	if e := c.file.AddRemoteSub(pk); e != nil {
		return e
	}
	return c.node.AddFeed(pk)
}

type Woops struct {
	Situation string
	Annoyance uint64
}

type Haha struct {
	Situation    string
	Humour       uint64
	Participants skyobject.Refs `skyobject:"schema=Participant"`
}

type Participant struct {
	Name string
	Age  uint64
}

func prepareDisruptor(t *testing.T, address string, discoveryAddresses []string, cb func(c *node.Conn, root *skyobject.Root)) *node.Node {
	c := node.NewConfig()
	{
		c.Skyobject.Registry = skyobject.NewRegistry(func(r *skyobject.Reg) {
			r.Register("Woops", Woops{})
			r.Register("Haha", Haha{})
			r.Register("Participant", Participant{})
		})
		c.InMemoryDB = true
		c.EnableListener = true
		c.PublicServer = true
		c.Listen = address
		c.EnableRPC = false
		c.RemoteClose = false
		c.DiscoveryAddresses = discoveryAddresses
		c.OnRootFilled = cb
	}
	cxo, e := node.NewNode(c)
	if e != nil {
		t.Fatal(e)
	}
	return cxo
}

func performDisruption(t *testing.T, n *node.Node, pk cipher.PubKey, sk cipher.SecKey) error {
	if e := n.AddFeed(pk); e != nil {
		return e
	}
	time.Sleep(time.Second)
	r, e := n.Container().LastRoot(pk)
	if e != nil {
		return e
	}
	goal := r.Seq

	p, e := n.Container().NewRoot(pk, sk, 0, n.Container().CoreRegistry().Types())
	if e != nil {
		return e
	}
	p.Append(
		&Woops{
			Situation: "Coffee Spill",
			Annoyance: 26,
		},
		&Haha{
			Situation: "Tripped Over Banana",
			Humour:    5,
			Participants: p.Refs(
				&Participant{
					Name: "Sudo",
					Age:  16,
				},
				&Participant{
					Name: "Bash",
					Age:  21,
				},
			),
		},
	)
	p.Root().Seq = goal + 1
	if e := p.Save(); e != nil {
		return e
	}

	n.Publish(p.Root())
	return nil
}
