// Calculator math. THIS IS THE FILE YOU EDIT.
//
// Every function here is plain Go: no browser, no syscall/js. main.go wraps
// them and exposes them to JavaScript, and calc_test.go checks them with
// `make test`.
//
// Rules for each function:
//   - return (result, nil) on success
//   - return (0, someError) when the input is invalid (e.g. divide by zero);
//     the error message is shown on the calculator's display
//   - leave `return 0, ErrNotImplemented` until you have written it
package main

import "errors"

// ErrNotImplemented marks a function you still have to write.
var ErrNotImplemented = errors.New("not implemented yet")

// The two shapes of calculator function: two numbers in, or one number in.
type binaryOp func(a, b float64) (float64, error)
type unaryOp func(a float64) (float64, error)

// ---------------------------------------------------------------------------
// Example (done for you)
// ---------------------------------------------------------------------------

// Add returns a + b.
func Add(a, b float64) (float64, error) {
	return a + b, nil
}

// ---------------------------------------------------------------------------
// TODO 1: basic operations
// ---------------------------------------------------------------------------

// Subtract returns a - b.
func Subtract(a, b float64) (float64, error) {
	return 0, ErrNotImplemented
}

// Multiply returns a * b.
func Multiply(a, b float64) (float64, error) {
	return 0, ErrNotImplemented
}

// Divide returns a / b.
// Hint: return an error such as errors.New("cannot divide by zero") when b == 0.
func Divide(a, b float64) (float64, error) {
	return 0, ErrNotImplemented
}

// ---------------------------------------------------------------------------
// TODO 2: more operations
// ---------------------------------------------------------------------------

// Power returns a raised to the power b.
// Hint: import "math" and use math.Pow.
func Power(a, b float64) (float64, error) {
	return 0, ErrNotImplemented
}

// Sqrt returns the square root of a.
// Hint: return an error for negative numbers, then use math.Sqrt.
func Sqrt(a float64) (float64, error) {
	return 0, ErrNotImplemented
}

// Percent returns a / 100 (so 50 % becomes 0.5).
func Percent(a float64) (float64, error) {
	return 0, ErrNotImplemented
}

// ---------------------------------------------------------------------------
// TODO 3: your own function
// ---------------------------------------------------------------------------
// Write a new function here (for example Modulo, Square or Factorial), then
// follow "Add your own function" in README.md to register it and give it a
// button.
