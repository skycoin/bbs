package state

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/keys"
	"github.com/skycoin/bbs/src/store/object/revisions/r0"
	"github.com/skycoin/bbs/src/store/state/pack"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

// NewThread triggers a request to add a new thread to the board.
// Returns sequence of root in which thread is to be available, or error on failure.
func (bi *BoardInstance) NewThread(thread *r0.Thread) (uint64, error) {

	var goalSeq uint64
	e := bi.EditPack(func(p *skyobject.Pack, h *pack.Headers) error {

		// Set goal seq.
		goalSeq = p.Root().Seq + 1

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

	return goalSeq, e
}

// NewPost triggers a request to add a new post to the board.
// Returns sequence of root in which post is to be available, or error on failure.
func (bi *BoardInstance) NewPost(post *r0.Post) (uint64, error) {
	pBody := post.GetBody()

	var goalSeq uint64
	e := bi.EditPack(func(p *skyobject.Pack, h *pack.Headers) error {

		// Set goal seq.
		goalSeq = p.Root().Seq + 1

		// Get post ref.
		pRef := p.Ref(post.Content)

		// Ensure thread exists.
		tpHash, has := h.GetThreadPageHash(pBody.OfThread)
		if !has {
			return boo.Newf(boo.NotFound,
				"thread of hash '%s' not found", pBody.OfThread)
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
		h.SetThread(pBody.OfThread, tpRef.Hash)

		// Add post to diff.
		if e := pages.DiffPage.Add(post.Content); e != nil {
			return e
		}

		// Save changes.
		return pages.Save(p)
	})

	return goalSeq, e
}

func (bi *BoardInstance) NewThreadVote(threadVote *r0.ThreadVote) (uint64, error) {
	tvBody := threadVote.GetBody()

	var goalSeq uint64
	e := bi.EditPack(func(p *skyobject.Pack, h *pack.Headers) error {

		// Check vote.
		if _, ok := h.GetThreadPageHash(tvBody.OfThread); !ok {
			return boo.Newf(boo.NotFound,
				"thread of hash '%s' is not found", tvBody.OfThread)
		}

		return addVoteToProfile(p, h, threadVote.Content, tvBody.Creator)
	})
	return goalSeq, e
}

func (bi *BoardInstance) NewPostVote(postVote *r0.PostVote) (uint64, error) {
	pvBody := postVote.GetBody()

	var goalSeq uint64
	e := bi.EditPack(func(p *skyobject.Pack, h *pack.Headers) error {
		return addVoteToProfile(p, h, postVote.Content, pvBody.Creator)
	})
	return goalSeq, e
}

func (bi *BoardInstance) NewUserVote(userVote *r0.UserVote) (uint64, error) {
	uvBody := userVote.GetBody()

	var goalSeq uint64
	e := bi.EditPack(func(p *skyobject.Pack, h *pack.Headers) error {
		return addVoteToProfile(p, h, userVote.Content, uvBody.Creator)
	})
	return goalSeq, e
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
