package lexer

import "fmt"

type Lexer struct {
	Tokens []Token
	Errors []error
	src    string
	pos    int
	lin    int
	col    int
}

func newLexer(src string) *Lexer {
	return &Lexer{
		Tokens: make([]Token, 0),
		Errors: make([]error, 0),
		src:    src,
		pos:    0,
		lin:    1,
		col:    1,
	}
}

func (l *Lexer) hasData() bool {
	return l.pos < len(l.src)
}

func (l *Lexer) valid() bool {
	return l.isInteger() || l.isExit() || l.isNewLine(l.src[l.pos]) || l.isWhiteSpaces(l.src[l.pos])
}

func (l *Lexer) advance() {
	l.pos++
}

func Tokenize(src string) []Token {
	l := newLexer(src)

	for l.hasData() {
		if l.valid() {
			continue
		}

		switch l.src[l.pos] {
		case ';':
			l.Tokens = append(l.Tokens, newToken("", SemiColon))
			l.advance()
		default:
			l.Errors = append(l.Errors, fmt.Errorf("invalid type %s", string(l.src[l.pos])))
			l.advance()
		}
	}

	l.Tokens = append(l.Tokens, newToken("eof", EOF))
	return l.Tokens
}
