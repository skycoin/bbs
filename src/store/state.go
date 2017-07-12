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
	enable bool
	store  map[cipher.PubKey]*State
}

// NewStateSaver creates a new state saver.
func NewStateSaver(enable bool) *StateSaver {
	return &StateSaver{
		enable: enable,
		store:  make(map[cipher.PubKey]*State),
	}
}

// Init initiates StateSaver.
func (s *StateSaver) Init(c *CXO, pks ...cipher.PubKey) {
	if !s.enable {
		return
	}
	s.Lock()
	defer s.Unlock()
	s.store = nil
	s.store = map[cipher.PubKey]*State{}
	for _, pk := range pks {
		r := c.node.Container().LastFullRoot(pk)
		if r == nil {
			continue
		}
		s.store[pk] = NewState().Fill(r)
	}
}

func (s *StateSaver) getState(bpk cipher.PubKey) *State {
	s.Lock()
	defer s.Unlock()
	state, has := s.store[bpk]
	if !has {
		state = NewState()
		s.store[bpk] = state
	}
	return state
}

// GetThreadVotes obtains thread votes.
func (s *StateSaver) GetThreadVotes(cxo *CXO, bpk cipher.PubKey, tRef skyobject.Reference) []typ.Vote {
	if !s.enable {
		return nil
	}
	votes, has := s.getState(bpk).GetThreadVotes(tRef)
	if !has {
		s.Fill(cxo.c.LastFullRoot(bpk))
		votes, has := s.getState(bpk).GetThreadVotes(tRef)
		if !has {
			log.Printf("thread of reference '%s:%s' not found in state saver", bpk.Hex(), tRef.String())
			return nil
		}
		return votes
	}
	return votes
}

// GetPostVotes obtains post votes.
func (s *StateSaver) GetPostVotes(cxo *CXO, bpk cipher.PubKey, pRef skyobject.Reference) []typ.Vote {
	if !s.enable {
		return nil
	}
	votes, has := s.getState(bpk).GetPostVotes(pRef)
	if !has {
		s.Fill(cxo.c.LastFullRoot(bpk))
		votes, has := s.getState(bpk).GetPostVotes(pRef)
		if !has {
			log.Printf("post of reference '%s:%s' not found in state saver", bpk.Hex(), pRef.String())
			return nil
		}
		return votes
	}
	return votes
}

// GetUserVotes obtains user votes.
func (s *StateSaver) GetUserVotes(cxo *CXO, bpk, upk cipher.PubKey) []typ.Vote {
	if !s.enable {
		return nil
	}
	votes, has := s.getState(bpk).GetUserVotes(upk)
	if !has {
		s.Fill(cxo.c.LastFullRoot(bpk))
		votes, has := s.getState(bpk).GetUserVotes(upk)
		if !has {
			log.Printf("user of reference '%s:%s' not found in state saver", bpk.Hex(), upk.Hex())
			return nil
		}
		return votes
	}
	return votes
}

// ReplaceThreadVotes replaces thread votes.
func (s *StateSaver) ReplaceThreadVotes(bpk cipher.PubKey, tRef skyobject.Reference, votes []typ.Vote) {
	if !s.enable {
		return
	}
	state := s.getState(bpk)
	state.tMux.Lock()
	state.tMap[tRef] = votes
	state.tMux.Unlock()
}

// ReplacePostVotes replaces post votes.
func (s *StateSaver) ReplacePostVotes(bpk cipher.PubKey, pRef skyobject.Reference, votes []typ.Vote) {
	if !s.enable {
		return
	}
	state := s.getState(bpk)
	state.pMux.Lock()
	state.pMap[pRef] = votes
	state.pMux.Unlock()
}

// ReplaceUserVotes replaces user votes.
func (s *StateSaver) ReplaceUserVotes(bpk, upk cipher.PubKey, votes []typ.Vote) {
	if !s.enable {
		return
	}
	state := s.getState(bpk)
	state.uMux.Lock()
	state.uMap[upk] = votes
	state.uMux.Unlock()
}

// Fill fills a state.
func (s *StateSaver) Fill(r *node.Root) {
	if !s.enable {
		return
	}
	s.Lock()
	state, has := s.store[r.Pub()]
	if !has {
		state = NewState()
		s.store[r.Pub()] = state
	}
	s.Unlock()

	state.Fill(r)
}

// Remove removes a board state.
func (s *StateSaver) Remove(bpk cipher.PubKey) {
	if !s.enable {
		return
	}
	s.Lock()
	defer s.Unlock()
	delete(s.store, bpk)
}

func (s *StateSaver) RemoveThreadVotes(bpk cipher.PubKey, tRef skyobject.Reference) {
	if !s.enable {
		return
	}
	state := s.getState(bpk)
	state.tMux.Lock()
	delete(state.tMap, tRef)
	state.tMux.Unlock()
}

func (s *StateSaver) RemovePostVotes(bpk cipher.PubKey, pRef skyobject.Reference) {
	if !s.enable {
		return
	}
	state := s.getState(bpk)
	state.pMux.Lock()
	delete(state.pMap, pRef)
	state.pMux.Unlock()
}

// State compiles maps for content and user votes.
type State struct {
	tMux sync.Mutex
	tMap map[skyobject.Reference][]typ.Vote
	pMux sync.Mutex
	pMap map[skyobject.Reference][]typ.Vote
	uMux sync.Mutex
	uMap map[cipher.PubKey][]typ.Vote
}

// NewState creates a new internal state.
func NewState() *State {
	return &State{
		tMap: make(map[skyobject.Reference][]typ.Vote),
		pMap: make(map[skyobject.Reference][]typ.Vote),
		uMap: make(map[cipher.PubKey][]typ.Vote),
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

func (s *State) GetUserVotes(pk cipher.PubKey) ([]typ.Vote, bool) {
	s.uMux.Lock()
	defer s.uMux.Unlock()
	votes, has := s.uMap[pk]
	return votes, has
}

// Fill fills the internal state.
func (s *State) Fill(r *node.Root) *State {
	tvsMap, e := s.fillThreadVotes(r.Walker())
	if e != nil {
		log.Printf(
			"[CONTAINER : INTERNAL STATE] fillThreadVotes failed for board '%s'. Error: '%s'",
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
			"[CONTAINER : INTERNAL STATE] fillPostVotes failed for board '%s'. Error: '%s'",
			r.Pub().Hex(), e.Error(),
		)
	} else {
		s.pMux.Lock()
		for k, v := range pvsMap {
			s.pMap[k] = v
		}
		s.pMux.Unlock()
	}
	uvsMap, e := s.fillUserVotes(r.Walker())
	if e != nil {
		log.Printf(
			"[CONTAINER : INTERNAL STATE] fillUserVotes failed for board '%s'. Error: '%s'",
			r.Pub().Hex(), e.Error(),
		)
	} else {
		s.uMux.Lock()
		for k, v := range uvsMap {
			s.uMap[k] = v
		}
		s.uMux.Unlock()
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

func (s *State) fillUserVotes(w *node.RootWalker) (map[cipher.PubKey][]typ.Vote, error) {
	vc := &typ.UserVotesContainer{}
	if e := w.AdvanceFromRoot(vc, makeUserVotesContainerFinder(w.Root())); e != nil {
		return nil, e
	}
	userVotesMap := make(map[cipher.PubKey][]typ.Vote)
	for _, userVotes := range vc.Users {
		votes := make([]typ.Vote, len(userVotes.Votes))
		for i, ref := range userVotes.Votes {
			if e := w.DeserializeFromRef(ref, &votes[i]); e != nil {
				return nil, errors.Errorf("failed to obtain vote '%s'", ref.String())
			}
		}
		userVotesMap[userVotes.User] = votes
	}
	return userVotesMap, nil
}
