package state

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/skyobject"
)

func (bi *BoardInstance) Submit(transport *object.Transport) (uint64, error) {

	var goal uint64

	switch transport.Body.Type {
	case object.V5ThreadType:
		if e := submitThread(bi, &goal, transport.Content); e != nil {
			return 0, e
		}
	case object.V5PostType:
		if e := submitPost(bi, &goal, transport.Content); e != nil {
			return 0, e
		}
	case object.V5ThreadVoteType:
		if e := submitThreadVote(bi, &goal, transport.Content); e != nil {
			return 0, e
		}
	case object.V5PostVoteType:
		if e := submitPostVote(bi, &goal, transport.Content); e != nil {
			return 0, e
		}
	case object.V5UserVoteType:
		if e := submitUserVote(bi, &goal, transport.Content); e != nil {
			return 0, e
		}
	default:
		return 0, boo.Newf(boo.InvalidInput,
			"content has invalid type '%s'", transport.Body.Type)
	}

	return goal, nil
}

func submitThread(bi *BoardInstance, goal *uint64, thread *object.Content) error {
	body := thread.GetBody()

	return bi.EditPack(func(p *skyobject.Pack, h *Headers) error {

		// Set goal sequence.
		*goal = p.Root().Seq + 1

		// Get thread ref.
		tRef := p.Ref(thread)

		// Ensure thread does not exist.
		if _, has := h.GetThreadPageHash(thread.GetHeader().Hash); has {
			return boo.Newf(boo.AlreadyExists,
				"Thread of ref %s already exists", tRef.String())
		}

		// Get root children pages.
		pages, e := object.GetPages(p, &object.GetPagesIn{
			RootPage:  false,
			BoardPage: true,
			DiffPage:  true,
			UsersPage: true,
		})
		if e != nil {
			return e
		}

		// Add thread to board page.
		if e := pages.BoardPage.AddThread(tRef); e != nil {
			return e
		}

		// Save changes.
		return addContentToDiffAndProfile(p, h, pages, thread, body.Creator)
	})
}

func submitPost(bi *BoardInstance, goal *uint64, post *object.Content) error {
	body := post.GetBody()

	return bi.EditPack(func(p *skyobject.Pack, h *Headers) error {

		// Set goal sequence.
		*goal = p.Root().Seq + 1

		// Get post ref.
		pRef := p.Ref(post)

		// Ensure thread exists.
		tpHash, has := h.GetThreadPageHash(body.OfThread)
		if !has {
			return boo.Newf(boo.NotFound,
				"thread of hash '%s' not found", body.OfThread)
		}

		// Get root pages.
		pages, e := object.GetPages(p, &object.GetPagesIn{
			RootPage:  false,
			BoardPage: true,
			DiffPage:  true,
			UsersPage: true,
		})
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

		// Save changes.
		return addContentToDiffAndProfile(p, h, pages, post, body.Creator)
	})
}

func submitThreadVote(bi *BoardInstance, goal *uint64, tVote *object.Content) error {
	body := tVote.GetBody()

	return bi.EditPack(func(p *skyobject.Pack, h *Headers) error {
		*goal = p.Root().Seq + 1

		// Check vote.
		if _, ok := h.GetThreadPageHash(body.OfThread); !ok {
			return boo.Newf(boo.NotFound,
				"thread of hash '%s' is not found", body.OfThread)
		}

		return addVoteToDiffAndProfile(p, h, tVote, body.Creator)
	})
}

func submitPostVote(bi *BoardInstance, goal *uint64, pVote *object.Content) error {
	body := pVote.GetBody()

	return bi.EditPack(func(p *skyobject.Pack, h *Headers) error {
		*goal = p.Root().Seq + 1
		return addVoteToDiffAndProfile(p, h, pVote, body.Creator)
	})
}

func submitUserVote(bi *BoardInstance, goal *uint64, uVote *object.Content) error {
	body := uVote.GetBody()
	return bi.EditPack(func(p *skyobject.Pack, h *Headers) error {
		*goal = p.Root().Seq + 1
		return addVoteToDiffAndProfile(p, h, uVote, body.Creator)
	})
}

func addContentToDiffAndProfile(p *skyobject.Pack, h *Headers,
	pages *object.Pages, content *object.Content, creator string,
) error {

	// Get profile page (create if not exist).
	profileHash, ok := h.GetUserProfileHash(creator)
	if !ok {
		var e error
		if profileHash, e = pages.UsersPage.NewUserProfile(creator); e != nil {
			return e
		}
		h.SetUser(creator, profileHash)
	}

	// Add content to appropriate user profile.
	newProfileHash, e := pages.UsersPage.AddUserSubmission(profileHash, content)
	if e != nil {
		return e
	}
	h.SetUser(creator, newProfileHash)

	// Add to diff page.
	if e := pages.DiffPage.Add(content); e != nil {
		return e
	}
	return pages.Save(p)
}

func addVoteToDiffAndProfile(p *skyobject.Pack, h *Headers, content *object.Content, creator string) error {

	// Get root children pages.
	pages, e := object.GetPages(p, &object.GetPagesIn{
		RootPage:  false,
		BoardPage: false,
		DiffPage:  true,
		UsersPage: true,
	})
	if e != nil {
		return e
	}

	return addContentToDiffAndProfile(p, h, pages, content, creator)
}

func (bi *BoardInstance) EnsureSubmissionKeys(subKeyTrans []*object.MessengerSubKeyTransport) (uint64, error) {
	bi.l.Println("ensuring submission keys as:", subKeyTrans)
	return bi.EditBoard(func(board *object.Content) (bool, error) {
		body := board.GetBody()
		body.SetSubKeys(subKeyTrans)
		board.SetBody(body)
		return true, nil
	})
}

func (bi *BoardInstance) GetSubmissionKeys() []*object.MessengerSubKeyTransport {
	var subKeys []*object.MessengerSubKeyTransport
	if e := bi.ViewBoard(func(board *object.Content) (bool, error) {
		subKeys = board.GetBody().GetSubKeys()
		return false, nil
	}); e != nil {
		bi.l.Println("error obtaining submission keys:", e)
		return nil
	}
	return subKeys
}

// BoardAction is a function in which board modification/viewing takes place.
// Returns a boolean that represents whether changes have been made and
// an error on failure.
type BoardAction func(board *object.Content) (bool, error)

// EditBoard triggers a board action.
func (bi *BoardInstance) EditBoard(action BoardAction) (uint64, error) {
	var goalSeq uint64
	e := bi.EditPack(func(p *skyobject.Pack, h *Headers) error {

		// Set goal seq.
		goalSeq = p.Root().Seq + 1

		// Get root children.
		pages, e := object.GetPages(p, &object.GetPagesIn{
			RootPage:  false,
			BoardPage: true,
			DiffPage:  false,
			UsersPage: false,
		})
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

func (bi *BoardInstance) ViewBoard(action BoardAction) error {
	return bi.ViewPack(func(p *skyobject.Pack, h *Headers) error {
		// Get root children.
		pages, e := object.GetPages(p, &object.GetPagesIn{
			RootPage:  false,
			BoardPage: true,
			DiffPage:  false,
			UsersPage: false,
		})
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
