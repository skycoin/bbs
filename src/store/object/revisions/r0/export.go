package r0

import (
	"github.com/skycoin/skycoin/src/cipher"
)

type ExpRoot struct {
	RootPage  RootPage
	BoardPage *ExpBoardPage
	UsersPage *ExpUsersPage
}

func (r *ExpRoot) Fill(rp *RootPage, bp *BoardPage, up *UsersPage) (error) {
	var e error

	//r.RootPage = *rp
	if r.BoardPage, e = ExportBoardPage(bp); e != nil {
		return e
	}
	if r.UsersPage, e = ExportUsersPage(up); e != nil {
		return e
	}
	return nil
}

type ExpBoardPage struct {
	Board Board
	Threads []*ExpThreadPage
}

func ExportBoardPage(bp *BoardPage) (*ExpBoardPage, error) {
	out := new(ExpBoardPage)

	board, e := bp.GetBoard(nil)
	if e != nil {
		return nil, e
	}
	out.Board = *board

	out.Threads = make([]*ExpThreadPage, bp.GetThreadCount())
	if e := bp.RangeThreadPages(func(i int, tp *ThreadPage) error {
		out.Threads[i], e = ExportThreadPage(tp)
		return e
	}, nil); e != nil {
		return nil, e
	}

	return out, nil
}

type ExpThreadPage struct {
	Thread Thread
	Posts []Post
}

func ExportThreadPage(tp *ThreadPage) (*ExpThreadPage, error) {
	var out = new(ExpThreadPage)

	thread, e := tp.GetThread(nil)
	if e != nil {
		return out, e
	}
	out.Thread = *thread

	out.Posts = make([]Post, tp.GetPostCount())
	if e := tp.RangePosts(func(i int, post *Post) error {
		out.Posts[i] = *post
		return nil
	}, nil); e != nil {
		return nil, e
	}

	return out, nil
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
	}, nil); e != nil {
		return nil, e
	}

	return out, nil
}

type ExpUserActivityPage struct {
	PK cipher.PubKey
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
	}, nil); e != nil {
		return nil, e
	}

	return out, nil
}