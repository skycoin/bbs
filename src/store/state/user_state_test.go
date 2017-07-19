package state

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"
)

func createUserState() (*UserState, func()) {
	memory := false
	configDir, e := ioutil.TempDir("", "skybbs")
	if e != nil {
		log.Panic(e)
	}
	us, e := NewUserState(&UserStateConfig{
		ConfigDir: &configDir,
		Memory:    &memory,
	})
	if e != nil {
		log.Panic(e)
	}
	return us, func() {
		us.Close()
		os.RemoveAll(configDir)
	}
}

func createUsers(t *testing.T, us *UserState, n int) {
	for i := 0; i < n; i++ {
		iStr := strconv.Itoa(i)
		_, e := us.NewUser(context.Background(), &NewUserInput{
			Alias:    "user" + iStr,
			Seed:     "user" + iStr,
			Password: "password" + iStr,
		})
		if e != nil {
			t.Error(e)
		}
	}
}

func TestUserState_GetUsers(t *testing.T) {
	us, quit := createUserState()
	defer quit()

	createUsers(t, us, 10)

	users, e := us.GetUsers(context.Background())
	if e != nil {
		t.Error(e)
	}
	t.Log(users)
}

func TestUserState_DeleteUser(t *testing.T) {
	us, quit := createUserState()
	defer quit()

	ctx := context.Background()
	createUsers(t, us, 1)

	users, e := us.GetUsers(ctx)
	if e != nil {
		t.Error(e)
	}

	{
		const expected = 1
		if count := len(users); count != expected {
			t.Errorf("expected %d users, got %d", expected, count)
		} else {
			t.Logf("got %d users as expected", expected)
		}
	}

	e = us.DeleteUser(ctx, users[0])
	if e != nil {
		t.Error(e)
	}

	users, e = us.GetUsers(ctx)
	if e != nil {
		t.Error(e)
	}

	{
		const expected = 0
		if count := len(users); count != expected {
			t.Errorf("expected %d users, got %d", expected, count)
		} else {
			t.Logf("got %d users as expected", expected)
		}
	}
}

func TestUserState_Login(t *testing.T) {
	us, quit := createUserState()
	defer quit()

	ctx := context.Background()
	createUsers(t, us, 3)

	view, e := us.Login(ctx, &LoginInput{
		Alias:    "user0",
		Password: "password0",
	})
	if e != nil {
		t.Error(e)
	}
	t.Log(*view)

	if e := us.Logout(ctx); e != nil {
		t.Error(e)
	}
}