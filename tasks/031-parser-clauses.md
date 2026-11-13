# 031: `fun` with several clauses

Depends on: 030.

## Task

Parse `fun` declarations that have more than one clause.

## Context

- Read: `parser/dec.go`, `ast/desugar.go`
- Change: `parser/dec.go`, `parser/parser_test.go`, `ast/desugar_test.go`
- Syntax: `fun name pats = body | name pats = body | ...`
- After the first clause, while the next token is `|`: consume it, require the same function name (else error `clause name mismatch`), parse one or more atomic patterns, `=`, and the body.
- Every clause must have the same number of patterns, else error `clauses have different numbers of arguments`.
- Clause patterns are atomic, so constructor and cons patterns need parentheses: `fun area (Circle r) = ...`, `fun sum (x :: xs) = ...`. With task 030 in place, `parseAtomicPat` already accepts them inside parentheses.
- `ast.FunToFn` (task 009) already handles several clauses. Add tests for it.

## Done when

`go test ./parser/ ./ast/ -timeout 30s` passes, with new cases:

- `fun fib 0 = 0 | fib 1 = 1 | fib n = n;` shows `(fun fib ((0) 0) ((1) 1) ((n) n))`
- `fun sum [] = 0 | sum (x :: xs) = x + sum xs;` parses, with 2 clauses
- `fun area (Circle r) = r | area (Rect (w, h)) = w * h;` parses, with 2 clauses
- `fun f 0 = 1 | g n = 2;` is an error
- `fun f a b = 1 | f a = 2;` is an error
- `FunToFn` of a two-clause, two-argument function `fun f 0 y = y | f x y = x` shows `(fn ($1 (fn ($2 (app (fn ((tuple 0 y) y) ((tuple x y) x)) (tuple $1 $2))))))`
