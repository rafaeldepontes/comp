package parser

import (
	"errors"

	"github.com/rafaeldepontes/comp/ast"
	"github.com/rafaeldepontes/comp/lexer"
)

func (p *Parser) parseStmt() {
	if p.advance().Kind == lexer.Exit {
		// p.advance() //consume '('

		if expr := p.parseExpr(); expr != nil {
			p.Program.Statements = append(
				p.Program.Statements,
				ast.ExitStmt{
					Expr: expr,
				},
			)

			_ = p.expect(lexer.SemiColon)
		} else {
			p.Errors = append(p.Errors, errors.New("[ERROR] invalid statement"))
		}
	}
}
