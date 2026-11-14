package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/holmes89/mole/ast"
	"github.com/holmes89/mole/lexer"
	"github.com/holmes89/mole/token"
)

func ParseExpr(src string) (ast.Expr, error) {
	toks, err := lexer.Tokenize(src)
	if err != nil {
		return nil, err
	}
	p := &parser{toks: toks}
	expr, _ := p.parseApp()
	if p.peek().Kind != token.EOF && p.peek().Text != ")" && p.peek().Text != "]" {
		return nil, Error{Pos: p.peek().Pos, Msg: "unexpected token after expression"}
	}
	return expr, nil
}

func (p *parser) parseExpr() (ast.Expr, bool) {
	expr, ok := p.parseApp()
	if !ok {
		return nil, false
	}
	return expr, true
}

func (p *parser) parseApp() (ast.Expr, bool) {
	t := p.peek()

	// Handle empty tuple () and list []
	if t.Text == "(" {
		p.next()
		fmt.Printf("DEBUG tuple: consumed (, next=%s\n", p.peek().Text)

		var elems []ast.Expr

		for p.peek().Kind != token.EOF && p.peek().Text != ")" {
			t2 := p.peek()

			switch t2.Kind {
			case token.INT:
				p.next()
				text := t2.Text
				var v int64
				var err error
				if strings.HasPrefix(text, "~") {
					v, err = strconv.ParseInt(text[1:], 10, 64)
					if err != nil {
						return nil, false
					}
					v = -v
				} else {
					v, err = strconv.ParseInt(text, 10, 64)
					if err != nil {
						return nil, false
					}
				}
				elems = append(elems, &ast.IntLit{Node: ast.Node{P: t2.Pos}, Value: v})

			case token.STRING:
				p.next()
				elems = append(elems, &ast.StrLit{Node: ast.Node{P: t2.Pos}, Value: t2.Text})

			case token.IDENT, token.LONGID, token.TYVAR:
				p.next()
				elems = append(elems, &ast.Var{Node: ast.Node{P: t2.Pos}, Name: t2.Text})

			case token.CON:
				p.next()
				elems = append(elems, &ast.Con{Node: ast.Node{P: t2.Pos}, Name: t2.Text})

			default:
				return nil, false
			}

			t = p.peek()
			if t.Text == "," {
				p.next() // consume comma
				continue
			}
			if t.Text != ")" {
				return nil, false
			}
			p.next() // consume )
		}

		if len(elems) == 0 {
			fmt.Printf("DEBUG tuple: empty tuple, returning Tuple node\n")
			return &ast.Tuple{Node: ast.Node{P: t.Pos}, Elems: []ast.Expr{}}, true
		}

		fmt.Printf("DEBUG tuple: returning Tuple with %d elems\n", len(elems))
		return &ast.Tuple{Node: ast.Node{P: t.Pos}, Elems: elems}, true

	} else if t.Text == "[" {
		p.next()
		if p.peek().Kind == token.EOF || p.peek().Text == "]" {
			if p.peek().Text == "]" {
				p.next()
			}
			return &ast.List{Node: ast.Node{P: t.Pos}, Elems: []ast.Expr{}}, true
		}

		var elems []ast.Expr
		for p.peek().Text != "]" {
			t2 := p.peek()

			switch t2.Kind {
			case token.INT:
				p.next()
				text := t2.Text
				var v int64
				var err error
				if strings.HasPrefix(text, "~") {
					v, err = strconv.ParseInt(text[1:], 10, 64)
					if err != nil {
						return nil, false
					}
					v = -v
				} else {
					v, err = strconv.ParseInt(text, 10, 64)
					if err != nil {
						return nil, false
					}
				}
				elems = append(elems, &ast.IntLit{Node: ast.Node{P: t2.Pos}, Value: v})

			case token.STRING:
				p.next()
				elems = append(elems, &ast.StrLit{Node: ast.Node{P: t2.Pos}, Value: t2.Text})

			case token.IDENT, token.LONGID, token.TYVAR:
				p.next()
				elems = append(elems, &ast.Var{Node: ast.Node{P: t2.Pos}, Name: t2.Text})

			case token.CON:
				p.next()
				elems = append(elems, &ast.Con{Node: ast.Node{P: t2.Pos}, Name: t2.Text})

			default:
				return nil, false
			}

			t = p.peek()
			if t.Text == "," {
				p.next() // consume comma
				continue
			}
			if t.Text != "]" {
				return nil, false
			}
			p.next() // consume ]
		}

		if len(elems) == 0 {
			return &ast.List{Node: ast.Node{P: t.Pos}, Elems: []ast.Expr{}}, true
		}
		// Consume the closing bracket
		t = p.peek()
		if t.Text == "]" {
			fmt.Printf("DEBUG list: consuming ]\n")
			p.next()
		}
		return &ast.List{Node: ast.Node{P: t.Pos}, Elems: elems}, true

	} else if t.Text == ";" {
		p.next()

		var expr ast.Expr
		t2 := p.peek()

		switch t2.Kind {
		case token.INT:
			p.next()
			text := t2.Text
			var v int64
			var err error
			if strings.HasPrefix(text, "~") {
				v, err = strconv.ParseInt(text[1:], 10, 64)
				if err != nil {
					return nil, false
				}
				v = -v
			} else {
				v, err = strconv.ParseInt(text, 10, 64)
				if err != nil {
					return nil, false
				}
			}
			expr = &ast.IntLit{Node: ast.Node{P: t2.Pos}, Value: v}

		case token.STRING:
			p.next()
			expr = &ast.StrLit{Node: ast.Node{P: t2.Pos}, Value: t2.Text}

		case token.IDENT, token.LONGID, token.TYVAR:
			p.next()
			expr = &ast.Var{Node: ast.Node{P: t2.Pos}, Name: t2.Text}

		case token.CON:
			p.next()
			expr = &ast.Con{Node: ast.Node{P: t2.Pos}, Name: t2.Text}

		default:
			return nil, false
		}

		for p.peek().Text != ";" && p.peek().Kind != token.EOF {
			t3 := p.peek()

			switch t3.Kind {
			case token.INT:
				p.next()
				text := t3.Text
				var v int64
				var err error
				if strings.HasPrefix(text, "~") {
					v, err = strconv.ParseInt(text[1:], 10, 64)
					if err != nil {
						return nil, false
					}
					v = -v
				} else {
					v, err = strconv.ParseInt(text, 10, 64)
					if err != nil {
						return nil, false
					}
				}
				expr = &ast.IntLit{Node: ast.Node{P: t3.Pos}, Value: v}

			case token.STRING:
				p.next()
				expr = &ast.StrLit{Node: ast.Node{P: t3.Pos}, Value: t3.Text}

			case token.IDENT, token.LONGID, token.TYVAR:
				p.next()
				expr = &ast.Var{Node: ast.Node{P: t3.Pos}, Name: t3.Text}

			case token.CON:
				p.next()
				expr = &ast.Con{Node: ast.Node{P: t3.Pos}, Name: t3.Text}

			default:
				return nil, false
			}

			t = p.peek()
		}

		return expr, true

	} else if t.Text == "," {
		p.next()

		var expr ast.Expr
		t2 := p.peek()

		switch t2.Kind {
		case token.INT:
			p.next()
			text := t2.Text
			var v int64
			var err error
			if strings.HasPrefix(text, "~") {
				v, err = strconv.ParseInt(text[1:], 10, 64)
				if err != nil {
					return nil, false
				}
				v = -v
			} else {
				v, err = strconv.ParseInt(text, 10, 64)
				if err != nil {
					return nil, false
				}
			}
			expr = &ast.IntLit{Node: ast.Node{P: t2.Pos}, Value: v}

		case token.STRING:
			p.next()
			expr = &ast.StrLit{Node: ast.Node{P: t2.Pos}, Value: t2.Text}

		case token.IDENT, token.LONGID, token.TYVAR:
			p.next()
			expr = &ast.Var{Node: ast.Node{P: t2.Pos}, Name: t2.Text}

		case token.CON:
			p.next()
			expr = &ast.Con{Node: ast.Node{P: t2.Pos}, Name: t2.Text}

		default:
			return nil, false
		}

		for p.peek().Text != "," && p.peek().Kind != token.EOF {
			t3 := p.peek()

			switch t3.Kind {
			case token.INT:
				p.next()
				text := t3.Text
				var v int64
				var err error
				if strings.HasPrefix(text, "~") {
					v, err = strconv.ParseInt(text[1:], 10, 64)
					if err != nil {
						return nil, false
					}
					v = -v
				} else {
					v, err = strconv.ParseInt(text, 10, 64)
					if err != nil {
						return nil, false
					}
				}
				expr = &ast.IntLit{Node: ast.Node{P: t3.Pos}, Value: v}

			case token.STRING:
				p.next()
				expr = &ast.StrLit{Node: ast.Node{P: t3.Pos}, Value: t3.Text}

			case token.IDENT, token.LONGID, token.TYVAR:
				p.next()
				expr = &ast.Var{Node: ast.Node{P: t3.Pos}, Name: t3.Text}

			case token.CON:
				p.next()
				expr = &ast.Con{Node: ast.Node{P: t3.Pos}, Name: t3.Text}

			default:
				return nil, false
			}

			t = p.peek()
		}

		return expr, true

	} else if t.Text == "|" {
		p.next()
		left, _ := p.parseExpr()
		if left == nil {
			return nil, false
		}

		t = p.peek()
		for t.Kind != token.EOF && t.Text != "|" {
			switch t.Kind {
			case token.INT:
				p.next()
				text := t.Text
				var v int64
				var err error
				if strings.HasPrefix(text, "~") {
					v, err = strconv.ParseInt(text[1:], 10, 64)
					if err != nil {
						return nil, false
					}
					v = -v
				} else {
					v, err = strconv.ParseInt(text, 10, 64)
					if err != nil {
						return nil, false
					}
				}
				right := &ast.IntLit{Node: ast.Node{P: t.Pos}, Value: v}

			case token.STRING:
				p.next()
				right := &ast.StrLit{Node: ast.Node{P: t.Pos}, Value: t.Text}

			case token.IDENT, token.LONGID, token.TYVAR:
				p.next()
				right := &ast.Var{Node: ast.Node{P: t.Pos}, Name: t.Text}

			case token.CON:
				p.next()
				right := &ast.Con{Node: ast.Node{P: t.Pos}, Name: t.Text}

			default:
				return nil, false
			}

			t = p.peek()
		}

		if right == nil {
			return nil, false
		}
		return &ast.OrElse{Node: ast.Node{P: t.Pos}, L: left, R: right}, true

	} else if t.Text == "=>" {
		p.next()
		left, _ := p.parseExpr()
		if left == nil {
			return nil, false
		}

		t = p.peek()
		for t.Kind != token.EOF && t.Text != "=>" {
			switch t.Kind {
			case token.INT:
				p.next()
				text := t.Text
				var v int64
				var err error
				if strings.HasPrefix(text, "~") {
					v, err = strconv.ParseInt(text[1:], 10, 64)
					if err != nil {
						return nil, false
					}
					v = -v
				} else {
					v, err = strconv.ParseInt(text, 10, 64)
					if err != nil {
						return nil, false
					}
				}
				right := &ast.IntLit{Node: ast.Node{P: t.Pos}, Value: v}

			case token.STRING:
				p.next()
				right := &ast.StrLit{Node: ast.Node{P: t.Pos}, Value: t.Text}

			case token.IDENT, token.LONGID, token.TYVAR:
				p.next()
				right := &ast.Var{Node: ast.Node{P: t.Pos}, Name: t.Text}

			case token.CON:
				p.next()
				right := &ast.Con{Node: ast.Node{P: t.Pos}, Name: t.Text}

			default:
				return nil, false
			}

			t = p.peek()
		}

		if right == nil {
			return nil, false
		}
		return &ast.AndAlso{Node: ast.Node{P: t.Pos}, L: left, R: right}, true

	} else if t.Text == "->" {
		p.next()
		left, _ := p.parseExpr()
		if left == nil {
			return nil, false
		}

		t = p.peek()
		for t.Kind != token.EOF && t.Text != "->" {
			switch t.Kind {
			case token.INT:
				p.next()
				text := t.Text
				var v int64
				var err error
				if strings.HasPrefix(text, "~") {
					v, err = strconv.ParseInt(text[1:], 10, 64)
					if err != nil {
						return nil, false
					}
					v = -v
				} else {
					v, err = strconv.ParseInt(text, 10, 64)
					if err != nil {
						return nil, false
					}
				}
				right := &ast.IntLit{Node: ast.Node{P: t.Pos}, Value: v}

			case token.STRING:
				p.next()
				right := &ast.StrLit{Node: ast.Node{P: t.Pos}, Value: t.Text}

			case token.IDENT, token.LONGID, token.TYVAR:
				p.next()
				right := &ast.Var{Node: ast.Node{P: t.Pos}, Name: t.Text}

			case token.CON:
				p.next()
				right := &ast.Con{Node: ast.Node{P: t.Pos}, Name: t.Text}

			default:
				return nil, false
			}

			t = p.peek()
		}

		if right == nil {
			return nil, false
		}
		return &ast.App{Node: ast.Node{P: t.Pos}, Fn: left, Arg: right}, true

	} else if t.Text == ":" || t.Text == "@" || t.Text == "<" || t.Text == ">" || t.Text == "<=" || t.Text == ">=" || t.Text == "+" || t.Text == "-" || t.Text == "*" || t.Text == "^" {
		p.next()
		left, _ := p.parseExpr()
		if left == nil {
			return nil, false
		}

		t = p.peek()
		for t.Kind != token.EOF && t.Text != ":" && t.Text != "@" && t.Text != "<" && t.Text != ">" && t.Text != "<=" && t.Text != ">=" && t.Text != "+" && t.Text != "-" && t.Text != "*" && t.Text != "^" {
			switch t.Kind {
			case token.INT:
				p.next()
				text := t.Text
				var v int64
				var err error
				if strings.HasPrefix(text, "~") {
					v, err = strconv.ParseInt(text[1:], 10, 64)
					if err != nil {
						return nil, false
					}
					v = -v
				} else {
					v, err = strconv.ParseInt(text, 10, 64)
					if err != nil {
						return nil, false
					}
				}
				right := &ast.IntLit{Node: ast.Node{P: t.Pos}, Value: v}

			case token.STRING:
				p.next()
				right := &ast.StrLit{Node: ast.Node{P: t.Pos}, Value: t.Text}

			case token.IDENT, token.LONGID, token.TYVAR:
				p.next()
				right := &ast.Var{Node: ast.Node{P: t.Pos}, Name: t.Text}

			case token.CON:
				p.next()
				right := &ast.Con{Node: ast.Node{P: t.Pos}, Name: t.Text}

			default:
				return nil, false
			}

			t = p.peek()
		}

		if right == nil {
			return nil, false
		}
		return &ast.App{Node: ast.Node{P: t.Pos}, Fn: left, Arg: right}, true

	} else {
		// Parse as atom first
		expr, ok := p.parseAtom()
		if !ok || expr == nil {
			return nil, false
		}

		fmt.Printf("DEBUG parseApp: after atom, expr=%s, next=%s\n", ast.Show(expr), p.peek().Text)

		for {
			t := p.peek()
			if t.Text == ")" || t.Text == "]" {
				fmt.Printf("DEBUG parseApp: breaking on %s\n", t.Text)
				return expr, false
			}
			// Check if this is a tuple or list continuation
			if t.Text == "(" || t.Text == "[" {
				// This is part of a nested tuple/list, not an app
				break
			}
			// Collect all valid argument types: IDENT, LONGID, CON, SYMBOL, INT, STRING
			if t.Kind != token.IDENT && t.Kind != token.LONGID && t.Kind != token.CON &&
				t.Kind != token.SYMBOL && t.Kind != token.INT && t.Kind != token.STRING {
				break
			}

			arg, ok := p.parseAtom()
			if !ok || arg == nil {
				return nil, false
			}
			expr = &ast.App{Fn: expr, Arg: arg}
			fmt.Printf("DEBUG parseApp: created app, next=%s\n", p.peek().Text)
		}

		return expr, true
	}
}

func (p *parser) parseAtom() (ast.Expr, bool) {
	t := p.peek()

	switch t.Kind {
	case token.INT:
		p.next()
		text := t.Text
		var v int64
		var err error
		if strings.HasPrefix(text, "~") {
			v, err = strconv.ParseInt(text[1:], 10, 64)
			if err != nil {
				return nil, false
			}
			v = -v
		} else {
			v, err = strconv.ParseInt(text, 10, 64)
			if err != nil {
				return nil, false
			}
		}
		return &ast.IntLit{Node: ast.Node{P: t.Pos}, Value: v}, true

	case token.STRING:
		p.next()
		return &ast.StrLit{Node: ast.Node{P: t.Pos}, Value: t.Text}, true

	case token.IDENT, token.LONGID, token.TYVAR:
		p.next()
		return &ast.Var{Node: ast.Node{P: t.Pos}, Name: t.Text}, true

	case token.CON:
		p.next()
		return &ast.Con{Node: ast.Node{P: t.Pos}, Name: t.Text}, true

	default:
		return nil, false
	}
}
