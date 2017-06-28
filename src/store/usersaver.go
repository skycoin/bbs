package store

import (
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/cmd/bbsnode/args"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/bbs/src/store/cxo"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
	"log"
	"os"
	"path/filepath"
	"sync"
)

// UserSaverFileName represents the filename of the users configuration file.
const UserSaverFileName = "bbs_users.json"

// UsersConfigFile represents the layout of a configuration file of boards.
type UsersConfigFile struct {
	Current string        `json:"current"`
	Users   []*UserConfig `json:"users"`
}

// UserConfig represents the config of a user.
type UserConfig struct {
	Alias  string `json:"alias,omitempty"`
	Master bool   `json:"master"`
	PubKey string `json:"public_key"`
	SecKey string `json:"secret_key,omitempty"`
}

func (uc *UserConfig) GetPK() cipher.PubKey {
	pk, _ := misc.GetPubKey(uc.PubKey)
	return pk
}

func (uc *UserConfig) GetSK() cipher.SecKey {
	sk, _ := misc.GetSecKey(uc.SecKey)
	return sk
}

// Check checks the validity of the UserConfig.
func (uc *UserConfig) Check() (cipher.PubKey, error) {
	pk, e := misc.GetPubKey(uc.PubKey)
	if e != nil {
		return pk, e
	}
	if uc.Master {
		sk, e := misc.GetSecKey(uc.SecKey)
		if e != nil {
			return pk, e
		}
		if pk != cipher.PubKeyFromSecKey(sk) {
			return pk, errors.New("invalid public-secret pair")
		}
	}
	return pk, nil
}

// UserSaver manages users.
type UserSaver struct {
	sync.Mutex
	config  *args.Config
	c       *cxo.Container
	store   map[cipher.PubKey]*UserConfig // All UserConfigs.
	masters map[cipher.PubKey]*UserConfig // UserConfigs of users we own.
	current cipher.PubKey                 // Currently active user.
}

// NewUserSaver creates a new UserSaver.
func NewUserSaver(config *args.Config, container *cxo.Container) (*UserSaver, error) {
	us := UserSaver{
		config:  config,
		c:       container,
		store:   make(map[cipher.PubKey]*UserConfig),
		masters: make(map[cipher.PubKey]*UserConfig),
	}
	us.checkUsers(us.load())
	if e := us.save(); e != nil {
		return nil, e
	}

	return &us, nil
}

func (us *UserSaver) absConfigDir() string {
	return filepath.Join(us.config.ConfigDir(), UserSaverFileName)
}

func (us *UserSaver) load() *UsersConfigFile {
	ucf := &UsersConfigFile{}
	// Don't load if specified not to.
	if !us.config.SaveConfig() {
		return ucf
	}
	log.Println("[USERSAVER] Loading configuration file...")
	// Load users from file.
	if e := util.LoadJSON(us.absConfigDir(), ucf); e != nil {
		log.Println("[USERSAVER]", e)
	}
	return ucf
}

// Checks if config is okay. Creates a random master user if not set.
func (us *UserSaver) checkUsers(ucf *UsersConfigFile) {
	// Check loaded users and intern in memory.
	for _, uc := range ucf.Users {
		log.Printf("\t- %v (master: %v)", uc.PubKey, uc.Master)
		upk, e := uc.Check()
		if e != nil {
			log.Println("\t\t config file check:", e)
			continue
		}
		us.store[upk] = uc
		if uc.Master == true {
			us.masters[upk] = uc
		}
		log.Println("\t\t loaded in memory")
	}
	// Check current user.
	upk, e := misc.GetPubKey(ucf.Current)
	if e != nil || upk == (cipher.PubKey{}) {
		log.Println("[USERSAVER] Current user invalid. Auto setting...")
		// Find one.
		if e := us.autoSetCurrent(); e != nil {
			// Make one.
			log.Println("[USERSAVER] Creating a random user.")
			upk, usk := cipher.GenerateKeyPair()
			uc := UserConfig{
				Alias:  misc.MakeRandomAlias(),
				Master: true,
				PubKey: upk.Hex(),
				SecKey: usk.Hex(),
			}
			us.store[upk] = &uc
			us.masters[upk] = &uc
			us.current = upk
			log.Println("[USERSAVER] Current user:", us.current.Hex())
			return
		}
		return
	}
	us.current = upk
	log.Println("[USERSAVER] Current user:", us.current.Hex())
	return
}

func (us *UserSaver) save() error {
	// Don't save if specified.
	if !us.config.SaveConfig() {
		return nil
	}
	ucf := UsersConfigFile{}
	ucf.Current = us.current.Hex()
	for _, uc := range us.store {
		ucf.Users = append(ucf.Users, uc)
	}
	return util.SaveJSON(us.absConfigDir(), ucf, os.FileMode(0700))
}

func (us *UserSaver) autoSetCurrent() error {
	for pk := range us.masters {
		us.current = pk
		log.Println("[USERSAVER] Current user:", us.current.Hex())
		return nil
	}
	return errors.New("no master users")
}

// List returns a list of all users that are in configuration.
func (us *UserSaver) List() []UserConfig {
	us.Lock()
	defer us.Unlock()
	list, i := make([]UserConfig, len(us.store)), 0
	for _, uc := range us.store {
		list[i] = *uc
		i += 1
	}
	return list
}

// ListMasters returns a list of all users that are master in configuration.
func (us *UserSaver) ListMasters() []UserConfig {
	us.Lock()
	defer us.Unlock()
	list, i := make([]UserConfig, len(us.store)), 0
	for _, uc := range us.masters {
		list[i] = *uc
		i += 1
	}
	return list
}

// Get gets a user configuration.
func (us *UserSaver) Get(upk cipher.PubKey) (UserConfig, bool) {
	us.Lock()
	defer us.Unlock()
	uc, has := us.store[upk]
	if has == false {
		return UserConfig{}, has
	}
	return *uc, has
}

// Get current returns the config of current user.
func (us *UserSaver) GetCurrent() UserConfig {
	us.Lock()
	defer us.Unlock()
	return *us.masters[us.current]
}

// Set Current sets the current user.
func (us *UserSaver) SetCurrent(upk cipher.PubKey) error {
	us.Lock()
	defer us.Unlock()
	if _, has := us.masters[upk]; has == false {
		return errors.New("not a master user")
	}
	us.current = upk
	return nil
}

// Add adds a user to configuration.
func (us *UserSaver) Add(alias string, upk cipher.PubKey) {
	us.Lock()
	defer us.Unlock()
	if uc, has := us.store[upk]; has {
		uc.Alias = alias
	} else {
		uc := UserConfig{Alias: alias, Master: false, PubKey: upk.Hex()}
		us.store[upk] = &uc
	}
	us.save()
}

// MasterAdd adds a board to configuration as master.
func (us *UserSaver) MasterAdd(alias string, upk cipher.PubKey, usk cipher.SecKey) {
	us.Lock()
	defer us.Unlock()
	uc := UserConfig{Alias: alias, Master: true, PubKey: upk.Hex(), SecKey: usk.Hex()}
	us.store[upk] = &uc
	us.masters[upk] = &uc
	us.save()
}

// Remove removes a user from configuration.
func (us *UserSaver) Remove(upk cipher.PubKey) error {
	us.Lock()
	defer us.Unlock()
	if _, has := us.masters[upk]; len(us.masters) == 1 && has == true {
		return errors.New("cannot remove only master user")
	}
	delete(us.store, upk)
	delete(us.masters, upk)
	if us.current == upk {
		us.autoSetCurrent()
	}
	us.save()
	return nil
}
