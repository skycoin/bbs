package state

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
)

type PackAction func(p *skyobject.Pack, h *PackHeaders) error

type PackInstance struct {
	pack    *skyobject.Pack
	headers *PackHeaders
}

func NewPackInstance(oldPI *PackInstance, pack *skyobject.Pack) (*PackInstance, error) {
	newPI := &PackInstance{pack: pack}
	var e error
	newPI.headers, e = NewPackInstanceHeaders(oldPI.headers, pack)
	if e != nil {
		return nil, e
	}
	return newPI, nil
}

func (pi *PackInstance) Do(action PackAction) error {
	return action(pi.pack, pi.headers)
}

/*
	<<< HEADERS >>>
*/

type PackHeaders struct {
	rootSeq uint64
	changes *object.Changes

	tMux    sync.Mutex
	threads map[cipher.SHA256]cipher.SHA256

	uMux  sync.Mutex
	users map[cipher.PubKey]cipher.SHA256
}

func NewPackInstanceHeaders(oldHeaders *PackHeaders, p *skyobject.Pack) (*PackHeaders, error) {
	if len(p.Root().Refs) != object.RootChildrenCount {
		return nil, boo.New(boo.InvalidRead,
			"invalid root")
	}
	headers := &PackHeaders{
		threads: make(map[cipher.SHA256]cipher.SHA256),
		users:   make(map[cipher.PubKey]cipher.SHA256),
	}

	// Get required root children.
	bPage, dPage, uPage, e := object.GetPages(p, nil,
		true, true, true)
	if e != nil {
		return nil, e
	}

	// Fill threads header data.
	e = bPage.Threads.Ascend(func(i int, tpRef *skyobject.Ref) error {
		tp, e := object.GetThreadPage(tpRef, nil)
		if e != nil {
			return e
		}
		headers.threads[tp.Thread.Hash] = tpRef.Hash
		return nil
	})
	if e != nil {
		return nil, e
	}

	// Fill users header data.
	e = uPage.Users.Ascend(func(i int, uapRef *skyobject.Ref) error {
		uap, e := object.GetUserActivityPage(uapRef, nil)
		if e != nil {
			return e
		}
		headers.users[uap.PubKey] = uapRef.Hash
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
	headers.changes, e = dPage.GetChanges(oldChanges, nil)
	if e != nil {
		return nil, e
	}

	return headers, nil
}

func (h *PackHeaders) GetRootSeq() uint64 {
	return h.rootSeq
}

func (h *PackHeaders) GetChanges() *object.Changes {
	return h.changes
}

func (h *PackHeaders) GetThreadPageHash(threadHash cipher.SHA256) (cipher.SHA256, bool) {
	h.tMux.Lock()
	defer h.tMux.Unlock()
	tpHash, has := h.threads[threadHash]
	return tpHash, has
}

func (h *PackHeaders) GetUserActivityPageHash(UserPubKey cipher.PubKey) (cipher.SHA256, bool) {
	h.uMux.Lock()
	defer h.uMux.Unlock()
	uapHash, has := h.users[UserPubKey]
	return uapHash, has
}
