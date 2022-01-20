package errors

import (
	"errors"
	"fmt"
)

// ErrNilWrap is returned when Wrap/Wrapf is called with nil error
var ErrNilWrap = errors.New("<nil>")

func New(text string) error {
	return addOrCreate(nil, errors.New(text),)
}

func Errorf(format string, a ...interface{}) error {
	return addOrCreate(nil, fmt.Errorf(format, a...),)
}

func Wrap(err error, msg string) error {
	if err == nil {
		err = ErrNilWrap
	}
	return addOrCreate(err, errors.New(msg))
}

func Wrapf(err error, format string, a ...interface{}) error {
	if err == nil {
		err = ErrNilWrap
	}
	return addOrCreate(err, fmt.Errorf(format, a...))
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}
