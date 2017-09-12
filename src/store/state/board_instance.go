package state

import (
	"context"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/misc/typ"
	"github.com/skycoin/bbs/src/store/state/pack"
	"github.com/skycoin/bbs/src/store/state/views"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"os"
	"sync"
	"time"
)

var (
	// ErrInstanceNotInitialized occurs when instance is not initialized.
	ErrInstanceNotInitialized = boo.New(boo.NotAllowed, "instance not initialized")

	// ErrNotEditable occurs when pack has flag 'ViewOnly' set.
	ErrNotEditable = boo.New(boo.NotAllowed, "cannot edit root")
)

type BoardInstance struct {
	l *log.Logger

	mux sync.RWMutex // Only use (RLock/RUnlock) with reading root sequence.
	n   *node.Node
	p   *skyobject.Pack
	h   *pack.Headers
	v   map[string]views.View

	needPublish typ.Bool // Whether there are changes that need to be published.
	isReceived  typ.Bool // Whether we have received this root.
	isReady     typ.Bool // Whether we have received a full root.
}

// Init initiates the  the board instance.
func (bi *BoardInstance) Init(n *node.Node, pk cipher.PubKey, adders ...views.Adder) *BoardInstance {
	bi.l = inform.NewLogger(true, os.Stdout, "INSTANCE:"+pk.Hex()[:5]+"...")
	bi.n = n
	bi.v = make(map[string]views.View)

	for _, adder := range adders {
		views.Add(bi.v, adder)
	}

	return bi
}

// Close closes the board instance.
func (bi *BoardInstance) Close() {
	bi.mux.Lock()
	defer bi.mux.Unlock()

	if bi.p != nil {
		bi.p.Close()
	}
}

// UpdateWithReceived updates pack header and views to reflect latest sequence of received root.
func (bi *BoardInstance) UpdateWithReceived(r *skyobject.Root, sk cipher.SecKey) error {
	bi.mux.Lock()
	defer bi.mux.Unlock()

	bi.isReceived.Set()
	bi.isReady.Set()

	var (
		master   = sk != cipher.SecKey{}    // Whether this node owns the board.
		ct       = bi.n.Container()         // CXO container.
		pFlags   = skyobject.HashTableIndex // Flags for unpacking root.
		firstRun = false                    // Is the first time running update.
	)

	// Preparation.
	if bi.p != nil {
		if bi.p.Root().Seq >= r.Seq {
			return nil
		}
		bi.p.Close()
	} else {
		firstRun = true
	}

	// Update pack, headers and views.
	{
		if !master {
			pFlags |= skyobject.ViewOnly
		}

		var e error
		if bi.p, e = ct.Unpack(r, pFlags, ct.CoreRegistry().Types(), sk); e != nil {
			return e
		}

		if bi.h, e = pack.NewHeaders(bi.h, bi.p); e != nil {
			return e
		}

		if firstRun {
			i := 1
			for name, view := range bi.v {
				if e := view.Init(bi.p, bi.h); e != nil {
					return boo.WrapType(e, boo.Internal, "failed to generate view")
				}
				bi.l.Printf("(%d/%d) Loaded '%s' view.", i, len(bi.v), name)
				i++
			}
		} else {
			i := 1
			for name, view := range bi.v {
				if e := view.Update(bi.p, bi.h); e != nil {
					return boo.WrapType(e, boo.Internal, "failed to update view")
				}
				bi.l.Printf("(%d/%d) Updated '%s' view.", i, len(bi.v), name)
				i++
			}
		}

		// TODO: Broadcast changes.
	}

	return nil
}

// PublishChanges publishes changes to CXO.
// (only if instance is initialised, changes were made, and the node owns the board.)
func (bi *BoardInstance) PublishChanges() error {
	defer bi.needPublish.Clear()

	if bi.needPublish.Value() == false {
		return nil
	}

	bi.mux.Lock()
	defer bi.mux.Unlock()

	if bi.p == nil || bi.p.Flags()&skyobject.ViewOnly > 0 {
		return nil
	}

	// Update CXO.
	if e := bi.p.Save(); e != nil {
		return e
	}
	bi.n.Publish(bi.p.Root())

	// Update headers.
	var e error
	if bi.h, e = pack.NewHeaders(bi.h, bi.p); e != nil {
		return e
	}

	// Update views.
	for _, view := range bi.v {
		if e := view.Update(bi.p, bi.h); e != nil {
			return boo.WrapType(e, boo.Internal, "failed to update view")
		}
	}

	// TODO: Broadcast changes.

	return nil
}

// Get obtains data from views.
func (bi *BoardInstance) Get(viewID, cmdID string, a ...interface{}) (interface{}, error) {
	bi.mux.Lock()
	defer bi.mux.Unlock()

	if bi.v == nil {
		return nil, ErrInstanceNotInitialized
	}

	view, has := bi.v[viewID]
	if !has {
		return nil, boo.Newf(boo.NotFound, "view of id '%s' is not found", viewID)
	}

	return view.Get(cmdID, a...)
}

// IsMaster determines if we are master.
func (bi *BoardInstance) IsMaster() bool {
	bi.mux.RLock()
	defer bi.mux.RUnlock()
	return bi.p != nil && bi.p.Flags()&skyobject.ViewOnly == 0
}

// GetSeq obtains the current sequence.
func (bi *BoardInstance) GetSeq() uint64 {
	bi.mux.RLock()
	defer bi.mux.RUnlock()

	if bi.p != nil {
		return bi.p.Root().Seq
	}
	return uint64(0)
}

// WaitSeq waits until sequence reaches or surpassed the goal.
func (bi *BoardInstance) WaitSeq(ctx context.Context, goal uint64) error {
	if bi.p == nil {
		return ErrInstanceNotInitialized
	} else if bi.GetSeq() >= goal {
		return nil
	}

	ctx, _ = context.WithTimeout(ctx, time.Second*30)

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

// PackAction represents an action applied to a root pack.
type PackAction func(p *skyobject.Pack, h *pack.Headers) error

// EditPack ensures safe modification to the pack.
func (bi *BoardInstance) EditPack(action PackAction) error {
	bi.mux.Lock()
	defer bi.mux.Unlock()

	if bi.p == nil {
		return ErrInstanceNotInitialized
	} else if bi.p.Flags()&skyobject.ViewOnly > 0 {
		return ErrNotEditable
	}

	bi.needPublish.Set()
	if e := action(bi.p, bi.h); e != nil {
		return e
	}
	return nil
}

// SetReceived set's the board as being received (however, not necessarily ready).
func (bi *BoardInstance) SetReceived() {
	bi.isReceived.Set()
}

// IsReceived determines whether board has been received.
func (bi *BoardInstance) IsReceived() bool {
	v := bi.isReceived.Value()
	//bi.l.Println("IsReceived:", v)
	return v
}

// IsReady determines whether board is received and ready.
func (bi *BoardInstance) IsReady() bool {
	v := bi.isReady.Value()
	//bi.l.Println("IsReady:", v)
	return v
}
