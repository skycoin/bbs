package r0

import (
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

type ExpRoot struct {
	RootPage  RootPage
	BoardPage *ExpBoardPage
	UsersPage *ExpUsersPage
}

func (r *ExpRoot) Fill(rp *RootPage, bp *BoardPage, up *UsersPage) error {
	var e error

	r.RootPage = *rp
	if r.BoardPage, e = ExportBoardPage(bp); e != nil {
		return e
	}
	if r.UsersPage, e = ExportUsersPage(up); e != nil {
		return e
	}
	return nil
}

func (r *ExpRoot) Dump(p *skyobject.Pack) error {

	if bp, e := ImportBoardPage(p, r.BoardPage); e != nil {
		return e
	} else if e := bp.Save(p); e != nil {
		return e
	} else if e := p.SetRefByIndex(IndexBoardPage, bp); e != nil {
		return e
	}

	if up, e := ImportUsersPage(p, r.UsersPage); e != nil {
		return e
	} else if e := up.Save(p); e != nil {
		return e
	} else if e := p.SetRefByIndex(IndexUsersPage, up); e != nil {
		return e
	}

	return nil
}

type ExpBoardPage struct {
	Board   Board
	Threads []*ExpThreadPage
}

func ExportBoardPage(bp *BoardPage) (*ExpBoardPage, error) {
	out := new(ExpBoardPage)

	board, e := bp.GetBoard()
	if e != nil {
		return nil, e
	}
	out.Board = *board

	out.Threads = make([]*ExpThreadPage, bp.GetThreadCount())
	if e := bp.RangeThreadPages(func(i int, tp *ThreadPage) error {
		out.Threads[i], e = ExportThreadPage(tp)
		return e
	}); e != nil {
		return nil, e
	}

	return out, nil
}

func ImportBoardPage(p *skyobject.Pack, in *ExpBoardPage) (*BoardPage, error) {
	bp := new(BoardPage)
	p.Ref(bp)
	bp.Board = p.Ref(&in.Board)
	for _, expTP := range in.Threads {
		tp, e := ImportThreadPage(p, expTP)
		if e != nil {
			return nil, e
		}
		if e := bp.Threads.Append(tp); e != nil {
			return nil, e
		}
	}
	return bp, nil
}

type ExpThreadPage struct {
	Thread Thread
	Posts  []Post
}

func ExportThreadPage(tp *ThreadPage) (*ExpThreadPage, error) {
	var out = new(ExpThreadPage)
	thread, e := tp.GetThread()
	if e != nil {
		return out, e
	}
	out.Thread = *thread

	out.Posts = make([]Post, tp.GetPostCount())
	if e := tp.RangePosts(func(i int, post *Post) error {
		out.Posts[i] = *post
		return nil
	}); e != nil {
		return nil, e
	}

	return out, nil
}

func ImportThreadPage(p *skyobject.Pack, in *ExpThreadPage) (*ThreadPage, error) {
	tp := new(ThreadPage)
	p.Ref(tp)
	tp.Thread = p.Ref(in.Thread)
	for _, post := range in.Posts {
		if e := tp.Posts.Append(post); e != nil {
			return nil, e
		}
	}
	return tp, nil
}

type ExpUsersPage struct {
	Users []*ExpUserActivityPage
}

func ExportUsersPage(up *UsersPage) (*ExpUsersPage, error) {
	out := new(ExpUsersPage)

	upLen, e := up.Users.Len()
	if e != nil {
		return nil, e
	}

	out.Users = make([]*ExpUserActivityPage, upLen)
	if e := up.RangeUserActivityPages(func(i int, uap *UserActivityPage) error {
		out.Users[i], e = ExportUserActivityPage(uap)
		return e
	}); e != nil {
		return nil, e
	}

	return out, nil
}

func ImportUsersPage(p *skyobject.Pack, in *ExpUsersPage) (*UsersPage, error) {
	up := new(UsersPage)
	p.Ref(up)
	for _, expUAP := range in.Users {
		uap, e := ImportUserActivityPage(p, expUAP)
		if e != nil {
			return nil, e
		}
		if e := up.Users.Append(uap); e != nil {
			return nil, e
		}
	}
	return up, nil
}

type ExpUserActivityPage struct {
	PK          cipher.PubKey
	VoteActions []Vote
}

func ExportUserActivityPage(uap *UserActivityPage) (*ExpUserActivityPage, error) {
	out := new(ExpUserActivityPage)
	out.PK = uap.PubKey

	vaLen, e := uap.VoteActions.Len()
	if e != nil {
		return nil, e
	}
	out.VoteActions = make([]Vote, vaLen)
	if e := uap.RangeVoteActions(func(i int, vote *Vote) error {
		out.VoteActions[i] = *vote
		return nil
	}); e != nil {
		return nil, e
	}

	return out, nil
}

func ImportUserActivityPage(p *skyobject.Pack, in *ExpUserActivityPage) (*UserActivityPage, error) {
	uap := &UserActivityPage{PubKey: in.PK}
	p.Ref(uap)
	for _, vote := range in.VoteActions {
		if e := uap.VoteActions.Append(vote); e != nil {
			return nil, e
		}
	}
	return uap, nil
}
