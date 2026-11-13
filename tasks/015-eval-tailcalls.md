# 015: Tail calls run in constant stack

Depends on: 014.

## Task

Prove with a test that tail calls do not grow the Go stack, and fix `Eval` if they do.

## Context

- Read: `eval/eval.go`
- Change: create `eval/tail_test.go`; edit `eval/eval.go` only if the test fails
- The test program:

```
fun loop n acc = if n = 0 then acc else loop (n - 1) (acc + 1);
loop 1000000 0;
```

- Expected result: `1000000`, in under 5 seconds.
- Why it should already work: in `loop (n - 1) (acc + 1)` the outer application is in tail position. `Eval` evaluates its function part `loop (n - 1)` with a normal recursive call that returns a closure, then handles the outer application by setting `e` and `en` and using `continue`.
- If the test overflows the stack or is slow, look for a tail position that calls `Eval` recursively and returns its result. Tail positions are: both branches of `if`, the body of a matched arm, the last expression of a `let` body, the last expression of a sequence, and the right side of `andalso` and `orelse`.

## Constraints

- Do not raise the stack limit or add goroutines. The fix is always "assign and `continue`".

## Done when

`go test ./eval/ -run TestTailCall -timeout 30s` passes. Also add a second case, `fun count n = if n = 0 then 0 else count (n - 1); count 1000000;`, giving `0`.
