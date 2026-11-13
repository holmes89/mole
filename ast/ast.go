package ast

import "github.com/holmes89/mole/token"

type Expr interface{ Pos() token.Pos }
type Pat interface{ Pos() token.Pos }
type Dec interface{ Pos() token.Pos }

type Node struct{ P token.Pos }

func (n Node) Pos() token.Pos { return n.P }

// Expressions
type IntLit struct { Node; Value int64 }
type StrLit struct { Node; Value string }
type Var struct { Node; Name string }      // x, +, Int.toString
type Con struct { Node; Name string }      // Leaf, true, nil
type Tuple struct { Node; Elems []Expr }   // zero elements is unit
type List struct { Node; Elems []Expr }
type App struct { Node; Fn, Arg Expr }
type If struct { Node; Cond, Then, Else Expr }
type Fn struct { Node; Arms []Arm }
type Let struct { Node; Decs []Dec; Body []Expr }
type Seq struct { Node; Exprs []Expr }
type AndAlso struct { Node; L, R Expr }
type OrElse struct { Node; L, R Expr }

type Arm struct { Pat Pat; Body Expr }

// Patterns
type PWild struct{ Node }
type PVar struct { Node; Name string }
type PTuple struct { Node; Elems []Pat }   // zero elements is unit

// Declarations
type ValDec struct { Node; Pat Pat; Expr Expr }
type FunDec struct { Node; Name string; Clauses []Clause }
type Clause struct { Pats []Pat; Body Expr }
