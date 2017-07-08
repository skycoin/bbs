package obj

import "github.com/skycoin/cxo/skyobject"

type BoardPage struct {
	Board       skyobject.Reference `skyobject:"schema=Board"`
	ThreadPages skyobject.Reference `skyobject:"schema=ThreadPage"`
}
