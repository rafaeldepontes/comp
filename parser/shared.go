package parser

import (
	"fmt"

	"github.com/rafaeldepontes/comp/ast"
	"github.com/rafaeldepontes/comp/lexer"
)

func parseFuncGeneric(p *parser) ast.Function {
	_ = p.expect(lexer.OpenParen)
	params := make([]ast.FuncParam, 0)
	for p.currentTokenType() != lexer.CloseParen && p.hasTokens() {
		paramName := parseType(p, "expected parameter name in function declaration")

		_ = p.expect(lexer.Colon)

		paramType := parseType(p, "")

		params = append(params, ast.FuncParam{
			Name: paramName,
			Type: paramType,
		})

		if p.currentTokenType() == lexer.Comma {
			p.advance()
		} else if p.currentTokenType() != lexer.CloseParen {
			panic(
				fmt.Sprintf(
					"expected ',' or ')' after parameter, but got %s",
					lexer.TokenTypeString(p.currentTokenType()),
				),
			)
		}
	}
	_ = p.expect(lexer.CloseParen)

	var rt ast.Type
	if p.currentTokenType() == lexer.Colon {
		p.advance() // consume ':'
		rt = parseType(p, "")
	}

	body := parseBlockStmt(p)
	return ast.Function{
		Params:     params,
		ReturnType: rt,
		Body:       body,
	}
}
