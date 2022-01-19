package errors_test

import (
	"fmt"

	"github.com/ShoshinNikita/errors"
)

var ErrGlobal = errors.New("global error")

func foo() error {
	return errors.Wrapf(bar(), "wrap in %q", "foo")
}

func bar() error {
	return fmt.Errorf("bar: %w", xyz())
}

func xyz() error {
	return errors.Wrap(ErrGlobal, "xyz")
}

func Example() {
	err := foo()

	fmt.Println(err)

	// Output:
	// wrap in "foo": bar: xyz: global error
}
