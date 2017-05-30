package cxo

import (
	"github.com/skycoin/cxo/skyobject"
	"github.com/evanlinjin/bbs/intern/typ"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/pkg/errors"
)

// GetVotesForThread obtains the votes for specified thread present in specified board.
func (c *Container) GetVotesForThread(bpk cipher.PubKey, tRef skyobject.Reference) ([]typ.Vote, error) {
	w := c.c.LastFullRoot(bpk).Walker()
	tvc := &typ.ThreadVoteContainer{}
	if e := w.AdvanceFromRoot(tvc, makeThreadVoteContainerFinder(w.Root())); e != nil {
		return nil, e
	}
	voteRefs, e := tvc.GetThreadVoteRefs(tRef)
	if e != nil {
		return nil, e
	}
	votes := make([]typ.Vote, len(voteRefs))
	for i, ref := range voteRefs {
		if e := w.DeserializeFromRef(ref, &votes[i]); e != nil {
			return nil, errors.Errorf("failed to obtain vote '%s'", ref.String())
		}
	}
	return votes, nil
}

// GetVotesForPost obtains the votes for specified post present in specified board.
func (c *Container) GetVotesForPost(bpk cipher.PubKey, pRef skyobject.Reference) ([]typ.Vote, error) {
	//w := c.c.LastFullRoot(bpk).Walker()
	// TODO: Complete!!!!
	return nil, nil
}