package io

import (
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/skycoin/src/cipher"
)

type BoardPageOut struct {
	Seq         uint64
	BoardPubKey cipher.PubKey
	Threads     []*object.Content
	ThreadVotes []*object.VotesSummary
}

type ThreadPageOut struct {
	Seq        uint64
	Thread     *object.Content
	ThreadVote *object.VotesSummary
	Posts      []*object.Content
	PostVotes  []*object.VotesSummary
}

type FollowPageOut struct {
	Seq        uint64
	FollowPage *object.FollowPage
}

type VoteUserOut struct {
	Seq         uint64
	UserPubKey  cipher.PubKey
	VoteSummary *object.VotesSummary
}

type VoteThreadOut struct {
	Seq         uint64
	ThreadRef   cipher.SHA256
	VoteSummary *object.VotesSummary
}

type VotePostOut struct {
	Seq         uint64
	ThreadRef   cipher.SHA256
	PostRef     cipher.SHA256
	VoteSummary *object.VotesSummary
}

type Delete struct {
	User       cipher.PubKey
	Sig        cipher.SHA256
	ContentRef cipher.SHA256
}
