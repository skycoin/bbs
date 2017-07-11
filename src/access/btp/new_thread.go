package btp

import (
	"fmt"
	"github.com/skycoin/bbs/src/boo"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/bbs/src/store/view"
)

// NewThreadInput is the configuration struct used when creating a new thread.
type NewThreadInput struct {
	BoardPublicKey string `json:"board_public_key"`
	Name           string `json:"name"`
	Desc           string `json:"description"`
}

// NewThread creates a new thread on specified board.
func (a *BoardAccessor) NewThread(in *NewThreadInput) (*view.Thread, error) {
	defer a.lock()()

	// Obtain board public key.
	bpk, e := misc.GetPubKey(in.BoardPublicKey)
	if e != nil {
		return nil, boo.New(boo.InvalidInput,
			"invalid public key provided:", e.Error())
	}

	// Check if master board.
	bInfo, has := a.bFile.MasterBoards[bpk.Hex()]
	if !has {
		// TODO: Remote submission.
		return nil, boo.Newf(boo.ObjectNotFound,
			"you do not own a board of public key '%s'", bpk.Hex())
	}

	// Obtain board secret key.
	bsk, e := misc.GetSecKey(bInfo.SecretKey)
	if e != nil {
		return nil, boo.New(boo.Internal,
			"invalid secret key retrieved:", e.Error())
	}

	fmt.Println(bsk.Hex())

	return nil, nil
}
