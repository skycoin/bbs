package content

import (
	"context"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/node"
	"sync"
	"time"
)

// GetBoardResult get's the specified board of public key.
func GetBoardResult(_ context.Context, root *node.Root) (*Result, error) {
	result := NewResult(root).
		GetPages(true, false, false, false).
		GetBoard()
	if e := result.Error(); e != nil {
		return nil, e
	}
	return result, nil
}

// NewBoard creates a new board and returns an error on failure.
func NewBoard(_ context.Context, root *node.Root, in *object.NewBoardIO) error {
	_, e := root.Append(
		root.MustDynamic("BoardPage", object.BoardPage{
			Board: root.Save(object.Board{
				Name:                in.Name,
				Desc:                in.Desc,
				Created:             time.Now().UnixNano(),
				SubmissionAddresses: in.SubmissionAddresses,
				Meta:                []byte("{}"), // TODO
			}),
		}),
		root.MustDynamic("ThreadVotesPages", object.ThreadVotesPages{}),
		root.MustDynamic("PostVotesPages", object.PostVotesPages{}),
		root.MustDynamic("UserVotesPages", object.UserVotesPages{}),
	)
	return boo.WrapType(e, boo.Internal, "failed to create board")
}

// DeleteBoard deletes a board.
func DeleteBoard(_ context.Context, root *node.Root, _ *object.BoardIO) error {
	_, e := root.Replace(nil)
	return boo.WrapType(e, boo.Internal, "failed on replacing root references")
}

// NewSubmissionAddress adds a new submission address to board.
func NewSubmissionAddress(_ context.Context, root *node.Root, in *object.AddressIO) error {
	result := NewResult(root).
		GetPages(true, false, false, false).
		GetBoard()
	if e := result.Error(); e != nil {
		return e
	}
	for _, address := range result.Board.SubmissionAddresses {
		if address == in.Address {
			return boo.Newf(boo.AlreadyExists,
				"submission address %s already exists in board %s", in.Address, in.PubKeyStr)
		}
	}
	result.Board.SubmissionAddresses = append(
		result.Board.SubmissionAddresses, in.Address)

	result.saveBoard().savePages(true, false, false, false)

	if e := result.Error(); e != nil {
		return boo.WrapType(e, boo.NotAuthorised, "secret key invalid")
	}
	return nil
}

// DeleteSubmissionAddress removes a specified submission address from board.
func DeleteSubmissionAddress(_ context.Context, root *node.Root, in *object.AddressIO) error {
	result := NewResult(root).
		GetPages(true, false, false, false).
		GetBoard()
	if e := result.Error(); e != nil {
		return e
	}
	for i, address := range result.Board.SubmissionAddresses {
		if address == in.Address {
			result.Board.SubmissionAddresses = append(
				result.Board.SubmissionAddresses[:i],
				result.Board.SubmissionAddresses[i+1:]...,
			)
			result.saveBoard().savePages(true, false, false, false)
			if e := result.Error(); e != nil {
				return boo.WrapType(e, boo.NotAuthorised, "secret key invalid")
			}
			return nil
		}
	}
	return boo.Newf(boo.NotFound,
		"submission address %s not found in board %s", in.Address, in.PubKeyStr)
}

// GetBoardPageResult gets the page of board of public key.
func GetBoardPageResult(_ context.Context, root *node.Root, _ *object.BoardIO) (*Result, error) {
	result := NewResult(root).
		GetPages(true, false, false, false).
		GetBoard().
		GetThreadPages().
		GetThreads()
	if e := result.Error(); e != nil {
		return nil, e
	}
	return result, nil
}

// NewThread creates a new thread on board of specified public key.
func NewThread(_ context.Context, root *node.Root, in *object.NewThreadIO) (*Result, error) {
	result := NewResult(root).
		GetPages(true, true, false, true).
		GetBoard().
		GetThreadPages().
		GetThreads()
	if e := result.Error(); e != nil {
		return nil, e
	}
	if e := in.Thread.Verify(); e != nil {
		return nil, e
	}
	if _, has := result.UserMap[in.Thread.User]; !has {
		result.UserVotesPages.Store = append(
			result.UserVotesPages.Store,
			object.UserVotesPage{PubKey: in.Thread.User},
		)
	}
	result.Thread = in.Thread
	result.ThreadPages = append(result.ThreadPages, result.ThreadPage)
	result.Threads = append(result.Threads, result.Thread)
	result.
		saveThread().
		saveThreadPage().
		savePages(true, true, false, true)
	if e := result.Error(); e != nil {
		return nil, e
	}
	return result, nil
}

// DeleteThread removes a thread of reference from board of public key.
func DeleteThread(_ context.Context, root *node.Root, in *object.ThreadIO) (*Result, error) {
	result := NewResult(root).
		GetPages(true, true, true, false).
		GetBoard().
		GetThreadPages().
		GetThreads()
	if e := result.Error(); e != nil {
		return nil, e
	}
	for i, tp := range result.ThreadPages {
		if tp.Thread == in.ThreadRef {
			var wg sync.WaitGroup
			wg.Add(3)
			go func() {
				defer wg.Done()
				result.deleteThreadVote(in.ThreadRef)
			}()
			go func() {
				defer wg.Done()
				result.deletePostVotes(tp.Posts)
			}()
			go func() {
				defer wg.Done()
				result.deleteThread(i)
			}()
			wg.Wait()
			result.
				savePages(true, true, true, false)
			if e := result.Error(); e != nil {
				return nil, e
			}
			return result, nil
		}
	}
	return nil, boo.Newf(boo.NotFound,
		"thread of reference %s not found in board %s",
		in.ThreadRefStr, in.BoardPubKeyStr)
}

// VoteThread adds/modifies/removes vote from thread.
func VoteThread(_ context.Context, root *node.Root, in *object.VoteThreadIO) (*Result, error) {
	result := NewResult(root).
		GetPages(false, true, false, true)
	if e := result.Error(); e != nil {
		return nil, e
	}
	if e := in.Vote.Verify(); e != nil {
		return nil, e
	}
	if _, has := result.UserMap[in.Vote.User]; !has {
		result.UserVotesPages.Store = append(
			result.UserVotesPages.Store,
			object.UserVotesPage{PubKey: in.Vote.User},
		)
	}
	result.ThreadVote = in.Vote
	result.
		saveThreadVote(in.ThreadRef).
		savePages(false, true, false, true)
	if e := result.Error(); e != nil {
		return nil, e
	}
	return result, nil
}

// GetThreadPageResult gets the page of thread of reference from board of public key.
func GetThreadPageResult(_ context.Context, root *node.Root, in *object.ThreadIO) (*Result, error) {
	result := NewResult(root).
		GetPages(true, false, false, false).
		GetBoard().
		GetThreadPage(in.ThreadRef).
		GetThread().
		GetPosts()
	if e := result.Error(); e != nil {
		return nil, e
	}
	return result, nil
}

// NewPost creates a new post on thread of reference from board of public key.
func NewPost(_ context.Context, root *node.Root, in *object.NewPostIO) (*Result, error) {
	result := NewResult(root).
		GetPages(true, true, true, true).
		GetBoard().
		GetThreadPage(in.ThreadRef).
		GetThread().
		GetPosts()
	if e := result.Error(); e != nil {
		return nil, e
	}
	if e := in.Post.Verify(); e != nil {
		return nil, e
	}
	if _, has := result.UserMap[in.Post.User]; !has {
		result.UserVotesPages.Store = append(
			result.UserVotesPages.Store,
			object.UserVotesPage{PubKey: in.Post.User},
		)
	}
	result.Post = in.Post
	result.Post.Created = time.Now().UnixNano()
	result.Posts = append(result.Posts, result.Post)
	result.
		savePost().
		saveThreadPage().
		savePages(true, true, true, true)
	if e := result.Error(); e != nil {
		return nil, e
	}
	return result, nil
}

// DeletePost removes a post of reference from thread of reference and board of public key.
func DeletePost(_ context.Context, root *node.Root, in *object.PostIO) (*Result, error) {
	result := NewResult(root).
		GetPages(true, false, true, false).
		GetBoard().
		GetThreadPage(in.ThreadRef).
		GetThread().
		GetPosts()
	if e := result.Error(); e != nil {
		return nil, e
	}
	for i, p := range result.Posts {
		if toRef(p.R) == in.PostRef {
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				result.deletePostVote(in.PostRef)
			}()
			go func() {
				defer wg.Done()
				result.deletePost(i)
			}()
			wg.Wait()
			result.
				saveThreadPage().
				savePages(true, false, true, false)
			if e := result.Error(); e != nil {
				return nil, e
			}
			return result, nil
		}
	}
	return nil, boo.Newf(boo.NotFound,
		"post of reference %s not found on thread %s of board %s",
		in.PostRefStr, in.ThreadRefStr, in.BoardPubKeyStr)
}

// VotePost adds/modifies/removes vote from thread.
func VotePost(_ context.Context, root *node.Root, in *object.VotePostIO) (*Result, error) {
	result := NewResult(root).
		GetPages(false, false, true, true)
	if e := result.Error(); e != nil {
		return nil, e
	}
	if e := in.Vote.Verify(); e != nil {
		return nil, e
	}
	if _, has := result.UserMap[in.Vote.User]; !has {
		result.UserVotesPages.Store = append(
			result.UserVotesPages.Store,
			object.UserVotesPage{PubKey: in.Vote.User},
		)
	}
	result.PostVote = in.Vote
	result.
		savePostVote(in.PostRef).
		savePages(false, false, true, true)
	if e := result.Error(); e != nil {
		return nil, e
	}
	return result, nil
}

// VoteUser places a vote on a user.
func VoteUser(_ context.Context, root *node.Root, in *object.VoteUserIO) (*Result, error) {
	result := NewResult(root).
		GetPages(false, false, false, true)
	if e := result.Error(); e != nil {
		return nil, e
	}
	if e := in.Vote.Verify(); e != nil {
		return nil, e
	}
	result.UserVote = in.Vote
	result.
		saveUserVote(in.UserRef).
		savePages(false, false, false, true)
	if e := result.Error(); e != nil {
		return nil, e
	}
	return result, nil
}
