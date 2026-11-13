# 032: Evaluate full patterns and `case`

Depends on: 031.

## Task

Extend `Match` to the new pattern forms and evaluate `case`.

## Context

- Read: `eval/match.go`, `eval/eval.go`
- Change: `eval/match.go`, `eval/eval.go`, create `eval/pattern_test.go`
- `Match` additions:
  - `PInt`, `PStr`: match an equal `value.Int` or `value.Str`.
  - `PList`: matches a list of exactly that length whose elements match.
  - `PCons`: matches a non-empty list; `Head` against the head, `Tail` against the tail (a `*value.List`, possibly nil).
  - `PCon`:
    - name `true` or `false`: matches the equal `value.Bool`.
    - otherwise look the name up in `en`. It must be a `value.Con` or `value.ConFn`; take its `ID`. The value must be a `value.Con` with the same `ID`. If the pattern has an `Arg`, match it against the value's `Arg`.
    - a name missing from `en` does not match.
- At every call site, pass the environment in which the pattern appears: the closure's `Env` for arms, the current env for `val` and `case`.
- `Case` in `Eval`: evaluate the scrutinee; try the arms in order; on a match make a frame on the current env, define the bindings, set `e = arm.Body`, `en = frame`, and `continue`. No match: `&value.Raise{Exn: value.ExnMatch}`.
- For tests, define constructors by hand in the env: `NONE` as `value.Con{Name: "NONE", ID: 50}` and `SOME` as `value.ConFn{Name: "SOME", ID: 51}`. Task 033 makes them real.

## Done when

`go test ./eval/ -timeout 30s` passes, with a new `TestPatterns` using the program helper:

- `fun sum [] = 0 | sum (x :: xs) = x + sum xs; sum [1, 2, 3, 4];` is `10`
- `fun fib 0 = 0 | fib 1 = 1 | fib n = fib (n - 1) + fib (n - 2); fib 20;` is `6765`
- `case SOME 3 of NONE => 0 | SOME n => n + 1;` is `4`
- `(fn [] => "empty" | _ => "nonempty") [1];` is `"nonempty"`
- `val [p, q] = [1, 2]; q;` is `2`
- `val h :: t = [7, 8, 9]; t;` is `[8,9]`
- `case true of false => 0 | true => 1;` is `1`
- `fun zero 0 = "zero"; zero 1;` raises `ExnMatch`
- `val [r] = [1, 2];` raises `ExnBind`
- `SOME (SOME ~1);` prints `SOME (SOME ~1)`
