package boo

import (
	"fmt"
)

// Type represents the type of boo.
type Type int

const (
	Unknown = Type(iota)
	Internal
	InvalidInput
	InvalidRead
	NotMaster
	NotAuthorised
	ObjectNotFound
	ObjectAlreadyExists
)

func Message(t Type) string {
	switch t {
	case Internal:
		return "An internal error has occurred."
	case InvalidInput:
		return "Received invalid input."
	case InvalidRead:
		return "Invalid read."
	case NotMaster:
		return "Action not permitted, node is not master."
	case NotAuthorised:
		return "You are not authorised to perform this action."
	case ObjectNotFound:
		return "Object is not found."
	case ObjectAlreadyExists:
		return "Object already exists."
	default:
		return "An unknown error has occurred."
	}
}

type root struct {
	details string
	what    Type
}

func (r *root) Error() string {
	return r.details
}

func New(what Type, v ...interface{}) error {
	return &root{
		details: fmt.Sprintln(v...),
		what:    what,
	}
}

func Newf(what Type, format string, v ...interface{}) error {
	return &root{
		details: fmt.Sprintf(format, v...),
		what:    what,
	}
}

func What(e error) Type {
	if r, ok := e.(*root); ok {
		return r.what
	}
	return Unknown
}
