package parser

import (
	"fmt"

	"github.com/holmes89/mole/token"
)

type Error struct {
	Pos token.Pos
	Msg string
}

func (e Error) Error() string {
	return fmt.Sprintf("%d:%d: %s", e.Pos.Line, e.Pos.Col, e.Msg)
}

type parser struct {
	toks []token.Token
	pos  int
}

func (p *parser) peek() token.Token {
	if p.pos >= len(p.toks) {
		return token.Token{Kind: token.EOF}
	}
	return p.toks[p.pos]
}

func (p *parser) next() token.Token {
	if p.pos < len(p.toks) {
		t := p.toks[p.pos]
		p.pos++
		return t
	}
	return token.Token{Kind: token.EOF}
}
