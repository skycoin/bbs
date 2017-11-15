package paginatedtypes

import (
	"fmt"
	"github.com/skycoin/bbs/src/misc/typ"
	"testing"
)

func TestSimple_Get(t *testing.T) {
	const count = 10

	p := NewSimple()
	for i := 0; i < count; i++ {
		p.Append(fmt.Sprintf("data_index(%d)", i))
	}

	t.Run("expect error when start_index > data_count", func(t *testing.T) {
		if _, e := p.Get(&typ.PaginatedInput{StartIndex: count + 1}); e != nil {
			t.Log("error as expected:", e)
		} else {
			t.Error("expected error when start_index > data_count")
		}
	})

	t.Run("expect error when max_count == 0", func(t *testing.T) {
		if _, e := p.Get(&typ.PaginatedInput{PageSize: 0}); e != nil {
			t.Log("error as expected:", e)
		} else {
			t.Error("expected error when max_count == 0")
		}
	})

	t.Run("get", func(t *testing.T) {
		out, e := p.Get(&typ.PaginatedInput{
			StartIndex: 2,
			PageSize:   6,
			Reverse:    false,
		})
		if e != nil {
			t.Error(e)
		} else {
			t.Log(out.Data, out.RemainingCount)
		}
	})
}
