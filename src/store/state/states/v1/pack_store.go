package v1

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/io"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"log"
	"sync"
)

type PackStore struct {
	l *log.Logger

	sync.Mutex
	pack *skyobject.Pack
	seq  uint64

	rss *RefsStore
	gcs *GotStore
	tvs *ContentVotesStore
	pvs *ContentVotesStore
	uvs *UserVotesStore
	fps *FollowPageStore
}

func NewPackStore(logger *log.Logger) *PackStore {
	return &PackStore{
		l:   logger,
		rss: NewRefsStore(),
		gcs: NewGotStore(),
		tvs: NewContentVotesStore("thread"),
		pvs: NewContentVotesStore("post"),
		uvs: NewUserVotesStore(),
		fps: NewFollowPageStore(),
	}
}

func (s *PackStore) Extract(pack *skyobject.Pack, seq uint64, full, recordChanges bool) (*io.Changes, error) {
	s.Lock()
	defer s.Unlock()

	// Prepare changes.
	var changes = io.NewChanges(pack.Root().Pub, recordChanges)

	// Extract direct root children.
	rootChildren, e := pack.RootRefs()
	if e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead,
			"failed to extract root")
	}
	if len(rootChildren) != countRootRefs {
		return nil, boo.Newf(boo.InvalidRead,
			"root has invalid ref count of %d when expecting %d",
			len(rootChildren), countRootRefs)
	}

	// Process Deleted.
	deleted, e := rootChildren[indexDeleted].(*object.Deleted)
	if e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead,
			"root child 'Deleted' is invalid")
	}
	s.processDeleted(deleted, changes)

	// Return if not expected to process full lot.
	if full {
		// Process Thread Votes.
		tvp, ok := rootChildren[indexThreadVotes].(*object.ThreadVotesPages)
		if !ok {
			return nil, boo.New(boo.InvalidRead,
				"root child 'ThreadVotesPages' is invalid")
		}
		s.processContentVotesPages(tvp.Threads, s.tvs, changes)

		// Process Post Votes.
		pvp, ok := rootChildren[indexPostVotes].(*object.PostVotesPages)
		if !ok {
			return nil, boo.New(boo.InvalidRead,
				"root child 'PostVotesPages' is invalid")
		}
		s.processContentVotesPages(pvp.Posts, s.pvs, changes)

		// Process User Votes.
		uvp, ok := rootChildren[indexUserVotes].(*object.UserVotesPages)
		if !ok {
			return nil, boo.New(boo.InvalidRead,
				"root child 'UserVotesPages' is invalid")
		}
		s.processUserVotesPages(uvp.Users, s.uvs)

		// Process Content.
		tpp, ok := rootChildren[indexContent].(*object.ThreadPages)
		if !ok {
			return nil, boo.New(boo.InvalidRead,
				"root child 'ThreadPages' is invalid")
		}
		if e := s.processContent(tpp, changes); e != nil {
			return nil, boo.WrapType(e, boo.InvalidRead,
				"root child 'ThreadPages' is corrupt")
		}

		// Store pack and seq.
		s.pack, s.seq = pack, seq
	}

	return changes, nil
}

func (s *PackStore) Seq() uint64 {
	s.Lock()
	defer s.Unlock()
	return s.seq
}

func (s *PackStore) Run(action func(pack *skyobject.Pack, seq uint64) error) error {
	s.Lock()
	defer s.Unlock()
	return action(s.pack, s.seq)
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func (s *PackStore) getThreadPages() (*object.ThreadPages, error) {
	tpsObj, e := s.pack.RefByIndex(indexContent)
	if e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead,
			"failed to obtain thread pages")
	}
	tps, e := tpsObj.(*object.ThreadPages)
	if e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead,
			"failed to obtain thread pages")
	}
	return tps, nil
}

func (s *PackStore) processDeleted(
	deleted *object.Deleted,
	changes *io.Changes,
) {
	for _, tRef := range deleted.Threads {
		changes.RecordDeleteThread(tRef)
		s.rss.Delete(tRef)
		s.tvs.Delete(tRef)
		s.gcs.DeleteThread(tRef)
	}
	for _, pRef := range deleted.Posts {
		tRef := s.gcs.GetPostOrigin(pRef)
		changes.RecordDeletePost(tRef, pRef)
		s.pvs.Delete(pRef)
		s.gcs.DeletePost(pRef)
	}
}

func (s *PackStore) processContentVotesPages(
	pages []object.ContentVotesPage,
	store *ContentVotesStore,
	changes *io.Changes,
) {
	for _, page := range pages {
		// Find hash of content votes page.
		newHash := cipher.SumSHA256(encoder.Serialize(page))

		// See if it already exists.
		vs, e := store.Get(page.Ref)

		// If exists and has the same hash -> continue.
		if e == nil && vs.Hash == newHash {
			continue
		}

		vs = s.generateVotesSummary(page.Votes, nil, nil)
		vs.R = page.Ref
		vs.Hash = newHash
		store.Set(page.Ref, vs)

		// Report changes if content already exists.
		if e == nil {
			switch store.what {
			case "thread":
				changes.RecordThreadVoteChanges(page.Ref, vs)

			case "post":
				tRef := s.gcs.GetPostOrigin(page.Ref)
				changes.RecordPostVoteChanges(tRef, vs)
			}
		}
	}
}

func (s *PackStore) processUserVotesPages(
	pages []object.UserVotesPage, store *UserVotesStore,
) {
	for _, page := range pages {
		store.Set(page.Ref, s.generateVotesSummary(
			page.Votes,
			// Up action.
			func(v *object.Vote) {
				fp := s.fps.Modify(v.Creator)
				fp.Lock()
				defer fp.Lock()
				fp.Yes[page.Ref.Hex()] = object.Tag{
					Mode: "+1",
					Text: string(v.Tag),
				}
			},
			// Down action.
			func(v *object.Vote) {
				fp := s.fps.Modify(v.Creator)
				fp.Lock()
				defer fp.Unlock()
				fp.No[page.Ref.Hex()] = object.Tag{
					Mode: "-1",
					Text: string(v.Tag),
				}
			},
		))
	}
}

func (s *PackStore) generateVotesSummary(
	vRefs skyobject.Refs,
	upAction func(v *object.Vote),
	downAction func(v *object.Vote),
) *object.VotesSummary {

	summary := &object.VotesSummary{
		Votes: make(map[cipher.PubKey]object.Vote),
	}
	if l, _ := vRefs.Len(); l == 0 {
		return summary
	}
	e := vRefs.Ascend(func(i int, ref *skyobject.Ref) error {
		obj, e := ref.Value()
		if e != nil {
			s.l.Printf("failed to obtain value of vote [%d]%s",
				i, ref.String())
			return nil
		}
		vote, e := obj.(*object.Vote)
		if e != nil {
			s.l.Printf("object of ref [%d]%s is not a vote",
				i, ref.String())
			return nil
		}
		if e := vote.Verify(); e != nil {
			log.Printf("failed to verify vote of ref [%d]%s : %#v",
				i, ref.String(), vote)
			return nil
		}
		summary.Votes[vote.Creator] = *vote
		switch vote.Mode {
		case +1:
			if upAction != nil {
				upAction(vote)
			}
			summary.Up += 1
		case -1:
			if downAction != nil {
				downAction(vote)
			}
			summary.Down += 1
		}
		return nil
	})
	if e != nil {
		s.l.Println("Error when generating vote summary:", e)
	}
	return summary
}

func (s *PackStore) processContent(
	tPages *object.ThreadPages,
	changes *io.Changes,
) error {
	return tPages.ThreadPages.Ascend(func(_ int, tPageRef *skyobject.Ref) error {
		tPageVal, e := tPageRef.Value()
		if e != nil {
			return e
		}
		tPage, e := tPageVal.(*object.ThreadPage)
		if e != nil {
			return e
		}
		return tPage.Posts.Ascend(func(_ int, pRef *skyobject.Ref) error {

			// Report changes if this is new.
			if s.gcs.Get(tPage.Thread.Hash, pRef.Hash) == false {

				// Add to got store.
				s.gcs.Set(tPage.Thread.Hash, pRef.Hash)

				// Get Post.
				postVal, e := pRef.Value()
				if e != nil {
					return e
				}
				post, e := postVal.(*object.Content)
				if e != nil {
					return e
				}

				// Record changes.
				changes.RecordNewPost(tPage.Thread.Hash, post)
			}
			return nil
		})
		return nil
	})
}
