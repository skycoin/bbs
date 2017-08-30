package object

import (
	"github.com/skycoin/bbs/src/misc/tag"
	"github.com/skycoin/skycoin/src/cipher"
	"time"
	"github.com/skycoin/bbs/src/store/object/revisions/r0"
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
	Board       *r0.Board
}

func (a *NewBoardIO) Process() error {
	if e := tag.Process(a); e != nil {
		return e
	}
	a.Board = &r0.Board{
		Created: time.Now().UnixNano(),
	}
	r0.SetData(a.Board, &r0.ContentData{
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
	Thread         *r0.Thread
}

func (a *NewThreadIO) Process(upk cipher.PubKey, usk cipher.SecKey) error {
	if e := tag.Process(a); e != nil {
		return e
	}
	a.Thread = &r0.Thread{
		OfBoard: a.BoardPubKey,
		Created: time.Now().UnixNano(),
		Creator: upk,
	}
	r0.SetData(a.Thread, &r0.ContentData{
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
	Post           *r0.Post
}

func (a *NewPostIO) Process(upk cipher.PubKey, usk cipher.SecKey) error {
	if e := tag.Process(a); e != nil {
		return e
	}
	a.Post = &r0.Post{
		OfBoard:  a.BoardPubKey,
		OfThread: a.ThreadRef,
		OfPost:   a.PostRef,
		Created:  time.Now().UnixNano(),
		Creator:  upk,
	}
	r0.SetData(a.Post, &r0.ContentData{
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
		User: r0.User{
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

// ConnectionIO represents input/output required when connection/disconnecting from address.
type ConnectionIO struct {
	Address string `bbs:"address"`
}

func (a *ConnectionIO) Process() error {
	if e := tag.Process(a); e != nil {
		return e
	}
	return nil
}

// SubmissionIO represents submission address input/output.
type SubmissionIO struct {
	BoardPubKeyStr string        `bbs:"bpkStr"`
	BoardPubKey    cipher.PubKey `bbs:"bpk"`
	SubAddress     string        `bbs:"address"`
}

func (a *SubmissionIO) Process() error {
	return tag.Process(a)
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
	Vote           *r0.Vote
}

func (a *UserVoteIO) Process(upk cipher.PubKey, usk cipher.SecKey) error {
	if e := tag.Process(a); e != nil {
		return e
	}
	a.Vote = &r0.Vote{
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
	Vote           *r0.Vote
}

func (a *ThreadVoteIO) Process(upk cipher.PubKey, usk cipher.SecKey) error {
	if e := tag.Process(a); e != nil {
		return e
	}
	a.Vote = &r0.Vote{
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
	Vote           *r0.Vote
}

func (a *PostVoteIO) Process(upk cipher.PubKey, usk cipher.SecKey) error {
	if e := tag.Process(a); e != nil {
		return e
	}
	a.Vote = &r0.Vote{
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

type UserIO struct {
	BoardPubKeyStr string        `bbs:"bpkStr"`
	BoardPubKey    cipher.PubKey `bbs:"bpk"`
	UserPubKeyStr  string        `bbs:"upkStr"`
	UserPubKey     cipher.PubKey `bbs:"upk"`
}

func (a *UserIO) Process() error {
	return tag.Process(a)
}
