package types

import (
	"errors"
	"fmt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
)

const UsersConfigFileName = "bbs_users.json"

// UsersConfig represents a configuration file for users.
type UsersConfig []*User

// UserManager managers users.
type UserManager struct {
	Masters []cipher.PubKey
	Users   map[cipher.PubKey]*User
}

// NewUserManager creates a new UserManager.
func NewUserManager() *UserManager {
	return &UserManager{
		Users: make(map[cipher.PubKey]*User),
	}
}

// Load loads data from config file.
func (m *UserManager) Load() error {
	fmt.Println(util.DataDir)
	// Load from file.
	uc := UsersConfig{}
	if e := util.LoadJSON(UsersConfigFileName, &uc); e != nil {
		return e
	}
	// Check loaded users. Append to m.Master, m.Users appropriately.
	for _, u := range uc {
		if e := u.PrepAndCheck(); e != nil {
			return e
		}
		m.Users[u.PublicKey] = u
		if u.Master {
			m.Masters = append(m.Masters, u.PublicKey)
		}
	}
	if len(m.Masters) == 0 {
		return errors.New("no masters")
	}
	return nil
}

// Save saves current in-memory config to file.
func (m *UserManager) Save() error {
	uc := make(UsersConfig, len(m.Users))
	i := 0
	for _, u := range m.Users {
		uc[i] = u
		i++
	}
	if e := util.SaveJSON(UsersConfigFileName, uc, 0600); e != nil {
		return e
	}
	return nil
}

// Clear clears the user configurations in memory.
func (m *UserManager) Clear() {
	m.Masters = []cipher.PubKey{}
	m.Users = make(map[cipher.PubKey]*User)
}

// AddNewRandomMaster adds a new random master.
func (m *UserManager) AddNewRandomMaster() error {
	pk, sk := cipher.GenerateKeyPair()
	u, e := NewUser(MakeRandomAlias(), pk, sk)
	if e != nil {
		return e
	}
	m.Users[pk] = u
	m.Masters = append(m.Masters, pk)
	return nil
}
