# 037: Parse `exception`, `raise`, and `handle`

Depends on: 036.

## Task

Add exception syntax to the AST and parser.

## Context

- Read: `parser/expr.go`, `parser/dec.go`
- Change: `ast/ast.go`, `ast/show.go`, `parser/expr.go`, `parser/dec.go`, `parser/parser_test.go`
- New nodes:

```go
type Raise struct { Node; Expr Expr }
type Handle struct { Node; Expr Expr; Arms []Arm }
type ExnDec struct { Node; Name string; Arg Ty } // Arg is nil when there is no "of"
```

- `ast.Show`: `(raise e)`, `(handle e (p body) (p body))`, `(exception Name)` or `(exception Name ty)`.
- `ast.ExnLine(d *ExnDec) string`: `exception Oops` or `exception Bad of string` (use `ast.TypeString`).
- `raise e`: handled at the top of `parseExpr` like `if`; the operand is `parseExpr`, so it extends as far right as possible.
- `e handle pat => body | pat => body`: `handle` binds looser than every infix operator. In `parseExpr`, after `parseBinary(1)` returns `left`, if the next token is `handle`, parse the arms (same arm loop as `fn`) and return `ast.Handle{Expr: left, Arms: arms}`.
- `exception Name [of ty]`: in `parseDec`.
- Remove `raise`, `handle`, and `exception` from the "not supported" lists.

## Done when

`go test ./parser/ ./ast/ -timeout 30s` passes, with a new `TestExceptions`:

- `raise Bad "x"` gives `(raise (app Bad "x"))`
- `f x handle Div => 0` gives `(handle (app f x) (Div 0))`
- `1 + 2 handle Oops => 3 | Bad _ => 4` gives `(handle (app + (tuple 1 2)) (Oops 3) ((pcon Bad _) 4))`
- `(raise Oops) handle Bad m => m` gives `(handle (raise Oops) ((pcon Bad m) m))`
- `Parse("exception Bad of string;")` shows `(exception Bad string)`; `ExnLine` gives `exception Bad of string`
- `Parse("fun f (_, 0) = raise Div | f (a, b) = a div b;")` gives 2 clauses
