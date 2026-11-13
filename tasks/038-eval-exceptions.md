# 038: Evaluate exceptions

Depends on: 037.

## Task

Evaluate `exception` declarations, `raise`, and `handle`.

## Context

- Read: `eval/eval.go`, `eval/dec.go`, `value/value.go`
- Change: `eval/eval.go`, `eval/dec.go`, `interp/interp.go`, create `eval/exception_test.go`
- `ExecDec` for `ExnDec`: `id := value.NextID()`; define the name as `value.Con{Name, ID: id}` (no argument) or `value.ConFn{Name, ID: id}`. Every execution of the declaration makes a new ID, so declaring `exception E` twice makes two different exceptions. Return no bindings.
- `Raise` in `Eval`: evaluate the operand. It must be a `value.Con`; return `&value.Raise{Exn: con}`.
- `Handle` in `Eval`:
  1. Evaluate `e.Expr` with a normal recursive `Eval` call. This is not a tail position.
  2. No error: return the value.
  3. The error is a `*value.Raise`: try the arms against `raise.Exn` with `Match`. On a match make a frame, define the bindings, set `e = arm.Body`, `en = frame`, and `continue`. No arm matches: return the same error.
  4. Any other error: return it unchanged.
- `Match` already compares constructor IDs (task 032), which is what makes `handle` match on identity.
- `interp.Eval`: for an `ExnDec`, add `Result{Kind: "exception", Text: ast.ExnLine(d)}`. `Result.Line` returns `Text` for this kind.
- Until task 039, make `infer.Dec` accept an `ExnDec` and do nothing, and make `infer.Expr` return a fresh variable for `Raise` and for `Handle`, so the eval tests can run through the program helper without the type checker in the way.

## Done when

`go test ./eval/ ./interp/ -timeout 30s` passes, with a new `TestExceptions`:

- `exception Oops; (raise Oops) handle Oops => 1;` is `1`
- `exception Bad of string; (raise Bad "x") handle Bad msg => msg ^ "!";` is `"x!"`
- `hd [] handle Empty => 0;` is `0`
- `1 div 0 handle Div => ~1;` is `~1`
- `exception Oops; exception Bad of string; (raise Oops) handle Bad _ => 1;` returns a `*value.Raise` whose `Error()` is `uncaught exception Oops`
- `raise Fail "nope";` returns a `*value.Raise` whose `Error()` is `uncaught exception Fail "nope"`
- `exception E; val old = E; exception E; (raise old) handle E => 1;` still raises (the two `E`s differ)
- `exception Bad of string; Bad "y";` prints `Bad "y"`
- interp: `Eval("exception Bad of string;")` gives one result whose `Line(true)` is `exception Bad of string`
