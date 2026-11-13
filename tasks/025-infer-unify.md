# 025: Unification

Depends on: 023.

## Task

Create the `infer` package with fresh variables and unification.

## Context

- Read: `types/types.go`
- Change: create `infer/infer.go`, `infer/unify.go`, `infer/unify_test.go`

```go
package infer

type Infer struct {
	level int // current let-nesting level
	next  int // next variable ID
}

func New() *Infer
func (in *Infer) Fresh() *types.TVar          // new variable at the current level
func (in *Infer) Unify(a, b types.Type) error
```

- `Unify(a, b)`:
  1. `a = types.Prune(a)`, `b = types.Prune(b)`.
  2. Same pointer: done.
  3. `a` is a `*TVar`: run the occurs check and level adjustment on `b` (below), then set `a.Link = b`. Done.
  4. `b` is a `*TVar`: same, the other way around.
  5. Both `*TCon`: names and argument counts must match, else error; then unify the arguments pairwise.
- Occurs check and level adjustment, `in.occurs(v *TVar, t Type) error`: walk `t` (pruning). If you meet `v` itself, return the error `circular type`. For every other unbound variable you meet, set its `Level` to `min(its Level, v.Level)`. This level adjustment is what makes generalization correct in task 026.
- Mismatch error text: `type mismatch: <types.String(a)> vs <types.String(b)>`.

## Done when

`go test ./infer/ -run TestUnify -timeout 30s` passes:

- fresh `a` with `int`: afterward `types.String(a)` is `int`
- `int` with `bool`: error
- `a -> int` with `bool -> b`: afterward `a` prints `bool` and `b` prints `int`
- `a` with `a list`: error containing `circular`
- `a` with `a`: no error
- `(a * b)` with `(int * string)`: both resolved
- tuple of 2 with tuple of 3: error
- level adjustment: `a` created at level 1, `b` created at level 2 (set `in.level` by hand in the test); after `Unify(a, List(b))`, `b.Level` is 1
