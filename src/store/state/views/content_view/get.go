package content_view

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/skycoin/src/cipher"
)

const (
	Board      = "Board"
	BoardPage  = "BoardPage"
	ThreadPage = "ThreadPage"
)

func (v *ContentView) Get(id string, a ...interface{}) (interface{}, error) {
	v.Lock()
	defer v.Unlock()

	switch {
	case id == Board:
		return v.getBoard()

	case id == BoardPage:
		return v.getBoardPage()

	case id == ThreadPage && len(a) == 1:
		return v.getThreadPage(a[0].(cipher.SHA256))

	default:
		return nil, boo.Newf(boo.NotAllowed,
			"invalid get request 's' (%v)", id, a)
	}
}

func (v *ContentView) getBoard() (*BoardRep, error) {
	return v.board, nil
}

type BoardPageOut struct {
	Board   *BoardRep    `json:"board"`
	Threads []*ThreadRep `json:"threads"`
}

func (v *ContentView) getBoardPage() (*BoardPageOut, error) {
	out := new(BoardPageOut)
	out.Board = v.board
	out.Threads = make([]*ThreadRep, len(v.board.Threads))
	for i, tHash := range v.board.Threads {
		out.Threads[i] = v.tMap[tHash]
	}
	return out, nil
}

type ThreadPageOut struct {
	Board  *BoardRep  `json:"board"`
	Thread *ThreadRep `json:"thread"`
	Posts  []*PostRep `json:"posts"`
}

func (v *ContentView) getThreadPage(threadHash cipher.SHA256) (*ThreadPageOut, error) {
	out := new(ThreadPageOut)
	out.Board = v.board
	out.Thread = v.tMap[threadHash]
	if out.Thread != nil {
		out.Posts = make([]*PostRep, len(out.Thread.Posts))
		for i, pHash := range out.Thread.Posts {
			out.Posts[i] = v.pMap[pHash]
		}
	}
	return out, nil
}
