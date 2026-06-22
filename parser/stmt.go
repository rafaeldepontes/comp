package parser

import (
	"errors"

	"github.com/rafaeldepontes/comp/ast"
	"github.com/rafaeldepontes/comp/lexer"
)

func parseIdentifier(p *parser, msg string) string {
	if msg == "" {
		msg = "expected identifier"
	}
	return p.expectError(lexer.Identifier, errors.New(msg)).Val
}

func parseType(p *parser, msg string) ast.Type {
	if msg == "" {
		msg = "expected type name"
	}

	if p.currentTokenType() == lexer.OpenBracket {
		p.advance()
		_ = p.expect(lexer.CloseBracket)
		return ast.ArrayType{
			ElementType: parseType(p, msg),
		}
	}

	token := p.expectError(lexer.Identifier, errors.New(msg))
	typeName := token.Val

	switch typeName {
	case "number":
		return ast.PrimitiveType{Type: ast.Number}
	case "string":
		return ast.PrimitiveType{Type: ast.String}
	case "boolean":
		return ast.PrimitiveType{Type: ast.Boolean}
	case "null":
		return ast.PrimitiveType{Type: ast.Null}
	case "void":
		return ast.PrimitiveType{Type: ast.Void}
	default:
		return ast.NamedType{Name: typeName}
	}
}

func parseStmt(p *parser) ast.Stmt {
	stmtFn, ok := stmtLT[p.currentTokenType()]
	if ok {
		return stmtFn(p)
	}

	expr := parseExpr(p, DefaultBP)
	_ = p.expect(lexer.SemiColon)

	return ast.ExpressionStmt{
		Expression: expr,
	}
}

func parseValDeclStmt(p *parser) ast.Stmt {
	isConstant := p.advance().Type == lexer.Const
	varName := parseIdentifier(p, "inside variable declaration expected to find variable name")

	var type_ ast.Type
	if p.currentTokenType() == lexer.Colon {
		p.advance()
		type_ = parseType(p, "")
	}

	var assignVal ast.Expr
	switch p.currentTokenType() {
	case lexer.Assignment:
		_ = p.expect(lexer.Assignment)
		assignVal = parseExpr(p, DefaultBP)
	case lexer.SemiColon:
		// No assignment, assignVal remains nil
	default:
		_ = p.expect(lexer.Assignment)
	}
	_ = p.expect(lexer.SemiColon)

	return ast.VarDeclStmt{
		VariableName:  varName,
		IsConstant:    isConstant,
		AssignedValue: assignVal,
		ExplicitType:  type_,
	}
}
