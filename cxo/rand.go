package cxo

import (
	"github.com/skycoin/skycoin/src/cipher"
	"strconv"
	"time"
)

// MakeTimeStampedRandomID makes a timestamped, random ID.
func MakeTimeStampedRandomID(n int) []byte {
	id := []byte(strconv.FormatInt(time.Now().UnixNano(), 10))
	id2 := cipher.RandByte(n - len(id))
	return append(id, id2...)
}
