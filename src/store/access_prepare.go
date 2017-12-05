package store

import (
	"encoding/json"
	"github.com/skycoin/bbs/src/misc/tag"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/skycoin/src/cipher"
	"time"
)

type PrepareOut struct {
	Hash string `json:"hash"`
	Raw  string `json:"raw"`
}

type PrepareThreadIn struct {
	OfBoardStr    string
	Name          string
	Body          string
	CreatorStr    string
	CreatorPubKey cipher.PubKey
	Data          *object.Body
}

func (a *PrepareThreadIn) Process() error {
	var e error
	if _, e = tag.GetPubKey(a.OfBoardStr); e != nil {
		return ErrProcess(e, "board public key")
	}
	if e = tag.CheckName(a.Name); e != nil {
		return ErrProcess(e, "name")
	}
	if e = tag.CheckBody(a.Body); e != nil {
		return ErrProcess(e, "body")
	}
	if a.CreatorPubKey, e = tag.GetPubKey(a.CreatorStr); e != nil {
		return ErrProcess(e, "creator public key")
	}
	a.Data = &object.Body{
		Type:    object.V5ThreadType,
		TS:      time.Now().UnixNano(),
		OfBoard: a.OfBoardStr,
		Name:    a.Name,
		Body:    a.Body,
		Creator: a.CreatorStr,
	}
	return nil
}

type PreparePostIn struct {
	OfBoardStr    string
	OfThreadStr   string
	OfPostStr     string
	Name          string
	Body          string
	ImagesStr     string
	CreatorStr    string
	CreatorPubKey cipher.PubKey
	Data          *object.Body
}

func (a *PreparePostIn) Process() error {
	var e error
	if _, e = tag.GetPubKey(a.OfBoardStr); e != nil {
		return ErrProcess(e, "board public key")
	}
	if _, e = tag.GetHash(a.OfThreadStr); e != nil {
		return ErrProcess(e, "thread hash")
	}
	if a.OfPostStr != "" {
		if _, e := tag.GetHash(a.OfPostStr); e != nil {
			return ErrProcess(e, "post hash")
		}
	}
	if e = tag.CheckName(a.Name); e != nil {
		return ErrProcess(e, "name")
	}
	if e = tag.CheckBody(a.Body); e != nil {
		return ErrProcess(e, "body")
	}
	var images []*object.ImageData
	if a.ImagesStr != "" {
		if e := json.Unmarshal([]byte(a.ImagesStr), &images); e != nil {
			return ErrProcess(e, "post images")
		}
	}
	if a.CreatorPubKey, e = tag.GetPubKey(a.CreatorStr); e != nil {
		return ErrProcess(e, "creator's public key")
	}
	a.Data = &object.Body{
		Type:     object.V5PostType,
		TS:       time.Now().UnixNano(),
		OfBoard:  a.OfBoardStr,
		OfThread: a.OfThreadStr,
		OfPost:   a.OfPostStr,
		Name:     a.Name,
		Body:     a.Body,
		Images:   images,
		Creator:  a.CreatorStr,
	}
	return nil
}

type PrepareThreadVoteIn struct {
	OfBoardStr    string
	OfThreadStr   string
	ValueStr      string
	TagsStr       string
	CreatorStr    string
	CreatorPubKey cipher.PubKey
	Data          *object.Body
}

func (a *PrepareThreadVoteIn) Process() error {
	var e error
	if _, e = tag.GetPubKey(a.OfBoardStr); e != nil {
		return ErrProcess(e, "board public key")
	}
	if _, e = tag.GetHash(a.OfThreadStr); e != nil {
		return ErrProcess(e, "thread hash")
	}
	var value int8
	if value, e = tag.GetVoteValue(a.ValueStr); e != nil {
		return ErrProcess(e, "vote value")
	}
	var tags []string
	if tags, e = tag.GetTags(a.TagsStr); e != nil {
		return ErrProcess(e, "vote tags")
	}
	if a.CreatorPubKey, e = tag.GetPubKey(a.CreatorStr); e != nil {
		return ErrProcess(e, "creator's public key")
	}
	a.Data = &object.Body{
		Type:     object.V5ThreadVoteType,
		TS:       time.Now().UnixNano(),
		OfBoard:  a.OfBoardStr,
		OfThread: a.OfThreadStr,
		Value:    int(value),
		Tags:     tags,
		Creator:  a.CreatorStr,
	}
	return nil
}

type PreparePostVoteIn struct {
	OfBoardStr    string
	OfPostStr     string
	ValueStr      string
	TagsStr       string
	CreatorStr    string
	CreatorPubKey cipher.PubKey
	Data          *object.Body
}

func (a *PreparePostVoteIn) Process() error {
	var e error
	if _, e = tag.GetPubKey(a.OfBoardStr); e != nil {
		return ErrProcess(e, "board public key")
	}
	if _, e = tag.GetHash(a.OfPostStr); e != nil {
		return ErrProcess(e, "post hash")
	}
	var value int8
	if value, e = tag.GetVoteValue(a.ValueStr); e != nil {
		return ErrProcess(e, "vote value")
	}
	var tags []string
	if tags, e = tag.GetTags(a.TagsStr); e != nil {
		return ErrProcess(e, "vote tags")
	}
	if a.CreatorPubKey, e = tag.GetPubKey(a.CreatorStr); e != nil {
		return ErrProcess(e, "creator's public key")
	}
	a.Data = &object.Body{
		Type:    object.V5PostVoteType,
		TS:      time.Now().UnixNano(),
		OfBoard: a.OfBoardStr,
		OfPost:  a.OfPostStr,
		Value:   int(value),
		Tags:    tags,
		Creator: a.CreatorStr,
	}
	return nil
}

type PrepareUserVoteIn struct {
	OfBoardStr    string
	OfUserStr     string
	ValueStr      string
	TagsStr       string
	CreatorStr    string
	CreatorPubKey cipher.PubKey
	Data          *object.Body
}

func (a *PrepareUserVoteIn) Process() error {
	var e error
	if _, e = tag.GetPubKey(a.OfBoardStr); e != nil {
		return ErrProcess(e, "board public key")
	}
	if _, e = tag.GetPubKey(a.OfUserStr); e != nil {
		return ErrProcess(e, "user public key")
	}
	var value int8
	if value, e = tag.GetVoteValue(a.ValueStr); e != nil {
		return ErrProcess(e, "vote value")
	}
	var tags []string
	if tags, e = tag.GetTags(a.TagsStr); e != nil {
		return ErrProcess(e, "vote tags")
	}
	if a.CreatorPubKey, e = tag.GetPubKey(a.CreatorStr); e != nil {
		return ErrProcess(e, "creator's public key")
	}
	a.Data = &object.Body{
		Type:    object.V5UserVoteType,
		TS:      time.Now().UnixNano(),
		OfBoard: a.OfBoardStr,
		OfUser:  a.OfUserStr,
		Value:   int(value),
		Tags:    tags,
		Creator: a.CreatorStr,
	}
	return nil
}

type FinalizeSubmissionIn struct {
	HashStr string
	Hash    cipher.SHA256
	SigStr  string
	Sig     cipher.Sig
}

func (a *FinalizeSubmissionIn) Process() error {
	var e error
	if a.Hash, e = tag.GetHash(a.HashStr); e != nil {
		return ErrProcess(e, "hash")
	}
	if a.Sig, e = tag.GetSig(a.SigStr); e != nil {
		return ErrProcess(e, "signature")
	}
	return nil
}
