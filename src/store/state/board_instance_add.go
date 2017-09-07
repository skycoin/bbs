package state

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object/revisions/r0"
	"github.com/skycoin/bbs/src/store/state/pack"
	"github.com/skycoin/cxo/skyobject"
)

func (bi *BoardInstance) NewThread(thread *r0.Thread) (uint64, error) {
	if e := thread.Verify(); e != nil {
		return 0, e
	}

	var goalSeq uint64
	e := bi.EditPack(func(p *skyobject.Pack, h *pack.Headers) error {

		// Set goal seq.
		goalSeq = p.Root().Seq + 1

		// Get thread ref.
		tRef := p.Ref(thread)

		// Ensure thread does not exist.
		if _, has := h.GetThreadPageHash(tRef.Hash); has {
			return boo.Newf(boo.AlreadyExists,
				"Thread of ref %s already exists", tRef.String())
		}

		// Get root children pages.
		pages, e := r0.GetPages(p, false, true, true, false)
		if e != nil {
			return e
		}

		// Add thread to board page.
		if e := pages.BoardPage.AddThread(tRef); e != nil {
			return e
		}

		// Add thread to diff page.
		if e := pages.DiffPage.Add(thread); e != nil {
			return e
		}

		// Save changes.
		return pages.Save(p)
	})

	return goalSeq, e
}

func (bi *BoardInstance) NewPost(post *r0.Post) (uint64, error) {
	if e := post.Verify(); e != nil {
		return 0, e
	}

	var goalSeq uint64
	e := bi.EditPack(func(p *skyobject.Pack, h *pack.Headers) error {

		// Set goal seq.
		goalSeq = p.Root().Seq + 1

		// Get post ref.
		pRef := p.Ref(post)

		// Ensure thread exists.
		tpHash, has := h.GetThreadPageHash(post.OfThread)
		if !has {
			return boo.Newf(boo.NotFound,
				"thread of hash '%s' not found", post.OfThread)
		}

		// Get root pages.
		pages, e := r0.GetPages(p, false, true, true, false)
		if e != nil {
			return e
		}

		// Add post to board page.
		tpRef, tPage, e := pages.BoardPage.GetThreadPage(tpHash)
		if e != nil {
			return e
		}
		if e := tPage.AddPost(pRef.Hash, post); e != nil {
			return e
		}
		if e := tPage.Save(tpRef); e != nil {
			return e
		}

		// Modify headers.
		h.SetThread(post.OfThread, tpRef.Hash)

		// Add post to diff.
		if e := pages.DiffPage.Add(post); e != nil {
			return e
		}

		// Save changes.
		return pages.Save(p)
	})

	return goalSeq, e
}

func (bi *BoardInstance) NewVote(vote *r0.Vote) (uint64, error) {
	if e := vote.Verify(); e != nil {
		return 0, e
	}

	var goalSeq uint64
	e := bi.EditPack(func(p *skyobject.Pack, h *pack.Headers) error {

		// Check vote.
		if e := checkVote(vote, h); e != nil {
			return e
		}

		// Set goal seq.
		goalSeq = p.Root().Seq + 1

		// Get root children pages.
		pages, e := r0.GetPages(p, false, false, true, true)
		if e != nil {
			return e
		}

		// Get users page. Create user activity page if not exist.
		uapHash, has := h.GetUserActivityPageHash(vote.Creator)
		if !has {
			uapHash, e = pages.UsersPage.NewUserActivityPage(vote.Creator)
			if e != nil {
				return e
			}
			h.SetUser(vote.Creator, uapHash)
		}

		// Add vote to appropriate user activity page.
		if uapHashNew, e := pages.UsersPage.AddUserActivity(uapHash, vote); e != nil {
			return e
		} else {
			h.SetUser(vote.Creator, uapHashNew)
		}

		// Add vote to diff page.
		if e := pages.DiffPage.Add(vote); e != nil {
			return e
		}

		// Save changes.
		return pages.Save(p)
	})

	return goalSeq, e
}

func checkVote(vote *r0.Vote, h *pack.Headers) error {
	switch vote.GetType() {
	case r0.UserVote:
		// TODO.

	case r0.ThreadVote:
		_, ok := h.GetThreadPageHash(vote.OfThread)
		if !ok {
			return boo.Newf(boo.NotFound,
				"thread of hash '%s' is not found",
				vote.OfThread.Hex())
		}

	case r0.PostVote:
		// TODO.

	default:
		return boo.Newf(boo.NotAllowed,
			"invalid vote type of '%s'",
			r0.VoteString[r0.UnknownVoteType])
	}

	return nil
}

type BoardAction func(board *r0.Board) (bool, error)

func (bi *BoardInstance) BoardAction(action BoardAction) (uint64, error) {
	var goalSeq uint64
	e := bi.EditPack(func(p *skyobject.Pack, h *pack.Headers) error {

		// Set goal seq.
		goalSeq = p.Root().Seq + 1

		// Get root children.
		pages, e := r0.GetPages(p, false, true, false, false)
		if e != nil {
			return e
		}

		// Get board.
		board, e := pages.BoardPage.GetBoard()
		if e != nil {
			return e
		}

		// Do action to board.
		if save, e := action(board); e != nil {
			return e
		} else if !save {
			return nil
		}

		// Save changes.
		if e := pages.BoardPage.Board.SetValue(board); e != nil {
			return e
		}

		return pages.Save(p)
	})
	return goalSeq, e
}
