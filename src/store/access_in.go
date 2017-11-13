package store

import (
	"encoding/json"
	"fmt"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/tag"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/skycoin/src/cipher"
	"time"
)

type SubmissionIn struct {
	Body   []byte
	SigStr string
	Sig    cipher.Sig
}

func (a *SubmissionIn) Process() error {
	// Convert signature.
	var e error
	if a.Sig, e = cipher.SigFromHex(a.SigStr); e != nil {
		return boo.WrapType(e, boo.InvalidInput,
			"invalid hex representation of signature:", a.SigStr)
	}
	return nil
}

// ConnectionIn represents the input required when connection/disconnecting from address.
type ConnectionIn struct {
	Address string
}

func (a *ConnectionIn) Process() error {
	if e := tag.CheckAddress(a.Address); e != nil {
		return e
	}
	return nil
}

type UserIn struct {
	BoardPubKeyStr string
	BoardPubKey    cipher.PubKey
	UserPubKeyStr  string
	UserPubKey     cipher.PubKey
}

func (a *UserIn) Process() error {
	var e error
	if a.BoardPubKey, e = tag.GetPubKey(a.BoardPubKeyStr); e != nil {
		return ErrProcess(e, "board public key")
	}
	if a.UserPubKey, e = tag.GetPubKey(a.UserPubKeyStr); e != nil {
		return ErrProcess(e, "user's public key")
	}
	return nil
}

// BoardIn represents a subscription input.
type BoardIn struct {
	PubKeyStr     string
	PubKey        cipher.PubKey
	UserPubKeyStr string
	UserPubKey    cipher.PubKey
}

func (a *BoardIn) Process() error {
	var e error
	if a.PubKey, e = tag.GetPubKey(a.PubKeyStr); e != nil {
		return ErrProcess(e, "board public key")
	}
	if a.UserPubKeyStr != "" {
		if a.UserPubKey, e = tag.GetPubKey(a.UserPubKeyStr); e != nil {
			return ErrProcess(e, "user public key")
		}
	}
	return nil
}

type ExportBoardIn struct {
	FilePath  string
	PubKeyStr string
	PubKey    cipher.PubKey
}

func (a *ExportBoardIn) Process() error {
	var e error
	if a.PubKey, e = tag.GetPubKey(a.PubKeyStr); e != nil {
		return ErrProcess(e, "board public key")
	}
	return nil
}

type ImportBoardIn struct {
	FilePath  string
}

func (a *ImportBoardIn) Process() error {
	if e := tag.CheckPath(a.FilePath); e != nil {
		return ErrProcess(e, "file path")
	}
	return nil
}

type NewBoardIn struct {
	Name        string
	Body        string
	Seed        string
	BoardPubKey cipher.PubKey
	BoardSecKey cipher.SecKey
	TS          int64
	Content     *object.Content
}

func (a *NewBoardIn) Process(subKeyTrans []*object.MessengerSubKeyTransport) error {
	if e := tag.CheckName(a.Name); e != nil {
		return ErrProcess(e, "name")
	}
	if e := tag.CheckBody(a.Body); e != nil {
		return ErrProcess(e, "body")
	}
	if a.TS == 0 {
		a.TS = time.Now().UnixNano()
	}

	a.BoardPubKey, a.BoardSecKey = cipher.GenerateDeterministicKeyPair([]byte(a.Seed))

	subKeys := make([]object.MessengerSubKey, len(subKeyTrans))
	for i, subKey := range subKeyTrans {
		subKeys[i] = subKey.ToMessengerSubKey()
	}

	a.Content = new(object.Content)
	a.Content.SetHeader(&object.ContentHeaderData{})
	a.Content.SetBody(&object.Body{
		Type:    object.V5BoardType,
		TS:      a.TS,
		Name:    a.Name,
		Body:    a.Body,
		SubKeys: subKeys,
		Tags:    []string{},
	})
	return nil
}

type ThreadIn struct {
	BoardPubKeyStr string
	BoardPubKey    cipher.PubKey
	ThreadRefStr   string
	ThreadRef      cipher.SHA256
	UserPubKeyStr  string
	UserPubKey     cipher.PubKey
}

func (a *ThreadIn) Process() error {
	var e error
	if a.BoardPubKey, e = tag.GetPubKey(a.BoardPubKeyStr); e != nil {
		return ErrProcess(e, "board public key")
	}
	if a.ThreadRef, e = tag.GetHash(a.ThreadRefStr); e != nil {
		return ErrProcess(e, "thread hash")
	}
	if a.UserPubKeyStr != "" {
		if a.UserPubKey, e = tag.GetPubKey(a.UserPubKeyStr); e != nil {
			return ErrProcess(e, "user's public key")
		}
	}
	return nil
}

type NewThreadIn struct {
	BoardPubKeyStr   string
	BoardPubKey      cipher.PubKey
	Name             string
	Body             string
	CreatorSecKeyStr string
	CreatorSecKey    cipher.SecKey
	CreatorPubKey    cipher.PubKey
	TS               int64
	Transport        *object.Transport
}

func (a *NewThreadIn) Process() error {
	var e error
	if a.BoardPubKey, e = tag.GetPubKey(a.BoardPubKeyStr); e != nil {
		return ErrProcess(e, "board public key")
	}
	if e = tag.CheckName(a.Name); e != nil {
		return ErrProcess(e, "name")
	}
	if e = tag.CheckBody(a.Body); e != nil {
		return ErrProcess(e, "body")
	}
	if a.CreatorSecKey, e = tag.GetSecKey(a.CreatorSecKeyStr); e != nil {
		return ErrProcess(e, "creator's secret key")
	}
	a.CreatorPubKey = cipher.PubKeyFromSecKey(a.CreatorSecKey)
	if a.TS == 0 {
		a.TS = time.Now().UnixNano()
	}
	data := &object.Body{
		Type:    object.V5ThreadType,
		TS:      a.TS,
		OfBoard: a.BoardPubKey.Hex(),
		Name:    a.Name,
		Body:    a.Body,
		Creator: a.CreatorPubKey.Hex(),
	}
	raw, e := json.Marshal(data)
	if e != nil {
		return ErrProcess(e, "thread body")
	}
	sig := cipher.SignHash(cipher.SumSHA256(raw), a.CreatorSecKey)
	if a.Transport, e = object.NewTransport(raw, sig); e != nil {
		return ErrProcess(e, "thread")
	}
	return nil
}

type NewPostIn struct {
	BoardPubKeyStr   string
	BoardPubKey      cipher.PubKey
	ThreadRefStr     string
	ThreadRef        cipher.SHA256
	PostRefStr       string
	PostRef          cipher.SHA256
	Name             string
	Body             string
	ImagesStr        string
	Images           []*object.ImageData
	CreatorSecKeyStr string
	CreatorSecKey    cipher.SecKey
	CreatorPubKey    cipher.PubKey
	TS               int64
	Transport        *object.Transport
}

func (a *NewPostIn) Process() error {
	var e error
	if a.BoardPubKey, e = tag.GetPubKey(a.BoardPubKeyStr); e != nil {
		return ErrProcess(e, "board public key")
	}
	if a.ThreadRef, e = tag.GetHash(a.ThreadRefStr); e != nil {
		return ErrProcess(e, "thread hash")
	}
	if a.PostRefStr != "" {
		if a.PostRef, e = tag.GetHash(a.PostRefStr); e != nil {
			return ErrProcess(e, "post hash")
		}
	}
	if e = tag.CheckName(a.Name); e != nil {
		return ErrProcess(e, "name")
	}
	if e = tag.CheckBody(a.Body); e != nil {
		return ErrProcess(e, "body")
	}
	if a.ImagesStr != "" {
		if e = json.Unmarshal([]byte(a.ImagesStr), &a.Images); e != nil {
			return ErrProcess(e, "post images")
		}
	}
	if a.CreatorSecKey, e = tag.GetSecKey(a.CreatorSecKeyStr); e != nil {
		return ErrProcess(e, "creator's secret key")
	}
	a.CreatorPubKey = cipher.PubKeyFromSecKey(a.CreatorSecKey)
	if a.TS == 0 {
		a.TS = time.Now().UnixNano()
	}
	data := &object.Body{
		Type:     object.V5PostType,
		TS:       a.TS,
		OfBoard:  a.BoardPubKeyStr,
		OfThread: a.ThreadRefStr,
		OfPost:   a.PostRefStr,
		Name:     a.Name,
		Body:     a.Body,
		Images:   a.Images,
		Creator:  a.CreatorPubKey.Hex(),
	}
	raw, e := json.Marshal(data)
	if e != nil {
		return ErrProcess(e, "post body")
	}
	sig := cipher.SignHash(cipher.SumSHA256(raw), a.CreatorSecKey)
	if a.Transport, e = object.NewTransport(raw, sig); e != nil {
		return ErrProcess(e, "post")
	}
	return nil
}

type VoteThreadIn struct {
	BoardPubKeyStr   string
	BoardPubKey      cipher.PubKey
	ThreadRefStr     string
	ThreadRef        cipher.SHA256
	ValueStr         string
	Value            int8
	TagsStr          string
	Tags             []string
	CreatorSecKeyStr string
	CreatorSecKey    cipher.SecKey
	CreatorPubKey    cipher.PubKey
	TS               int64
	Transport        *object.Transport
}

func (a *VoteThreadIn) Process() error {
	var e error
	if a.BoardPubKey, e = tag.GetPubKey(a.BoardPubKeyStr); e != nil {
		return ErrProcess(e, "board public key")
	}
	if a.ThreadRef, e = tag.GetHash(a.ThreadRefStr); e != nil {
		return ErrProcess(e, "thread hash")
	}
	if a.Value, e = tag.GetVoteValue(a.ValueStr); e != nil {
		return ErrProcess(e, "vote value")
	}
	if a.Tags, e = tag.GetTags(a.TagsStr); e != nil {
		return ErrProcess(e, "tags")
	}
	if a.CreatorSecKey, e = tag.GetSecKey(a.CreatorSecKeyStr); e != nil {
		return ErrProcess(e, "creator's secret key")
	}
	a.CreatorPubKey = cipher.PubKeyFromSecKey(a.CreatorSecKey)
	if a.TS == 0 {
		a.TS = time.Now().UnixNano()
	}
	data := &object.Body{
		Type:     object.V5ThreadVoteType,
		TS:       a.TS,
		OfBoard:  a.BoardPubKeyStr,
		OfThread: a.ThreadRefStr,
		Value:    int(a.Value),
		Creator:  a.CreatorPubKey.Hex(),
	}
	raw, e := json.Marshal(data)
	if e != nil {
		return ErrProcess(e, "thread vote body")
	}
	sig := cipher.SignHash(cipher.SumSHA256(raw), a.CreatorSecKey)
	if a.Transport, e = object.NewTransport(raw, sig); e != nil {
		return ErrProcess(e, "thread vote")
	}
	return nil
}

type VoteUserIn struct {
	BoardPubKeyStr   string
	BoardPubKey      cipher.PubKey
	UserPubKeyStr    string
	UserPubKey       cipher.PubKey
	ValueStr         string
	Value            int8
	TagsStr          string
	Tags             []string
	CreatorSecKeyStr string
	CreatorSecKey    cipher.SecKey
	CreatorPubKey    cipher.PubKey
	TS               int64
	Transport        *object.Transport
}

func (a *VoteUserIn) Process() error {
	var e error
	if a.BoardPubKey, e = tag.GetPubKey(a.BoardPubKeyStr); e != nil {
		return ErrProcess(e, "board public key")
	}
	if a.UserPubKey, e = tag.GetPubKey(a.UserPubKeyStr); e != nil {
		return ErrProcess(e, "user public key")
	}
	if a.Value, e = tag.GetVoteValue(a.ValueStr); e != nil {
		return ErrProcess(e, "vote value")
	}
	if a.Tags, e = tag.GetTags(a.TagsStr); e != nil {
		return ErrProcess(e, "tags")
	}
	if a.CreatorSecKey, e = tag.GetSecKey(a.CreatorSecKeyStr); e != nil {
		return ErrProcess(e, "creator's secret key")
	}
	a.CreatorPubKey = cipher.PubKeyFromSecKey(a.CreatorSecKey)
	if a.TS == 0 {
		a.TS = time.Now().UnixNano()
	}
	data := &object.Body{
		Type:    object.V5UserVoteType,
		TS:      a.TS,
		OfBoard: a.BoardPubKeyStr,
		OfUser:  a.UserPubKeyStr,
		Value:   int(a.Value),
		Tags:    a.Tags,
		Creator: a.CreatorPubKey.Hex(),
	}
	raw, e := json.Marshal(data)
	if e != nil {
		return ErrProcess(e, "user vote body")
	}
	sig := cipher.SignHash(cipher.SumSHA256(raw), a.CreatorSecKey)
	if a.Transport, e = object.NewTransport(raw, sig); e != nil {
		return ErrProcess(e, "user vote")
	}
	return nil
}

type VotePostIn struct {
	BoardPubKeyStr   string
	BoardPubKey      cipher.PubKey
	PostRefStr       string
	PostRef          cipher.SHA256
	ValueStr         string
	Value            int8
	TagsStr          string
	Tags             []string
	CreatorSecKeyStr string
	CreatorSecKey    cipher.SecKey
	CreatorPubKey    cipher.PubKey
	TS               int64
	Transport        *object.Transport
}

func (a *VotePostIn) Process() error {
	var e error
	if a.BoardPubKey, e = tag.GetPubKey(a.BoardPubKeyStr); e != nil {
		return ErrProcess(e, "board public key")
	}
	if a.PostRef, e = tag.GetHash(a.PostRefStr); e != nil {
		return ErrProcess(e, "post hash")
	}
	if a.Value, e = tag.GetVoteValue(a.ValueStr); e != nil {
		return ErrProcess(e, "vote value")
	}
	if a.Tags, e = tag.GetTags(a.TagsStr); e != nil {
		return ErrProcess(e, "tags")
	}
	if a.CreatorSecKey, e = tag.GetSecKey(a.CreatorSecKeyStr); e != nil {
		return ErrProcess(e, "creator's secret key")
	}
	a.CreatorPubKey = cipher.PubKeyFromSecKey(a.CreatorSecKey)
	if a.TS == 0 {
		a.TS = time.Now().UnixNano()
	}
	data := &object.Body{
		Type:    object.V5PostVoteType,
		TS:      a.TS,
		OfBoard: a.BoardPubKeyStr,
		OfPost:  a.PostRefStr,
		Value:   int(a.Value),
		Tags:    a.Tags,
		Creator: a.CreatorPubKey.Hex(),
	}
	raw, e := json.Marshal(data)
	if e != nil {
		return ErrProcess(e, "vote body")
	}
	sig := cipher.SignHash(cipher.SumSHA256(raw), a.CreatorSecKey)
	if a.Transport, e = object.NewTransport(raw, sig); e != nil {
		return ErrProcess(e, "vote")
	}
	return nil
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func ErrProcess(e error, what string) error {
	msg := fmt.Sprintf("failed to process %s", what)
	if e == nil {
		return boo.New(boo.InvalidInput, msg)
	} else {
		return boo.WrapType(e, boo.InvalidInput, msg)
	}
}
