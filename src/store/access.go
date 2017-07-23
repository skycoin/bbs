package store

import (
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state"
)

// Access allows access to store.
type Access struct {
	Session *state.Session
}

type UsersOutput struct {
	Users []object.UserView `json:"users"`
}
