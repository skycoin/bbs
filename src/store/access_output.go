package store

import (
	"context"
	"github.com/skycoin/bbs/src/store/object"
)

type UsersOutput struct {
	Users []object.UserView `json:"users"`
}

func getUsers(_ context.Context, aliases []string) *UsersOutput {
	out := &UsersOutput{
		Users: make([]object.UserView, len(aliases)),
	}
	for i, alias := range aliases {
		out.Users[i] = object.UserView{
			User: object.User{Alias: alias},
		}
	}
	return out
}