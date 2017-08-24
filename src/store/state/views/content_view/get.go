package content_view

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/skycoin/src/cipher"
)

const (
	Board        = "Board"
	BoardPage    = "BoardPage"
	ThreadPage   = "ThreadPage"
	SubAddresses = "SubAddresses"
	ContentVotes = "ContentVotes"
)

func (v *ContentView) Get(id string, a ...interface{}) (interface{}, error) {
	v.Lock()
	defer v.Unlock()

	switch {
	case id == Board:
		return v.getBoard()

	case id == BoardPage && len(a) == 1:
		return v.getBoardPage(a[0].(cipher.PubKey))

	case id == ThreadPage && len(a) == 2:
		return v.getThreadPage(a[0].(cipher.PubKey), a[1].(cipher.SHA256))

	case id == SubAddresses:
		return v.getSubAddresses()

	case id == ContentVotes && len(a) == 2:
		return v.getVotes(a[0].(cipher.PubKey), a[1].(cipher.SHA256))

	default:
		return nil, boo.Newf(boo.NotAllowed,
			"invalid get request 's' (%v)", id, a)
	}
}

func (v *ContentView) getBoard() (*BoardRepView, error) {
	return v.board.View(), nil
}

type BoardPageOut struct {
	Board   *BoardRepView    `json:"board"`
	Threads []*ThreadRepView `json:"threads"`
}

func (v *ContentView) getBoardPage(perspective cipher.PubKey) (*BoardPageOut, error) {
	out := new(BoardPageOut)
	out.Board = v.board.View()
	out.Threads = make([]*ThreadRepView, len(v.board.Threads))
	for i, tHash := range v.board.Threads {
		out.Threads[i] = v.tMap[tHash].View(v.vMap[tHash].View(perspective))
	}
	return out, nil
}

type ThreadPageOut struct {
	Board  *BoardRepView  `json:"board"`
	Thread *ThreadRepView `json:"thread"`
	Posts  []*PostRepView `json:"posts"`
}

func (v *ContentView) getThreadPage(perspective cipher.PubKey, threadHash cipher.SHA256) (*ThreadPageOut, error) {
	out := new(ThreadPageOut)
	out.Board = v.board.View()

	threadRep := v.tMap[threadHash]

	if threadRep != nil {
		out.Thread = threadRep.View(v.vMap[threadHash].View(perspective))

		out.Posts = make([]*PostRepView, len(threadRep.Posts))
		for i, pHash := range threadRep.Posts {
			out.Posts[i] = v.pMap[pHash].View(v.vMap[pHash].View(perspective))
		}
	}

	return out, nil
}

func (v *ContentView) getSubAddresses() ([]string, error) {
	sa := v.board.SubAddresses
	if len(sa) == 0 {
		return nil, boo.Newf(boo.NotFound,
			"board of public key '%s' has no submission addresses",
			v.board.PubKey)
	}
	return sa, nil
}

type ContentVotesOut struct {
	Votes *VoteRepView `json:"votes"`
}

func (v *ContentView) getVotes(perspective cipher.PubKey, cHash cipher.SHA256) (*ContentVotesOut, error) {
	out := new(ContentVotesOut)
	out.Votes = v.vMap[cHash].View(perspective)
	return out, nil
}
