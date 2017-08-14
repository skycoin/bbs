package v1

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
)

type RefsStore struct {
	sync.Mutex
	Threads map[cipher.SHA256]cipher.SHA256
}

func NewRefsStore() *RefsStore {
	return &RefsStore{
		Threads: make(map[cipher.SHA256]cipher.SHA256),
	}
}

func (s *RefsStore) Get(tRef cipher.SHA256) (cipher.SHA256, error) {
	s.Lock()
	defer s.Unlock()
	tPageRef, has := s.Threads[tRef]
	if !has {
		return tPageRef, boo.Newf(boo.NotFound,
			"cannot find thread of ref '%s'", tRef.Hex())
	}
	return tPageRef, nil
}

func (s *RefsStore) Set(tRef, tPageRef cipher.SHA256) {
	s.Lock()
	defer s.Unlock()
	s.Threads[tRef] = tPageRef
}

func (s *RefsStore) Delete(tRef cipher.SHA256) {
	s.Lock()
	defer s.Unlock()
	delete(s.Threads, tRef)
}
