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
	expr, ok := p.parseExpr()
	if !ok {
		return nil, Error{Pos: p.peek().Pos, Msg: "expected expression"}
	}
	if p.peek().Kind != token.EOF {
		return nil, Error{Pos: p.peek().Pos, Msg: "unexpected token after expression"}
	}
	return expr, nil
}

func (p *parser) parseExpr() (ast.Expr, bool) {
	t := p.peek()
	switch t.Kind {
	case token.KEYWORD:
		switch t.Text {
		case "if":
			fmt.Printf("DEBUG: parsing if expr at pos %d (%s)\n", t.Pos, t.Text)
			p.next() // consume "if"
			fmt.Printf("DEBUG: after consuming if, at pos %d (%v %q)\n", p.peek().Pos, p.peek().Kind, p.peek().Text)
			cond, ok := p.parseExpr()
			if !ok {
				return nil, false
			} else if p.peek().Text != "then" {
				fmt.Printf("DEBUG: after cond, at pos %d (%v %q), expecting 'then' but got '%s'\n", p.peek().Pos, p.peek().Kind, p.peek().Text, p.peek().Text)
				p.next() // consume "then"
				thenExpr, ok := p.parseExpr()
				if !ok {
					return nil, false
				} else if p.peek().Text != "else" {
					fmt.Printf("DEBUG: after thenExpr, at pos %d (%v %q), expecting 'else' but got '%s'\n", p.peek().Pos, p.peek().Kind, p.peek().Text, p.peek().Text)
					p.next() // consume "else"
					elseExpr, ok := p.parseExpr()
					if !ok {
						return nil, false
					} else {
						expr := &ast.If{Node: ast.Node{P: t.Pos}, Cond: cond, Then: thenExpr, Else: elseExpr}
						return expr, true
					}
				}
			}
		case "fn":
			fmt.Printf("DEBUG: parsing fn expr at pos %d (%s)\n", t.Pos, t.Text)
			p.next() // consume "fn"
			fmt.Printf("DEBUG: after consuming fn, at pos %d (%v %q)\n", p.peek().Pos, p.peek().Kind, p.peek().Text)
			pats, ok := p.parsePats()
			if !ok {
				return nil, false
			} else {
				fmt.Printf("DEBUG: after parsePats(), at pos %d (%v %q), got %d patterns\n", p.peek().Pos, p.peek().Kind, p.peek().Text, len(pats))
				body, ok := p.parseExpr()
				if !ok {
					return nil, false
				} else {
					fmt.Printf("DEBUG: after parseExpr(body), at pos %d (%v %q)\n", p.peek().Pos, p.peek().Kind, p.peek().Text)
					expr := &ast.Fn{Node: ast.Node{P: t.Pos}, Pats: pats, Body: body}
					return expr, true
				}
			}
		case "let":
			fmt.Printf("DEBUG: parsing let expr at pos %d (%s)\n", t.Pos, t.Text)
			p.next() // consume "let"
			var e ast.Expr = nil
			var b bool = false
			e = nil
			b = false
			return e, b
		default:
			fmt.Printf("DEBUG: default case for KEYWORD at pos %d (%s)\n", t.Pos, t.Text)
			var e ast.Expr = nil
			var b bool = false
			e = nil
			b = false
			return e, b
		}
	case token.SYMBOL:
		fmt.Printf("DEBUG: SYMBOL case at pos %d (%v %q)\n", t.Pos, t.Kind, t.Text)
		if t.Text == "(" {
			return p.parseBinary(1)
		} else {
			return p.parseBinary(1)
		}
	}
	return p.parseBinary(1)
}

func (p *parser) parsePats() ([]ast.Pat, bool) {
	var pats []ast.Pat
	if p.peek().Text == "(" {
		p.next() // consume "("
		for p.peek().Text != ")" {
			pat, ok := p.parsePat()
			if !ok {
				return nil, false
			}
			pats = append(pats, pat)
			if p.peek().Text == "," {
				p.next() // consume ","
			} else if p.peek().Text == "=>" {
				break
			} else {
				break
			}
		}
		p.next() // consume ")"
		if len(pats) == 0 {
			return []ast.Pat{}, true
		}
		return pats, true
	} else {
		pat, ok := p.parsePat()
		if !ok {
			return nil, false
		}
		p.next() // consume "=>"
		return []ast.Pat{pat}, true
	}
}

func (p *parser) parsePat() (ast.Pat, bool) {
	t := p.peek()
	switch t.Kind {
	case token.SYMBOL:
		if t.Text == "(" {
			return p.parseParenPat()
		}
		if t.Text == "_" {
			p.next()
			return &ast.PWild{Node: ast.Node{P: t.Pos}}, true
		}
	case token.IDENT, token.LONGID:
		p.next()
		return &ast.PVar{Node: ast.Node{P: t.Pos}, Name: t.Text}, true
	}

	expr, ok := p.parseExpr()
	if ok {
		return expr, false
	}

	return nil, false
}

func (p *parser) parseParenPat() (ast.Pat, bool) {
	p.next() // consume "("
	var elems []ast.Pat
	for p.peek().Text != ")" {
		pat, ok := p.parsePat()
		if !ok {
			return nil, false
		}
		elems = append(elems, pat)
		if p.peek().Text == "," {
			p.next() // consume ","
		} else {
			break
		}
	}
	p.next() // consume ")"

	if len(elems) == 0 {
		return &ast.PTuple{Node: ast.Node{}, Elems: []ast.Pat{}}, true
	}
	if len(elems) == 1 {
		return elems[0], true
	}
	return &ast.PTuple{Node: ast.Node{P: elems[0].Pos()}, Elems: elems}, true
}

// operatorLevel returns the precedence level of an operator.
// Higher number means higher precedence (binds tighter).
func operatorLevel(t token.Token) int {
	switch t.Kind {
	case token.KEYWORD:
		switch t.Text {
		case "orelse":
			return 1
		case "andalso":
			return 2
		case "o":
			return 3
		case "div", "mod":
			return 7
		default:
			return 0
		}
	case token.SYMBOL:
		switch t.Text {
		case "::", "@":
			return 5
		case "+", "-", "^":
			return 6
		case "*", "div", "mod":
			return 7
		case "=", "<>", "<", ">", "<=", ">=":
			return 4
		default:
			return 0
		}
	default:
		return 0
	}
}

// isRightAssociative returns true if the operator at level l is right-associative.
func isRightAssociative(l int) bool {
	switch l {
	case 5: // ::, @
		return true
	default:
		return false
	}
}

func (p *parser) parseBinary(minLevel int) (ast.Expr, bool) {
	fmt.Printf("DEBUG: parseBinary(%d) at pos %d (%v %q)\n", minLevel, p.peek().Pos, p.peek().Kind, p.peek().Text)
	left, ok := p.parseApp()
	if !ok {
		return nil, false
	}

	for {
		t := p.peek()
		l := operatorLevel(t)

		if l < minLevel {
			break
		}

		fmt.Printf("DEBUG: consuming operator %s at pos %d (%v %q)\n", t.Text, t.Pos, t.Kind, t.Text)
		p.next() // consume the operator

		// For left-associative operators, parse right at level+1
		// For right-associative operators (like ::, @), parse right at same level
		parseRightLevel := l + 1
		if isRightAssociative(l) {
			parseRightLevel = l
		}

		fmt.Printf("DEBUG: about to call parseBinary(%d) for right side\n", parseRightLevel)
		right, ok := p.parseBinary(parseRightLevel)
		if !ok {
			return nil, false
		}

		var expr ast.Expr

		switch t.Text {
		case "orelse":
			expr = &ast.OrElse{Node: ast.Node{P: t.Pos}, L: left, R: right}
		case "andalso":
			expr = &ast.AndAlso{Node: ast.Node{P: t.Pos}, L: left, R: right}
		default:
			// For all other operators, create an App node with a tuple of arguments
			expr = &ast.App{
				Node: ast.Node{P: t.Pos},
				Fn:   &ast.Var{Node: ast.Node{P: t.Pos}, Name: t.Text},
				Arg:  &ast.Tuple{Node: ast.Node{P: t.Pos}, Elems: []ast.Expr{left, right}},
			}
		}

		left = expr
	}

	return left, true
}

func startsAtom(t token.Token) bool {
	switch t.Kind {
	case token.INT, token.STRING, token.IDENT, token.LONGID, token.CON:
		return true
	case token.SYMBOL:
		// Include => so that fn expressions work correctly
		if t.Text == "=>" {
			return true
		}
		return t.Text == "(" || t.Text == "["
	default:
		return false
	}
}

func (p *parser) parseApp() (ast.Expr, bool) {
	fmt.Printf("DEBUG: parseApp at pos %d (%v %q)\n", p.peek().Pos, p.peek().Kind, p.peek().Text)
	expr, ok := p.parseAtom()
	if !ok {
		return nil, false
	}
	for startsAtom(p.peek()) {
		arg, ok := p.parseAtom()
		if !ok {
			return nil, false
		}
		expr = &ast.App{Fn: expr, Arg: arg}
	}
	return expr, true
}

func (p *parser) parseAtom() (ast.Expr, bool) {
	fmt.Printf("DEBUG: parseAtom at pos %d (%v %q)\n", p.peek().Pos, p.peek().Kind, p.peek().Text)
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
	case token.SYMBOL:
		if t.Text == "(" || t.Text == "[" {
			return p.parseDelimited()
		}
	}
	return nil, false
}

func (p *parser) parseDelimited() (ast.Expr, bool) {
	open := p.next()
	close := ")"
	if open.Text == "[" {
		close = "]"
	}
	node := ast.Node{P: open.Pos}
	var elems []ast.Expr
	separator := ""
	if p.peek().Text != close {
		for {
			expr, ok := p.parseExpr()
			if !ok {
				return nil, false
			}
			elems = append(elems, expr)
			if p.peek().Text == close {
				break
			}
			if separator == "" {
				separator = p.peek().Text
				if separator != "," && (separator != ";" || open.Text != "(") {
					return nil, false
				}
			}
			if p.peek().Text != separator {
				return nil, false
			}
			p.next()
		}
	}
	p.next() // Consume the closing delimiter.
	if open.Text == "[" {
		return &ast.List{Node: node, Elems: elems}, true
	}
	if separator == ";" {
		return &ast.Seq{Node: node, Exprs: elems}, true
	}
	if len(elems) == 1 {
		return elems[0], true
	}
	return &ast.Tuple{Node: node, Elems: elems}, true
}
