package content_view

import (
	"github.com/skycoin/cxo/skyobject"
	"sync"
)

type ContentView struct {
}

func (v *ContentView) Init(pack *skyobject.Pack, mux *sync.Mutex) error {
	return nil
}

func (v *ContentView) Update(oldPack, newPack *skyobject.Pack, oldMux, newMux *sync.Mutex) error {
	return nil
}

func (v *ContentView) Get(id string, in interface{}) (interface{}, error) {
	return nil, nil
}
