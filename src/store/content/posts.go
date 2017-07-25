package content

import (
	"context"
	"github.com/skycoin/bbs/src/misc/verify"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state"
	"time"
	"github.com/skycoin/bbs/src/misc/boo"
)

// GetThreadPageResult gets the page of thread of reference from board of public key.
func GetThreadPageResult(_ context.Context, cxo *state.CXO, in *object.ThreadIO) (*Result, error) {
	result := NewResult(cxo, in.GetBoardPK()).
		getBoardPage().
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
		getBoardPage().
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

	result.savePost().saveThreadPage().saveBoardPage()

	if e := result.Error(); e != nil {
		return nil, e
	}

	return result, nil
}

// DeletePost removes a post of reference from thread of reference and board of public key.
func DeletePost(_ context.Context, cxo *state.CXO, in *object.PostIO) (*Result, error) {
	result := NewResult(cxo, in.GetBoardPK(), in.BoardSecKey).
		getBoardPage().
		getBoard().
		getThreadPage(in.GetThreadRef()).
		getThread().
		getPosts()
	defer cxo.Lock()()

	for i, p := range result.Posts {
		if toRef(p.R) == in.GetPostRef() {
			result.ThreadPage.Posts = append(
				result.ThreadPage.Posts[:i],
				result.ThreadPage.Posts[i+1:]...,
			)
			result.ThreadPage.Deleted = append(
				result.ThreadPage.Deleted,
				p.R,
			)
			result.Posts = append(
				result.Posts[:i],
				result.Posts[i+1:]...,
			)
			result.saveThreadPage().saveBoardPage()
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
