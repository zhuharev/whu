package errors

import "fmt"

type ErrorType int

func (e ErrorType) New(message string) *Error {
	return New(e, message)
}

// Error implements error interface
type Error struct {
	typ ErrorType
	msg string

	// stash contain line, looks like `key=value key2=value2`
	stash string
}

func New(errType ErrorType, message string) *Error {
	return &Error{
		typ: errType,
		msg: message,
	}
}

func (e *Error) Error() string {
	if e.stash == "" {
		return e.msg
	}
	return fmt.Sprintf("%s %s", e.msg, e.stash)
}

func (e *Error) Is(err error) bool {
	if x, ok := err.(*Error); ok && x.typ == e.typ {
		return true
	}
	return false
}

func (e *Error) String(key string, value string) *Error {
	field := fmt.Sprintf("%s=%s", key, value)
	if e.stash != "" {
		e.stash += " "
	}
	e.stash += field
	return e
}

func (e *Error) Int(key string, value int) *Error {
	field := fmt.Sprintf("%s=%d", key, value)
	if e.stash != "" {
		e.stash += " "
	}
	e.stash += field
	return e
}

func Wrap(err error, message string) *Error {
	if x, ok := err.(*Error); ok {
		x.msg = fmt.Sprintf("%s: %s", message, x.msg)
		return x
	}
	return New(0, message)
}
