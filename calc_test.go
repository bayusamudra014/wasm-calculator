// Tests for calc.go. Run with `make test` (or `go test -v .`).
// They run on your machine, not in the browser, so they are a fast way to
// check your functions before building the .wasm file.
package main

import (
	"errors"
	"math"
	"testing"
)

type binaryCase struct {
	a, b, want float64
	wantErr    bool
}

type unaryCase struct {
	a, want float64
	wantErr bool
}

func checkBinary(t *testing.T, name string, fn binaryOp, cases []binaryCase) {
	t.Helper()
	for _, c := range cases {
		got, err := fn(c.a, c.b)
		check(t, name, []float64{c.a, c.b}, got, err, c.want, c.wantErr)
	}
}

func checkUnary(t *testing.T, name string, fn unaryOp, cases []unaryCase) {
	t.Helper()
	for _, c := range cases {
		got, err := fn(c.a)
		check(t, name, []float64{c.a}, got, err, c.want, c.wantErr)
	}
}

func check(t *testing.T, name string, in []float64, got float64, err error, want float64, wantErr bool) {
	t.Helper()
	if errors.Is(err, ErrNotImplemented) {
		t.Skipf("%s is not implemented yet", name)
	}
	switch {
	case wantErr && err == nil:
		t.Errorf("%s%v = %v, want an error", name, in, got)
	case !wantErr && err != nil:
		t.Errorf("%s%v returned error %q, want %v", name, in, err, want)
	case !wantErr && math.Abs(got-want) > 1e-9:
		t.Errorf("%s%v = %v, want %v", name, in, got, want)
	}
}

func TestAdd(t *testing.T) {
	checkBinary(t, "Add", Add, []binaryCase{
		{a: 2, b: 3, want: 5},
		{a: -4, b: 1.5, want: -2.5},
	})
}

func TestSubtract(t *testing.T) {
	checkBinary(t, "Subtract", Subtract, []binaryCase{
		{a: 10, b: 4, want: 6},
		{a: 1, b: 2.5, want: -1.5},
	})
}

func TestMultiply(t *testing.T) {
	checkBinary(t, "Multiply", Multiply, []binaryCase{
		{a: 6, b: 7, want: 42},
		{a: -2, b: 0.5, want: -1},
	})
}

func TestDivide(t *testing.T) {
	checkBinary(t, "Divide", Divide, []binaryCase{
		{a: 9, b: 3, want: 3},
		{a: 1, b: 4, want: 0.25},
		{a: 5, b: 0, wantErr: true},
	})
}

func TestPower(t *testing.T) {
	checkBinary(t, "Power", Power, []binaryCase{
		{a: 2, b: 10, want: 1024},
		{a: 9, b: 0.5, want: 3},
	})
}

func TestSqrt(t *testing.T) {
	checkUnary(t, "Sqrt", Sqrt, []unaryCase{
		{a: 16, want: 4},
		{a: 2, want: math.Sqrt2},
		{a: -1, wantErr: true},
	})
}

func TestPercent(t *testing.T) {
	checkUnary(t, "Percent", Percent, []unaryCase{
		{a: 50, want: 0.5},
		{a: 5, want: 0.05},
	})
}
