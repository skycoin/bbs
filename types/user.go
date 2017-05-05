package types

import (
	"errors"
	"github.com/skycoin/skycoin/src/cipher"
)

// User represents a user.
type User struct {
	Alias        string        `json:"alias"`
	Master       bool          `json:"master"`
	PublicKey    cipher.PubKey `json:"-"`
	SecretKey    cipher.SecKey `json:"-"`
	PublicKeyStr string        `json:"public_key"`
	SecretKeyStr string        `json:"private_key,omitempty"`
}

// NewUser creates a new user.
// If no secret key is provided, user is not master.
func NewUser(alias string, pk cipher.PubKey, sk ...cipher.SecKey) (*User, error) {
	if len(sk) > 1 {
		return nil, errors.New("invalid number of secret keys provided")
	}
	var user User
	user.Alias = alias
	user.Master = false

	// Check Public Key. If okay, add to user.
	if e := pk.Verify(); e != nil {
		return nil, e
	}
	user.PublicKey = pk
	user.PublicKeyStr = pk.Hex()

	// If secret key is provided, check it. If okay, add to user.
	if len(sk) == 1 {
		user.Master = true
		e := sk[0].Verify()
		if e != nil {
			return nil, e
		}
		user.SecretKey = sk[0]
		user.SecretKeyStr = sk[0].Hex()
	}
	return &user, nil
}

// PrepAndCheck prepares the user and checks whether it's valid.
func (u *User) PrepAndCheck() (e error) {
	// Prepare and check PublicKey.
	u.PublicKey, e = cipher.PubKeyFromHex(u.PublicKeyStr)
	// If master, prepare and check SecretKey.
	if u.Master {
		u.SecretKey, e = cipher.SecKeyFromHex(u.SecretKeyStr)
		// See if SecretKey generates expected PublicKey.
		pk := cipher.PubKeyFromSecKey(u.SecretKey)
		if pk != u.PublicKey {
			e = errors.New("secret key does not match public key")
		}
	}
	return
}
