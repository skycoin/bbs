package pack

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object/revisions/r0"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
)

/*
	<<< HEADERS >>>
*/

type Headers struct {
	rootSeq uint64
	changes *r0.Changes

	tMux    sync.Mutex
	threads map[cipher.SHA256]cipher.SHA256 // key(thread hash), value(thread page hash)

	uMux  sync.Mutex
	users map[cipher.PubKey]cipher.SHA256 // key(user's public key), value(user's page hash)
}

func NewHeaders(oldHeaders *Headers, p *skyobject.Pack) (*Headers, error) {
	if len(p.Root().Refs) != r0.RootChildrenCount {
		return nil, boo.New(boo.InvalidRead,
			"invalid root")
	}
	headers := &Headers{
		threads: make(map[cipher.SHA256]cipher.SHA256),
		users:   make(map[cipher.PubKey]cipher.SHA256),
	}

	// Get required root children.
	pages, e := r0.GetPages(p, false, true, true, true)
	if e != nil {
		return nil, e
	}

	// Fill threads header data.
	e = pages.BoardPage.Threads.Ascend(func(i int, tpElem *skyobject.RefsElem) error {
		tp, e := r0.GetThreadPage(tpElem)
		if e != nil {
			return e
		}
		headers.threads[tp.Thread.Hash] = tpElem.Hash
		return nil
	})
	if e != nil {
		return nil, e
	}

	// Fill users header data.
	e = pages.UsersPage.Users.Ascend(func(i int, uapElem *skyobject.RefsElem) error {
		uap, e := r0.GetUserActivityPage(uapElem)
		if e != nil {
			return e
		}
		headers.users[uap.PubKey] = uapElem.Hash
		return nil
	})
	if e != nil {
		return nil, e
	}

	// Fill initial changes object.
	var oldChanges *r0.Changes
	if oldHeaders != nil {
		oldChanges = oldHeaders.GetChanges()
	}
	headers.changes, e = pages.DiffPage.GetChanges(oldChanges)
	if e != nil {
		return nil, e
	}

	return headers, nil
}

func (h *Headers) GetRootSeq() uint64 {
	return h.rootSeq
}

func (h *Headers) GetChanges() *r0.Changes {
	return h.changes
}

func (h *Headers) GetThreadPageHash(threadHash cipher.SHA256) (cipher.SHA256, bool) {
	h.tMux.Lock()
	defer h.tMux.Unlock()
	tpHash, has := h.threads[threadHash]
	return tpHash, has
}

func (h *Headers) GetUserActivityPageHash(UserPubKey cipher.PubKey) (cipher.SHA256, bool) {
	h.uMux.Lock()
	defer h.uMux.Unlock()
	uapHash, has := h.users[UserPubKey]
	return uapHash, has
}

func (h *Headers) SetUser(upk cipher.PubKey, uapHash cipher.SHA256) {
	h.uMux.Lock()
	defer h.uMux.Unlock()
	h.users[upk] = uapHash
}

func (h *Headers) SetThread(tRef, tpRef cipher.SHA256) {
	h.tMux.Lock()
	defer h.tMux.Unlock()

	h.threads[tRef] = tpRef
}

// RangeThreadFunc is the function used to range the threads.
// Quits range on error.
type RangeThreadFunc func(tHash, tpHash cipher.SHA256) error

func (h *Headers) RangeThreads(action RangeThreadFunc) error {
	h.tMux.Lock()
	defer h.tMux.Unlock()

	for tHash, tpHash := range h.threads {
		if e := action(tHash, tpHash); e != nil {
			return e
		}
	}
	return nil
}
