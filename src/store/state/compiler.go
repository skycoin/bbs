package state

import (
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"os"
	"sync"
)

type CompilerConfig struct {
	Workers *int
}

// Compiler compiles board states.
type Compiler struct {
	user    cipher.PubKey
	c       *CompilerConfig
	l       *log.Logger
	mux     sync.Mutex
	boards  map[cipher.PubKey]*BoardState
	workers chan func()
	quit    chan struct{}
	wg      sync.WaitGroup
}

func NewCompiler(config *CompilerConfig) *Compiler {
	compiler := &Compiler{
		c:       config,
		l:       inform.NewLogger(true, os.Stdout, "COMPILER"),
		boards:  make(map[cipher.PubKey]*BoardState),
		workers: make(chan func()),
		quit:    make(chan struct{}),
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
	for _, bs := range c.boards {
		bs.Close()
	}
	c.boards = make(map[cipher.PubKey]*BoardState)

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
	c.getBoardState(root.Pub()).
		newRoots <- root
}

func (c *Compiler) DeleteBoard(bpk cipher.PubKey) {
	c.deleteBoardState(bpk)
}

func (c *Compiler) GetBoard(bpk cipher.PubKey) *BoardState {
	return c.getBoardState(bpk)
}

func (c *Compiler) getBoardState(bpk cipher.PubKey) *BoardState {
	c.mux.Lock()
	defer c.mux.Unlock()
	bs, has := c.boards[bpk]
	if !has {
		c.boards[bpk] = NewBoardState(bpk, c.user, c.workers)
		bs = c.boards[bpk]
	}
	return bs
}

func (c *Compiler) deleteBoardState(bpk cipher.PubKey) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if bs, has := c.boards[bpk]; has {
		bs.Close()
		delete(c.boards, bpk)
	}
}
