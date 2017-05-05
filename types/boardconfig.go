package types

import (
	"errors"
	"github.com/skycoin/skycoin/src/cipher"
)

// BoardConfig represents a board's configuration as stored on a local file.
type BoardConfig struct {
	Name         string        `json:"name,omitempty"`
	Master       bool          `json:"master"`
	URL          string        `json:"url"`
	PublicKey    cipher.PubKey `json:"-"`
	SecretKey    cipher.SecKey `json:"-"`
	PublicKeyStr string        `json:"public_key"`
	SecretKeyStr string        `json:"secret_key,omitempty"` // Empty if Master = false.
}

// NewBoardConfig creates a new non-master board config from public key.
func NewBoardConfig(pk cipher.PubKey, url string) (*BoardConfig, error) {
	if e := pk.Verify(); e != nil {
		return nil, e
	}
	bc := BoardConfig{
		Master:       false,
		URL:          url,
		PublicKey:    pk,
		PublicKeyStr: pk.Hex(),
	}
	return &bc, nil
}

// NewMasterBoardConfig creates a new master BoardConfig from a seed.
func NewMasterBoardConfig(name, seed, url string) *BoardConfig {
	pk, sk := cipher.GenerateDeterministicKeyPair([]byte(seed))
	bc := BoardConfig{
		Name:         name,
		Master:       true,
		URL:          url,
		PublicKey:    pk,
		SecretKey:    sk,
		PublicKeyStr: pk.Hex(),
		SecretKeyStr: sk.Hex(),
	}
	return &bc
}

// PrepAndCheck prepares the BoardConfig and checks whether it's valid.
func (c *BoardConfig) PrepAndCheck() (e error) {
	// Prepare and check PublicKey.
	c.PublicKey, e = cipher.PubKeyFromHex(c.PublicKeyStr)
	// If master, prepare and check SecretKey.
	if c.Master {
		c.SecretKey, e = cipher.SecKeyFromHex(c.SecretKeyStr)
		// See if SecretKey generates expected PublicKey.
		pk := cipher.PubKeyFromSecKey(c.SecretKey)
		if pk != c.PublicKey {
			e = errors.New("secret key does not match public key")
		}
	}
	return
}
