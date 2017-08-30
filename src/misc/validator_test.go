package misc

import (
	"testing"
	"fmt"
)

func TestCheckAlias(t *testing.T) {

	_, e := CheckAlias("aa")
	if  e == nil {
		t.Error("min length check failed")
		fmt.Print(e.Error())

	}

	_, e = CheckAlias("aaddddddddddddddddddddddddddddddddddddddddddddddd")
	if  e == nil {
		t.Error("max length check failed")
	}

	_, e = CheckAlias("_____")
	if  e == nil {
		t.Error("username failed")
	}

}

