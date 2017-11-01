package state

import (
	"context"
	"encoding/json"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/misc/typ"
	"github.com/skycoin/bbs/src/misc/typ/paginatedtypes"
	"github.com/skycoin/bbs/src/store/object"
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

	//adders []views.Adder // generates views.

	mux sync.RWMutex // Only use (RLock/RUnlock) with reading root sequence.
	n   *node.Node
	p   *skyobject.Pack
	h   *Headers
	v   *Viewer

	needPublish typ.Bool // Whether there are changes that need to be published.
	needReset   typ.Bool // Whether a reset is needed.
	isReceived  typ.Bool // Whether we have received this root.
	isReady     typ.Bool // Whether we have received a full root.
}

// Init initiates the  the board instance.
func (bi *BoardInstance) Init(n *node.Node, pk cipher.PubKey) *BoardInstance {
	bi.l = inform.NewLogger(true, os.Stdout, "INSTANCE:"+pk.Hex()[:5]+"...")
	bi.n = n
	//bi.v

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

	bi.l.Printf("TRIGGERED: UpdateWithReceived()")

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

	if !master {
		pFlags |= skyobject.ViewOnly
	}

	newPack, e := ct.Unpack(r, pFlags, ct.CoreRegistry().Types(), sk)
	if e != nil {
		bi.l.Println(" - root unpack failed with error:", e)
		if newPack, e = bi.fixRoot(firstRun, pFlags, r.Seq, r.Pub, sk); e != nil {
			bi.l.Println("\t- FAILED:", e)
			return e
		} else {
			bi.l.Println("\t- SUCCESS!", newPack.Root().Seq)
			bi.needPublish.Set()
		}
	}

	bi.l.Println(" - root unpack succeeded.")
	bi.p = newPack

	newHeaders, e := NewHeaders(bi.h, bi.p)
	if e != nil {
		bi.l.Println(" - failed to generate new headers:", e)
		return e
	}

	bi.l.Println(" - new headers successfully generated.")
	bi.h = newHeaders

	if firstRun {
		if bi.v, e = NewViewer(bi.p, paginatedtypes.NewSimple); e != nil {
			return e
		}
	} else {
		if e := bi.v.Update(bi.p, bi.h); e != nil {
			return e
		}
	}

	// TODO: Broadcast changes.

	return nil
}

func (bi *BoardInstance) fixRoot(firstRun bool, flags skyobject.Flag, goal uint64, pk cipher.PubKey, sk cipher.SecKey) (*skyobject.Pack, error) {
	var (
		ct        = bi.n.Container()
		isMaster  = sk != (cipher.SecKey{})
		validPack *skyobject.Pack
	)

	// If we don't have old, find it.
	if firstRun == false {
		for i := goal; i >= 0; i-- {
			if tempRoot, e := ct.Root(pk, i); e != nil || len(tempRoot.Refs) != object.RootChildrenCount {
				continue
			} else if tempPack, e := ct.Unpack(tempRoot, flags, ct.CoreRegistry().Types(), sk); e != nil {
				continue
			} else {
				// TODO (evanlinjin) : Need to check root.
				validPack = tempPack
				break
			}
		}
	}
	if validPack == nil {
		return nil, boo.New(boo.InvalidRead,
			"failed to find a valid root that can represent a board")
	}

	// Return if we are unable to change most recent root.
	if isMaster == false {
		return validPack, nil
	}

	// Surpass sequence.
	oldSeq := validPack.Root().Seq
	for i := oldSeq; i < goal; i++ {
		if e := validPack.Save(); e != nil {
			return nil, boo.WrapTypef(e, boo.Internal, "failed to surpass seq(%d)", oldSeq)
		}
	}

	return validPack, nil
}

// PublishChanges publishes changes to CXO.
// Only use if instance is initialised, changes were made, and the node owns the board.
// Should be triggered by compiler based on an interval.
func (bi *BoardInstance) PublishChanges() error {

	if bi.needPublish.Value() == false {
		return nil
	}
	defer bi.needPublish.Clear()

	bi.mux.Lock()
	defer bi.mux.Unlock()

	if bi.p == nil || bi.p.Flags()&skyobject.ViewOnly > 0 {
		return nil
	}

	// Update CXO.
	if e := bi.p.Save(); e != nil {
		return boo.WrapType(e, boo.Internal, "failed to save in cxo db")
	}
	bi.n.Publish(bi.p.Root())

	// Reset header and views if needed.
	if bi.needReset.Value() {

		// Reset headers.
		var e error
		if bi.h, e = NewHeaders(nil, bi.p); e != nil {
			return boo.WrapType(e, boo.Internal, "failed to reset headers")
		}

		// Reset views.
		if bi.v, e = NewViewer(bi.p, paginatedtypes.NewSimple); e != nil {
			return boo.WrapType(e, boo.Internal, "failed to reset view")
		}

		// End the need to reset.
		bi.needReset.Clear()

	} else {

		// Update headers.
		var e error
		if bi.h, e = NewHeaders(bi.h, bi.p); e != nil {
			return boo.WrapType(e, boo.Internal, "failed to generate new headers")
		}

		// Update views.
		if e := bi.v.Update(bi.p, bi.h); e != nil {
			return boo.WrapType(e, boo.Internal, "failed to update view")
		}
	}

	// TODO: Broadcast changes.

	return nil
}

// Viewer obtains the viewer.
func (bi *BoardInstance) Viewer() *Viewer {
	return bi.v
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

// GetSummary returns the board's summary in encoded json and signed with board's public key.
func (bi *BoardInstance) GetSummary(pk cipher.PubKey, sk cipher.SecKey) (*object.BoardSummaryWrap, error) {
	v, e := bi.Viewer().GetBoard()
	if e != nil {
		return nil, e
	}
	raw, e := json.Marshal(v)
	if e != nil {
		return nil, e
	}
	out := &object.BoardSummaryWrap{Raw: raw}
	out.Sign(pk, sk)
	return out, nil
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
type PackAction func(p *skyobject.Pack, h *Headers) error

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

// ViewPack views a pack.
func (bi *BoardInstance) ViewPack(action PackAction) error {
	bi.mux.Lock()
	defer bi.mux.Unlock()

	if bi.p == nil {
		return ErrInstanceNotInitialized
	}

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
	return bi.isReceived.Value()
}

// IsReady determines whether board is received and ready.
func (bi *BoardInstance) IsReady() bool {
	return bi.isReady.Value()
}

func (bi *BoardInstance) Export(pk cipher.PubKey, sk cipher.SecKey) (*object.PagesJSON, error) {
	var out *object.PagesJSON
	var e = bi.ViewPack(func(p *skyobject.Pack, h *Headers) error {
		pages, e := object.GetPages(p, &object.GetPagesIn{
			RootPage:  true,
			BoardPage: true,
			DiffPage:  true,
			UsersPage: true,
		})
		if e != nil {
			return e
		}
		out, e = pages.ToJSON(pk, sk)
		return e
	})
	return out, e
}

func (bi *BoardInstance) Import(in *object.PagesJSON) (uint64, error) {
	var goal uint64
	e := bi.EditPack(func(p *skyobject.Pack, h *Headers) error {
		goal = p.Root().Seq + 1
		pages, e := object.NewPages(p, in)
		if e != nil {
			return e
		}
		return pages.Save(p)
	})
	bi.needReset.Set()
	return goal, e
}
