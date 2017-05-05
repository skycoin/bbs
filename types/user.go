package types

import (
	"encoding/json"
	"github.com/skycoin/skycoin/src/cipher"
	"io/ioutil"
	"errors"
)

// User represents a user.
type User struct {
	Alias         string        `json:"alias"`
	PublicKey     cipher.PubKey `json:"-"`
	SecretKey     cipher.SecKey `json:"-"`
	PublicKeyStr  string        `json:"public_key"`
	PrivateKeyStr string        `json:"private_key,omitempty"`
}

// NewUser creates a new user.
func NewUser(alias string, pk cipher.PubKey, sk ...cipher.SecKey) (*User, error) {
	if len(sk) > 1 {
		return nil, errors.New("invalid number of secret keys provided")
	}
	var user User
	user.Alias = alias

	// Check Public Key. If okay, add to user.
	if e := pk.Verify(); e != nil {
		return nil, e
	}
	user.PublicKey = pk
	user.PublicKeyStr = pk.Hex()

	// If secret key is provided, check it. If okay, add to user.
	if len(sk) == 1 {
		e := sk[0].Verify()
		if e != nil {
			return nil, e
		}
		user.SecretKey = sk[0]
		user.PublicKeyStr = sk[0].Hex()
	}
	return &user, nil
}



// UserConfigFile represents a user configuration file.
type UserConfigFile struct {
	Master *User   `json:"master"` // Master user.
	Others []*User `json:"others"` // Other users.
}

// UserConfig represents a user configuration.
type UserConfig struct {
	Alias     string
	PublicKey cipher.PubKey
	SecretKey cipher.SecKey
}

type UserConfigRaw struct {
	Alias     string `json:"alias"`
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
