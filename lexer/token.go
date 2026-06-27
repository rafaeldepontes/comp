package lexer

type TokenKind int

const (
	_ TokenKind = iota
	Exit
	Int
	SemiColon
	EOF
	Unknown
)

var kwd = map[string]TokenKind{
	"exit": Exit,
}

type Token struct {
	Value string
	Kind  TokenKind
}

func newToken(val string, kind TokenKind) Token {
	return Token{
		Value: val,
		Kind:  kind,
	}
}

func (t Token) Debug() string {
	switch t.Kind {
	case Exit:
		return "exit"
	case Int:
		return "integer"
	case SemiColon:
		return "semi colon"
	case EOF:
		return "eof"
	default:
		return "unknown token"
	}
}

func TokenKindString(t TokenKind) string {
	switch t {
	case Exit:
		return "exit"
	case Int:
		return "integer"
	case SemiColon:
		return "semi colon"
	case EOF:
		return "eof"
	default:
		return "unknown token"
	}
}
