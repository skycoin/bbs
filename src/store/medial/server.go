package medial

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
	"time"
)

type Item struct {
	Raw []byte
	PK  cipher.PubKey
	TS  int64
}

func (i *Item) IsTimedOut(interval time.Duration) bool {
	return time.Now().UnixNano()-i.TS > int64(interval)
}

type ServerConfig struct {
	GarbageCollectionInterval time.Duration
	ItemTimeoutInterval       time.Duration
}

type Server struct {
	c     *ServerConfig
	mux   sync.Mutex
	store map[cipher.SHA256]*Item
	quit  chan struct{}
	wg    sync.WaitGroup
}

func NewServer(config *ServerConfig) *Server {
	s := &Server{
		c:     config,
		store: make(map[cipher.SHA256]*Item),
		quit:  make(chan struct{}),
	}
	go s.service()
	return s
}

func (s *Server) lock() func() {
	s.mux.Lock()
	return s.mux.Unlock
}

func (s *Server) Close() {
	if s != nil {
		defer s.lock()()
		if s.quit != nil {
			close(s.quit)
			s.wg.Wait()
		}
	}
}

func (s *Server) service() {
	s.wg.Add(1)
	defer s.wg.Done()

	gcTicker := time.NewTicker(s.c.GarbageCollectionInterval)
	defer gcTicker.Stop()

	for {
		select {
		case <-s.quit:
			return
		case <-gcTicker.C:
			s.collectGarbage()
		}
	}
}

func (s *Server) collectGarbage() {
	defer s.lock()()
	for k, v := range s.store {
		if v.IsTimedOut(s.c.ItemTimeoutInterval) {
			delete(s.store, k)
		}
	}
}

func (s *Server) Add(upk cipher.PubKey, body *object.Body) (cipher.SHA256, []byte, error) {
	hash, raw := body.ToRaw()

	defer s.lock()()
	if _, ok := s.store[hash]; ok {
		return cipher.SHA256{}, []byte{}, boo.New(boo.AlreadyExists,
			"already waiting for finalization of submission of hash %s",
			hash.Hex())
	}

	s.store[hash] = &Item{
		Raw: raw,
		TS:  body.TS,
	}
	return hash, raw, nil
}

func (s *Server) Satisfy(hash cipher.SHA256, sig cipher.Sig) ([]byte, error) {
	defer s.lock()()

	if item, ok := s.store[hash]; !ok {
		return nil, boo.Newf(boo.NotFound,
			"submission of hash %s is not found in medial", hash.Hex())
	} else {
		return item.Raw, nil
	}
}
