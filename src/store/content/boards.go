package content

import (
	"context"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/skycoin/src/cipher"
	"time"
)

// GetBoardResult get's the specified board of public key.
func GetBoardResult(_ context.Context, cxo *state.CXO, pk cipher.PubKey) (*Result, error) {
	result := NewResult(cxo, pk).
		getPages(true, false, false).
		getBoard()

	if e := result.Error(); e != nil {
		return nil, e
	}
	return result, nil
}

// NewBoard creates a new board and returns an error on failure.
func NewBoard(_ context.Context, cxo *state.CXO, in *object.NewBoardIO) error {
	e := cxo.NewRoot(in.GetPK(), in.GetSK(), func(r *node.Root) error {
		_, e := r.Append(
			r.MustDynamic("BoardPage", object.BoardPage{
				Board: r.Save(object.Board{
					Name:                in.Name,
					Desc:                in.Desc,
					Created:             time.Now().UnixNano(),
					SubmissionAddresses: in.GetSubmissionAddresses(),
					Meta:                []byte("{}"), // TODO
				}),
			}),
			r.MustDynamic("ThreadVotesPage", object.ThreadVotesPage{}),
			r.MustDynamic("PostVotesPage", object.PostVotesPage{}),
		)
		return e
	})
	return boo.WrapType(e, boo.Internal, "failed to create board")
}

// DeleteBoard deletes a board.
func DeleteBoard(_ context.Context, cxo *state.CXO, in *object.BoardIO) error {
	e := cxo.ModifyRoot(in.GetPK(), in.GetSK(), func(r *node.Root) error {
		_, e := r.Replace(nil)
		return e
	})
	return boo.WrapType(e, boo.Internal, "failed on replacing root references")
}

// NewSubmissionAddress adds a new submission address to board.
func NewSubmissionAddress(_ context.Context, cxo *state.CXO, in *object.AddressIO) error {
	result := NewResult(cxo, in.GetPK(), in.SecKey).
		getPages(true, false, false).
		getBoard()
	defer cxo.Lock()()

	for _, address := range result.Board.SubmissionAddresses {
		if address == in.Address {
			return boo.Newf(boo.AlreadyExists,
				"submission address %s already exists in board %s", in.Address, in.PubKey)
		}
	}
	result.Board.SubmissionAddresses = append(
		result.Board.SubmissionAddresses, in.Address)

	result.saveBoard().savePages(true, false, false)

	if e := result.Error(); e != nil {
		return boo.WrapType(e, boo.NotAuthorised, "secret key invalid")
	}
	return nil
}

// DeleteSubmissionAddress removes a specified submission address from board.
func DeleteSubmissionAddress(_ context.Context, cxo *state.CXO, in *object.AddressIO) error {
	result := NewResult(cxo, in.GetPK(), in.SecKey).
		getPages(true, false, false).
		getBoard()
	defer cxo.Lock()()

	for i, address := range result.Board.SubmissionAddresses {
		if address == in.Address {
			result.Board.SubmissionAddresses = append(
				result.Board.SubmissionAddresses[:i],
				result.Board.SubmissionAddresses[i+1:]...,
			)

			result.saveBoard().savePages(true, false, false)

			if e := result.Error(); e != nil {
				return boo.WrapType(e, boo.NotAuthorised, "secret key invalid")
			}
			return nil
		}
	}
	return boo.Newf(boo.NotFound,
		"submission address %s not found in board %s", in.Address, in.PubKey)
}
