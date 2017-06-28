package misc

import (
	"github.com/pkg/errors"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"math/rand"
	"strconv"
	"time"
)

// MakeTimeStampedRandomID makes a timestamped, random ID.
func MakeTimeStampedRandomID(n int) cipher.PubKey {
	s1 := []byte(strconv.FormatInt(time.Now().UnixNano(), 10))
	s2 := cipher.RandByte(n - len(s1))
	seed := append(s1, s2...)
	pk, _ := cipher.GenerateDeterministicKeyPair(seed)
	return pk
}

func MakeRandomAlias() string {
	out := "anonymous_"
	animals := []string{
		"cat",
		"bat",
		"bison",
		"dolphin",
		"eagle",
		"pony",
		"ape",
		"lobster",
		"monkey",
		"dog",
		"parrot",
		"cow",
		"sheep",
		"deer",
		"duck",
		"rabbit",
		"spider",
		"wolf",
		"turkey",
	}
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(animals) - 1)
	return out + animals[i]
}

func MakeDeterministicRef(seed byte) skyobject.Reference {
	shaOut := cipher.SHA256{}
	shaBytes := make([]byte, 32)
	for i := 0; i < len(shaBytes); i++ {
		shaBytes[i] = seed
	}
	shaOut.Set(shaBytes)
	return skyobject.Reference(shaOut)
}

func MakeIntBetween(min, max int) (int, error) {
	if min > max {
		return -1, errors.New("invalid range provided")
	}
	return min + rand.Intn(max-min+1), nil
}
