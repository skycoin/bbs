package state

import (
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// Instruction represents an instruction given to worker.
type Instruction struct {
	user    *cipher.PubKey      // Current User.
	data    []byte              // Raw data of single vote.
	summary *object.VoteSummary // Summary of votes.
}

func (i *Instruction) Run() error {

	var vote object.Vote
	if e := encoder.DeserializeRaw(i.data, &vote); e != nil {
		return e
	}
	if e := vote.Verify(); e != nil {
		return e
	}
	isUser := vote.User == *i.user
	switch vote.Mode {
	case -1:
		i.summary.Down.Lock()
		i.summary.Down.Count += 1
		if isUser {
			i.summary.Down.Voted = true
		}
		i.summary.Down.Unlock()
	case +1:
		i.summary.Up.Lock()
		i.summary.Up.Count += 1
		if isUser {
			i.summary.Up.Voted = true
		}
		i.summary.Up.Unlock()
	}
	switch string(vote.Tag) {
	case "spam":
		i.summary.Spam.Lock()
		i.summary.Spam.Count += 1
		if isUser {
			i.summary.Spam.Voted = true
		}
		i.summary.Spam.Unlock()
	}
	return nil
}
