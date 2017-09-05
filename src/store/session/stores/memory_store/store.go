package memory_store

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object"
	"sync"
)

type Store struct {
	sync.Mutex
	users   map[string]*object.UserFile
	current *string
}

func NewStore() *Store {
	return &Store{
		users: make(map[string]*object.UserFile),
	}
}

func (s *Store) GetUsers() ([]string, error) {
	s.Lock()
	defer s.Unlock()

	out := make([]string, len(s.users))
	i := 0
	for alias := range s.users {
		out[i] = alias
		i += 1
	}
	return out, nil
}

func (s *Store) GetUser(alias string) (*object.UserFile, bool) {
	s.Lock()
	defer s.Unlock()

	f, ok := s.users[alias]
	return f, ok
}

func (s *Store) NewUser(alias string, file *object.UserFile) error {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.users[alias]; ok {
		return boo.Newf(boo.AlreadyExists,
			"user of alias '%s' already exists", alias)
	}

	s.users[alias] = file
	return nil
}

func (s *Store) DeleteUser(alias string) error {
	s.Lock()
	defer s.Unlock()
	delete(s.users, alias)
	return nil
}
