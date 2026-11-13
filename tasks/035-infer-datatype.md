# 035: Infer `datatype` declarations

Depends on: 034.

## Task

Give user datatypes and their constructors types.

## Context

- Read: `infer/dec.go`, `infer/scheme.go`
- Change: `infer/dec.go`, create `infer/datatype_test.go`
- Replace the do-nothing `DatatypeDec` case from task 033. For `datatype ['a] name = C1 [of t1] | ...`:
  1. Register the type constructor before looking at the constructors, so the type can refer to itself: `in.tycons[name] = 0` or `1` (1 when there is a type variable).
  2. Make one `vars` map for the whole declaration. If there is a type variable, put a fresh variable in it under that name.
  3. The result type is `&types.TCon{Name: name}` with the type variable as its one argument when there is one.
  4. For each constructor: without an argument its type is the result type; with an argument convert it with `FromAst(arg, vars)` and its type is `Arrow(argType, resultType)`.
  5. A type variable in a constructor argument that is not the declared one is an error: `unbound type variable`. Check by comparing the size of `vars` before and after.
  6. Mark every variable in each constructor type as generic (`Level = types.Generic`) and define the constructor name in `tenv`.
- Return no bindings.

## Done when

`go test ./infer/ -timeout 30s` passes, with a new `TestDatatype` using the `typesOf` helper:

- `datatype shape = Circle of int | Rect of int * int; fun area (Circle r) = 3 * r * r | area (Rect (w, h)) = w * h;` ends with `shape -> int`
- then `Circle 2;` is `shape` and `[Circle 1, Rect (2, 3)];` is `shape list`
- `datatype 'a tree = Leaf | Node of 'a tree * 'a * 'a tree; fun insert (x, Leaf) = Node (Leaf, x, Leaf) | insert (x, Node (l, v, r)) = if x < v then Node (insert (x, l), v, r) else Node (l, v, insert (x, r));` ends with `int * int tree -> int tree`
- with the same datatype: `fun toList Leaf = [] | toList (Node (l, v, r)) = toList l @ [v] @ toList r;` is `'a tree -> 'a list`
- `Leaf;` is `'a tree`
- errors: `Circle "big";`; `Nope 3;`; `datatype t = A of 'b;`; `datatype t = A of foo;`
