package v1

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/obtain"
	"github.com/skycoin/bbs/src/store/io"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"log"
	"sync"
)

/*
	<<< CONTENT VOTES STORE >>>
*/

// ContentVotesStore stores the votes of content (i.e. Threads, Posts).
type ContentVotesStore struct {
	sync.Mutex
	what      string
	pagesHash cipher.SHA256
	store     map[cipher.SHA256]*object.VotesSummary
}

func NewContentVotesStore(
	oldStore *ContentVotesStore,
	contentName string,
	pagesHash cipher.SHA256,
	pages []object.ContentVotesPage,
	changes *io.Changes,
) (*ContentVotesStore, error) {

	// If old votes pages is the same, return old.
	if oldStore != nil && oldStore.pagesHash == pagesHash {
		return oldStore, nil
	}

	newStore := &ContentVotesStore{
		what:      contentName,
		pagesHash: pagesHash,
		store:     make(map[cipher.SHA256]*object.VotesSummary),
	}
	for _, votesPage := range pages {

		votesPageHash := cipher.SumSHA256(encoder.Serialize(votesPage))

		// If no changes from old votes page, copy over and continue.
		var vs *object.VotesSummary
		if oldStore != nil {
			vs, _ = oldStore.Get(votesPage.Ref)
		}
		if vs != nil && vs.Hash == votesPageHash {
			goto SaveVoteSummary
		}

		vs = generateSummary(&votesPage.Votes, nil, nil)
		vs.OfContent = votesPage.Ref
		vs.Hash = votesPageHash

		// Record changes in vote summary.
		switch contentName {
		case namePost:
			changes.RecordPostVoteChanges(votesPage.Ref, vs)
		case nameThread:
			changes.RecordThreadVoteChanges(votesPage.Ref, vs)
		}

	SaveVoteSummary:
		newStore.store[votesPage.Ref] = vs
	}

	return newStore, nil
}

// Get gets content of reference.
func (s *ContentVotesStore) Get(ref cipher.SHA256) (*object.VotesSummary, error) {
	s.Lock()
	defer s.Unlock()

	vs, has := s.store[ref]
	if !has {
		return nil, boo.Newf(boo.NotFound,
			"cannot find %s of reference %s", s.what, ref.Hex())
	}
	return vs, nil
}

// Set sets content of reference.
func (s *ContentVotesStore) Set(ref cipher.SHA256, vs *object.VotesSummary) {
	s.Lock()
	defer s.Unlock()

	s.store[ref] = vs
}

// Delete removes content of reference.
func (s *ContentVotesStore) Delete(ref cipher.SHA256) {
	s.Lock()
	defer s.Unlock()

	delete(s.store, ref)
}

/*
	<<< USER VOTES STORE >>>
*/

// UserVotesStore stores the votes of users.
type UserVotesStore struct {
	sync.Mutex
	pagesHash cipher.SHA256
	store     map[cipher.PubKey]*object.VotesSummary
}

// NewUserVotesStore creates a new UserVotesStore.
func NewUserVotesStore(
	oldStore *UserVotesStore,
	pagesHash cipher.SHA256,
	pages []object.UserVotesPage,
	onUp, onDown func(v *object.Vote),
	_ *io.Changes,
) (*UserVotesStore, error) {

	// If old votes pages is the same, return old.
	if oldStore != nil && oldStore.pagesHash == pagesHash {
		return oldStore, nil
	}

	newStore := &UserVotesStore{
		pagesHash: pagesHash,
		store:     make(map[cipher.PubKey]*object.VotesSummary),
	}
	for _, votesPage := range pages {

		votesPageHash := cipher.SumSHA256(encoder.Serialize(votesPage))

		// If no changes from old votes page, copy over and continue.
		var vs *object.VotesSummary
		if oldStore != nil {
			vs, _ = oldStore.Get(votesPage.Ref)
		}
		if vs != nil && vs.Hash == votesPageHash {
			goto SaveVotesSummary
		}

		vs = generateSummary(&votesPage.Votes, onUp, onDown)
		vs.OfUser = votesPage.Ref
		vs.Hash = votesPageHash

	SaveVotesSummary:
		newStore.store[votesPage.Ref] = vs
	}

	return newStore, nil
}

// Get obtains votes of user public key.
func (s *UserVotesStore) Get(pk cipher.PubKey) (*object.VotesSummary, error) {
	s.Lock()
	defer s.Unlock()

	vs, has := s.store[pk]
	if !has {
		return nil, boo.Newf(boo.NotFound,
			"cannot find user of public key %s", pk.Hex())
	}
	return vs, nil
}

// Set sets user votes of user public key.
func (s *UserVotesStore) Set(pk cipher.PubKey, vs *object.VotesSummary) {
	s.Lock()
	defer s.Unlock()

	s.store[pk] = vs
}

// Delete removes votes of user public key.
func (s *UserVotesStore) Delete(pk cipher.PubKey) {
	s.Lock()
	defer s.Unlock()

	delete(s.store, pk)
}

/*
	<<< FOLLOW PAGE STORE >>
*/

type FollowPageStore struct {
	sync.Mutex
	store map[cipher.PubKey]*object.FollowPage
}

func NewFollowPageStore() *FollowPageStore {
	return &FollowPageStore{
		store: make(map[cipher.PubKey]*object.FollowPage),
	}
}

func (s *FollowPageStore) Get(pk cipher.PubKey) (*object.FollowPage, error) {
	s.Lock()
	defer s.Unlock()

	fp, has := s.store[pk]
	if !has {
		return nil, boo.Newf(boo.NotFound,
			"cannot find follow page for user of public key %s", pk.Hex())
	}
	return fp, nil
}

func (s *FollowPageStore) Modify(pk cipher.PubKey) *object.FollowPage {
	s.Lock()
	defer s.Unlock()

	fp, has := s.store[pk]
	if !has {
		fp = &object.FollowPage{
			UserPubKey: pk.Hex(),
			Yes:        make(map[string]object.Tag),
			No:         make(map[string]object.Tag),
		}
		s.store[pk] = fp
	}
	return fp
}

/*
	<<< GOT STORE >>>
*/

// Represents a thread and associated references.
type GotThread struct {
	sync.Mutex
	tPageHash  cipher.SHA256          // Reference of thread page that contains the thread.
	PostHashes map[cipher.SHA256]bool // Record of post references of the thread.
}

// NewGotThread creates a new GotThread, recording changes (if any with the old GotThread).
// Input   'oldGT' can be nil (if there is not older version of this GotThread).
// Input 'changes' can be nil (if there is not need to record changes).
func NewGotThread(
	oldGT *GotThread,
	tPageHash cipher.SHA256,
	tPage *object.ThreadPage,
	changes *io.Changes,
	postOrigins map[cipher.SHA256]cipher.SHA256,
) (*GotThread, error) {

	// If old ThreadPage's hash is the same as new;
	// no changes needed - just copy.
	if oldGT != nil && oldGT.tPageHash == tPageHash {
		return oldGT, nil
	}

	// Create new GotThread.
	newGT := &GotThread{
		tPageHash:  tPageHash,
		PostHashes: make(map[cipher.SHA256]bool),
	}
	e := tPage.Posts.Ascend(func(_ int, postRef *skyobject.Ref) error {
		// If old thread does not have this post, record changes.
		if oldGT.HasPost(postRef.Hash) {
			post, e := obtain.Content(postRef)
			if e != nil {
				return e
			}
			changes.RecordNewPost(tPage.Thread.Hash, post)
		}
		// Record post in GotThread.
		newGT.PostHashes[postRef.Hash] = true
		// Record post origins (ref of thread that post came from).
		postOrigins[postRef.Hash] = tPage.Thread.Hash
		return nil
	})
	if e != nil {
		return nil, e
	}
	return newGT, nil
}

func (g *GotThread) Do(action func(g *GotThread) error) error {
	g.Lock()
	defer g.Unlock()
	return action(g)
}

func (g *GotThread) HasPost(postHash cipher.SHA256) bool {
	g.Lock()
	defer g.Unlock()
	return g.PostHashes[postHash]
}

// GotStore keeps the posts that we have, per thread,
// and the thread page reference of each thread.
type GotStore struct {
	sync.Mutex
	tPagesHash cipher.SHA256                   // Hash of the entirety of the ThreadPages.
	threads    map[cipher.SHA256]*GotThread    // Tells what posts are in a thread.
	posts      map[cipher.SHA256]cipher.SHA256 // Tells which thread, post originates from.
}

// NewGotStore creates a new GotStore.
// Input   'oldGS' can be nil (if there is no previous GotStore).
// Input 'changes' can be nil (if changes don't need to be recorded).
func NewGotStore(
	oldGS *GotStore,
	tPagesHash cipher.SHA256,
	tPages *object.ThreadPages,
	changes *io.Changes,
) (*GotStore, error) {
	// If no changes, return old GotStore.
	if oldGS != nil && oldGS.tPagesHash == tPagesHash {
		return oldGS, nil
	}
	newGS := &GotStore{
		tPagesHash: tPagesHash,
		threads:    make(map[cipher.SHA256]*GotThread),
		posts:      make(map[cipher.SHA256]cipher.SHA256),
	}
	e := tPages.ThreadPages.Ascend(func(_ int, tPageRef *skyobject.Ref) error {
		tPage, e := obtain.ThreadPage(tPageRef)
		if e != nil {
			return e
		}
		oldGotThread, has := oldGS.GetThread(tPage.Thread.Hash)
		if !has {
			// TODO: Changes - record new thread.
		}
		newGotThread, e := NewGotThread(
			oldGotThread, tPageRef.Hash, tPage, changes, newGS.posts)
		if e != nil {
			return e
		}
		newGS.threads[tPage.Thread.Hash] = newGotThread
		return nil
	})
	if e != nil {
		return nil, e
	}
	return newGS, nil
}

func (s *GotStore) GetThread(tHash cipher.SHA256) (*GotThread, bool) {
	s.Lock()
	defer s.Unlock()
	gt, has := s.threads[tHash]
	return gt, has
}

func (s *GotStore) GetPostOrigin(pHash cipher.SHA256) cipher.SHA256 {
	s.Lock()
	defer s.Unlock()
	return s.posts[pHash]
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func generateSummary(
	voteRefs *skyobject.Refs,
	onUp, onDown func(v *object.Vote),
) *object.VotesSummary {

	summary := &object.VotesSummary{
		Votes: make(map[cipher.PubKey]object.Vote),
	}
	if len, _ := voteRefs.Len(); len == 0 {
		return summary
	}
	e := voteRefs.Ascend(func(i int, voteRef *skyobject.Ref) error {
		vote, e := obtain.Vote(voteRef)
		if e != nil {
			return boo.Wrapf(e, "failed on vote of index %d", i)
		}
		if e := vote.Verify(); e != nil {
			log.Printf("failed to verify vote of ref [%d]%s : %#v",
				i, voteRef.String(), vote)
			return nil
		}
		summary.Votes[vote.Creator] = *vote
		switch vote.Mode {
		case +1:
			if onUp != nil {
				onUp(vote)
			}
			summary.Up += 1
		case -1:
			if onDown != nil {
				onDown(vote)
			}
			summary.Down += 1
		}
		return nil
	})
	if e != nil {
		log.Println("Error when generating vote summary:", e)
	}
	return summary
}
