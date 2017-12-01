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
	BoardPubKeyStr   string
	Name             string
	Body             string
	CreatorPubKeyStr string
	CreatorPubKey    cipher.PubKey
	Data             *object.Body
}

func (a *PrepareThreadIn) Process() error {
	var e error
	if _, e = tag.GetPubKey(a.BoardPubKeyStr); e != nil {
		return ErrProcess(e, "board public key")
	}
	if e = tag.CheckName(a.Name); e != nil {
		return ErrProcess(e, "name")
	}
	if e = tag.CheckBody(a.Body); e != nil {
		return ErrProcess(e, "body")
	}
	if a.CreatorPubKey, e = tag.GetPubKey(a.CreatorPubKeyStr); e != nil {
		return ErrProcess(e, "creator public key")
	}
	a.Data = &object.Body{
		Type:    object.V5ThreadType,
		TS:      time.Now().UnixNano(),
		OfBoard: a.BoardPubKeyStr,
		Name:    a.Name,
		Body:    a.Body,
		Creator: a.CreatorPubKeyStr,
	}
	return nil
}

type PreparePostIn struct {
	BoardPubKeyStr   string
	ThreadHashStr    string
	PostHashStr      string
	Name             string
	Body             string
	ImagesStr        string
	CreatorPubKeyStr string
	CreatorPubKey    cipher.PubKey
	Data             *object.Body
}

func (a *PreparePostIn) Process() error {
	var e error
	if _, e = tag.GetPubKey(a.BoardPubKeyStr); e != nil {
		return ErrProcess(e, "board public key")
	}
	if _, e = tag.GetHash(a.ThreadHashStr); e != nil {
		return ErrProcess(e, "thread hash")
	}
	if a.PostHashStr != "" {
		if _, e := tag.GetHash(a.PostHashStr); e != nil {
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
	if a.CreatorPubKey, e = tag.GetPubKey(a.CreatorPubKeyStr); e != nil {
		return ErrProcess(e, "creator's public key")
	}
	a.Data = &object.Body{
		Type:     object.V5PostType,
		TS:       time.Now().UnixNano(),
		OfBoard:  a.BoardPubKeyStr,
		OfThread: a.ThreadHashStr,
		OfPost:   a.PostHashStr,
		Name:     a.Name,
		Body:     a.Body,
		Images:   images,
		Creator:  a.CreatorPubKeyStr,
	}
	return nil
}

type PrepareThreadVoteIn struct {
	BoardPubKeyStr   string
	ThreadHashStr    string
	ValueStr         string
	TagsStr          string
	CreatorPubKeyStr string
	CreatorPubKey    cipher.PubKey
	Data             *object.Body
}

func (a *PrepareThreadVoteIn) Process() error {
	var e error
	if _, e = tag.GetPubKey(a.BoardPubKeyStr); e != nil {
		return ErrProcess(e, "board public key")
	}
	if _, e = tag.GetHash(a.ThreadHashStr); e != nil {
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
	if a.CreatorPubKey, e = tag.GetPubKey(a.CreatorPubKeyStr); e != nil {
		return ErrProcess(e, "creator's public key")
	}
	a.Data = &object.Body{
		Type:     object.V5ThreadVoteType,
		TS:       time.Now().UnixNano(),
		OfBoard:  a.BoardPubKeyStr,
		OfThread: a.ThreadHashStr,
		Value:    int(value),
		Tags:     tags,
		Creator:  a.CreatorPubKeyStr,
	}
	return nil
}

type PreparePostVoteIn struct {
	BoardPubKeyStr   string
	PostHashStr      string
	ValueStr         string
	TagsStr          string
	CreatorPubKeyStr string
	CreatorPubKey    cipher.PubKey
	Data             *object.Body
}

func (a *PreparePostVoteIn) Process() error {
	var e error
	if _, e = tag.GetPubKey(a.BoardPubKeyStr); e != nil {
		return ErrProcess(e, "board public key")
	}
	if _, e = tag.GetHash(a.PostHashStr); e != nil {
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
	if a.CreatorPubKey, e = tag.GetPubKey(a.CreatorPubKeyStr); e != nil {
		return ErrProcess(e, "creator's public key")
	}
	a.Data = &object.Body{
		Type:    object.V5ThreadVoteType,
		TS:      time.Now().UnixNano(),
		OfBoard: a.BoardPubKeyStr,
		OfPost:  a.PostHashStr,
		Value:   int(value),
		Tags:    tags,
		Creator: a.CreatorPubKeyStr,
	}
	return nil
}

type PrepareUserVoteIn struct {
	BoardPubKeyStr   string
	UserPubKeyStr    string
	ValueStr         string
	TagsStr          string
	CreatorPubKeyStr string
	CreatorPubKey    cipher.PubKey
	Data             *object.Body
}

func (a *PrepareUserVoteIn) Process() error {
	var e error
	if _, e = tag.GetPubKey(a.BoardPubKeyStr); e != nil {
		return ErrProcess(e, "board public key")
	}
	if _, e = tag.GetPubKey(a.UserPubKeyStr); e != nil {
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
	if a.CreatorPubKey, e = tag.GetPubKey(a.CreatorPubKeyStr); e != nil {
		return ErrProcess(e, "creator's public key")
	}
	a.Data = &object.Body{
		Type:    object.V5ThreadVoteType,
		TS:      time.Now().UnixNano(),
		OfBoard: a.BoardPubKeyStr,
		OfUser:  a.UserPubKeyStr,
		Value:   int(value),
		Tags:    tags,
		Creator: a.CreatorPubKeyStr,
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
