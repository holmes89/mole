# 029: Type check in the session; pass the milestone 2 golden test

Depends on: 028.

## Task

Run type inference before evaluation in `interp`, report types, and make `m2_types` pass.

## Context

- Read: `interp/interp.go`, `infer/dec.go`, `tests/golden/m2_types.mole`, `tests/golden/m2_types.out`
- Change: `interp/interp.go`, `interp/interp_test.go`
- `Interpreter` gains `inf *infer.Infer` and `tenv *env.Env`.
- `New`: also `Declare` every builtin entry's type in `tenv`, plus `use` as `string -> unit`.
- `Eval`, for each declaration:
  1. `tframe := env.New(in.tenv)`; `tbinds, err := in.inf.Dec(dec, tframe)`. On an error return it; nothing is evaluated and nothing is bound.
  2. Evaluate as before in a value frame. On an error return it; neither frame is kept.
  3. On success keep both frames. Fill `Result.Type` with `types.String` of the binding's type. Value bindings and type bindings come in the same order.
- The `use` change from task 020 (re-parenting the frame when a nested `Eval` ran) applies to the type frame too.
- Set `const Milestone = 2`.
- Then run `go test ./tests/ -timeout 60s` and fix failures one at a time, adding a unit test for each fix in the package you changed. `m1_core` must keep passing.

## Constraints

- Never edit anything under `tests/golden/` or `tests/golden_test.go`.

## Done when

`make vet && make test` passes, including `TestGolden/m1_core.mole` and `TestGolden/m2_types.mole`, and after `make build`, `printf 'fun id x = x;\n' | ./bin/mole` prints `val id = fn : 'a -> 'a`.
