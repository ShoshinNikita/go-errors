package errors

import (
	"errors"
	"runtime"
	"strconv"
	"strings"
)

type StackTrace []Frame

func (s StackTrace) String() string {
	if len(s) == 0 {
		return ""
	}

	var res string
	for _, frame := range s {
		res += frame.String() + "\n"
	}
	return res[:len(res)-1]
}

type Frame struct {
	Funcion string
	File    string
	Line    int
}

func (f Frame) String() string {
	return f.Funcion + "\n\t" + f.File + ":" + strconv.Itoa(f.Line)
}

func ExtractStackTrace(err error) StackTrace {
	var e *Error
	if errors.As(err, &e) {
		return e.StackTrace()
	}
	return nil
}

type programCounters []uintptr

func getProgramCounters(skip int) programCounters {
	// Skip runtime.Callers and getProgramCounters
	skip += 2

	var pc [32]uintptr
	n := runtime.Callers(skip, pc[:])
	return pc[:n]
}

// toStackTrace converts program counters to a list of frames and filters all runtime call
func (c programCounters) toStackTrace() StackTrace {
	res := make([]Frame, 0, len(c))

	frames := runtime.CallersFrames(c)
	for {
		f, ok := frames.Next()
		if !ok {
			break
		}
		if strings.HasPrefix(f.Function, "runtime.") {
			break
		}

		res = append(res, Frame{
			Funcion: f.Function,
			File:    f.File,
			Line:    f.Line,
		})
	}
	return res
}
