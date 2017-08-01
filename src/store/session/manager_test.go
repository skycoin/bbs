package session

import (
	"context"
	"github.com/skycoin/bbs/src/misc/keys"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state"
	"github.com/skycoin/skycoin/src/cipher"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"testing"
)

const (
	removeTempDir = false
)

var (
	master          = false
	testMode        = false
	memoryMode      = false
	cxoPort         = 8998
	cxoRPCEnable    = false
	cxoRPCPort      = 8997
	compilerWorkers = 10
)

func createUserState() (*Manager, func()) {
	configDir, e := ioutil.TempDir("", "skybbs")
	if e != nil {
		log.Panic(e)
	}
	us, e := NewManager(
		&ManagerConfig{
			Master:       &master,
			TestMode:     &testMode,
			MemoryMode:   &memoryMode,
			ConfigDir:    &configDir,
			CXOPort:      &cxoPort,
			CXORPCEnable: &cxoRPCEnable,
			CXORPCPort:   &cxoRPCPort,
		},
		&state.CompilerConfig{
			Workers: &compilerWorkers,
		},
	)
	if e != nil {
		log.Panic(e)
	}
	return us, func() {
		us.Close()
		if removeTempDir {
			os.RemoveAll(configDir)
		}
	}
}

func createUsers(t *testing.T, us *Manager, n int) {
	for i := 0; i < n; i++ {
		iStr := strconv.Itoa(i)
		_, e := us.NewUser(context.Background(), &object.NewUserIO{
			Alias:    "person" + iStr,
			Seed:     "user" + iStr,
			Password: "password" + iStr,
		})
		if e != nil {
			t.Error(e)
		}
	}
}

func createPublicKeys(n int) []cipher.PubKey {
	keys := make([]cipher.PubKey, n)
	for i := 0; i < n; i++ {
		pk, _ := cipher.GenerateKeyPair()
		keys[i] = pk
	}
	return keys
}

func createSeeds(n int) []string {
	seeds := make([]string, n)
	for i := 0; i < n; i++ {
		var e error
		seeds[i], e = keys.GenerateSeed()
		if e != nil {
			panic(e)
		}
	}
	return seeds
}

func checkCount(t *testing.T, got, exp int) {
	if got != exp {
		t.Errorf("got %d, expected %d", got, exp)
	} else {
		t.Logf("got %d as expected", got)
	}
}

func login(ctx context.Context, t *testing.T, us *Manager, n int) {
	file, e := us.Login(ctx, &object.LoginIO{
		Alias:    "person" + strconv.Itoa(n),
		Password: "password" + strconv.Itoa(n),
	})
	if e != nil {
		t.Fatal("failed to login:", e)
	}
	t.Log("User File:", *file.GenerateView(us.cxo))
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
	createUsers(t, us, 10)

	users, e := us.GetUsers(ctx)
	if e != nil {
		t.Error(e)
	}

	{
		const expected = 10
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
		const expected = 9
		if count := len(users); count != expected {
			t.Errorf("expected %d users, got %d", expected, count)
		} else {
			t.Logf("got %d users as expected", expected)
		}
	}
}

func TestUserState_Login(t *testing.T) {
	const count = 4
	us, quit := createUserState()
	defer quit()

	ctx := context.Background()
	createUsers(t, us, count)

	for i := 0; i < count; i++ {
		login(ctx, t, us, i)

		file, e := us.GetInfo(ctx)
		if e != nil {
			t.Error(e)
		} else {
			t.Log("GetInfo:", *file.GenerateView(us.cxo))
		}

		if e := us.Logout(ctx); e != nil {
			t.Error(e)
		}

		file, e = us.GetInfo(ctx)
		if e != nil {
			t.Log("Got error as expected:", *file.GenerateView(us.cxo), e)
		} else {
			t.Error("Didn't get error:", *file.GenerateView(us.cxo))
		}
	}
}

func TestSessionManager_NewSubscription(t *testing.T) {
	const initCount = 10

	us, quit := createUserState()
	defer quit()

	ctx := context.Background()
	pks := createPublicKeys(initCount)

	createUsers(t, us, 1)
	login(ctx, t, us, 0)

	for i := 0; i < initCount; i++ {
		in := object.BoardIO{PubKeyStr: pks[i].Hex()}
		object.Process(in)
		file, e := us.NewSubscription(ctx, &in)
		if e != nil {
			t.Error("Failed to subscribe:", e)
		} else {
			checkCount(t, len(file.Subscriptions), i+1)
		}
	}

	for i := initCount - 1; i >= 0; i-- {
		in := object.BoardIO{PubKeyStr: pks[i].Hex()}
		object.Process(in)
		file, e := us.DeleteSubscription(ctx, &in)
		if e != nil {
			t.Error("Failed to unsubscribe:", e)
		} else {
			checkCount(t, len(file.Subscriptions), i)
		}
	}
}

func TestSessionManager_NewMaster(t *testing.T) {
	const initCount = 10

	us, quit := createUserState()
	defer quit()

	ctx := context.Background()
	pks := make([]cipher.PubKey, initCount)
	seeds := createSeeds(initCount)

	createUsers(t, us, 1)
	login(ctx, t, us, 0)

	for i := 0; i < initCount; i++ {
		in := object.NewBoardIO{
			Seed: seeds[i],
			Name: "Master Board " + strconv.Itoa(i),
			Desc: "A generated test board of index " + strconv.Itoa(i),
		}
		object.Process(in)
		file, e := us.NewMaster(ctx, &in)
		if e != nil {
			t.Error(e)
		} else {
			checkCount(t, len(file.Masters), i+1)
			pks[i] = file.Masters[i].PubKey
		}
	}

	for i := initCount - 1; i >= 0; i-- {
		in := object.BoardIO{
			PubKeyStr: pks[i].Hex(),
		}
		object.Process(in)
		file, e := us.DeleteMaster(ctx, &in)
		if e != nil {
			t.Error("Failed to delete master subscription:", e)
		} else {
			checkCount(t, len(file.Masters), i)
		}
	}
}
