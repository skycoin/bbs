package tag

import (
	"reflect"
	"time"
)

const transferKey = "transfer"

// Transfer performs transfer.
func Transfer(src, dst interface{}) {
	dstVal, dstTyp := getReflectPair(dst)
	dstMap := makeTagMap(transferKey, dstVal, dstTyp)

	if v, has := dstMap["time"]; has {
		v.Set(reflect.ValueOf(time.Now().UnixNano()))
	}

	srcVal, srcTyp := getReflectPair(src)

	for i := 0; i < srcTyp.NumField(); i++ {
		if vTag, has := getTagKey(transferKey, srcTyp, i); has {
			if dstMap[vTag].CanSet() {
				dstMap[vTag].Set(srcVal.Field(i))
			}
		}
	}
}
