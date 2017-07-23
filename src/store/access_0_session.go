package store

import (
	"context"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state"
)

func (a *Access) GetUsers(ctx context.Context) (*UsersOutput, error) {
	aliases, e := a.Session.GetUsers(ctx)
	if e != nil {
		return nil, e
	}
	out := &UsersOutput{
		Users: make([]object.UserView, len(aliases)),
	}
	for i, alias := range aliases {
		out.Users[i] = object.UserView{
			User: object.User{Alias: alias},
		}
	}
	return out, nil
}

func (a *Access) NewUser(ctx context.Context, in *object.NewUserIO) (*UsersOutput, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	if _, e := a.Session.NewUser(ctx, in); e != nil {
		return nil, e
	}
	return a.GetUsers(ctx)
}

func (a *Access) DeleteUser(ctx context.Context, alias string) (*UsersOutput, error) {
	if e := a.Session.DeleteUser(ctx, alias); e != nil {
		return nil, e
	}
	return a.GetUsers(ctx)
}

func (a *Access) Login(ctx context.Context, in *object.LoginIO) (*object.UserView, error) {
	if e := in.Process(); e != nil {
		return nil, e
	}
	file, e := a.Session.Login(ctx, in)
	if e != nil {
		return nil, e
	}
	out := &object.UserView{
		User:      object.User{Alias: file.User.Alias},
		PublicKey: file.User.PublicKey.Hex(),
		SecretKey: file.User.SecretKey.Hex(),
	}
	return out, nil
}

func (a *Access) Logout(ctx context.Context) error {
	return a.Session.Logout(ctx)
}

func (a *Access) GetSession(ctx context.Context) (*state.UserFileView, error) {
	file, e := a.Session.GetInfo(ctx)
	if e != nil {
		return nil, e
	}
	view := file.GenerateView(a.Session.GetCXO())
	return view, nil
}
