package content_view

import (
	"github.com/skycoin/bbs/src/store/object/revisions/r0"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
)

/*
	<<< INDEXED HASH >>>
*/

type IndexHash struct {
	h cipher.SHA256
	i int
}

/*
	<<< BOARD >>>
*/

type BoardRep struct {
	PubKey     cipher.PubKey
	Name       string
	Body       string
	Created    int64
	Tags       []string
	SubPubKeys []string
	Threads    []IndexHash
}

func (r *BoardRep) Fill(pk cipher.PubKey, board *r0.Board) *BoardRep {
	data := board.GetData()
	r.PubKey = pk
	r.Name = data.Name
	r.Body = data.Body
	r.Created = data.Created
	r.Tags = data.Tags
	r.SubPubKeys = data.SubKeys
	return r
}

type BoardRepView struct {
	PubKey      string   `json:"public_key"`
	Name        string   `json:"name"`
	Body        string   `json:"body"`
	Created     int64    `json:"created"`
	Tags        []string `json:"tags"`
	ThreadCount int      `json:"thread_count"`
}

func (r *BoardRep) View() *BoardRepView {
	if r == nil {
		return nil
	}
	return &BoardRepView{
		PubKey:      r.PubKey.Hex(),
		Name:        r.Name,
		Body:        r.Body,
		Created:     r.Created,
		Tags:        r.Tags,
		ThreadCount: len(r.Threads),
	}
}

/*
	<<< THREAD >>>
*/

type ThreadRep struct {
	Ref     cipher.SHA256
	Name    string
	Body    string
	Created int64
	Creator cipher.PubKey
	Posts   []IndexHash
}

func (r *ThreadRep) FillThread(thread *r0.Thread) *ThreadRep {
	data := thread.GetData()
	r.Ref = thread.R
	r.Name = data.Name
	r.Body = data.Body
	r.Created = data.Created
	r.Creator = data.GetCreator()
	return r
}

func (r *ThreadRep) FillThreadPage(tPage *r0.ThreadPage) *ThreadRep {
	t, e := tPage.GetThread()
	if e != nil {
		log.Println("ThreadRep.FillThreadPage() Error:", e)
		return nil
	}
	data := t.GetData()
	r.Ref = t.R
	r.Name = data.Name
	r.Body = data.Body
	r.Created = data.Created
	r.Creator = data.GetCreator()
	return r
}

type ThreadRepView struct {
	Seq       int          `json:"seq"`
	Ref       string       `json:"ref"`
	Name      string       `json:"name"`
	Body      string       `json:"body"`
	Created   int64        `json:"created"`
	Creator   string       `json:"creator"`
	Votes     *VoteRepView `json:"votes,omitempty"`
	PostCount int          `json:"post_count"`
}

func (r *ThreadRep) View(i int, votes *VoteRepView) *ThreadRepView {
	if r == nil {
		return nil
	}
	return &ThreadRepView{
		Seq:       i,
		Ref:       r.Ref.Hex(),
		Name:      r.Name,
		Body:      r.Body,
		Created:   r.Created,
		Creator:   r.Creator.Hex(),
		Votes:     votes,
		PostCount: len(r.Posts),
	}
}

/*
	<<< POST >>>
*/

type PostRep struct {
	Ref     cipher.SHA256
	Name    string
	Body    string
	Images  []*r0.ImageData
	Created int64
	Creator cipher.PubKey
}

func (r *PostRep) Fill(post *r0.Post) *PostRep {
	data := post.GetData()
	r.Ref = post.R
	r.Name = data.Name
	r.Body = data.Body
	r.Images = data.Images
	r.Created = data.Created
	r.Creator = data.GetCreator()
	return r
}

type PostRepView struct {
	Seq     int             `json:"seq"`
	Ref     string          `json:"ref"`
	Type    string          `json:"type"`
	Name    string          `json:"name"`
	Body    string          `json:"body"`
	Images  []*r0.ImageData `json:"images,omitempty"`
	Created int64           `json:"created"`
	Creator string          `json:"creator"`
	Votes   *VoteRepView    `json:"votes,omitempty"`
}

func (r *PostRep) View(i int, votes *VoteRepView) *PostRepView {
	if r == nil {
		return nil
	}
	return &PostRepView{
		Seq:     i,
		Ref:     r.Ref.Hex(),
		Name:    r.Name,
		Body:    r.Body,
		Images:  r.Images,
		Created: r.Created,
		Creator: r.Creator.Hex(),
		Votes:   votes,
	}
}

/*
	<<< VOTES >>>
*/

type VotesRep struct {
	Ref       cipher.SHA256
	Votes     map[cipher.PubKey]*r0.Vote
	UpCount   int
	DownCount int
}

func (r *VotesRep) Fill(hash cipher.SHA256) *VotesRep {
	r.Ref = hash
	r.Votes = make(map[cipher.PubKey]*r0.Vote)
	return r
}

func (r *VotesRep) Add(vote *r0.Vote) {
	if oldVote, has := r.Votes[vote.Creator]; has {
		switch oldVote.Mode {
		case +1:
			r.UpCount--
		case -1:
			r.DownCount--
		}
	}
	r.Votes[vote.Creator] = vote
	switch vote.Mode {
	case +1:
		r.UpCount++
	case -1:
		r.DownCount++
	case 0:
		delete(r.Votes, vote.Creator)
	}
}

type X struct {
	Voted bool `json:"voted"`
	Count int  `json:"count"`
}

type VoteRepView struct {
	Ref  string `json:"ref"`
	Up   X      `json:"up_votes"`
	Down X      `json:"down_votes"`
}

func (r *VotesRep) View(perspective cipher.PubKey) *VoteRepView {
	if r == nil {
		return nil
	}
	vote := r.Votes[perspective]
	return &VoteRepView{
		Ref: r.Ref.Hex(),
		Up: X{
			Voted: vote != nil && vote.Mode == +1,
			Count: r.UpCount,
		},
		Down: X{
			Voted: vote != nil && vote.Mode == -1,
			Count: r.DownCount,
		},
	}
}
