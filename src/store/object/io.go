package object

import (
	"encoding/json"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/keys"
	"github.com/skycoin/bbs/src/misc/tag"
	"github.com/skycoin/bbs/src/store/object/revisions/r0"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"time"
)

// NewBoard represents io required to create a new board.
type NewBoardIO struct {
	Name        string        `bbs:"name"`
	Body        string        `bbs:"body"`
	Seed        string        `bbs:"bSeed"`
	BoardPubKey cipher.PubKey `bbs:"bpk"`
	BoardSecKey cipher.SecKey `bbs:"bsk"`
	Board       *r0.Board
}

func (a *NewBoardIO) Process(subPKs []cipher.PubKey) error {
	log.Println("Processing board, got submissions pks:", subPKs)
	if e := tag.Process(a); e != nil {
		return e
	}
	a.Board = new(r0.Board)
	a.Board.SetData(&r0.BoardData{
		Name:    a.Name,
		Body:    a.Body,
		Created: time.Now().UnixNano(),
		SubKeys: keys.PubKeyArrayToStringArray(subPKs),
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
	a.Thread = new(r0.Thread)
	a.Thread.SetData(&r0.ThreadData{
		OfBoard: a.BoardPubKey.Hex(),
		Name:    a.Name,
		Body:    a.Body,
		Created: time.Now().UnixNano(),
		Creator: upk.Hex(),
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
	ImagesStr      string
	Images         []*r0.ImageData
	Post           *r0.Post
}

func (a *NewPostIO) Process(upk cipher.PubKey, usk cipher.SecKey) error {
	if e := tag.Process(a); e != nil {
		return e
	}
	if a.ImagesStr != "" {
		if e := json.Unmarshal([]byte(a.ImagesStr), &a.Images); e != nil {
			return boo.WrapType(e, boo.InvalidInput, "failed to read 'images' form value")
		}
	}
	a.Post = new(r0.Post)
	a.Post.SetData(&r0.PostData{
		OfBoard:  a.BoardPubKey.Hex(),
		OfThread: a.ThreadRef.Hex(),
		OfPost:   a.PostRef.Hex(),
		Name:     a.Name,
		Body:     a.Body,
		Images:   a.Images,
		Created:  time.Now().UnixNano(),
		Creator:  upk.Hex(),
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

type ExportBoardIO struct {
	PubKeyStr string        `bbs:"bpkStr"`
	PubKey    cipher.PubKey `bbs:"bpk"`
	Name      string        `bbs:"alias"`
}

func (a *ExportBoardIO) Process() error {
	if e := tag.Process(a); e != nil {
		return e
	}
	return nil
}
