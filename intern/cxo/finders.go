package cxo

import "github.com/skycoin/cxo/skyobject"

func makeBoardContainerFinder() func(v *skyobject.Value) bool {
	return func(v *skyobject.Value) bool {
		return v.Schema().Name() == "BoardContainer"
	}
}

func makeThreadPageFinder(tRef skyobject.Reference) func(v *skyobject.Value) bool {
	return func(v *skyobject.Value) bool {
		tVal, e := v.FieldByName("Thread")
		if e != nil {
			return false
		}
		ref, e := tVal.Static()
		if e != nil {
			return false
		}
		return ref == tRef
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
