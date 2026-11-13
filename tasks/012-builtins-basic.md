# 012: Builtins, part 1 (operators and first-order functions)

Depends on: 010.

## Task

Create the builtin table with the operators and the simple builtins.

## Context

- Read: `value/value.go`
- Change: create `builtins/builtins.go`, `builtins/ops.go`, `builtins/builtins_test.go`

```go
package builtins

type ApplyFunc func(fn, arg value.Value) (value.Value, error)

type Entry struct {
	Name  string
	Type  string      // ML type, used from task 026 on
	Value value.Value
}

// Table returns every builtin. apply calls an ML function value; it is used in task 016.
func Table(apply ApplyFunc) []Entry
```

- An operator takes one `value.Tuple` of two elements. Wrong argument shapes return a plain error such as `+: expected int * int`.

| Name | Type | Notes |
|---|---|---|
| `+` `-` `*` | `int * int -> int` | |
| `div` `mod` | `int * int -> int` | round toward negative infinity: `~7 div 2 = ~4`, `~7 mod 2 = 1`. Divisor 0 returns `&value.Raise{Exn: value.ExnDiv}` |
| `<` `>` `<=` `>=` | `int * int -> bool` | |
| `=` `<>` | `'a * 'a -> bool` | structural equality on Int, Str, Bool, Unit, Tuple, List, Con (same ID and equal Arg). Anything else: error `cannot compare functions` |
| `^` | `string * string -> string` | |
| `::` | `'a * 'a list -> 'a list` | |
| `@` | `'a list * 'a list -> 'a list` | |
| `~` | `int -> int` | |
| `not` | `bool -> bool` | |
| `size` | `string -> int` | |
| `hd` | `'a list -> 'a` | empty list raises `ExnEmpty` |
| `tl` | `'a list -> 'a list` | empty list raises `ExnEmpty` |
| `null` | `'a list -> bool` | |
| `length` | `'a list -> int` | |
| `rev` | `'a list -> 'a list` | |
| `Int.toString` | `int -> string` | negatives print with `~` |
| `Int.abs` | `int -> int` | |
| `Int.min` `Int.max` | `int * int -> int` | |
| `String.concat` | `string list -> string` | |
| `List.nth` | `'a list * int -> 'a` | out of range raises `ExnSubscript` |

- Also add these plain values: `true` and `false` (`value.Bool`, type `bool`), `nil` (`(*value.List)(nil)`, type `'a list`), `Div`, `Match`, `Bind`, `Empty`, `Subscript` (the `value.Exn...` variables, type `exn`), and `Fail` (`value.ConFn{Name: "Fail", ID: value.FailID}`, type `string -> exn`).

## Done when

`go test ./builtins/ -timeout 30s` passes. Write a test helper `call(name, arg)` that finds the entry in `Table(nil)` and calls its `Fn`. Cases: `+ (1,2)` is 3; `div (~7,2)` is ~4; `mod (~7,2)` is 1; `div (1,0)` returns a `*value.Raise` with `ExnDiv`; `= ([1,2],[1,2])` is true; `:: (0,[1])` prints `[0,1]`; `hd []` raises `ExnEmpty`; `Int.toString ~42` is `"~42"`; every entry has a non-empty `Type`.
