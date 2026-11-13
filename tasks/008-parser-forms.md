# 008: Parser, part 3 (`if`, `fn`, `let`, simple patterns)

Depends on: 007.

## Task

Parse `if`, `fn`, and `let` expressions, and the simple pattern forms.

## Context

- Change: `parser/expr.go`, create `parser/pat.go`, `parser/parser_test.go`
- These forms start with a keyword and extend as far right as possible. Handle them at the top of `parseExpr`, before `parseBinary`:
  - `if` e1 `then` e2 `else` e3: `ast.If`. Each part is `parseExpr`.
  - `fn` pat `=>` e: `ast.Fn` with one arm. The body is `parseExpr`.
  - `let` decs `in` e1 `;` e2 ... `end`: `ast.Let`. `let` is an atom: add it to `parseAtom` and to the "can start an atom" check.
- For `let`, declarations are parsed by `parseDec`, which task 009 writes. For now add a stub `parseDec` that returns the error `declarations not supported yet`. Declarations may be separated by an optional `;`.
- Patterns (`parsePat`, in `pat.go`), simple forms only:
  - `_`: `ast.PWild`
  - IDENT: `ast.PVar`
  - `(` `)`: `ast.PTuple` with no elements
  - `(` p `)`: just p
  - `(` p1 `,` p2 ... `)`: `ast.PTuple`
- Add `func parseAtomicPat` with those forms and have `parsePat` call it. Later tasks extend both.
- These keywords are not supported yet. When `parseExpr` sees one, return the error `<keyword> not supported`: `case`, `raise`, `handle`.

## Done when

`go test ./parser/ -run TestForms -timeout 30s` passes:

- `if x > 5 then 1 else 2` gives `(if (app > (tuple x 5)) 1 2)`
- `fn a => fn b => a + b` gives `(fn (a (fn (b (app + (tuple a b))))))`
- `fn (a, _) => a` gives `(fn ((tuple a _) a))`
- `fn () => 1` gives `(fn (() 1))`
- `f (fn x => x) 1` gives `(app (app f (fn (x x))) 1)`
- `case x of y => y` returns an error containing `not supported`
