package v1

import (
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// Instruction represents an instruction given to worker.
type Instruction struct {
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

	i.summary.Lock()
	defer i.summary.Unlock()

	i.summary.Votes[vote.User] = vote

	switch vote.Mode {
	case -1:
		i.summary.Downs += 1
	case +1:
		i.summary.Ups += 1
	}

	switch string(vote.Tag) {
	case "spam":
		i.summary.Spams += 1
	}
	return nil
}
