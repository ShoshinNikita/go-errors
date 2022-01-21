package tests

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/ShoshinNikita/go-errors"
	pkgerrors "github.com/pkg/errors"
)

var errForBenchmarks error

func BenchmarkWrap(b *testing.B) {
	b.Run("go-errors", func(b *testing.B) {
		benchmarkWrap(b, goErrorsWrapper{})
	})
	b.Run("fmt", func(b *testing.B) {
		benchmarkWrap(b, fmtWrapper{})
	})
	b.Run("pkg_errors", func(b *testing.B) {
		benchmarkWrap(b, pkgErrorsWrapper{})
	})
}

func benchmarkWrap(b *testing.B, w interface {
	Wrap2() error
	Wrap5() error
	Wrap10() error
}) {
	b.Run("few (2)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			errForBenchmarks = w.Wrap2()
		}
	})
	b.Run("medium (5)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			errForBenchmarks = w.Wrap5()
		}
	})
	b.Run("many (10)", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			errForBenchmarks = w.Wrap10()
		}
	})
}

func BenchmarkFormat(b *testing.B) {
	const n = 10

	b.Run("go-errors", func(b *testing.B) {
		err := goErrorsWrapper{}.Wrap10()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			fmt.Fprintf(ioutil.Discard, "%+v", err)
		}
	})
	b.Run("fmt", func(b *testing.B) {
		err := fmtWrapper{}.Wrap10()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			fmt.Fprintf(ioutil.Discard, "%+v", err)
		}
	})
	b.Run("pkg_errors", func(b *testing.B) {
		err := pkgErrorsWrapper{}.Wrap10()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			fmt.Fprintf(ioutil.Discard, "%+v", err)
		}
	})
}

type goErrorsWrapper struct{}

func (e goErrorsWrapper) Wrap1() error  { return errors.Wrap(sql.ErrNoRows, "1") }
func (e goErrorsWrapper) Wrap2() error  { return errors.Wrap(e.Wrap1(), "2") }
func (e goErrorsWrapper) Wrap3() error  { return errors.Wrap(e.Wrap2(), "3") }
func (e goErrorsWrapper) Wrap4() error  { return errors.Wrap(e.Wrap3(), "4") }
func (e goErrorsWrapper) Wrap5() error  { return errors.Wrap(e.Wrap4(), "5") }
func (e goErrorsWrapper) Wrap6() error  { return errors.Wrap(e.Wrap5(), "6") }
func (e goErrorsWrapper) Wrap7() error  { return errors.Wrap(e.Wrap6(), "7") }
func (e goErrorsWrapper) Wrap8() error  { return errors.Wrap(e.Wrap7(), "8") }
func (e goErrorsWrapper) Wrap9() error  { return errors.Wrap(e.Wrap8(), "9") }
func (e goErrorsWrapper) Wrap10() error { return errors.Wrap(e.Wrap9(), "10") }

type pkgErrorsWrapper struct{}

func (e pkgErrorsWrapper) Wrap1() error  { return pkgerrors.Wrap(sql.ErrNoRows, "1") }
func (e pkgErrorsWrapper) Wrap2() error  { return pkgerrors.Wrap(e.Wrap1(), "2") }
func (e pkgErrorsWrapper) Wrap3() error  { return pkgerrors.Wrap(e.Wrap2(), "3") }
func (e pkgErrorsWrapper) Wrap4() error  { return pkgerrors.Wrap(e.Wrap3(), "4") }
func (e pkgErrorsWrapper) Wrap5() error  { return pkgerrors.Wrap(e.Wrap4(), "5") }
func (e pkgErrorsWrapper) Wrap6() error  { return pkgerrors.Wrap(e.Wrap5(), "6") }
func (e pkgErrorsWrapper) Wrap7() error  { return pkgerrors.Wrap(e.Wrap6(), "7") }
func (e pkgErrorsWrapper) Wrap8() error  { return pkgerrors.Wrap(e.Wrap7(), "8") }
func (e pkgErrorsWrapper) Wrap9() error  { return pkgerrors.Wrap(e.Wrap8(), "9") }
func (e pkgErrorsWrapper) Wrap10() error { return pkgerrors.Wrap(e.Wrap9(), "10") }

type fmtWrapper struct{}

func (e fmtWrapper) Wrap1() error  { return fmt.Errorf("1: %w", sql.ErrNoRows) }
func (e fmtWrapper) Wrap2() error  { return fmt.Errorf("2: %w", e.Wrap1()) }
func (e fmtWrapper) Wrap3() error  { return fmt.Errorf("3: %w", e.Wrap2()) }
func (e fmtWrapper) Wrap4() error  { return fmt.Errorf("4: %w", e.Wrap3()) }
func (e fmtWrapper) Wrap5() error  { return fmt.Errorf("5: %w", e.Wrap4()) }
func (e fmtWrapper) Wrap6() error  { return fmt.Errorf("6: %w", e.Wrap5()) }
func (e fmtWrapper) Wrap7() error  { return fmt.Errorf("7: %w", e.Wrap6()) }
func (e fmtWrapper) Wrap8() error  { return fmt.Errorf("8: %w", e.Wrap7()) }
func (e fmtWrapper) Wrap9() error  { return fmt.Errorf("9: %w", e.Wrap8()) }
func (e fmtWrapper) Wrap10() error { return fmt.Errorf("10: %w", e.Wrap9()) }
