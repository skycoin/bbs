package obj

import "github.com/skycoin/cxo/skyobject"

type ThreadPage struct {
	Thread skyobject.Reference  `skyobject:"schema=Thread"`
	Posts  skyobject.References `skyobject:"schema=Post"`
}
