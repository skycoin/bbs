package state

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/store/object"
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
	file *object.CXOFileManager

	mux    sync.Mutex
	boards map[cipher.PubKey]*BoardInstance
	adders []views.Adder

	newRoots chan *skyobject.Root
	quit     chan struct{}
	wg       sync.WaitGroup
}

func NewCompiler(
	config *CompilerConfig,
	file *object.CXOFileManager,
	newRoots chan *skyobject.Root,
	node *node.Node,
	adders ...views.Adder,
) *Compiler {
	compiler := &Compiler{
		c:        config,
		l:        inform.NewLogger(true, os.Stdout, LogPrefix),
		node:     node,
		file:     file,
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
	c.file.RangeMasterSubs(func(pk cipher.PubKey, sk cipher.SecKey) {
		c.mux.Lock()
		bi, ok := c.boards[pk]
		c.mux.Unlock()

		if !ok {
			c.l.Printf(" - ['%s'] Initialising in compiler.", pk.Hex()[:5]+"...")
			if e := c.InitBoard(false, pk, sk); e != nil {
				c.l.Println(" - - (ERROR)", e)
				return
			} else {
				c.l.Println(" - - (OKAY)")
				c.mux.Lock()
				bi = c.boards[pk]
				c.mux.Unlock()
			}
		}

		if bi.UpdateNeeded() {
			if e := bi.Update(c.node, nil); e != nil {
				c.l.Printf(" - ['%s'] Update failed with error: %v", pk.Hex()[:5]+"...", e)
				c.DeleteBoard(pk)
				c.l.Println(" - - (RESET) Result:", c.InitBoard(false, pk, sk))
			}
		}
	})
}

func (c *Compiler) doRemoteUpdate(root *skyobject.Root) {
	bi, e := c.GetBoard(root.Pub)
	if e != nil {
		c.l.Println("Board '%s' not compiled.", root.Pub.Hex()[:5]+"...")

		bsk, master := c.file.GetMasterSubSecKey(root.Pub)

		if master {
			if e := c.InitBoard(true, root.Pub, bsk); e != nil {
				c.l.Println("Init board error:", e)
				return
			}
		} else {
			if e := c.InitBoard(true, root.Pub); e != nil {
				c.l.Println("Init board error:", e)
				return
			}
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

func (c *Compiler) InitBoard(checkFile bool, pk cipher.PubKey, sk ...cipher.SecKey) error {
	if checkFile && !(c.file.HasMasterSub(pk) || c.file.HasRemoteSub(pk))  {
		return boo.Newf(boo.NotFound,
			"Not subscribed to feed '%s'", pk.Hex()[:5]+"...")
	}

	c.mux.Lock()
	defer c.mux.Unlock()

	root, e := c.node.Container().LastRoot(pk)
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
	bi, ok := c.boards[pk]
	c.mux.Unlock()
	if !ok {
		c.l.Printf("First time compiling board '%s'", pk.Hex()[:5]+"...")
		if e := c.InitBoard(false, pk); e != nil {
			return nil, e
		}
		return c.GetBoard(pk)
	}
	return bi, nil
}

func (c *Compiler) GetBoardForce(pk cipher.PubKey) (*BoardInstance, error) {
	bi, e := c.GetBoard(pk)
	if e != nil {
		if e := c.InitBoard(false, pk); e != nil {
			return nil, e
		}
		return c.GetBoard(pk)
	} else {
		return bi, nil
	}
}
