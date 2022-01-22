package errors

import (
	"errors"
	"fmt"
)

// ErrNilWrap is used when Wrap/Wrapf is called with nil error
var ErrNilWrap = errors.New("<nil>")

// New is like error.New but adds a stack trace to the error
func New(text string) error {
	return newError(nil, errors.New(text))
}

// New is like fmt.Errorf but adds a stack trace to the error
func Errorf(format string, a ...interface{}) error {
	return newError(nil, fmt.Errorf(format, a...))
}

// Wrap wraps the original error and adds a stack trace if needed. If Wrap is called with nil error,
// ErrNilWrap will be wrapped instead
func Wrap(err error, msg string) error {
	if err == nil {
		err = ErrNilWrap
	}
	return newError(err, errors.New(msg))
}

// Wrapf is like Wrap but formats the message
func Wrapf(err error, format string, a ...interface{}) error {
	if err == nil {
		err = ErrNilWrap
	}
	return newError(err, fmt.Errorf(format, a...))
}

// Is is a wrapper for errors.Is
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As is a wrapper for errors.As
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Unwrap is a wrapper for errors.Unwrap
func Unwrap(err error) error {
	return errors.Unwrap(err)
}
