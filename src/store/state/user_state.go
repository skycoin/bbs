package state

import (
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/src/store/obj"
	"time"
	"path/filepath"
	"os"
	"strings"
	"log"
	"context"
)

const (
	extension = ".json"
	timeout = time.Second * 5
	saveDuration = time.Second * 5
)

var (
	ErrNoUserLoaded = errors.New("no user has been loaded")
)

// UserFile represents a file of user configuration.
type UserFile struct {
	User          obj.User           `json:"user"`
	Subscriptions []obj.Subscription `json:"subscriptions"`
}

// UserStateConfig configures a UserState.
type UserStateConfig struct {
	ConfigDir *string `json:"config_dir"`
	Memory    *bool   `json:"memory"`
}

// UserState represents a user state.
type UserState struct {
	cf       *UserStateConfig
	file     *UserFile
	key      string
	requests chan interface{}
	quit     chan struct{}
}

// NewUserState creates a new UserState with configuration.
func NewUserState(config *UserStateConfig) (*UserState, error) {
	us := &UserState{
		cf:       config,
		requests: make(chan interface{}),
		quit:     make(chan struct{}),
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
}

func (s *UserState) service() {
	ticker := time.NewTicker(saveDuration)
	defer ticker.Stop()

	for {
		select {
		case req := <-s.requests:
			s.processRequest(req)
		case <-ticker.C:
		case <-s.quit:
			return
		}
	}
}

func (s *UserState) processRequest(req interface{}) {
	switch req.(type) {
	case *reqGetUsers:
		s.processGetUsers(req.(*reqGetUsers))
	}
}

/*
	<<< GET USERS >>>
*/

type reqGetUsers struct {
	users chan []string
	e chan error
}

func (s *UserState) processGetUsers(req *reqGetUsers) {
	var users []string
	e := filepath.Walk(*s.cf.ConfigDir, func(_ string, info os.FileInfo, e error) error {
		if e != nil || info.IsDir() || !strings.HasSuffix(info.Name(), extension) {
			return nil
		}
		name := strings.TrimSuffix(info.Name(), extension)
		log.Printf("[USERSTATE] Found User: '%s'.", name)
		users = append(users, name)
		return nil
	})
	if e != nil {
		req.e <- e
	} else {
		req.users <- users
	}
}

// GetUsers obtains list of available users.
func (s *UserState) GetUsers(ctx context.Context) ([]string, error) {
	ctx, _ = context.WithTimeout(ctx, timeout)
	req := &reqGetUsers{
		users: make(chan []string),
		e: make(chan error),
	}
	s.requests <- req
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case users := <-req.users:
		return users, nil
	case e := <- req.e:
		return nil, e
	}
}