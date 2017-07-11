package store

import (
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/src/store/typ"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"sync"
)

// StateSaver stores states for public keys.
type StateSaver struct {
	sync.Mutex
	store map[cipher.PubKey]*State
}

// NewStateSaver creates a new state saver.
func NewStateSaver() *StateSaver {
	return &StateSaver{
		store: make(map[cipher.PubKey]*State),
	}
}

// Init initiates StateSaver.
func (s *StateSaver) Init(c *CXO, pks ...cipher.PubKey) {
	s.Lock()
	defer s.Unlock()
	for _, pk := range pks {
		r := c.node.Container().LastFullRoot(pk)
		if r == nil {
			continue
		}
		s.store[pk] = NewState().Fill(r)
	}
}

func (s *StateSaver) getState(bpk cipher.PubKey) (*State, error) {
	s.Lock()
	defer s.Unlock()
	state, has := s.store[bpk]
	if !has {
		return nil, errors.Errorf(
			"board of public key '%s' not found in state saver", bpk.Hex())
	}
	return state, nil
}

// GetThreadVotes obtains thread votes.
func (s *StateSaver) GetThreadVotes(bpk cipher.PubKey, tRef skyobject.Reference) ([]typ.Vote, error) {
	state, e := s.getState(bpk)
	if e != nil {
		return nil, e
	}
	votes, has := state.GetThreadVotes(tRef)
	if !has {
		return nil, errors.Errorf(
			"thread of reference '%s:%s' not found in state saver", bpk.Hex(), tRef.String())
	}
	return votes, nil
}

// GetPostVotes obtains post votes.
func (s *StateSaver) GetPostVotes(bpk cipher.PubKey, pRef skyobject.Reference) ([]typ.Vote, error) {
	state, e := s.getState(bpk)
	if e != nil {
		return nil, e
	}
	votes, has := state.GetPostVotes(pRef)
	if !has {
		return nil, errors.Errorf(
			"post of reference '%s:%s' not found in state saver", bpk.Hex(), pRef.String())
	}
	return votes, nil
}

// ReplaceThreadVotes replaces thread votes.
func (s *StateSaver) ReplaceThreadVotes(bpk cipher.PubKey, tRef skyobject.Reference, votes []typ.Vote) {
	state, e := s.getState(bpk)
	if e != nil {
		state = NewState()
		s.Lock()
		s.store[bpk] = state
		s.Unlock()
	}
	state.tMux.Lock()
	state.tMap[tRef] = votes
	state.tMux.Unlock()
}

// ReplacePostVotes replaces post votes.
func (s *StateSaver) ReplacePostVotes(bpk cipher.PubKey, pRef skyobject.Reference, votes []typ.Vote) {
	state, e := s.getState(bpk)
	if e != nil {
		state = NewState()
		s.Lock()
		s.store[bpk] = state
		s.Unlock()
	}
	state.pMux.Lock()
	state.pMap[pRef] = votes
	state.pMux.Unlock()
}

// Fill fills a state.
func (s *StateSaver) Fill(r *node.Root) {
	s.Lock()
	state, has := s.store[r.Pub()]
	if !has {
		state = NewState()
		s.store[r.Pub()] = state
	}
	s.Unlock()

	state.Fill(r)
}

// State compiles maps for content and user votes.
type State struct {
	tMux sync.Mutex
	tMap map[skyobject.Reference][]typ.Vote
	pMux sync.Mutex
	pMap map[skyobject.Reference][]typ.Vote
}

// NewState creates a new internal state.
func NewState() *State {
	return &State{
		tMap: make(map[skyobject.Reference][]typ.Vote),
		pMap: make(map[skyobject.Reference][]typ.Vote),
	}
}

func (s *State) GetThreadVotes(ref skyobject.Reference) ([]typ.Vote, bool) {
	s.tMux.Lock()
	defer s.tMux.Unlock()
	votes, has := s.tMap[ref]
	return votes, has
}

func (s *State) GetPostVotes(ref skyobject.Reference) ([]typ.Vote, bool) {
	s.pMux.Lock()
	defer s.pMux.Unlock()
	votes, has := s.pMap[ref]
	return votes, has
}

// Fill fills the internal state.
func (s *State) Fill(r *node.Root) *State {
	tvsMap, e := s.fillThreadVotes(r.Walker())
	if e != nil {
		log.Printf(
			"[CONTAINER : INTERNAL STATE] Failed for board '%s'. Error: '%s'",
			r.Pub().Hex(), e.Error(),
		)
	} else {
		s.tMux.Lock()
		for k, v := range tvsMap {
			s.tMap[k] = v
		}
		s.tMux.Unlock()
	}
	pvsMap, e := s.fillPostVotes(r.Walker())
	if e != nil {
		log.Printf(
			"[CONTAINER : INTERNAL STATE] Failed for board '%s'. Error: '%s'",
			r.Pub().Hex(), e.Error(),
		)
	} else {
		s.pMux.Lock()
		for k, v := range pvsMap {
			s.pMap[k] = v
		}
		s.pMux.Unlock()
	}
	return s
}

func (s *State) fillThreadVotes(w *node.RootWalker) (map[skyobject.Reference][]typ.Vote, error) {
	vc := &typ.ThreadVotesContainer{}
	if e := w.AdvanceFromRoot(vc, makeThreadVotesContainerFinder(w.Root())); e != nil {
		return nil, e
	}
	threadVotesMap := make(map[skyobject.Reference][]typ.Vote)
	for _, threadVotes := range vc.Threads {
		votes := make([]typ.Vote, len(threadVotes.Votes))
		for i, ref := range threadVotes.Votes {
			if e := w.DeserializeFromRef(ref, &votes[i]); e != nil {
				return nil, errors.Errorf("failed to obtain vote '%s'", ref.String())
			}
		}
		threadVotesMap[threadVotes.Thread] = votes
	}
	return threadVotesMap, nil
}

func (s *State) fillPostVotes(w *node.RootWalker) (map[skyobject.Reference][]typ.Vote, error) {
	vc := &typ.PostVotesContainer{}
	if e := w.AdvanceFromRoot(vc, makePostVotesContainerFinder(w.Root())); e != nil {
		return nil, e
	}
	postVotesMap := make(map[skyobject.Reference][]typ.Vote)
	for _, postVotes := range vc.Posts {
		votes := make([]typ.Vote, len(postVotes.Votes))
		for i, ref := range postVotes.Votes {
			if e := w.DeserializeFromRef(ref, &votes[i]); e != nil {
				return nil, errors.Errorf("failed to obtain vote '%s'", ref.String())
			}
		}
		postVotesMap[postVotes.Post] = votes
	}
	return postVotesMap, nil
}
