package typ

import "github.com/skycoin/bbs/src/misc/boo"

type PaginatedCreator func() Paginated

type Paginated interface {
	Append(v string)
	Has(v string) bool
	Get(in *PaginatedInput) (*PaginatedOutput, error)
	Len() int
	Clear()
}

type PaginatedInput struct {
	StartIndex uint `json:"start_index"` // index to start with.
	PageSize   uint `json:"page_size"`   // max number of elements in a page.
	Reverse    bool `json:"reverse"`     // whether to get elements in the opposite direction.
}

type PaginatedOutput struct {
	RecordCount uint     `json:"record_count"`
	StartIndex  uint     `json:"start_index"`
	PageSize    uint     `json:"page_size"`
	IsReversed  bool     `json:"is_reversed"`
	Data        []string `json:"-"`
}

func NewPaginatedOutput(in *PaginatedInput, dataCount uint) (*PaginatedOutput, error) {
	if in.StartIndex < 0 || in.StartIndex >= dataCount {
		return nil, boo.Newf(boo.InvalidInput,
			"invalid 'start_index' provided, valid values are between %d and %d inclusive",
			0, dataCount-1)
	}
	if in.PageSize <= 0 {
		return nil, boo.New(boo.InvalidInput,
			"invalid 'max_count' provided, valid values are in range '>= 0'")
	}

	var obtainedCount uint
	if in.Reverse {
		if in.PageSize > in.StartIndex {
			obtainedCount = in.StartIndex + 1
		} else {
			obtainedCount = in.PageSize
		}
	} else {
		if in.StartIndex+in.PageSize > dataCount {
			obtainedCount = dataCount - in.StartIndex
		} else {
			obtainedCount = in.PageSize
		}
	}

	return &PaginatedOutput{
		RecordCount: obtainedCount,
		StartIndex:  in.StartIndex,
		PageSize:    in.PageSize,
		IsReversed:  in.Reverse,
		Data:        make([]string, obtainedCount),
	}, nil
}
