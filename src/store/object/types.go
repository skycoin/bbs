package object

import (
	"github.com/skycoin/bbs/src/misc/tag"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
)

type ThreadPages struct {
	Board       skyobject.Ref  `skyobject:"schema=bbs.Board"`
	ThreadPages skyobject.Refs `skyobject:"schema=bbs.ThreadPage"`
}

type ThreadPage struct {
	Thread skyobject.Ref  `skyobject:"schema=bbs.Content"`
	Posts  skyobject.Refs `skyobject:"schema=bbs.Content"`
}

type ThreadVotesPages struct {
	Store   []ContentVotesPage
	Deleted []cipher.SHA256
}

type PostVotesPages struct {
	Store   []ContentVotesPage
	Deleted []cipher.SHA256
}

type ContentVotesPage struct {
	Ref   cipher.SHA256
	Votes skyobject.Refs `skyobject:"schema=bbs.Vote"`
}

type UserVotesPages struct {
	Store   []UserVotesPage
	Deleted []cipher.PubKey
}

type UserVotesPage struct {
	Ref   cipher.PubKey
	Votes skyobject.Refs `skyobject:"schema=bbs.Vote"`
}

/*
	<<< BOARD >>>
*/

type Board struct {
	Name    string `json:"name" trans:"heading"`
	Desc    string `json:"description" trans:"body"`
	Created int64  `json:"created" trans:"time"`
	Meta    []byte `json:"-"` // TODO: Recommended Submission Addresses.
}

type BoardView struct {
	Board
	BoardHash cipher.SHA256 `json:"-"`
	PubKey    string        `json:"public_key"`
}

/*
	<<< CONTENT >>>
*/

type Content struct {
	Title   string        `json:"title" trans:"heading"`
	Body    string        `json:"body" trans:"body"`
	Created int64         `json:"created" trans:"time"`
	Creator cipher.PubKey `json:"-" verify:"upk" trans:"upk"`
	Sig     cipher.Sig    `json:"-" verify:"sig"`
	Meta    []byte        `json:"-"`
}

func (c Content) Verify() error { return tag.Verify(&c) }

type ContentView struct {
	Content
	Ref     string `json:"reference"`
	Creator User   `json:"creator"`
}

/*
	<<< VOTES >>>
*/

type Vote struct {
	Mode    int8          `json:"mode" trans:"mode"`
	Tag     string        `json:"-" trans:"tag"` // TODO: Fix transfer.
	Created int64         `json:"created" trans:"time"`
	Creator cipher.PubKey `json:"-" verify:"upk" trans:"upk"`
	Sig     cipher.Sig    `json:"-" verify:"sig"`
}

func (v Vote) Verify() error { return tag.Verify(&v) }

type ContentVotesSummary struct {
	sync.Mutex
	Hash  cipher.SHA256
	Votes map[cipher.PubKey]Vote
	Up    CompiledVotes
	Down  CompiledVotes
}

func (s *ContentVotesSummary) View(perspective cipher.PubKey) ContentVotesSummaryView {
	s.Lock()
	defer s.Unlock()
	vote := s.Votes[perspective]
	return ContentVotesSummaryView{
		Up: CompiledVotesView{
			CompiledVotes: s.Up,
			Voted:         vote.Mode == +1,
		},
		Down: CompiledVotesView{
			CompiledVotes: s.Down,
			Voted:         vote.Mode == -1,
		},
	}
}

type ContentVotesSummaryView struct {
	Up   CompiledVotesView `json:"up"`
	Down CompiledVotesView `json:"down"`
}

type CompiledVotes struct {
	Count int            `json:"count"`
	Tags  map[string]int `json:"tags"`
}

type CompiledVotesView struct {
	CompiledVotes
	Voted bool `json:"voted"`
}

/*
	<<< USER >>>
*/

type User struct {
	Alias  string        `json:"alias" trans:"alias"`
	PubKey cipher.PubKey `json:"-" trans:"upk"`
	SecKey cipher.SecKey `json:"-" trans:"usk"`
}

type UserView struct {
	User
	PubKey string `json:"public_key"`
	SecKey string `json:"secret_key,omitempty"`
}
