package states

import (
	"context"
	"github.com/skycoin/bbs/src/store/io"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

// StateConfig configures the state on a per-board basis.
type StateConfig struct {
	Master bool          // Whether this node owns the board.
	PubKey cipher.PubKey // Public key of board.
	SecKey cipher.SecKey // Secret key of board (only if master).
}

type PublishFunc func(node *node.Node, pack *skyobject.Pack) error

// State represents the compiled state of a board.
// It takes in a root that represents a board and compiles it into a state
// where data can easily be extracted.
type State interface {
	// Init initiates the state. Enabling access to worker channel.
	Init(config *StateConfig, node *node.Node) error

	// Close closes the state.
	Close()

	// IncomingChan obtains a channel used to update the board state.
	IncomingChan() chan<- *io.State

	// ChangesChan obtains a channel used for updating via WebSocket.
	ChangesChan() <-chan *io.Changes

	// Publish publishes the root.
	Publish(ctx context.Context, publish PublishFunc) error

	// GetBoardPage obtains a BoardPage.
	GetBoardPage(ctx context.Context) (
		*io.BoardPageOut, error)

	// GetThreadPage obtains a ThreadPage.
	GetThreadPage(ctx context.Context, tRef cipher.SHA256) (
		*io.ThreadPageOut, error)

	// GetFollowPage obtains a FollowPage.
	GetFollowPage(ctx context.Context, upk cipher.PubKey) (
		*io.FollowPageOut, error)

	// GetUserVotes obtains a user's votes.
	GetUserVotes(ctx context.Context, upk cipher.PubKey) (
		*io.VoteUserOut, error)

	// NewThread creates a new thread, and returns BoardPage.
	NewThread(ctx context.Context, thread *object.Content) (
		*io.BoardPageOut, error)

	// NewPost creates a new post on specified thread and returns ThreadPage.
	NewPost(ctx context.Context, post *object.Content) (
		*io.ThreadPageOut, error)

	// DeleteThread removes a thread and returns BoardPage.
	DeleteThread(ctx context.Context, tRef cipher.SHA256) (
		*io.BoardPageOut, error)

	// DeletePost removes a post and returns ThreadPage.
	DeletePost(ctx context.Context, pRef cipher.SHA256) (
		*io.ThreadPageOut, error)

	// VoteThread adds/removes/modifies a vote on a thread.
	// (Remove if vote is nil).
	VoteThread(ctx context.Context, vote *cipher.SHA256) (
		*io.VoteThreadOut, error)

	// VotePost adds/removes/modifies a vote on a post.
	// (Remove if vote is nil).
	VotePost(ctx context.Context, vote *cipher.SHA256) (
		*io.VotePostOut, error)

	// VoteUser adds/removes/modifies a vote on a user.
	// (Remove if vote is nil).
	VoteUser(ctx context.Context, vote *cipher.SHA256) (
		*io.VoteUserOut, error)
}
