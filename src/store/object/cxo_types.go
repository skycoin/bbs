package object

import (
	"github.com/skycoin/bbs/src/misc/tag"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

const (
	IndexBoardPage    = 0
	IndexContentPage  = 1
	IndexActivityPage = 2
)

var indexString = [...]string{
	IndexBoardPage: "BoardPage",
	IndexContentPage: "ContentPage",
	IndexActivityPage: "ActivityPage",
}

/*
	<<< ROOT CHILDREN >>>
*/

type BoardPage struct {
	Board skyobject.Ref `skyobject:"schema=bbs.Board"`
	Meta  []byte
}

func GetPages(p *skyobject.Pack, mux *sync.Mutex, get ...bool) (
	bp *BoardPage, cp *ContentPage, ap *ActivityPage, e error,
) {
	defer dynamicLock(mux)()

	if len(get) > IndexBoardPage && get[IndexBoardPage] {
		// TODO: Get BoardPage.
	}
	if len(get) > IndexContentPage && get[IndexContentPage] {
		if cp, e = GetContentPage(p, nil); e != nil {
			return
		}
	}
	if len(get) > IndexActivityPage && get[IndexActivityPage] {
		if ap, e = GetActivityPage(p, nil); e != nil {
			return
		}
	}
	return
}

/*
	<<< CONTENT PAGE >>>
*/

type ContentPage struct {
	Threads skyobject.Refs `skyobject:"schema=bbs.Content"`
	Content skyobject.Refs `skyobject:"schema=bbs.Content"`
	Deleted skyobject.Refs `skyobject:"schema=bbs.Content"`
	Votes   skyobject.Refs `skyobject:"schema=bbs.Vote"`
}

func GetContentPage(p *skyobject.Pack, mux *sync.Mutex) (*ContentPage, error) {
	defer dynamicLock(mux)()
	cpValue, e := p.RefByIndex(IndexContentPage)
	if e != nil {
		return nil, getRootChildErr(e, IndexContentPage)
	}
	cp, ok := cpValue.(*ContentPage)
	if !ok {
		return nil, extRootChildErr(IndexContentPage)
	}
	return cp, nil
}

func (cp *ContentPage) HasContent(cHash cipher.SHA256, mux *sync.Mutex) bool {
	defer dynamicLock(mux)()
	_, e := cp.Content.RefByHash(cHash)
	return e == nil
}

func (cp *ContentPage) HasThread(tHash cipher.SHA256, mux *sync.Mutex) bool {
	defer dynamicLock(mux)()
	_, e := cp.Threads.RefByHash(tHash)
	return e == nil
}

func (cp *ContentPage) HasDeleted(cHash cipher.SHA256, mux *sync.Mutex) bool {
	defer dynamicLock(mux)()
	_, e := cp.Deleted.RefByHash(cHash)
	return e == nil
}

func (cp *ContentPage) Save(p *skyobject.Pack, mux *sync.Mutex) error {
	defer dynamicLock(mux)()
	if e := p.SetRefByIndex(IndexContentPage, cp); e != nil {
		return saveRootChildErr(e, IndexContentPage)
	}
	return nil
}

/*
	<<< ACTIVITY PAGE >>>
*/

type ActivityPage struct {
	Users skyobject.Refs `skyobject:"schema=bbs.UserActivity"`
}

func GetActivityPage(p *skyobject.Pack, mux *sync.Mutex) (*ActivityPage, error) {
	defer dynamicLock(mux)()
	apValue, e := p.RefByIndex(IndexActivityPage)
	if e != nil {
		return nil, getRootChildErr(e, IndexActivityPage)
	}
	ap, ok := apValue.(*ActivityPage)
	if !ok {
		return nil, extRootChildErr(IndexActivityPage)
	}
	return ap, nil
}

func (ap *ActivityPage) Save(p *skyobject.Pack, mux *sync.Mutex) error {
	defer dynamicLock(mux)()
	if e := p.SetRefByIndex(IndexActivityPage, ap); e != nil {
		return saveRootChildErr(e, IndexActivityPage)
	}
	return nil
}

func (ap *ActivityPage) AppendUserActivity(upk cipher.PubKey) (cipher.SHA256, error) {
	ua := &UserActivity{PubKey: upk}
	if e := ap.Users.Append(ua); e != nil {
		return cipher.SHA256{}, appendErr(e, ua, "ActivityPage.Users")
	}
	return cipher.SumSHA256(encoder.Serialize(ua)), nil
}

func (ap *ActivityPage) AddToUserActivity(userActivityHash cipher.SHA256, v interface{}, isCreation ...bool) error {
	uaRef, e := ap.Users.RefByHash(userActivityHash)
	if e != nil {
		return getArrayHashErr(e, userActivityHash, "Users")
	}
	ua, e := GetUserActivity(uaRef)
	if e != nil {
		return e
	}
	switch v.(type) {
	case *Content:
		if len(isCreation) > 0 && !isCreation[0] {
			if e := ua.ContentDeletions.Append(v); e != nil {
				return appendErr(e, v, "ActivityPage.ContentDeletions")
			}
		} else {
			if e := ua.ContentCreations.Append(v); e != nil {
				return appendErr(e, v, "ActivityPage.ContentCreations")
			}
		}
	case *Vote:
		if e := ua.VoteActions.Append(v); e != nil {
			return appendErr(e, v, "ActivityPage.VoteActions")
		}
	default:
		return boo.Newf(boo.NotAllowed,
			"invalid type '%T' provided", v)
	}

	// Save.
	if e := uaRef.SetValue(ua); e != nil {
		return boo.Newf(boo.NotAllowed,
			"failed to save")
	}
	return nil
}

/*
	<<< USER ACTIVITY >>>
*/

type UserActivity struct {
	PubKey           cipher.PubKey
	VoteActions      skyobject.Refs `skyobject:"schema=bbs.Vote"`
	ContentCreations skyobject.Refs `skyobject:"schema=bbs.Content"`
	ContentDeletions skyobject.Refs `skyobject:"schema=bbs.Content"`
}

func GetUserActivity(ref *skyobject.Ref) (*UserActivity, error) {
	uaValue, e := ref.Value()
	if e != nil {
		return nil, valErr(e, ref)
	}
	ua, ok := uaValue.(*UserActivity)
	if !ok {
		return nil, extErr(ref)
	}
	return ua, nil
}

type Board struct {
	Name     string
	Desc     string
	SubAddrs []string
	Created  int64
}

type Vote struct {
	OfUser   cipher.PubKey
	OfThread cipher.SHA256
	OfPost   cipher.SHA256

	Mode int8
	Tag  []byte

	Created int64         `verify:"time"`
	Creator cipher.PubKey `verify:"upk"`
	Sig     cipher.Sig    `verify:"sig"`
}

func ToVote(v interface{}) *Vote {
	if vote, ok := v.(*Vote); ok {
		return vote
	}
	return nil
}

func (v Vote) Verify() error { return tag.Verify(&v) }

/*
	<<< USER >>>
*/

type User struct {
	Alias  string        `json:"alias" trans:"alias"`
	PubKey cipher.PubKey `json:"-" trans:"upk"`
	SecKey cipher.SecKey `json:"-" trans:"usk"`
}

type UserView struct {
	User
	PubKey string `json:"public_key"`
	SecKey string `json:"secret_key,omitempty"`
}

/*
	<<< CONNECTION >>>
*/

type Connection struct {
	Address string `json:"address"`
	State   string `json:"state"`
}

func dynamicLock(mux *sync.Mutex) func() {
	if mux != nil {
		mux.Lock()
		return mux.Unlock
	}
	return func() {}
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func getRootChildErr(e error, i int) error {
	return boo.WrapTypef(e, boo.InvalidRead,
		"failed to get root child '%s' of index %d",
		indexString[i], i)
}

func extRootChildErr(i int) error {
	return boo.Newf(boo.InvalidRead,
		"failed to extract root child '%s' of index %d",
		indexString[i], i)
}

func saveRootChildErr(e error, i int) error {
	return boo.WrapTypef(e, boo.NotAllowed,
		"failed to save root child '%s' of index %d",
		indexString[i], i)
}

func valErr(e error, ref *skyobject.Ref) error {
	return boo.WrapTypef(e, boo.InvalidRead,
		"failed to obtain value from object of ref '%s'",
		ref.String())
}

func extErr(ref *skyobject.Ref) error {
	return boo.Newf(boo.InvalidRead,
		"failed to extract object from ref '%s'",
		ref.String())
}

func getArrayHashErr(e error, hash cipher.SHA256, what string) error {
	return boo.WrapTypef(e, boo.NotFound,
		"hash '%s' not found in '%s' array",
		hash.Hex(), what)
}

func appendErr(e error, v interface{}, what string) error {
	return boo.WrapTypef(e, boo.NotAllowed,
		"failed to append object '%v' to '%s' array",
		v, what)
}