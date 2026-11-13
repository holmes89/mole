# 013: Evaluator, part 1 (expressions without closures)

Depends on: 009, 011, 012.

## Task

Create `eval.Eval` for literals, variables, tuples, lists, `if`, `andalso`, `orelse`, sequences, and calls to builtins.

## Context

- Read: `ast/ast.go`, `value/value.go`, `env/env.go`
- Change: create `eval/eval.go`, `eval/eval_test.go`
- Signatures:

```go
func Eval(e ast.Expr, en *env.Env) (value.Value, error)
func Apply(fn, arg value.Value) (value.Value, error)
```

- Write `Eval` as a loop from the start: `for { switch e := e.(type) { ... } }`. A case either returns, or assigns a new `e` (and maybe `en`) and uses `continue`. That is how tail positions avoid growing the Go stack. Task 015 depends on this shape.
- Cases:
  - `IntLit`, `StrLit`: return `value.Int`, `value.Str`.
  - `Var`, `Con`: `en.Get(name)`. Missing: error `unbound variable: x` (or `unbound constructor: X`).
  - `Tuple`: zero elements is `value.Unit{}`; otherwise evaluate left to right into `value.Tuple`.
  - `List`: evaluate left to right, build with `value.FromSlice`.
  - `If`: evaluate `Cond` (must be `value.Bool`, else error), then set `e` to the chosen branch and `continue`.
  - `AndAlso`: evaluate `L`; if false return false; else `e = R; continue`. `OrElse` is the mirror image.
  - `Seq`: evaluate all but the last; then `e = last; continue`.
  - `App`: evaluate `Fn`, then `Arg`. If the function is a `value.Builtin`, return `b.Fn(arg)`. If it is a `value.ConFn`, return `value.Con{Name, ID, Arg: arg}`. Anything else: error `not a function`.
- `Apply`: the same builtin and `ConFn` logic, as a standalone function.
- Test helper: `newEnv()` makes `env.New(nil)` and defines every entry of `builtins.Table(Apply)`. `run(src)` uses `parser.ParseExpr` then `Eval`, and returns `value.String` of the result.

## Constraints

- Do not handle `Fn` or `Let` yet (task 014). Return the error `not supported yet` for them.

## Done when

`go test ./eval/ -run TestBasic -timeout 30s` passes: `1 + 2 * 3` is `7`; `"mo" ^ "le"` is `"mole"`; `if 10 > 5 then "big" else "small"` is `"big"`; `(1, "two", true)` is `(1,"two",true)`; `0 :: [1, 2] @ [3]` is `[0,1,2,3]`; `true andalso (false orelse true)` is `true`; `false andalso (1 div 0 = 1)` is `false`; `()` is `()`; `(1; 2)` is `2`; `1 div 0` returns a `*value.Raise`; `nope` returns an error; `Fail "x"` is `Fail "x"`.
