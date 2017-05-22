package cxo

import (
	"github.com/skycoin/cxo/skyobject"
	"fmt"
)

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

type ThreadPageCatcher struct {
	Index int
	Data  []byte
}

func (c *ThreadPageCatcher) ViaThreadRef(tRef skyobject.Reference) func(v *skyobject.Value) bool {
	return func(v *skyobject.Value) bool {
		c.Index += 1
		fmt.Println("Index:", c.Index)
		vRef, _ := v.Static()
		fmt.Println("Ref:", vRef.String())
		tVal, e := v.FieldByName("Thread")
		if e != nil {
			return false
		}
		ref, e := tVal.Static()
		if e != nil {
			return false
		}
		if ref == tRef {
			fmt.Println("/t Gotcha!")
			c.Data = v.Data()
			return true
		}
		return false
	}
}