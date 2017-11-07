package state

import (
	"encoding/json"
	"github.com/skycoin/bbs/src/store/object"
)

type VotesRep struct {
	Ref  string
	Type object.ContentType

	Votes     map[string]*object.Content // Key: pk string, Value: vote.
	UpCount   int
	DownCount int
}

func (r *VotesRep) String() string {
	raw, _ := json.MarshalIndent(r, "", "    ")
	return string(raw)
}

func (r *VotesRep) GetValue(c *object.Content) int {
	var value int
	switch r.Type {
	case object.V5ThreadVoteType:
		value = c.GetBody().Value
	case object.V5PostVoteType:
		value = c.GetBody().Value
	case object.V5UserVoteType:
		value = c.GetBody().Value
	}
	return value
}

func (r *VotesRep) Fill(refType object.ContentType, refHash string) *VotesRep {
	r.Ref = refHash
	r.Type = refType
	r.Votes = make(map[string]*object.Content)
	return r
}

func (r *VotesRep) Add(c *object.Content) {
	creator := c.GetBody().Creator
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
