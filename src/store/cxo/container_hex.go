package cxo

import (
	"encoding/hex"
	"errors"
	"github.com/skycoin/bbs/src/store/typ"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

// HexObj represents an object stored in cxo as hex.
type HexObj struct {
	Ref  string `json:"ref"`
	Data string `json:"data"`
}

// ThreadPageHex represents a ThreadPage with data as hex.
type ThreadPageHex struct {
	ThreadPage HexObj   `json:"thread_page"`
	Thread     HexObj   `json:"thread"`
	Posts      []HexObj `json:"posts"`
}

// GetThreadPageAsHex retrieves a ThreadPage with data as hex.
func (c *Container) GetThreadPageAsHex(bpk cipher.PubKey, tRef skyobject.Reference) (tph *ThreadPageHex, e error) {
	c.Lock(c.GetThreadPageAsHex)
	defer c.Unlock()

	tph = new(ThreadPageHex)
	w := c.c.LastFullRoot(bpk).Walker()

	bc := &typ.BoardContainer{}
	if e = w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return
	}

	catcher := &ThreadPageCatcher{Index: -1}
	tp := &typ.ThreadPage{}
	if e = w.AdvanceFromRefsField("ThreadPages", tp, catcher.ViaThreadRef(w, tRef)); e != nil {
		return
	}

	// Load in thread page data.
	tph.ThreadPage.Ref = bc.ThreadPages[catcher.Index].String()
	tph.ThreadPage.Data = hex.EncodeToString(catcher.Data)

	// Load in thread data.
	tData, has := c.c.Get(tp.Thread)
	if !has {
		e = errors.New("thread not found")
		return
	}
	tph.Thread.Ref = tRef.String()
	tph.Thread.Data = hex.EncodeToString(tData)

	// Load in boards data.
	tph.Posts = make([]HexObj, len(tp.Posts))
	for i, ref := range tp.Posts {
		pData, has := c.c.Get(ref)
		if !has {
			continue
		}
		tph.Posts[i].Ref = ref.String()
		tph.Posts[i].Data = hex.EncodeToString(pData)
	}
	return
}

// GetThreadPageWithTpRefAsHex retrieves a ThreadPage with data as hex.
// This function uses the reference of a ThreadPage.
func (c *Container) GetThreadPageWithTpRefAsHex(tpRef skyobject.Reference) (tph *ThreadPageHex, e error) {
	c.Lock(c.GetThreadPageWithTpRefAsHex)
	defer c.Unlock()

	tph = new(ThreadPageHex)

	tp := &typ.ThreadPage{}
	tpData, has := c.c.Get(tpRef)
	if !has {
		e = errors.New("thread page not found")
		return
	}
	tp.Deserialize(tpData)
	tph.ThreadPage.Ref = tpRef.String()
	tph.ThreadPage.Data = hex.EncodeToString(tpData)

	tData, has := c.c.Get(tp.Thread)
	if !has {
		e = errors.New("thread not found")
		return
	}
	tph.Thread.Ref = tp.Thread.String()
	tph.Thread.Data = hex.EncodeToString(tData)

	tph.Posts = make([]HexObj, len(tp.Posts))
	for i, ref := range tp.Posts {
		pData, has := c.c.Get(ref)
		if !has {
			continue
		}
		tph.Posts[i].Ref = ref.String()
		tph.Posts[i].Data = hex.EncodeToString(pData)
	}
	return
}

// NewThreadWithHex attempts to creates a new thread in a board with hex data of thread.
func (c *Container) NewThreadWithHex(bpk cipher.PubKey, bsk cipher.SecKey, tData []byte) (e error) {
	c.Lock(c.NewThreadWithHex)
	defer c.Unlock()

	// Obtain thread.
	t := &typ.Thread{}
	if e = t.Deserialize(tData); e != nil {
		return
	}
	t.MasterBoard = bpk.Hex()
	// Save to board.
	w := c.c.LastRootSk(bpk, bsk).Walker()
	bc := &typ.BoardContainer{}
	if e = w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return
	}
	var tRef skyobject.Reference
	if tRef, e = w.AppendToRefsField("Threads", *t); e != nil {
		return e
	}
	_, e = w.AppendToRefsField("ThreadPages", typ.ThreadPage{Thread: tRef})
	//t.Ref = cipher.SHA256(tRef).Hex()
	return
}

// NewPostWithHex attempts to create a new post in a given board and thread with hex data of post.
// The post better be properly signed otherwise other nodes will not accept it.
func (c *Container) NewPostWithHex(bpk cipher.PubKey, bsk cipher.SecKey, tRef skyobject.Reference, pData []byte) error {
	c.Lock(c.NewPostWithHex)
	defer c.Unlock()

	// Obtain post.
	p := &typ.Post{}
	if e := p.Deserialize(pData); e != nil {
		return e
	}
	// Save.
	w := c.c.LastRootSk(bpk, bsk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return e
	}
	tp := &typ.ThreadPage{}
	if e := w.AdvanceFromRefsField("ThreadPages", tp, makeThreadPageFinder(w, tRef)); e != nil {
		return e
	}
	_, e := w.AppendToRefsField("Posts", *p)
	return e
}
