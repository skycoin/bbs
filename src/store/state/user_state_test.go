package state

import (
	"testing"
	"log"
	"context"
)

func createUserState() *UserState {
	var (
		configDir = "/home/evan/.skybbs"
		memory = false
	)
	us, e := NewUserState(&UserStateConfig{
		ConfigDir: &configDir,
		Memory: &memory,
	})
	if e != nil {
		log.Panic(e)
	}
	return us
}

func TestUserState_GetUsers(t *testing.T) {
	us := createUserState()
	defer us.Close()
	users, e := us.GetUsers(context.Background())
	if e != nil {
		t.Error(e)
	}
	t.Log(users)
}