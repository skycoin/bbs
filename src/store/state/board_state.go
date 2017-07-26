package state

import (
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/skyobject"
)

// BoardState represents an internal state of a board.
type BoardState struct {
	newRoots chan *skyobject.Root
	t        map[skyobject.Reference]*object.VoteSummary
	b        map[skyobject.Reference]*object.VoteSummary
}

func NewBoardState() BoardState {
	bs := BoardState{
		newRoots: make(chan *skyobject.Root),
		t:        make(map[skyobject.Reference]*object.VoteSummary),
		b:        make(map[skyobject.Reference]*object.VoteSummary),
	}
	return bs
}
