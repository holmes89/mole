# 036: Pass the milestone 3 golden test

Depends on: 035.

## Task

Make `m3_patterns` pass.

## Context

- Read: `tests/golden/m3_patterns.mole`, `tests/golden/m3_patterns.out`
- Change: `interp/interp.go` (set `const Milestone = 3`); then whatever package the failures point to
- Run `go test ./tests/ -timeout 60s`. Fix failures one at a time, adding a unit test for each fix in the package you changed.
- Things this file checks that earlier tasks may have missed:
  - A `datatype` line prints exactly as `ast.DatatypeLine` gives it.
  - `val [p, q] = [1, 2];` prints two lines, in order.
  - `zero 1;` prints `uncaught exception Match`; `val [r] = [1, 2];` prints `uncaught exception Bind`.
  - The last four inputs are type errors and print a line starting with `Error:`.

## Constraints

- Never edit anything under `tests/golden/` or `tests/golden_test.go`.
- `m1_core` and `m2_types` must keep passing.

## Done when

`make vet && make test` passes, including all three of `TestGolden/m1_core.mole`, `m2_types.mole`, and `m3_patterns.mole`.
