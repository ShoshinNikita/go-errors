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

// TODO: comment about "runtime"
func getStackTrace(skip int) StackTrace {
	const depth = 32

	// Skip runtime.Callers and GetStackTrace
	skip += 2

	pc := make([]uintptr, depth)
	n := runtime.Callers(skip, pc)
	pc = pc[:n]

	res := make([]Frame, 0, len(pc))

	frames := runtime.CallersFrames(pc)
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

func ExtractStackTrace(err error) StackTrace {
	var e *Error
	if errors.As(err, &e) {
		return e.StackTrace()
	}
	return nil
}
