package store

import (
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/src/store/typ"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

// GetVotesForThread obtains the votes for specified thread present in specified board.
func (c *CXO) GetVotesForThread(bpk cipher.PubKey, tRef skyobject.Reference) []typ.Vote {
	c.Lock(c.GetVotesForThread)
	defer c.Unlock()
	return c.ss.GetThreadVotes(c, bpk, tRef)
}

// GetVotesForPost obtains the votes for specified post present in specified board.
func (c *CXO) GetVotesForPost(bpk cipher.PubKey, pRef skyobject.Reference) []typ.Vote {
	c.Lock(c.GetVotesForPost)
	defer c.Unlock()
	return c.ss.GetPostVotes(c, bpk, pRef)
}

func (c *CXO) GetVotesForUser(bpk, upk cipher.PubKey) []typ.Vote {
	c.Lock(c.GetVotesForUser)
	defer c.Unlock()
	return c.ss.GetUserVotes(c, bpk, upk)
}

// VoteForThread adds a vote for a thread on a specified board.
func (c *CXO) AddVoteForThread(bpk cipher.PubKey, bsk cipher.SecKey, tRef skyobject.Reference, newVote *typ.Vote) error {
	c.Lock(c.AddVoteForThread)
	defer c.Unlock()

	w := c.c.LastRootSk(bpk, bsk).Walker()
	vc := &typ.ThreadVotesContainer{}
	if e := w.AdvanceFromRoot(vc, makeThreadVotesContainerFinder(w.Root())); e != nil {
		return e
	}
	// Obtain vote references.
	voteRefs, e := vc.GetThreadVotes(tRef)
	if e != nil {
		return e
	}
	// Loop through votes to see if user has already voted.
	for i, vRef := range voteRefs.Votes {
		tempVote := &typ.Vote{}
		if e := w.DeserializeFromRef(vRef, tempVote); e != nil {
			return errors.Wrap(e, "failed to obtain vote")
		}
		// Replace vote if already voted.
		if tempVote.User == newVote.User {
			voteRefs.Votes[i] = w.Root().Save(*newVote)
			goto SaveThreadVotesContainer
		}
	}
	// If user has not already voted - add a new vote.
	voteRefs.Votes = append(voteRefs.Votes, w.Root().Save(*newVote))

SaveThreadVotesContainer:
	defer c.ss.Fill(w.Root()) // TODO: Optimise.
	return w.ReplaceCurrent(*vc)
}

// VoteForPost adds a vote for a post on a specified board.
func (c *CXO) AddVoteForPost(bpk cipher.PubKey, bsk cipher.SecKey, pRef skyobject.Reference, newVote *typ.Vote) error {
	c.Lock(c.AddVoteForPost)
	defer c.Unlock()

	w := c.c.LastRootSk(bpk, bsk).Walker()
	vc := &typ.PostVotesContainer{}
	if e := w.AdvanceFromRoot(vc, makePostVotesContainerFinder(w.Root())); e != nil {
		return e
	}
	// Obtain vote references.
	voteRefs, e := vc.GetPostVotes(pRef)
	if e != nil {
		return e
	}
	// Loop through votes to see if user has voted already.
	for i, vRef := range voteRefs.Votes {
		tempVote := &typ.Vote{}
		if e := w.DeserializeFromRef(vRef, tempVote); e != nil {
			return errors.Wrap(e, "failed to obtain vote")
		}
		// Replace vote if already voted.
		if tempVote.User == newVote.User {
			voteRefs.Votes[i] = w.Root().Save(*newVote)
			goto SavePostVotesContainer
		}
	}
	// If user has not already voted - add a new vote.
	voteRefs.Votes = append(voteRefs.Votes, w.Root().Save(*newVote))

SavePostVotesContainer:
	defer c.ss.Fill(w.Root()) // TODO: Optimise.
	return w.ReplaceCurrent(*vc)
}

func (c *CXO) AddVoteForUser(bpk, upk cipher.PubKey, bsk cipher.SecKey, newVote *typ.Vote) error {
	c.Lock(c.AddVoteForUser)
	defer c.Unlock()

	w := c.c.LastRootSk(bpk, bsk).Walker()
	vc := &typ.UserVotesContainer{}
	if e := w.AdvanceFromRoot(vc, makeUserVotesContainerFinder(w.Root())); e != nil {
		return e
	}
	// Obtain vote references.
	voteRefs := vc.GetUser(upk)
	// Loop through votes to see if user has already voted.
	for i, vRef := range voteRefs.Votes {
		tempVote := &typ.Vote{}
		if e := w.DeserializeFromRef(vRef, tempVote); e != nil {
			return errors.Wrap(e, "failed to obtain vote")
		}
		// Replace vote if already voted.
		if tempVote.User == newVote.User {
			voteRefs.Votes[i] = w.Root().Save(*newVote)
			goto SaveUserVotesContainer
		}
	}
	// If user has not already voted - add a new vote.
	voteRefs.Votes = append(voteRefs.Votes, w.Root().Save(*newVote))

SaveUserVotesContainer:
	defer c.ss.Fill(w.Root()) // TODO: Optimise.
	return w.ReplaceCurrent(*vc)
}

// RemoveVoteForThread removes a vote completely for a thread and specified user.
func (c *CXO) RemoveVoteForThread(upk, bpk cipher.PubKey, bsk cipher.SecKey, tRef skyobject.Reference) error {
	c.Lock(c.RemoveVoteForThread)
	defer c.Unlock()

	w := c.c.LastRootSk(bpk, bsk).Walker()
	vc := &typ.ThreadVotesContainer{}
	if e := w.AdvanceFromRoot(vc, makeThreadVotesContainerFinder(w.Root())); e != nil {
		return e
	}
	// Obtain vote references.
	voteRefs, e := vc.GetThreadVotes(tRef)
	if e != nil {
		return errors.Wrap(e, "failed to remove vote for thread")
	}
	// loop through votes to see what to remove.
	for i, vRef := range voteRefs.Votes {
		tempVote := &typ.Vote{}
		if e := w.DeserializeFromRef(vRef, tempVote); e != nil {
			return errors.Wrap(e, "failed to obtain vote")
		}
		if tempVote.User == upk {
			// Delete of index i.
			voteRefs.Votes[i], voteRefs.Votes[0] = voteRefs.Votes[0], voteRefs.Votes[i]
			voteRefs.Votes = voteRefs.Votes[1:]
			defer c.ss.Fill(w.Root()) // TODO: Optimise.
			return w.ReplaceCurrent(*vc)
		}
	}
	return nil
}

// RemoveVoteForPost removes a vote completely for a post and specified user.
func (c *CXO) RemoveVoteForPost(upk, bpk cipher.PubKey, bsk cipher.SecKey, pRef skyobject.Reference) error {
	c.Lock(c.RemoveVoteForPost)
	defer c.Unlock()

	w := c.c.LastRootSk(bpk, bsk).Walker()
	vc := &typ.PostVotesContainer{}
	if e := w.AdvanceFromRoot(vc, makePostVotesContainerFinder(w.Root())); e != nil {
		return e
	}
	// Obtain vote references.
	voteRefs, e := vc.GetPostVotes(pRef)
	if e != nil {
		return e
	}
	// loop through votes to see what to remove.
	for i, vRef := range voteRefs.Votes {
		tempVote := &typ.Vote{}
		if e := w.DeserializeFromRef(vRef, tempVote); e != nil {
			return errors.Wrap(e, "failed to obtain vote")
		}
		if tempVote.User == upk {
			// Delete of index i.
			voteRefs.Votes[i], voteRefs.Votes[0] = voteRefs.Votes[0], voteRefs.Votes[i]
			voteRefs.Votes = voteRefs.Votes[1:]
			defer c.ss.Fill(w.Root()) // TODO: Optimise.
			return w.ReplaceCurrent(*vc)
		}
	}
	return nil
}

func (c *CXO) RemoveVoteForUser(currentUPK, bpk, upk cipher.PubKey, bsk cipher.SecKey) error {
	c.Lock(c.RemoveVoteForUser)
	defer c.Unlock()

	w := c.c.LastRootSk(bpk, bsk).Walker()
	vc := &typ.UserVotesContainer{}
	if e := w.AdvanceFromRoot(vc, makeUserVotesContainerFinder(w.Root())); e != nil {
		return e
	}
	// Obtain vote references.
	voteRefs := vc.GetUser(upk)
	// Loop through votes to see what to remove.
	for i, vRef := range voteRefs.Votes {
		tempVote := &typ.Vote{}
		if e := w.DeserializeFromRef(vRef, tempVote); e != nil {
			return errors.Wrap(e, "failed to obtain vote")
		}
		if tempVote.User == currentUPK {
			// Delete of index i.
			voteRefs.Votes[i], voteRefs.Votes[0] = voteRefs.Votes[0], voteRefs.Votes[i]
			voteRefs.Votes = voteRefs.Votes[1:]
			defer c.ss.Fill(w.Root()) // TODO: Optimise.
			return w.ReplaceCurrent(*vc)
		}
	}
	return nil
}
