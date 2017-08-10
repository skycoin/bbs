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
	Alias string `bbs:"alias"`
}

func (a *Login) Process() error {
	if e := tag.Process(a); e != nil {
		return e
	}
	return nil
}

// Connection represents input/output required when connection/disconnecting from address.
type Connection struct {
	Address string `bbs:"address"`
}

func (a *Connection) Process() error {
	if e := tag.Process(a); e != nil {
		return e
	}
	return nil
}

// Subscription represents a subscription input.
type Subscription struct {
	PubKeyStr string        `bbs:"bpkStr"`
	PubKey    cipher.PubKey `bbs:"bpk"`
	SecKey    cipher.SecKey `bbs:"bsk"`
}

func (a *Subscription) Process() error {
	if e := tag.Process(a); e != nil {
		return e
	}
	return nil
}
