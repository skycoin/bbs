package session

import "github.com/skycoin/bbs/src/store/object"

type Store interface {
	GetUsers() ([]string, error)
	GetUser(alias string) (*object.UserFile, bool)
	NewUser(alias string, file *object.UserFile) error
	DeleteUser(alias string) error
}
