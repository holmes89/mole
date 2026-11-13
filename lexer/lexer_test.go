package lexer

import (
	"testing"

	"github.com/holmes89/mole/token"
)

func TestBasic(t *testing.T) {
	tests := []struct {
		name   string
		src    string
		expect []token.Token
		err    bool
	}{
		{
			name: "semicolon statement",
			src:  "val x = 10;",
			expect: []token.Token{
				{Kind: token.KEYWORD, Text: "val", Pos: posAt("val x = 10;", 0)},
				{Kind: token.IDENT, Text: "x", Pos: posAt("val x = 10;", 4)},
				{Kind: token.SYMBOL, Text: "=", Pos: posAt("val x = 10;", 6)},
				{Kind: token.INT, Text: "10", Pos: posAt("val x = 10;", 8)},
				{Kind: token.SYMBOL, Text: ";", Pos: posAt("val x = 10;", 10)},
				{Kind: token.EOF, Text: "", Pos: posAt("val x = 10;", 11)},
			},
			err: false,
		},
		{
			name: "less equal",
			src:  "a<=b",
			expect: []token.Token{
				{Kind: token.IDENT, Text: "a", Pos: posAt("a<=b", 0)},
				{Kind: token.SYMBOL, Text: "<=", Pos: posAt("a<=b", 1)},
				{Kind: token.IDENT, Text: "b", Pos: posAt("a<=b", 3)},
				{Kind: token.EOF, Text: "", Pos: posAt("a<=b", 4)},
			},
			err: false,
		},
		{
			name: "cons operator",
			src:  "x :: xs",
			expect: []token.Token{
				{Kind: token.IDENT, Text: "x", Pos: posAt("x :: xs", 0)},
				{Kind: token.SYMBOL, Text: "::", Pos: posAt("x :: xs", 2)},
				{Kind: token.IDENT, Text: "xs", Pos: posAt("x :: xs", 5)},
				{Kind: token.EOF, Text: "", Pos: posAt("x :: xs", 7)},
			},
			err: false,
		},
		{
			name: "lambda expression",
			src:  "fn a => a",
			expect: []token.Token{
				{Kind: token.KEYWORD, Text: "fn", Pos: posAt("fn a => a", 0)},
				{Kind: token.IDENT, Text: "a", Pos: posAt("fn a => a", 3)},
				{Kind: token.SYMBOL, Text: "=>", Pos: posAt("fn a => a", 5)},
				{Kind: token.IDENT, Text: "a", Pos: posAt("fn a => a", 8)},
				{Kind: token.EOF, Text: "", Pos: posAt("fn a => a", 9)},
			},
			err: false,
		},
		{
			name: "constructor and identifier",
			src:  "Leaf true foo'",
			expect: []token.Token{
				{Kind: token.CON, Text: "Leaf", Pos: posAt("Leaf true foo'", 0)},
				{Kind: token.CON, Text: "true", Pos: posAt("Leaf true foo'", 5)},
				{Kind: token.IDENT, Text: "foo'", Pos: posAt("Leaf true foo'", 10)},
				{Kind: token.EOF, Text: "", Pos: posAt("Leaf true foo'", 14)},
			},
			err: false,
		},
		{
			name: "newline tracking",
			src:  "a\n  b",
			expect: []token.Token{
				{Kind: token.IDENT, Text: "a", Pos: token.Pos{Line: 1, Col: 1, Offset: 0}},
				{Kind: token.IDENT, Text: "b", Pos: token.Pos{Line: 2, Col: 3, Offset: 4}},
				{Kind: token.EOF, Text: "", Pos: token.Pos{Line: 2, Col: 4, Offset: 5}},
			},
			err: false,
		},
		{
			name: "unexpected character",
			src:  "#",
			expect: nil,
			err:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toks, err := Tokenize(tt.src)
			if (err != nil) != tt.err {
				t.Errorf("Tokenize(%q): error = %v, wantErr = %v", tt.src, err, tt.err)
				return
			}
			if tt.err {
				return
			}
			if len(toks) != len(tt.expect) {
				t.Errorf("len(tokens) = %d, want %d", len(toks), len(tt.expect))
				for i, tok := range toks {
					t.Logf("token[%d] = %+v", i, tok)
				}
				return
			}
			for i, tok := range toks {
				want := tt.expect[i]
				if tok.Kind != want.Kind {
					t.Errorf("token[%d].Kind = %v, want %v", i, tok.Kind, want.Kind)
				}
				if tok.Text != want.Text {
					t.Errorf("token[%d].Text = %q, want %q", i, tok.Text, want.Text)
				}
				if tok.Pos.Line != want.Pos.Line || tok.Pos.Col != want.Pos.Col || tok.Pos.Offset != want.Pos.Offset {
					t.Errorf("token[%d].Pos = %+v, want %+v", i, tok.Pos, want.Pos)
				}
			}
		})
	}
}

func TestRest(t *testing.T) {
	tests := []struct {
		name   string
		src    string
		expect []token.Token
		err    bool
		incomplete bool
	}{
		{
			name: "string with escapes",
			src:  `"mo\"le\n"`,
			expect: []token.Token{
				{Kind: token.STRING, Text: `mo"le
`, Pos: posAt(`"mo\"le\n"`, 1)},
				{Kind: token.EOF, Text: "", Pos: posAt(`"mo\"le\n"`, 10)},
			},
			err: false,
		},
		{
			name: "comment",
			src:  `(* a (* b *) c *) 5`,
			expect: []token.Token{
				{Kind: token.INT, Text: "5", Pos: posAt(`(* a (* b *) c *) 5`, 18)},
				{Kind: token.EOF, Text: "", Pos: posAt(`(* a (* b *) c *) 5`, 19)},
			},
			err: false,
		},
		{
			name: "unclosed comment",
			src:  `(* open`,
			expect: nil,
			err:  true,
		},
		{
			name: "unclosed string",
			src:  `"abc`,
			expect: nil,
			err:  true,
		},
		{
			name: "tilde prefix int",
			src:  `~5 + ~ x`,
			expect: []token.Token{
				{Kind: token.INT, Text: "~5", Pos: posAt(`~5 + ~ x`, 0)},
				{Kind: token.SYMBOL, Text: "+", Pos: posAt(`~5 + ~ x`, 3)},
				{Kind: token.IDENT, Text: "~", Pos: posAt(`~5 + ~ x`, 5)},
				{Kind: token.IDENT, Text: "x", Pos: posAt(`~5 + ~ x`, 7)},
				{Kind: token.EOF, Text: "", Pos: posAt(`~5 + ~ x`, 8)},
			},
			err: false,
		},
		{
			name: "long id",
			src:  `Int.toString ~42`,
			expect: []token.Token{
				{Kind: token.LONGID, Text: "Int.toString", Pos: posAt(`Int.toString ~42`, 0)},
				{Kind: token.INT, Text: "~42", Pos: posAt(`Int.toString ~42`, 13)},
				{Kind: token.EOF, Text: "", Pos: posAt(`Int.toString ~42`, 16)},
			},
			err: false,
		},
		{
			name: "type variable",
			src:  `'a list`,
			expect: []token.Token{
				{Kind: token.TYVAR, Text: "'a", Pos: posAt(`'a list`, 0)},
				{Kind: token.IDENT, Text: "list", Pos: posAt(`'a list`, 3)},
				{Kind: token.EOF, Text: "", Pos: posAt(`'a list`, 7)},
			},
			err: false,
		},
		{
			name: "parens around int",
			src:  `(1)`,
			expect: []token.Token{
				{Kind: token.SYMBOL, Text: "(", Pos: posAt(`(1)`, 0)},
				{Kind: token.INT, Text: "1", Pos: posAt(`(1)`, 1)},
				{Kind: token.SYMBOL, Text: ")", Pos: posAt(`(1)`, 2)},
				{Kind: token.EOF, Text: "", Pos: posAt(`(1)`, 3)},
			},
			err: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			toks, err := Tokenize(tt.src)
			if (err != nil) != tt.err {
				t.Errorf("Tokenize(%q): error = %v, wantErr = %v", tt.src, err, tt.err)
				return
			}
			if tt.err {
				if !tt.incomplete {
					t.Logf("error: %v", err)
				}
				return
			}
			if len(toks) != len(tt.expect) {
				t.Errorf("len(tokens) = %d, want %d", len(toks), len(tt.expect))
				for i, tok := range toks {
					t.Logf("token[%d] = %+v", i, tok)
				}
				return
			}
			for i, tok := range toks {
				want := tt.expect[i]
				if tok.Kind != want.Kind {
					t.Errorf("token[%d].Kind = %v, want %v", i, tok.Kind, want.Kind)
				}
				if tok.Text != want.Text {
					t.Errorf("token[%d].Text = %q, want %q", i, tok.Text, want.Text)
				}
				if tok.Pos.Line != want.Pos.Line || tok.Pos.Col != want.Pos.Col || tok.Pos.Offset != want.Pos.Offset {
					t.Errorf("token[%d].Pos = %+v, want %+v", i, tok.Pos, want.Pos)
				}
			}
		})
	}
}
