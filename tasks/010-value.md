# 010: Runtime values

Depends on: 005.

## Task

Create the `value` package: runtime value types, the raised-exception error, and value printing.

## Context

- Change: create `value/value.go`, `value/print.go`, `value/print_test.go`
- Types, exactly:

```go
package value

type Value interface{}

type Int int64
type Str string
type Bool bool
type Unit struct{}
type Tuple []Value

// List is an immutable cons cell. A nil *List is the empty list.
type List struct {
	Head Value
	Tail *List
}

// Con is a constructor value or exception value: Leaf, SOME 1, Div, Fail "x".
type Con struct {
	Name string
	ID   int
	Arg  Value // nil when the constructor takes no argument
}

// ConFn is a constructor that still needs its argument: SOME, Fail.
type ConFn struct {
	Name string
	ID   int
}

// Builtin is a Go function taking one argument. Curried builtins return another Builtin.
type Builtin struct {
	Name string
	Fn   func(arg Value) (Value, error)
}

// Raise is the error for a raised exception.
type Raise struct{ Exn Con }

func (r *Raise) Error() string { return "uncaught exception " + String(r.Exn) }
```

- Built-in exceptions, as package variables with fixed IDs: `ExnDiv` (ID 1), `ExnMatch` (2), `ExnBind` (3), `ExnEmpty` (4), `ExnSubscript` (5), all `Con` with nil `Arg`; and `FailID = 6`.
- `func NextID() int`: returns 100, 101, 102, ... on successive calls.
- Helpers: `func FromSlice(vs []Value) *List` and `func (l *List) Slice() []Value`.
- `func String(v Value) string`:
  - `Int`: decimal; negatives use `~`: `~3`
  - `Str`: in double quotes, with `\n`, `\t`, `\\`, `\"` escaped again
  - `Bool`: `true` or `false`. `Unit`: `()`
  - `Tuple`: `(1,"two",true)`. `*List`: `[1,2,3]` and `[]`. No space after commas.
  - `Con` without `Arg`: the name. With `Arg`: name, space, argument: `Circle 2`, `Rect (2,3)`, `Fail "x"`. When the argument is itself a `Con` with an `Arg`, wrap it in parentheses: `SOME (SOME 1)`.
  - `ConFn`, `Builtin`, and any other type (closures come later): `fn`

## Done when

`go test ./value/ -timeout 30s` passes, with a table over `String` covering each line above, including `Int(-3)` giving `~3`, `Str("a\nb")` giving `"a\nb"` (backslash n, 6 characters), an empty `(*List)(nil)` giving `[]`, and `SOME (SOME ~1)`.
