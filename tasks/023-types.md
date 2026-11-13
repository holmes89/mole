# 023: Type representation and printing

## Task

Create the `types` package: how types are represented and printed.

## Context

- Change: create `types/types.go`, `types/print.go`, `types/print_test.go`
- Representation, exactly:

```go
package types

type Type interface{}

// TVar is a type variable. When Link is not nil the variable has been
// unified with Link, and stands for it.
type TVar struct {
	ID    int
	Level int
	Link  Type
}

// TCon is every other type. Name is "int", "bool", "string", "unit", "exn",
// "list", "option", a datatype name, "*" for tuples, or "->" for functions.
type TCon struct {
	Name string
	Args []Type
}

// Generic is the Level of a generalized (polymorphic) variable.
const Generic = 1 << 30
```

- Helpers: `Int()`, `Bool()`, `String_()`, `Unit()`, `Exn()` returning `*TCon` with no args; `List(t)`, `Arrow(a, b)`, `Tuple(ts ...Type)`.
- `func Prune(t Type) Type`: while `t` is a `*TVar` with a non-nil `Link`, follow the link. Return the result.
- `func String(t Type) string` prints a type. Prune at every step.
  - Base types: the name. `int`.
  - Type variables: renamed in order of first appearance, left to right, within this one call: `'a`, `'b`, `'c`, ...
  - Postfix constructors: argument, space, name: `int list`, `int list list`, `'a option`.
  - Tuples: parts joined with ` * `.
  - Functions: `a -> b`, right associative.
  - Parentheses, using three precedence levels (function 1, tuple 2, postfix 3):
    - left side of `->` at level 1 gets parentheses: `('a -> 'b) -> 'a`; the right side never does.
    - a tuple part that is a tuple or function gets parentheses: `int * (int * int)`, `(int -> int) * int`.
    - a postfix argument that is a tuple or function gets parentheses: `(int * bool) list`.

## Done when

`go test ./types/ -timeout 30s` passes, with a table over `String`:

- `Arrow(Tuple(Int(), Int()), Int())` gives `int * int -> int`
- `Arrow(Arrow(a, b), Arrow(List(a), List(b)))` gives `('a -> 'b) -> 'a list -> 'b list`
- `List(List(Int()))` gives `int list list`
- `List(Tuple(Int(), Bool()))` gives `(int * bool) list`
- `Arrow(Tuple(a, b), Tuple(b, a))` gives `'a * 'b -> 'b * 'a`
- two unrelated variables with IDs 7 and 3, in that order, give `'a` then `'b`
- a variable linked to `Int()` prints `int`
