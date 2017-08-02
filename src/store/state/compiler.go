package state

import (
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"os"
	"sync"
	"context"
	"time"
	"github.com/skycoin/bbs/src/store/state/states"
)

type CompilerConfig struct {
	Workers *int
}

// Compiler compiles board states.
type Compiler struct {
	user      cipher.PubKey
	c         *CompilerConfig
	l         *log.Logger
	mux       sync.Mutex
	newBState states.NewState
	bStates   map[cipher.PubKey]states.State
	workers   chan func()
	quit      chan struct{}
	wg        sync.WaitGroup
}

func NewCompiler(config *CompilerConfig, options ...Option) *Compiler {
	compiler := &Compiler{
		c:       config,
		l:       inform.NewLogger(true, os.Stdout, "COMPILER"),
		bStates: make(map[cipher.PubKey]states.State),
		workers: make(chan func()),
		quit:    make(chan struct{}),
	}
	for _, option := range options {
		if e := option(compiler); e != nil {
			compiler.l.Fatal(e)
		}
	}
	if compiler.newBState == nil {
		compiler.l.Fatal("newBState not set")
	}
	return compiler
}

// Open opens the compiler.
func (c *Compiler) Open(user cipher.PubKey) {
	c.l.Printf("Opening for user %s ...", user.Hex())
	c.user = user
	for i := 0; i < *c.c.Workers; i++ {
		go c.workerLoop()
	}
}

// Close closes this Compiler.
func (c *Compiler) Close() {
	c.mux.Lock()
	defer c.mux.Unlock()
	for _, bs := range c.bStates {
		bs.Close()
	}
	c.bStates = make(map[cipher.PubKey]states.State)

	for {
		select {
		case c.quit <- struct{}{}:
		default:
			c.wg.Wait()
			return
		}
	}
}

func (c *Compiler) workerLoop() {
	c.wg.Add(1)
	defer c.wg.Done()
	for {
		select {
		case worker := <-c.workers:
			worker()
		case <-c.quit:
			return
		}
	}
}

func (c *Compiler) Trigger(root *node.Root) {
	ctx, _ := context.WithTimeout(context.Background(), time.Second * 10)
	c.getBoardState(root.Pub()).Trigger(ctx, root)
}

func (c *Compiler) DeleteBoard(bpk cipher.PubKey) {
	c.deleteBoardState(bpk)
}

func (c *Compiler) GetBoard(bpk cipher.PubKey) states.State {
	return c.getBoardState(bpk)
}

func (c *Compiler) getBoardState(bpk cipher.PubKey) states.State {
	c.mux.Lock()
	defer c.mux.Unlock()
	bs, has := c.bStates[bpk]
	if !has {
		c.bStates[bpk] = c.newBState(bpk, c.user, c.workers)
		bs = c.bStates[bpk]
	}
	return bs
}

func (c *Compiler) deleteBoardState(bpk cipher.PubKey) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if bs, has := c.bStates[bpk]; has {
		bs.Close()
		delete(c.bStates, bpk)
	}
}
