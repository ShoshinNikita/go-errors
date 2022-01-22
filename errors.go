package errors

import (
	"errors"
	"fmt"
)

type Error struct {
	errors []error
	pcs    []programCounters
}

var (
	_ error         = (*Error)(nil)
	_ fmt.Formatter = (*Error)(nil)
)

// TODO: comment
func addOrCreate(err, new error) *Error {
	// Skip addOrCreate and New/Errorf/Wrap/etc.
	const skip = 2

	pc := getProgramCounters(skip)

	if err == nil {
		return &Error{
			errors: []error{new},
			pcs:    []programCounters{pc},
		}
	}

	// Don't use errors.As to not lost any errors
	if e, ok := err.(*Error); ok {
		e.errors = append(e.errors, new)
		e.pcs = append(e.pcs, pc)
		return e
	}

	// Try to extract program counters
	pcs := []programCounters{pc}
	if e := (*Error)(nil); errors.As(err, &e) {
		pcs = append(pcs, e.pcs...)
	}
	return &Error{
		errors: []error{err, new},
		pcs:    pcs,
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
	var deepestStackTrace StackTrace
	for i := range e.pcs {
		s := e.pcs[i].toStackTrace()
		if len(s) > len(deepestStackTrace) {
			deepestStackTrace = s
		}
	}
	return deepestStackTrace
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
		f.Write([]byte(e.StackTrace().String()))
	}
}
