package state

import (
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/bbs/src/misc/boo"
)

func (bi *BoardInstance) NewThread(thread *object.Thread) (uint64, error) {
	if e := thread.Verify(); e != nil {
		return 0, e
	}

	// TODO: Check user permissions.

	var goalSeq uint64
	e := bi.PackDo(func(p *skyobject.Pack, h *PackHeaders) error {

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
		pages, e := object.GetPages(p, nil, true, true, false)
		if e != nil {
			return e
		}

		// Add thread to board page.
		if e := pages.BoardPage.AddThread(tRef, nil); e != nil {
			return e
		}

		// Add thread to diff page.
		if e := pages.DiffPage.Add(thread); e != nil {
			return e
		}

		// Save changes.
		return pages.Save(p, nil)
	})

	return goalSeq, e
}

func (bi *BoardInstance) NewPost(post *object.Post) (uint64, error) {
	if e := post.Verify(); e != nil {
		return 0, e
	}

	// TODO: Check user permissions.

	var goalSeq uint64
	e := bi.PackDo(func(p *skyobject.Pack, h *PackHeaders) error {

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
		pages, e := object.GetPages(p, nil, true, true, false)
		if e != nil {
			return e
		}

		// Add post to board page.
		tpRef, tPage, e := pages.BoardPage.GetThreadPage(tpHash, nil)
		if e != nil {
			return e
		}
		if e := tPage.AddPost(pRef.Hash, post, nil); e != nil {
			return e
		}
		if e := tPage.Save(tpRef); e != nil {
			return e
		}

		// Add post to diff.
		if e := pages.DiffPage.Add(post); e != nil {
			return e
		}

		// Save changes.
		return pages.Save(p, nil)
	})

	return goalSeq, e
}