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

type SubmissionIO struct {
	Body   []byte
	SigStr string
	Sig    cipher.Sig
}

func (a *SubmissionIO) Process() error {
	// Convert signature.
	var e error
	if a.Sig, e = cipher.SigFromHex(a.SigStr); e != nil {
		return boo.WrapType(e, boo.InvalidInput,
			"invalid hex representation of signature:", a.SigStr)
	}
	return nil
}

// NewBoard represents io required to create a new board.
type NewBoardIO struct {
	Name        string        `bbs:"name"`
	Body        string        `bbs:"body"`
	Seed        string        `bbs:"bSeed"`
	BoardPubKey cipher.PubKey `bbs:"bpk"`
	BoardSecKey cipher.SecKey `bbs:"bsk"`
	Content   *r0.Content
}

func (a *NewBoardIO) Process(subPKs []cipher.PubKey) error {
	log.Println("Processing board, got submissions pks:", subPKs)
	if e := tag.Process(a); e != nil {
		return e
	}
	a.Content = new(r0.Content)
	a.Content.SetHeader(&r0.ContentHeaderData{})
	a.Content.SetBody(&r0.Body{
		Type: r0.V5BoardType,
		TS: time.Now().UnixNano(),
		Name:    a.Name,
		Body:    a.Body,
		SubKeys: keys.PubKeyArrayToStringArray(subPKs),
		Tags:    []string{},
	})
	return nil
}

// NewThread represents io required to create a new thread.
type NewThreadIO struct {
	BoardPubKeyStr string        `bbs:"bpkStr"`
	BoardPubKey    cipher.PubKey `bbs:"bpk"`
	Name           string        `bbs:"name"`
	Body           string        `bbs:"body"`
	Transport      *r0.Transport
}

func (a *NewThreadIO) Process(upk cipher.PubKey, usk cipher.SecKey) error {
	if e := tag.Process(a); e != nil {
		return e
	}

	tData := &r0.ThreadData{
		OfBoard: a.BoardPubKey.Hex(),
		Name:    a.Name,
		Body:    a.Body,
		Creator: upk.Hex(),
	}
	tDataRaw, e := json.Marshal(tData)
	if e != nil {
		return e
	}
	tSig := cipher.SignHash(cipher.SumSHA256(tDataRaw), usk)

	if a.Transport, e = r0.NewTransport(tDataRaw, tSig); e != nil {
		return e
	}
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
	Transport      *r0.Transport
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

	pData := &r0.PostData{
		OfBoard:  a.BoardPubKey.Hex(),
		OfThread: a.ThreadRef.Hex(),
		OfPost:   a.PostRef.Hex(),
		Name:     a.Name,
		Body:     a.Body,
		Images:   a.Images,
		Creator:  upk.Hex(),
	}
	pDataRaw, e := json.Marshal(pData)
	if e != nil {
		return e
	}
	pSig := cipher.SignHash(cipher.SumSHA256(pDataRaw), usk)

	if a.Transport, e = r0.NewTransport(pDataRaw, pSig); e != nil {
		return e
	}

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
	Transport      *r0.Transport
}

func (a *UserVoteIO) Process(upk cipher.PubKey, usk cipher.SecKey) error {
	if e := tag.Process(a); e != nil {
		return e
	}

	vData := &r0.UserVoteData{
		VoteData: r0.VoteData{
			OfBoard: a.BoardPubKeyStr,
			Value:   int(a.Mode),
			Tag:     string(a.Tag),
			Creator: upk.Hex(),
		},
		OfUser: a.UserPubKeyStr,
	}
	vDataRaw, e := json.Marshal(vData)
	if e != nil {
		return e
	}
	vSig := cipher.SignHash(cipher.SumSHA256(vDataRaw), usk)

	if a.Transport, e = r0.NewTransport(vDataRaw, vSig); e != nil {
		return e
	}

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
	Transport      *r0.Transport
}

func (a *ThreadVoteIO) Process(upk cipher.PubKey, usk cipher.SecKey) error {
	if e := tag.Process(a); e != nil {
		return e
	}

	vData := &r0.ThreadVoteData{
		VoteData: r0.VoteData{
			OfBoard: a.BoardPubKeyStr,
			Value:   int(a.Mode),
			Tag:     string(a.Tag),
			Creator: upk.Hex(),
		},
		OfThread: a.ThreadRefStr,
	}
	vDataRaw, e := json.Marshal(vData)
	if e != nil {
		return e
	}
	vSig := cipher.SignHash(cipher.SumSHA256(vDataRaw), usk)

	if a.Transport, e = r0.NewTransport(vDataRaw, vSig); e != nil {
		return e
	}

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
	Transport      *r0.Transport
}

func (a *PostVoteIO) Process(upk cipher.PubKey, usk cipher.SecKey) error {
	if e := tag.Process(a); e != nil {
		return e
	}

	vData := &r0.PostVoteData{
		VoteData: r0.VoteData{
			OfBoard: a.BoardPubKeyStr,
			Value:   int(a.Mode),
			Tag:     string(a.Tag),
			Creator: upk.Hex(),
		},
		OfPost: a.PostRefStr,
	}
	vDataRaw, e := json.Marshal(vData)
	if e != nil {
		return e
	}
	vSig := cipher.SignHash(cipher.SumSHA256(vDataRaw), usk)

	if a.Transport, e = r0.NewTransport(vDataRaw, vSig); e != nil {
		return e
	}
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
