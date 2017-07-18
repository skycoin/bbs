package state

import (
	"github.com/skycoin/bbs/src/store/obj"
	"github.com/skycoin/cxo/skyobject"
)

// BoardState represents an internal state of a board.
type BoardState struct {
	t map[skyobject.Reference]*obj.VoteSummary
	b map[skyobject.Reference]*obj.VoteSummary
}
