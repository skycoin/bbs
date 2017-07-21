package store

import (
	"github.com/skycoin/bbs/src/store/obj"
	"github.com/skycoin/bbs/src/store/state"
)

// Access allows access to store.
type Access struct {
	Session *state.Session
}

type UsersOutput struct {
	Users []obj.UserView `json:"users"`
}
