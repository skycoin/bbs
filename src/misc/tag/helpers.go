package tag

import (
	"fmt"
	"github.com/skycoin/bbs/src/misc/boo"
	"reflect"
	"strings"
)

type tMap map[string]reflect.Value

func (tm tMap) set(key string, v interface{}) {
	tm[key].Set(reflect.ValueOf(v))
}

func makeTagMap(tagKey string, rVal reflect.Value, rTyp reflect.Type) tMap {
	out := make(tMap)
	for i := 0; i < rTyp.NumField(); i++ {
		if vTag, has := getTagKey(tagKey, rTyp, i); has {
			field := rVal.Field(i)
			// TODO: LATER: Have nested fields additions.
			out[vTag] = field
		}
	}
	return out
}

func getReflectPair(v interface{}) (reflect.Value, reflect.Type) {
	rVal := reflect.ValueOf(v).Elem()
	return rVal, rVal.Type()
}

func getTagKey(tagKey string, rTyp reflect.Type, i int) (string, bool) {
	return rTyp.Field(i).Tag.Lookup(tagKey)
}

func clearField(fv reflect.Value) {
	fv.Set(reflect.Zero(fv.Type()))
}

func splitStr(str string, check func(v string) bool) ([]string, error) {
	out := strings.Split(str, ",")
	for i := len(out) - 1; i >= 0; i-- {
		out[i] = strings.TrimSpace(out[i])
		if !check(out[i]) || out[0] == "" {
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
