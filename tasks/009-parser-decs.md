# 009: Parser, part 4 (declarations and programs)

Depends on: 008.

## Task

Parse `val` and `fun` declarations and whole programs, and add `ast.FunToFn`.

## Context

- Change: create `parser/dec.go`, edit `parser/parser.go`, `parser/parser_test.go`, create `ast/desugar.go`, `ast/desugar_test.go`
- `parseDec` (replace the stub from task 008):
  - `val` pat `=` e: `ast.ValDec`.
  - `fun` name apat apat ... `=` e: `ast.FunDec` with one `Clause`. At least one pattern. Each pattern is `parseAtomicPat`.
  - `datatype` and `exception`: return the error `<keyword> not supported`.
- `func Parse(src string) ([]ast.Dec, error)`: tokenize, then repeat until `EOF`:
  - If the next token is `val`, `fun`, `datatype`, or `exception`: `parseDec`, then require `;`.
  - Otherwise: parse an expression, require `;`, and wrap it as `ast.ValDec{Pat: ast.PVar{Name: "it"}, Expr: e}`.
- `let` now works, because `parseDec` is real.
- `ast/desugar.go`: `func FunToFn(d *FunDec) Expr` turns a function declaration into a `Fn` expression:
  - Every clause has the same number of patterns, n.
  - n = 1: one `Fn` whose arms are the clauses (`Pat: clause.Pats[0]`).
  - n > 1, one clause: nested single-arm `Fn`s, one per pattern: `fn p1 => fn p2 => body`.
  - n > 1, several clauses: `fn $1 => fn $2 => (fn (p1, p2) => body1 | ...) ($1, $2)`. The parameters are `PVar`s named `$1`, `$2`, ...; the inner `Fn` has one arm per clause with a `PTuple` pattern; it is applied to a `Tuple` of `Var`s `$1`, `$2`.

## Done when

`go test ./parser/ ./ast/ -timeout 30s` passes, with new cases:

- `Parse("val x = 10;")` shows `(val x 10)`
- `Parse("1 + 2;")` shows `(val it (app + (tuple 1 2)))`
- `Parse("fun add a b = a + b;")` shows `(fun add ((a b) (app + (tuple a b))))`
- `Parse("val (a, b) = (1, 2); a;")` returns 2 declarations
- `ParseExpr("let val y = 5 val z = 6 in y * z end")` shows `(let ((val y 5) (val z 6)) (app * (tuple y z)))`
- `Parse("val x = 1")` (no `;`) returns an error
- `Parse("datatype t = A;")` returns an error containing `not supported`
- `Show(FunToFn(...))` for `fun add a b = a + b` gives `(fn (a (fn (b (app + (tuple a b))))))`
