package typ

import "github.com/skycoin/cxo/skyobject"

// BoardContainer represents the first branch object from root, in cxo.
type BoardContainer struct {
	Board       skyobject.Reference  `skyobject:"schema=Board"`
	Threads     skyobject.References `skyobject:"schema=Thread"`
	ThreadPages skyobject.References `skyobject:"schema=ThreadPage"`
}
