package parser

import (
	"fmt"

	"github.com/rafaeldepontes/comp/ast"
	"github.com/rafaeldepontes/comp/lexer"
)

type parser struct {
	tokens []lexer.Token
	errors []error
	pos    int
}

func (p *parser) hasTokens() bool {
	return p.pos < len(p.tokens) && p.currentTokenType() != lexer.EOF
}

func (p *parser) currentTokenType() lexer.TokenType {
	return p.currentToken().Type
}

func (p *parser) currentToken() lexer.Token {
	return p.tokens[p.pos]
}

func (p *parser) advance() lexer.Token {
	i := p.pos
	p.pos++
	return p.tokens[i]
}

func (p *parser) expectError(expectedType lexer.TokenType, err error) lexer.Token {
	token := p.currentToken()
	type_ := token.Type

	if type_ != expectedType {
		if err == nil {
			err = fmt.Errorf("[ERROR] expected %s but received %s instead\n",
				lexer.TokenTypeString(expectedType),
				lexer.TokenTypeString(type_),
			)
		}
		panic(err)
	}

	return p.advance()
}

func (p *parser) expect(type_ lexer.TokenType) lexer.Token {
	return p.expectError(type_, nil)
}

func Parse(tokens []lexer.Token) ast.BlockStmt {
	createTokensLookups()
	p := &parser{
		pos:    0,
		tokens: tokens,
		errors: make([]error, 0),
	}

	body := make([]ast.Stmt, 0)
	for p.hasTokens() {
		body = append(body, parseStmt(p))
	}

	return ast.BlockStmt{
		Body: body,
	}
}
