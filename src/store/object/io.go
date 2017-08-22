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
	BoardPubKeyStr string        `bbs:"bpkStr"`
	BoardPubKey    cipher.PubKey `bbs:"bpk"`
	Name           string        `bbs:"name"`
	Body           string        `bbs:"body"`
	Thread         *Thread
}

func (a *NewThreadIO) Process(upk cipher.PubKey, usk cipher.SecKey) error {
	if e := tag.Process(a); e != nil {
		return e
	}
	a.Thread = &Thread{
		OfBoard: a.BoardPubKey,
		Created: time.Now().UnixNano(),
		Creator: upk,
	}
	SetData(a.Thread, &ContentData{
		Name: a.Name,
		Body: a.Body,
	})
	tag.Sign(a.Thread, upk, usk)
	return nil
}

type NewPostIO struct {
	BoardPubKeyStr string        `bbs:"bpkStr"`
	BoardPubKey    cipher.PubKey `bbs:"bpk"`
	ThreadRefStr   string        `bbs:"tRefStr"`
	ThreadRef      cipher.SHA256 `bbs:"tRef"`
	PostRefStr     string        `bbs:"pRefStr"`
	PostRef        cipher.SHA256 `bbs:"pRef"`
	Name           string        `bbs:"name"`
	Body           string        `bbs:"body"`
	Post           *Post
}

func (a *NewPostIO) Process(upk cipher.PubKey, usk cipher.SecKey) error {
	if e := tag.Process(a); e != nil {
		return e
	}
	a.Post = &Post{
		OfBoard:  a.BoardPubKey,
		OfThread: a.ThreadRef,
		OfPost:   a.PostRef,
		Created:  time.Now().UnixNano(),
		Creator:  upk,
	}
	SetData(a.Post, &ContentData{
		Name: a.Name,
		Body: a.Body,
	})
	tag.Sign(a.Post, upk, usk)
	return nil
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

// BoardIO represents a subscription input.
type BoardIO struct {
	PubKeyStr string        `bbs:"bpkStr"`
	PubKey    cipher.PubKey `bbs:"bpk"`
}

func (a *BoardIO) Process() error {
	if e := tag.Process(a); e != nil {
		return e
	}
	return nil
}

type ThreadIO struct {
	BoardPubKeyStr string        `bbs:"bpkStr"`
	BoardPubKey    cipher.PubKey `bbs:"bpk"`
	ThreadRefStr   string        `bbs:"tRefStr"`
	ThreadRef      cipher.SHA256 `bbs:"tRef"`
}

func (a *ThreadIO) Process() error {
	if e := tag.Process(a); e != nil {
		return e
	}
	return nil
}

type UserVoteIO struct {
	BoardPubKeyStr string        `bbs:"bpkStr"`
	BoardPubKey    cipher.PubKey `bbs:"bpk"`
	UserPubKeyStr  string        `bbs:"upkStr"`
	UserPubKey     cipher.PubKey `bbs:"upk"`
	ModeStr        string        `bbs:"modeStr"`
	Mode           int8          `bbs:"mode"`
	TagStr         string        `bbs:"tagStr"`
	Tag            []byte        `bbs:"tag"`
	Vote           *Vote
}

func (a *UserVoteIO) Process(upk cipher.PubKey, usk cipher.SecKey) error {
	if e := tag.Process(a); e != nil {
		return e
	}
	a.Vote = &Vote{
		OfBoard: a.BoardPubKey,
		OfUser:  a.UserPubKey,
		Mode:    a.Mode,
		Tag:     a.Tag,
		Created: time.Now().UnixNano(),
		Creator: upk,
	}
	tag.Sign(a.Vote, upk, usk)
	return nil
}

type ThreadVoteIO struct {
	BoardPubKeyStr string        `bbs:"bpkStr"`
	BoardPubKey    cipher.PubKey `bbs:"bpk"`
	ThreadRefStr   string        `bbs:"tRefStr"`
	ThreadRef      cipher.SHA256 `bbs:"tRef"`
	ModeStr        string        `bbs:"modeStr"`
	Mode           int8          `bbs:"mode"`
	TagStr         string        `bbs:"tagStr"`
	Tag            []byte        `bbs:"tag"`
	Vote           *Vote
}

func (a *ThreadVoteIO) Process(upk cipher.PubKey, usk cipher.SecKey) error {
	if e := tag.Process(a); e != nil {
		return e
	}
	a.Vote = &Vote{
		OfBoard:  a.BoardPubKey,
		OfThread: a.ThreadRef,
		Mode:     a.Mode,
		Tag:      a.Tag,
		Created:  time.Now().UnixNano(),
		Creator:  upk,
	}
	tag.Sign(a.Vote, upk, usk)
	return nil
}

type PostVoteIO struct {
	BoardPubKeyStr string        `bbs:"bpkStr"`
	BoardPubKey    cipher.PubKey `bbs:"bpk"`
	PostRefStr     string        `bbs:"pRefStr"`
	PostRef        cipher.SHA256 `bbs:"pRef"`
	ModeStr        string        `bbs:"modeStr"`
	Mode           int8          `bbs:"mode"`
	TagStr         string        `bbs:"tagStr"`
	Tag            []byte        `bbs:"tag"`
	Vote           *Vote
}

func (a *PostVoteIO) Process(upk cipher.PubKey, usk cipher.SecKey) error {
	if e := tag.Process(a); e != nil {
		return e
	}
	a.Vote = &Vote{
		OfBoard: a.BoardPubKey,
		OfPost:  a.PostRef,
		Mode:    a.Mode,
		Tag:     a.Tag,
		Created: time.Now().UnixNano(),
		Creator: upk,
	}
	tag.Sign(a.Vote, upk, usk)
	return nil
}
