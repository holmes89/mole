# 026: Generalize, instantiate, and builtin types

Depends on: 024, 025.

## Task

Add let-polymorphism helpers and load builtin type strings into a type environment.

## Context

- Read: `infer/infer.go`, `types/types.go`, `ast/ast.go` (the `Ty` nodes)
- Change: create `infer/scheme.go`, `infer/scheme_test.go`
- A type environment is an `*env.Env` whose values are `types.Type`. A polymorphic type simply contains variables whose `Level` is `types.Generic`.
- `func (in *Infer) Generalize(t types.Type)`: walk `t` (pruning). Every unbound variable whose `Level` is greater than `in.level` gets `Level = types.Generic`.
- `func (in *Infer) Instantiate(t types.Type) types.Type`: return a copy of `t` in which each distinct generic variable is replaced by one fresh variable (use a map from old variable to new). Non-generic variables and constructors without generic variables inside may be shared.
- `func (in *Infer) FromAst(ty ast.Ty, vars map[string]*types.TVar) (types.Type, error)`: convert syntax to a type.
  - `TyVar`: look up the name in `vars`; if missing, create a fresh variable and store it.
  - `TyCon`: check the name and argument count against the known type constructors, else error `unknown type: <name>`. Keep them in a map on `Infer`: `tycons map[string]int`, starting with `int`, `bool`, `string`, `unit`, `exn` (0 arguments) and `list`, `option` (1 argument).
  - `TyTuple`, `TyArrow`: convert the parts.
- `func (in *Infer) Declare(tenv *env.Env, name, typeString string) error`: `parser.ParseType`, then `FromAst` with a new `vars` map, then mark every variable in the result as generic, then `tenv.Define(name, t)`.

## Done when

`go test ./infer/ -run TestScheme -timeout 30s` passes:

- `Declare(tenv, "map", "('a -> 'b) -> 'a list -> 'b list")`, then two `Instantiate` calls on the stored type give types with different variable pointers; unifying the first one's `'a` with `int` leaves the second printing `('a -> 'b) -> 'a list -> 'b list`
- with `in.level = 0`: a variable made at level 1 becomes generic after `Generalize`; a variable made at level 0 does not
- `Declare(tenv, "x", "foo")` returns an error containing `unknown type`
- every entry of `builtins.Table(nil, nil)` declares without error
