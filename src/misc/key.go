package misc

import (
	"encoding/hex"
	"github.com/skycoin/bbs/src/boo"
	"github.com/skycoin/skycoin/src/cipher"
)

// GetPubKey obtains the public key from string, avoiding panics.
func GetPubKey(s string) (cipher.PubKey, error) {
	b, e := hex.DecodeString(s)
	if e != nil {
		return cipher.PubKey{}, boo.New(boo.InvalidInput,
			"invalid public key hex string:", e.Error())
	} else if len(b) != len(cipher.PubKey{}) {
		return cipher.PubKey{}, boo.New(boo.InvalidInput,
			"invalid public key hex string length")
	}
	return cipher.NewPubKey(b), nil
}

// GetSecKey obtains the secret key from string, avoiding panics.
func GetSecKey(s string) (cipher.SecKey, error) {
	b, e := hex.DecodeString(s)
	if e != nil {
		return cipher.SecKey{}, boo.New(boo.InvalidInput,
			"invalid secret key hex string:", e.Error())
	} else if len(b) != len(cipher.SecKey{}) {
		return cipher.SecKey{}, boo.New(boo.InvalidInput,
			"invalid secret key hex string length")
	}
	return cipher.NewSecKey(b), nil
}
