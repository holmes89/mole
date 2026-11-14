package parser

import (
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
	return p.parseBinary(1)
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

		p.next() // consume the operator

		// For left-associative operators, parse right at level+1
		// For right-associative operators (like ::, @), parse right at same level
		parseRightLevel := l + 1
		if isRightAssociative(l) {
			parseRightLevel = l
		}

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
		return t.Text == "(" || t.Text == "["
	default:
		return false
	}
}

func (p *parser) parseApp() (ast.Expr, bool) {
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
