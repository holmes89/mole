# 001: Scaffold the repository

## Task

Create the Go module, Makefile, and a placeholder CLI.

## Context

- Change: create `go.mod`, `Makefile`, `main.go`, `.gitignore`
- `go.mod`: module `github.com/holmes89/mole`, with the `go` directive matching the local `go version` (major.minor).
- `Makefile` targets (all `.PHONY`):
  - `build`: `go build -o bin/mole .`
  - `test`: `go test -timeout 120s ./...`
  - `vet`: `go vet ./...`
  - `clean`: `rm -rf bin/`
- `main.go`: package `main`; `main()` prints `mole: not implemented` and a newline.
- `.gitignore`: one line, `bin/`.
- `tests/golden/` and `docs/` already exist. Leave them alone.

## Constraints

- No dependencies yet.

## Done when

`make build && make vet && ./bin/mole` prints `mole: not implemented`.
