package parser

import (
	"strings"
	"testing"

	"github.com/holmes89/mole/ast"
)

func TestAtoms(t *testing.T) {
	tests := []struct {
		src  string
		want string
	}{
		{"42", "42"},
		{"~5", "-5"},
		{"f x y", "(app (app f x) y)"},
		{"()", "()"},
		{"(1, \"two\", true)", "(tuple 1 \"two\" true)"},
		{`[1, 2]`, "(list 1 2)"},
		{"[]", "(list)"},
		{"(f 1; g 2)", "(seq (app f 1) (app g 2))"},
		{"Int.toString 3", "(app Int.toString 3)"},
	}

	for _, tc := range tests {
		t.Run(tc.src, func(t *testing.T) {
			expr, err := ParseExpr(tc.src)
			if err != nil && strings.HasPrefix(tc.want, "<error>") {
				return
			}
			if err != nil {
				t.Fatalf("ParseExpr(%q) error: %v", tc.src, err)
				return
			}
			got := ast.Show(expr)
			if got != tc.want {
				t.Errorf("Show() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestErrors(t *testing.T) {
	tests := []string{
		"(1",
		"1 )",
	}

	for _, src := range tests {
		t.Run(src, func(t *testing.T) {
			_, err := ParseExpr(src)
			if err == nil {
				t.Errorf("ParseExpr(%q) should return error", src)
			}
		})
	}
}

func TestInfix(t *testing.T) {
	tests := []struct {
		src  string
		want string
	}{
		{"1 + 2 * 3", "(app + (tuple 1 (app * (tuple 2 3))))"},
		{"1 - 2 - 3", "(app - (tuple (app - (tuple 1 2)) 3))"},
		{"0 :: [1] @ [2]", "(app :: (tuple 0 (app @ (tuple (list 1) (list 2)))))"},
		{"f x + 1", "(app + (tuple (app f x) 1))"},
		{"a andalso b orelse c", "(orelse (andalso a b) c)"},
		{"x div 3 = 1", "(app = (tuple (app div (tuple x 3)) 1))"},
		{"f o g", "(app o (tuple f g))"},
	}

	for _, tc := range tests {
		t.Run(tc.src, func(t *testing.T) {
			expr, err := ParseExpr(tc.src)
			if err != nil {
				t.Fatalf("ParseExpr(%q) error: %v", tc.src, err)
				return
			}
			got := ast.Show(expr)
			if got != tc.want {
				t.Errorf("Show() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestForms(t *testing.T) {
	tests := []struct {
		src  string
		want string
	}{
		{"if x > 5 then 1 else 2", "(if (app > (tuple x 5)) 1 2)"},
		{"fn a => fn b => a + b", "(fn (a (fn (b (app + (tuple a b))))))"},
		{"fn (a, _) => a", "(fn ((tuple a _) a))"},
		{"fn () => 1", "(fn (() 1))"},
		{"f (fn x => x) 1", "(app (app f (fn (x x))) 1)"},
	}

	for _, tc := range tests {
		t.Run(tc.src, func(t *testing.T) {
			expr, err := ParseExpr(tc.src)
			if err != nil && strings.HasPrefix(tc.want, "<error>") {
				return
			}
			if err != nil {
				t.Fatalf("ParseExpr(%q) error: %v", tc.src, err)
				return
			}
			got := ast.Show(expr)
			if got != tc.want {
				t.Errorf("Show() = %q, want %q", got, tc.want)
			}
		})
	}

	t.Run("case not supported", func(t *testing.T) {
		_, err := ParseExpr("case x of y => y")
		if err == nil {
			t.Errorf("ParseExpr(%q) should return error", "case x of y => y")
			return
		}
		if !strings.Contains(err.Error(), "not supported") {
			t.Errorf("error message = %q, want to contain 'not supported'", err.Error())
		}
	})
}
