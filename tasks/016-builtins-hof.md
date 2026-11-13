# 016: Builtins, part 2 (higher-order functions)

Depends on: 014.

## Task

Add the builtins that call ML functions.

## Context

- Read: `builtins/builtins.go`
- Change: create `builtins/hof.go`, edit `builtins/builtins.go` (to include the new entries), create `eval/hof_test.go`
- These use the `apply` argument of `Table(apply)` to call a closure or builtin.
- A curried builtin returns another `value.Builtin` for each argument. For example `map` takes `f` and returns a `Builtin` that takes the list.

| Name | Type | Behavior |
|---|---|---|
| `map` | `('a -> 'b) -> 'a list -> 'b list` | left to right |
| `foldl` | `('a * 'b -> 'b) -> 'b -> 'a list -> 'b` | calls `f (x, acc)` from the left |
| `foldr` | `('a * 'b -> 'b) -> 'b -> 'a list -> 'b` | calls `f (x, acc)` from the right |
| `List.filter` | `('a -> bool) -> 'a list -> 'a list` | |
| `List.exists` | `('a -> bool) -> 'a list -> bool` | stops at the first true |
| `List.all` | `('a -> bool) -> 'a list -> bool` | stops at the first false |
| `o` | `('b -> 'c) * ('a -> 'b) -> 'a -> 'c` | takes a tuple `(f, g)`; returns a `Builtin` computing `f (g x)` |
| `print` | `string -> unit` | see below |

- `print` needs an output stream. Change the signature to `func Table(apply ApplyFunc, out io.Writer) []Entry` and update the callers (the `eval` test helper passes `os.Stdout` or a buffer). `print` writes the string to `out` and returns `value.Unit{}`.
- Propagate any error from `apply` immediately.

## Done when

`go test ./eval/ ./builtins/ -timeout 30s` passes, with `eval/hof_test.go` cases run through the program helper:

- `map (fn n => n * n) [1, 2, 3];` is `[1,4,9]`
- `foldl (fn (x, acc) => x + acc) 0 [1, 2, 3];` is `6`
- `foldr (fn (x, acc) => x :: acc) [] [1, 2, 3];` is `[1,2,3]`
- `foldl (fn (x, acc) => x :: acc) [] [1, 2, 3];` is `[3,2,1]`
- `List.filter (fn n => n > 1) [1, 2, 3];` is `[2,3]`
- `List.exists (fn n => n > 2) [1, 2, 3];` is `true`
- `fun double x = x * 2; val quad = double o double; quad 3;` is `12`
- `print "hi\n";` writes `hi` and a newline to the buffer and gives `()`
- `map (fn n => 1 div n) [1, 0];` returns a `*value.Raise`
