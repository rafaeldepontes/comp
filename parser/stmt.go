package parser

import (
	"github.com/rafaeldepontes/comp/ast"
	"github.com/rafaeldepontes/comp/lexer"
)

func parseStmt(p *parser) ast.Stmt {
	stmtFn, ok := stmtLT[p.currentTokenType()]
	if ok {
		return stmtFn(p)
	}

	expr := parseExpr(p, DefaltBP)
	_ = p.expect(lexer.SemiColon)

	return ast.ExpressionStmt{
		Expression: expr,
	}
}
