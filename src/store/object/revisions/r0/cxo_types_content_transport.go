package r0

import (
	"encoding/json"
	"fmt"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/skycoin/src/cipher"
)

type ContentTransport struct {
	raw    []byte
	header *ContentHeaderData
}

type ThreadTransport struct {
	*ContentTransport
	body *ThreadData
}

type CheckThreadFunc func(td *ThreadData) error

func NewThreadTransport(raw []byte, sig cipher.Sig, check CheckThreadFunc) (*ThreadTransport, error) {
	body := new(ThreadData)
	if e := json.Unmarshal(raw, body); e != nil {
		return nil, genErrInvalidJSON(e, "thread")
	}
	if check != nil {
		if e := check(body); e != nil {
			return nil, e
		}
	}

	header := &ContentHeaderData{
		Type: V5ThreadType,
		Hash: cipher.SumSHA256(raw).Hex(),
		PK:   body.GetCreator().Hex(),
		Sig:  sig.Hex(),
	}
	if e := header.Verify(); e != nil {
		return nil, e
	}

	return &ThreadTransport{
		ContentTransport: &ContentTransport{raw: raw, header: header},
		body:             body,
	}, nil
}

type PostTransport struct {
	*ContentTransport
	body *PostData
}

type CheckPostFunc func(pd *PostData) error

func NewPostTransport(raw []byte, sig cipher.Sig, check CheckPostFunc) (*PostTransport, error) {
	body := new(PostData)
	if e := json.Unmarshal(raw, body); e != nil {
		return nil, genErrInvalidJSON(e, "post")
	}
	if check != nil {
		if e := check(body); e != nil {
			return nil, e
		}
	}

	header := &ContentHeaderData{
		Type: V5PostType,
		Hash: cipher.SumSHA256(raw).Hex(),
		PK:   body.GetCreator().Hex(),
		Sig:  sig.Hex(),
	}
	if e := header.Verify(); e != nil {
		return nil, e
	}

	return &PostTransport{
		ContentTransport: &ContentTransport{raw: raw, header: header},
		body:             body,
	}, nil
}

type ThreadVoteTransport struct {
	*ContentTransport
	body *ThreadVoteData
}

type CheckThreadVoteFunc func(tvd *ThreadVoteData) error

func NewThreadVoteTransport(raw []byte, sig cipher.Sig, check CheckThreadVoteFunc) (*ThreadVoteTransport, error) {
	body := new(ThreadVoteData)
	if e := json.Unmarshal(raw, body); e != nil {
		return nil, genErrInvalidJSON(e, "thread vote")
	}
	if check != nil {
		if e := check(body); e != nil {
			return nil, e
		}
	}

	header := &ContentHeaderData{
		Type: V5ThreadVoteType,
		Hash: cipher.SumSHA256(raw).Hex(),
		PK:   body.GetCreator().Hex(),
		Sig:  sig.Hex(),
	}
	if e := header.Verify(); e != nil {
		return nil, e
	}

	return &ThreadVoteTransport{
		ContentTransport: &ContentTransport{raw: raw, header: header},
		body:             body,
	}, nil
}

type PostVoteTransport struct {
	*ContentTransport
	body *PostVoteData
}

type CheckPostVoteFunc func(pvd *PostVoteData) error

func NewPostVoteTransport(raw []byte, sig cipher.Sig, check CheckPostVoteFunc) (*PostVoteTransport, error) {
	body := new(PostVoteData)
	if e := json.Unmarshal(raw, body); e != nil {
		return nil, genErrInvalidJSON(e, "post vote")
	}
	if check != nil {
		if e := check(body); e != nil {
			return nil, e
		}
	}

	header := &ContentHeaderData{
		Type: V5PostVoteType,
		Hash: cipher.SumSHA256(raw).Hex(),
		PK:   body.GetCreator().Hex(),
		Sig:  sig.Hex(),
	}
	if e := header.Verify(); e != nil {
		return nil, e
	}

	return &PostVoteTransport{
		ContentTransport: &ContentTransport{raw: raw, header: header},
		body:             body,
	}, nil
}

type UserVoteTransport struct {
	*ContentTransport
	body *UserVoteData
}

type CheckUserVoteFunc func(uvd *UserVoteData) error

func NewUserVoteTransport(raw []byte, sig cipher.Sig, check CheckUserVoteFunc) (*UserVoteTransport, error) {
	body := new(UserVoteData)
	if e := json.Unmarshal(raw, body); e != nil {
		return nil, genErrInvalidJSON(e, "user vote")
	}
	if check != nil {
		if e := check(body); e != nil {
			return nil, e
		}
	}

	header := &ContentHeaderData{
		Type: V5UserVoteType,
		Hash: cipher.SumSHA256(raw).Hex(),
		PK:   body.GetCreator().Hex(),
		Sig:  sig.Hex(),
	}
	if e := header.Verify(); e != nil {
		return nil, e
	}

	return &UserVoteTransport{
		ContentTransport: &ContentTransport{raw: raw, header: header},
		body:             body,
	}, nil
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func genErrInvalidJSON(e error, what string) error {
	return boo.WrapType(e, boo.InvalidInput,
		fmt.Sprintf("failed to read '%s' data", what))
}

func genErrHeaderUnverified(e error, what string) error {
	return boo.WrapType(e, boo.NotAuthorised,
		fmt.Sprintf("failed to verify '%s'", what))
}
