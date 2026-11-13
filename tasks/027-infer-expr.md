# 027: Infer expression types

Depends on: 026.

## Task

Implement type inference for expressions.

## Context

- Read: `infer/infer.go`, `infer/scheme.go`, `ast/ast.go`
- Change: create `infer/expr.go`, `infer/pat.go`, `infer/expr_test.go`
- `func (in *Infer) Expr(e ast.Expr, tenv *env.Env) (types.Type, error)`:

| Node | Rule |
|---|---|
| `IntLit`, `StrLit` | `int`, `string` |
| `Var`, `Con` | look up in `tenv`; missing is the error `unbound variable: x`; return `Instantiate` of it |
| `Tuple` | no elements: `unit`; else `Tuple(types of elements)` |
| `List` | `elem := Fresh()`; unify each element's type with `elem`; result `List(elem)` |
| `App` | `tf`, `ta` from the two parts; `r := Fresh()`; `Unify(tf, Arrow(ta, r))`; result `r` |
| `If` | condition unifies with `bool`; the branches unify with each other; result is the branch type |
| `AndAlso`, `OrElse` | both sides unify with `bool`; result `bool` |
| `Seq` | infer each; result is the last one's type |
| `Fn` | `arg := Fresh()`, `res := Fresh()`; for each arm: `frame := env.New(tenv)`; `pt, binds := in.Pat(arm.Pat, frame)`; `Unify(pt, arg)`; define each binding in `frame` (not generalized); infer the body in `frame` and unify with `res`. Result `Arrow(arg, res)` |
| `Let` | `frame := env.New(tenv)`; for each declaration call `in.Dec(dec, frame)` (task 028; add a stub returning `not supported yet`); result is the type of the last body expression in `frame` |

- `pat.go`:

```go
type TBinding struct { Name string; Type types.Type }
func (in *Infer) Pat(p ast.Pat, tenv *env.Env) (types.Type, []TBinding, error)
```

  - `PWild`: `Fresh()`, no bindings. `PVar`: `t := Fresh()`, one binding. `PTuple`: no elements is `unit`; else the tuple of the parts, bindings in order.
  - If the same variable name appears twice in one pattern, return the error `duplicate variable in pattern: x`.
- Wrap unification errors with the position: `line:col: <message>`.
- Test helper: `tenv` built with `Declare` for every entry of `builtins.Table(nil, nil)`; `typeOf(src)` uses `parser.ParseExpr`, `in.Expr`, and `types.String`.

## Done when

`go test ./infer/ -run TestExpr -timeout 30s` passes:

- `42` is `int`; `"a" ^ "b"` is `string`; `(1, true)` is `int * bool`
- `[]` is `'a list`; `[[1, 2], [3]]` is `int list list`
- `fn x => x` is `'a -> 'a`
- `fn (a, b) => (b, a)` is `'a * 'b -> 'b * 'a`
- `fn x => x + 1` is `int -> int`
- `map (fn n => n * 2) [1, 2]` is `int list`
- `List.filter` is `('a -> bool) -> 'a list -> 'a list`
- `foldr (fn (x, acc) => x :: acc) [] [1, 2, 3]` is `int list`
- errors: `if true then 1 else "no"`; `1 + true`; `fn x => x x`; `fn f => (f 1, f "a")`; `fn (x, x) => x`; `nope`
