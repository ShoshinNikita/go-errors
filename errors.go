package errors

import (
	"errors"
	"fmt"
)

// Errors implements error interface and methods Is and As to be compatible with errors.Is and errors.As
type Error struct {
	errors []error
	// pcs is a list of program counters with equal and maximum size. This state is maintained by addNewPCs
	pcs []programCounters
}

var (
	_ error         = (*Error)(nil)
	_ fmt.Formatter = (*Error)(nil)
)

// newError creates a new Error or reuse the existing one
func newError(err, new error) *Error {
	// Skip newError and New/Errorf/Wrap/etc.
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
		addNewPCs(&e.pcs, pc)
		return e
	}

	// Try to extract program counters
	pcs := []programCounters{pc}
	if e := (*Error)(nil); errors.As(err, &e) {
		addNewPCs(&pcs, e.pcs...)
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

// addNewPCs adds new program counters reusing pcs. After the call, pcs will contain elements
// with the maximum number of program counters, and all elements will have the same size
func addNewPCs(pcs *[]programCounters, newPCs ...programCounters) {
	if len(*pcs) == 0 {
		l := len(newPCs)
		if l == 0 {
			// Just in case
			return
		}

		*pcs = append(*pcs, newPCs[0])
		newPCs = newPCs[1:]
		if l == 1 {
			// Nothing to filter
			return
		}
	}

	maxDepth := len((*pcs)[0])
	for _, pc := range newPCs {
		if l := len(pc); l > maxDepth {
			maxDepth = l
		}
	}
	if maxDepth > len((*pcs)[0]) {
		// All current pcs are less
		*pcs = (*pcs)[:0]
	}
	for _, pc := range newPCs {
		if len(pc) == maxDepth {
			*pcs = append(*pcs, pc)
		}
	}
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
