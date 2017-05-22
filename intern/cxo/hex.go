package cxo

import (
	"encoding/hex"
	"errors"
	"github.com/evanlinjin/bbs/intern/typ"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

type HexObj struct {
	Ref  string `json:"ref"`
	Data string `json:"data"`
}

type ThreadPageHex struct {
	ThreadPage HexObj   `json:"thread_page"`
	Thread     HexObj   `json:"thread"`
	Posts      []HexObj `json:"posts"`
}

func (c *Container) GetThreadPageAsHex(bpk cipher.PubKey, tRef skyobject.Reference) (tph *ThreadPageHex, e error) {
	tph = new(ThreadPageHex)
	w := c.c.LastRoot(bpk).Walker()

	bc := &typ.BoardContainer{}
	if e = w.AdvanceFromRoot(bc, makeBoardContainerFinder()); e != nil {
		return
	}

	catcher := &ThreadPageCatcher{Index: -1}
	tp := &typ.ThreadPage{}
	if e = w.AdvanceFromRefsField("ThreadPages", tp, catcher.ViaThreadRef(tRef)); e != nil {
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

func (c *Container) GetThreadPageWithTpRefAsHex(tpRef skyobject.Reference) (tph *ThreadPageHex, e error) {
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

// NewThreadPageWithHex adds a new thread page to a board using hex data.
func (c *Container) NewThreadPageWithHex(bpk cipher.PubKey, tph *ThreadPageHex) error {
	return nil
}
