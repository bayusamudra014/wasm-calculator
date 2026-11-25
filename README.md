# WASM Calculator (template)

A calculator web app whose math is written in Go and compiled to WebAssembly.
The UI is finished. Your job is to write the math functions in Go.

## Files

| File | What it is | Do you edit it? |
|---|---|---|
| `calc.go` | The math functions (`Add`, `Subtract`, …). Plain Go. | **Yes, this is your main file** |
| `calc_test.go` | Unit tests that check `calc.go`. | Optional (add tests for your own functions) |
| `main.go` | Bridge: exposes the Go functions to JS as `window.calcWasm`. | Only to register new functions |
| `web/index.html` | Page and keypad. | Only to add new buttons |
| `web/app.js` | Calculator UI logic (typing, display, calling WASM). | No |
| `web/style.css` | Styles. | No |
| `Makefile` | `build`, `serve`, `test`, `clean`. | No |

## Run it

Requires Go 1.21+ and Python 3 (for the static file server).

```bash
make serve
```

Then open http://localhost:8083. You need a web server because browsers won't load a `.wasm` file from `file://`.

Under the keypad, the **WASM functions** list shows the state of each function:
✓ implemented, ○ still a TODO in `calc.go`, ✗ has a button but isn't registered in `main.go`.

Only `+` works at first. Every other key shows "not implemented yet" until you write it.

## How it works

```
 ┌─────────── browser ───────────┐
 │ index.html ─ buttons          │
 │ app.js     ─ calcWasm.add(2,3)├──► main.go  binary(Add) ──► calc.go  Add(2,3)
 │            ◄─ {ok:true,value:5}◄──           result(...) ◄──          5, nil
 └───────────────────────────────┘
```

1. `make build` compiles the Go code with `GOOS=js GOARCH=wasm` into `web/main.wasm`, and copies Go's JS glue file `wasm_exec.js` into `web/`.
2. `app.js` loads `main.wasm` and runs it.
3. Go's `main()` creates `window.calcWasm` (one JS function per Go function) and fires a `wasm-ready` event.
4. When you press `=`, `app.js` calls something like `calcWasm.divide(9, 3)`.
5. The bridge in `main.go` converts the JS numbers to `float64`, calls `Divide(9, 3)`, and turns the `(value, error)` result into `{ ok: true, value: 3 }` or `{ ok: false, error: "..." }`.
6. `app.js` shows the value, or the error message, on the display.

## Exercises

Work in `calc.go`. After each change, check it two ways:

```bash
make test
```

```bash
make serve
```

`make test` runs the tests on your machine without a browser. Tests for functions you haven't written yet show as **SKIP**. A function is done when its test shows **PASS**.

### TODO 1: basic operations
Write `Subtract`, `Multiply` and `Divide`.
- `Divide` must return an error when dividing by zero:
  ```go
  if b == 0 {
      return 0, errors.New("cannot divide by zero")
  }
  ```
  Try `5 ÷ 0 =` in the browser. Your error message appears on the display.

### TODO 2: more operations
Write `Power`, `Sqrt` and `Percent`.
- Add `"math"` to the imports and use `math.Pow` and `math.Sqrt`.
- `Sqrt` of a negative number should return an error.

### TODO 3: add your own function
For example, `Modulo` (remainder of `a / b`):

1. **Write it in `calc.go`:**
   ```go
   func Modulo(a, b float64) (float64, error) {
       if b == 0 {
           return 0, errors.New("cannot divide by zero")
       }
       return math.Mod(a, b), nil
   }
   ```
2. **Register it in `main.go`** inside the `calc` map. Use `binary(...)` for two numbers or `unary(...)` for one:
   ```go
   "modulo": binary(Modulo),
   ```
3. **Add a button in `web/index.html`** inside `<div class="keys">`. `data-op` must match the name from step 2:
   ```html
   <button data-op="modulo" data-arity="2" data-symbol="mod" class="fn">mod</button>
   ```
   Use `data-arity="1"` for one-number functions (they run right away, like `√x`).
4. **(Optional) Add a test** in `calc_test.go`, copying the pattern of `TestDivide`.
5. Run `make serve` again and try it.

More ideas: `Square` (x²), `Reciprocal` (1/x), `Factorial` (n!, with an error for negative or non-integer input), `Abs`, `Log10`.

## Troubleshooting

- **"Could not load main.wasm"**: run `make build` (or `make serve`), and open the page through http://localhost:8083, not by double-clicking the file.
- **Changes don't show up**: rebuild with `make build`, then hard-refresh the browser (Cmd+Shift+R) so it doesn't use the cached `main.wasm`.
- **A button shows `"xyz" is not registered in main.go`**: the `data-op` name in `index.html` doesn't match a key in the `calc` map in `main.go`.
- **Build error `"math" imported and not used`**: Go doesn't allow unused imports. Only add `"math"` once you use it.
