package object

import (
	"encoding/json"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/keys"
	"github.com/skycoin/bbs/src/misc/tag"
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
	Content     *Content
}

func (a *NewBoardIO) Process(subPKs []cipher.PubKey) error {
	log.Println("Processing board, got submissions pks:", subPKs)
	if e := tag.Process(a); e != nil {
		return e
	}
	a.Content = new(Content)
	a.Content.SetHeader(&ContentHeaderData{})
	a.Content.SetBody(&Body{
		Type:    V5BoardType,
		TS:      time.Now().UnixNano(),
		Name:    a.Name,
		Body:    a.Body,
		SubKeys: keys.PubKeyArrayToStringArray(subPKs),
		Tags:    []string{},
	})
	return nil
}

// NewThread represents io required to create a new thread.
type NewThreadIO struct {
	BoardPubKeyStr   string        `bbs:"bpkStr"`
	BoardPubKey      cipher.PubKey `bbs:"bpk"`
	Name             string        `bbs:"name"`
	Body             string        `bbs:"body"`
	CreatorSecKeyStr string        `bbs:"uskStr"`
	CreatorSecKey    cipher.SecKey `bbs:"usk"`
	CreatorPubKey    cipher.PubKey
	Transport        *Transport
}

func (a *NewThreadIO) Process() error {
	if e := tag.Process(a); e != nil {
		return e
	}
	a.CreatorPubKey = cipher.PubKeyFromSecKey(a.CreatorSecKey)
	tData := &Body{
		Type: V5ThreadType,
		TS: time.Now().UnixNano(),
		OfBoard: a.BoardPubKey.Hex(),
		Name:    a.Name,
		Body:    a.Body,
		Creator: a.CreatorPubKey.Hex(),
	}
	tDataRaw, e := json.Marshal(tData)
	if e != nil {
		return e
	}
	tSig := cipher.SignHash(cipher.SumSHA256(tDataRaw), a.CreatorSecKey)
	if a.Transport, e = NewTransport(tDataRaw, tSig); e != nil {
		return e
	}
	return nil
}

type NewPostIO struct {
	BoardPubKeyStr   string        `bbs:"bpkStr"`
	BoardPubKey      cipher.PubKey `bbs:"bpk"`
	ThreadRefStr     string        `bbs:"tRefStr"`
	ThreadRef        cipher.SHA256 `bbs:"tRef"`
	PostRefStr       string        `bbs:"pRefStr"`
	PostRef          cipher.SHA256 `bbs:"pRef"`
	Name             string        `bbs:"name"`
	Body             string        `bbs:"body"`
	ImagesStr        string
	Images           []*ImageData
	CreatorSecKeyStr string        `bbs:"uskStr"`
	CreatorSecKey    cipher.SecKey `bbs:"usk"`
	CreatorPubKey    cipher.PubKey
	Transport        *Transport
}

func (a *NewPostIO) Process() error {
	if e := tag.Process(a); e != nil {
		return e
	}
	if a.ImagesStr != "" {
		if e := json.Unmarshal([]byte(a.ImagesStr), &a.Images); e != nil {
			return boo.WrapType(e, boo.InvalidInput, "failed to read 'images' form value")
		}
	}
	a.CreatorPubKey = cipher.PubKeyFromSecKey(a.CreatorSecKey)
	pData := &Body{
		Type: V5PostType,
		TS: time.Now().UnixNano(),
		OfBoard:  a.BoardPubKey.Hex(),
		OfThread: a.ThreadRef.Hex(),
		OfPost:   a.PostRef.Hex(),
		Name:     a.Name,
		Body:     a.Body,
		Images:   a.Images,
		Creator:  a.CreatorPubKey.Hex(),
	}
	pDataRaw, e := json.Marshal(pData)
	if e != nil {
		return e
	}
	pSig := cipher.SignHash(cipher.SumSHA256(pDataRaw), a.CreatorSecKey)
	if a.Transport, e = NewTransport(pDataRaw, pSig); e != nil {
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
	PubKeyStr     string        `bbs:"bpkStr"`
	PubKey        cipher.PubKey `bbs:"bpk"`
	UserPubKeyStr string        `bbs:"upkStr"`
	UserPubKey    cipher.PubKey `bbs:"upk"`
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
	UserPubKeyStr  string        `bbs:"upkStr"`
	UserPubKey     cipher.PubKey `bbs:"upk"`
}

func (a *ThreadIO) Process() error {
	if e := tag.Process(a); e != nil {
		return e
	}
	return nil
}

type UserVoteIO struct {
	BoardPubKeyStr   string        `bbs:"bpkStr"`
	BoardPubKey      cipher.PubKey `bbs:"bpk"`
	UserPubKeyStr    string        `bbs:"upkStr"`
	UserPubKey       cipher.PubKey `bbs:"upk"`
	ModeStr          string        `bbs:"modeStr"`
	Mode             int8          `bbs:"mode"`
	TagStr           string        `bbs:"tagStr"`
	Tag              []byte        `bbs:"tag"`
	CreatorSecKeyStr string        `bbs:"uskStr"`
	CreatorSecKey    cipher.SecKey `bbs:"usk"`
	CreatorPubKey    cipher.PubKey
	Transport        *Transport
}

func (a *UserVoteIO) Process() error {
	if e := tag.Process(a); e != nil {
		return e
	}
	a.CreatorPubKey = cipher.PubKeyFromSecKey(a.CreatorSecKey)
	vData := &Body{
		Type: V5UserVoteType,
		TS: time.Now().UnixNano(),
		OfBoard: a.BoardPubKeyStr,
		OfUser:  a.UserPubKeyStr,
		Value:   int(a.Mode),
		Tag:     string(a.Tag),
		Creator: a.CreatorPubKey.Hex(),

	}
	vDataRaw, e := json.Marshal(vData)
	if e != nil {
		return e
	}
	vSig := cipher.SignHash(cipher.SumSHA256(vDataRaw), a.CreatorSecKey)

	if a.Transport, e = NewTransport(vDataRaw, vSig); e != nil {
		return e
	}

	return nil
}

type ThreadVoteIO struct {
	BoardPubKeyStr   string        `bbs:"bpkStr"`
	BoardPubKey      cipher.PubKey `bbs:"bpk"`
	ThreadRefStr     string        `bbs:"tRefStr"`
	ThreadRef        cipher.SHA256 `bbs:"tRef"`
	ModeStr          string        `bbs:"modeStr"`
	Mode             int8          `bbs:"mode"`
	TagStr           string        `bbs:"tagStr"`
	Tag              []byte        `bbs:"tag"`
	CreatorSecKeyStr string        `bbs:"uskStr"`
	CreatorSecKey    cipher.SecKey `bbs:"usk"`
	CreatorPubKey    cipher.PubKey
	Transport        *Transport
}

func (a *ThreadVoteIO) Process() error {
	if e := tag.Process(a); e != nil {
		return e
	}
	a.CreatorPubKey = cipher.PubKeyFromSecKey(a.CreatorSecKey)
	vData := &Body{
		Type: V5ThreadVoteType,
		TS: time.Now().UnixNano(),
		OfBoard: a.BoardPubKeyStr,
		OfThread: a.ThreadRefStr,
		Value:   int(a.Mode),
		Tag:     string(a.Tag),
		Creator: a.CreatorPubKey.Hex(),
	}
	vDataRaw, e := json.Marshal(vData)
	if e != nil {
		return e
	}
	vSig := cipher.SignHash(cipher.SumSHA256(vDataRaw), a.CreatorSecKey)
	if a.Transport, e = NewTransport(vDataRaw, vSig); e != nil {
		return e
	}
	return nil
}

type PostVoteIO struct {
	BoardPubKeyStr   string        `bbs:"bpkStr"`
	BoardPubKey      cipher.PubKey `bbs:"bpk"`
	PostRefStr       string        `bbs:"pRefStr"`
	PostRef          cipher.SHA256 `bbs:"pRef"`
	ModeStr          string        `bbs:"modeStr"`
	Mode             int8          `bbs:"mode"`
	TagStr           string        `bbs:"tagStr"`
	Tag              []byte        `bbs:"tag"`
	CreatorSecKeyStr string        `bbs:"uskStr"`
	CreatorSecKey    cipher.SecKey `bbs:"usk"`
	CreatorPubKey    cipher.PubKey
	Transport        *Transport
}

func (a *PostVoteIO) Process() error {
	if e := tag.Process(a); e != nil {
		return e
	}
	a.CreatorPubKey = cipher.PubKeyFromSecKey(a.CreatorSecKey)
	vData := &Body{
		Type: V5PostVoteType,
		TS: time.Now().UnixNano(),
		OfBoard: a.BoardPubKeyStr,
		OfPost: a.PostRefStr,
		Value:   int(a.Mode),
		Tag:     string(a.Tag),
		Creator: a.CreatorPubKey.Hex(),
	}
	vDataRaw, e := json.Marshal(vData)
	if e != nil {
		return e
	}
	vSig := cipher.SignHash(cipher.SumSHA256(vDataRaw), a.CreatorSecKey)
	if a.Transport, e = NewTransport(vDataRaw, vSig); e != nil {
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
