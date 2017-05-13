package cxo

import (
	"github.com/evanlinjin/bbs/typ"
	"github.com/skycoin/cxo/skyobject"
)

// FindBoardContainer finds a board container.
func FindBoardContainer(_ *skyobject.Value) bool {
	return true
}

func MakeFindThread(thread *typ.Thread) func(*skyobject.Value) bool {
	return func(v *skyobject.Value) bool {
		fN, _ := v.FieldByName("Name")
		fD, _ := v.FieldByName("Desc")
		n, _ := fN.String()
		d, _ := fD.String()
		return n == thread.Name && d == thread.Desc
	}
}
