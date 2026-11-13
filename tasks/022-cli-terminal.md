# 022: CLI, interactive mode

Depends on: 021.

## Task

Implement the interactive REPL with readline.

## Context

- Read: `main.go`
- Change: `main.go`, `go.mod`, `go.sum`
- Add the dependency: `go get github.com/chzyer/readline`.
- `runTerminal(in)`:
  - Print the banner `Mole ML (Ctrl-D to exit)`.
  - Create readline with prompt `- ` and history file `~/.mole_history` (build the path with `os.UserHomeDir`; if that fails, use no history file).
  - Loop: read a line. Append it plus a newline to the buffer. Call `feed` (from task 021) and keep the rest.
  - Prompt is `- ` when the trimmed buffer is empty, and `= ` when it holds an unfinished declaration.
  - `readline.ErrInterrupt` (Ctrl-C): clear the buffer, reset the prompt, continue.
  - `io.EOF` (Ctrl-D): exit with status 0.
- Save each completed input to history.

## Constraints

- Do not change `feed` or `runPiped`.

## Done when

`make vet && make build && make test` pass, and the piped checks from task 021 still give the same output. Then report that milestone 1 is complete. (Interactive behavior needs a human to try: `./bin/mole`, type `fun f x =`, see the `= ` prompt, type `x + 1;`, see `val f = fn`.)
