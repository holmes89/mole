package token

type Kind int

const (
	EOF     Kind = iota
	INT          // 42 or ~5 (Text keeps the ~)
	STRING       // Text is the unescaped contents, without the quotes
	IDENT        // x, foo', ~
	CON          // Leaf, SOME, true, false, nil
	LONGID       // Int.toString
	TYVAR        // 'a
	KEYWORD      // see Keywords
	SYMBOL       // ( ) [ ] , ; | _ => -> = <> < > <= >= + - * ^ :: @
)

// Pos is a source position. Line and Col start at 1. Offset is a byte offset from 0.
type Pos struct{ Line, Col, Offset int }

type Token struct {
	Kind Kind
	Text string
	Pos  Pos
}

var Keywords = map[string]bool{
	"andalso": true,
	"orelse":   true,
	"case":     true,
	"of":       true,
	"datatype": true,
	"exception": true,
	"raise":    true,
	"handle":   true,
	"div":      true,
	"mod":      true,
	"o":        true,
	"val":      true,
	"fun":      true,
	"fn":       true,
	"let":      true,
	"in":       true,
	"end":      true,
	"if":       true,
	"then":     true,
	"else":     true,
}

func (k Kind) String() string {
	switch k {
	case EOF:
		return "EOF"
	case INT:
		return "INT"
	case STRING:
		return "STRING"
	case IDENT:
		return "IDENT"
	case CON:
		return "CON"
	case LONGID:
		return "LONGID"
	case TYVAR:
		return "TYVAR"
	case KEYWORD:
		return "KEYWORD"
	case SYMBOL:
		return "SYMBOL"
	}
	return ""
}
