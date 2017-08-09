package store

import (
	"github.com/skycoin/bbs/src/store/cxo"
	"github.com/skycoin/bbs/src/store/session"
	"context"
	"github.com/skycoin/bbs/src/store/io"
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

/*
	<<<  >>>
*/