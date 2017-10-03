package r0

import (
	"fmt"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object/transfer"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"sync"
)

const (
	RootPageName    = "bbs.r0.RootPage"
	BoardPageName   = "bbs.r0.BoardPage"
	ThreadPageName  = "bbs.r0.ThreadPage"
	DiffPageName    = "bbs.r0.DiffPage"
	UsersPageName   = "bbs.r0.UsersPage"
	UserProfileName = "bbs.r0.UserProfile"
	ContentName     = "bbs.r0.Content"
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

func GetPages(p *skyobject.Pack, get ...bool) (out *Pages, e error) {
	out = &Pages{PK: p.Root().Pub}

	if len(get) > IndexRootPage && get[IndexRootPage] {
		if out.RootPage, e = GetRootPage(p); e != nil {
			return
		}
	}
	if len(get) > IndexBoardPage && get[IndexBoardPage] {
		if out.BoardPage, e = GetBoardPage(p); e != nil {
			return
		}
	}
	if len(get) > IndexDiffPage && get[IndexDiffPage] {
		if out.DiffPage, e = GetDiffPage(p); e != nil {
			return
		}
	}
	if len(get) > IndexUsersPage && get[IndexUsersPage] {
		if out.UsersPage, e = GetUsersPage(p); e != nil {
			return
		}
	}
	return
}

func (p *Pages) Save(pack *skyobject.Pack) error {
	if p.BoardPage != nil {
		if e := p.BoardPage.Save(pack); e != nil {
			return e
		}
	}
	if p.DiffPage != nil {
		if e := p.DiffPage.Save(pack); e != nil {
			return e
		}
	}
	if p.UsersPage != nil {
		if e := p.UsersPage.Save(pack); e != nil {
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
	Sum []byte // Summary of the root.
}

func GetRootPage(p *skyobject.Pack) (*RootPage, error) {
	rpVal, e := p.RefByIndex(IndexRootPage)
	if e != nil {
		return nil, getRootChildErr(e, IndexRootPage)
	}

	rp, ok := rpVal.(*RootPage)
	if !ok {
		return nil, extRootChildErr(IndexRootPage)
	}
	return rp, nil
}

func (rp *RootPage) ToRep() (*transfer.RootPageRep, error) {
	out := &transfer.RootPageRep{
		Type:     rp.Typ,
		Revision: rp.Rev,
		Deleted:  rp.Del,
		Summary:  transfer.RootPageSummary{}, // TODO: Implement summary.
	}
	return out, nil
}

func (rp *RootPage) FromRep(rpRep *transfer.RootPageRep) error {
	*rp = RootPage{
		Typ: rpRep.Type,
		Rev: rpRep.Revision,
		Del: rpRep.Deleted,
		Sum: nil, // TODO: Implement summary.
	}
	return nil
}

/*
	<<< BOARD PAGE >>>
*/

type BoardPage struct {
	Board   skyobject.Ref  `skyobject:"schema=bbs.r0.Content"`
	Threads skyobject.Refs `skyobject:"schema=bbs.r0.ThreadPage"`
}

func GetBoardPage(p *skyobject.Pack) (*BoardPage, error) {
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

func (bp *BoardPage) Save(p *skyobject.Pack) error {
	if e := p.SetRefByIndex(IndexBoardPage, bp); e != nil {
		return saveRootChildErr(e, IndexBoardPage)
	}
	return nil
}

func (bp *BoardPage) GetBoard() (*Board, error) {
	c, e := GetContentFromRef(&bp.Board)
	if e != nil {
		return nil, e
	}
	return c.ToBoard(), nil
}

func (bp *BoardPage) GetThreadCount() int {
	l, _ := bp.Threads.Len()
	return l
}

func (bp *BoardPage) RangeThreadPages(action func(i int, tp *ThreadPage) error) error {
	return bp.Threads.Ascend(func(i int, tpElem *skyobject.RefsElem) error {
		tp, e := GetThreadPage(tpElem)
		if e != nil {
			return e
		}
		return action(i, tp)
	})
}

func (bp *BoardPage) GetThreadPage(tpHash cipher.SHA256) (*skyobject.RefsElem, *ThreadPage, error) {
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

func (bp *BoardPage) AddThread(tRef skyobject.Ref) error {
	e := bp.Threads.Append(ThreadPage{Thread: tRef})
	if e != nil {
		return boo.Wrap(e, "failed to append thread to 'BoardPage.Threads'")
	}
	return nil
}

func (bp *BoardPage) ExportBoard() (transfer.Board, error) {
	return nil, nil
}

func (bp *BoardPage) ExportThreadPages() ([]transfer.ThreadPage, error) {
	l, e := bp.Threads.Len()
	if e != nil {
		return nil, e
	}
	out := make([]transfer.ThreadPage, l)
	e = bp.RangeThreadPages(func(i int, tp *ThreadPage) error {
		out[i] = tp
		return nil
	})
	return out, e
}

func (bp *BoardPage) ImportBoard(p *skyobject.Pack, b transfer.Board) error {
	return bp.Board.SetValue(b)
}

func (bp *BoardPage) ImportThreadPages(p *skyobject.Pack, tps []transfer.ThreadPage) error {
	bp.Threads.Clear()
	for _, tp := range tps {
		if e := bp.Threads.Append(tp); e != nil {
			return e
		}
	}
	return nil
}

/*
	<<< THREAD PAGE >>>
*/

type ThreadPage struct {
	Thread skyobject.Ref  `skyobject:"schema=bbs.r0.Content"`
	Posts  skyobject.Refs `skyobject:"schema=bbs.r0.Content"`
}

func GetThreadPage(tpElem *skyobject.RefsElem) (*ThreadPage, error) {
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

func (tp *ThreadPage) GetThread() (*Thread, error) {
	c, e := GetContentFromRef(&tp.Thread)
	if e != nil {
		return nil, e
	}
	return c.ToThread(), nil
}

func (tp *ThreadPage) GetPostCount() int {
	l, _ := tp.Posts.Len()
	return l
}

func (tp *ThreadPage) RangePosts(action func(i int, post *Post) error) error {
	return tp.Posts.Ascend(func(i int, pElem *skyobject.RefsElem) error {
		post, e := GetPost(pElem)
		if e != nil {
			return e
		}
		return action(i, post)
	})
}

func (tp *ThreadPage) AddPost(cxoPostHash cipher.SHA256, post *Post) error {
	if elem, _ := tp.Posts.RefByHash(cxoPostHash); elem != nil {
		return boo.Newf(boo.AlreadyExists,
			"post of hash '%s' already exists in 'ThreadPage.Posts'", cxoPostHash.Hex())
	}
	if e := tp.Posts.Append(post.Content); e != nil {
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

func (tp *ThreadPage) ExportThread() (transfer.Thread, error) {
	return nil, nil
}

func (tp *ThreadPage) ExportPosts() ([]transfer.Post, error) {
	return nil, nil
}

func (tp *ThreadPage) ImportThread(p *skyobject.Pack, t transfer.Thread) error {
	return tp.Thread.SetValue(t)
}

func (tp *ThreadPage) ImportPosts(p *skyobject.Pack, ps []transfer.Post) error {
	tp.Posts.Clear()
	for _, p := range ps {
		if e := tp.Posts.Append(p); e != nil {
			return e
		}
	}
	return nil
}

/*
	<<< DIFF PAGE >>>
*/

type DiffPage struct {
	Submissions skyobject.Refs `skyobject:"schema=bbs.r0.Content"`
}

func GetDiffPage(p *skyobject.Pack) (*DiffPage, error) {
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

func (dp *DiffPage) Save(p *skyobject.Pack) error {
	if e := p.SetRefByIndex(IndexDiffPage, dp); e != nil {
		return saveRootChildErr(e, IndexDiffPage)
	}
	return nil
}

func (dp *DiffPage) Add(c *Content) error {
	if e := dp.Submissions.Append(c); e != nil {
		return boo.Newf(boo.Internal,
			"failed to append %v to 'DiffPage.Submissions'", c)
	}
	return nil
}

func (dp *DiffPage) GetOfIndex(i int) (*Content, error) {
	cElem, e := dp.Submissions.RefByIndex(i)
	if e != nil {
		return nil, refByIndexErr(e, i, "DiffPage.Submissions")
	}
	return GetContentFromElem(cElem)
}

type Changes struct {
	NeedReset bool
	Total     int
	New       []*Content
}

func (dp *DiffPage) GetChanges(oldC *Changes) (*Changes, error) {
	newC := new(Changes)

	// Get counts.
	newC.Total, _ = dp.Submissions.Len()

	// Return if no old changes.
	if oldC == nil {
		return newC, nil
	}

	// Check if reset needed.
	if oldC.Total > newC.Total {
		newC.NeedReset = true
		return newC, nil
	}

	// Get content.
	if oldC.Total < newC.Total {
		newC.New = make([]*Content, newC.Total-oldC.Total)
		for i := oldC.Total; i < newC.Total; i++ {
			var e error
			newC.New[i-oldC.Total], e = dp.GetOfIndex(i)
			if e != nil {
				return nil, e
			}
		}
	}
	fmt.Printf("\t- (CHANGES) content: %d(+%d)\n",
		newC.Total, len(newC.New))

	return newC, nil
}

/*
	<<< USERS PAGE >>>
*/

type UsersPage struct {
	Users skyobject.Refs `skyobject:"schema=bbs.r0.UserProfile"`
}

func GetUsersPage(p *skyobject.Pack) (*UsersPage, error) {
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

func (up *UsersPage) Save(p *skyobject.Pack) error {
	if e := p.SetRefByIndex(IndexUsersPage, up); e != nil {
		return saveRootChildErr(e, IndexUsersPage)
	}
	return nil
}

func (up *UsersPage) NewUserProfile(upk string) (cipher.SHA256, error) {
	ua := &UserProfile{PubKey: upk}
	if e := up.Users.Append(ua); e != nil {
		return cipher.SHA256{}, appendErr(e, ua, "UsersPage.Users")
	}
	return cipher.SumSHA256(encoder.Serialize(ua)), nil
}

func (up *UsersPage) AddUserSubmission(uapHash cipher.SHA256, c *Content) (cipher.SHA256, error) {
	uapElem, e := up.Users.RefByHash(uapHash)
	if e != nil {
		return cipher.SHA256{}, refByHashErr(e, uapHash, "Users")
	}
	uap, e := GetUserActivityPage(uapElem)
	if e != nil {
		return cipher.SHA256{}, e
	}
	if e := uap.Submissions.Append(c); e != nil {
		return cipher.SHA256{}, appendErr(e, c, "UsersPage.Submissions")
	}

	// Save.
	if e := uapElem.SetValue(uap); e != nil {
		return cipher.SHA256{}, boo.Newf(boo.NotAllowed,
			"failed to save")
	}
	return uapElem.Hash, nil
}

func (up *UsersPage) RangeUserProfiles(action func(i int, uap *UserProfile) error) error {
	return up.Users.Ascend(func(i int, uapElem *skyobject.RefsElem) error {
		uap, e := GetUserActivityPage(uapElem)
		if e != nil {
			return e
		}
		return action(i, uap)
	})
}

/*
	<<< USER PROFILE >>>
*/

type UserProfile struct {
	R           cipher.SHA256 `enc:"-"`
	PubKey      string
	Submissions skyobject.Refs `skyobject:"schema=bbs.r0.Content"`
}

func GetUserActivityPage(uapElem *skyobject.RefsElem) (*UserProfile, error) {
	uapVal, e := uapElem.Value()
	if e != nil {
		return nil, elemValueErr(e, uapElem)
	}
	uap, ok := uapVal.(*UserProfile)
	if !ok {
		return nil, elemExtErr(uapElem)
	}
	uap.R = uapElem.Hash
	return uap, nil
}

func (uap *UserProfile) RangeSubmissions(action func(i int, c *Content) error) error {
	return uap.Submissions.Ascend(func(i int, vElem *skyobject.RefsElem) error {
		vote, e := GetContentFromElem(vElem)
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
	<<< GENERATOR (IMPORT/EXPORT) >>>
*/

type Generator struct {
	p *skyobject.Pack
}

func NewGenerator(p *skyobject.Pack) *Generator {
	return &Generator{p: p}
}

func (g *Generator) Pack() *skyobject.Pack {
	return g.p
}

func (g *Generator) NewRootPage() transfer.RootPage {
	return new(RootPage)
}

func (g *Generator) NewBoardPage() transfer.BoardPage {
	return new(BoardPage)
}

func (g *Generator) NewThreadPage() transfer.ThreadPage {
	return new(ThreadPage)
}

func (g *Generator) NewBoard() transfer.Board {
	return nil // new(Board)
}

func (g *Generator) NewThread() transfer.Thread {
	return nil // new(Thread)
}

func (g *Generator) NewPost() transfer.Post {
	return nil // new(Post)
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
