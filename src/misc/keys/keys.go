package keys

import (
	"encoding/hex"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/skycoin/src/cipher"
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

func PubKeyToSlice(pk cipher.PubKey) []byte {
	out := make([]byte, 33)
	for i, v := range [33]byte(pk) {
		out[i] = v
	}
	return out
}