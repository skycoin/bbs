package rpc

import (
	"context"
	"encoding/json"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store"
	"github.com/skycoin/bbs/src/store/object"
)

var (
	ErrEmptyInput = boo.New(boo.InvalidInput, "empty input error")
)

type Gateway struct {
	Access   *store.Access
	QuitChan chan int
}

func (g *Gateway) Quit(code int, ok *bool) error {
	if ok == nil {
		return ErrEmptyInput
	}
	g.QuitChan <- code
	*ok = true
	return nil
}

/*
	<<< CONNECTIONS >>>
*/

func (g *Gateway) GetConnections(_ *struct{}, out *string) error {
	return send(out)(g.Access.GetConnections(context.Background()))
}

func (g *Gateway) NewConnection(in *object.ConnectionIO, out *string) error {
	return send(out)(g.Access.NewConnection(context.Background(), in))
}

func (g *Gateway) DeleteConnection(in *object.ConnectionIO, out *string) error {
	return send(out)(g.Access.DeleteConnection(context.Background(), in))
}

/*
	<<< SUBSCRIPTIONS >>>
*/

func (g *Gateway) GetSubscriptions(_ *struct{}, out *string) error {
	return send(out)(g.Access.GetSubscriptions(context.Background()))
}

func (g *Gateway) NewSubscription(in *object.BoardIO, out *string) error {
	return send(out)(g.Access.NewSubscription(context.Background(), in))
}

func (g *Gateway) DeleteSubscription(in *object.BoardIO, out *string) error {
	return send(out)(g.Access.DeleteSubscription(context.Background(), in))
}

/*
	<<< CONTENT : ADMIN >>>
*/

func (g *Gateway) NewBoard(in *object.NewBoardIO, out *string) error {
	return send(out)(g.Access.NewBoard(context.Background(), in))
}

func (g *Gateway) DeleteBoard(in *object.BoardIO, out *string) error {
	return send(out)(g.Access.DeleteBoard(context.Background(), in))
}

//func (g *Gateway) ExportBoard(in *object.ExportBoardIO, out *string) error {
//	return send(out)(g.Access.ExportBoard(context.Background(), in))
//}
//
//func (g *Gateway) ImportBoard(in *object.ExportBoardIO, out *string) error {
//	return send(out)(g.Access.ImportBoard(context.Background(), in))
//}

/*
	<<< CONTENT >>>
*/

func (g *Gateway) GetBoards(_ *struct{}, out *string) error {
	return send(out)(g.Access.GetBoards(context.Background()))
}

func (g *Gateway) GetBoard(in *object.BoardIO, out *string) error {
	return send(out)(g.Access.GetBoard(context.Background(), in))
}

func (g *Gateway) GetBoardPage(in *object.BoardIO, out *string) error {
	return send(out)(g.Access.GetBoardPage(context.Background(), in))
}

func (g *Gateway) GetThreadPage(in *object.ThreadIO, out *string) error {
	return send(out)(g.Access.GetThreadPage(context.Background(), in))
}

func (g *Gateway) GetFollowPage(in *object.UserIO, out *string) error {
	return send(out)(g.Access.GetFollowPage(context.Background(), in))
}

/*
	<<< CONTENT : SUBMISSION >>>
*/

func (g *Gateway) NewThread(in *object.NewThreadIO, out *string) error {
	return send(out)(g.Access.NewThread(context.Background(), in))
}

func (g *Gateway) NewPost(in *object.NewPostIO, out *string) error {
	return send(out)(g.Access.NewPost(context.Background(), in))
}

func (g *Gateway) VoteThread(in *object.ThreadVoteIO, out *string) error {
	return send(out)(g.Access.VoteThread(context.Background(), in))
}

func (g *Gateway) VotePost(in *object.PostVoteIO, out *string) error {
	return send(out)(g.Access.VotePost(context.Background(), in))
}

func (g *Gateway) VoteUser(in *object.UserVoteIO, out *string) error {
	return send(out)(g.Access.VoteUser(context.Background(), in))
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func send(out *string) func(v interface{}, e error) error {
	return func(v interface{}, e error) error {
		if e != nil {
			return e
		} else if data, e := json.MarshalIndent(v, "", "  "); e != nil {
			return e
		} else {
			*out = string(data)
			return nil
		}
	}
}
