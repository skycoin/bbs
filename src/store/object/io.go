package object

import (
	"github.com/skycoin/bbs/src/misc/tag"
	"github.com/skycoin/skycoin/src/cipher"
	"time"
)

// NewBoard represents io required to create a new board.
type NewBoardIO struct {
	Name        string        `bbs:"name"`
	Body        string        `bbs:"body"`
	SubAddrsStr string        `bbs:"subAddrsStr"`
	SubAddrs    []string      `bbs:"subAddrs"`
	Seed        string        `bbs:"bSeed"`
	BoardPubKey cipher.PubKey `bbs:"bpk"`
	BoardSecKey cipher.SecKey `bbs:"bsk"`
	Board       *Board
}

func (a *NewBoardIO) Process() error {
	if e := tag.Process(a); e != nil {
		return e
	}
	a.Board = &Board{
		Created: time.Now().UnixNano(),
	}
	SetData(a.Board, &ContentData{
		Name:         a.Name,
		Body:         a.Body,
		SubAddresses: a.SubAddrs,
	})
	return nil
}

// NewThread represents io required to create a new thread.
type NewThreadIO struct {
	Title string `bbs:"heading"`
	Body  string `bbs:"body"`
}

// NewUser represents io required when creating a new user.
type NewUserIO struct {
	Alias      string        `bbs:"alias" trans:"alias"`
	Seed       string        `bbs:"uSeed"`
	Password   string        `bbs:"password"`
	UserPubKey cipher.PubKey `bbs:"upk" trans:"upk"`
	UserSecKey cipher.SecKey `bbs:"usk" trans:"usk"`
	File       *UserFile
}

func (a *NewUserIO) Process() error {
	if e := tag.Process(a); e != nil {
		return e
	}
	a.File = &UserFile{
		User: User{
			Alias:  a.Alias,
			PubKey: a.UserPubKey,
			SecKey: a.UserSecKey,
		},
		Seed: a.Seed,
	}
	return nil
}

// Login represents input required when logging io.
type LoginIO struct {
	Alias string `bbs:"alias"`
}

func (a *LoginIO) Process() error {
	if e := tag.Process(a); e != nil {
		return e
	}
	return nil
}

// Connection represents input/output required when connection/disconnecting from address.
type ConnectionIO struct {
	Address string `bbs:"address"`
}

func (a *ConnectionIO) Process() error {
	if e := tag.Process(a); e != nil {
		return e
	}
	return nil
}

// Subscription represents a subscription input.
type SubscriptionIO struct {
	PubKeyStr string        `bbs:"bpkStr"`
	PubKey    cipher.PubKey `bbs:"bpk"`
	SecKey    cipher.SecKey `bbs:"bsk"`
}

func (a *SubscriptionIO) Process() error {
	if e := tag.Process(a); e != nil {
		return e
	}
	return nil
}
