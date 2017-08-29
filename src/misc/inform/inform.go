package inform

import (
	"io"
	"log"
)

type empty struct{}

// Write is for empty.
func (e *empty) Write(_ []byte) (int, error) { return 0, nil }

// NewLogger creates a new logger.
func NewLogger(show bool, dst io.Writer, prefix string) *log.Logger {
	if !show {
		dst = &empty{}
	}
	return log.New(
		dst,
		"["+prefix+"] ",
		log.Ldate+log.Ltime+log.Lshortfile,
	)
}
