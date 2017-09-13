package state

import "github.com/skycoin/cxo/skyobject"

type RootWrap struct {
	Done chan struct{}
	Root *skyobject.Root
}
