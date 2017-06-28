package cxo

import (
	"github.com/skycoin/bbs/src/store/typ"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
)

func makeBoardContainerFinder(r *node.Root) func(_ int, dRef skyobject.Dynamic) bool {
	return func(_ int, dRef skyobject.Dynamic) bool {
		schema, e := r.SchemaByReference(dRef.Schema)
		if e != nil {
			return false
		}
		return schema.Name() == "BoardContainer"
	}
}

func makeThreadVotesContainerFinder(r *node.Root) func(_ int, dRef skyobject.Dynamic) bool {
	return func(_ int, dRef skyobject.Dynamic) bool {
		schema, e := r.SchemaByReference(dRef.Schema)
		if e != nil {
			return false
		}
		return schema.Name() == "ThreadVotesContainer"
	}
}

func makePostVotesContainerFinder(r *node.Root) func(_ int, dRef skyobject.Dynamic) bool {
	return func(_ int, dRef skyobject.Dynamic) bool {
		schema, e := r.SchemaByReference(dRef.Schema)
		if e != nil {
			return false
		}
		return schema.Name() == "PostVotesContainer"
	}
}

func makeThreadPageFinder(w *node.RootWalker, tRef skyobject.Reference) func(i int, ref skyobject.Reference) bool {
	return func(_ int, ref skyobject.Reference) bool {
		threadPage := &typ.ThreadPage{}
		w.DeserializeFromRef(ref, threadPage)
		return threadPage.Thread == tRef
	}
}

type ThreadPageCatcher struct {
	Index int
	Data  []byte
}

func (c *ThreadPageCatcher) ViaThreadRef(w *node.RootWalker, tRef skyobject.Reference) func(int, skyobject.Reference) bool {
	return func(i int, tpRef skyobject.Reference) bool {
		threadPage := &typ.ThreadPage{}
		w.DeserializeFromRef(tpRef, threadPage)
		if threadPage.Thread == tRef {
			c.Index = i
			c.Data, _ = w.Root().Get(tpRef)
			return true
		}
		return false
	}
}
