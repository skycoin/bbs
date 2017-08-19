package state

import (
	"context"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"os"
	"sync"
	"time"
	"github.com/skycoin/cxo/node"
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

	piMux sync.RWMutex
	pi    *PackInstance

	changesChan chan *object.Changes // Changes to tree (for output - web socket).

	needUpdateMux sync.RWMutex
	needUpdate    bool
}

func NewBoardInstance(config *BoardInstanceConfig, ct *skyobject.Container, root *skyobject.Root) (
	*BoardInstance, error,
) {
	// Prepare output.
	bi := &BoardInstance{
		c:           config,
		l:           inform.NewLogger(true, os.Stdout, "INSTANCE:"+config.PK.Hex()),
		changesChan: make(chan *object.Changes, 10),
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
	bi.pi, e = NewPackInstance(nil, pack)
	if e != nil {
		return nil, e
	}

	// Output.
	return bi, nil
}

// Update updates the board instance. (External trigger).
func (bi *BoardInstance) Update(node *node.Node, root *skyobject.Root) error {
	return bi.SetPack(func(oldPI *PackInstance) (*PackInstance, error) {

		// If master of board. First update last changes.
		if bi.c.Master && bi.UpdateNeeded() {
			root = oldPI.pack.Root()
			node.Publish(root)
		}

		// Prepare new pack instance.
		ct := node.Container()
		newPack, e := ct.Unpack(root, bi.flag, ct.CoreRegistry().Types(), bi.c.SK)
		if e != nil {
			return nil, e
		}
		newPI, e := NewPackInstance(oldPI, newPack)
		if e != nil {
			return nil, e
		}

		// Broadcast changes.
		changes := newPI.headers.GetChanges()
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
		// Set new pack instance.
		return newPI, nil
	})
}

// ChangesChan for WebSocket goodness.
func (bi *BoardInstance) ChangesChan() chan *object.Changes {
	return bi.changesChan
}

/*
	<<< Update? >>>
	>>> Whether a call to (*BoardInstance).Update() is needed.
*/

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

/*
	<<< PACK INSTANCE >>>
*/

func (bi *BoardInstance) SetPack(set func(oldPI *PackInstance) (*PackInstance, error)) error {
	bi.piMux.Lock()
	defer bi.piMux.Unlock()
	if pi, e := set(bi.pi); e != nil {
		return e
	} else {
		bi.pi = pi
		return nil
	}
}

func (bi *BoardInstance) EditPack(edit func(pi *PackInstance) error) error {
	bi.piMux.Lock()
	defer bi.piMux.Unlock()
	return edit(bi.pi)
}

func (bi *BoardInstance) ReadPack(read func(pi *PackInstance) error) error {
	bi.piMux.RLock()
	defer bi.piMux.RUnlock()
	return read(bi.pi)
}

func (bi *BoardInstance) GetSeq() uint64 {
	var seq uint64
	bi.ReadPack(func(pi *PackInstance) error {
		seq = pi.GetSeq()
		return nil
	})
	return seq
}

func (bi *BoardInstance) WaitSeq(ctx context.Context, goal uint64) error {
	if bi.GetSeq() >= goal {
		return nil
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if bi.GetSeq() >= goal {
				return nil
			}
		}
	}
}
