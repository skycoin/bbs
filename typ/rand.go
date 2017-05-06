package typ

import (
	"github.com/skycoin/skycoin/src/cipher"
	"math/rand"
	"strconv"
	"time"
)

// MakeTimeStampedRandomID makes a timestamped, random ID.
func MakeTimeStampedRandomID(n int) []byte {
	id := []byte(strconv.FormatInt(time.Now().UnixNano(), 10))
	id2 := cipher.RandByte(n - len(id))
	return append(id, id2...)
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
