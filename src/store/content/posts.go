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

// GetThreadPageResult gets the page of thread of reference from board of public key.
func GetThreadPageResult(_ context.Context, cxo *state.CXO, in *object.ThreadIO) (*Result, error) {
	result := NewResult(cxo, in.GetBoardPK()).
		getPages(true, false, false).
		getBoard().
		getThreadPage(in.GetThreadRef()).
		getThread().
		getPosts()

	if e := result.Error(); e != nil {
		return nil, e
	}

	return result, nil
}

// NewPost creates a new post on thread of reference from board of public key.
func NewPost(_ context.Context, cxo *state.CXO, in *object.NewPostIO) (*Result, error) {
	result := NewResult(cxo, in.GetBoardPK(), in.BoardSecKey).
		getPages(true, true, true).
		getBoard().
		getThreadPage(in.GetThreadRef()).
		getThread().
		getPosts()
	defer cxo.Lock()()

	result.Post = &object.Post{
		Title: in.Title,
		Body:  in.Body,
	}

	_, e := verify.Sign(result.Post, in.UserPubKey, in.UserSecKey)
	if e != nil {
		return nil, e
	}

	result.Post.Created = time.Now().UnixNano()
	result.Posts = append(result.Posts, result.Post)

	result.savePost().saveThreadPage().savePages(true, true, true)

	if e := result.Error(); e != nil {
		return nil, e
	}

	return result, nil
}

// DeletePost removes a post of reference from thread of reference and board of public key.
func DeletePost(_ context.Context, cxo *state.CXO, in *object.PostIO) (*Result, error) {
	result := NewResult(cxo, in.GetBoardPK(), in.BoardSecKey).
		getPages(true, false, true).
		getBoard().
		getThreadPage(in.GetThreadRef()).
		getThread().
		getPosts()
	defer cxo.Lock()()

	for i, p := range result.Posts {
		if toRef(p.R) == in.GetPostRef() {
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				result.deletePostVote(in.GetPostRef())
			} ()
			go func() {
				defer wg.Done()
				result.deletePost(i)
			} ()
			wg.Wait()

			result.saveThreadPage().savePages(true, false, true)
			if e := result.Error(); e != nil {
				return nil, e
			}
			return result, nil
		}
	}
	return nil, boo.Newf(boo.NotFound,
		"post of reference %s not found on thread %s of board %s",
		in.PostRef, in.ThreadRef, in.BoardPubKey)
}
