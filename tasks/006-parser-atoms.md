# 006: Parser, part 1 (atoms and application)

Depends on: 004, 005.

## Task

Create the parser and parse atomic expressions and function application.

## Context

- Read: `token/token.go`, `ast/ast.go`
- Change: create `parser/parser.go`, `parser/expr.go`, `parser/parser_test.go`
- `parser.go`: a `parser` struct holding `toks []token.Token` and `pos int`, with helpers `peek()`, `next()`, `isSym(text)`, `isKw(text)`, `expectSym(text) error`.
- Error type: `type Error struct { Pos token.Pos; Msg string }`, printing as `line:col: msg`.
- Entry point for now: `func ParseExpr(src string) (ast.Expr, error)`. It calls `lexer.Tokenize`, parses one expression, and requires `EOF` after it.
- Atoms (`parseAtom`):
  - INT: `ast.IntLit`. Text `~5` means -5.
  - STRING: `ast.StrLit`.
  - IDENT or LONGID: `ast.Var`.
  - CON: `ast.Con`.
  - `(` `)`: `ast.Tuple` with no elements.
  - `(` e `)`: just e.
  - `(` e1 `,` e2 ... `)`: `ast.Tuple`.
  - `(` e1 `;` e2 ... `)`: `ast.Seq`.
  - `[` `]` and `[` e1 `,` e2 ... `]`: `ast.List`.
- Application (`parseApp`): one atom followed by zero or more atoms, left associative. `f x y` is `(app (app f x) y)`. Keep consuming while the next token can start an atom: INT, STRING, IDENT, LONGID, CON, `(`, `[`.
- Inside parentheses and brackets call `parseExpr`. For now `parseExpr` just calls `parseApp`; later tasks extend it.

## Constraints

- Hand-written recursive descent. No parser generator.

## Done when

`go test ./parser/ -run TestAtoms -timeout 30s` passes, comparing `ast.Show`:

- `42` gives `42`; `~5` gives `-5`
- `f x y` gives `(app (app f x) y)`
- `()` gives `()`
- `(1, "two", true)` gives `(tuple 1 "two" true)`
- `[1, 2]` gives `(list 1 2)`; `[]` gives `(list)`
- `(f 1; g 2)` gives `(seq (app f 1) (app g 2))`
- `Int.toString 3` gives `(app Int.toString 3)`
- `(1` and `1 )` return errors
