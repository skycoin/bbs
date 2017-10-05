package store

import (
	"context"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/cxo"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/object/revisions/r0"
	"github.com/skycoin/bbs/src/store/session"
	"github.com/skycoin/bbs/src/store/state/views"
	"github.com/skycoin/bbs/src/store/state/views/content_view"
	"github.com/skycoin/bbs/src/store/state/views/follow_view"
	"log"
	"time"
)

type Access struct {
	CXO     *cxo.Manager
	Session *session.Manager
}

func (a *Access) SubmitContent(ctx context.Context, in *object.SubmissionIO) (interface{}, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	switch in.Type {
	case r0.V5ThreadType:
		transport, e := r0.NewThreadTransport(in.Body, in.Sig, nil)
		if e != nil {
			return nil, e
		}
		bi, e := a.CXO.GetBoardInstance(transport.OfBoard)
		if e != nil {
			return nil, e
		}
		thread := new(r0.Thread).Fill(transport)
		var goal uint64
		if bi.IsMaster() {
			if goal, e = bi.NewThread(thread); e != nil {
				return nil, e
			}
		} else {
			goal, e = a.CXO.Relay().NewContent(ctx, bi.GetSubmissionKeys(), thread.Content)
			if e != nil {
				return nil, e
			}
		}
		if e := bi.WaitSeq(ctx, goal); e != nil {
			return nil, e
		}
		return bi.Get(views.Content, content_view.BoardPage, &content_view.BoardPageIn{
			Perspective: transport.Header.PK,
		})

	case r0.V5PostType:
		transport, e := r0.NewPostTransport(in.Body, in.Sig, nil)
		if e != nil {
			return nil, e
		}
		bi, e := a.CXO.GetBoardInstance(transport.OfBoard)
		if e != nil {
			return nil, e
		}
		post := new(r0.Post).Fill(transport)
		var goal uint64
		if bi.IsMaster() {
			if goal, e = bi.NewPost(post); e != nil {
				return nil, e
			}
		} else {
			goal, e = a.CXO.Relay().NewContent(ctx, bi.GetSubmissionKeys(), post.Content)
			if e != nil {
				return nil, e
			}
		}
		if e := bi.WaitSeq(ctx, goal); e != nil {
			return nil, e
		}
		return bi.Get(views.Content, content_view.ThreadPage, &content_view.ThreadPageIn{
			Perspective: post.GetBody().Creator,
			ThreadHash:  post.GetBody().OfThread,
		})

	case r0.V5ThreadVoteType:
		transport, e := r0.NewThreadVoteTransport(in.Body, in.Sig, nil)
		if e != nil {
			return nil, e
		}
		bi, e := a.CXO.GetBoardInstance(transport.OfBoard)
		if e != nil {
			return nil, e
		}
		threadVote := new(r0.ThreadVote).Fill(transport)
		var goal uint64
		if bi.IsMaster() {
			if goal, e = bi.NewThreadVote(threadVote); e != nil {
				return nil, e
			}
		} else {
			goal, e = a.CXO.Relay().NewContent(ctx, bi.GetSubmissionKeys(), threadVote.Content)
			if e != nil {
				return nil, e
			}
		}
		if e := bi.WaitSeq(ctx, goal); e != nil {
			return nil, e
		}
		return bi.Get(views.Content, content_view.ContentVotes, &content_view.ContentVotesIn{
			Perspective: threadVote.GetBody().Creator,
			ContentHash: threadVote.GetHeader().Hash,
		})

	case r0.V5PostVoteType:
		transport, e := r0.NewPostVoteTransport(in.Body, in.Sig, nil)
		if e != nil {
			return nil, e
		}
		bi, e := a.CXO.GetBoardInstance(transport.OfBoard)
		if e != nil {
			return nil, e
		}
		postVote := new(r0.PostVote).Fill(transport)
		var goal uint64
		if bi.IsMaster() {
			if goal, e = bi.NewPostVote(postVote); e != nil {
				return nil, e
			}
		} else {
			goal, e = a.CXO.Relay().NewContent(ctx, bi.GetSubmissionKeys(), postVote.Content)
			if e != nil {
				return nil, e
			}
		}
		if e := bi.WaitSeq(ctx, goal); e != nil {
			return nil, e
		}
		return bi.Get(views.Content, content_view.ContentVotes, &content_view.ContentVotesIn{
			Perspective: postVote.GetBody().Creator,
			ContentHash: postVote.GetHeader().Hash,
		})

	case r0.V5UserVoteType:
		transport, e := r0.NewUserVoteTransport(in.Body, in.Sig, nil)
		if e != nil {
			return nil, e
		}
		bi, e := a.CXO.GetBoardInstance(transport.OfBoard)
		if e != nil {
			return nil, e
		}
		userVote := new(r0.UserVote).Fill(transport)
		var goal uint64
		if bi.IsMaster() {
			if goal, e = bi.NewUserVote(userVote); e != nil {
				return nil, e
			}
		} else {
			goal, e = a.CXO.Relay().NewContent(ctx, bi.GetSubmissionKeys(), userVote.Content)
			if e != nil {
				return nil, e
			}
		}
		if e := bi.WaitSeq(ctx, goal); e != nil {
			return nil, e
		}
		out, e := bi.Get(views.Follow, follow_view.FollowPage, userVote.GetBody().Creator)
		if e != nil {
			return nil, e
		}
		return getFollowPageOutput(out), nil

	default:
		return nil, boo.Newf(boo.InvalidInput,
			"content submission of type '%s' is invalid", in.Type)
	}
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
	<<< CONTENT : ADMIN >>>
*/

func (a *Access) GetBoards(ctx context.Context) (*BoardsOutput, error) {
	m, r, e := a.CXO.GetBoards(ctx)
	if e != nil {
		return nil, e
	}
	return getBoardsOutput(ctx, m, r), nil
}

func (a *Access) NewBoard(ctx context.Context, in *object.NewBoardIO) (*BoardsOutput, error) {
	if e := in.Process(a.CXO.Relay().GetKeys()); e != nil {
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
		if e := a.CXO.UnsubscribeRemote(in.PubKey); e != nil {
			return nil, e
		}
	}
	return a.GetBoards(ctx)
}

func (a *Access) GetBoard(ctx context.Context, in *object.BoardIO) (*BoardOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.PubKey)
	if e != nil {
		return nil, e
	}
	board, e := bi.Get(views.Content, content_view.Board)
	if e != nil {
		return nil, e
	}
	return getBoardOutput(board), nil
}

func (a *Access) ExportBoard(ctx context.Context, in *object.ExportBoardIO) (*ExportBoardOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	path, data, e := a.CXO.ExportBoard(in.PubKey, in.Name)
	if e != nil {
		return nil, e
	}
	return getExportBoardOutput(path, data), nil
}

func (a *Access) ImportBoard(ctx context.Context, in *object.ExportBoardIO) (*ExportBoardOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	path, out, e := a.CXO.ImportBoard(in.PubKey, in.Name)
	if e != nil {
		return nil, e
	}
	return getExportBoardOutput(path, out), nil
}

func (a *Access) GetDiscoveredBoards(ctx context.Context) ([]string, error) {
	return a.CXO.GetDiscoveredBoards(), nil
}

/*
	<<< CONTENT >>>
*/

func (a *Access) GetBoardPage(ctx context.Context, in *object.BoardIO) (interface{}, error) {
	var perspective string
	if uf, _ := a.Session.GetCurrentFile(); uf != nil {
		perspective = uf.User.PubKey.Hex()
	}
	if e := in.Process(); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.PubKey)
	if e != nil {
		return nil, e
	}
	return bi.Get(views.Content, content_view.BoardPage, &content_view.BoardPageIn{
		Perspective: perspective,
	})
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
		log.Println("NewThread: subs:", bi.GetSubmissionKeys())
		goal, e = a.CXO.Relay().NewContent(ctx, bi.GetSubmissionKeys(), in.Thread.Content)
		if e != nil {
			return nil, e
		}
	}
	if e := bi.WaitSeq(ctx, goal); e != nil {
		return nil, e
	}
	return bi.Get(views.Content, content_view.BoardPage, &content_view.BoardPageIn{
		Perspective: uf.User.PubKey.Hex(),
	})
}

func (a *Access) GetThreadPage(ctx context.Context, in *object.ThreadIO) (interface{}, error) {
	var perspective string
	if uf, _ := a.Session.GetCurrentFile(); uf != nil {
		perspective = uf.User.PubKey.Hex()
	}
	if e := in.Process(); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.BoardPubKey)
	if e != nil {
		return nil, e
	}
	return bi.Get(views.Content, content_view.ThreadPage, &content_view.ThreadPageIn{
		Perspective: perspective,
		ThreadHash:  in.ThreadRefStr,
	})
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
		goal, e = a.CXO.Relay().NewContent(ctx, bi.GetSubmissionKeys(), in.Post.Content)
		if e != nil {
			return nil, e
		}
	}
	if e := bi.WaitSeq(ctx, goal); e != nil {
		return nil, e
	}
	return bi.Get(views.Content, content_view.ThreadPage, &content_view.ThreadPageIn{
		Perspective: uf.User.PubKey.Hex(),
		ThreadHash:  in.ThreadRefStr,
	})
}

/*
	<<< VOTES >>>
*/

func (a *Access) GetFollowPage(ctx context.Context, in *object.UserIO) (interface{}, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.BoardPubKey)
	if e != nil {
		return nil, e
	}
	out, e := bi.Get(views.Follow, follow_view.FollowPage, in.UserPubKey)
	if e != nil {
		return nil, e
	}
	return getFollowPageOutput(out), nil
}

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
		if goal, e = bi.NewUserVote(in.Vote); e != nil {
			return nil, e
		}
	} else {
		goal, e = a.CXO.Relay().NewContent(ctx, bi.GetSubmissionKeys(), in.Vote.Content)
		if e != nil {
			return nil, e
		}
	}
	if e := bi.WaitSeq(ctx, goal); e != nil {
		return nil, e
	}
	out, e := bi.Get(views.Follow, follow_view.FollowPage, uf.User.PubKey)
	if e != nil {
		return nil, e
	}
	return getFollowPageOutput(out), nil
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
		if goal, e = bi.NewThreadVote(in.Vote); e != nil {
			return nil, e
		}
	} else {
		goal, e = a.CXO.Relay().NewContent(ctx, bi.GetSubmissionKeys(), in.Vote.Content)
		if e != nil {
			return nil, e
		}
	}
	if e := bi.WaitSeq(ctx, goal); e != nil {
		return nil, e
	}
	return bi.Get(views.Content, content_view.ContentVotes, &content_view.ContentVotesIn{
		Perspective: uf.User.PubKey.Hex(),
		ContentHash: in.ThreadRefStr,
	})
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
		if goal, e = bi.NewPostVote(in.Vote); e != nil {
			return nil, e
		}
	} else {
		goal, e = a.CXO.Relay().NewContent(ctx, bi.GetSubmissionKeys(), in.Vote.Content)
		if e != nil {
			return nil, e
		}
	}
	if e := bi.WaitSeq(ctx, goal); e != nil {
		return nil, e
	}
	return bi.Get(views.Content, content_view.ContentVotes, &content_view.ContentVotesIn{
		Perspective: uf.User.PubKey.Hex(),
		ContentHash: in.PostRefStr,
	})
}
