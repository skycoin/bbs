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
	Threads []ContentVotesPage
}

type PostVotesPages struct {
	Posts []ContentVotesPage
}

type ContentVotesPage struct {
	Ref   cipher.SHA256
	Votes skyobject.Refs `skyobject:"schema=bbs.Vote"`
}

type UserVotesPages struct {
	Users []UserVotesPage
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
	R       cipher.SHA256 `json:"-" enc:"-"` // Stores the content's hash for easier processing.
	Refer   cipher.SHA256 `json:"refer"`
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

type Deleted struct {
	Threads []cipher.SHA256
	Posts   []cipher.SHA256
}

/*
	<<< VOTES >>>
*/

type Vote struct {
	OfUser    cipher.PubKey `json:"-"`
	OfContent cipher.SHA256 `json:"-"`
	Mode      int8          `json:"mode" trans:"mode"`
	Tag       []byte        `json:"-" trans:"tag"`
	Created   int64         `json:"created" trans:"time"`
	Creator   cipher.PubKey `json:"-" verify:"upk" trans:"upk"`
	Sig       cipher.Sig    `json:"-" verify:"sig"`
}

func (v Vote) Verify() error { return tag.Verify(&v) }

type VotesSummary struct {
	sync.Mutex
	Index     int // Index in VotesPages.
	OfUser    cipher.PubKey
	OfContent cipher.SHA256
	Hash      cipher.SHA256 // VotesPages' hash.
	Votes     map[cipher.PubKey]Vote
	Up        int
	Down      int
}

func (s *VotesSummary) View(perspective cipher.PubKey) VotesSummaryView {
	s.Lock()
	defer s.Unlock()
	vote := s.Votes[perspective]
	return VotesSummaryView{
		Up: CompiledVotesView{
			Count: s.Up,
			Voted: vote.Mode == +1,
		},
		Down: CompiledVotesView{
			Count: s.Down,
			Voted: vote.Mode == -1,
		},
	}
}

type VotesSummaryView struct {
	Up   CompiledVotesView `json:"up"`
	Down CompiledVotesView `json:"down"`
}

type CompiledVotesView struct {
	Count int  `json:"count"`
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

/*
	<<< CONNECTION >>>
*/

type Connection struct {
	Address string `json:"address"`
	State   string `json:"state"`
}

/*
	<<< NEW >>>
*/

type VoteStore struct {
	Users []cipher.PubKey
	Votes skyobject.Refs `skyobject:"schema=bbs.Vote"`
}

// ThreadStore stores threads and associated objects.
// The fields "Threads", "ThreadVotes" and "ThreadPosts" are all of the same length.
// i.e. ThreadVotes[i] is the votes of thread Threads[i].
//      ThreadPosts[i] is the posts in thread Threads[i].
//
// BENEFITS:
// - As threads are in `skyobject.Refs`, so their hashes are directly indexed in Radix true by CXO.
// - Easily find associated Votes and Posts without creating a map.
// - Easy deletion of threads and associated content.
type ThreadStore struct {
	Threads     skyobject.Refs `skyobject:"schema=bbs.Thread"`
	ThreadVotes skyobject.Refs `skyobject:"schema=bbs.VoteStore"`
	ThreadPosts skyobject.Refs `skyobject:"schema=bbs.PostStore"`
}
