package state

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
)

/*
	<<< HEADERS >>>
*/

type Headers struct {
	rootSeq uint64
	changes *object.Changes

	tMux    sync.Mutex
	threads map[string]cipher.SHA256 // key(thread content hash), value(thread page hash)

	uMux  sync.Mutex
	users map[string]cipher.SHA256 // key(user's public key), value(user profile hash)
}

func NewHeaders(oldHeaders *Headers, p *skyobject.Pack) (*Headers, error) {
	if len(p.Root().Refs) != object.RootChildrenCount {
		return nil, boo.New(boo.InvalidRead,
			"invalid root")
	}
	headers := &Headers{
		threads: make(map[string]cipher.SHA256),
		users:   make(map[string]cipher.SHA256),
	}

	// Get required root children.
	pages, e := object.GetPages(p, &object.GetPagesIn{
		RootPage:  false,
		BoardPage: true,
		DiffPage:  true,
		UsersPage: true,
	})
	if e != nil {
		return nil, e
	}

	// Fill threads header data.
	e = pages.BoardPage.Threads.Ascend(func(i int, tpElem *skyobject.RefsElem) error {
		tp, e := object.GetThreadPage(tpElem)
		if e != nil {
			return e
		}
		t, e := tp.GetThread()
		if e != nil {
			return e
		}
		headers.threads[t.GetHeader().Hash] = tpElem.Hash
		return nil
	})
	if e != nil {
		return nil, e
	}

	// Fill users header data.
	e = pages.UsersPage.Users.Ascend(func(i int, uapElem *skyobject.RefsElem) error {
		uap, e := object.GetUserProfile(uapElem)
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
	var oldChanges *object.Changes
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

func (h *Headers) GetChanges() *object.Changes {
	return h.changes
}

func (h *Headers) GetThreadPageHash(threadHash string) (cipher.SHA256, bool) {
	h.tMux.Lock()
	defer h.tMux.Unlock()
	tpHash, has := h.threads[threadHash]
	return tpHash, has
}

func (h *Headers) GetUserProfileHash(UserPubKey string) (cipher.SHA256, bool) {
	h.uMux.Lock()
	defer h.uMux.Unlock()
	uapHash, has := h.users[UserPubKey]
	return uapHash, has
}

func (h *Headers) SetUser(upk string, profileHash cipher.SHA256) {
	h.uMux.Lock()
	defer h.uMux.Unlock()
	h.users[upk] = profileHash
}

func (h *Headers) SetThread(threadHash string, tpRef cipher.SHA256) {
	h.tMux.Lock()
	defer h.tMux.Unlock()

	h.threads[threadHash] = tpRef
}

// RangeThreadFunc is the function used to range the threads.
// Quits range on error.
type RangeThreadFunc func(threadHash string, tpHash cipher.SHA256) error

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
