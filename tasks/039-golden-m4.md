# 039: Infer exceptions; pass the milestone 4 golden test

Depends on: 038.

## Task

Type check exceptions and make `m4_exceptions` pass.

## Context

- Read: `infer/dec.go`, `infer/expr.go`, `tests/golden/m4_exceptions.mole`, `tests/golden/m4_exceptions.out`
- Change: `infer/dec.go`, `infer/expr.go`, create `infer/exception_test.go`, `interp/interp.go`
- Replace the placeholders from task 038:
  - `ExnDec` without an argument: define the name in `tenv` as `exn`. With an argument: convert it with `FromAst` (a type variable in it is the error `unbound type variable`) and define the name as `Arrow(argType, exn)`.
  - `Raise`: the operand's type must unify with `exn`. The result is a fresh variable, so `raise` fits anywhere.
  - `Handle`: `t := type of e.Expr`. For each arm: frame, `Pat`, unify the pattern type with `exn`, define the bindings, infer the body, unify with `t`. Result `t`.
- Set `const Milestone = 4` in `interp/interp.go`.
- Then run `go test ./tests/ -timeout 60s` and fix failures one at a time, adding a unit test for each fix.

## Constraints

- Never edit anything under `tests/golden/` or `tests/golden_test.go`.
- The three earlier golden files must keep passing.

## Done when

`make vet && make test` passes, including all four `TestGolden` files, and `infer/exception_test.go` covers:

- `exception Bad of string; Bad "y";` is `exn`
- `fun check n = if n < 0 then raise Fail "negative" else n;` is `int -> int`
- `hd [] handle Empty => 0;` is `int`
- `exception Bad of string; (raise Bad "x") handle Bad msg => msg ^ "!";` is `string`
- errors: `raise 5;`; `exception Oops; 1 handle Oops => "one";`; `1 handle 5 => 1;`
