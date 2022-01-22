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

## Benchmarks

You can use this command `BENCH=. make bench` to run all benchmarks

```plain
name                           time/op
Wrap/go-errors/few (2)         4.92µs ± 7%
Wrap/go-errors/medium (5)      12.5µs ± 3%
Wrap/go-errors/many (10)       27.3µs ± 2%

Wrap/pkg_errors/few (2)        4.27µs ± 8%
Wrap/pkg_errors/medium (5)     11.5µs ± 8%
Wrap/pkg_errors/many (10)      25.6µs ± 2%

Wrap/fmt/few (2)               1.01µs ± 4%
Wrap/fmt/medium (5)            2.61µs ± 2%
Wrap/fmt/many (10)             5.43µs ± 9%

Format/go-errors               18.3µs ± 5%
Format/pkg_errors               108µs ± 2%
Format/fmt                      208ns ± 2%

name                           alloc/op
Wrap/go-errors/few (2)           720B ± 0%
Wrap/go-errors/medium (5)      1.66kB ± 0%
Wrap/go-errors/many (10)       3.28kB ± 0%

Wrap/pkg_errors/few (2)          672B ± 0%
Wrap/pkg_errors/medium (5)     1.68kB ± 0%
Wrap/pkg_errors/many (10)      3.36kB ± 0%

Wrap/fmt/few (2)                 128B ± 0%
Wrap/fmt/medium (5)              368B ± 0%
Wrap/fmt/many (10)               816B ± 0%

Format/go-errors               17.9kB ± 0%
Format/pkg_errors              1.42kB ± 0%
Format/fmt                      0.00B     

name                           allocs/op
Wrap/go-errors/few (2)           9.00 ± 0%
Wrap/go-errors/medium (5)        16.0 ± 0%
Wrap/go-errors/many (10)         27.0 ± 0%

Wrap/pkg_errors/few (2)          8.00 ± 0%
Wrap/pkg_errors/medium (5)       20.0 ± 0%
Wrap/pkg_errors/many (10)        40.0 ± 0%

Wrap/fmt/few (2)                 4.00 ± 0%
Wrap/fmt/medium (5)              10.0 ± 0%
Wrap/fmt/many (10)               20.0 ± 0%

Format/go-errors                 53.0 ± 0%
Format/pkg_errors                 155 ± 0%
Format/fmt                       0.00     
```

