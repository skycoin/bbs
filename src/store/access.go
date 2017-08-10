package store

import (
	"context"
	"github.com/skycoin/bbs/src/store/cxo"
	"github.com/skycoin/bbs/src/store/io"
	"github.com/skycoin/bbs/src/store/session"
	"time"
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

func (a *Access) NewUser(ctx context.Context, in *io.NewUser) (*UsersOutput, error) {
	if e := in.Process(); e != nil {
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

func (a *Access) Login(ctx context.Context, in *io.Login) (*SessionOutput, error) {
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

func (a *Access) NewConnection(ctx context.Context, in *io.Connection) (*ConnectionsOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	if e := a.CXO.Connect(in.Address); e != nil {
		return nil, e
	}
	time.Sleep(time.Second)
	return a.GetConnections(ctx)
}

func (a *Access) DeleteConnection(ctx context.Context, in *io.Connection) (*ConnectionsOutput, error) {
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

func (a *Access) NewSubscription(ctx context.Context, in *io.Subscription) (*SubscriptionsOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	if e := a.CXO.SubscribeRemote(in.PubKey); e != nil {
		return nil, e
	}
	return a.GetSubscriptions(ctx)
}

func (a *Access) DeleteSubscription(ctx context.Context, in *io.Subscription) (*SubscriptionsOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	if e := a.CXO.UnsubscribeRemote(in.PubKey); e != nil {
		return nil, e
	}
	return a.GetSubscriptions(ctx)
}
