package errors

import (
	"errors"
	"fmt"
)

type Error struct {
	errors []error
	stack  StackTrace
}

var (
	_ error         = (*Error)(nil)
	_ fmt.Formatter = (*Error)(nil)
)

// TODO: comment
func addOrCreate(err, new error, additionalSkips int) *Error {
	// Skip addOrCreate and New/Errorf/Wrap/etc.
	skip := 2 + additionalSkips

	stack := getStackTrace(skip)

	if err == nil {
		return &Error{
			errors: []error{new},
			stack:  stack,
		}
	}

	// Don't use errors.As to not lost any errors
	if e, ok := err.(*Error); ok {
		e.errors = append(e.errors, new)
		if len(e.stack) < len(stack) {
			e.stack = stack
		}
		return e
	}

	// Try to extract the deepest stack
	var e *Error
	if errors.As(err, &e) && len(e.stack) >= len(stack) {
		stack = e.stack
	}
	return &Error{
		errors: []error{err, new},
		stack:  stack,
	}
}

func (e *Error) Error() string {
	var res string
	for i := len(e.errors) - 1; i >= 0; i-- {
		res += e.errors[i].Error()
		if i > 0 {
			res += ": "
		}
	}
	return res
}

func (e *Error) StackTrace() StackTrace {
	res := make(StackTrace, len(e.stack))
	copy(res, e.stack)
	return res
}

func (e *Error) Is(target error) bool {
	for _, err := range e.errors {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func (e *Error) As(target interface{}) bool {
	for _, err := range e.errors {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}

func (e *Error) Format(f fmt.State, verb rune) {
	f.Write([]byte(e.Error()))
	if verb == 'v' && (f.Flag('+') || f.Flag('#')) {
		f.Write([]byte("\n"))
		f.Write([]byte(e.stack.String()))
	}
}
