# 033: `datatype` declarations

Depends on: 032.

## Task

Parse, evaluate, and report `datatype` declarations, and add the built-in `option` constructors.

## Context

- Read: `parser/dec.go`, `parser/ty.go`, `eval/dec.go`, `interp/interp.go`
- Change: `ast/ast.go`, `ast/show.go`, create `ast/typestring.go`, `parser/dec.go`, `eval/dec.go`, `builtins/builtins.go`, `interp/interp.go`, plus tests in each package
- AST:

```go
type DatatypeDec struct { Node; TyVar string; Name string; Cons []ConBind } // TyVar is "" or "'a"
type ConBind struct { Name string; Arg Ty }                                 // Arg is nil when there is no "of"
```

- Syntax: `datatype [tyvar] name = Con [of ty] | Con [of ty] ...`. At most one type variable.
- `ast.TypeString(t Ty) string` prints type syntax as written, with the same parenthesis rules as `types.String` (function loosest, then tuple, then postfix): `int * int`, `'a tree * 'a * 'a tree`, `(int * bool) list`, `(int -> int) -> int`.
- `ast.DatatypeLine(d *DatatypeDec) string`: the declaration on one line: `datatype shape = Circle of int | Rect of int * int`, `datatype 'a tree = Leaf | Node of 'a tree * 'a * 'a tree`.
- `eval.ExecDec` for `DatatypeDec`: for each constructor take `id := value.NextID()`; define the name as `value.Con{Name, ID: id}` when there is no argument, or `value.ConFn{Name, ID: id}` when there is. Return no bindings.
- `builtins`: add `NONE` (`value.Con`, ID 7, type `'a option`) and `SOME` (`value.ConFn`, ID 8, type `'a -> 'a option`).
- `interp.Eval`: for a `DatatypeDec`, add `Result{Kind: "datatype", Text: ast.DatatypeLine(d)}`. `Result.Line` returns `Text` for this kind.
- Type checking of datatypes comes in task 035. Until then make `infer.Dec` accept a `DatatypeDec` and do nothing, so the session can run it.

## Done when

`go test ./ast/ ./parser/ ./eval/ ./interp/ -timeout 30s` passes, with new cases:

- parser: `datatype shape = Circle of int | Rect of int * int;` parses to a `DatatypeDec` with 2 constructors; `datatype 'a tree = Leaf | Node of 'a tree * 'a * 'a tree;` has `TyVar` `'a`
- ast: `DatatypeLine` gives back exactly the two lines above
- eval: after `datatype shape = Circle of int | Rect of int * int;`, `Rect (2, 3);` prints `Rect (2,3)` and `[Circle 1, Rect (2, 3)];` prints `[Circle 1,Rect (2,3)]`
- eval: `datatype t = A | B; case B of A => 1 | B => 2;` is `2`
- interp: `Eval("datatype t = A | B;")` gives one result whose `Line(true)` is `datatype t = A | B`
