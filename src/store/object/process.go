package object

import (
	"fmt"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/keys"
	"github.com/skycoin/skycoin/src/cipher"
	"reflect"
	"strconv"
	"strings"
)

const (
	tagKey          = "bbs"
	valUserPKStr    = "upkStr"
	valUserPK       = "upk"
	valUserSK       = "usk"
	valBoardPKStr   = "bpkStr"
	valBoardPK      = "bpk"
	valBoardSK      = "bsk"
	valThreadRefStr = "tRefStr"
	valThreadRef    = "tRef"
	valPostRefStr   = "pRefStr"
	valPostRef      = "pRef"
	valBoardSeed    = "bSeed"
	valUserSeed     = "uSeed"
	valHeading      = "heading"
	valBody         = "body"
	valAlias        = "alias"
	valPassword     = "password"
	valAddress      = "address"
	valSubAddrsStr  = "subAddrsStr"
	valSubAddrs     = "subAddrs"
	valConsStr      = "consStr"
	valCons         = "cons"
	valModeStr      = "modeStr"
	valMode         = "mode"
	valTagStr       = "tagStr"
	valTag          = "tag"
)

type TagMap map[string]reflect.Value

func (tm TagMap) Set(key string, v interface{}) {
	tm[key].Set(reflect.ValueOf(v))
}

func Process(obj interface{}) error {
	if e := process(makeTagMap(getReflectPair(obj))); e != nil {
		return e
	}
	return nil
}

func process(tm TagMap) error {
	// User public key.
	if upkStr, has := tm[valUserPKStr]; has {
		upk, e := keys.GetPubKey(upkStr.String())
		if e != nil {
			return wrapErr(e, "user public key")
		}
		tm.Set(valUserPK, upk)
	}
	// Board public key.
	if bpkStr, has := tm[valBoardPKStr]; has {
		bpk, e := keys.GetPubKey(bpkStr.String())
		if e != nil {
			return wrapErr(e, "board public key")
		}
		tm.Set(valBoardPK, bpk)
	}
	// Thread reference.
	if tRefStr, has := tm[valThreadRefStr]; has {
		tRef, e := keys.GetReference(tRefStr.String())
		if e != nil {
			return wrapErr(e, "thread reference")
		}
		tm.Set(valThreadRef, tRef)
	}
	// Post reference.
	if pRefStr, has := tm[valPostRefStr]; has {
		pRef, e := keys.GetReference(pRefStr.String())
		if e != nil {
			return wrapErr(e, "post reference")
		}
		tm.Set(valPostRef, pRef)
	}
	// Submission addresses.
	if subAddrsStr, has := tm[valSubAddrsStr]; has {
		subAddrs, e := splitStr(subAddrsStr.String(), func(v string) bool {
			return true
		})
		if e != nil {
			return wrapErr(e, "submission addresses")
		}
		tm.Set(valSubAddrs, subAddrs)
	}
	// Connections.
	if consStr, has := tm[valConsStr]; has {
		cons, e := splitStr(consStr.String(), func(v string) bool {
			return true
		})
		if e != nil {
			return wrapErr(e, "connections")
		}
		tm.Set(valCons, cons)
	}
	// Vote mode.
	if modeStr, has := tm[valModeStr]; has {
		mode, e := strconv.Atoi(modeStr.String())
		if e != nil {
			return boo.WrapType(e, boo.InvalidInput,
				"invalid vote mode input")
		}
		switch int8(mode) {
		case -1, 0, +1:
		default:
			return boo.New(boo.InvalidInput,
				"invalid vote mode input")
		}
		tm.Set(valMode, int8(mode))
	}
	// Vote tag.
	if tagStr, has := tm[valTagStr]; has {
		tag := tagStr.String()
		switch tag {
		case "", "spam":
		default:
			return boo.New(boo.InvalidInput,
				"invalid vote tag input")
		}
		tm.Set(valTag, []byte(tag))
	}
	// User seed.
	if uSeed, has := tm[valUserSeed]; has {
		if e := CheckSeed(uSeed.String()); e != nil {
			return e
		}
		pk, sk := cipher.GenerateDeterministicKeyPair(
			[]byte(uSeed.String()))
		tm.Set(valUserPK, pk)
		tm.Set(valUserSK, sk)
	}
	// Board seed.
	if bSeed, has := tm[valBoardSeed]; has {
		if e := CheckSeed(bSeed.String()); e != nil {
			return e
		}
		pk, sk := cipher.GenerateDeterministicKeyPair(
			[]byte(bSeed.String()))
		tm.Set(valBoardPK, pk)
		tm.Set(valBoardSK, sk)
	}
	// Heading text.
	if heading, has := tm[valHeading]; has {
		if e := CheckHeading(heading.String()); e != nil {
			return e
		}
	}
	// Body text.
	if body, has := tm[valBody]; has {
		if e := CheckBody(body.String()); e != nil {
			return e
		}
	}
	// Alias text.
	if alias, has := tm[valAlias]; has {
		if e := CheckAlias(alias.String()); e != nil {
			return e
		}
	}
	// Password text.
	if password, has := tm[valPassword]; has {
		if e := CheckPassword(password.String()); e != nil {
			return e
		}
	}
	// Address text.
	if address, has := tm[valAddress]; has {
		if e := CheckAddress(address.String()); e != nil {
			return e
		}
	}
	return nil
}

/*
	<<< CHECKING FUNCTIONS >>>
*/

// CheckSeed ensures validity of seed. TODO
func CheckSeed(seed string) error {
	return nil
}

// CheckHeading ensures validity of board/thread/post name. TODO
func CheckHeading(heading string) error {
	return nil
}

// CheckBody ensures validity of board/thread/post description. TODO
func CheckBody(body string) error {
	return nil
}

// CheckAlias ensures validity of user alias. TODO
func CheckAlias(alias string) error {
	return nil
}

// CheckPassword ensures validity of password. TODO
func CheckPassword(password string) error {
	return nil
}

// CheckAddress ensures validity of address. TODO
func CheckAddress(address string) error {
	return nil
}

// CheckMode check's the vote's mode.
func CheckMode(mode int8) error {
	switch mode {
	case -1, 0, +1:
		return nil
	default:
		return boo.Newf(boo.InvalidInput,
			"invalid vote mode of %d provided", mode)
	}
}

// CheckTag check's the vote's tag.
func CheckTag(tag []byte) error {
	switch string(tag) {
	case "", "spam":
		return nil
	default:
		return boo.Newf(boo.InvalidInput,
			"invalid vote tag of %s provided", string(tag))
	}
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func splitStr(str string, check func(v string) bool) ([]string, error) {
	out := strings.Split(str, ",")
	for i := len(out) - 1; i >= 0; i-- {
		out[i] = strings.TrimSpace(out[i])
		if !check(out[i]) {
			out[i], out[0] = out[0], out[i]
			out = out[1:]
		}
	}
	return out, nil
}

func wrapErr(e error, what string) error {
	return boo.WrapType(e, boo.Internal,
		fmt.Sprintln("failed to process", what))
}

func makeTagMap(rVal reflect.Value, rTyp reflect.Type) TagMap {
	out := make(TagMap)
	for i := 0; i < rTyp.NumField(); i++ {
		if tagVal, has := getTagKey(rTyp, i); has {
			out[tagVal] = rVal.Field(i)
		}
	}
	return out
}

func getReflectPair(v interface{}) (reflect.Value, reflect.Type) {
	rVal := reflect.ValueOf(v).Elem()
	return rVal, rVal.Type()
}

func getTagKey(rTyp reflect.Type, i int) (string, bool) {
	return rTyp.Field(i).Tag.Lookup(tagKey)
}
