package content

import (
	"context"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/verify"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state"
	"time"
)

// GetBoardPageResult gets the page of board of public key.
func GetBoardPageResult(_ context.Context, cxo *state.CXO, in *object.BoardIO) (*Result, error) {
	result := NewResult(cxo, in.GetPK()).
		getBoardPage().
		getBoard().
		getThreadPages().
		getThreads()

	if e := result.Error(); e != nil {
		return nil, e
	}

	return result, nil
}

// NewThread creates a new thread on board of specified public key.
func NewThread(_ context.Context, cxo *state.CXO, in *object.NewThreadIO) (*Result, error) {
	result := NewResult(cxo, in.GetBoardPK(), in.BoardSecKey).
		getBoardPage().
		getBoard().
		getThreadPages().
		getThreads()
	defer cxo.Lock()()

	result.Thread = &object.Thread{
		Post: object.Post{
			Title: in.Title,
			Body:  in.Body,
		},
	}

	_, e := verify.Sign(&result.Thread.Post, in.UserPubKey, in.UserSecKey)
	if e != nil {
		return nil, e
	}

	result.Thread.Post.Created = time.Now().UnixNano()

	result.ThreadPages = append(result.ThreadPages, result.ThreadPage)
	result.Threads = append(result.Threads, result.Thread)

	result.saveThread().saveThreadPage().saveBoardPage()

	if e := result.Error(); e != nil {
		return nil, e
	}

	return result, nil
}

// DeleteThread removes a thread of reference from board of public key.
func DeleteThread(_ context.Context, cxo *state.CXO, in *object.ThreadIO) (*Result, error) {
	result := NewResult(cxo, in.GetBoardPK(), in.BoardSecKey).
		getBoardPage().
		getBoard().
		getThreadPages().
		getThreads()
	defer cxo.Lock()()

	for i, tp := range result.ThreadPages {
		if tp.Thread == in.GetThreadRef() {
			result.BoardPage.ThreadPages = append(
				result.BoardPage.ThreadPages[:i],
				result.BoardPage.ThreadPages[i+1:]...,
			)
			result.BoardPage.Deleted = append(
				result.BoardPage.Deleted,
				toSHA256(result.ThreadPages[i].Thread),
			)
			result.ThreadPages = append(
				result.ThreadPages[:i],
				result.ThreadPages[i+1:]...,
			)
			result.Threads = append(
				result.Threads[:i],
				result.Threads[i+1:]...,
			)
			result.saveBoardPage()
			if e := result.Error(); e != nil {
				return nil, e
			}
			return result, nil
		}
	}
	return nil, boo.Newf(boo.NotFound,
		"thread of reference %s not found in board %s",
		in.ThreadRef, in.BoardPubKey)
}
