package users

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/misc/tag"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/util/file"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
)

const (
	usersLogPrefix = "USERS"
	usersSubDir    = "/users"
	usersFileExt   = ".user.v2"
)

var (
	ErrAlreadyLoggedIn = boo.New(boo.NotAllowed, "already logged in")
	ErrNotLoggedIn     = boo.New(boo.NotAllowed, "not logged in")
)

type ManagerConfig struct {
	MemoryMode *bool
	ConfigDir  *string
}

type Manager struct {
	c *ManagerConfig
	l *log.Logger

	mux  sync.Mutex
	file *File
}

func NewManager(config *ManagerConfig) (*Manager, error) {
	m := &Manager{
		c: config,
		l: inform.NewLogger(true, os.Stdout, usersLogPrefix),
	}
	if e := os.MkdirAll(m.folderPath(), os.FileMode(0700)); e != nil {
		return nil, e
	}
	return m, nil
}

func (m *Manager) GetUsers() ([]string, error) {
	defer m.lock()()
	files, e := ioutil.ReadDir(m.folderPath())
	if e != nil {
		return nil, boo.WrapType(e, boo.Internal,
			"failed to read directory")
	}
	var aliases []string
	for _, info := range files {
		if !info.IsDir() && strings.HasSuffix(info.Name(), usersFileExt) {
			name := strings.TrimSuffix(info.Name(), usersFileExt)
			m.l.Printf("Found user: '%s'.", name)
			aliases = append(aliases, name)
		}
	}
	return aliases, nil
}

func (m *Manager) NewUser(in *object.NewUserIO) error {
	defer m.lock()()
	if *m.c.MemoryMode {
		return nil
	}
	f := &File{
		User: object.User{
			Alias:     in.Alias,
			PublicKey: in.UserPubKey,
			SecretKey: in.UserSecKey,
		},
		Seed: in.Seed,
	}
	return file.SaveBinary(
		m.filePath(f.User.Alias), encoder.Serialize(f), os.FileMode(0600))
}

func (m *Manager) DeleteUser(alias string) error {
	defer m.lock()()
	if e := os.Remove(m.filePath(alias)); e != nil {
		return e
	}
	if m.file != nil && m.file.User.Alias == alias {
		return ErrAlreadyLoggedIn
	}
	return nil
}

func (m *Manager) GetCurrentFile() (*File, error) {
	defer m.lock()()
	if m.file == nil {
		return nil, ErrNotLoggedIn
	}
	return &(*m.file), nil
}

func (m *Manager) GetUPK() cipher.PubKey {
	defer m.lock()()
	if m.file == nil {
		return cipher.PubKey{}
	}
	return m.file.User.PublicKey
}

func (m *Manager) Sign(obj interface{}) error {
	defer m.lock()()
	if m.file == nil {
		return ErrNotLoggedIn
	}
	tag.Sign(
		obj,
		m.file.User.PublicKey,
		m.file.User.SecretKey,
	)
	return nil
}

func (m *Manager) Login(in *object.LoginIO) (*File, error) {
	defer m.lock()()
	if m.file != nil {
		return nil, ErrAlreadyLoggedIn
	}
	m.file = new(File)

	data, e := ioutil.ReadFile(m.filePath(in.Alias))
	if e != nil {
		if os.IsNotExist(e) {
			return nil, boo.WrapType(e, boo.NotFound,
				"user not found")
		}
		return nil, boo.WrapType(e, boo.InvalidRead,
			"failed to read user file")
	}
	if e := encoder.DeserializeRaw(data, m.file); e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead,
			"corrupt user file")
	}
	return &(*m.file), nil
}

func (m *Manager) Logout() error {
	defer m.lock()()
	if m.file == nil {
		return ErrNotLoggedIn
	}
	m.file = nil
	return nil
}

func (m *Manager) IsLoggedIn() bool {
	defer m.lock()()
	return m.file != nil
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func (m *Manager) lock() func() {
	m.mux.Lock()
	return m.mux.Unlock
}

func (m *Manager) folderPath() string {
	return path.Join(*m.c.ConfigDir, usersSubDir)
}

func (m *Manager) filePath(name string) string {
	return path.Join(m.folderPath(), name+usersFileExt)
}
