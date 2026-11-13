# 004: Lexer, part 2 (strings, comments, `~`, long names, type variables)

Depends on: 003.

## Task

Extend `lexer.Tokenize` with the remaining token forms.

## Context

- Change: `lexer/lexer.go`, `lexer/lexer_test.go`
- Comments: `(*` to `*)`. They nest. Skip them. A comment that never closes is an `*Error` with `Incomplete: true`.
- STRING: starts and ends with `"`. Escapes: `\n`, `\t`, `\\`, `\"`. `Text` is the unescaped contents. Any other escape is an error. A string that reaches end of input (or a newline) without closing is an `*Error` with `Incomplete: true`.
- `~`: when the next character is a digit, lex one INT token whose `Text` includes the `~` (`~5`). Otherwise emit IDENT `~`.
- LONGID: a name starting with an uppercase letter, then `.`, then a name, with no spaces (`Int.toString`, `List.filter`). One token, kind `LONGID`, `Text` is the whole thing.
- TYVAR: `'` followed by a letter and then letters or digits (`'a`, `'key`). `Text` includes the `'`.
- Check for `(*` before treating `(` as a symbol.

## Constraints

- Keep every test from task 003 passing.

## Done when

`go test ./lexer/ -timeout 30s` passes, with a new `TestRest` table:

- `"mo\"le\n"` gives one STRING with text `mo"le` plus a newline character
- `(* a (* b *) c *) 5` gives INT `5`, EOF
- `(* open` gives an error with `Incomplete` true
- `"abc` gives an error with `Incomplete` true
- `~5 + ~ x` gives INT `~5`, SYMBOL `+`, IDENT `~`, IDENT `x`, EOF
- `Int.toString ~42` gives LONGID `Int.toString`, INT `~42`, EOF
- `'a list` gives TYVAR `'a`, IDENT `list`, EOF
- `(1)` gives SYMBOL `(`, INT, SYMBOL `)`, EOF
