package obj

import (
	"github.com/skycoin/cxo/skyobject"
)

type Thread struct {
	Post
	MasterBoardRef skyobject.Reference `json:"-" skyobject:"schema=Board"`
}
