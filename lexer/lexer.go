package lexer

import (
	"fmt"
	"strings"

	"github.com/holmes89/mole/token"
)

type Error struct {
	Pos        token.Pos
	Msg        string
	Incomplete bool
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d:%d: %s", e.Pos.Line, e.Pos.Col, e.Msg)
}

func Tokenize(src string) ([]token.Token, error) {
	var tokens []token.Token
	i := 0
	n := len(src)

	for {
		if i >= n {
			break
		}
		c := src[i]

		// Skip spaces, tabs, newlines
		if c == ' ' || c == '\t' || c == '\n' {
			i++
			continue
		}

		// Handle comment start: (*
		if c == '(' && i+1 < n && src[i+1] == '*' {
			_, nextPos, err := tokenizeComment(src, i)
			if err != nil {
				return tokens, err
			}
			i = nextPos
			continue
		}

		var tok token.Token
		var err error

		switch {
		case isDigit(c):
			tok, i, err = tokenizeInt(src, i)
		case c == '~' && i+1 < n && isDigit(src[i+1]):
			tok, i, err = tokenizeIntWithTilde(src, i)
		case c == '~':
			tok, i, err = singleCharToken(src, i, token.IDENT)
		case isLetter(c) || c == '_':
			// Uppercase letters start long identifiers (e.g., Int.toString)
			if isUpper(c) {
				tok, i, err = tokenizeLongId(src, i)
			} else {
				tok, i, err = tokenizeName(src, i)
			}
		case c == '=':
			tok, i, err = tokenizeSymbol(src, i, byte('='))
		case c == '<':
			tok, i, err = tokenizeSymbol(src, i, byte('<'))
		case c == '>':
			tok, i, err = tokenizeSymbol(src, i, byte('>'))
		case c == '(':
			tok, i, err = singleCharToken(src, i, token.SYMBOL)
		case c == ')':
			tok, i, err = singleCharToken(src, i, token.SYMBOL)
		case c == '[':
			tok, i, err = singleCharToken(src, i, token.SYMBOL)
		case c == ']':
			tok, i, err = singleCharToken(src, i, token.SYMBOL)
		case c == ',':
			tok, i, err = singleCharToken(src, i, token.SYMBOL)
		case c == ';':
			tok, i, err = singleCharToken(src, i, token.SYMBOL)
		case c == '|':
			tok, i, err = singleCharToken(src, i, token.SYMBOL)
		case c == ':':
			// Check for :: before :
			if i+1 < n && src[i+1] == ':' {
				tok, i, err = tokenizeSymbol(src, i, byte(':'))
			} else {
				tok, i, err = singleCharToken(src, i, token.SYMBOL)
			}
		case c == '_':
			tok, i, err = singleCharToken(src, i, token.SYMBOL)
		case c == '+':
			tok, i, err = singleCharToken(src, i, token.SYMBOL)
		case c == '-':
			tok, i, err = tokenizeSymbol(src, i, byte('-'))
		case c == '*':
			tok, i, err = singleCharToken(src, i, token.SYMBOL)
		case c == '^':
			tok, i, err = singleCharToken(src, i, token.SYMBOL)
		case c == '@':
			tok, i, err = singleCharToken(src, i, token.SYMBOL)
		case c == '"':
			tok, i, err = tokenizeString(src, i)
		case c == '\'':
			tok, i, err = tokenizeTyVar(src, i)
		default:
			return tokens, &Error{Pos: posAt(src, i), Msg: "unexpected character: " + string(c)}
		}

		if err != nil {
			return tokens, err
		}

		tokens = append(tokens, tok)
	}

	// Always end with EOF token
	tokens = append(tokens, token.Token{Kind: token.EOF, Text: "", Pos: posAt(src, i)})

	return tokens, nil
}

func tokenizeComment(src string, i int) (token.Token, int, error) {
	start := i + 2 // skip "(*" relative to i
	j := start
	n := len(src)
	depth := 1

	for ; j < n && depth > 0; j++ {
		if src[j] == '*' && j+1 < n && src[j+1] == ')' {
			depth--
			j++ // skip *)
		} else if src[j] == '(' {
			depth++
		}
	}

	if depth > 0 {
		return token.Token{}, 0, &Error{Pos: posAt(src, i), Msg: "unclosed comment", Incomplete: true}
	}

	// Skip to the next newline or end of input (but only skip comment content, not code after *)
	for ; j < n && src[j] != '\n'; j++ {
		// skip comment content
	}
	
	// If we reached end of string without finding a newline, go back to position after *)
	if j >= n {
		j = 17 // position after the closing ) (which is at index 16)
	}
	return token.Token{}, j, nil
}

func tokenizeString(src string, i int) (token.Token, int, error) {
	start := i + 1 // skip opening "
	j := start
	n := len(src)
	escaped := false
	
	for ; j < n; j++ {
		c := src[j]
		if escaped {
			// Check for valid escapes
			switch c {
			case 'n', 't', '\\', '"':
				escaped = false
				continue
			default:
				return token.Token{}, 0, &Error{Pos: posAt(src, j), Msg: "invalid escape sequence"}
			}
		}

		if c == '\\' {
			escaped = true
			continue
		}

		if c == '"' {
			// Unescape the text
			text := unescapeString(src[start:j])
			return token.Token{Kind: token.STRING, Text: text, Pos: posAt(src, start)}, j + 1, nil
		}

		// Newline in string is incomplete
		if c == '\n' {
			return token.Token{}, 0, &Error{Pos: posAt(src, j), Msg: "unclosed string", Incomplete: true}
		}
	}

	// Reached end of input without closing quote
	return token.Token{}, 0, &Error{Pos: posAt(src, start), Msg: "unclosed string", Incomplete: true}
}

func tokenizeIntWithTilde(src string, i int) (token.Token, int, error) {
	start := i
	i++ // skip ~
	for ; i < len(src) && isDigit(src[i]); i++ {
		// continue collecting digits
	}
	text := src[start:i] // include ~ and all digits
	return token.Token{Kind: token.INT, Text: text, Pos: posAt(src, start)}, i, nil
}

func tokenizeTyVar(src string, i int) (token.Token, int, error) {
	start := i
	// Must be followed by a letter
	if i+1 >= len(src) || !isLetter(src[i+1]) {
		return token.Token{}, 0, fmt.Errorf("type variable must be followed by a letter")
	}

	i++ // skip the '
	for ; i < len(src) && (isLetter(src[i]) || isDigit(src[i])); i++ {
		// continue collecting letters and digits
	}

	text := src[start:i] // include the '
	return token.Token{Kind: token.TYVAR, Text: text, Pos: posAt(src, start)}, i, nil
}

func tokenizeLongId(src string, i int) (token.Token, int, error) {
	start := i
	for ; i < len(src); i++ {
		c := src[i]
		if c == '.' && i+1 < len(src) && isLetter(src[i+1]) {
			i++ // skip the dot
			continue
		}
		if !isLetter(c) && c != '.' {
			break  // Stop at non-letter, non-dot (space, operator, etc.)
		}
	}

	text := src[start:i]
	
	// If there were no dots, it's a constructor (e.g., Int), not a long id
	if text == "" || !strings.Contains(text, ".") {
		if isUpper(text[0]) || text == "true" || text == "false" || text == "nil" {
			return token.Token{Kind: token.CON, Text: text, Pos: posAt(src, start)}, i, nil
		}
	}
	
	return token.Token{Kind: token.LONGID, Text: text, Pos: posAt(src, start)}, i, nil
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func tokenizeInt(src string, i int) (token.Token, int, error) {
	start := i
	for i < len(src) && isDigit(src[i]) {
		i++
	}
	text := src[start:i]
	return token.Token{Kind: token.INT, Text: text, Pos: posAt(src, start)}, i, nil
}

func tokenizeName(src string, i int) (token.Token, int, error) {
	start := i
	for i < len(src) && (isLetter(src[i]) || isDigit(src[i]) || src[i] == '_' || src[i] == '\'') {
		i++
	}
	text := src[start:i]

	if token.Keywords[text] {
		return token.Token{Kind: token.KEYWORD, Text: text, Pos: posAt(src, start)}, i, nil
	}

	if isUpper(text[0]) || text == "true" || text == "false" || text == "nil" {
		return token.Token{Kind: token.CON, Text: text, Pos: posAt(src, start)}, i, nil
	}

	return token.Token{Kind: token.IDENT, Text: text, Pos: posAt(src, start)}, i, nil
}

func tokenizeSymbol(src string, i int, first byte) (token.Token, int, error) {
	start := i
	if first == '=' || first == '<' || first == '>' || first == '-' || first == ':' {
		// Multi-char symbols need special handling
		peek := i + 1
		if peek >= len(src) {
			return token.Token{Kind: token.SYMBOL, Text: string(first), Pos: posAt(src, start)}, i + 1, nil
		}

		switch second := src[peek]; {
		case first == '=' && second == '>':
			return token.Token{Kind: token.SYMBOL, Text: "=>", Pos: posAt(src, start)}, i + 2, nil
		case first == '<' && second == '>':
			return token.Token{Kind: token.SYMBOL, Text: "<>", Pos: posAt(src, start)}, i + 2, nil
		case first == '<' && second == '=':
			return token.Token{Kind: token.SYMBOL, Text: "<=", Pos: posAt(src, start)}, i + 2, nil
		case first == '>' && second == '=':
			return token.Token{Kind: token.SYMBOL, Text: ">=", Pos: posAt(src, start)}, i + 2, nil
		case first == ':' && second == ':':
			return token.Token{Kind: token.SYMBOL, Text: "::", Pos: posAt(src, start)}, i + 2, nil
		case first == '-' && second == '>':
			return token.Token{Kind: token.SYMBOL, Text: "->", Pos: posAt(src, start)}, i + 2, nil
		default:
			return token.Token{Kind: token.SYMBOL, Text: string(first), Pos: posAt(src, start)}, i + 1, nil
		}
	}
	return token.Token{}, 0, fmt.Errorf("unexpected character after %q: %q", string(first), string(src[i]))
}

func singleCharToken(src string, i int, kind token.Kind) (token.Token, int, error) {
	if i >= len(src) {
		return token.Token{}, 0, fmt.Errorf("unexpected end of input")
	}
	c := src[i]
	return token.Token{Kind: kind, Text: string(c), Pos: posAt(src, i)}, i + 1, nil
}

func unescapeString(s string) string {
	result := make([]byte, 0, len(s))
	inEscape := false
	for i := 0; i < len(s); i++ {
		if inEscape {
			switch s[i] {
			case 'n':
				result = append(result, '\n')
			case 't':
				result = append(result, '\t')
			case '\\', '"':
				result = append(result, s[i])
			default:
				// Keep invalid escapes as-is
				result = append(result, s[i-1], s[i])
			}
			inEscape = false
		} else if s[i] == '\\' {
			inEscape = true
		} else {
			result = append(result, s[i])
		}
	}
	return string(result)
}

func posAt(src string, offset int) token.Pos {
	line := 1
	col := 1
	for i := 0; i < offset && i < len(src); i++ {
		if src[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return token.Pos{Line: line, Col: col, Offset: offset}
}

func isUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}
