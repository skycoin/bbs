package r0

import (
	"fmt"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/skycoin/src/cipher"
	"encoding/json"
)

type Transport struct {
	Header  *ContentHeaderData
	Body    *Body
	Content *Content
}

func NewTransport(rawBody []byte, sig cipher.Sig) (*Transport, error) {
	out := new(Transport)

	var e error
	if out.Body, e = NewBody(rawBody); e != nil {
		return nil, e
	}

	creator, e := out.Body.GetCreator()
	if e != nil {
		return nil, e
	}

	out.Header = &ContentHeaderData{
		Hash: cipher.SumSHA256(rawBody).Hex(),
		Sig: sig.Hex(),
	}
	if e := out.Header.Verify(creator); e != nil {
		return nil, e
	}

	out.Content = new(Content)
	if out.Content.Header, e = json.Marshal(out.Header); e != nil {
		return nil, e
	}
	out.Content.Body = rawBody

	return out, nil
}

func (t *Transport) GetOfBoard() cipher.PubKey {
	pk, _ := t.Body.GetOfBoard()
	return pk
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func genErrInvalidJSON(e error, what string) error {
	return boo.WrapType(e, boo.InvalidInput,
		fmt.Sprintf("failed to read '%s' data", what))
}

func genErrHeaderUnverified(e error, hash string) error {
	return boo.WrapType(e, boo.NotAuthorised,
		fmt.Sprintf("failed to verify content of hash '%s'", hash))
}
