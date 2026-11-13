# 002: Token kinds

## Task

Create the `token` package.

## Context

- Change: create `token/token.go`, `token/token_test.go`
- Write exactly these declarations:

```go
package token

type Kind int

const (
	EOF     Kind = iota
	INT          // 42 or ~5 (Text keeps the ~)
	STRING       // Text is the unescaped contents, without the quotes
	IDENT        // x, foo', ~
	CON          // Leaf, SOME, true, false, nil
	LONGID       // Int.toString
	TYVAR        // 'a
	KEYWORD      // see Keywords
	SYMBOL       // ( ) [ ] , ; | _ => -> = <> < > <= >= + - * ^ :: @
)

// Pos is a source position. Line and Col start at 1. Offset is a byte offset from 0.
type Pos struct{ Line, Col, Offset int }

type Token struct {
	Kind Kind
	Text string
	Pos  Pos
}

var Keywords = map[string]bool{ /* fill in */ }
```

- `Keywords` holds: `val fun fn let in end if then else andalso orelse case of datatype exception raise handle div mod o`.
- Add `func (k Kind) String() string` returning the constant's name (`"INT"`, `"SYMBOL"`, ...).

## Done when

`go test ./token/ -timeout 30s` passes, with a test that `Keywords["andalso"]` is true, `Keywords["foo"]` is false, `len(Keywords)` is 20, and `INT.String()` is `"INT"`.
