package state

import (
	"context"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/object/revisions/r0"
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

// CompilerConfig configure the Compiler.
type CompilerConfig struct {
	UpdateInterval *int // In seconds.
}

// Compiler compiles views for boards.
type Compiler struct {
	c *CompilerConfig
	l *log.Logger

	node *node.Node
	file *object.CXOFileManager

	mux    sync.Mutex
	boards map[cipher.PubKey]*BoardInstance
	adders []views.Adder

	newRoots chan RootWrap
	quit     chan struct{}
	wg       sync.WaitGroup
}

// NewCompiler creates a new compiler.
func NewCompiler(
	config *CompilerConfig,
	file *object.CXOFileManager,
	newRoots chan RootWrap,
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

// Close closes the compiler.
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
			c.publishAllMasters()

		case rootWrap := <-c.newRoots:
			c.updateSingle(rootWrap.Root)
			select {
			case rootWrap.Done <- struct{}{}:
			default:
			}

		case <-c.quit:
			for _, bi := range c.boards {
				bi.Close()
			}
			return
		}
	}
}

func (c *Compiler) publishAllMasters() {
	c.file.RangeMasterSubs(func(pk cipher.PubKey, sk cipher.SecKey) {
		bi := c.ensureBoard(pk)

		if e := bi.PublishChanges(); e != nil {
			c.l.Printf(" - [%s] Publish failed with error: %v", pk.Hex()[:5]+"...", e)
		}
	})
}

func (c *Compiler) updateSingle(root *skyobject.Root) {

	isRemote := c.file.HasRemoteSub(root.Pub)
	sk, isMaster := c.file.GetMasterSubSecKey(root.Pub)

	if !isRemote && !isMaster {
		return
	}

	bi := c.ensureBoard(root.Pub)

	if root.IsFull == false {
		c.l.Printf("received root '%s' is not full, returning.", root.Pub.Hex()[:5]+"...")
		return
	}

	if isMaster && bi.needPublish.Value() == true {
		return
	}

	c.l.Printf("compiling '%s' : remote(%v) master(%v)", root.Pub.Hex()[:5]+"...", isRemote, isMaster)
	bi.UpdateWithReceived(root, sk)
}

// EnsureSubmissionKeys ranges through masters and ensures that their specified
// submission public keys are as specified.
func (c *Compiler) EnsureSubmissionKeys(keys []cipher.PubKey) error {
	return c.file.RangeMasterSubs(func(pk cipher.PubKey, sk cipher.SecKey) {
		if bi, e := c.GetBoard(pk); e != nil {
			c.l.Println(e)
		} else {
			if _, e := bi.EnsureSubmissionKeys(keys); e != nil {
				c.l.Println(e)
			}
		}
	})
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

	switch {
	case ok == false:
		return nil, boo.Newf(boo.NotFound,
			"board '%s' not found", pk.Hex()[:5]+"...")

	case bi.IsReceived() == false:
		return nil, boo.Newf(boo.NotFound,
			"board '%s' has not been received", pk.Hex()[:5]+"...")
	}

	return bi, nil
}

func (c *Compiler) UpdateBoard(root *skyobject.Root) {
	c.newRoots <- RootWrap{Root: root}
}

func (c *Compiler) UpdateBoardWithContext(ctx context.Context, root *skyobject.Root) error {
	done := make(chan struct{})
	c.newRoots <- RootWrap{Root: root, Done: done}
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Compiler) GetMasterSummaries() []*r0.BoardSummaryWrap {
	var out []*r0.BoardSummaryWrap
	c.file.RangeMasterSubs(func(pk cipher.PubKey, sk cipher.SecKey) {
		if bi, e := c.GetBoard(pk); e != nil {
			c.l.Println(e)
		} else {
			summary, e := bi.GetSummary(pk, sk)
			if e != nil {
				c.l.Println(e)
				return
			}
			out = append(out, summary)
		}
	})
	return out
}

func (c *Compiler) RangeMasterSubs(action object.MasterSubAction) error {
	return c.file.RangeMasterSubs(action)
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
	bi.SetReceived()
	return bi
}
