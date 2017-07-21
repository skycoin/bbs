package boo

import (
	"strconv"
	"testing"
)

func TestElem_Error(t *testing.T) {
	msg := "something went wrong"
	e := New(Internal, msg)
	if e.Error() != msg {
		t.Errorf("got '%s', expected '%s'.", e.Error(), msg)
	}
	t.Log(e.Error())
}

func TestWrap(t *testing.T) {
	msgs := []string{"1", "2", "3", "4"}
	getExpected := func(f int) string {
		exp := "0"
		for i := 1; i <= f; i++ {
			exp = strconv.Itoa(i) + ": " + exp
		}
		return exp
	}
	e := New(Internal, "0")
	for i, m := range msgs {
		if i%2 == 0 {
			e = Wrap(e, m)
		} else {
			e = Wrapf(e, "%d", i+1)
		}
		got, exp := e.Error(), getExpected(i+1)
		if got != exp {
			t.Errorf("got '%s', expected '%s'", got, exp)
		} else {
			t.Logf("got '%s' as expected", got)
		}
	}
}

func TestType(t *testing.T) {
	t.Run("type at end", func(t *testing.T) {
		e := New(NotMaster, "this is a problem")
		e = Wrap(e, "something went wrong")
		e = Wrap(e, "woops")
		if Type(e) != NotMaster {
			t.Error("didn't get NotMaster")
		} else {
			t.Log("got NotMaster")
		}
	})
	t.Run("type change", func(t *testing.T) {
		e := New(NotMaster, "this is a problem")
		e = WrapType(e, InvalidRead, "how unfortunate")
		e = Wrap(e, "something went wrong")
		e = Wrap(e, "woops")
		if Type(e) != InvalidRead {
			t.Error("didn't get InvalidRead")
		} else {
			t.Log("got InvalidRead")
		}
	})
}
