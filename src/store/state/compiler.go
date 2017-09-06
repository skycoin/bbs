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
		bi := c.ensureBoard(pk)

		r, e := c.node.Container().LastRoot(pk)
		if e != nil {
			c.l.Printf(" - [%s] LastRoot failed with error : %v", pk.Hex()[:5]+"...", e)
		}

		if e := bi.UpdateWithReceived(r, sk); e != nil {
			c.l.Printf(" - [%s] Update failed with error: %v", pk.Hex()[:5]+"...", e)
		}
	})
}

func (c *Compiler) doRemoteUpdate(root *skyobject.Root) {

	isRemote := c.file.HasRemoteSub(root.Pub)
	sk, isMaster := c.file.GetMasterSubSecKey(root.Pub)

	if !isRemote && !isMaster {
		return
	}

	c.ensureBoard(root.Pub).UpdateWithReceived(root, sk)
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
		return nil, boo.Newf(boo.NotFound, "board '%s' not found", pk.Hex()[:5]+"...")
	}
	return bi, nil
}

func (c *Compiler) NewRootsChan() chan *skyobject.Root {
	return c.newRoots
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func (c *Compiler) ensureBoard(pk cipher.PubKey) *BoardInstance {
	c.mux.Lock()
	defer c.mux.Unlock()

	bi, has := c.boards[pk]
	if !has {
		bi = new(BoardInstance).Init(c.node, pk, c.adders...)
		c.boards[pk] = bi
	}
	return bi
}
