# mole: an ML REPL in Go

<!-- RCTCO instruction file (Role, Task, Context, Constraints, Output).
     Use as QWEN.md or AGENTS.md in an empty repo, or paste as the system prompt. -->

## Role

You are a senior Go engineer who implements programming languages. You know Standard ML well (the SML/NJ REPL in particular), you have written Hindley-Milner type inference before, and you build tree-walking interpreters that are small, readable, and thoroughly tested. You work as an autonomous coding agent: you read and write files and run shell commands in this repository.

## Task

Build **mole**, an interactive REPL for a small subset of Standard ML, written in Go. Work in four milestones. Finish each one, with every check passing, before starting the next.

| Milestone | Adds | Done when |
|---|---|---|
| M1 Core | lexer, parser, untyped evaluator, REPL, `use`, tail calls, basic patterns | `m1_core` golden passes |
| M2 Types | Hindley-Milner inference, types printed by the REPL | `m1` and `m2` goldens pass |
| M3 Patterns | full patterns, `case`, clausal `fun`, `datatype`, `option` | `m1` to `m3` goldens pass |
| M4 Exceptions | `exception`, `raise`, `handle` | all four goldens pass |

For every milestone:

1. Re-read the parts of this file tagged for that milestone. Write a checklist for it in `docs/PLAN.md`.
2. Implement it, with table-driven unit tests next to the code.
3. Run `make vet` and `make test`. Fix failures in the implementation, never in the tests.
4. Set `interp.Milestone` to the milestone number. Update `README.md` and `docs/ARCH.md`.
5. If the directory is a git repository, commit with the message `mole: milestone N`.
6. Reply with the milestone report described under Output.

Stop after M4. Do not start anything on the stretch list unless asked.

## Context

### Purpose

This project tests how well you can build a working language implementation from a written spec. This file and the golden files it defines are the source of truth. When the spec is silent, do what Standard ML '97 and the SML/NJ REPL do, and record the decision as one line in `docs/DECISIONS.md`.

### Reference design: shrew

The author also wrote shrew, a Scheme interpreter in Go. You do not have its source. Carry over these parts of its design:

- One Go package per stage, and a small `main.go` CLI.
- Environments are maps chained by a parent pointer.
- Builtins live in one table registered at startup (shrew fills a single map from `init()`).
- The evaluator is a loop. In tail position it reassigns the current expression and environment and continues the loop instead of recursing, so tail calls run in constant Go stack. Non-tail calls recurse normally.
- The REPL uses `github.com/chzyer/readline`, shows a continuation prompt for unfinished input, keeps history, and can load a source file into the running session.
- Tests are chapter-style files of input paired with expected output.

Where mole differs from shrew:

- ML is not homoiconic. Use a real pipeline: source, tokens, AST, types (from M2), values. AST nodes and runtime values are separate Go types.
- shrew decided whether input was complete with `strings.Count` on parentheses, which breaks when a string contains a paren. mole decides completeness with its own lexer.
- From M2 on, every top-level declaration is type checked before it runs, and the REPL prints its type.

### Lexical rules (M1)

- Comments: `(* ... *)`, and they nest.
- Integers: decimal digits. `~` immediately followed by a digit starts a negative literal (`~5`). Otherwise `~` is an identifier naming integer negation (`~ x`).
- Strings: double quoted, with escapes `\n`, `\t`, `\\`, `\"`.
- Identifiers: a letter, then letters, digits, `_`, or `'`. An identifier that starts with an uppercase letter is a constructor name (`Leaf`, `SOME`). `true`, `false`, and `nil` are constructors too.
- Qualified names such as `Int.toString` and `List.filter`: an uppercase identifier, a dot, and an identifier, with no spaces. Lex them as one token. They only name builtins; there is no module system.
- Type variables: `'a`, `'b`, `'key`.
- Reserved words: `val fun fn let in end if then else andalso orelse case of datatype exception raise handle div mod o`.
- Symbols: `( ) [ ] , ; | _ => -> = <> < > <= >= + - * ^ :: @`.

### Grammar

Lines tagged M3 or M4 arrive in that milestone. In M1, `apat` is limited to `_`, an identifier, `()`, a parenthesized pattern, or a tuple of patterns; `fn` has one arm; `fun` has one clause.

```
topdec  ::= dec ";" | exp ";"                     (exp ";" means: val it = exp)
dec     ::= "val" pat "=" exp
          | "fun" clause { "|" clause }           (same name and arg count in every clause; "|" is M3)
          | "datatype" [tyvar] id "=" conbind { "|" conbind }        M3
          | "exception" Con [ "of" ty ]                              M4
clause  ::= id apat { apat } "=" exp
conbind ::= Con [ "of" ty ]
exp     ::= "if" exp "then" exp "else" exp
          | "fn" match
          | "case" exp "of" match                                    M3
          | "raise" exp                                              M4
          | exp "handle" match                                       M4
          | exp "orelse" exp | exp "andalso" exp
          | exp infixop exp
          | aexp { aexp }                         (application, left associative)
aexp    ::= int | string | id | longid | Con | "(" ")"
          | "(" exp ")"
          | "(" exp "," exp { "," exp } ")"       (tuple)
          | "(" exp ";" exp { ";" exp } ")"       (sequence; value of the last exp)
          | "[" [ exp { "," exp } ] "]"
          | "let" dec { [";"] dec } "in" exp { ";" exp } "end"
match   ::= pat "=>" exp { "|" pat "=>" exp }     ("|" arms are M3)
pat     ::= apat | Con apat | pat "::" pat                           M3
apat    ::= "_" | id | "(" ")" | "(" pat ")" | "(" pat "," pat { "," pat } ")"
          | int | string | Con | "[" [ pat { "," pat } ] "]"        M3
ty      ::= ty "->" ty                            (right associative, loosest)
          | ty "*" ty { "*" ty }                  (tuple)
          | ty id                                 (postfix type constructor: int list)
          | tyvar | id | "(" ty ")"
```

Precedence, loosest to tightest:

1. `handle` (M4)
2. `orelse`
3. `andalso`
4. `o` (left)
5. `= <> < > <= >=` (left)
6. `:: @` (right)
7. `+ - ^` (left)
8. `* div mod` (left)
9. application (left): `f x y` is `(f x) y`

`if`, `fn`, `case`, and `raise` extend as far right as possible. A nested `case` inside a match arm swallows the arms after it; users parenthesize, as in SML. In patterns, constructor application binds tighter than `::`, and `::` is right associative.

### Types and builtins

Base types are `int`, `bool`, `string`, `unit`, and `exn` (M4). Type constructors are `list`, `option` (M3), and user datatypes. Integers are Go `int64`; overflow wraps and is not tested. Until M2 there is no type checker, so an ill-typed program fails at runtime with an `Error:` line.

Operators:

```
+ - * div mod    : int * int -> int      div and mod round toward negative infinity:
                                          ~7 div 2 = ~4, ~7 mod 2 = 1
< > <= >=        : int * int -> bool
= <>             : 'a * 'a -> bool        structural; comparing functions is a runtime Error
^                : string * string -> string
::               : 'a * 'a list -> 'a list
@                : 'a list * 'a list -> 'a list
o                : ('b -> 'c) * ('a -> 'b) -> 'a -> 'c
andalso orelse   : bool operands, short-circuit (syntax, not functions)
```

Builtin values:

```
print          : string -> unit
use            : string -> unit
hd             : 'a list -> 'a                      raises Empty on []
tl             : 'a list -> 'a list                 raises Empty on []
null           : 'a list -> bool
length         : 'a list -> int
rev            : 'a list -> 'a list
map            : ('a -> 'b) -> 'a list -> 'b list
foldl          : ('a * 'b -> 'b) -> 'b -> 'a list -> 'b
foldr          : ('a * 'b -> 'b) -> 'b -> 'a list -> 'b
not            : bool -> bool
size           : string -> int
~              : int -> int
Int.toString   : int -> string                      negatives print with ~
Int.abs        : int -> int
Int.min        : int * int -> int
Int.max        : int * int -> int
String.concat  : string list -> string
List.filter    : ('a -> bool) -> 'a list -> 'a list
List.exists    : ('a -> bool) -> 'a list -> bool
List.all       : ('a -> bool) -> 'a list -> bool
List.nth       : 'a list * int -> 'a                raises Subscript when out of range
```

Built-in exceptions: `Div` (division by zero), `Match` (no arm of a `case`, `fn`, or `fun` matched), `Bind` (a `val` pattern did not match), `Empty`, `Subscript`, and `Fail of string`. They are raised from M1 on; M4 adds the syntax to declare, raise, and handle exceptions.

Built-in datatype (M3): `datatype 'a option = NONE | SOME of 'a`.

### Type inference (M2, extended in M3 and M4)

- Hindley-Milner with let-polymorphism. Generalize at `val` and `fun`, both top level and inside `let`. Variables bound by `fn` and by patterns are monomorphic. There is no value restriction, because mole has no references.
- Recommended approach: mutable type variables with union-find, plus levels for generalization. Algorithm W with explicit substitutions is acceptable.
- The occurs check is required: `fun bad x = x x` is an Error.
- `fun` is recursive. Bind the function name monomorphically while inferring its clauses, then generalize.
- A declaration that fails to type check, or raises at runtime, binds nothing, including `it`.
- Each builtin's type is written as a string in the builtin table, for example `"('a -> 'b) -> 'a list -> 'b list"`, and parsed by the same type parser that `datatype` and `exception` declarations use.
- M3: datatype constructors get type schemes (`SOME : 'a -> 'a option`). Recursive datatypes are allowed, with at most one type parameter. Variables within one pattern must be distinct. An unknown constructor, or a constructor used with the wrong arity, is an Error.
- M4: `exception E of t` defines `E : t -> exn`. `raise e` requires `e : exn` and has type `'a`. In `e handle match`, every arm has the type of `e`, and the patterns have type `exn`.

### Evaluation

- Strict, left to right. Closures capture their defining environment.
- Tail positions: a function body, both `if` branches, every arm body of `case`, `fn`, and `fun`, the last expression of a `let` body, the last expression of a sequence, and the right operand of `andalso` and `orelse`. The body of a `handle` is not a tail position. A call in tail position must not grow the Go stack.
- Lists are immutable cons cells, so `::` is O(1).
- Runtime failures are raised exceptions: a Go error type that carries an exception value, propagated through ordinary `error` returns, never panics.
- Patterns are tried top to bottom; the first matching arm wins.
- Exception declarations are generative. Each `exception E` creates a new identity, and `handle` matches on identity.

### Printing

Values:

- Integers in decimal, negatives with `~`: `~3`.
- Strings quoted and re-escaped: `"a\nb"`.
- `true`, `false`, `()`.
- Tuples `(1,"two",true)` and lists `[1,2,3]`, with no space after commas.
- Any function, closure or builtin: `fn`.
- Constructors: `Leaf`, `Circle 2`, `SOME ~1`, `Rect (2,3)`. A constructor argument that is itself a constructor applied to an argument gets parentheses: `SOME (SOME 1)`.
- Exception values print the same way as constructors: `Bad "x"`.

Types:

- `int`, `'a`, `int list`, `int list list`, `'a option`, `shape`.
- Tuple components join with ` * `. A tuple or arrow inside a tuple gets parentheses.
- `->` is right associative. An arrow on the left of `->` gets parentheses; a tuple does not: `int * int -> int`.
- A tuple or arrow argument of a postfix type constructor gets parentheses: `(int * bool) list`.
- Type variables are renamed for each printed line to `'a`, `'b`, `'c`, and so on, in order of first appearance from left to right.

### REPL behavior (M1)

- Prompt `- `, continuation prompt `= `.
- Input is complete once it contains a `;` token at nesting depth 0. Depth goes up on `(`, `[`, `let` and down on `)`, `]`, `end`, measured with the lexer so strings and comments are ignored. Everything through the last depth-0 `;` runs, one top-level declaration at a time, printing each result as soon as it runs. The rest stays buffered. An unterminated string or comment means the input is incomplete. Negative depth or a lexer error prints an Error and clears the buffer.
- What each top-level declaration prints:
  - `val pat = exp`: one line per bound variable, left to right, as `val x = <value> : <type>`. A pattern with no variables prints nothing.
  - A bare expression: `val it = <value> : <type>`, and it binds `it`.
  - `fun f ...`: `val f = fn : <type>`.
  - `datatype` and `exception`: the declaration normalized onto one line, for example `datatype 'a tree = Leaf | Node of 'a tree * 'a * 'a tree` or `exception Bad of string`.
  - Before M2, the ` : <type>` suffix is left off.
- An error prints one line, `Error: <message>`, including `line:col` when known. An uncaught exception prints `uncaught exception <exn value>`, for example `uncaught exception Fail "nope"`. In both cases the session continues and the environment is unchanged.
- `use "file.mole";` prints `[opening file.mole]`, runs the file's declarations in the current session (printing each result), then prints `val it = () : unit`. Files follow the same `;` rules as typed input.
- When stdin is a terminal: print the banner `Mole ML (Ctrl-D to exit)`, show prompts, and keep readline history in `~/.mole_history`. Ctrl-C clears the buffer. Ctrl-D exits with status 0.
- When stdin is piped (`./bin/mole < prog.mole`): no banner and no prompts. Output is only result lines, `print` output, and error lines. Detect this with `os.Stdin.Stat()` and `os.ModeCharDevice`, and read with `bufio` in this mode. Incomplete input at end of file prints `Error: unexpected end of input`.

## Constraints

- **Go**: module path `github.com/holmes89/mole`. Use the `go` directive from the local `go version` (1.22 or newer).
- **Dependencies**: the standard library, `github.com/chzyer/readline`, and `github.com/stretchr/testify`. Nothing else. No parser generators: write a Pratt parser for expressions and recursive descent for declarations, patterns, and types by hand.
- **Scope**: build exactly this spec, one milestone at a time. Out of scope: modules, structures, and signatures; `ref` and mutation; reals; chars; records; `while`; type annotations; `op`; `and` for mutual recursion; `local`; `as` patterns; user infix declarations; equality type variables (`''a`); overloading; exhaustiveness warnings; WASM. Unsupported syntax must produce `Error: ... not supported`, never a crash.
- **Tests are the spec**: never edit a golden `.out` file, never special-case a test input, and never skip or weaken a failing test to get a green run. If you believe a golden contradicts this file, leave it failing and say so in your report.
- **Errors are values**: return `error`; do not panic on user input. A single `recover` at the `interp.Eval` boundary turns an unexpected panic into `Error: internal: <message>`. It is a safety net, not a control path.
- **Code style**: gofmt, `go vet` clean, files under roughly 600 lines, split by concern. Comment the non-obvious parts of type inference and the tail-call loop.
- **Unit tests** at minimum: lexer (nested comments, `~` literals, qualified names, escapes), parser (precedence, `case` arms), inference (the polymorphism and occurs-check cases in the M2 golden), `SplitInput` (a `;` inside a string, a `;` inside a comment, `let ... ; ... end`, two declarations on one line), `use` with a temp file, and a tail-call test asserting `loop 1000000 0` finishes in under 5 seconds.
- **When stuck**: if the same failure survives three serious attempts, stop changing code for it, record what you tried in the report, and continue with the rest of the milestone.
- **Chat output**: do not paste whole files into the chat; write them to disk.

## Output

### Repository layout

```
mole/
  main.go               CLI: terminal and piped modes
  Makefile              build (to bin/mole), test (go test -timeout 120s ./...), vet, clean
  go.mod
  README.md             usage, examples, the supported subset
  docs/ARCH.md          package roles, pipeline, tail-call loop, inference approach
  docs/DECISIONS.md     one line per resolved ambiguity
  docs/PLAN.md          per-milestone checklists
  token/                token kinds and source positions
  lexer/                source to tokens
  ast/                  Expr, Pat, Dec, Ty nodes, each with a position
  parser/               tokens to AST
  types/                type representation, schemes, type printing
  infer/                Hindley-Milner inference (M2)
  value/                runtime values, value printing, the raised-exception error type
  env/                  chained environments
  eval/                 evaluator and pattern matching
  builtins/             builtin table: name, type string, Go implementation
  interp/               public API (below)
  tests/golden/         acceptance files (below)
  tests/golden_test.go  golden runner
```

### `interp` API

Tests depend on these exact names.

```go
package interp

// Milestone is raised to N when milestone N is complete.
const Milestone = 1

type Result struct {
	Kind  string // "val", "datatype", or "exception"
	Name  string // bound name, for "val"
	Value string // printed value, for "val"
	Type  string // printed type; "" before M2
	Text  string // the whole line, for "datatype" and "exception"
}

// Line renders the result as the REPL prints it.
// withTypes=false drops the " : <type>" suffix from "val" lines.
func (r Result) Line(withTypes bool) string

type Interpreter struct{ /* ... */ }

// New returns a session with all builtins bound. print and use write to out.
func New(out io.Writer) *Interpreter

// Eval runs each complete top-level declaration in src, in order, in the session.
// It stops at the first error; declarations before the error stay bound.
// Trailing incomplete input is an error.
func (in *Interpreter) Eval(src string) ([]Result, error)

// SplitInput returns the complete ";"-terminated chunks of src and the incomplete
// remainder. The REPL and the golden runner both use it.
func SplitInput(src string) (chunks []string, rest string, err error)

// FormatError renders an error the way the REPL prints it:
// "Error: ..." or "uncaught exception ...".
func FormatError(err error) string
```

### Golden tests

Create the eight files below verbatim. `tests/golden_test.go` must do this for each `tests/golden/mN_*.mole`:

1. Skip the file when N is greater than `interp.Milestone`.
2. Create a fresh `interp.New(&buf)`.
3. Split the file with `SplitInput`. For each chunk, call `Eval`, then append each result's `Line(withTypes)` and a newline to `buf`, or `FormatError(err)` and a newline if it failed. `withTypes` is false when the file's first line is `(* golden: untyped *)`, and true otherwise. `print` output reaches `buf` on its own, in order.
4. Compare `buf` with the matching `.out` file line by line. An expected line that is exactly `Error:` matches any actual line starting with `Error:`. Every other line must match exactly. On a mismatch, report the file, line number, expected line, and actual line.

#### M1

`tests/golden/m1_core.mole`

```
(* golden: untyped *)
1 + 2 * 3;
val x = 10;
x div 3;
x mod 3;
~7 div 2;
~7 mod 2;
~5 + 2;
"mo" ^ "le";
size "mole";
if x > 5 then "big" else "small";
val add = fn a => fn b => a + b;
add 2 3;
fun fact n = if n = 0 then 1 else n * fact (n - 1);
fact 10;
let val y = 5 val z = 6 in y * z end;
(1, "two", true);
val (a, b) = (1, "x");
[1, 2, 3];
0 :: [1, 2] @ [3];
true andalso (false orelse true);
();
6 * 7;
it + 1;
(* a comment (* nested *) still a comment *) 5;
print "hello\n";
(print "a"; print "b\n"; 5);
foldl (fn (x, acc) => x + acc) 0 [1, 2, 3];
map (fn n => n * n) [1, 2, 3];
rev [1, 2, 3];
Int.toString ~42;
List.filter (fn n => n > 1) [1, 2, 3];
fun loop n acc = if n = 0 then acc else loop (n - 1) (acc + 1);
loop 1000000 0;
1 div 0;
hd [];
undefinedName;
x;
```

`tests/golden/m1_core.out`

```
val it = 7
val x = 10
val it = 3
val it = 1
val it = ~4
val it = 1
val it = ~3
val it = "mole"
val it = 4
val it = "big"
val add = fn
val it = 5
val fact = fn
val it = 3628800
val it = 30
val it = (1,"two",true)
val a = 1
val b = "x"
val it = [1,2,3]
val it = [0,1,2,3]
val it = true
val it = ()
val it = 42
val it = 43
val it = 5
hello
val it = ()
ab
val it = 5
val it = 6
val it = [1,4,9]
val it = [3,2,1]
val it = "~42"
val it = [2,3]
val loop = fn
val it = 1000000
uncaught exception Div
uncaught exception Empty
Error:
val it = 10
```

#### M2

`tests/golden/m2_types.mole`

```
42;
val s = "mole";
fun double x = x * 2;
fun id x = x;
id 3;
id "three";
(id 1, id true);
fun compose f g x = f (g x);
val quad = double o double;
quad 3;
fun len xs = if null xs then 0 else 1 + len (tl xs);
len [true, false, true];
map double [1, 2, 3];
[];
[[1, 2], [3]];
fn (a, b) => (b, a);
let val f = fn x => x in (f 1, f "a") end;
fun eq (a, b) = a = b;
print "hi\n";
foldr (fn (x, acc) => x :: acc) [] [1, 2, 3];
List.filter;
if true then 1 else "no";
1 + true;
fun bad x = x x;
fn f => (f 1, f "a");
val y = 1 + true;
y;
double 21;
```

`tests/golden/m2_types.out`

```
val it = 42 : int
val s = "mole" : string
val double = fn : int -> int
val id = fn : 'a -> 'a
val it = 3 : int
val it = "three" : string
val it = (1,true) : int * bool
val compose = fn : ('a -> 'b) -> ('c -> 'a) -> 'c -> 'b
val quad = fn : int -> int
val it = 12 : int
val len = fn : 'a list -> int
val it = 3 : int
val it = [2,4,6] : int list
val it = [] : 'a list
val it = [[1,2],[3]] : int list list
val it = fn : 'a * 'b -> 'b * 'a
val it = (1,"a") : int * string
val eq = fn : 'a * 'a -> bool
hi
val it = () : unit
val it = [1,2,3] : int list
val it = fn : ('a -> bool) -> 'a list -> 'a list
Error:
Error:
Error:
Error:
Error:
Error:
val it = 42 : int
```

#### M3

`tests/golden/m3_patterns.mole`

```
datatype shape = Circle of int | Rect of int * int;
fun area (Circle r) = 3 * r * r
  | area (Rect (w, h)) = w * h;
area (Rect (3, 5));
area (Circle 2);
Circle 2;
[Circle 1, Rect (2, 3)];
datatype 'a tree = Leaf | Node of 'a tree * 'a * 'a tree;
fun insert (x, Leaf) = Node (Leaf, x, Leaf)
  | insert (x, Node (l, v, r)) =
      if x < v then Node (insert (x, l), v, r)
      else Node (l, v, insert (x, r));
fun toList Leaf = []
  | toList (Node (l, v, r)) = toList l @ [v] @ toList r;
toList (foldl insert Leaf [5, 2, 8, 1]);
insert (1, Leaf);
fun sum [] = 0
  | sum (x :: xs) = x + sum xs;
sum [1, 2, 3, 4];
fun fib 0 = 0 | fib 1 = 1 | fib n = fib (n - 1) + fib (n - 2);
fib 20;
case SOME 3 of NONE => 0 | SOME n => n + 1;
fun first [] = NONE | first (x :: _) = SOME x;
first ["a", "b"];
SOME (SOME ~1);
val [p, q] = [1, 2];
val h :: t = [7, 8, 9];
fun swap (a, b) = (b, a);
(fn [] => "empty" | _ => "nonempty") [1];
fun zero 0 = "zero";
zero 1;
val [r] = [1, 2];
Circle "big";
case 1 of true => 0 | false => 1;
fun dup (x, x) = x;
Nope 3;
```

`tests/golden/m3_patterns.out`

```
datatype shape = Circle of int | Rect of int * int
val area = fn : shape -> int
val it = 15 : int
val it = 12 : int
val it = Circle 2 : shape
val it = [Circle 1,Rect (2,3)] : shape list
datatype 'a tree = Leaf | Node of 'a tree * 'a * 'a tree
val insert = fn : int * int tree -> int tree
val toList = fn : 'a tree -> 'a list
val it = [1,2,5,8] : int list
val it = Node (Leaf,1,Leaf) : int tree
val sum = fn : int list -> int
val it = 10 : int
val fib = fn : int -> int
val it = 6765 : int
val it = 4 : int
val first = fn : 'a list -> 'a option
val it = SOME "a" : string option
val it = SOME (SOME ~1) : int option option
val p = 1 : int
val q = 2 : int
val h = 7 : int
val t = [8,9] : int list
val swap = fn : 'a * 'b -> 'b * 'a
val it = "nonempty" : string
val zero = fn : int -> string
uncaught exception Match
uncaught exception Bind
Error:
Error:
Error:
Error:
```

#### M4

`tests/golden/m4_exceptions.mole`

```
exception Oops;
exception Bad of string;
fun safeDiv (_, 0) = raise Div
  | safeDiv (a, b) = a div b;
safeDiv (10, 0) handle Div => ~1;
safeDiv (10, 3);
(raise Bad "x") handle Bad msg => msg ^ "!";
hd [] handle Empty => 0;
Bad "y";
fun check n = if n < 0 then raise Fail "negative" else n;
check 5;
check ~5 handle Fail m => size m;
raise Bad "boom";
(raise Oops) handle Bad _ => 1;
raise Fail "nope";
fun firstPos [] = raise Oops
  | firstPos (x :: xs) = if x > 0 then x else firstPos xs;
firstPos [~1, 0, 3] handle Oops => 0;
firstPos [~1] handle Oops => 0;
raise 5;
1 handle Oops => "one";
```

`tests/golden/m4_exceptions.out`

```
exception Oops
exception Bad of string
val safeDiv = fn : int * int -> int
val it = ~1 : int
val it = 3 : int
val it = "x!" : string
val it = 0 : int
val it = Bad "y" : exn
val check = fn : int -> int
val it = 5 : int
val it = 8 : int
uncaught exception Bad "boom"
uncaught exception Oops
uncaught exception Fail "nope"
val firstPos = fn : int list -> int
val it = 3 : int
val it = 0 : int
Error:
Error:
```

### Milestone report

After each milestone, reply with only this report. No preamble.

```
## Milestone N: <name>
Status: COMPLETE | INCOMPLETE
make vet: PASS | FAIL
make test: PASS | FAIL
<last 15 lines of make test output>
Goldens: <passing>/<run>; failing: <file:line, one per line, or none>
Built: <3 to 6 bullets on what now works>
Decisions: <new docs/DECISIONS.md entries, or none>
Known gaps: <bullets, or none>
```

### Stretch list (only when asked)

Type annotations `(e : ty)`, `and` for mutual recursion, `op+`, `as` patterns, `local`, records, `ref` with `!` and `:=` (which brings in the value restriction), chars with `explode` and `implode`, reals, exhaustiveness warnings, a WASM build like shrew's.
