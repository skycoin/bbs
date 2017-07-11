package access

import (
	"github.com/skycoin/bbs/src/boo"
	"github.com/skycoin/bbs/src/store/view"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
)

// StateSaver saves the internal exposed data structure.
type StateSaver struct {
	mux    sync.Mutex
	boards map[cipher.PubKey]*State
}

// NewStateSaver creates a new state saver.
func NewStateSaver() *StateSaver {
	return &StateSaver{
		boards: make(map[cipher.PubKey]*State),
	}
}

func (s *StateSaver) lock() func() {
	s.mux.Lock()
	return s.mux.Unlock
}

// Close closes the StateSaver.
func (s *StateSaver) Close() {
	for _, state := range s.boards {
		state.Close()
	}
}

func (s *StateSaver) Get(pk cipher.PubKey) (view.Board, error) {
	defer s.lock()()
	state, has := s.boards[pk]
	if !has {
		return view.Board{}, boo.Newf(boo.ObjectNotFound,
			"board of public key '%s' not found in internal state", pk.Hex())
	}
	return state.GetView(), nil
}

// Update updates a board state of StateSaver.
func (s *StateSaver) Update(r *node.Root) {
	defer s.lock()()
	s.update(r.Root)
}

// UpdateLocal updates a board state of StateSaver as local.
func (s *StateSaver) UpdateLocal(r *skyobject.Root) {
	defer s.lock()()
	s.update(r)
}

// Remove removes a board state of public key.
func (s *StateSaver) Remove(pk cipher.PubKey) {
	defer s.lock()()
	delete(s.boards, pk)
}

func (s *StateSaver) update(r *skyobject.Root) {
	bpk := r.Pub()
	state, has := s.boards[bpk]
	if !has {
		state = NewState()
		s.boards[bpk] = state
	}
	state.PushNewBoardRoot(r)
}
