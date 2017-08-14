package v1

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
)

/*
	<<< CONTENT VOTES STORE >>>
*/

// ContentVotesStore stores the votes of content (i.e. Threads, Posts).
type ContentVotesStore struct {
	sync.Mutex
	what  string
	store map[cipher.SHA256]*object.VotesSummary
}

// NewContentVotesStore creates a new ContentVotesStore.
func NewContentVotesStore(what string) *ContentVotesStore {
	return &ContentVotesStore{
		what:  what,
		store: make(map[cipher.SHA256]*object.VotesSummary),
	}
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
	store map[cipher.PubKey]*object.VotesSummary
}

// NewUserVotesStore creates a new UserVotesStore.
func NewUserVotesStore() *UserVotesStore {
	return &UserVotesStore{
		store: make(map[cipher.PubKey]*object.VotesSummary),
	}
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

// GotStore keeps the posts that we have, per thread.
type GotStore struct {
	sync.Mutex
	Threads map[cipher.SHA256]map[cipher.SHA256]bool // Tells what posts are in a thread.
	Posts   map[cipher.SHA256]cipher.SHA256          // Tells which thread, post originates from.
}

func NewGotStore() *GotStore {
	return &GotStore{
		Threads: make(map[cipher.SHA256]map[cipher.SHA256]bool),
		Posts:   make(map[cipher.SHA256]cipher.SHA256),
	}
}

func (s *GotStore) Set(tRef, pRef cipher.SHA256) {
	s.Lock()
	defer s.Unlock()

	pMap, has := s.Threads[tRef]
	if !has {
		pMap = make(map[cipher.SHA256]bool)
		s.Threads[tRef] = pMap
	}
	pMap[pRef] = true
	s.Posts[pRef] = tRef
}

func (s *GotStore) Get(tRef, pRef cipher.SHA256) bool {
	s.Lock()
	defer s.Unlock()

	pMap, has := s.Threads[tRef]
	if !has {
		return false
	}
	return pMap[pRef]
}

func (s *GotStore) GetPostOrigin(pRef cipher.SHA256) cipher.SHA256 {
	s.Lock()
	defer s.Unlock()

	return s.Posts[pRef]
}

func (s *GotStore) DeleteThread(tRef cipher.SHA256) {
	s.Lock()
	defer s.Unlock()

	pMap, has := s.Threads[tRef]
	if !has {
		return
	}

	for pRef := range pMap {
		delete(s.Posts, pRef)
	}

	delete(s.Threads, tRef)
}

func (s *GotStore) DeletePost(pRef cipher.SHA256) {
	s.Lock()
	defer s.Unlock()

	tRef, has := s.Posts[pRef]
	if !has {
		return
	}

	tMap, has := s.Threads[tRef]
	if !has {
		return
	}

	delete(tMap, pRef)
}
