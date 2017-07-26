package content

import (
	"context"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/verify"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state"
	"time"
	"sync"
)

// GetBoardPageResult gets the page of board of public key.
func GetBoardPageResult(_ context.Context, cxo *state.CXO, in *object.BoardIO) (*Result, error) {
	result := NewResult(cxo, in.GetPK()).
		getPages(true, false, false).
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
		getPages(true, true, false).
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

	result.saveThread().saveThreadPage().savePages(true, true, false)

	if e := result.Error(); e != nil {
		return nil, e
	}

	return result, nil
}

// DeleteThread removes a thread of reference from board of public key.
func DeleteThread(_ context.Context, cxo *state.CXO, in *object.ThreadIO) (*Result, error) {
	result := NewResult(cxo, in.GetBoardPK(), in.BoardSecKey).
		getPages(true, true, true).
		getBoard().
		getThreadPages().
		getThreads()
	defer cxo.Lock()()

	for i, tp := range result.ThreadPages {
		if tp.Thread == in.GetThreadRef() {
			var wg sync.WaitGroup
			wg.Add(3)
			go func() {
				defer wg.Done()
				result.deleteThreadVote(in.GetThreadRef())
			} ()
			go func() {
				defer wg.Done()
				result.deletePostVotes(tp.Posts)
			} ()
			go func() {
				defer wg.Done()
				result.deleteThread(i)
			} ()
			wg.Wait()

			// Save.
			result.savePages(true, true, true)
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

// VoteThread adds/modifies/removes vote from thread.
func VoteThread(_ context.Context, cxo *state.CXO, in *object.VoteThreadIO) (*Result, error) {
	result := NewResult(cxo, in.GetBoardPK(), in.BoardSecKey).
		getPages(false, true, false)
	defer cxo.Lock()()

	tvi, has := result.ThreadRefMap[toSHA256(in.GetThreadRef())]
	if !has {
		return nil, boo.Newf(boo.NotFound,
			"thread of reference %s not found in board %s",
			in.ThreadRef, in.BoardPubKey)
	}

	var vote object.Vote
	for i, vRef := range result.ThreadVotesPage.Store[tvi].Votes {
		if e := result.deserialize(vRef, &vote); e != nil {
			return nil, boo.WrapTypef(e, boo.InvalidRead,
				"vote %d from thread %s of board %s is corrupt",
				i, in.ThreadRef, in.BoardPubKey)
		}
		if vote.User == in.UserPubKey {
			vote.Mode = in.GetMode()
			vote.Tag = in.GetTag()
			vote.Created = time.Now().UnixNano()
			if _, e := verify.Sign(&vote, in.UserPubKey, in.UserSecKey); e != nil {
				return nil, e
			}
			result.ThreadVotesPage.Store[tvi].Votes[i] =
				result.root.Save(vote)
			e := result.savePages(false, true, false).Error()
			if e != nil {
				return nil, e
			}
			return result, nil
		}
	}
	vote.Mode = in.GetMode()
	vote.Tag = in.GetTag()
	vote.Created = time.Now().UnixNano()
	if _, e := verify.Sign(&vote, in.UserPubKey, in.UserSecKey); e != nil {
		return nil, e
	}
	result.ThreadVotesPage.Store[tvi].Votes = append(
		result.ThreadVotesPage.Store[tvi].Votes,
		result.root.Save(vote),
	)
	e := result.savePages(false, true, false).Error()
	if e != nil {
		return nil, e
	}
	return result, nil
}
