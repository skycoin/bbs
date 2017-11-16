package paginatedtypes

import "github.com/skycoin/bbs/src/misc/typ"

func NewMapped() typ.Paginated {
	return &Mapped{
		dict: make(map[string]struct{}),
	}
}

type Mapped struct {
	list []string
	dict map[string]struct{}
}

func (p *Mapped) Append(v string) {
	if p.Has(v) {
		return
	}
	p.dict[v] = struct{}{}
	p.list = append(p.list, v)
}

func (p *Mapped) Has(v string) bool {
	_, ok := p.dict[v]
	return ok
}

func (p *Mapped) Get(in *typ.PaginatedInput) (*typ.PaginatedOutput, error) {
	out, e := typ.NewPaginatedOutput(in, uint(len(p.list)))
	if e != nil {
		return nil, e
	}

	var action func(v uint) uint
	if in.Reverse {
		action = func(v uint) uint { return v - 1 }
	} else {
		action = func(v uint) uint { return v + 1 }
	}
	for i, j := uint(0), in.StartIndex; i < uint(len(out.Data)); i, j = i+1, action(j) {
		out.Data[i] = p.list[j]
	}

	return out, nil
}

func (p *Mapped) Len() int {
	return len(p.list)
}

func (p *Mapped) Clear() {
	p.list = []string{}
	p.dict = make(map[string]struct{})
}
