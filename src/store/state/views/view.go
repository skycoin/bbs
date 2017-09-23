package views

import (
	"github.com/skycoin/bbs/src/store/state/pack"
	"github.com/skycoin/bbs/src/store/state/views/content_view"
	"github.com/skycoin/bbs/src/store/state/views/follow_view"
	"github.com/skycoin/cxo/skyobject"
)

type View interface {

	// Init initiates the view.
	Init(pack *skyobject.Pack, headers *pack.Headers) error

	// Update updates the view.
	Update(pack *skyobject.Pack, headers *pack.Headers) error

	// Get obtains information from the view.
	Get(id string, a ...interface{}) (interface{}, error)
}

type Adder func() (string, View)

func Add(viewsMap map[string]View, add Adder) {
	id, view := add()
	viewsMap[id] = view
}

const (
	Content = "content"
	Follow  = "follow"
)

func AddContent() Adder {
	return func() (string, View) {
		return Content, new(content_view.ContentView)
	}
}

func AddFollow() Adder {
	return func() (string, View) {
		return Follow, new(follow_view.FollowView)
	}
}
