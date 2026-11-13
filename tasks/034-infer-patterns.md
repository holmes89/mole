# 034: Infer full patterns and `case`

Depends on: 033.

## Task

Extend type inference to the new patterns and to `case`.

## Context

- Read: `infer/pat.go`, `infer/expr.go`
- Change: `infer/pat.go`, `infer/expr.go`, create `infer/pattern_test.go`
- `Pat` additions:

| Pattern | Type | Notes |
|---|---|---|
| `PInt`, `PStr` | `int`, `string` | no bindings |
| `PList` | `List(elem)` | `elem := Fresh()`; unify each element's type with `elem` |
| `PCons` | `List(th)` | `th` is the head's type; unify the tail's type with `List(th)` |
| `PCon` without `Arg` | instantiated type of the constructor | look the name up in `tenv`; missing: error `unbound constructor: X`. If the instantiated type is a function type, error `constructor X needs an argument` |
| `PCon` with `Arg` | `r` | instantiated constructor type must unify with `Arrow(argType, r)` where `r := Fresh()`; on failure: error `constructor X takes no argument` |

- `true` and `false` are already in `tenv` as `bool`, and `NONE` and `SOME` are there from task 033.
- The duplicate-variable check from task 027 must cover variables nested anywhere in the pattern.
- `Case` in `Expr`: `ts := type of scrutinee`; `res := Fresh()`; for each arm: frame on `tenv`, `Pat`, unify the pattern type with `ts`, define the bindings (not generalized), infer the body, unify with `res`. Result `res`. This is the same as `Fn` with the scrutinee type in place of the argument type.

## Done when

`go test ./infer/ -timeout 30s` passes, with a new `TestPatternTypes`:

- `fun sum [] = 0 | sum (x :: xs) = x + sum xs;` is `int list -> int`
- `fun first [] = NONE | first (x :: _) = SOME x;` is `'a list -> 'a option`
- `fun fib 0 = 0 | fib 1 = 1 | fib n = fib (n - 1) + fib (n - 2);` is `int -> int`
- `fun swap (a, b) = (b, a);` is `'a * 'b -> 'b * 'a`
- `fun zero 0 = "zero";` is `int -> string`
- `case SOME 3 of NONE => 0 | SOME n => n + 1;` is `int`
- `SOME (SOME ~1);` is `int option option`
- `val h :: t = [7, 8, 9];` gives `int` and `int list`
- errors: `case 1 of true => 0 | false => 1;`; `fun dup (x, x) = x;`; `case 1 of Nope => 0;`; `fn (SOME) => 1;`; `fn (NONE x) => 1;`
