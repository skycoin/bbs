package state

import (
	"context"
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/src/boo"
	"github.com/skycoin/bbs/src/misc/hide"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/store/obj"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/util/file"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

const (
	logPrefix    = "USERSTATE"
	subDir       = "/users"
	extension    = ".dat"
	timeout      = time.Second * 5
	saveDuration = time.Second * 5
)

var (
	ErrNoUserLoaded    = errors.New("no user has been loaded")
	ErrAlreadyLoggedIn = errors.New("already logged in")
)

// UserStateConfig configures a UserState.
type UserStateConfig struct {
	ConfigDir *string `json:"config_dir"`
	Memory    *bool   `json:"memory"`
}

// UserState represents a user state.
type UserState struct {
	c        *UserStateConfig
	l        *log.Logger
	file     *UserFile
	key      []byte
	changes  bool
	requests chan interface{}
	quit     chan struct{}
	wg       sync.WaitGroup
}

// NewUserState creates a new UserState with configuration.
func NewUserState(config *UserStateConfig) (*UserState, error) {
	us := &UserState{
		c:        config,
		l:        inform.NewLogger(true, os.Stdout, logPrefix),
		changes: false,
		requests: make(chan interface{}),
		quit:     make(chan struct{}),
	}
	e := os.MkdirAll(us.folderPath(), os.FileMode(0700))
	if e != nil {
		return nil, e
	}
	go us.service()
	return us, nil
}

// Close closes the UserState service.
func (s *UserState) Close() {
	select {
	case s.quit <- struct{}{}:
	default:
	}
	s.wg.Wait()
}

func (s *UserState) folderPath() string {
	return path.Join(*s.c.ConfigDir, subDir)
}

func (s *UserState) filePath(user string) string {
	return path.Join(s.folderPath(), user+extension)
}

func (s *UserState) timedContext(ctx context.Context) context.Context {
	ctx, _ = context.WithTimeout(ctx, timeout)
	return ctx
}

func (s *UserState) service() {
	s.wg.Add(1)
	defer s.wg.Done()

	ticker := time.NewTicker(saveDuration)
	defer ticker.Stop()

	for {
		select {
		case req := <-s.requests:
			s.processRequest(req)
		case <-ticker.C:
			s.processSave()
		case <-s.quit:
			s.processSave()
			return
		}
	}
}

/*
	<<< TRIGGERS : SAVE >>>
*/

func (s *UserState) processSave() {
	s.changes = false
	if *s.c.Memory || s.file == nil {
		return
	}
	data, e := hide.Encrypt(s.key, encoder.Serialize(*s.file))
	if e != nil {
		s.l.Printf("Error: %v", e)
		return
	}
	e = file.SaveBinary(
		s.filePath(s.file.User.Alias),
		data, os.FileMode(0600),
	)
	if e != nil {
		s.l.Printf("Error: %v", e)
		return
	}
}

/*
	<<< REQUESTS : PROCESS >>>
*/

func (s *UserState) processRequest(req interface{}) {
	switch req.(type) {
	case *reqGetUsers:
		s.processGetUsers(req.(*reqGetUsers))

	case *reqNewUser:
		s.processNewUser(req.(*reqNewUser))

	case *reqDeleteUser:
		s.processDeleteUser(req.(*reqDeleteUser))

	case *reqLogin:
		s.processLogin(req.(*reqLogin))

	case *reqLogout:
		s.processLogout(req.(*reqLogout))

	default:
		s.l.Printf("Unprocessed request of type '%T'", req)
	}
}

/*
	<<< REQUESTS : GET USERS >>>
*/

type reqGetUsers struct {
	out chan []string
	e   chan error
}

func (s *UserState) processGetUsers(req *reqGetUsers) {
	files, e := ioutil.ReadDir(s.folderPath())
	if e != nil {
		req.e <- e
		return
	}
	var users []string
	for _, info := range files {
		if info.IsDir() || !strings.HasSuffix(info.Name(), extension) {
			continue
		}
		name := strings.TrimSuffix(info.Name(), extension)
		s.l.Printf("Found User: '%s'.", name)
		users = append(users, name)
	}
	req.out <- users
}

// GetUsers obtains list of available out.
func (s *UserState) GetUsers(ctx context.Context) ([]string, error) {
	ctx = s.timedContext(ctx)
	req := reqGetUsers{
		out: make(chan []string),
		e:   make(chan error),
	}
	go func() { s.requests <- &req }()
	select {
	case <-ctx.Done():
		return nil, boo.WrapType(ctx.Err(), boo.Internal)
	case users := <-req.out:
		return users, nil
	case e := <-req.e:
		return nil, e
	}
}

/*
	<<< REQUESTS : NEW USER >>>
*/

type NewUserInput struct {
	Alias    string `json:"alias"`
	Seed     string `json:"seed"`
	Password string `json:"password"`
}

type reqNewUser struct {
	in  *NewUserInput
	out chan UserFile
	e   chan error
}

func (s *UserState) processNewUser(req *reqNewUser) {
	pk, sk := cipher.GenerateDeterministicKeyPair([]byte(req.in.Seed))
	key := []byte(req.in.Password)
	uFile := UserFile{
		User: obj.User{
			Alias:     req.in.Alias,
			PublicKey: pk,
			SecretKey: sk,
		},
	}
	data, e := hide.Encrypt(key, encoder.Serialize(uFile))
	if e != nil {
		req.e <- e
		return
	}
	e = file.SaveBinary(s.filePath(req.in.Alias), data, os.FileMode(0600))
	if e != nil {
		req.e <- e
		return
	}
	req.out <- uFile
}

// NewUser creates a new user.
func (s *UserState) NewUser(ctx context.Context, in *NewUserInput) (*UserFileView, error) {
	ctx = s.timedContext(ctx)
	// TODO: Check inputs.

	req := reqNewUser{
		in:  in,
		out: make(chan UserFile),
		e:   make(chan error),
	}
	go func() { s.requests <- &req }()
	select {
	case <-ctx.Done():
		return nil, boo.WrapType(ctx.Err(), boo.Internal)
	case out := <-req.out:
		return out.GenerateView(), nil
	case e := <-req.e:
		return nil, e
	}
}

/*
	<<< REQUESTS : DELETE USER >>>
*/

type reqDeleteUser struct {
	alias string
	e     chan error
}

func (s *UserState) processDeleteUser(req *reqDeleteUser) {
	if s.file != nil {
		req.e <- ErrAlreadyLoggedIn
		return
	}
	req.e <- os.Remove(s.filePath(req.alias))
}

// DeleteUser deletes a user configuration file.
func (s *UserState) DeleteUser(ctx context.Context, alias string) error {
	ctx = s.timedContext(ctx)
	req := reqDeleteUser{
		alias: alias,
		e:     make(chan error),
	}
	go func() { s.requests <- &req }()
	select {
	case <-ctx.Done():
		return boo.WrapType(ctx.Err(), boo.Internal)
	case e := <-req.e:
		return e
	}
}

/*
	<<< REQUESTS : LOGIN >>>
*/

type LoginInput struct {
	Alias    string `json:"alias"`
	Password string `json:"password"`
}

type reqLogin struct {
	in  *LoginInput
	out chan *UserFileView
	e   chan error
}

func (s *UserState) processLogin(req *reqLogin) {
	if s.file != nil {
		req.e <- ErrAlreadyLoggedIn
		return
	}

	fName := s.filePath(req.in.Alias)
	key := []byte(req.in.Password)

	data, e := ioutil.ReadFile(fName)
	if e != nil {
		req.e <- boo.WrapType(e, boo.ObjectNotFound, "user file not found")
		return
	}

	data, e = hide.Decrypt(key, data)
	if e != nil {
		req.e <- boo.WrapType(e, boo.NotAuthorised, "decryption failed")
		return
	}

	uFile := &UserFile{}
	if e := encoder.DeserializeRaw(data, uFile); e != nil {
		req.e <- boo.WrapType(e, boo.InvalidRead, "corrupt user file")
		return
	}

	s.file, s.key, s.changes = uFile, key, true
	req.out <- uFile.GenerateView()
}

// Login logs in a user.
func (s *UserState) Login(ctx context.Context, in *LoginInput) (*UserFileView, error) {
	ctx = s.timedContext(ctx)
	req := reqLogin{
		in:  in,
		out: make(chan *UserFileView),
		e:   make(chan error),
	}
	go func() { s.requests <- &req }()
	select {
	case <-ctx.Done():
		return nil, boo.WrapType(ctx.Err(), boo.Internal)
	case out := <-req.out:
		return out, nil
	case e := <-req.e:
		return nil, e
	}
}

/*
	<<< REQUESTS : LOGOUT >>>
*/

type reqLogout struct {
	e chan error
}

func (s *UserState) processLogout(req *reqLogout) {
	s.processSave()
	s.file, s.key = nil, nil
	req.e <- nil
}

func (s *UserState) Logout(ctx context.Context) error {
	ctx = s.timedContext(ctx)
	req := reqLogout{
		e: make(chan error),
	}
	go func() {s.requests <- &req}()
	select {
	case <-ctx.Done():
		return boo.WrapType(ctx.Err(), boo.Internal)
	case e := <- req.e:
		return e
	}
	return nil
}

/*
	<<< REQUESTS : GET SESSION >>>
*/

type reqGetSession struct {

}

// TODO: add master, remove master, add subscription, remove subscription.
