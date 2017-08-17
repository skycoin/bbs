package state

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/obtain"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
)

type PackAction func(
	p *skyobject.Pack,
	h *ActivityHeader,
) error

type PackInstance struct {
	mux    sync.Mutex
	pack   *skyobject.Pack
	header *ActivityHeader
}

func NewPackInstance(oldPI *PackInstance, pack *skyobject.Pack) (*PackInstance, *Changes, error) {
	header, e := NewPackInstanceHeader(pack)
	if e != nil {
		return nil, nil, e
	}
	pi := &PackInstance{
		pack:   pack,
		header: header,
	}
	var changes *Changes
	if oldPI != nil {
		if changes, e = NewChanges(oldPI, pi); e != nil {
			return nil, nil, e
		}
	}
	return pi, changes, nil
}

func (pi *PackInstance) Do(action PackAction) error {
	pi.mux.Lock()
	defer pi.mux.Unlock()
	return action(pi.pack, pi.header)
}

/*
	<<< HEADERS >>>
*/

type ActivityHeader struct {
	sync.Mutex
	store map[cipher.PubKey]cipher.SHA256
}

func NewPackInstanceHeader(p *skyobject.Pack) (*ActivityHeader, error) {
	if len(p.Root().Refs) != 3 {
		return nil, boo.New(boo.InvalidRead, "invalid root") // TODO: Rename.
	}
	header := &ActivityHeader{
		store: make(map[cipher.PubKey]cipher.SHA256),
	}

	// Fill activity header.
	aPage, e := object.GetActivityPage(p, nil)
	if e != nil {
		return nil, e
	}
	e = aPage.Users.Ascend(func(i int, uaRef *skyobject.Ref) error {
		ua, e := object.GetUserActivity(uaRef)
		if e != nil {
			return e
		}
		header.store[ua.PubKey] = uaRef.Hash
		return nil
	})
	if e != nil {
		return nil, e
	}

	return header, nil
}

func (h *ActivityHeader) Len() int {
	h.Lock()
	defer h.Unlock()
	return len(h.store)
}

func (h *ActivityHeader) Get(upk cipher.PubKey) (cipher.SHA256, bool) {
	h.Lock()
	defer h.Unlock()
	uaHash, has := h.store[upk]
	return uaHash, has
}

func (h *ActivityHeader) Add(upk cipher.PubKey, uaHash cipher.SHA256) {
	h.Lock()
	defer h.Unlock()
	h.store[upk] = uaHash
}

/*
	<<< CHANGES >>>
*/

type Changes struct {
	NewContent     []*object.Content
	DeletedContent []cipher.SHA256
	NewVotes       []*object.Vote
}

func NewChanges(oldPI, newPI *PackInstance) (*Changes, error) {
	changes := new(Changes)

	// Get old data.
	var oldContentLen, oldDeletedLen, oldVotesLen int
	if e := oldPI.Do(func(p *skyobject.Pack, h *ActivityHeader) (e error) {
		cPage, e := object.GetContentPage(p, nil)
		if e != nil {
			return e
		}
		oldContentLen, e = cPage.Content.Len()
		if e != nil {
			return e
		}
		oldDeletedLen, e = cPage.Deleted.Len()
		if e != nil {
			return e
		}
		oldVotesLen, e = cPage.Votes.Len()
		if e != nil {
			return e
		}
		return nil
	}); e != nil {
		return nil, e
	}

	// Get new data.
	var newContentLen, newDeletedLen, newVotesLen int
	if e := newPI.Do(func(p *skyobject.Pack, h *ActivityHeader) error {
		cPage, e := object.GetContentPage(p, nil)
		if e != nil {
			return e
		}
		newContentLen, e = cPage.Content.Len()
		if e != nil {
			return e
		}
		newDeletedLen, e = cPage.Deleted.Len()
		if e != nil {
			return e
		}
		newVotesLen, e = cPage.Votes.Len()
		if e != nil {
			return e
		}

		// Append to changes.
		if newContentLen > oldContentLen {
			changes.NewContent = make([]*object.Content, newContentLen-oldContentLen)
			for i := oldContentLen; i < newContentLen; i++ {
				ref, e := cPage.Content.RefBiIndex(i)
				if e != nil {
					return e
				}
				changes.NewContent[i], e = obtain.Content(ref)
				if e != nil {
					return e
				}
			}
		}
		if newVotesLen > oldVotesLen {
			changes.NewVotes = make([]*object.Vote, newVotesLen-oldVotesLen)
			for i := oldVotesLen; i < newVotesLen; i++ {
				ref, e := cPage.Votes.RefBiIndex(i)
				if e != nil {
					return e
				}
				changes.NewVotes[i], e = obtain.Vote(ref)
				if e != nil {
					return e
				}
			}
		}
		if newDeletedLen > oldDeletedLen {
			changes.DeletedContent = make([]cipher.SHA256, newDeletedLen-oldDeletedLen)
			for i := oldDeletedLen; i < newDeletedLen; i++ {
				ref, e := cPage.Deleted.RefBiIndex(i)
				if e != nil {
					return e
				}
				changes.DeletedContent[i] = ref.Hash
			}
		}
		return nil
	}); e != nil {
		return nil, e
	}

	return changes, nil
}
