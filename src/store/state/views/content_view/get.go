package content_view

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"log"
	"github.com/skycoin/bbs/src/store/object"
)

const (
	Board        = "Board"
	BoardPage    = "BoardPage"
	ThreadPage   = "ThreadPage"
	ContentVotes = "ContentVotes"
)

func (v *ContentView) Get(id string, a ...interface{}) (interface{}, error) {
	v.Lock()
	defer v.Unlock()

	switch {
	case id == Board:
		return v.getBoard()

	case id == BoardPage && len(a) == 1:
		return v.getBoardPage(a[0].(*BoardPageIn))

	case id == ThreadPage && len(a) == 1:
		return v.getThreadPage(a[0].(*ThreadPageIn))

	case id == ContentVotes && len(a) == 1:
		return v.getVotes(a[0].(*ContentVotesIn))

	default:
		return nil, boo.Newf(boo.NotAllowed,
			"invalid get request 's' (%v)", id, a)
	}
}

func (v *ContentView) getBoard() (*object.ContentRep, error) {
	return v.c[v.i.Board], nil
}

type BoardPageIn struct {
	Perspective string
}

type BoardPageOut struct {
	Board   *object.ContentRep   `json:"board"`
	Threads []*object.ContentRep `json:"threads"`
}

func (v *ContentView) getBoardPage(in *BoardPageIn) (*BoardPageOut, error) {
	out := new(BoardPageOut)
	out.Board = v.c[v.i.Board]
	out.Threads = make([]*object.ContentRep, len(v.i.Threads))
	for i, tHash := range v.i.Threads {
		out.Threads[i] = v.c[tHash]
		if votes, ok := v.v[tHash]; ok {
			out.Threads[i].Votes = votes.View(in.Perspective)
		}
	}
	return out, nil
}

type ThreadPageIn struct {
	Perspective string
	ThreadHash  string
}

type ThreadPageOut struct {
	Board  *object.ContentRep   `json:"board"`
	Thread *object.ContentRep   `json:"thread"`
	Posts  []*object.ContentRep `json:"posts"`
}

func (v *ContentView) getThreadPage(in *ThreadPageIn) (*ThreadPageOut, error) {
	out := new(ThreadPageOut)
	out.Board = v.c[v.i.Board]
	out.Thread = v.c[in.ThreadHash]

	if out.Thread == nil {
		return nil, boo.Newf(boo.NotFound,
			"thread of hash '%s' is not found in board '%s'",
			in.ThreadHash, v.pk.Hex())
	}
	if votes, ok := v.v[in.ThreadHash]; ok {
		out.Thread.Votes = votes.View(in.Perspective)
	}

	postHashes := v.i.Posts[in.ThreadHash]
	out.Posts = make([]*object.ContentRep, len(postHashes))
	for i, pHash := range postHashes {
		out.Posts[i] = v.c[pHash]
		if votes, ok := v.v[pHash]; ok {
			out.Posts[i].Votes = votes.View(in.Perspective)
		}
	}

	return out, nil
}

type ContentVotesIn struct {
	Perspective string
	ContentHash string
}

type ContentVotesOut struct {
	Votes *VoteRepView `json:"votes"`
}

func (v *ContentView) getVotes(in *ContentVotesIn) (*ContentVotesOut, error) {
	out := new(ContentVotesOut)

	if votes, ok := v.v[in.ContentHash]; ok {
		out.Votes = votes.View(in.Perspective)
		return out, nil
	}

	if _, ok := v.c[in.ContentHash]; ok {
		out.Votes = &VoteRepView{
			Ref: in.ContentHash,
		}
		return out, nil
	}

	for i, hash := range v.c {
		log.Printf("[%d] %s", i, hash)
	}

	return nil, boo.Newf(boo.NotFound,
		"content of hash '%s' is not found", in.ContentHash)
}
