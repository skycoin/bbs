package views

import (
	"github.com/skycoin/cxo/skyobject"
	"sync"
)

type View interface {

	// Init initiates the view.
	Init(pack *skyobject.Pack, mux *sync.Mutex) error

	// Update updates the view.
	Update(oldPack, newPack *skyobject.Pack, oldMux, newMux *sync.Mutex) error

	// Get obtains information from the view.
	Get(id string, in interface{}) (interface{}, error)
}
