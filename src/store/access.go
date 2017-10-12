package store

import (
	"context"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/cxo"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/object/revisions/r0"
	"github.com/skycoin/bbs/src/store/state"
	"github.com/skycoin/bbs/src/store/state/views"
	"github.com/skycoin/bbs/src/store/state/views/content_view"
	"github.com/skycoin/bbs/src/store/state/views/follow_view"
	"log"
	"time"
)

type Access struct {
	CXO *cxo.Manager
}

func (a *Access) SubmitContent(ctx context.Context, in *object.SubmissionIO) (interface{}, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}

	transport, e := r0.NewTransport(in.Body, in.Sig)
	if e != nil {
		return nil, e
	}

	bi, e := submitAndWait(ctx, a, transport)
	if e != nil {
		return nil, e
	}

	switch transport.Body.Type {
	case r0.V5ThreadType:
		return bi.Get(views.Content, content_view.BoardPage, &content_view.BoardPageIn{
			Perspective: transport.Body.Creator,
		})

	case r0.V5PostType:
		return bi.Get(views.Content, content_view.ThreadPage, &content_view.ThreadPageIn{
			Perspective: transport.Body.Creator,
			ThreadHash:  transport.Body.OfThread,
		})

	case r0.V5ThreadVoteType:
		return bi.Get(views.Content, content_view.ContentVotes, &content_view.ContentVotesIn{
			Perspective: transport.Body.Creator,
			ContentHash: transport.Body.OfThread,
		})

	case r0.V5PostVoteType:
		return bi.Get(views.Content, content_view.ContentVotes, &content_view.ContentVotesIn{
			Perspective: transport.Body.Creator,
			ContentHash: transport.Body.OfPost,
		})

	case r0.V5UserVoteType:
		out, e := bi.Get(views.Follow, follow_view.FollowPage, transport.Body.Creator)
		if e != nil {
			return nil, e
		}
		return getFollowPageOutput(out), nil

	default:
		return nil, boo.Newf(boo.InvalidInput,
			"content submission of type '%s' is invalid", transport.Body.Type)
	}
}

func submitAndWait(ctx context.Context, a *Access, transport *r0.Transport) (*state.BoardInstance, error) {
	ofBoard := transport.GetOfBoard()

	bi, e := a.CXO.GetBoardInstance(ofBoard)
	if e != nil {
		return nil, e
	}
	var goal uint64
	if bi.IsMaster() {
		if goal, e = bi.Submit(transport); e != nil {
			return nil, e
		}
	} else {
		if goal, e = a.CXO.Relay().NewContent(ctx, bi.GetSubmissionKeys(), transport.Content); e != nil {
			return nil, e
		}
	}
	return bi, bi.WaitSeq(ctx, goal)
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

func (a *Access) GetBoards(ctx context.Context) (*BoardsOutput, error) {
	m, r, e := a.CXO.GetBoards(ctx)
	if e != nil {
		return nil, e
	}
	return getBoardsOutput(ctx, m, r), nil
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

func (a *Access) GetBoardPage(ctx context.Context, in *object.BoardIO) (interface{}, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.PubKey)
	if e != nil {
		return nil, e
	}
	return bi.Get(views.Content, content_view.BoardPage, &content_view.BoardPageIn{
		Perspective: in.UserPubKeyStr,
	})
}

func (a *Access) NewThread(ctx context.Context, in *object.NewThreadIO) (interface{}, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.BoardPubKey)
	if e != nil {
		return nil, e
	}
	var goal uint64
	if bi.IsMaster() {
		if goal, e = bi.Submit(in.Transport); e != nil {
			return nil, e
		}
	} else {
		log.Println("NewThread: subs:", bi.GetSubmissionKeys())
		goal, e = a.CXO.Relay().NewContent(ctx, bi.GetSubmissionKeys(), in.Transport.Content)
		if e != nil {
			return nil, e
		}
	}
	if e := bi.WaitSeq(ctx, goal); e != nil {
		return nil, e
	}
	return bi.Get(views.Content, content_view.BoardPage, &content_view.BoardPageIn{
		Perspective: in.UserPubKeyStr,
	})
}

func (a *Access) GetThreadPage(ctx context.Context, in *object.ThreadIO) (interface{}, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.BoardPubKey)
	if e != nil {
		return nil, e
	}
	return bi.Get(views.Content, content_view.ThreadPage, &content_view.ThreadPageIn{
		Perspective: in.UserPubKeyStr,
		ThreadHash:  in.ThreadRefStr,
	})
}

func (a *Access) NewPost(ctx context.Context, in *object.NewPostIO) (interface{}, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.BoardPubKey)
	if e != nil {
		return nil, e
	}
	var goal uint64
	if bi.IsMaster() {
		if goal, e = bi.Submit(in.Transport); e != nil {
			return nil, e
		}
	} else {
		goal, e = a.CXO.Relay().NewContent(ctx, bi.GetSubmissionKeys(), in.Transport.Content)
		if e != nil {
			return nil, e
		}
	}
	if e := bi.WaitSeq(ctx, goal); e != nil {
		return nil, e
	}
	return bi.Get(views.Content, content_view.ThreadPage, &content_view.ThreadPageIn{
		Perspective: in.UserPubKeyStr,
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
	out, e := bi.Get(views.Follow, follow_view.FollowPage, in.UserPubKeyStr)
	if e != nil {
		return nil, e
	}
	return getFollowPageOutput(out), nil
}

func (a *Access) VoteUser(ctx context.Context, in *object.UserVoteIO) (interface{}, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.BoardPubKey)
	if e != nil {
		return nil, e
	}
	var goal uint64
	if bi.IsMaster() {
		if goal, e = bi.Submit(in.Transport); e != nil {
			return nil, e
		}
	} else {
		goal, e = a.CXO.Relay().NewContent(ctx, bi.GetSubmissionKeys(), in.Transport.Content)
		if e != nil {
			return nil, e
		}
	}
	if e := bi.WaitSeq(ctx, goal); e != nil {
		return nil, e
	}
	out, e := bi.Get(views.Follow, follow_view.FollowPage, in.UserPubKeyStr)
	if e != nil {
		return nil, e
	}
	return getFollowPageOutput(out), nil
}

func (a *Access) VoteThread(ctx context.Context, in *object.ThreadVoteIO) (interface{}, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.BoardPubKey)
	if e != nil {
		return nil, e
	}
	var goal uint64
	if bi.IsMaster() {
		if goal, e = bi.Submit(in.Transport); e != nil {
			return nil, e
		}
	} else {
		goal, e = a.CXO.Relay().NewContent(ctx, bi.GetSubmissionKeys(), in.Transport.Content)
		if e != nil {
			return nil, e
		}
	}
	if e := bi.WaitSeq(ctx, goal); e != nil {
		return nil, e
	}
	return bi.Get(views.Content, content_view.ContentVotes, &content_view.ContentVotesIn{
		Perspective: in.UserPubKeyStr,
		ContentHash: in.ThreadRefStr,
	})
}

func (a *Access) VotePost(ctx context.Context, in *object.PostVoteIO) (interface{}, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.BoardPubKey)
	if e != nil {
		return nil, e
	}
	var goal uint64
	if bi.IsMaster() {
		if goal, e = bi.Submit(in.Transport); e != nil {
			return nil, e
		}
	} else {
		goal, e = a.CXO.Relay().NewContent(ctx, bi.GetSubmissionKeys(), in.Transport.Content)
		if e != nil {
			return nil, e
		}
	}
	if e := bi.WaitSeq(ctx, goal); e != nil {
		return nil, e
	}
	return bi.Get(views.Content, content_view.ContentVotes, &content_view.ContentVotesIn{
		Perspective: in.UserPubKeyStr,
		ContentHash: in.PostRefStr,
	})
}
