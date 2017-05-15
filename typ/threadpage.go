package typ

import "github.com/skycoin/cxo/skyobject"

// ThreadPage represents a ThreadPage as stored in cxo.
type ThreadPage struct {
	Thread skyobject.Reference  `skyobject:"schema=Thread"`
	Posts  skyobject.References `skyobject:"schema=Post"`
}
