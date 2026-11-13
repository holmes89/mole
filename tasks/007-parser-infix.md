# 007: Parser, part 2 (infix operators)

Depends on: 006.

## Task

Add infix operators to `parseExpr` using precedence climbing.

## Context

- Change: `parser/expr.go`, `parser/parser_test.go`
- Precedence levels, loosest first:

| Level | Operators | Associativity |
|---|---|---|
| 1 | `orelse` | left |
| 2 | `andalso` | left |
| 3 | `o` | left |
| 4 | `=` `<>` `<` `>` `<=` `>=` | left |
| 5 | `::` `@` | right |
| 6 | `+` `-` `^` | left |
| 7 | `*` `div` `mod` | left |

- Application (task 006) binds tighter than every operator.
- `orelse`, `andalso`, `o`, `div`, `mod` are KEYWORD tokens. The rest are SYMBOL tokens.
- Algorithm: `parseBinary(minLevel)`: parse the left side with `parseApp`; while the next token is an operator with level >= `minLevel`, consume it and parse the right side with `parseBinary(level+1)` for left associative operators or `parseBinary(level)` for right associative ones.
- `parseExpr` now calls `parseBinary(1)`.
- AST produced:
  - `a orelse b`: `ast.OrElse`. `a andalso b`: `ast.AndAlso`.
  - Every other operator: `ast.App{Fn: ast.Var{Name: op}, Arg: ast.Tuple{Elems: [a, b]}}`.

## Done when

`go test ./parser/ -run TestInfix -timeout 30s` passes:

- `1 + 2 * 3` gives `(app + (tuple 1 (app * (tuple 2 3))))`
- `1 - 2 - 3` gives `(app - (tuple (app - (tuple 1 2)) 3))`
- `0 :: [1] @ [2]` gives `(app :: (tuple 0 (app @ (tuple (list 1) (list 2)))))`
- `f x + 1` gives `(app + (tuple (app f x) 1))`
- `a andalso b orelse c` gives `(orelse (andalso a b) c)`
- `x div 3 = 1` gives `(app = (tuple (app div (tuple x 3)) 1))`
- `f o g` gives `(app o (tuple f g))`
