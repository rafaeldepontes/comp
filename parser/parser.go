package parser

import (
	"fmt"

	"github.com/rafaeldepontes/comp/ast"
	"github.com/rafaeldepontes/comp/lexer"
)

type Parser struct {
	Tokens  []lexer.Token
	Errors  []error
	Program ast.Program
	pos     int
}

func newParser(tokens []lexer.Token) *Parser {
	return &Parser{
		Tokens: tokens,
		Errors: make([]error, 0),
		pos:    0,
	}
}

func (p *Parser) hasTokens() bool {
	return p.pos < len(p.Tokens) && p.currentTokenType() != lexer.EOF
}

func (p *Parser) currentTokenType() lexer.TokenKind {
	return p.currentToken().Kind
}

func (p *Parser) currentToken() lexer.Token {
	return p.Tokens[p.pos]
}

func (p *Parser) advance() lexer.Token {
	i := p.pos
	p.pos++
	return p.Tokens[i]
}

func (p *Parser) expectError(expectedType lexer.TokenKind, err error) lexer.Token {
	token := p.currentToken()
	type_ := token.Kind

	if type_ != expectedType {
		if err == nil {
			err = fmt.Errorf("[ERROR] expected %s but received %s instead\n",
				lexer.TokenKindString(expectedType),
				lexer.TokenKindString(type_),
			)
		}
		p.Errors = append(p.Errors, err)
		return lexer.Token{
			Kind:  expectedType,
			Value: "",
		}
	}

	return p.advance()
}

func (p *Parser) expect(type_ lexer.TokenKind) lexer.Token {
	return p.expectError(type_, nil)
}

func Parse(tokens []lexer.Token) []ast.Stmt {
	p := newParser(tokens)

	for p.hasTokens() {
		if p.currentTokenType() == lexer.EOF {
			continue
		}

		p.parseStmt()
	}

	return p.Program.Statements
}
