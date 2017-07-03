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

// Quit exits the node.
func Quit() ClientFunc {
	return gen("quit", Values{})
}

/*
	<<< FOR STATS >>>
*/

// StatsGet obtains stats.
func StatsGet() ClientFunc {
	return gen("stats/get", Values{})
}

/*
	<<< FOR CONNECTIONS >>>
*/

// ConnectionsGetAll gets all connections.
func ConnectionsGetAll() ClientFunc {
	return gen("connections/get_all", Values{})
}

// ConnectionsAdd adds a connection.
func ConnectionsAdd(address *string) ClientFunc {
	return gen("connections/add", Values{
		"address": address,
	})
}

// ConnectionsRemove removes a connections.
func ConnectionsRemove(address *string) ClientFunc {
	return gen("connections/remove", Values{
		"address": address,
	})
}

/*
	<<< FOR SUBSCRIPTIONS >>>
*/

// SubscriptionsGetAll gets all subscriptions.
func SubscriptionsGetAll() ClientFunc {
	return gen("subscriptions/get_all", Values{})
}

// SubscriptionsGet obtains a subscription of board.
func SubscriptionsGet(board *string) ClientFunc {
	return gen("subscriptions/get", Values{
		"board": board,
	})
}

// SubscriptionsAdd adds a subscription of board and address.
func SubscriptionsAdd(board, address *string) ClientFunc {
	return gen("subscriptions/add", Values{
		"board":   board,
		"address": address,
	})
}

// SubscriptionsRemove removes a subscription of board.
func SubscriptionsRemove(board *string) ClientFunc {
	return gen("subscriptions/remove", Values{
		"board": board,
	})
}

/*
	<<< FOR USERS >>>
*/

// UsersGetAll gets all users.
func UsersGetAll() ClientFunc {
	return gen("users/get_all", Values{})
}

// UsersAdd adds a user in which we are not master.
func UsersAdd(user, alias *string) ClientFunc {
	return gen("users/add", Values{
		"user":  user,
		"alias": alias,
	})
}

// UsersRemove removes a user.
func UsersRemove(user *string) ClientFunc {
	return gen("users/remove", Values{
		"user": user,
	})
}

// UsersMastersGetAll gets all master users.
func UsersMastersGetAll() ClientFunc {
	return gen("users/masters/get_all", Values{})
}

// UsersMastersAdd adds a master user.
func UsersMastersAdd(alias, seed *string) ClientFunc {
	return gen("users/masters/add", Values{
		"alias": alias,
		"seed":  seed,
	})
}

// UsersMastersCurrentGet gets the current master user.
func UsersMastersCurrentGet() ClientFunc {
	return gen("users/masters/current/get", Values{})
}

// UsersMastersCurrentSet sets the current master user.
func UsersMastersCurrentSet(user *string) ClientFunc {
	return gen("users/masters/current/set", Values{
		"user": user,
	})
}

/*
	<<< FOR BOARDS >>>
*/

func BoardsGetAll() ClientFunc {
	return gen("boards/get_all", Values{})
}

func BoardsGet(board *string) ClientFunc {
	return gen("boards/get", Values{
		"board": board,
	})
}

func BoardsAdd(boardName, boardDescription, boardSubmissionAddresses, seed *string) ClientFunc {
	return gen("boards/add", Values{
		"name":                 boardName,
		"description":          boardDescription,
		"submission_addresses": boardSubmissionAddresses,
		"seed":                 seed,
	})
}

func BoardsRemove(board *string) ClientFunc {
	return gen("boards/remove", Values{
		"board": board,
	})
}

func BoardsMetaGet(board *string) ClientFunc {
	return gen("boards/meta/get", Values{
		"board": board,
	})
}

func BoardsMetaSubmissionAddressesGetAll(board *string) ClientFunc {
	return gen("boards/meta/submission_addresses/get_all", Values{
		"board": board,
	})
}

func BoardsMetaSubmissionAddressesAdd(board, address *string) ClientFunc {
	return gen("boards/meta/submission_addresses/add", Values{
		"board":   board,
		"address": address,
	})
}

func BoardsMetaSubmissionAddressesRemove(board, address *string) ClientFunc {
	return gen("boards/meta/submission_addresses/remove", Values{
		"board":   board,
		"address": address,
	})
}

func BoardsPageGet(board *string) ClientFunc {
	return gen("boards/page/get", Values{
		"board": board,
	})
}

/*
	<<< FOR THREADS >>>
*/

func ThreadsGetAll(board *string) ClientFunc {
	return gen("threads/get_all", Values{
		"board": board,
	})
}

func ThreadsAdd(board, threadName, threadDescription *string) ClientFunc {
	return gen("threads/add", Values{
		"board":       board,
		"name":        threadName,
		"description": threadDescription,
	})
}

func ThreadsRemove(board, thread *string) ClientFunc {
	return gen("threads/remove", Values{
		"board":  board,
		"thread": thread,
	})
}

func ThreadsImport(fromBoard, thread, toBoard *string) ClientFunc {
	return gen("threads/import", Values{
		"from_board": fromBoard,
		"thread":     thread,
		"to_board":   toBoard,
	})
}

func ThreadsPageGet(board, thread *string) ClientFunc {
	return gen("threads/page/get", Values{
		"board":  board,
		"thread": thread,
	})
}

func ThreadsVotesGet(board, thread *string) ClientFunc {
	return gen("threads/votes/get", Values{
		"board":  board,
		"thread": thread,
	})
}

func ThreadsVotesAdd(board, thread, voteMode, voteTag *string) ClientFunc {
	return gen("threads/votes/add", Values{
		"board":  board,
		"thread": thread,
		"mode":   voteMode,
		"tag":    voteTag,
	})
}

/*
	<<< FOR POSTS >>>
*/

func PostsGetAll(board, thread *string) ClientFunc {
	return gen("posts/get_all", Values{
		"board":  board,
		"thread": thread,
	})
}

func PostsAdd(board, thread, postTitle, postBody *string) ClientFunc {
	return gen("posts/add", Values{
		"board":  board,
		"thread": thread,
		"title":  postTitle,
		"body":   postBody,
	})
}

func PostsRemove(board, thread, post *string) ClientFunc {
	return gen("posts/remove", Values{
		"board":  board,
		"thread": thread,
		"post":   post,
	})
}

func PostsVotesGet(board, post *string) ClientFunc {
	return gen("posts/votes/get", Values{
		"board": board,
		"post":  post,
	})
}

func PostsVotesAdd(board, post, voteMode, voteTag *string) ClientFunc {
	return gen("posts/votes/add", Values{
		"board": board,
		"post":  post,
		"mode":  voteMode,
		"tag":   voteTag,
	})
}

/*
	<<< FOR TESTS >>>
*/

func TestsAddFilledBoard(seed, threads, minPosts, maxPosts *string) ClientFunc {
	return gen("tests/add_filled_board", Values{
		"seed":      seed,
		"threads":   threads,
		"min_posts": minPosts,
		"max_posts": maxPosts,
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
