package session

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/store/object"
	"log"
	"os"
	"path"
	"sync"
	"github.com/skycoin/bbs/src/store/session/stores/memory_store"
	"github.com/skycoin/bbs/src/store/session/stores/drive_store"
)

const (
	usersLogPrefix = "USERS"
	usersSubDir    = "/users"
	usersFileExt   = ".ucf"
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
	current *string
	db Store
}

func NewManager(config *ManagerConfig) *Manager {
	m := &Manager{
		c: config,
		l: inform.NewLogger(true, os.Stdout, usersLogPrefix),
	}
	if *m.c.MemoryMode {
		m.db = memory_store.NewStore()
	} else {
		m.db = drive_store.NewStore(*m.c.ConfigDir, usersSubDir, usersFileExt)
	}
	return m
}

func (m *Manager) GetUsers() ([]string, error) {
	defer m.lock()()
	aliases, e := m.db.GetUsers()
	if e != nil {
		return nil, boo.WrapType(e, boo.Internal,
			"failed to get all users")
	}
	return aliases, nil
}

func (m *Manager) NewUser(in *object.NewUserIO) error {
	defer m.lock()()
	return m.db.NewUser(in.Alias, in.File)
}

func (m *Manager) DeleteUser(alias string) error {
	defer m.lock()()
	if m.current != nil && *m.current == alias {
		return ErrAlreadyLoggedIn
	}
	return m.db.DeleteUser(alias)
}

func (m *Manager) GetCurrentFile() (*object.UserFile, error) {
	defer m.lock()()
	if m.current == nil {
		return nil, ErrNotLoggedIn
	}
	f, ok := m.db.GetUser(*m.current)
	if !ok {
		return nil, boo.Newf(boo.Internal,
			"failed to get current user file")
	}
	return f, nil
}

func (m *Manager) Login(in *object.LoginIO) (*object.UserFile, error) {
	defer m.lock()()
	if m.current != nil {
		return nil, ErrAlreadyLoggedIn
	}
	if f, ok := m.db.GetUser(in.Alias); !ok {
		return nil, boo.Newf(boo.NotFound,
			"user of alias '%s' not found", in.Alias)
	} else {
		alias := in.Alias
		m.current = &alias
		return f, nil
	}
}

func (m *Manager) Logout() error {
	defer m.lock()()
	if m.current == nil {
		return ErrNotLoggedIn
	} else {
		m.current = nil
		return nil
	}
}

func (m *Manager) IsLoggedIn() bool {
	defer m.lock()()
	return m.current != nil
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
