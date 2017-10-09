package state

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/keys"
	"github.com/skycoin/bbs/src/store/object/revisions/r0"
	"github.com/skycoin/bbs/src/store/state/pack"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

func (bi *BoardInstance) Submit(content *r0.Content) (uint64, error) {

	var (
		goal uint64
		contentType = content.GetHeader().Type
	)

	switch contentType {
	case r0.V5ThreadType:
		if e := submitThread(bi, &goal, content.ToThread()); e != nil {
			return 0, e
		}
	case r0.V5PostType:
		if e := submitPost(bi, &goal, content.ToPost()); e != nil {
			return 0, e
		}
	case r0.V5ThreadVoteType:
		if e := submitThreadVote(bi, &goal, content.ToThreadVote()); e != nil {
			return 0, e
		}
	case r0.V5PostVoteType:
		if e := submitPostVote(bi, &goal, content.ToPostVote()); e != nil {
			return 0, e
		}
	case r0.V5UserVoteType:
		if e := submitUserVote(bi, &goal, content.ToUserVote()); e != nil {
			return 0, e
		}
	default:
		return 0, boo.Newf(boo.InvalidInput,
			"content has invalid type '%s'", contentType)
	}

	return goal, nil
}

func submitThread(bi *BoardInstance, goal *uint64, thread *r0.Thread) error {
	return bi.EditPack(func(p *skyobject.Pack, h *pack.Headers) error {

		// Set goal sequence.
		*goal = p.Root().Seq + 1

		// Get thread ref.
		tRef := p.Ref(thread.Content)

		// Ensure thread does not exist.
		if _, has := h.GetThreadPageHash(thread.GetHeader().Hash); has {
			return boo.Newf(boo.AlreadyExists,
				"Thread of ref %s already exists", tRef.String())
		}

		// Get root children pages.
		pages, e := r0.GetPages(p, false, true, true, true)
		if e != nil {
			return e
		}

		// Add thread to board page.
		if e := pages.BoardPage.AddThread(tRef); e != nil {
			return e
		}

		// Add thread to diff page.
		if e := pages.DiffPage.Add(thread.Content); e != nil {
			return e
		}

		// Save changes.
		return pages.Save(p)
	})
}

func submitPost(bi *BoardInstance, goal *uint64, post *r0.Post) error {
	body := post.GetBody()

	return bi.EditPack(func(p *skyobject.Pack, h *pack.Headers) error {

		// Set goal sequence.
		*goal = p.Root().Seq + 1

		// Get post ref.
		pRef := p.Ref(post.Content)

		// Ensure thread exists.
		tpHash, has := h.GetThreadPageHash(body.OfThread)
		if !has {
			return boo.Newf(boo.NotFound,
				"thread of hash '%s' not found", body.OfThread)
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
		h.SetThread(body.OfThread, tpRef.Hash)

		// Add post to diff.
		if e := pages.DiffPage.Add(post.Content); e != nil {
			return e
		}

		// Save changes.
		return pages.Save(p)
	})
}

func submitThreadVote(bi *BoardInstance, goal *uint64, tVote *r0.ThreadVote) error {
	body := tVote.GetBody()

	return bi.EditPack(func(p *skyobject.Pack, h *pack.Headers) error {
		*goal = p.Root().Seq + 1

		// Check vote.
		if _, ok := h.GetThreadPageHash(body.OfThread); !ok {
			return boo.Newf(boo.NotFound,
				"thread of hash '%s' is not found", body.OfThread)
		}

		return addVoteToProfile(p, h, tVote.Content, body.Creator)
	})
}

func submitPostVote(bi *BoardInstance, goal *uint64, pVote *r0.PostVote) error {
	body := pVote.GetBody()

	return bi.EditPack(func(p *skyobject.Pack, h *pack.Headers) error {
		*goal = p.Root().Seq + 1
		return addVoteToProfile(p, h, pVote.Content, body.Creator)
	})
}

func submitUserVote(bi *BoardInstance, goal *uint64, uVote *r0.UserVote) error {
	body := uVote.GetBody()
	return bi.EditPack(func(p *skyobject.Pack, h *pack.Headers) error {
		*goal = p.Root().Seq + 1
		return addVoteToProfile(p, h, uVote.Content, body.Creator)
	})
}

func addVoteToProfile(p *skyobject.Pack, h *pack.Headers, content *r0.Content, creator string) error {

	// Get root children pages.
	pages, e := r0.GetPages(p, false, false, true, true)
	if e != nil {
		return e
	}

	// Get profile page (create if not exist).
	profileHash, ok := h.GetUserProfileHash(creator)
	if !ok {
		if profileHash, e = pages.UsersPage.NewUserProfile(creator); e != nil {
			return e
		}
		h.SetUser(creator, profileHash)
	}

	// Add vote to appropriate user profile.
	newProfileHash, e := pages.UsersPage.AddUserSubmission(profileHash, content)
	if e != nil {
		return e
	}
	h.SetUser(creator, newProfileHash)

	// Add to diff.
	if e := pages.DiffPage.Add(content); e != nil {
		return e
	}

	return pages.Save(p)
}

func (bi *BoardInstance) EnsureSubmissionKeys(pks []cipher.PubKey) (uint64, error) {
	bi.l.Println("ensuring submission keys as:", keys.PubKeyArrayToString(pks))
	return bi.EditBoard(func(board *r0.Board) (bool, error) {
		body := board.GetBody()
		if keys.ComparePubKeyArrays(body.GetSubKeys(), pks) {
			return false, nil
		}
		body.SetSubKeys(pks)
		board.SetBody(body)
		return true, nil
	})
}

func (bi *BoardInstance) GetSubmissionKeys() []cipher.PubKey {
	var pks []cipher.PubKey
	if e := bi.ViewBoard(func(board *r0.Board) (bool, error) {
		pks = board.GetBody().GetSubKeys()
		return false, nil
	}); e != nil {
		bi.l.Println("error obtaining submission keys:", e)
		return nil
	}
	return pks
}

// BoardAction is a function in which board modification/viewing takes place.
// Returns a boolean that represents whether changes have been made and
// an error on failure.
type BoardAction func(board *r0.Board) (bool, error)

// EditBoard triggers a board action.
func (bi *BoardInstance) EditBoard(action BoardAction) (uint64, error) {
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
		if e := pages.BoardPage.Board.SetValue(board.Content); e != nil {
			return e
		}

		return pages.Save(p)
	})
	return goalSeq, e
}

func (bi *BoardInstance) ViewBoard(action BoardAction) error {
	return bi.ViewPack(func(p *skyobject.Pack, h *pack.Headers) error {
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
		if _, e := action(board); e != nil {
			return e
		}
		return nil
	})
}
