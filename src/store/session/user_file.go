package session

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/skycoin/src/cipher"
	"reflect"
)

const (
	tagKey       = "bbs"
	boardPKValue = "bpk"
	boardSKValue = "bsk"
	userPKValue  = "upk"
	userSKValue  = "usk"
)

var (
	// ErrEmpty occurs when object is nil.
	ErrEmpty = boo.New(boo.NotFound, "nil error")
)

func corruptWrap(e error) error {
	return boo.WrapType(e, boo.InvalidRead, "corrupt user file")
}

// UserFileView represents a user user as displayed to end user.
type UserFileView struct {
	User          object.UserView           `json:"user"`
	Subscriptions []object.SubscriptionView `json:"subscriptions"`
	Masters       []object.SubscriptionView `json:"master_subscriptions"`
	Connections   []object.ConnectionView   `json:"connections"`
}

// UserFile represents a user of user configuration.
type UserFile struct {
	User          object.User           `json:"user"`
	Subscriptions []object.Subscription `json:"subscriptions"`
	Masters       []object.Subscription `json:"master_subscriptions"`
	Connections   []string              `json:"connections"`
}

// Check ensures the validity of the UserFile.
func (f *UserFile) Check() error {
	if f == nil {
		return ErrEmpty
	}
	if e := f.User.PublicKey.Verify(); e != nil {
		return corruptWrap(e)
	}
	if e := f.User.SecretKey.Verify(); e != nil {
		return corruptWrap(e)
	}
	for i, sub := range f.Subscriptions {
		if e := sub.PubKey.Verify(); e != nil {
			return corruptWrap(e)
		}
		f.Subscriptions[i].SecKey = cipher.SecKey{}
	}
	for _, sub := range f.Masters {
		if e := sub.PubKey.Verify(); e != nil {
			return corruptWrap(e)
		}
		if e := sub.SecKey.Verify(); e != nil {
			return corruptWrap(e)
		}
	}
	return nil
}

// GenerateView generates something readable for front end.
func (f *UserFile) GenerateView(cxo *CXO) *UserFileView {
	view := new(UserFileView)

	if e := f.Check(); e != nil {
		return view
	}

	// Fill "User".
	view.User = object.UserView{
		User:      object.User{Alias: f.User.Alias},
		PublicKey: f.User.PublicKey.Hex(),
		SecretKey: f.User.SecretKey.Hex(),
	}

	// Fill "Subscriptions".
	subscriptions := make([]object.SubscriptionView, len(f.Subscriptions))
	for i, s := range f.Subscriptions {
		subscriptions[i] = object.SubscriptionView{
			PubKey: s.PubKey.Hex(),
			SecKey: s.SecKey.Hex(),
		}
	}
	view.Subscriptions = subscriptions

	// Fill "Masters".
	masters := make([]object.SubscriptionView, len(f.Masters))
	for i, m := range f.Masters {
		masters[i] = object.SubscriptionView{
			PubKey: m.PubKey.Hex(),
			SecKey: m.SecKey.Hex(),
		}
	}
	view.Masters = masters

	// Fill "Connections".
	connections := make([]object.ConnectionView, len(f.Connections))
	activeConnectionsMap := make(map[string]bool)
	activeConnections, _ := cxo.GetConnections()
	for _, address := range activeConnections {
		activeConnectionsMap[address] = true
	}
	for i, address := range f.Connections {
		connections[i] = object.ConnectionView{
			Address: address,
			Active:  activeConnectionsMap[address],
		}
	}
	view.Connections = connections

	return view
}

// FindMaster finds the index of a master subscription.
// If not found, returns an error.
func (f *UserFile) FindMaster(pk cipher.PubKey) (int, error) {
	for i, sub := range f.Masters {
		if sub.PubKey == pk {
			return i, nil
		}
	}
	return -1, boo.Newf(boo.NotFound,
		"board %s not found as master", pk.Hex())
}

func (f *UserFile) FillMaster(v interface{}) error {
	rVal, rTyp := getReflectPair(v)

	var e error
	var mIndex = -1

	for i := 0; i < rTyp.NumField(); i++ {
		if tagVal, has := getTagKey(rTyp, i); has {
			field := rVal.Field(i)
			switch tagVal {
			case boardPKValue:
				mIndex, e = f.FindMaster(field.Interface().(cipher.PubKey))
				if e != nil {
					return e
				}
			case boardSKValue:
				if mIndex == -1 {
					panic(boo.New(boo.Internal,
						"struct has no field with '%s' tag", boardPKValue))
				}
				field.Set(reflect.ValueOf(
					f.Masters[mIndex].SecKey))
				return nil
			}
		}
	}
	return nil
}

func (f *UserFile) FillUser(v interface{}) {
	rVal, rTyp := getReflectPair(v)
	for i := 0; i < rTyp.NumField(); i++ {
		if tagVal, has := getTagKey(rTyp, i); has {
			field := rVal.Field(i)
			switch tagVal {
			case userPKValue:
				field.Set(reflect.ValueOf(f.User.PublicKey))
			case userSKValue:
				field.Set(reflect.ValueOf(f.User.SecretKey))
				return
			}
		}
	}
}

func getReflectPair(v interface{}) (reflect.Value, reflect.Type) {
	rVal := reflect.ValueOf(v).Elem()
	return rVal, rVal.Type()
}

func getTagKey(rTyp reflect.Type, i int) (string, bool) {
	return rTyp.Field(i).Tag.Lookup(tagKey)
}
