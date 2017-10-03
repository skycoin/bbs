package content_view

import (
	"github.com/skycoin/bbs/src/store/object/revisions/r0"
	"encoding/json"
)

/*
	<<< VOTES >>>
*/

type VotesRep struct {
	Ref  string
	Type r0.ContentType

	Votes     map[string]*r0.Content // Key: pk string, Value: vote.
	UpCount   int
	DownCount int
}

func (r *VotesRep) String() string {
	raw, _ := json.MarshalIndent(r, "", "    ")
	return string(raw)
}

func (r *VotesRep) GetValue(c *r0.Content) int {
	var value int
	switch r.Type {
	case r0.V5ThreadVoteType:
		value = c.ToThreadVote().GetBody().Value
	case r0.V5PostVoteType:
		value = c.ToPostVote().GetBody().Value
	case r0.V5UserVoteType:
		value = c.ToUserVote().GetBody().Value
	}
	return value
}

func (r *VotesRep) Fill(refType r0.ContentType, refHash string) *VotesRep {
	r.Ref = refHash
	r.Type = refType
	r.Votes = make(map[string]*r0.Content)
	return r
}

func (r *VotesRep) Add(c *r0.Content) {
	creator := c.GetHeader().PK
	if oldC, has := r.Votes[creator]; has {
		switch r.GetValue(oldC) {
		case +1:
			r.UpCount--
		case -1:
			r.DownCount--
		}
	}
	r.Votes[creator] = c

	switch r.GetValue(c) {
	case +1:
		r.UpCount++
	case -1:
		r.DownCount++
	case 0:
		delete(r.Votes, creator)
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

func (r *VotesRep) View(user string) *VoteRepView {
	if r == nil {
		return nil
	}
	c := r.Votes[user]
	return &VoteRepView{
		Ref: r.Ref,
		Up: X{
			Voted: c != nil && r.GetValue(c) == +1,
			Count: r.UpCount,
		},
		Down: X{
			Voted: c != nil && r.GetValue(c) == -1,
			Count: r.DownCount,
		},
	}
}
