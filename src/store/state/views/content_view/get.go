package content_view

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/skycoin/src/cipher"
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

func (v *ContentView) getBoard() (*BoardRepView, error) {
	return v.board.View(), nil
}

type BoardPageIn struct {
	Perspective cipher.PubKey
}

type BoardPageOut struct {
	Board   *BoardRepView    `json:"board"`
	Threads []*ThreadRepView `json:"threads"`
}

func (v *ContentView) getBoardPage(in *BoardPageIn) (*BoardPageOut, error) {
	out := new(BoardPageOut)
	out.Board = v.board.View()
	out.Threads = make([]*ThreadRepView, len(v.board.Threads))
	for i, tHash := range v.board.Threads {
		out.Threads[i] = v.tMap[tHash.h].View(tHash.i, v.vMap[tHash.h].View(in.Perspective))
	}
	return out, nil
}

type ThreadPageIn struct {
	Perspective cipher.PubKey
	ThreadHash  cipher.SHA256
}

type ThreadPageOut struct {
	Board  *BoardRepView  `json:"board"`
	Thread *ThreadRepView `json:"thread"`
	Posts  []*PostRepView `json:"posts"`
}

func (v *ContentView) getThreadPage(in *ThreadPageIn) (*ThreadPageOut, error) {
	out := new(ThreadPageOut)
	out.Board = v.board.View()

	threadRep := v.tMap[in.ThreadHash]

	if threadRep != nil {
		out.Thread = threadRep.View(0, v.vMap[in.ThreadHash].View(in.Perspective))

		out.Posts = make([]*PostRepView, len(threadRep.Posts))
		for i, pHash := range threadRep.Posts {
			out.Posts[i] = v.pMap[pHash.h].View(pHash.i, v.vMap[pHash.h].View(in.Perspective))
		}
	}

	return out, nil
}

type ContentVotesIn struct {
	Perspective cipher.PubKey
	ContentHash cipher.SHA256
}

type ContentVotesOut struct {
	Votes *VoteRepView `json:"votes"`
}

func (v *ContentView) getVotes(in *ContentVotesIn) (*ContentVotesOut, error) {
	out := new(ContentVotesOut)
	out.Votes = v.vMap[in.ContentHash].View(in.Perspective)
	return out, nil
}
