package state

import (
	"context"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"os"
	"sync"
)

type InitBoardInstance func(ct *skyobject.Container, root *skyobject.Root) (
	*BoardInstance, error)

type BoardInstanceConfig struct {
	Master bool
	PK     cipher.PubKey
	SK     cipher.SecKey
}

type BoardInstance struct {
	c *BoardInstanceConfig
	l *log.Logger

	flag skyobject.Flag // Used for compiling pack.

	seqMux  sync.RWMutex
	seqChan chan uint64 // Triggered when sequence increments (for output).
	seq     uint64      // Last completed sequence.

	piMux sync.RWMutex
	pi    *PackInstance

	changesChan chan *Changes // Changes to tree (for web socket).

	needUpdateMux sync.RWMutex
	needUpdate    bool
}

func NewBoardInstance(config *BoardInstanceConfig) InitBoardInstance {
	return func(ct *skyobject.Container, root *skyobject.Root) (*BoardInstance, error) {

		// Prepare output.
		bi := &BoardInstance{
			c:           config,
			l:           inform.NewLogger(true, os.Stdout, "INSTANCE:"+config.PK.Hex()),
			seqChan:     make(chan uint64),
			changesChan: make(chan *Changes, 10),
		}

		// Prepare flags.
		bi.flag = skyobject.HashTableIndex | skyobject.EntireTree
		if !bi.c.Master {
			bi.flag |= skyobject.ViewOnly
		}

		// Prepare pack instance.
		pack, e := ct.Unpack(root, bi.flag, ct.CoreRegistry().Types(), config.SK)
		if e != nil {
			return nil, e
		}
		bi.pi, _, e = NewPackInstance(nil, pack)
		if e != nil {
			return nil, e
		}

		// Output.
		return bi, nil
	}
}

func (bi *BoardInstance) compileViews(changes *Changes) error {
	// Broadcast changes.
	for {
		select {
		case bi.changesChan <- changes:
			goto FinishBroadcast
		default:
			// Empty if too full.
			<-bi.changesChan
		}
	}
FinishBroadcast:

	// TODO: Implement - compiles "VIEWS".

	// Success.
	bi.GetPack().Do(func(p *skyobject.Pack, h *ActivityHeader) error {
		bi.SetSeq(p.Root().Seq)
		return nil
	})
	return nil
}

// Update updates the board instance.
func (bi *BoardInstance) Update() InitBoardInstance {
	return func(ct *skyobject.Container, root *skyobject.Root) (*BoardInstance, error) {

		// Prepare new pack instance.
		newPack, e := ct.Unpack(root, bi.flag, ct.CoreRegistry().Types(), bi.c.SK)
		if e != nil {
			return nil, e
		}
		newPI, changes, e := NewPackInstance(bi.GetPack(), newPack)
		if e != nil {
			return nil, e
		}

		// Compile views.
		if e := bi.compileViews(changes); e != nil {
			bi.l.Println("Compilation of views failed with error:", e)
		}

		// Set new pack instance.
		bi.SetPack(newPI)
		return bi, nil
	}
}

// ChangesChan for WebSocket goodness.
func (bi *BoardInstance) ChangesChan() chan *Changes {
	return bi.changesChan
}

func (bi *BoardInstance) UpdateNeeded() bool {
	bi.needUpdateMux.RLock()
	defer bi.needUpdateMux.RUnlock()
	return bi.needUpdate
}

func (bi *BoardInstance) SetUpdateNeeded() {
	bi.needUpdateMux.Lock()
	defer bi.needUpdateMux.Unlock()
	bi.needUpdate = true
}

func (bi *BoardInstance) ClearUpdateNeeded() {
	bi.needUpdateMux.Lock()
	defer bi.needUpdateMux.Unlock()
	bi.needUpdate = false
}

func (bi *BoardInstance) GetPack() *PackInstance {
	bi.piMux.RLock()
	defer bi.piMux.RUnlock()
	return bi.pi
}

func (bi *BoardInstance) SetPack(pi *PackInstance) {
	bi.piMux.Lock()
	defer bi.piMux.Unlock()
	bi.pi = pi
}

func (bi *BoardInstance) GetSeq() uint64 {
	bi.seqMux.RLock()
	defer bi.seqMux.RUnlock()
	return bi.seq
}

func (bi *BoardInstance) SetSeq(seq uint64) {
	bi.seqMux.Lock()
	defer bi.seqMux.Unlock()
	if seq > bi.seq {
		bi.seq = seq
		for {
			select {
			case bi.seqChan <- seq:
			default:
				return
			}
		}
	}
}

func (bi *BoardInstance) WaitSeq(ctx context.Context, goal uint64) error {
	if bi.GetSeq() >= goal {
		return nil
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case seq := <-bi.seqChan:
			if seq >= goal {
				return nil
			}
		}
	}
}
