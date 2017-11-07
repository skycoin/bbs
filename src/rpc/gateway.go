package rpc

import (
	"context"
	"encoding/json"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store"
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
	<<< CONNECTIONS : MESSENGER >>>
*/

func (g *Gateway) GetMessengerConnections(_ *struct{}, out *string) error {
	return send(out)(g.Access.GetMessengerConnections(context.Background()))
}

func (g *Gateway) NewMessengerConnection(in *store.ConnectionIn, out *string) error {
	return send(out)(g.Access.NewMessengerConnection(context.Background(), in))
}

func (g *Gateway) DeleteMessengerConnection(in *store.ConnectionIn, out *string) error {
	return send(out)(g.Access.DeleteMessengerConnection(context.Background(), in))
}

func (g *Gateway) Discover(_ *struct{}, out *string) error {
	return send(out)(g.Access.GetAvailableBoards(context.Background()))
}

/*
	<<< CONNECTIONS >>>
*/

func (g *Gateway) GetConnections(_ *struct{}, out *string) error {
	return send(out)(g.Access.GetConnections(context.Background()))
}

func (g *Gateway) NewConnection(in *store.ConnectionIn, out *string) error {
	return send(out)(g.Access.NewConnection(context.Background(), in))
}

func (g *Gateway) DeleteConnection(in *store.ConnectionIn, out *string) error {
	return send(out)(g.Access.DeleteConnection(context.Background(), in))
}

/*
	<<< SUBSCRIPTIONS >>>
*/

func (g *Gateway) GetSubscriptions(_ *struct{}, out *string) error {
	return send(out)(g.Access.GetSubscriptions(context.Background()))
}

func (g *Gateway) NewSubscription(in *store.BoardIn, out *string) error {
	return send(out)(g.Access.NewSubscription(context.Background(), in))
}

func (g *Gateway) DeleteSubscription(in *store.BoardIn, out *string) error {
	return send(out)(g.Access.DeleteSubscription(context.Background(), in))
}

/*
	<<< CONTENT : ADMIN >>>
*/

func (g *Gateway) NewBoard(in *store.NewBoardIn, out *string) error {
	return send(out)(g.Access.NewBoard(context.Background(), in))
}

func (g *Gateway) DeleteBoard(in *store.BoardIn, out *string) error {
	return send(out)(g.Access.DeleteBoard(context.Background(), in))
}

func (g *Gateway) ExportBoard(in *store.ExportBoardIn, out *string) error {
	return send(out)(g.Access.ExportBoard(context.Background(), in))
}

func (g *Gateway) ImportBoard(in *store.ImportBoardIn, out *string) error {
	return send(out)(g.Access.ImportBoard(context.Background(), in))
}

/*
	<<< CONTENT >>>
*/

func (g *Gateway) GetBoards(_ *struct{}, out *string) error {
	return send(out)(g.Access.GetBoards(context.Background()))
}

func (g *Gateway) GetBoard(in *store.BoardIn, out *string) error {
	return send(out)(g.Access.GetBoard(context.Background(), in))
}

func (g *Gateway) GetBoardPage(in *store.BoardIn, out *string) error {
	return send(out)(g.Access.GetBoardPage(context.Background(), in))
}

func (g *Gateway) GetThreadPage(in *store.ThreadIn, out *string) error {
	return send(out)(g.Access.GetThreadPage(context.Background(), in))
}

func (g *Gateway) GetFollowPage(in *store.UserIn, out *string) error {
	return send(out)(g.Access.GetFollowPage(context.Background(), in))
}

/*
	<<< CONTENT : SUBMISSION >>>
*/

func (g *Gateway) NewThread(in *store.NewThreadIn, out *string) error {
	return send(out)(g.Access.NewThread(context.Background(), in))
}

func (g *Gateway) NewPost(in *store.NewPostIn, out *string) error {
	return send(out)(g.Access.NewPost(context.Background(), in))
}

func (g *Gateway) VoteThread(in *store.VoteThreadIn, out *string) error {
	return send(out)(g.Access.VoteThread(context.Background(), in))
}

func (g *Gateway) VotePost(in *store.VotePostIn, out *string) error {
	return send(out)(g.Access.VotePost(context.Background(), in))
}

func (g *Gateway) VoteUser(in *store.VoteUserIn, out *string) error {
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
