package types

import (
	"encoding/json"
	"github.com/skycoin/skycoin/src/cipher"
	"io/ioutil"
)

// UserConfig represents a user configuration.
type UserConfig struct {
	PublicKey cipher.PubKey
	SecretKey cipher.SecKey
}

type UserConfigRaw struct {
	PublicKey string `json:"public_key"`
	SecretKey string `json:"secret_key"`
}

// NewUserConfigFromData obtains a new UserConfig from data.
func NewUserConfigFromData(data []byte) (*UserConfig, error) {
	// Read json file to 'rawConfig'.
	rawConfig := UserConfigRaw{}
	if e := json.Unmarshal(data, &rawConfig); e != nil {
		return nil, e
	}

	// Create UserConfig from UserConfigRaw.
	var e error
	userConfig := UserConfig{}
	userConfig.PublicKey, e = cipher.PubKeyFromHex(rawConfig.PublicKey)
	if e != nil {
		return nil, e
	}
	userConfig.SecretKey, e = cipher.SecKeyFromHex(rawConfig.SecretKey)
	if e != nil {
		return nil, e
	}

	// Check keys.
	if e := userConfig.CheckKeys(); e != nil {
		return nil, e
	}
	return &userConfig, nil
}

// NewUserConfigFromFile obtains a new UserConfig from a file.
func NewUserConfigFromFile(filename string) (*UserConfig, error) {
	data, e := ioutil.ReadFile(filename)
	if e != nil {
		return nil, e
	}
	return NewUserConfigFromData(data)
}

func (c *UserConfig) CheckKeys() error {
	if e := c.PublicKey.Verify(); e != nil {
		return e
	}
	if e := c.SecretKey.Verify(); e != nil {
		return e
	}
	return nil
}
