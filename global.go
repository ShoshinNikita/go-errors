package errors

import (
	"errors"
	"fmt"
)

func New(text string) error {
	return addOrCreate(nil, errors.New(text), 0)
}

func Errorf(format string, a ...interface{}) error {
	return addOrCreate(nil, fmt.Errorf(format, a...), 0)
}

func Wrap(err error, msg string) error {
	return addOrCreate(err, errors.New(msg), 0)
}

func Wrapf(err error, format string, a ...interface{}) error {
	return addOrCreate(err, fmt.Errorf(format, a...), 0)
}

func WrapWithSkip(skip int, err error, msg string) error {
	return addOrCreate(err, errors.New(msg), skip)
}

func WrapfWithSkip(skip int, err error, format string, a ...interface{}) error {
	return addOrCreate(err, fmt.Errorf(format, a...), skip)
}
