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
	MaxCount   uint `json:"max_count"`   // max number of elements to get.
	Reverse    bool `json:"reverse"`     // whether to get elements in the opposite direction.
}

type PaginatedOutput struct {
	StartIndex     uint     `json:"start_index"`
	ObtainedCount  uint     `json:"obtained_count"`
	RemainingCount uint     `json:"remaining_count"`
	IsReversed     bool     `json:"is_reversed"`
	Data           []string `json:"-"`
}

func NewPaginatedOutput(in *PaginatedInput, dataCount uint) (*PaginatedOutput, error) {
	if in.StartIndex < 0 || in.StartIndex >= dataCount {
		return nil, boo.Newf(boo.InvalidInput,
			"invalid 'start_index' provided, valid values are between %d and %d inclusive",
			0, dataCount-1)
	}
	if in.MaxCount <= 0 {
		return nil, boo.New(boo.InvalidInput,
			"invalid 'max_count' provided, valid values are in range '>= 0'")
	}

	var obtainedCount uint
	if in.Reverse {
		if in.MaxCount > in.StartIndex {
			obtainedCount = in.StartIndex + 1
		} else {
			obtainedCount = in.MaxCount
		}
	} else {
		if in.StartIndex+in.MaxCount > dataCount {
			obtainedCount = dataCount - in.StartIndex
		} else {
			obtainedCount = in.MaxCount
		}
	}

	var remainingCount uint
	if in.Reverse {
		remainingCount = in.StartIndex + 1 - obtainedCount
	} else {
		remainingCount = dataCount - in.StartIndex - obtainedCount
	}

	return &PaginatedOutput{
		StartIndex:     in.StartIndex,
		ObtainedCount:  obtainedCount,
		RemainingCount: remainingCount,
		IsReversed:     in.Reverse,
		Data:           make([]string, obtainedCount),
	}, nil
}
