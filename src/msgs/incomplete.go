package msgs

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
)

type Elem struct {
	res *NewContentResponse
	e   error
}

type Incomplete struct {
	mux   sync.Mutex
	store map[cipher.SHA256]chan *NewContentResponse
}

func NewIncomplete() *Incomplete {
	return &Incomplete{
		store: make(map[cipher.SHA256]chan *NewContentResponse),
	}
}

func (i *Incomplete) Add(hash cipher.SHA256) (chan *NewContentResponse, error) {
	i.mux.Lock()
	defer i.mux.Unlock()
	if _, has := i.store[hash]; has {
		return nil, boo.Newf(boo.AlreadyExists,
			"request %s already exists", hash.Hex())
	}
	out := make(chan *NewContentResponse, 1)
	i.store[hash] = out
	return out, nil
}

func (i *Incomplete) Remove(hash cipher.SHA256) {
	i.mux.Lock()
	defer i.mux.Unlock()
	if c, ok := i.store[hash]; ok {
		close(c)
		delete(i.store, hash)
	}
}

func (i *Incomplete) Satisfy(res *NewContentResponse) {
	// TODO: return error "nothing to satisfy"
	i.mux.Lock()
	defer i.mux.Unlock()
	if c, ok := i.store[res.Hash]; ok {
		select {
		case c <- res:
		default:
		}
	}
}
