package state

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/store/state/views"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"os"
	"sync"
	"time"
)

const (
	LogPrefix = "COMPILER"
)

type CompilerConfig struct {
	UpdateInterval *int // In seconds.
}

type Compiler struct {
	c *CompilerConfig
	l *log.Logger

	node *node.Node

	mux    sync.Mutex
	boards map[cipher.PubKey]*BoardInstance
	adders []views.Adder

	newRoots chan *skyobject.Root
	quit     chan struct{}
	wg       sync.WaitGroup
}

func NewCompiler(config *CompilerConfig, newRoots chan *skyobject.Root, node *node.Node, adders ...views.Adder) *Compiler {
	compiler := &Compiler{
		c:        config,
		l:        inform.NewLogger(true, os.Stdout, LogPrefix),
		node:     node,
		boards:   make(map[cipher.PubKey]*BoardInstance),
		adders:   adders,
		newRoots: newRoots,
		quit:     make(chan struct{}),
	}
	go compiler.updateLoop()
	return compiler
}

func (c *Compiler) Close() {
	for {
		select {
		case c.quit <- struct{}{}:
		default:
			c.wg.Wait()
			return
		}
	}
}

// Only for master boards.
func (c *Compiler) updateLoop() {
	c.wg.Add(1)
	defer c.wg.Done()

	ticker := time.NewTicker(time.Second * time.Duration(*c.c.UpdateInterval))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.doMasterUpdate()

		case root := <-c.newRoots:
			c.doRemoteUpdate(root)

		case <-c.quit:
			for _, bi := range c.boards {
				bi.Close()
			}
			return
		}
	}
}

func (c *Compiler) doMasterUpdate() {
	c.mux.Lock()
	defer c.mux.Unlock()

	for _, bi := range c.boards {
		if bi.IsMaster() && bi.UpdateNeeded() {
			c.l.Println("Update needed for master board:", bi.c.PK.Hex())
			if e := bi.Update(c.node, nil); e != nil {
				c.l.Println("Error on update instance:", e)
			}
		}
	}
}

func (c *Compiler) doRemoteUpdate(root *skyobject.Root) {
	bi, e := c.GetBoard(root.Pub)
	if e != nil {
		c.l.Println("Received root error:", e)
		if e := c.InitBoard(root.Pub); e != nil {
			c.l.Println("Init board error:", e)
			return
		}
		bi, e = c.GetBoard(root.Pub)
		if e != nil {
			c.l.Println("Failed to obtain board after init:", e)
			return
		}
	}

	if e := bi.Update(c.node, root); e != nil {
		c.l.Println("Update board instance error:",
			e.Error())
	}
}

func (c *Compiler) InitBoard(pk cipher.PubKey, sk ...cipher.SecKey) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	root, e := c.node.Container().LastFull(pk)
	if e != nil {
		return e
	}

	switch len(sk) {
	case 0:
		bi, e := NewBoardInstance(
			&BoardInstanceConfig{Master: false, PK: pk},
			c.node.Container(), root, c.adders...,
		)
		if e != nil {
			return e
		}
		c.boards[pk] = bi

	case 1:
		bi, e := NewBoardInstance(
			&BoardInstanceConfig{Master: true, PK: pk, SK: sk[0]},
			c.node.Container(), root, c.adders...,
		)
		if e != nil {
			return e
		}
		c.boards[pk] = bi

	default:
		return boo.Newf(boo.Internal,
			"invalid secret key count provided of %d", len(sk))
	}
	return nil
}

func (c *Compiler) DeleteBoard(bpk cipher.PubKey) {
	c.mux.Lock()
	defer c.mux.Unlock()

	bi, has := c.boards[bpk]
	if !has {
		return
	}

	bi.Close()
	delete(c.boards, bpk)
}

func (c *Compiler) GetBoard(pk cipher.PubKey) (*BoardInstance, error) {
	c.mux.Lock()
	defer c.mux.Unlock()
	bi, ok := c.boards[pk]
	if !ok {
		return nil, boo.Newf(boo.NotFound,
			"board of public key '%s' is not found in compiler", pk.Hex())
	}
	return bi, nil
}

func (c *Compiler) GetBoardForce(pk cipher.PubKey) (*BoardInstance, error) {
	bi, e := c.GetBoard(pk)
	if e != nil {
		if e := c.InitBoard(pk); e != nil {
			return nil, e
		}
		return c.GetBoard(pk)
	} else {
		return bi, nil
	}
}
