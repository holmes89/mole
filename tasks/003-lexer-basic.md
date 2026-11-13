# 003: Lexer, part 1 (integers, names, symbols)

Depends on: 002.

## Task

Implement `lexer.Tokenize` for integers, identifiers, keywords, constructors, and symbols.

## Context

- Read: `token/token.go`
- Change: create `lexer/lexer.go`, `lexer/lexer_test.go`
- Signature: `func Tokenize(src string) ([]token.Token, error)`
- The result always ends with one `EOF` token. On an error, return the tokens lexed so far (no `EOF`) together with the error.
- Skip spaces, tabs, and newlines. Track `Line`, `Col`, and `Offset` for the first character of each token.
- INT: one or more decimal digits.
- Names: a letter, then letters, digits, `_`, or `'`.
  - In `token.Keywords`: kind `KEYWORD`.
  - `true`, `false`, `nil`, or starts with an uppercase letter: kind `CON`.
  - Otherwise: kind `IDENT`.
- SYMBOL, longest match first: `=>` `->` `<>` `<=` `>=` `::` then `( ) [ ] , ; | _ = < > + - * ^ @`.
- Any other character is an error. Define:

```go
type Error struct {
	Pos        token.Pos
	Msg        string
	Incomplete bool // set in task 004
}
func (e *Error) Error() string // "line:col: msg"
```

## Constraints

- No regexp. Do not handle strings, comments, `~`, type variables, or dots yet (task 004).

## Done when

`go test ./lexer/ -run TestBasic -timeout 30s` passes. Table cases, comparing kinds and texts:

- `val x = 10;` gives KEYWORD `val`, IDENT `x`, SYMBOL `=`, INT `10`, SYMBOL `;`, EOF
- `a<=b` gives IDENT, SYMBOL `<=`, IDENT, EOF
- `x :: xs` gives IDENT, SYMBOL `::`, IDENT, EOF
- `fn a => a` gives KEYWORD, IDENT, SYMBOL `=>`, IDENT, EOF
- `Leaf true foo'` gives CON, CON, IDENT `foo'`, EOF
- `a\n  b`: the second token has Line 2, Col 3, Offset 4
- `#` returns an error
