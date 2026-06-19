package parser

import (
	"errors"

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

func parseValDeclStmt(p *parser) ast.Stmt {
	isConstant := p.advance().Type == lexer.Const
	varName := p.expectError(
		lexer.Identifier,
		errors.New("inside variable declaration expected to find variable name"),
	).Val

	_ = p.expect(lexer.Assignment)
	assignVal := parseExpr(p, Assignment)
	_ = p.expect(lexer.SemiColon)

	return ast.VarDeclStmt{
		VariableName:  varName,
		IsConstant:    isConstant,
		AssignedValue: assignVal,
	}
}
