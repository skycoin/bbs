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
	ObjectNotFound
	ObjectAlreadyExists
)

func Message(t int) string {
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

func (e *elem) Error() string {
	msg, typ := "", Unknown
	e.dig(&msg, &typ)
	return msg[2:]
}

func New(t int, m string) error {
	return &elem{e: nil, m: m, t: t}
}

func Newf(t int, f string, v ...interface{}) error {
	return &elem{e: nil, m: fmt.Sprintf(f, v...), t: t}
}

func Wrap(e error, m string) error {
	return &elem{e: e, m: m, t: Unknown}
}

func Wrapf(e error, f string, v ...interface{}) error {
	return &elem{e: e, m: fmt.Sprintf(f, v...), t: Unknown}
}

func WrapType(e error, t int, m string) error {
	return &elem{e: e, m: m, t: t}
}

func WrapTypef(e error, t int, f string, v ...interface{}) error {
	return &elem{e: e, m: fmt.Sprintf(f, v...), t: t}
}

func Type(e error) int {
	switch e.(type) {
	case *elem:
		msg, typ := "", Unknown
		e.(*elem).dig(&msg, &typ)
		return typ
	default:
		return Unknown
	}
}
