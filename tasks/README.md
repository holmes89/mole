# mole tasks

39 task files that build mole in order. Each is sized for one fresh agent session.

| Tasks | Milestone | Ends with |
|---|---|---|
| 001 to 022 | M1: untyped core, REPL, `use` | `m1_core` golden passes |
| 023 to 029 | M2: Hindley-Milner types | `m2_types` golden passes |
| 030 to 036 | M3: patterns, `case`, `datatype` | `m3_patterns` golden passes |
| 037 to 039 | M4: exceptions | all four goldens pass |

## How to run one

Fresh session, builder agent: `Do tasks/003-lexer-basic.md`

Commit after each task that reports DONE. If a task reports BLOCKED, improve the task file (or hand that one task to a stronger model) and run it again in a new session.

## Where a small model is most likely to struggle

- 014 (closures and declarations) and 020 (`use`): the most moving parts in M1.
- 025 to 028 (unification, levels, generalization): the hardest reasoning in the project.
- 019, 029, 036, 039 (golden tasks): open-ended debugging. Each failure is small, but there may be several.

## Design decisions baked into the tasks

- Infix operators are builtins applied to a tuple, so the evaluator and the type checker need no operator cases.
- `fun` is turned into nested `fn` by `ast.FunToFn`; both the evaluator and the type checker use it.
- Each top-level declaration runs in a new environment frame that is kept only on success, so a failed declaration binds nothing.
- Types are `TVar` and `TCon` only; tuples and functions are `TCon` named `*` and `->`. Polymorphism uses levels, with `types.Generic` marking generalized variables.
- Constructors and exceptions are `value.Con` with an integer ID; patterns match on ID.
- The golden runner (`tests/golden/runner.go.txt`) is copied in by task 019, not written by the agent, so no model grades itself.
