package views

import (
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state/pack"
	"github.com/skycoin/bbs/src/store/state/views/content_view"
	"github.com/skycoin/cxo/skyobject"
	"sync"
)

type View interface {

	// Init initiates the view.
	Init(pack *skyobject.Pack, headers *pack.Headers, mux *sync.Mutex) error

	// Update updates the view.
	Update(pack *skyobject.Pack, headers *pack.Headers, mux *sync.Mutex) error

	// Get obtains information from the view.
	Get(id string, a ...interface{}) (object.Lockable, error)
}

type Adder func() (string, View)

func Add(viewsMap map[string]View, add Adder) {
	id, view := add()
	viewsMap[id] = view
}

const (
	Content = "content"
)

func AddContent() Adder {
	return func() (string, View) {
		return Content, new(content_view.ContentView)
	}
}
