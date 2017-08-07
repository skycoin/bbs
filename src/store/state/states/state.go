package states

import (
	"context"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

// NewState creates a new State.
type NewState func(bpk cipher.PubKey, workerChan chan<- func()) State

// State represents the compiled state of a board.
type State interface {
	// Close closes the state.
	Close()

	// Trigger is used to update the state.
	Trigger(ctx context.Context, root *node.Root)

	// GetThreadVotes obtains the votes of thread of reference.
	GetThreadVotes(ref skyobject.Reference) *object.VoteSummary

	// GetThreadVotesSeq obtains the thread votes above the specified sequence.
	GetThreadVotesSeq(ctx context.Context, ref skyobject.Reference, seq uint64) *object.VoteSummary

	// GetPostVotes obtains the votes of post of reference.
	GetPostVotes(ref skyobject.Reference) *object.VoteSummary

	// GetPostVotesSeq obtains the post votes above the specified sequence.
	GetPostVotesSeq(ctx context.Context, ref skyobject.Reference, seq uint64) *object.VoteSummary

	// GetUserVotes obtains the votes of users of public key.
	GetUserVotes(upk cipher.PubKey) *object.VoteSummary

	// GetUserVotesSeq obtains the user votes above the specified sequence.
	GetUserVotesSeq(ctx context.Context, upk cipher.PubKey, seq uint64) *object.VoteSummary
}
