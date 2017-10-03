package content_view

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object/revisions/r0"
	"fmt"
	"encoding/json"
	"log"
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

func (v *ContentView) getBoard() (*r0.ContentRep, error) {
	return v.c[v.i.Board], nil
}

type BoardPageIn struct {
	Perspective string
}

type BoardPageOut struct {
	Board   *r0.ContentRep   `json:"board"`
	Threads []*r0.ContentRep `json:"threads"`
}

func (v *ContentView) getBoardPage(in *BoardPageIn) (*BoardPageOut, error) {
	out := new(BoardPageOut)
	out.Board = v.c[v.i.Board]
	out.Threads = make([]*r0.ContentRep, len(v.i.Threads))
	for i, tHash := range v.i.Threads {
		out.Threads[i] = v.c[tHash]
		out.Threads[i].Votes = v.v[tHash].View(in.Perspective)
	}
	return out, nil
}

type ThreadPageIn struct {
	Perspective string
	ThreadHash  string
}

type ThreadPageOut struct {
	Board  *r0.ContentRep   `json:"board"`
	Thread *r0.ContentRep   `json:"thread"`
	Posts  []*r0.ContentRep `json:"posts"`
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
	out.Thread.Votes = v.v[in.ThreadHash].View(in.Perspective)

	postHashes := v.i.Posts[in.ThreadHash]
	out.Posts = make([]*r0.ContentRep, len(postHashes))
	for i, pHash := range postHashes {
		out.Posts[i] = v.c[pHash]
		out.Posts[i].Votes = v.v[pHash].View(in.Perspective)
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
	raw, _ := json.MarshalIndent(in, "", "    ")
	fmt.Println("<<< GET VOTES : INPUT >>>", string(raw))

	if _, ok := v.v[in.ContentHash]; !ok {
		log.Println("DOES NOT HAVE VOTES FOR CONTENT:", in.ContentHash)
	}

	out.Votes = v.v[in.ContentHash].View(in.Perspective)
	return out, nil
}
