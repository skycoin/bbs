package store

import (
	"github.com/skycoin/bbs/src/store/view"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
)

// State represents a state of a board.
type State struct {
	board *view.Board
	has   map[cipher.PubKey]struct{}
}

// StateSaver saves the internal exposed data structure.
type StateSaver struct {
	sync.Mutex
	boards map[cipher.PubKey]*State
}

func NewStateSaver() *StateSaver {
	return &StateSaver{
		boards: make(map[cipher.PubKey]*State),
	}
}

// Update updates a board state of StateSaver.
func (s *StateSaver) Update(root *node.Root) {
	s.Lock()
	defer s.Unlock()
}
