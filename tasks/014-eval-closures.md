# 014: Evaluator, part 2 (closures, patterns, declarations)

Depends on: 013.

## Task

Add closures, pattern matching, `let`, and declarations to the evaluator.

## Context

- Read: `eval/eval.go`, `ast/ast.go`, `ast/desugar.go`
- Change: `eval/eval.go`, create `eval/match.go`, `eval/dec.go`, `eval/closure_test.go`
- Closure type, in package `eval`:

```go
type Closure struct {
	Arms []ast.Arm
	Env  *env.Env
}
```

- `match.go`:

```go
type Binding struct { Name string; Value value.Value }

// Match reports whether v matches p, and the variables bound, left to right.
func Match(p ast.Pat, v value.Value, en *env.Env) ([]Binding, bool)
```

  - `PWild`: matches. `PVar`: matches and binds. `PTuple` with zero elements: matches `value.Unit`. `PTuple`: matches a `value.Tuple` of the same length when every element matches.
  - `en` is unused for now; later tasks need it.
- `Eval` additions:
  - `Fn`: return `&Closure{Arms: e.Arms, Env: en}`.
  - `App` with a `*Closure`: try the arms in order with `Match`. On the first match, make `frame := env.New(closure.Env)`, define the bindings in it, then set `e = arm.Body`, `en = frame`, and `continue`. No arm matches: return `&value.Raise{Exn: value.ExnMatch}`.
  - `Let`: `frame := env.New(en)`; run each declaration with `ExecDec(dec, frame)`; evaluate all body expressions but the last; then `e = last`, `en = frame`, `continue`.
- `Apply` for a `*Closure`: match, build the frame, and call `Eval(arm.Body, frame)`.
- `dec.go`:

```go
// ExecDec runs a declaration, defines its names in en, and returns them in order.
func ExecDec(d ast.Dec, en *env.Env) ([]Binding, error)
```

  - `ValDec`: evaluate the expression, then `Match`. No match: `&value.Raise{Exn: value.ExnBind}`. Define each binding.
  - `FunDec`: build `&Closure{Arms: ast.FunToFn(d).(*ast.Fn).Arms, Env: en}` and define the function's name in `en`. Because the closure's `Env` is `en` itself, the function can call itself. (Adjust if your `Fn` nodes are values rather than pointers.)

## Done when

`go test ./eval/ -timeout 30s` passes. Add a helper that runs a whole program with `parser.Parse`, calling `ExecDec` for each declaration in one env, and returns the printed value of the last binding. Cases:

- `val add = fn a => fn b => a + b; add 2 3;` is `5`
- `fun fact n = if n = 0 then 1 else n * fact (n - 1); fact 10;` is `3628800`
- `let val y = 5 val z = 6 in y * z end;` is `30`
- `val (a, b) = (1, "x"); b;` is `"x"`
- `(fn (x, _) => x) (1, 2);` is `1`
- `fun k x y = x; k 1 2;` is `1`
- `val f = fn x => x;` prints `fn`
- `val (a, b) = (1, 2, 3);` returns a `*value.Raise` with `ExnBind`
