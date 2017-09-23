package boo

import (
	"fmt"
)

const (
	Unknown = iota
	Internal
	InvalidInput
	InvalidRead
	NotMaster
	NotAuthorised
	NotAllowed
	NotFound
	AlreadyExists
)

// Message returns the error message of type 't'.
func Message(t int) string {
	switch t {
	case Internal:
		return "Internal Error"
	case InvalidInput:
		return "Invalid Input"
	case InvalidRead:
		return "Invalid Read"
	case NotMaster:
		return "Node is Not Master"
	case NotAuthorised:
		return "Not Authorised"
	case NotAllowed:
		return "Not Allowed"
	case NotFound:
		return "Not Found"
	case AlreadyExists:
		return "Already Exists"
	default:
		return "Unknown Error"
	}
}

// This satisfies the 'error' interface.
type elem struct {
	e error
	m string
	t int
}

func (e *elem) dig(msg *string, typ *int) {
	*msg += ": " + e.m
	if *typ == Unknown {
		*typ = e.t
	}
	if e.e == nil {
		return
	}
	switch e.e.(type) {
	case *elem:
		e.e.(*elem).dig(msg, typ)
	default:
		*msg += ": " + e.e.Error()
	}
}

// Error returns the error message.
func (e *elem) Error() string {
	msg, typ := "", Unknown
	e.dig(&msg, &typ)
	return msg[2:]
}

// New creates a new error of type 't'.
func New(t int, m ...interface{}) error {
	return &elem{e: nil, m: fmt.Sprint(m...), t: t}
}

// Newf creates a new error of type 't' with formatted text.
func Newf(t int, f string, v ...interface{}) error {
	return &elem{e: nil, m: fmt.Sprintf(f, v...), t: t}
}

// Wrap wraps an error with an additional message.
func Wrap(e error, m string) error {
	if e == nil {
		return nil
	}
	return &elem{e: e, m: m, t: Unknown}
}

// Wrapf wraps an error with an additional formatted message.
func Wrapf(e error, f string, v ...interface{}) error {
	if e == nil {
		return nil
	}
	return &elem{e: e, m: fmt.Sprintf(f, v...), t: Unknown}
}

// WrapType wraps an error with a type and an additional message.
func WrapType(e error, t int, m ...interface{}) error {
	if e == nil {
		return nil
	}
	return &elem{e: e, m: fmt.Sprint(m...), t: t}
}

// WrapTypef wraps an error with a type and an additional formatted message.
func WrapTypef(e error, t int, f string, v ...interface{}) error {
	if e == nil {
		return nil
	}
	return &elem{e: e, m: fmt.Sprintf(f, v...), t: t}
}

// Type returns the type of the error.
func Type(e error) int {
	v, ok := e.(*elem)
	if !ok {
		return Unknown
	}

	msg, typ := "", Unknown
	v.dig(&msg, &typ)
	return typ
}
