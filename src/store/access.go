package store

import (
	"context"
	"github.com/skycoin/bbs/src/store/cxo"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/session"
	"github.com/skycoin/bbs/src/store/state/views"
	"github.com/skycoin/bbs/src/store/state/views/content_view"
	"time"
	"github.com/skycoin/bbs/src/rpc"
)

type Access struct {
	CXO     *cxo.Manager
	Session *session.Manager
}

/*
	<<< SESSION >>>
*/

func (a *Access) GetUsers(ctx context.Context) (*UsersOutput, error) {
	aliases, e := a.Session.GetUsers()
	if e != nil {
		return nil, e
	}
	return getUsers(ctx, aliases), nil
}

func (a *Access) NewUser(ctx context.Context, in *object.NewUserIO) (*UsersOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	if e := a.Session.NewUser(in); e != nil {
		return nil, e
	}
	return a.GetUsers(ctx)
}

func (a *Access) DeleteUser(ctx context.Context, alias string) (*UsersOutput, error) {
	if e := a.Session.DeleteUser(alias); e != nil {
		return nil, e
	}
	return a.GetUsers(ctx)
}

func (a *Access) GetSession(ctx context.Context) (*SessionOutput, error) {
	f, e := a.Session.GetCurrentFile()
	if e != nil && e != session.ErrNotLoggedIn {
		return nil, e
	}
	return getSession(ctx, f), nil
}

func (a *Access) Login(ctx context.Context, in *object.LoginIO) (*SessionOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	f, e := a.Session.Login(in)
	if e != nil {
		return nil, e
	}
	return getSession(ctx, f), nil
}

func (a *Access) Logout(ctx context.Context) (*SessionOutput, error) {
	if e := a.Session.Logout(); e != nil {
		return nil, e
	}
	return getSession(ctx, nil), nil
}

/*
	<<< CONNECTIONS >>>
*/

func (a *Access) GetConnections(ctx context.Context) (*ConnectionsOutput, error) {
	return getConnections(ctx, a.CXO.GetConnections()), nil
}

func (a *Access) NewConnection(ctx context.Context, in *object.ConnectionIO) (*ConnectionsOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	if e := a.CXO.Connect(in.Address); e != nil {
		return nil, e
	}
	time.Sleep(time.Second)
	return a.GetConnections(ctx)
}

func (a *Access) DeleteConnection(ctx context.Context, in *object.ConnectionIO) (*ConnectionsOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	if e := a.CXO.Disconnect(in.Address); e != nil {
		return nil, e
	}
	return a.GetConnections(ctx)
}

/*
	<<< SUBSCRIPTIONS >>>
*/

func (a *Access) GetSubscriptions(ctx context.Context) (*SubscriptionsOutput, error) {
	return getSubscriptions(ctx, a.CXO.GetSubscriptions()), nil
}

func (a *Access) NewSubscription(ctx context.Context, in *object.BoardIO) (*SubscriptionsOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	if e := a.CXO.SubscribeRemote(in.PubKey); e != nil {
		return nil, e
	}
	return a.GetSubscriptions(ctx)
}

func (a *Access) DeleteSubscription(ctx context.Context, in *object.BoardIO) (*SubscriptionsOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	if e := a.CXO.UnsubscribeRemote(in.PubKey); e != nil {
		return nil, e
	}
	return a.GetSubscriptions(ctx)
}

/*
	<<< CONTENT >>>
*/

func (a *Access) GetBoards(ctx context.Context) (*BoardsOutput, error) {
	m, r, e := a.CXO.GetBoards()
	if e != nil {
		return nil, e
	}
	return getBoardsOutput(ctx, m, r), nil
}

func (a *Access) NewBoard(ctx context.Context, in *object.NewBoardIO) (*BoardsOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	if e := a.CXO.NewBoard(in); e != nil {
		return nil, e
	}
	return a.GetBoards(ctx)
}

func (a *Access) DeleteBoard(ctx context.Context, in *object.BoardIO) (*BoardsOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	if e := a.CXO.UnsubscribeMaster(in.PubKey); e != nil {
		return nil, e
	}
	return a.GetBoards(ctx)
}

func (a *Access) GetBoardPage(ctx context.Context, in *object.BoardIO) (interface{}, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.PubKey)
	if e != nil {
		return nil, e
	}
	return bi.Get(views.Content, content_view.BoardPage)
}

func (a *Access) NewThread(ctx context.Context, in *object.NewThreadIO) (interface{}, error) {
	uf, e := a.Session.GetCurrentFile()
	if e != nil {
		return nil, e
	}
	if e := in.Process(uf.User.PubKey, uf.User.SecKey); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.BoardPubKey)
	if e != nil {
		return nil, e
	}
	var goal uint64
	if bi.IsMaster() {
		if goal, e = bi.NewThread(in.Thread); e != nil {
			return nil, e
		}
	} else {
		sa, e := bi.Get(views.Content, content_view.SubAddresses)
		if e != nil {
			return nil, e
		}
		goal, e = rpc.Send(ctx, sa, rpc.NewThread(in.Thread))
		if e != nil {
			return nil, e
		}
	}
	if e := bi.WaitSeq(ctx, goal); e != nil {
		return nil, e
	}
	return bi.Get(views.Content, content_view.BoardPage)
}

func (a *Access) GetThreadPage(ctx context.Context, in *object.ThreadIO) (interface{}, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.BoardPubKey)
	if e != nil {
		return nil, e
	}
	return bi.Get(views.Content, content_view.ThreadPage, in.ThreadRef)
}

func (a *Access) NewPost(ctx context.Context, in *object.NewPostIO) (interface{}, error) {
	uf, e := a.Session.GetCurrentFile()
	if e != nil {
		return nil, e
	}
	if e := in.Process(uf.User.PubKey, uf.User.SecKey); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.BoardPubKey)
	if e != nil {
		return nil, e
	}
	var goal uint64
	if bi.IsMaster() {
		if goal, e = bi.NewPost(in.Post); e != nil {
			return nil, e
		}
	} else {
		sa, e := bi.Get(views.Content, content_view.SubAddresses)
		if e != nil {
			return nil, e
		}
		goal, e = rpc.Send(ctx, sa, rpc.NewPost(in.Post))
		if e != nil {
			return nil, e
		}
	}
	if e := bi.WaitSeq(ctx, goal); e != nil {
		return nil, e
	}
	return bi.Get(views.Content, content_view.ThreadPage, in.ThreadRef)
}

/*
	<<< VOTES >>>
*/

func (a *Access) VoteUser(ctx context.Context, in *object.UserVoteIO) (interface{}, error) {
	uf, e := a.Session.GetCurrentFile()
	if e != nil {
		return nil, e
	}
	if e := in.Process(uf.User.PubKey, uf.User.SecKey); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.BoardPubKey)
	if e != nil {
		return nil, e
	}
	var goal uint64
	if bi.IsMaster() {
		if goal, e = bi.NewVote(in.Vote); e != nil {
			return nil, e
		}
	} else {
		sa, e := bi.Get(views.Content, content_view.SubAddresses)
		if e != nil {
			return nil, e
		}
		goal, e = rpc.Send(ctx, sa, rpc.NewVote(in.Vote))
		if e != nil {
			return nil, e
		}
	}
	if e := bi.WaitSeq(ctx, goal); e != nil {
		return nil, e
	}
	// TODO: Complete.
	return nil, nil
}

func (a *Access) VoteThread(ctx context.Context, in *object.ThreadVoteIO) (interface{}, error) {
	uf, e := a.Session.GetCurrentFile()
	if e != nil {
		return nil, e
	}
	if e := in.Process(uf.User.PubKey, uf.User.SecKey); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.BoardPubKey)
	if e != nil {
		return nil, e
	}
	var goal uint64
	if bi.IsMaster() {
		if goal, e = bi.NewVote(in.Vote); e != nil {
			return nil, e
		}
	} else {
		sa, e := bi.Get(views.Content, content_view.SubAddresses)
		if e != nil {
			return nil, e
		}
		goal, e = rpc.Send(ctx, sa, rpc.NewVote(in.Vote))
		if e != nil {
			return nil, e
		}
	}
	if e := bi.WaitSeq(ctx, goal); e != nil {
		return nil, e
	}
	// TODO: Complete.
	return nil, nil
}

func (a *Access) VotePost(ctx context.Context, in *object.PostVoteIO) (interface{}, error) {
	uf, e := a.Session.GetCurrentFile()
	if e != nil {
		return nil, e
	}
	if e := in.Process(uf.User.PubKey, uf.User.SecKey); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.BoardPubKey)
	if e != nil {
		return nil, e
	}
	var goal uint64
	if bi.IsMaster() {
		if goal, e = bi.NewVote(in.Vote); e != nil {
			return nil, e
		}
	} else {
		sa, e := bi.Get(views.Content, content_view.SubAddresses)
		if e != nil {
			return nil, e
		}
		goal, e = rpc.Send(ctx, sa, rpc.NewVote(in.Vote))
		if e != nil {
			return nil, e
		}
	}
	if e := bi.WaitSeq(ctx, goal); e != nil {
		return nil, e
	}
	// TODO: Complete.
	return nil, nil
}
