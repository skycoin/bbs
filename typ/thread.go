package typ

import (
	"strings"
	"errors"
)

// Thread represents a thread stored in cxo.
type Thread struct {
	Name        string `json:"name"`
	Desc        string `json:"description"`
	MasterBoard string `json:"master_board"`
	Hash        string `json:"hash" enc:"-"`
}

func (t *Thread) Check() error {
	if t == nil {
		return errors.New("thread is nil")
	}
	name := strings.TrimSpace(t.Name)
	desc := strings.TrimSpace(t.Desc)
	if len(name) < 3 && len(desc) < 3 {
		return errors.New("thread content too short")
	}
	return nil
}