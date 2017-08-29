package drive_store

import (
	"github.com/skycoin/bbs/src/store/object"
	"path"
	"sync"
	"io/ioutil"
	"strings"
	"github.com/skycoin/skycoin/src/util/file"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"os"
)

type Store struct {
	configDir, subDir, ext string

	sync.Mutex
}

func NewStore(configDir, subDir, ext string) *Store {
	store := &Store{
		configDir: configDir,
		subDir:    subDir,
		ext:       ext,
	}
	if e := os.MkdirAll(store.folderPath(), os.FileMode(0700)); e != nil {
		panic(e)
	}
	return store
}

func (s *Store) GetUsers() ([]string, error) {
	s.Lock()
	defer s.Unlock()

	files, e := ioutil.ReadDir(s.folderPath())
	if e != nil {
		return nil, e
	}

	var aliases []string
	for _, info := range files {
		if !info.IsDir() && strings.HasSuffix(info.Name(), s.ext) {
			name := strings.TrimSuffix(info.Name(), s.ext)
			aliases = append(aliases, name)
		}
	}

	return aliases, nil
}

func (s *Store) GetUser(alias string) (*object.UserFile, bool) {
	s.Lock()
	defer s.Unlock()

	data, e := ioutil.ReadFile(s.filePath(alias))
	if e != nil {
		return nil, false
	}

	out := new(object.UserFile)
	if e := encoder.DeserializeRaw(data, out); e != nil {
		return nil, false
	}

	return out, true
}

func (s *Store) NewUser(alias string, f *object.UserFile) error {
	s.Lock()
	defer s.Unlock()

	return file.SaveBinary(
		s.filePath(alias),
		encoder.Serialize(f),
		os.FileMode(0600),
	)
}

func (s *Store) DeleteUser(alias string) error {
	s.Lock()
	defer s.Unlock()

	return os.Remove(s.filePath(alias))
}

func (s *Store) folderPath() string {
	return path.Join(s.configDir, s.subDir)
}

func (s *Store) filePath(alias string) string {
	return path.Join(s.folderPath(), alias+s.ext)
}
