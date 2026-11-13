# 024: Parse type expressions

Depends on: 009.

## Task

Add type syntax to the AST and parser.

## Context

- Read: `parser/parser.go`
- Change: `ast/ast.go`, `ast/show.go`, create `parser/ty.go`, `parser/ty_test.go`
- New AST nodes:

```go
type Ty interface{ Pos() token.Pos }

type TyVar struct { Node; Name string }            // 'a
type TyCon struct { Node; Name string; Args []Ty } // int, 'a list, int tree
type TyTuple struct { Node; Elems []Ty }           // int * bool
type TyArrow struct { Node; From, To Ty }          // int -> bool
```

- `ast.Show` additions: `TyVar` shows its name; `TyCon` shows `int` or `(list int)`; `TyTuple` shows `(* a b)`; `TyArrow` shows `(-> a b)`.
- Grammar, loosest first:
  - arrow: `tuple [ "->" arrow ]` (right associative)
  - tuple: `postfix { "*" postfix }` (two or more parts make a `TyTuple`)
  - postfix: `atom { IDENT }` (each following identifier wraps what came before: `int list list`)
  - atom: TYVAR, IDENT (a type name with no arguments), or `(` arrow `)`
- Add method `parseType()` on the parser and `func ParseType(src string) (ast.Ty, error)` that tokenizes, parses one type, and requires `EOF`.

## Done when

`go test ./parser/ -run TestType -timeout 30s` passes:

- `int` gives `int`
- `int list list` gives `(list (list int))`
- `int * int -> int` gives `(-> (* int int) int)`
- `('a -> 'b) -> 'a list -> 'b list` gives `(-> (-> 'a 'b) (-> (list 'a) (list 'b)))`
- `'a tree * 'a * 'a tree` gives `(* (tree 'a) 'a (tree 'a))`
- `(int * bool) list` gives `(list (* int bool))`
- `int ->` returns an error
