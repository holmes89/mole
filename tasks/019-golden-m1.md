# 019: Pass the milestone 1 golden test

Depends on: 018.

## Task

Install the golden test runner and make `m1_core` pass.

## Context

- Read: `tests/golden/runner.go.txt`, `tests/golden/m1_core.mole`, `tests/golden/m1_core.out`
- Change: create `tests/golden_test.go`; fix whatever package the failures point to
- Step 1: `cp tests/golden/runner.go.txt tests/golden_test.go`. Do not change its contents.
- Step 2: `go test ./tests/ -timeout 60s`. Each failure names a line of `m1_core.out` with the wanted and actual text.
- Step 3: fix the implementation, one failure at a time, and add a unit test for each fix in the package you changed.
- The golden file is compared without types (its first line is `(* golden: untyped *)`).
- An expected line `Error:` matches any output line starting with `Error:`.

## Constraints

- Never edit anything under `tests/golden/` and never edit `tests/golden_test.go` after copying it.
- Never special-case a golden input in the implementation.

## Done when

`make vet && make test` passes, including `TestGolden/m1_core.mole`.
