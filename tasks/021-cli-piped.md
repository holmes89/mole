# 021: CLI, piped mode

Depends on: 020.

## Task

Make `main.go` run a program read from standard input when input is not a terminal.

## Context

- Read: `interp/interp.go`, `interp/split.go`
- Change: `main.go`
- Structure `main.go` as:
  - `main()`: creates `interp.New(os.Stdout)`; if stdin is a terminal, calls `runTerminal(in)` (a stub for now that prints `interactive mode: see task 022`); otherwise calls `runPiped(in, os.Stdin)`.
  - Terminal check: `fi, _ := os.Stdin.Stat(); isTerminal := fi.Mode()&os.ModeCharDevice != 0`.
  - `func feed(in *interp.Interpreter, buf string) (rest string)`: shared by both modes. `SplitInput(buf)`; on an error print `interp.FormatError(err)` and return `""`. For each chunk: `Eval`, print each result's `Line(withTypes)`, and on an error print `FormatError(err)`. Return the rest.
  - `withTypes := interp.Milestone >= 2`.
  - `runPiped`: read lines with `bufio.Scanner` (raise its buffer limit to 1 MB), append each line plus a newline to a buffer, call `feed`, keep the rest. At end of input, if the trimmed rest is not empty, print `Error: unexpected end of input`.
- In piped mode print no banner and no prompts.

## Done when

`make vet && make build` pass and these hold:

- `printf '1 + 2;\nval x = 10;\nx * 2;\n' | ./bin/mole` prints exactly `val it = 3`, `val x = 10`, `val it = 20`
- `printf 'fun f x =\n  x + 1;\nf 1;\n' | ./bin/mole` prints `val f = fn`, `val it = 2`
- `printf '1 div 0;\n2;\n' | ./bin/mole` prints `uncaught exception Div`, `val it = 2`
- `printf '1 +' | ./bin/mole` prints `Error: unexpected end of input`
- `./bin/mole < tests/golden/m1_core.mole | diff - tests/golden/m1_core.out` shows a difference only on the one `Error:` line
