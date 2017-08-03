package session

import (
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/skycoin/src/cipher"
	"math/rand"
	"testing"
	"time"
)

func generateKeyPairs(n int) []object.Subscription {
	out := make([]object.Subscription, n)
	for i := 0; i < n; i++ {
		pk, sk := cipher.GenerateKeyPair()
		out[i] = object.Subscription{
			PubKey: pk,
			SecKey: sk,
		}
	}
	return out
}

func TestUserFile_FillMaster(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	const mCount = 10
	var mPick = rand.Intn(mCount)

	board := &struct {
		PK cipher.PubKey `bbs:"bpk"`
		SK cipher.SecKey `bbs:"bsk"`
	}{}
	file := &File{
		Masters: generateKeyPairs(mCount),
	}
	board.PK = file.Masters[mPick].PubKey
	file.FillMaster(board)

	if board.PK == file.Masters[mPick].PubKey &&
		board.SK == file.Masters[mPick].SecKey {
		t.Logf("Success! PK(%s) SK(%s)",
			board.PK.Hex(), board.SK.Hex())
	} else {
		t.Errorf("Failed! PK(%s) SK(%s)",
			board.PK.Hex(), board.SK.Hex())
	}
}

func TestUserFile_FillUser(t *testing.T) {
	user := &struct {
		PK cipher.PubKey `bbs:"upk"`
		SK cipher.SecKey `bbs:"usk"`
	}{}
	pk, sk := cipher.GenerateKeyPair()
	file := &File{
		User: object.User{
			PublicKey: pk,
			SecretKey: sk,
		},
	}
	file.FillUser(user)
	if user.PK == pk && user.SK == sk {
		t.Logf("Success! PK(%s) SK(%s)", pk.Hex(), sk.Hex())
	} else {
		t.Error("Failed!")
	}
}
