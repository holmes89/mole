# 017: `interp.SplitInput`

Depends on: 004.

## Task

Write the function that cuts source text into complete top-level declarations.

## Context

- Read: `lexer/lexer.go`, `token/token.go`
- Change: create `interp/split.go`, `interp/split_test.go`
- Signature: `func SplitInput(src string) (chunks []string, rest string, err error)`
- A chunk ends at a `;` token at depth 0, and includes that `;`. `rest` is whatever follows the last such `;`.
- Method: call `lexer.Tokenize(src)`. Walk the tokens it returns:
  - SYMBOL `(` or `[`, or KEYWORD `let`: depth + 1.
  - SYMBOL `)` or `]`, or KEYWORD `end`: depth - 1.
  - SYMBOL `;` at depth 0: cut. The chunk is `src[start : tok.Pos.Offset+1]`; then `start` moves past it.
- Because the lexer skips strings and comments, a `;` inside them never counts.
- If `Tokenize` returned an error:
  - `Incomplete` is true (unterminated string or comment): not an error here. Return the chunks found so far and the rest.
  - Otherwise: return the error.
- If depth goes below 0: return the error `unexpected closing bracket`.
- Trim leading whitespace from each chunk. `rest` may be only whitespace; callers trim it.

## Done when

`go test ./interp/ -run TestSplit -timeout 30s` passes:

- `1 + 2;` gives 1 chunk, empty rest
- `val x = 1; x;` gives 2 chunks: `val x = 1;` and `x;`
- `val x = 1; x` gives 1 chunk and rest `x` (after trimming)
- `(print "a"; 5);` gives 1 chunk
- `let val a = 1; val b = 2 in a end;` gives 1 chunk
- `print "a;b";` gives 1 chunk
- `(* x; y *) 5;` gives 1 chunk
- `"abc` gives 0 chunks, rest `"abc`, nil error
- `1 );` gives an error
- `fun f x =\n  x;\nf 1;` gives 2 chunks
