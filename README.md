# errors

`github.com/ShoshinNikita/go-errors` is like `github.com/pkg/errors` but with major improvements:

- it prints only the deepest stack trace
- it filters `runtime` calls in stack traces
- it wraps a special error (see `ErrNilWrap`) if `Wrap`/`Wrapf` is called with `nil` error

## Example

Consider the following code

```go
package main

import (
	"fmt"

	"github.com/pkg/errors"
)

func main() {
	fmt.Printf("%+v\n", foo())
}

func foo() error { return errors.Wrap(bar(), "foo") }
func bar() error { return errors.Wrap(errors.New("original err"), "bar") }
```

It prints 3 stack traces. And most of the stack frames are just repeated

```plain
original err
main.bar
	./main.go:14
main.foo
	./main.go:13
main.main
	./main.go:10
runtime.main
	runtime/proc.go:255
runtime.goexit
	runtime/asm_amd64.s:1581
bar
main.bar
	./main.go:14
main.foo
	./main.go:13
main.main
	./main.go:10
runtime.main
	runtime/proc.go:255
runtime.goexit
	runtime/asm_amd64.s:1581
foo
main.foo
	./main.go:13
main.main
	./main.go:10
runtime.main
	runtime/proc.go:255
runtime.goexit
	runtime/asm_amd64.s:1581
```

If we replace `github.com/pkg/errors` with `github.com/ShoshinNikita/go-errors`, we will get only the deepest stack trace. Additionally, all `runtime` calls will be filtered

```plain
foo: bar: original err
main.bar
	./main.go:14
main.foo
	./main.go:13
main.main
	./main.go:10
```

## TODO

- Add benchmarks
