package rpc

import (
	"context"
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

func (g *Gateway) GetConnections(_ *struct{}, out *store.ConnectionsOutput) error {
	if a, e := g.Access.GetConnections(context.Background()); e != nil {
		return e
	} else {
		*out = *a
		return nil
	}
}

func (g *Gateway) NewConnection(in *object.ConnectionIO, out *store.ConnectionsOutput) error {
	if a, e := g.Access.NewConnection(context.Background(), in); e != nil {
		return e
	} else {
		*out = *a
		return nil
	}
}

func (g *Gateway) DeleteConnection(in *object.ConnectionIO, out *store.ConnectionsOutput) error {
	if a, e := g.Access.DeleteConnection(context.Background(), in); e != nil {
		return e
	} else {
		*out = *a
		return nil
	}
}

/*
	<<< SUBSCRIPTIONS >>>
*/

func (g *Gateway) GetSubscriptions(_ *struct{}, out *store.SubscriptionsOutput) error {
	if a, e := g.Access.GetSubscriptions(context.Background()); e != nil {
		return e
	} else {
		*out = *a
		return nil
	}
}

func (g *Gateway) NewSubscription(in *object.BoardIO, out *store.SubscriptionsOutput) error {
	if a, e := g.Access.NewSubscription(context.Background(), in); e != nil {
		return e
	} else {
		*out = *a
		return nil
	}
}

func (g *Gateway) DeleteSubscription(in *object.BoardIO, out *store.SubscriptionsOutput) error {
	if a, e := g.Access.DeleteSubscription(context.Background(), in); e != nil {
		return e
	} else {
		*out = *a
		return nil
	}
}

/*
	<<< CONTENT : ADMIN >>>
*/

func (g *Gateway) NewBoard(in *object.NewBoardIO, out *store.BoardsOutput) error {
	if a, e := g.Access.NewBoard(context.Background(), in); e != nil {
		return e
	} else {
		*out = *a
		return nil
	}
}

func (g *Gateway) DeleteBoard(in *object.BoardIO, out *store.BoardsOutput) error {
	if a, e := g.Access.DeleteBoard(context.Background(), in); e != nil {
		return e
	} else {
		*out = *a
		return nil
	}
}

func (g *Gateway) ExportBoard(in *object.ExportBoardIO, out *store.ExportBoardOutput) error {
	if a, e := g.Access.ExportBoard(context.Background(), in); e != nil {
		return e
	} else {
		*out = *a
		return nil
	}
}

func (g *Gateway) ImportBoard(in *object.ExportBoardIO, out *store.ExportBoardOutput) error {
	if a, e := g.Access.ImportBoard(context.Background(), in); e != nil {
		return e
	} else {
		*out = *a
		return nil
	}
}
