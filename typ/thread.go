package typ

import (
	"time"
)

// Thread represents a thread.
type Thread struct {
	ID           []byte `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Created      int64  `json:"created"`
	LastModified int64  `json:"last_modified"`
	Version      uint64 `json:"version"`
}

// NewThread creates a new thread.
func NewThread(name, desc string) *Thread {
	now := time.Now().UnixNano()
	thread := Thread{
		ID:           MakeTimeStampedRandomID(128),
		Name:         name,
		Description:  desc,
		Created:      now,
		LastModified: now,
		Version:      0,
	}
	return &thread
}
