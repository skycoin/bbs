package tag

import (
	"encoding/hex"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/skycoin/src/cipher"
	"strconv"
	"strings"
	"log"
)

// GetPubKey obtains the public key from string, avoiding panics.
func GetPubKey(s string) (cipher.PubKey, error) {
	b, e := hex.DecodeString(s)
	if e != nil {
		return cipher.PubKey{}, boo.WrapType(e, boo.InvalidInput,
			"invalid public key hex string")
	} else if len(b) != len(cipher.PubKey{}) {
		return cipher.PubKey{}, boo.New(boo.InvalidInput,
			"invalid public key hex string length")
	}
	pk := cipher.NewPubKey(b)
	if e := pk.Verify(); e != nil {
		return cipher.PubKey{}, boo.WrapType(e, boo.InvalidRead,
			"failed to verify public key")
	}
	return pk, nil
}

// GetSecKey obtains the secret key from string, avoiding panics.
func GetSecKey(s string) (cipher.SecKey, error) {
	b, e := hex.DecodeString(s)
	if e != nil {
		return cipher.SecKey{}, boo.WrapType(e, boo.InvalidInput,
			"invalid secret key hex string")
	} else if len(b) != len(cipher.SecKey{}) {
		return cipher.SecKey{}, boo.New(boo.InvalidInput,
			"invalid secret key hex string length")
	}
	sk := cipher.NewSecKey(b)
	if e := sk.Verify(); e != nil {
		return cipher.SecKey{}, boo.WrapType(e, boo.InvalidRead,
			"failed to verify secret key")
	}
	return cipher.NewSecKey(b), nil
}

// GetHash obtains a skyobject reference from hex string.
func GetHash(s string) (cipher.SHA256, error) {
	h, e := cipher.SHA256FromHex(s)
	if e != nil {
		return cipher.SHA256{}, boo.WrapType(e, boo.InvalidInput,
			"invalid reference")
	}
	return h, e
}

// GetSig obtains a signature from hex string.
func GetSig(s string) (cipher.Sig, error) {
	sig, e := cipher.SigFromHex(s)
	if e != nil {
		return cipher.Sig{}, boo.WrapType(e, boo.InvalidInput,
			"invalid sig string")
	}
	return sig, nil
}

func GetVoteValue(v string) (int8, error) {
	value, e := strconv.Atoi(v)
	if e != nil {
		return 0, boo.WrapType(e, boo.InvalidInput, "invalid vote value")
	}
	switch value {
	case -1, 0, +1:
	default:
		return 0, boo.New(boo.InvalidInput, "invalid vote value")
	}
	return int8(value), nil
}

func GetTags(v string) ([]string, error) {
	tags := strings.Split(v, ",")
	for i := len(tags) - 1; i >= 0; i-- {
		tags[i] = strings.TrimSpace(tags[i])
		if tags[i] == "" {
			tags = append(tags[:i], tags[i+1:]...)
		}
	}
	log.Println("GetTags(v string) got:", tags)
	return tags, nil
}
