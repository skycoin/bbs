package typ

import (
	"errors"
	"strings"
	"time"
)

// Thread represents a thread stored in cxo.
type Thread struct {
	Name    string `json:"name"`
	Desc    string `json:"description"`
	Created int64  `json:"created"`
}

func NewThread(name, desc string) *Thread {
	return &Thread{
		Name:    name,
		Desc:    desc,
		Created: time.Now().UnixNano(),
	}
}

func (t *Thread) CheckAndPrep() error {
	if t == nil {
		return errors.New("nil thread")
	}
	t.Name = strings.TrimSpace(t.Name)
	t.Desc = strings.TrimSpace(t.Desc)
	t.Created = time.Now().UnixNano()
	if len(t.Name) == 0 {
		return errors.New("invalid thread name")
	}
	return nil
}
