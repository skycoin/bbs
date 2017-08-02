package states

import (
	"context"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

type NewState func(bpk, upk cipher.PubKey, workerChan chan<- func()) State

type State interface {
	Close()
	Trigger(ctx context.Context, root *node.Root)
	GetThreadVotes(ref skyobject.Reference) *object.VoteSummary
	GetThreadVotesSeq(ctx context.Context, ref skyobject.Reference, seq uint64) *object.VoteSummary
	GetPostVotes(ref skyobject.Reference) *object.VoteSummary
	GetPostVotesSeq(ctx context.Context, ref skyobject.Reference, seq uint64) *object.VoteSummary
}

