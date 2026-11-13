# mole

An interactive REPL for a small subset of Standard ML, written in Go. Pipeline: source, tokens, AST, types (from task 023), values.

## Commands

- Build: `make build` (binary at `bin/mole`)
- Test: `make test`
- Vet: `make vet`
- One test: `go test ./lexer/ -run TestName -timeout 30s`

Run `make vet` and the task's own check before reporting done. Always pass `-timeout` to `go test`.

## Layout

Module path: `github.com/holmes89/mole`.

- `token/`: token kinds and positions
- `lexer/`: source to tokens
- `ast/`: Expr, Pat, Dec, Ty nodes, plus `ast.Show` for tests
- `parser/`: tokens to AST (hand-written Pratt and recursive descent)
- `value/`: runtime values and how they print
- `env/`: chained environments (`map[string]any` with a parent pointer)
- `builtins/`: one table of builtin values with their type strings
- `eval/`: evaluator loop and pattern matching
- `types/`, `infer/`: types and Hindley-Milner inference (tasks 023 on)
- `interp/`: public API used by the CLI and the golden tests
- `tests/golden/`: acceptance files. Never edit them.
- `main.go`: CLI

## Conventions

- Standard library only, plus `github.com/chzyer/readline` in `main.go`. Tests use the standard `testing` package with table-driven cases.
- Return `error`. Never panic on user input.
- AST nodes and types are always pointers: build `&ast.IntLit{...}` and switch on `*ast.IntLit`; same for `*types.TVar` and `*types.TCon`.
- Runtime values are plain values (`value.Int(3)`, `value.Con{...}`), except `*value.List`, `*value.Raise`, and `*eval.Closure`.
- Infix operators are ordinary builtins: `a + b` parses to `(app + (tuple a b))`.
- Runtime failures such as division by zero are `*value.Raise` errors carrying an exception value.
- gofmt. Keep files under about 400 lines; split by concern.

## Tasks

Work comes from files in `tasks/`, done in number order, one task file per session. Each task file contains everything needed for that task. `docs/SPEC.md` is the full language reference for humans: do not read it unless a task tells you to, and when a task and the spec differ, the task wins.
