package object

import (
	"github.com/skycoin/skycoin/src/cipher"
	"strings"
	"sync"
)

type BoardView struct {
	*Board
	PublicKey     string              `json:"public_key"`
	ExternalRoots []*ExternalRootView `json:"external_roots,omitempty"`
}

type ExternalRootView struct {
	ExternalRoot
	PublicKey string `json:"public_key"`
}

type ThreadView struct {
	*Thread
	Ref         string           `json:"reference"`
	AuthorRef   string           `json:"author_reference,omitempty"`
	AuthorAlias string           `json:"author_alias,omitempty"`
	Votes       *VoteSummaryView `json:"votes,omitempty"`
}

type PostView struct {
	*Post
	Ref         string           `json:"reference"`
	AuthorRef   string           `json:"author_reference,omitempty"`
	AuthorAlias string           `json:"author_alias,omitempty"`
	Votes       *VoteSummaryView `json:"votes"`
}

/*
	<<< VOTING >>>
*/

type VoteSummary struct {
	sync.Mutex
	Hash  cipher.SHA256
	Votes map[cipher.PubKey]Vote
	Ups   int
	Downs int
	Spams int
}

func NewVoteSummary() *VoteSummary {
	return &VoteSummary{
		Votes: make(map[cipher.PubKey]Vote),
	}
}

func (s *VoteSummary) GenerateView(upk cipher.PubKey) *VoteSummaryView {
	s.Lock()
	uVote := s.Votes[upk]
	s.Unlock()
	return &VoteSummaryView{
		Up: VoteView{
			Voted: uVote.Mode == +1,
			Count: s.Ups,
		},
		Down: VoteView{
			Voted: uVote.Mode == -1,
			Count: s.Downs,
		},
		Spam: VoteView{
			Voted: strings.Contains(
				string(uVote.Tag), "spam"),
			Count: s.Spams,
		},
	}
}

type VoteSummaryView struct {
	Up   VoteView `json:"up"`
	Down VoteView `json:"down"`
	Spam VoteView `json:"spam"`
}

type VoteView struct {
	Voted bool `json:"voted"`
	Count int  `json:"count"`
}

type UserView struct {
	User
	PublicKey string           `json:"public_key,omitempty"`
	SecretKey string           `json:"secret_key,omitempty"`
	Votes     *VoteSummaryView `json:"votes,omitempty"`
}

type SubscriptionView struct {
	PubKey string `json:"public_key"`
	SecKey string `json:"secret_key,omitempty"`
}

type ConnectionView struct {
	Address string `json:"address"`
	Active  bool   `json:"active"`
}
