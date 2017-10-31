package paginatedtypes

import "github.com/skycoin/bbs/src/misc/typ"

func NewSimple() typ.Paginated {
	return new(Simple)
}

type Simple struct {
	list []string
}

func (p *Simple) Append(v string) {
	p.list = append(p.list, v)
}

func (p *Simple) Get(in *typ.PaginatedInput) (*typ.PaginatedOutput, error) {
	out, e := typ.NewPaginatedOutput(in, uint(len(p.list)))
	if e != nil {
		return nil, e
	}

	var action func(v uint) uint
	if in.Reverse {
		action = func(v uint) uint { return v-1 }
	} else {
		action = func(v uint) uint { return v+1 }
	}
	for i, j := uint(0), in.StartIndex; i < out.ObtainedCount; i, j = i+1, action(j) {
		out.Data[i] = p.list[j]
	}

	return out, nil
}

func (p *Simple) Len() int {
	return len(p.list)
}

func (p *Simple) Clear() {
	p.list = []string{}
}