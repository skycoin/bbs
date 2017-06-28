package misc

import "testing"

func TestMakeIntBetween(t *testing.T) {
	for i := 0; i < 100; i++ {
		n, e := MakeIntBetween(i, i)
		if e != nil {
			t.Error(e)
		}
		t.Log(n)
	}
}
