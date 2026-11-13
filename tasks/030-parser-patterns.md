# 030: Full patterns, multi-arm `fn`, and `case`

Depends on: 029.

## Task

Parse the remaining pattern forms, `fn` with several arms, and `case`.

## Context

- Read: `parser/pat.go`, `parser/expr.go`
- Change: `ast/ast.go`, `ast/show.go`, `parser/pat.go`, `parser/expr.go`, `parser/parser_test.go`
- New nodes:

```go
type PInt struct { Node; Value int64 }
type PStr struct { Node; Value string }
type PCon struct { Node; Name string; Arg Pat }  // Arg is nil for Leaf, NONE, true
type PCons struct { Node; Head, Tail Pat }       // x :: xs
type PList struct { Node; Elems []Pat }          // [], [a, b]

type Case struct { Node; Scrut Expr; Arms []Arm }
```

- `ast.Show`: `PInt` and `PStr` like literals; `PCon` shows `Leaf` or `(pcon SOME x)`; `PCons` shows `(pcons h t)`; `PList` shows `(plist a b)`; `Case` shows `(case e (p body) (p body))`.
- Pattern grammar:
  - `parsePat`: `conPat [ "::" parsePat ]` (right associative; builds `PCons`)
  - `conPat`: CON followed by an atomic pattern gives `PCon` with `Arg`; otherwise `parseAtomicPat`
  - `parseAtomicPat` gains: INT (`~5` is -5), STRING, CON alone (`PCon` with nil `Arg`; `nil` gives an empty `PList`), `[` patterns `]`
  - A CON is followed by an argument when the next token can start an atomic pattern: `_`, IDENT, INT, STRING, CON, `(`, `[`.
- `fn`: after the first arm, while the next token is `|`, parse another `pat => expr` arm.
- `case e of pat => expr | pat => expr ...`: `ast.Case`. Remove `case` from the "not supported" list.
- An arm body is `parseExpr`, so it takes everything up to the next `|`. A `case` nested inside an arm swallows the arms that follow it; users add parentheses.

## Done when

`go test ./parser/ -timeout 30s` passes, with a new `TestPatterns`:

- `fn [] => 0 | _ => 1` gives `(fn ((plist) 0) (_ 1))`
- `fn (x :: xs) => x` gives `(fn ((pcons x xs) x))`
- `fn (a :: b :: rest) => a` gives `(fn ((pcons a (pcons b rest)) a))`
- `case v of NONE => 0 | SOME n => n + 1` gives `(case v (NONE 0) ((pcon SOME n) (app + (tuple n 1))))`
- `fn (Rect (w, h)) => w` gives `(fn ((pcon Rect (tuple w h)) w))`
- `fn 0 => "zero" | ~1 => "neg"` gives `(fn (0 "zero") (-1 "neg"))`
- `Parse("val h :: t = [7, 8];")` shows `(val (pcons h t) (list 7 8))`
- `Parse("val [p, q] = [1, 2];")` shows `(val (plist p q) (list 1 2))`
