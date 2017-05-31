package typ

import (
	"github.com/skycoin/skycoin/src/cipher"
	"strconv"
	"testing"
)

func runVote(i int) error {
	vote := Vote{}
	if i%2 == 0 {
		vote.Up()
	} else {
		vote.Down()
	}
	if e := vote.Sign(cipher.GenerateDeterministicKeyPair([]byte(strconv.Itoa(i)))); e != nil {
		return e
	}
	return vote.Verify()
}

func TestVote_Sign(t *testing.T) {
	for i := 0; i < 100; i++ {
		if e := runVote(i); e != nil {
			t.Error(e)
		}
	}
}
