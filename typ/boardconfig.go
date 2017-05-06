package typ

import (
	"errors"
	"github.com/skycoin/skycoin/src/cipher"
)

// BoardConfig represents a board configuration on disk and in memory.
type BoardConfig struct {
	Name      string        `json:"name,omitempty"`
	Master    bool          `json:"master"`
	URL       string        `json:"url"`
	PubKey    cipher.PubKey `json:"-"`
	SecKey    cipher.SecKey `json:"-"`
	PubKeyStr string        `json:"public_key"`
	SecKeyStr string        `json:"secret_key,omitempty"`
}

// NewBoardConfig creates a new non-master board config from public key.
func NewBoardConfig(pk cipher.PubKey, url string) (*BoardConfig, error) {
	if e := pk.Verify(); e != nil {
		return nil, e
	}
	bc := BoardConfig{
		Master:    false,
		URL:       url,
		PubKey:    pk,
		PubKeyStr: pk.Hex(),
	}
	return &bc, nil
}

// NewMasterBoardConfig creates a new master BoardConfig from a seed.
func NewMasterBoardConfig(name, seed, url string) *BoardConfig {
	pk, sk := cipher.GenerateDeterministicKeyPair([]byte(seed))
	bc := BoardConfig{
		Name:      name,
		Master:    true,
		URL:       url,
		PubKey:    pk,
		SecKey:    sk,
		PubKeyStr: pk.Hex(),
		SecKeyStr: sk.Hex(),
	}
	return &bc
}

// PrepAndCheck prepares the BoardConfig and checks whether it's valid.
func (c *BoardConfig) PrepAndCheck() (e error) {
	// Prepare and check PublicKey.
	c.PubKey, e = cipher.PubKeyFromHex(c.PubKeyStr)
	// If master, prepare and check SecretKey.
	if c.Master {
		c.SecKey, e = cipher.SecKeyFromHex(c.SecKeyStr)
		// See if SecretKey generates expected PublicKey.
		pk := cipher.PubKeyFromSecKey(c.SecKey)
		if pk != c.PubKey {
			e = errors.New("secret key does not match public key")
		}
	}
	return
}
