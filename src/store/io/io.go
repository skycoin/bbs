package io

import (
	"github.com/skycoin/bbs/src/misc/tag"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/skycoin/src/cipher"
)

// NewUser represents io required when creating a new user.
type NewUser struct {
	Alias      string        `bbs:"alias" trans:"alias"`
	Seed       string        `bbs:"uSeed"`
	Password   string        `bbs:"password"`
	UserPubKey cipher.PubKey `bbs:"upk" trans:"upk"`
	UserSecKey cipher.SecKey `bbs:"usk" trans:"usk"`
	File       *object.UserFile
}

func (a *NewUser) Process() error {
	if e := tag.Process(a); e != nil {
		return e
	}
	a.File = &object.UserFile{
		User: object.User{
			Alias:  a.Alias,
			PubKey: a.UserPubKey,
			SecKey: a.UserSecKey,
		},
		Seed: a.Seed,
	}
	return nil
}

// Login represents input required when logging io.
type Login struct {
	Alias    string `bbs:"alias"`
	Password string `bbs:"password"`
}

func (a *Login) Process() error {
	if e := tag.Process(a); e != nil {
		return e
	}
	return nil
}
