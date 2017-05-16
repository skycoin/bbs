package misc

import (
	"encoding/hex"
	"errors"
	"github.com/skycoin/skycoin/src/cipher"
)

// GetPubKey obtains a public key from string, avoiding panics.
func GetPubKey(s string) (cipher.PubKey, error) {
	b, e := hex.DecodeString(s)
	if e != nil || len(b) != len(cipher.PubKey{}) {
		return cipher.PubKey{}, errors.New("invalid public key")
	}
	return cipher.NewPubKey(b), nil
}

// GetSecKey obtains a secret key from string, avoiding panics.
func GetSecKey(s string) (cipher.SecKey, error) {
	b, e := hex.DecodeString(s)
	if e != nil || len(b) != len(cipher.SecKey{}) {
		return cipher.SecKey{}, errors.New("invalid secret key")
	}
	return cipher.NewSecKey(b), nil
}
