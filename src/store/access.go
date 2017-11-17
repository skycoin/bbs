package store

import (
	"context"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/typ"
	"github.com/skycoin/bbs/src/store/cxo"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state"
	"github.com/skycoin/skycoin/src/util/file"
	"log"
	"math"
	"os"
	"time"
)

type Access struct {
	CXO *cxo.Manager
}

func (a *Access) SubmitContent(ctx context.Context, in *SubmissionIn) (interface{}, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}

	transport, e := object.NewTransport(in.Body, in.Sig)
	if e != nil {
		return nil, e
	}

	bi, e := submitAndWait(ctx, a, transport)
	if e != nil {
		return nil, e
	}

	switch transport.Body.Type {
	case object.V5ThreadType, object.V5PostType:
		return getNewContentOut(transport)

	case object.V5ThreadVoteType:
		return getThreadVoteOut(transport, bi)

	case object.V5PostVoteType:
		return getPostVoteOut(transport, bi)

	case object.V5UserVoteType:
		return getUserVoteOut(transport, bi)

	default:
		return nil, boo.Newf(boo.InvalidInput,
			"content submission of type '%s' is invalid", transport.Body.Type)
	}
}

func submitAndWait(ctx context.Context, a *Access, transport *object.Transport) (*state.BoardInstance, error) {
	bi, e := a.CXO.GetBoardInstance(transport.GetOfBoard())
	if e != nil {
		return nil, e
	}
	var goal uint64
	if bi.IsMaster() {
		if goal, e = bi.Submit(transport); e != nil {
			return nil, e
		}
	} else {
		if goal, e = a.CXO.SubmitToRemote(ctx, bi.GetSubmissionKeys(), transport); e != nil {
			return nil, e
		}
	}
	return bi, bi.WaitSeq(ctx, goal)
}

/*
	<<< CONNECTIONS : MESSENGER >>>
*/

func (a *Access) GetMessengerConnections(ctx context.Context) (*MessengersOut, error) {
	return getMessengersOut(ctx, a.CXO.GetMessengers()), nil
}

func (a *Access) NewMessengerConnection(ctx context.Context, in *ConnectionIn) (*MessengersOut, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	if e := a.CXO.ConnectToMessenger(in.Address); e != nil {
		return nil, e
	}
	return a.GetMessengerConnections(ctx)
}

func (a *Access) DeleteMessengerConnection(ctx context.Context, in *ConnectionIn) (*MessengersOut, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	if e := a.CXO.DisconnectFromMessenger(in.Address); e != nil {
		return nil, e
	}
	return a.GetMessengerConnections(ctx)
}

func (a *Access) GetAvailableBoards(ctx context.Context) (*AvailableBoardsOut, error) {
	return getAvailableBoardsOut(a.CXO.GetAvailableBoards()), nil
}

/*
	<<< CONNECTIONS >>>
*/

func (a *Access) GetConnections(ctx context.Context) (*ConnectionsOut, error) {
	return getConnectionsOut(ctx, a.CXO.GetActiveConnections(), a.CXO.GetSavedConnections()), nil
}

func (a *Access) NewConnection(ctx context.Context, in *ConnectionIn) (*ConnectionsOut, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	if e := a.CXO.Connect(in.Address); e != nil {
		return nil, e
	}
	time.Sleep(time.Second)
	return a.GetConnections(ctx)
}

func (a *Access) DeleteConnection(ctx context.Context, in *ConnectionIn) (*ConnectionsOut, error) {
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

func (a *Access) GetSubscriptions(ctx context.Context) (*SubscriptionsOut, error) {
	return getSubscriptionsOut(ctx, a.CXO.GetSubscriptions()), nil
}

func (a *Access) NewSubscription(ctx context.Context, in *BoardIn) (*SubscriptionsOut, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	if e := a.CXO.SubscribeRemote(in.PubKey); e != nil {
		return nil, e
	}
	return a.GetSubscriptions(ctx)
}

func (a *Access) DeleteSubscription(ctx context.Context, in *BoardIn) (*SubscriptionsOut, error) {
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

func (a *Access) NewBoard(ctx context.Context, in *NewBoardIn) (*BoardsOut, error) {
	if e := in.Process(a.CXO.Relay().SubmissionKeys()); e != nil {
		return nil, e
	}
	if e := a.CXO.NewBoard(in.Content, in.BoardPubKey, in.BoardSecKey); e != nil {
		return nil, e
	}
	return a.GetBoards(ctx)
}

func (a *Access) DeleteBoard(ctx context.Context, in *BoardIn) (*BoardsOut, error) {
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

func (a *Access) ExportBoard(ctx context.Context, in *ExportBoardIn) (*ExportBoardOut, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}

	out, e := a.CXO.ExportBoard(in.PubKey, in.FilePath)
	if e != nil {
		return nil, e
	}
	if e := file.SaveJSON(in.FilePath, out, os.FileMode(0600)); e != nil {
		return nil, e
	}
	return getExportBoardOut(in.FilePath, out), nil
}

func (a *Access) ImportBoard(ctx context.Context, in *ImportBoardIn) (*ExportBoardOut, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	pagesIn := new(object.PagesJSON)
	if e := file.LoadJSON(in.FilePath, pagesIn); e != nil {
		return nil, e
	}
	if e := a.CXO.ImportBoard(ctx, pagesIn); e != nil {
		return nil, e
	}
	return getExportBoardOut(in.FilePath, pagesIn), nil
}

/*
	<<< CONTENT >>>
*/

func (a *Access) GetBoards(ctx context.Context) (*BoardsOut, error) {
	m, r, e := a.CXO.GetBoards(ctx)
	if e != nil {
		return nil, e
	}
	return getBoardsOut(ctx, m, r), nil
}

func (a *Access) GetBoard(ctx context.Context, in *BoardIn) (*BoardOut, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.PubKey)
	if e != nil {
		return nil, e
	}
	board, e := bi.Viewer().GetBoard()
	if e != nil {
		return nil, e
	}
	return getBoardOut(board), nil
}

func (a *Access) GetBoardPage(ctx context.Context, in *BoardIn) (interface{}, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.PubKey)
	if e != nil {
		return nil, e
	}
	return bi.Viewer().GetBoardPage(&state.BoardPageIn{
		Perspective:    in.UserPubKeyStr,
		PaginatedInput: typ.PaginatedInput{PageSize: math.MaxUint64},
	})
}

func (a *Access) NewThread(ctx context.Context, in *NewThreadIn) (interface{}, error) {
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
		goal, e = a.CXO.SubmitToRemote(ctx, bi.GetSubmissionKeys(), in.Transport)
		if e != nil {
			return nil, e
		}
	}
	if e := bi.WaitSeq(ctx, goal); e != nil {
		return nil, e
	}
	return getNewContentOut(in.Transport)
}

func (a *Access) GetThreadPage(ctx context.Context, in *ThreadIn) (interface{}, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.BoardPubKey)
	if e != nil {
		return nil, e
	}
	return bi.Viewer().GetThreadPage(&state.ThreadPageIn{
		Perspective:    in.UserPubKeyStr,
		ThreadHash:     in.ThreadRefStr,
		PaginatedInput: typ.PaginatedInput{PageSize: math.MaxUint64},
	})
}

func (a *Access) NewPost(ctx context.Context, in *NewPostIn) (interface{}, error) {
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
		goal, e = a.CXO.SubmitToRemote(ctx, bi.GetSubmissionKeys(), in.Transport)
		if e != nil {
			return nil, e
		}
	}
	if e := bi.WaitSeq(ctx, goal); e != nil {
		return nil, e
	}
	return getNewContentOut(in.Transport)
}

func (a *Access) GetParticipants(ctx context.Context, in *BoardIn) (interface{}, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.PubKey)
	if e != nil {
		return nil, e
	}
	return bi.Viewer().GetParticipants()
}

/*
	<<< VOTES >>>
*/

func (a *Access) GetFollowPage(ctx context.Context, in *UserIn) (interface{}, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	bi, e := a.CXO.GetBoardInstance(in.BoardPubKey)
	if e != nil {
		return nil, e
	}
	return bi.Viewer().GetUserProfile(&state.UserProfileIn{
		UserPubKey: in.UserPubKeyStr,
	})
}

func (a *Access) VoteUser(ctx context.Context, in *VoteUserIn) (interface{}, error) {
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
		goal, e = a.CXO.SubmitToRemote(ctx, bi.GetSubmissionKeys(), in.Transport)
		if e != nil {
			return nil, e
		}
	}
	if e := bi.WaitSeq(ctx, goal); e != nil {
		return nil, e
	}
	return getUserVoteOut(in.Transport, bi)
}

func (a *Access) VoteThread(ctx context.Context, in *VoteThreadIn) (interface{}, error) {
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
		goal, e = a.CXO.SubmitToRemote(ctx, bi.GetSubmissionKeys(), in.Transport)
		if e != nil {
			return nil, e
		}
	}
	if e := bi.WaitSeq(ctx, goal); e != nil {
		return nil, e
	}
	return getThreadVoteOut(in.Transport, bi)
}

func (a *Access) VotePost(ctx context.Context, in *VotePostIn) (interface{}, error) {
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
		goal, e = a.CXO.SubmitToRemote(ctx, bi.GetSubmissionKeys(), in.Transport)
		if e != nil {
			return nil, e
		}
	}
	if e := bi.WaitSeq(ctx, goal); e != nil {
		return nil, e
	}
	return getPostVoteOut(in.Transport, bi)
}
