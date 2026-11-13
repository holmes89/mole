# 028: Infer declarations (let-polymorphism)

Depends on: 027.

## Task

Implement type inference for `val` and `fun` declarations, with generalization.

## Context

- Read: `infer/expr.go`, `infer/scheme.go`, `ast/desugar.go`
- Change: create `infer/dec.go`, `infer/dec_test.go` (replace the stub from task 027)
- `func (in *Infer) Dec(d ast.Dec, tenv *env.Env) ([]TBinding, error)` defines the declared names in `tenv` and returns them in order.
- `ValDec`:
  1. `in.level++`
  2. `t := in.Expr(d.Expr, tenv)`; `pt, binds := in.Pat(d.Pat, tenv)`; `Unify(pt, t)`
  3. `in.level--` (also on every error path; use `defer` or a helper)
  4. For each binding: `in.Generalize(binding.Type)`, then `tenv.Define`.
- `FunDec`:
  1. `in.level++`
  2. `self := in.Fresh()`; `frame := env.New(tenv)`; `frame.Define(d.Name, self)`. The name is not generic inside its own body.
  3. `t := in.Expr(ast.FunToFn(d), frame)`; `Unify(self, t)`
  4. `in.level--`
  5. `in.Generalize(t)`; `tenv.Define(d.Name, t)`; return one binding.
- Why levels: variables created while inferring the right-hand side have a level above `in.level` after step 3 or 4, unless unification tied them to something outside (task 025 lowers their level). Only the untied ones become generic.
- `Let` in `expr.go` now works through `in.Dec`.
- Test helper: `typesOf(src)` parses a program, runs `in.Dec` on each declaration in one `tenv`, and returns the printed type of each binding, or the error.

## Done when

`go test ./infer/ -timeout 30s` passes, with a new `TestDec`:

- `fun id x = x;` is `'a -> 'a`
- `fun id x = x; (id 1, id true);` ends with `int * bool`
- `fun double x = x * 2;` is `int -> int`
- `fun compose f g x = f (g x);` is `('a -> 'b) -> ('c -> 'a) -> 'c -> 'b`
- `fun len xs = if null xs then 0 else 1 + len (tl xs);` is `'a list -> int`
- `fun eq (a, b) = a = b;` is `'a * 'a -> bool`
- `let val f = fn x => x in (f 1, f "a") end;` is `int * string`
- `val (a, b) = (1, "x");` gives `int` and `string`
- `fun double x = x * 2; val quad = double o double;` ends with `int -> int`
- `fun bad x = x x;` is an error
- `val y = 1 + true;` is an error, and a following `y;` is an error too
