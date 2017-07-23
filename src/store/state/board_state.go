package state

import (
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/skyobject"
)

// BoardState represents an internal state of a board.
type BoardState struct {
	t map[skyobject.Reference]*object.VoteSummary
	b map[skyobject.Reference]*object.VoteSummary
}
