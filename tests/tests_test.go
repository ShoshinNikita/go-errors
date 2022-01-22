package tests

import (
	"database/sql"
	stderrors "errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/ShoshinNikita/go-errors"
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

const modulePath = "github.com/ShoshinNikita/go-errors/tests"

var (
	ErrGlobal = errors.New("global error")
)

func TestStackTraceRestore(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	a := func() error { return errors.Errorf("some %q", "error") }
	b := func() error { return pkgerrors.Wrap(a(), "b") }
	c := func() error { return errors.Wrapf(b(), "%s", "c") }
	d := func() error { return fmt.Errorf("d: %w", c()) }
	e := func() error { return errors.Wrap(d(), "e") } // Stack trace should be restored
	err := e()

	require.EqualError(err, `e: d: c: b: some "error"`)
	trace := errors.ExtractStackTrace(err)
	checkStackTraces(t, trace,
		"TestStackTraceRestore.func1", // wrap in a
		"TestStackTraceRestore.func2", // wrap in b
		"TestStackTraceRestore.func3", // wrap in c
		"TestStackTraceRestore.func4", // wrap in d
		"TestStackTraceRestore.func5", // wrap in e
		"TestStackTraceRestore",       // call to e
	)
}

func TestOverwriteGlobalErrorStackTrace(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	a := func() error { return errors.Wrap(ErrGlobal, "a") }
	b := func() error { return errors.Wrap(a(), "b") }
	err := b()

	require.EqualError(err, "b: a: global error")
	trace := errors.ExtractStackTrace(err)
	checkStackTraces(t, trace,
		"TestOverwriteGlobalErrorStackTrace.func1", // wrap in a
		"TestOverwriteGlobalErrorStackTrace.func2", // wrap in b
		"TestOverwriteGlobalErrorStackTrace",       // call to b
	)
}

func TestWrapNil(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	a := func() error { return nil }
	b := func() error { return errors.Wrap(a(), "b") }
	c := func() error { return errors.Wrapf(a(), "%d", 1) }

	err := b()
	require.EqualError(err, "b: <nil>")
	require.True(errors.Is(err, errors.ErrNilWrap))
	trace := errors.ExtractStackTrace(err)
	checkStackTraces(t, trace,
		"TestWrapNil.func2", // wrap in b
		"TestWrapNil",       // call to b
	)

	err = c()
	require.EqualError(err, "1: <nil>")
	require.True(errors.Is(err, errors.ErrNilWrap))
	trace = errors.ExtractStackTrace(err)
	checkStackTraces(t, trace,
		"TestWrapNil.func3", // wrapf in c
		"TestWrapNil",       // call to c
	)
}

// TODO: comment
func checkStackTraces(t *testing.T, trace errors.StackTrace, expectedFuncs ...string) {
	t.Helper()

	for i := range expectedFuncs {
		expectedFuncs[i] = modulePath + "." + expectedFuncs[i]
	}

	// All stack traces have call to testing.tRunner
	expectedFuncs = append(expectedFuncs, "testing.tRunner")

	if len(expectedFuncs) != len(trace) {
		t.Fatalf(
			"stack trace has wrong number of functions: expected %d, got %d, stack trace: %s",
			len(expectedFuncs), len(trace), trace,
		)
	}

	funcs := make([]string, 0, len(trace))
	for _, frame := range trace {
		funcs = append(funcs, frame.Function)
	}

	for i := range funcs {
		if expectedFuncs[i] != funcs[i] {
			t.Fatalf("wrong #%d stack trace frame: %q, got %q, stack trace: %s", i, expectedFuncs[i], funcs[i], trace)
		}
	}
}

func TestIs(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	err := errors.Wrap(sql.ErrNoRows, "wrap")
	err = errors.Wrapf(err, "wrap%s", "f")
	err = fmt.Errorf("fmt: %w", err)
	err = errors.Wrap(err, "wrap")

	require.True(errors.Is(err, sql.ErrNoRows))
	require.False(errors.Is(err, sql.ErrTxDone))

	require.True(stderrors.Is(err, sql.ErrNoRows))
	require.False(stderrors.Is(err, sql.ErrTxDone))
}

func TestAs(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	err := errors.Wrap(&MyError{}, "wrap")
	err = errors.Wrapf(err, "wrap%s", "f")
	err = fmt.Errorf("fmt: %w", err)
	err = &os.PathError{Op: "op", Path: "path", Err: err}
	err = errors.Wrap(err, "wrap")

	var (
		pathError *os.PathError
		myError   *MyError
		linkError *os.LinkError
	)
	require.True(errors.As(err, &pathError))
	require.True(errors.As(err, &myError))
	require.False(errors.As(err, &linkError))

	require.True(stderrors.As(err, &pathError))
	require.True(stderrors.As(err, &myError))
	require.False(stderrors.As(err, &linkError))
}

type MyError struct{}

func (*MyError) Error() string {
	return "my error"
}

func TestFormat(t *testing.T) {
	t.Parallel()

	require := require.New(t)

	const errorMsg = "func4: func3: func2: func1"

	modulePath := strings.ReplaceAll(modulePath, ".", "\\.")

	_, filepath, _, ok := runtime.Caller(0)
	require.True(ok)
	filepath = strings.ReplaceAll(filepath, ".", "\\.")

	fullErrorPattern := "^" + errorMsg + "\n"
	for _, f := range []string{"func1", "func2", "func3", "func4", "TestFormat"} {
		fullErrorPattern += modulePath + "\\." + f + "\n\t" + filepath + ":\\d+" + "\n"
	}
	fullErrorPattern += "testing.tRunner" + "\n\t" + ".*/src/testing/testing.go:\\d+$"

	err := func4()
	require.Equal("func4: func3: func2: func1", fmt.Sprint(err))
	require.Equal("func4: func3: func2: func1", fmt.Sprintf("%s", err))
	require.Equal("func4: func3: func2: func1", fmt.Sprintf("%v", err))
	require.Equal("func4: func3: func2: func1", fmt.Sprintf("%q", err))
	require.Regexp(fullErrorPattern, fmt.Sprintf("%+v", err))
	require.Regexp(fullErrorPattern, fmt.Sprintf("%#v", err))
}

func func1() error { return errors.New("func1") }
func func2() error { return errors.Wrap(func1(), "func2") }
func func3() error { return errors.Wrap(func2(), "func3") }
func func4() error { return errors.Wrap(func3(), "func4") }
