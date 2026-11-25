//go:build js && wasm

// Bridge between JavaScript and the Go math in calc.go.
//
// main() puts an object on window.calcWasm whose methods call the Go
// functions. Each method returns a plain JS object:
//
//	{ ok: true,  value: 42 }
//	{ ok: false, error: "cannot divide by zero" }
package main

import "syscall/js"

// result turns a Go (value, error) pair into the JS object described above.
func result(v float64, err error) any {
	if err != nil {
		return map[string]any{"ok": false, "error": err.Error()}
	}
	return map[string]any{"ok": true, "value": v}
}

// numberArgs checks that JS passed at least n numbers.
func numberArgs(args []js.Value, n int) bool {
	if len(args) < n {
		return false
	}
	for _, a := range args[:n] {
		if a.Type() != js.TypeNumber {
			return false
		}
	}
	return true
}

// binary exposes a two-number function to JS: calcWasm.name(a, b).
func binary(fn binaryOp) js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) any {
		if !numberArgs(args, 2) {
			return map[string]any{"ok": false, "error": "expected 2 numbers"}
		}
		return result(fn(args[0].Float(), args[1].Float()))
	})
}

// unary exposes a one-number function to JS: calcWasm.name(a).
func unary(fn unaryOp) js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) any {
		if !numberArgs(args, 1) {
			return map[string]any{"ok": false, "error": "expected 1 number"}
		}
		return result(fn(args[0].Float()))
	})
}

func main() {
	// The keys are the names JavaScript uses (data-op="..." in index.html).
	// TODO 3: register your own function here, e.g. "modulo": binary(Modulo),
	calc := map[string]any{
		"add":      binary(Add),
		"subtract": binary(Subtract),
		"multiply": binary(Multiply),
		"divide":   binary(Divide),
		"power":    binary(Power),
		"sqrt":     unary(Sqrt),
		"percent":  unary(Percent),
	}

	js.Global().Set("calcWasm", js.ValueOf(calc))
	js.Global().Call("dispatchEvent", js.Global().Get("Event").New("wasm-ready"))

	// Keep the Go program alive so JS can keep calling the functions.
	select {}
}
