package r0

import (
	"fmt"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"sync"
)

const (
	RootPageName         = "bbs.r0.RootPage"
	BoardPageName        = "bbs.r0.BoardPage"
	ThreadPageName       = "bbs.r0.ThreadPage"
	DiffPageName         = "bbs.r0.DiffPage"
	UsersPageName        = "bbs.r0.UsersPage"
	UserActivityPageName = "bbs.r0.UserActivityPage"
	BoardName            = "bbs.r0.Board"
	ThreadName           = "bbs.r0.Thread"
	PostName             = "bbs.r0.Post"
	VoteName             = "bbs.r0.Vote"
	UserName             = "bbs.r0.User"
)

const (
	IndexRootPage     = 0
	IndexBoardPage    = 1
	IndexDiffPage     = 2
	IndexUsersPage    = 3
	RootChildrenCount = 4
)

var indexString = [...]string{
	IndexRootPage:  "RootPage",
	IndexBoardPage: "BoardPage",
	IndexDiffPage:  "DiffPage",
	IndexUsersPage: "UsersPage",
}

/*
	<<< ROOT CHILDREN >>>
*/

type Pages struct {
	PK        cipher.PubKey
	RootPage  *RootPage
	BoardPage *BoardPage
	DiffPage  *DiffPage
	UsersPage *UsersPage
}

func GetPages(p *skyobject.Pack, mux *sync.Mutex, get ...bool) (out *Pages, e error) {
	defer dynamicLock(mux)()
	out = &Pages{PK: p.Root().Pub}

	if len(get) > IndexRootPage && get[IndexRootPage] {
		// TODO: Implement.
	}

	if len(get) > IndexBoardPage && get[IndexBoardPage] {
		if out.BoardPage, e = GetBoardPage(p, nil); e != nil {
			return
		}
	}
	if len(get) > IndexDiffPage && get[IndexDiffPage] {
		if out.DiffPage, e = GetDiffPage(p, nil); e != nil {
			return
		}
	}
	if len(get) > IndexUsersPage && get[IndexUsersPage] {
		if out.UsersPage, e = GetUsersPage(p, nil); e != nil {
			return
		}
	}
	return
}

func (p *Pages) Save(pack *skyobject.Pack, mux *sync.Mutex) error {
	defer dynamicLock(mux)()
	if p.BoardPage != nil {
		if e := p.BoardPage.Save(pack, nil); e != nil {
			return e
		}
	}
	if p.DiffPage != nil {
		if e := p.DiffPage.Save(pack, nil); e != nil {
			return e
		}
	}
	if p.UsersPage != nil {
		if e := p.UsersPage.Save(pack, nil); e != nil {
			return e
		}
	}
	return nil
}

/*
	<<< ROOT PAGE >>>
*/

const (
	RootTypeBoard = "board"
)

// RootPage helps determine the type, version of the root, and whether the root has been deleted.
type RootPage struct {
	Typ string // Type of root.
	Rev uint64 // Revision of root type.
	Del bool   // Whether root is deleted.
}

/*
	<<< BOARD PAGE >>>
*/

type BoardPage struct {
	Board   skyobject.Ref  `skyobject:"schema=bbs.r0.Board"`
	Threads skyobject.Refs `skyobject:"schema=bbs.r0.ThreadPage"`
}

func GetBoardPage(p *skyobject.Pack, mux *sync.Mutex) (*BoardPage, error) {
	defer dynamicLock(mux)()
	bpVal, e := p.RefByIndex(IndexBoardPage)
	if e != nil {
		return nil, getRootChildErr(e, IndexBoardPage)
	}
	bp, ok := bpVal.(*BoardPage)
	if !ok {
		return nil, extRootChildErr(IndexBoardPage)
	}
	return bp, nil
}

func (bp *BoardPage) Save(p *skyobject.Pack, mux *sync.Mutex) error {
	defer dynamicLock(mux)()
	if e := p.SetRefByIndex(IndexBoardPage, bp); e != nil {
		return saveRootChildErr(e, IndexBoardPage)
	}
	return nil
}

func (bp *BoardPage) GetBoard(mux *sync.Mutex) (*Board, error) {
	bVal, e := bp.Board.Value()
	if e != nil {
		return nil, valueErr(e, &bp.Board)
	}
	b, ok := bVal.(*Board)
	if !ok {
		return nil, extErr(&bp.Board)
	}
	return b, nil
}

func (bp *BoardPage) GetThreadCount() int {
	l, _ := bp.Threads.Len()
	return l
}

func (bp *BoardPage) RangeThreadPages(action func(i int, tp *ThreadPage) error, mux *sync.Mutex) error {
	defer dynamicLock(mux)()
	return bp.Threads.Ascend(func(i int, tpElem *skyobject.RefsElem) error {
		tp, e := GetThreadPage(tpElem, nil)
		if e != nil {
			return e
		}
		return action(i, tp)
	})
}

func (bp *BoardPage) GetThreadPage(tpHash cipher.SHA256, mux *sync.Mutex) (*skyobject.RefsElem, *ThreadPage, error) {
	defer dynamicLock(mux)()
	tpElem, e := bp.Threads.RefByHash(tpHash)
	if e != nil {
		return nil, nil, refByHashErr(e, tpHash, "BoardPage.Threads")
	}
	tpVal, e := tpElem.Value()
	if e != nil {
		return nil, nil, elemValueErr(e, tpElem)
	}
	tp, ok := tpVal.(*ThreadPage)
	if !ok {
		return nil, nil, elemExtErr(tpElem)
	}
	return tpElem, tp, nil
}

func (bp *BoardPage) AddThread(tRef skyobject.Ref, mux *sync.Mutex) error {
	defer dynamicLock(mux)()
	e := bp.Threads.Append(ThreadPage{Thread: tRef})
	if e != nil {
		return boo.Wrap(e, "failed to append thread to 'BoardPage.Threads'")
	}
	return nil
}

/*
	<<< THREAD PAGE >>>
*/

type ThreadPage struct {
	Thread skyobject.Ref  `skyobject:"schema=bbs.r0.Thread"`
	Posts  skyobject.Refs `skyobject:"schema=bbs.r0.Post"`
}

func GetThreadPage(tpElem *skyobject.RefsElem, mux *sync.Mutex) (*ThreadPage, error) {
	defer dynamicLock(mux)()
	tpVal, e := tpElem.Value()
	if e != nil {
		return nil, elemValueErr(e, tpElem)
	}
	tp, ok := tpVal.(*ThreadPage)
	if !ok {
		return nil, elemExtErr(tpElem)
	}
	return tp, nil
}

func (tp *ThreadPage) GetThread(mux *sync.Mutex) (*Thread, error) {
	defer dynamicLock(mux)()
	tVal, e := tp.Thread.Value()
	if e != nil {
		return nil, valueErr(e, &tp.Thread)
	}
	t, ok := tVal.(*Thread)
	if !ok {
		return nil, extErr(&tp.Thread)
	}
	t.R = tp.Thread.Hash
	return t, nil
}

func (tp *ThreadPage) GetPostCount() int {
	l, _ := tp.Posts.Len()
	return l
}

func (tp *ThreadPage) RangePosts(action func(i int, post *Post) error, mux *sync.Mutex) error {
	defer dynamicLock(mux)()
	return tp.Posts.Ascend(func(i int, pElem *skyobject.RefsElem) error {
		post, e := GetPost(pElem, nil)
		if e != nil {
			return e
		}
		return action(i, post)
	})
}

func (tp *ThreadPage) AddPost(postHash cipher.SHA256, post *Post, mux *sync.Mutex) error {
	//if elem, e := tp.Posts.RefByHash(postHash); e == nil {
	//	return boo.Newf(boo.AlreadyExists,
	//		"post of hash '%s' already exists in 'ThreadPage.Posts'", postHash.Hex())
	//
	//	fmt.Println(elem.String())
	//}
	// TODO: Redo
	if e := tp.Posts.Append(post); e != nil {
		return boo.WrapTypef(e, boo.Internal,
			"failed to append %v to 'ThreadPage.Posts'", post)
	}
	return nil
}

func (tp *ThreadPage) Save(tpElem *skyobject.RefsElem) error {
	if e := tpElem.SetValue(tp); e != nil {
		return boo.WrapType(e, boo.Internal, "failed to save 'ThreadPage'")
	}
	return nil
}

/*
	<<< DIFF PAGE >>>
*/

type DiffPage struct {
	Threads skyobject.Refs `skyobject:"schema=bbs.r0.Thread"`
	Posts   skyobject.Refs `skyobject:"schema=bbs.r0.Post"`
	Votes   skyobject.Refs `skyobject:"schema=bbs.r0.Vote"`
}

func GetDiffPage(p *skyobject.Pack, mux *sync.Mutex) (*DiffPage, error) {
	defer dynamicLock(mux)()
	dpVal, e := p.RefByIndex(IndexDiffPage)
	if e != nil {
		return nil, getRootChildErr(e, IndexDiffPage)
	}
	dp, ok := dpVal.(*DiffPage)
	if !ok {
		return nil, extRootChildErr(IndexDiffPage)
	}
	return dp, nil
}

func (dp *DiffPage) Save(p *skyobject.Pack, mux *sync.Mutex) error {
	defer dynamicLock(mux)()
	if e := p.SetRefByIndex(IndexDiffPage, dp); e != nil {
		return saveRootChildErr(e, IndexDiffPage)
	}
	return nil
}

func (dp *DiffPage) Add(v interface{}) error {
	switch v.(type) {
	case *Thread:
		if e := dp.Threads.Append(v); e != nil {
			return boo.Newf(boo.Internal,
				"failed to append %v to 'DiffPage.Threads'", v)
		}
	case *Post:
		if e := dp.Posts.Append(v); e != nil {
			return boo.Newf(boo.Internal,
				"failed to append %v to 'DiffPage.Posts'", v)
		}
	case *Vote:
		if e := dp.Votes.Append(v); e != nil {
			return boo.Newf(boo.Internal,
				"failed to append %v to 'DiffPage.Votes'", v)
		}
	default:
		return boo.Newf(boo.Internal,
			"failed to add object of type %T to 'DiffPage'", v)
	}
	return nil
}

func (dp *DiffPage) GetThreadOfIndex(i int, mux *sync.Mutex) (*Thread, error) {
	defer dynamicLock(mux)()
	tElem, e := dp.Threads.RefByIndex(i)
	if e != nil {
		return nil, refByIndexErr(e, i, "DiffPage.Threads")
	}
	tVal, e := tElem.Value()
	if e != nil {
		return nil, elemValueErr(e, tElem)
	}
	t, ok := tVal.(*Thread)
	if !ok {
		return nil, elemExtErr(tElem)
	}
	t.R = tElem.Hash
	return t, nil
}

func (dp *DiffPage) GetPostOfIndex(i int, mux *sync.Mutex) (*Post, error) {
	defer dynamicLock(mux)()
	pElem, e := dp.Posts.RefByIndex(i)
	if e != nil {
		return nil, refByIndexErr(e, i, "DiffPage.Posts")
	}
	pVal, e := pElem.Value()
	if e != nil {
		return nil, elemValueErr(e, pElem)
	}
	p, ok := pVal.(*Post)
	if !ok {
		return nil, elemExtErr(pElem)
	}
	p.R = pElem.Hash
	return p, nil
}

func (dp *DiffPage) GetVoteOfIndex(i int, mux *sync.Mutex) (*Vote, error) {
	defer dynamicLock(mux)()
	vElem, e := dp.Votes.RefByIndex(i)
	if e != nil {
		return nil, refByIndexErr(e, i, "DiffPage.Votes")
	}
	vVal, e := vElem.Value()
	if e != nil {
		return nil, elemValueErr(e, vElem)
	}
	v, ok := vVal.(*Vote)
	if !ok {
		return nil, elemExtErr(vElem)
	}
	return v, nil
}

type Changes struct {
	ThreadCount int
	PostCount   int
	VoteCount   int

	NewThreads []*Thread
	NewPosts   []*Post
	NewVotes   []*Vote
}

func (dp *DiffPage) GetChanges(oldC *Changes, mux *sync.Mutex) (*Changes, error) {
	defer dynamicLock(mux)()
	newC := new(Changes)

	// Get counts.
	newC.ThreadCount, _ = dp.Threads.Len()
	newC.PostCount, _ = dp.Posts.Len()
	newC.VoteCount, _ = dp.Votes.Len()

	// Return if no old changes.
	if oldC == nil {
		return newC, nil
	}

	// Get content.
	if oldC.ThreadCount < newC.ThreadCount {
		newC.NewThreads = make([]*Thread, newC.ThreadCount-oldC.ThreadCount)
		for i := oldC.ThreadCount; i < newC.ThreadCount; i++ {
			var e error
			newC.NewThreads[i-oldC.ThreadCount], e = dp.GetThreadOfIndex(i, nil)
			if e != nil {
				return nil, e
			}
		}
	}
	if oldC.PostCount < newC.PostCount {
		newC.NewPosts = make([]*Post, newC.PostCount-oldC.PostCount)
		for i := oldC.PostCount; i < newC.PostCount; i++ {
			var e error
			newC.NewPosts[i-oldC.PostCount], e = dp.GetPostOfIndex(i, nil)
			if e != nil {
				return nil, e
			}
		}
	}
	if oldC.VoteCount < newC.VoteCount {
		newC.NewVotes = make([]*Vote, newC.VoteCount-oldC.VoteCount)
		for i := oldC.VoteCount; i < newC.VoteCount; i++ {
			var e error
			newC.NewVotes[i-oldC.VoteCount], e = dp.GetVoteOfIndex(i, nil)
			if e != nil {
				return nil, e
			}
		}
	}
	fmt.Printf("CHANGES(total count): t(%d) p(%d) v(%d)\n",
		newC.ThreadCount, newC.PostCount, newC.VoteCount)
	fmt.Printf("CHANGES(added count): t(%d) p(%d) v(%d)\n\n",
		len(newC.NewThreads), len(newC.NewPosts), len(newC.NewVotes))
	return newC, nil
}

/*
	<<< USERS PAGE >>>
*/

type UsersPage struct {
	Users skyobject.Refs `skyobject:"schema=bbs.r0.UserActivityPage"`
}

func GetUsersPage(p *skyobject.Pack, mux *sync.Mutex) (*UsersPage, error) {
	defer dynamicLock(mux)()
	upVal, e := p.RefByIndex(IndexUsersPage)
	if e != nil {
		return nil, getRootChildErr(e, IndexUsersPage)
	}
	up, ok := upVal.(*UsersPage)
	if !ok {
		return nil, extRootChildErr(IndexUsersPage)
	}
	return up, nil
}

func (up *UsersPage) Save(p *skyobject.Pack, mux *sync.Mutex) error {
	defer dynamicLock(mux)()
	if e := p.SetRefByIndex(IndexUsersPage, up); e != nil {
		return saveRootChildErr(e, IndexUsersPage)
	}
	return nil
}

func (up *UsersPage) NewUserActivityPage(upk cipher.PubKey) (cipher.SHA256, error) {
	ua := &UserActivityPage{PubKey: upk}
	if e := up.Users.Append(ua); e != nil {
		return cipher.SHA256{}, appendErr(e, ua, "UsersPage.Users")
	}
	return cipher.SumSHA256(encoder.Serialize(ua)), nil
}

func (up *UsersPage) AddUserActivity(uapHash cipher.SHA256, v interface{}) error {
	uapElem, e := up.Users.RefByHash(uapHash)
	if e != nil {
		return refByHashErr(e, uapHash, "Users")
	}
	uap, e := GetUserActivityPage(uapElem, nil)
	if e != nil {
		return e
	}
	switch v.(type) {
	case *Vote:
		if e := uap.VoteActions.Append(v); e != nil {
			return appendErr(e, v, "UsersPage.VoteActions")
		}
	default:
		return boo.Newf(boo.NotAllowed,
			"invalid type '%T' provided", v)
	}

	// Save.
	if e := uapElem.SetValue(uap); e != nil {
		return boo.Newf(boo.NotAllowed,
			"failed to save")
	}
	return nil
}

func (up *UsersPage) RangeUserActivityPages(action func(i int, uap *UserActivityPage) error, mux *sync.Mutex) error {
	defer dynamicLock(mux)()
	return up.Users.Ascend(func(i int, uapElem *skyobject.RefsElem) error {
		uap, e := GetUserActivityPage(uapElem, nil)
		if e != nil {
			return e
		}
		return action(i, uap)
	})
}

/*
	<<< USER ACTIVITY PAGE >>>
*/

type UserActivityPage struct {
	R           cipher.SHA256 `enc:"-"`
	PubKey      cipher.PubKey
	VoteActions skyobject.Refs `skyobject:"schema=bbs.r0.Vote"`
}

func GetUserActivityPage(uapElem *skyobject.RefsElem, mux *sync.Mutex) (*UserActivityPage, error) {
	defer dynamicLock(mux)()
	uapVal, e := uapElem.Value()
	if e != nil {
		return nil, elemValueErr(e, uapElem)
	}
	uap, ok := uapVal.(*UserActivityPage)
	if !ok {
		return nil, elemExtErr(uapElem)
	}
	uap.R = uapElem.Hash
	return uap, nil
}

func (uap *UserActivityPage) RangeVoteActions(action func(i int, vote *Vote) error, mux *sync.Mutex) error {
	defer dynamicLock(mux)()
	return uap.VoteActions.Ascend(func(i int, vElem *skyobject.RefsElem) error {
		vote, e := GetVote(vElem, nil)
		if e != nil {
			return e
		}
		return action(i, vote)
	})
}

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
	PubKey string `json:"public_key,omitempty"`
	SecKey string `json:"secret_key,omitempty"`
}

/*
	<<< CONNECTION >>>
*/

type Connection struct {
	Address string `json:"address"`
	State   string `json:"state"`
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func dynamicLock(mux *sync.Mutex) func() {
	if mux != nil {
		mux.Lock()
		return mux.Unlock
	}
	return func() {}
}

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

func valueErr(e error, ref *skyobject.Ref) error {
	return boo.WrapTypef(e, boo.InvalidRead,
		"failed to obtain value from object of ref '%s'",
		ref.String())
}

func elemValueErr(e error, elem *skyobject.RefsElem) error {
	return boo.WrapType(e, boo.InvalidRead,
		"failed to obtain value from elem object of ref '%s'",
		elem.String())
}

func extErr(ref *skyobject.Ref) error {
	return boo.Newf(boo.InvalidRead,
		"failed to extract object from ref '%s'",
		ref.String())
}

func elemExtErr(elem *skyobject.RefsElem) error {
	return boo.Newf(boo.InvalidRead,
		"failed to extract object from elem '%s'",
		elem.String())
}

func refByHashErr(e error, hash cipher.SHA256, what string) error {
	return boo.WrapTypef(e, boo.NotFound,
		"failed to get hash '%s' from '%s' array",
		hash.Hex(), what)
}

func refByIndexErr(e error, i int, what string) error {
	return boo.WrapTypef(e, boo.NotFound,
		"failed to get '%s[%d]'", what, i)
}

func appendErr(e error, v interface{}, what string) error {
	return boo.WrapTypef(e, boo.NotAllowed,
		"failed to append object '%v' to '%s' array",
		v, what)
}
