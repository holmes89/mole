# 018: `interp` session API

Depends on: 016, 017.

## Task

Create the public interpreter API: a session that evaluates source and reports results.

## Context

- Read: `eval/dec.go`, `builtins/builtins.go`, `interp/split.go`
- Change: create `interp/interp.go`, `interp/interp_test.go`

```go
package interp

// Milestone is raised when a milestone's golden file passes.
const Milestone = 1

type Result struct {
	Kind  string // "val" (later: "datatype", "exception")
	Name  string
	Value string // printed with value.String
	Type  string // "" until task 029
	Text  string // unused for "val"
}

func (r Result) Line(withTypes bool) string
type Interpreter struct { /* out io.Writer; env *env.Env */ }
func New(out io.Writer) *Interpreter
func (in *Interpreter) Eval(src string) ([]Result, error)
func FormatError(err error) string
```

- `Line`: for kind `val`, `val <Name> = <Value>`, plus ` : <Type>` when `withTypes` is true and `Type` is not empty.
- `New`: `env.New(nil)` with every entry of `builtins.Table(eval.Apply, out)` defined.
- `Eval`: `parser.Parse(src)`. For each declaration, in order:
  - `frame := env.New(in.env)`; `bindings, err := eval.ExecDec(dec, frame)`.
  - On an error: return the results so far and the error. `in.env` is not changed, so a failed declaration binds nothing.
  - On success: `in.env = frame`; add one `Result` per binding.
- A parse error returns no results and the error.
- Wrap the body of `Eval` in a deferred `recover` that turns a panic into the error `internal: <message>`.
- `FormatError`: a `*value.Raise` gives `uncaught exception <value.String(exn)>`. Anything else gives `Error: <err.Error()>`.

## Done when

`go test ./interp/ -timeout 30s` passes, with a new `TestEval`:

- `1 + 2;` gives one result whose `Line(false)` is `val it = 3`
- `val (a, b) = (1, "x");` gives two results: `val a = 1`, `val b = "x"`
- `6 * 7; it + 1;` gives `val it = 42` then `val it = 43`
- `val x = 10; undefinedName; ` returns one result and an error; a later `Eval("x;")` still gives `val it = 10`
- `FormatError` of the `1 div 0;` error is `uncaught exception Div`
- `FormatError` of the `undefinedName;` error starts with `Error:`
- `print "hello\n";` writes `hello` and a newline to the writer passed to `New`
