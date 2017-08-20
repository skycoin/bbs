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

		// Add thread.
		pages, e := object.GetPages(p, nil, true, true, false)
		if e != nil {
			return e
		}
		if e := pages.BoardPage.AddThread(tRef, nil); e != nil {
			return e
		}
		if e := pages.DiffPage.Add(thread); e != nil {
			return e
		}
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

		// Add post to thread.
		pages, e := object.GetPages(p, nil, true, true, false)
		if e != nil {
			return e
		}
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
		return pages.Save(p, nil)
	})

	return goalSeq, e
}

//func (bi *BoardInstance) NewContent(content *object.Content) (uint64, error) {
//
//	if e := content.Verify(); e != nil {
//		return 0, e
//	}
//
//	var goalSeq uint64
//	e := bi.GetPack().Do(func(p *skyobject.Pack, h *PackHeaders) error {
//
//		// Set goal sequence.
//		goalSeq = p.Root().Seq + 1
//
//		// Get content ref.
//		cRef := p.Ref(content)
//
//		// Get content page and activity page.
//		_, cPage, aPage, e := object.GetPages(p, nil)
//		if e != nil {
//			return e
//		}
//
//		// Check if content already exists.
//		if cPage.HasContent(cRef.Hash, nil) {
//			return boo.Newf(boo.AlreadyExists,
//				"Content of ref %s already exists", cRef.String())
//		}
//
//		// Add content to content page.
//		if content.IsThread() {
//			if e := cPage.Threads.Append(content); e != nil {
//				return e
//			}
//		} else {
//			// Content is post - additional checks.
//			// Make sure referred thread exists and is not deleted.
//			tHash := content.OfContent[0]
//			if !cPage.HasThread(tHash, nil) {
//				return boo.Newf(boo.NotFound,
//					"thread of hash '%s' not found", tHash.Hex())
//			}
//			if cPage.HasDeleted(tHash, nil) {
//				return boo.Newf(boo.NotAllowed,
//					"thread of hash '%s' is deleted", tHash.Hex())
//			}
//			// If post refers another post. Make sure that post is not deleted.
//			if rpHash, yes := content.RefersPost(); yes {
//				if cPage.HasDeleted(rpHash, nil) {
//					return boo.Newf(boo.NotAllowed,
//						"post of hash '%s' is deleted", rpHash.Hex())
//				}
//			}
//		}
//		if e := cPage.Content.Append(content); e != nil {
//			return e
//		}
//		if e := cPage.Save(p, nil); e != nil {
//			return e
//		}
//
//		// Add content to activity page.
//		uaHash, has := h.GetUserActivityPageHash(content.Creator)
//		if !has {
//			// Add user activity page if not exist.
//			uaHash, e = aPage.NewUserActivityPage(content.Creator)
//			if e != nil {
//				return e
//			}
//		}
//		if e := aPage.AddUserActivity(uaHash, content); e != nil {
//			return e
//		}
//		if e := aPage.Save(p, nil); e != nil {
//			return e
//		}
//
//		return nil
//	})
//	if e != nil {
//		return 0, e
//	}
//
//	// Done.
//	bi.SetUpdateNeeded()
//	return goalSeq, nil
//}
//
//func (bi *BoardInstance) NewVote(post *object.Content) error {
//	return nil
//}
//
//func (bi *BoardInstance) DeleteContent(cHash cipher.SHA256) error {
//	return nil
//}
