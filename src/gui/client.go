package gui

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	ErrTimeout = errors.New("timeout")
)

type (
	// ClientFunc is a client function.
	ClientFunc func(ctx context.Context, port int) ([]byte, error)

	// Values represents key/value pairs of form.
	Values map[string]*string
)

/*
	<<< FOR BOARD META >>>
*/

// GetSubmissionAddresses obtains the submission addresses of specified board.
func GetSubmissionAddresses(board *string) ClientFunc {
	return gen("board_meta/get_submission_addresses", Values{
		"board": board,
	})
}

// AddSubmissionAddress adds a submission address to specified board.
func AddSubmissionAddress(board, address *string) ClientFunc {
	return gen("board_meta/add_submission_address", Values{
		"board":   board,
		"address": address,
	})
}

// RemoveSubmissionAddress removes a submission address from specified board.
func RemoveSubmissionAddress(board, address *string) ClientFunc {
	return gen("board_meta/remove_submission_address", Values{
		"board":   board,
		"address": address,
	})
}

/*
	<<< FOR BOARDS, THREADS & POSTS >>>
*/

// GetBoards obtains boards in which the bbsnode is subscribed.
func GetBoards() ClientFunc {
	return gen("get_boards", nil)
}

// AddBoard adds a new board.
func AddBoard(boardName, boardDescription, boardSubmissionAddresses, seed *string) ClientFunc {
	return gen("add_board", Values{
		"name":                 boardName,
		"description":          boardDescription,
		"submission_addresses": boardSubmissionAddresses,
		"seed":                 seed,
	})
}

// RemoveBoard removes a board.
func RemoveBoard(board *string) ClientFunc {
	return gen("remove_board", Values{
		"board": board,
	})
}

// GetBoardPage obtains the board page of specified board of public key.
func GetBoardPage(board *string) ClientFunc {
	return gen("get_board_page", Values{
		"board": board,
	})
}

// GetThreads obtains threads of a specified board of public key.
func GetThreads(board *string) ClientFunc {
	return gen("get_threads", Values{
		"board": board,
	})
}

// AddThread adds a new thread on the specified board.
func AddThread(board, threadName, threadDescription *string) ClientFunc {
	return gen("add_thread", Values{
		"board":       board,
		"name":        threadName,
		"description": threadDescription,
	})
}

// RemoveThread removes a thread from the specified board.
func RemoveThread(board, thread *string) ClientFunc {
	return gen("remove_thread", Values{
		"board":  board,
		"thread": thread,
	})
}

// GetThreadPage obtains a thread page of specified board and thread.
func GetThreadPage(board, thread *string) ClientFunc {
	return gen("get_thread_page", Values{
		"board":  board,
		"thread": thread,
	})
}

// GetPosts obtains the posts of a thread of specified board and thread.
func GetPosts(board, thread *string) ClientFunc {
	return gen("get_posts", Values{
		"board":  board,
		"thread": thread,
	})
}

// AddPost adds a new post on specified board and thread.
func AddPost(board, thread, postTitle, postBody *string) ClientFunc {
	return gen("add_post", Values{
		"board":  board,
		"thread": thread,
		"title":  postTitle,
		"body":   postBody,
	})
}

// RemovePost removes a post in specified board, thread and post reference.
func RemovePost(board, thread, post *string) ClientFunc {
	return gen("remove_post", Values{
		"board":  board,
		"thread": thread,
		"post":   post,
	})
}

// ImportThread imports a thread from a board to another.
func ImportThread(fromBoard, thread, toBoard *string) ClientFunc {
	return gen("import_thread", Values{
		"from_board": fromBoard,
		"thread":     thread,
		"to_board":   toBoard,
	})
}

/*
	<<< FOR VOTES >>>
*/

// GetThreadVotes obtains votes for specified board and thread.
func GetThreadVotes(board, thread *string) ClientFunc {
	return gen("get_thread_votes", Values{
		"board":  board,
		"thread": thread,
	})
}

// GetPostVotes obtains votes for specified board and post.
func GetPostVotes(board, post *string) ClientFunc {
	return gen("get_post_votes", Values{
		"board": board,
		"post":  post,
	})
}

// AddThreadVote adds a vote to a thread of specified board.
func AddThreadVote(board, thread, voteMode, voteTag *string) ClientFunc {
	return gen("add_thread_vote", Values{
		"board":  board,
		"thread": thread,
		"mode":   voteMode,
		"tag":    voteTag,
	})
}

// AddPostVote adds a vote to post of specified board.
func AddPostVote(board, post, voteMode, voteTag *string) ClientFunc {
	return gen("add_post_vote", Values{
		"board": board,
		"post":  post,
		"mode":  voteMode,
		"tag":   voteTag,
	})
}

/*
	<<< HELPER FUNCTIONS >>>
*/

// Asynchronously requests from api.
func request(port int, path string, data url.Values) (chan []byte, chan error) {
	bChan, eChan := make(chan []byte), make(chan error)
	go func() {
		resp, e := http.PostForm(
			fmt.Sprintf("http://127.0.0.1:%d/api/%s", port, path),
			data,
		)
		if e != nil {
			eChan <- e
			return
		}
		defer resp.Body.Close()
		body, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			eChan <- e
			return
		}
		bChan <- body
		return
	}()
	return bChan, eChan
}

// Generates a method of requesting data from api.
func gen(path string, values Values) ClientFunc {
	return func(ctx context.Context, port int) ([]byte, error) {
		// Get form values.
		urlValues := url.Values{}
		for k, v := range values {
			urlValues[k] = []string{*v}
		}

		// Send request.
		bChan, eChan := request(port, path, urlValues)

		// Await reply.
		select {
		case <-ctx.Done():
			return nil, ErrTimeout
		case e := <-eChan:
			return nil, e
		case body := <-bChan:
			return body, nil
		}
	}
}
