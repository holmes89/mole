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
