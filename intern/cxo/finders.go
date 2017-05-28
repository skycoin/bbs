package cxo

import (
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/cxo/node"
	"github.com/evanlinjin/bbs/intern/typ"
)

func makeBoardContainerFinder(r *node.Root) func(_ int, dRef skyobject.Dynamic) bool {
	return func(i int, dRef skyobject.Dynamic) bool {
		schema, e := r.SchemaByReference(dRef.Schema)
		if e != nil {
			return false
		}
		return schema.Name() == "BoardContainer"
	}
}

func makeThreadPageFinder(w *node.RootWalker, tRef skyobject.Reference) func(i int, ref skyobject.Reference) bool {
	return func(_ int, ref skyobject.Reference) bool {
		threadPage := &typ.ThreadPage{}
		w.DeserializeFromRef(ref, threadPage)
		return threadPage.Thread == tRef
	}
}

func makeThreadFinder(tRef skyobject.Reference) func(v *skyobject.Value) bool {
	return func(v *skyobject.Value) bool {
		if r, _ := v.Static(); r == tRef {
			return true
		}
		return false
	}
}

type ThreadPageCatcher struct {
	Index int
	Data  []byte
}

func (c *ThreadPageCatcher) ViaThreadRef(tRef skyobject.Reference) func(int, skyobject.Reference) bool {
	return func(i int, ref skyobject.Reference) bool {
		if ref == tRef {
			c.Index = i
			return true
		}
		return false
	}
}
