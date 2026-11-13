package ast

import "testing"

func TestShow(t *testing.T) {
	tests := []struct {
		name string
		n    any
		want string
	}{
		{
			name: "IntLit positive",
			n:    &IntLit{Node: Node{}, Value: 5},
			want: "5",
		},
		{
			name: "IntLit negative",
			n:    &IntLit{Node: Node{}, Value: -5},
			want: "-5",
		},
		{
			name: "StrLit",
			n:    &StrLit{Node: Node{}, Value: `hi`},
			want: `"hi"`,
		},
		{
			name: "Var",
			n:    &Var{Node: Node{}, Name: "x"},
			want: "x",
		},
		{
			name: "Var plus",
			n:    &Var{Node: Node{}, Name: "+"},
			want: "+",
		},
		{
			name: "Con true",
			n:    &Con{Node: Node{}, Name: "true"},
			want: "true",
		},
		{
			name: "Con nil",
			n:    &Con{Node: Node{}, Name: "nil"},
			want: "nil",
		},
		{
			name: "Tuple empty",
			n:    &Tuple{Node: Node{}, Elems: []Expr{}},
			want: "()",
		},
		{
			name: "Tuple one element",
			n:    &Tuple{Node: Node{}, Elems: []Expr{&IntLit{Node: Node{}, Value: 1}}},
			want: "(tuple 1)",
		},
		{
			name: "Tuple two elements",
			n:    &Tuple{Node: Node{}, Elems: []Expr{&IntLit{Node: Node{}, Value: 1}, &IntLit{Node: Node{}, Value: 2}}},
			want: "(tuple 1 2)",
		},
		{
			name: "List empty",
			n:    &List{Node: Node{}, Elems: []Expr{}},
			want: "(list)",
		},
		{
			name: "List one element",
			n:    &List{Node: Node{}, Elems: []Expr{&IntLit{Node: Node{}, Value: 1}}},
			want: "(list 1)",
		},
		{
			name: "List two elements",
			n:    &List{Node: Node{}, Elems: []Expr{&IntLit{Node: Node{}, Value: 1}, &IntLit{Node: Node{}, Value: 2}}},
			want: "(list 1 2)",
		},
		{
			name: "App",
			n:    &App{Fn: &Var{Node: Node{}, Name: "+"}, Arg: &Tuple{Elems: []Expr{&IntLit{Value: 1}, &IntLit{Value: 2}}}},
			want: "(app + (tuple 1 2))",
		},
		{
			name: "If",
			n:    &If{Cond: &Var{Name: "x"}, Then: &IntLit{Value: 1}, Else: &IntLit{Value: 2}},
			want: "(if x 1 2)",
		},
		{
			name: "Fn",
			n:    &Fn{Arms: []Arm{{Pat: &PVar{Name: "x"}, Body: &Var{Name: "x"}}}},
			want: "(fn (x x))",
		},
		{
			name: "Let",
			n:    &Let{Decs: []Dec{&ValDec{Pat: &PVar{Name: "x"}, Expr: &IntLit{Value: 1}}}, Body: []Expr{&IntLit{Value: 2}}},
			want: "(let (x 1) 2)",
		},
		{
			name: "Seq",
			n:    &Seq{Exprs: []Expr{&IntLit{Value: 1}, &IntLit{Value: 2}}},
			want: "(seq 1 2)",
		},
		{
			name: "AndAlso",
			n:    &AndAlso{L: &Var{Name: "a"}, R: &Var{Name: "b"}},
			want: "(andalso a b)",
		},
		{
			name: "OrElse",
			n:    &OrElse{L: &Var{Name: "a"}, R: &Var{Name: "b"}},
			want: "(orelse a b)",
		},
		{
			name: "PWild",
			n:    &PWild{Node: Node{}},
			want: "_",
		},
		{
			name: "PVar",
			n:    &PVar{Name: "x"},
			want: "x",
		},
		{
			name: "PTuple empty",
			n:    &PTuple{Node: Node{}, Elems: []Pat{}},
			want: "()",
		},
		{
			name: "PTuple one element",
			n:    &PTuple{Node: Node{}, Elems: []Pat{&PVar{Name: "x"}}},
			want: "(tuple x)",
		},
		{
			name: "PTuple two elements",
			n:    &PTuple{Node: Node{}, Elems: []Pat{&PVar{Name: "a"}, &PWild{Node: Node{}}}},
			want: "(tuple a _)",
		},
		{
			name: "ValDec",
			n:    &ValDec{Pat: &PTuple{Elems: []Pat{&PVar{Name: "a"}, &PWild{Node: Node{}}}}, Expr: &IntLit{Value: 1}},
			want: "(val (tuple a _) 1)",
		},
		{
			name: "FunDec",
			n:    &FunDec{Name: "f", Clauses: []Clause{{Pats: []Pat{&PVar{Name: "x"}}, Body: &Var{Name: "x"}}}},
			want: "(fun f ((x x)))",
		},
		{
			name: "FunDec multiple clauses",
			n:    &FunDec{Name: "f", Clauses: []Clause{
				{Pats: []Pat{&PVar{Name: "x"}}, Body: &IntLit{Value: 1}},
				{Pats: []Pat{&PWild{Node: Node{}}}, Body: &IntLit{Value: 2}},
			}},
			want: "(fun f ((x 1) (_ 2)))",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Show(tt.n)
			if got != tt.want {
				t.Errorf("Show() = %q, want %q", got, tt.want)
			}
		})
	}
}
