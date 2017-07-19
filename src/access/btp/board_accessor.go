package btp

import (
	"fmt"
	"github.com/skycoin/bbs/src/access"
	"github.com/skycoin/bbs/src/boo"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/bbs/src/store"
	"path"
	"sync"
	"time"
)

const (
	BoardsFileName    = "boards.json"
	SaveCheckDuration = 10 * time.Second
)

// BoardAccessorConfig configures a BoardAccessor.
type BoardAccessorConfig struct {
	MemoryMode *bool   // Whether to use local storage on runtime.
	ConfigDir  *string // Configuration directory.
}

// BoardAccessor controls access to a board.
type BoardAccessor struct {
	mux        sync.Mutex
	wg         sync.WaitGroup
	config     *BoardAccessorConfig
	cxo        *store.CXO
	stateSaver *access.StateSaver
	bFile      *BoardsFile
	quit       chan struct{}
}

// NewBoardAccessor creates a new BoardAccessor.
func NewBoardAccessor(config *BoardAccessorConfig, cxo *store.CXO, stateSaver *access.StateSaver) *BoardAccessor {
	boardAccessor := &BoardAccessor{
		config:     config,
		cxo:        cxo,
		stateSaver: stateSaver,
		bFile:      NewBoardsFile(),
		quit:       make(chan struct{}),
	}
	go boardAccessor.service()
	return boardAccessor
}

// Close closes the BoardAccessor.
func (a *BoardAccessor) Close() {
	go func() { a.quit <- struct{}{} }()
	a.wg.Wait()
}

func (a *BoardAccessor) lock() func() {
	a.mux.Lock()
	return a.mux.Unlock
}

func (a *BoardAccessor) service() {
	a.wg.Add(1)

	ticker := time.NewTicker(SaveCheckDuration)
	boardsFileLoc := path.Join(*a.config.ConfigDir, BoardsFileName)

	if e := a.prepareBoardsFile(&boardsFileLoc); e != nil {
		panic(e)
	}

	for {
		select {
		case <-ticker.C:
			a.saveBoardsFile(&boardsFileLoc)

		case <-a.quit:
			a.saveBoardsFile(&boardsFileLoc)
			a.wg.Done()
			return
		}
	}
}

func (a *BoardAccessor) prepareBoardsFile(loc *string) error {
	defer a.lock()()
	a.bFile.Load(*loc)
	for pkStr, bInfo := range a.bFile.Boards {
		pk, e := misc.GetPubKey(pkStr)
		if e != nil {
			return boo.Newf(boo.InvalidRead,
				"invalid public key '%s' from file '%s'", pkStr, *loc)
		}
		if e := a.cxo.Subscribe(bInfo.Connection, pk); e != nil {
			fmt.Printf("Failed to subscribe to '[%s] (%s)'\n", bInfo.Connection, pk.Hex())
		}
		a.stateSaver.Update(a.cxo.GetRoot(pk))
	}
	for pkStr := range a.bFile.MasterBoards {
		pk, e := misc.GetPubKey(pkStr)
		if e != nil {
			return boo.Newf(boo.InvalidRead,
				"invalid public key '%s' from file '%s'", pkStr, *loc)
		}
		if e := a.cxo.Subscribe("", pk); e != nil {
			fmt.Printf("Failed to subscribe to '[localhost] (%s)'\n", pk.Hex())
		}
		a.stateSaver.Update(a.cxo.GetRoot(pk))
	}
	return nil
}

func (a *BoardAccessor) saveBoardsFile(loc *string) {
	if !a.bFile.unsaved {
		return
	}
	defer a.lock()()
	a.bFile.Save(*loc)
}
