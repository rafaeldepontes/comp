package analyser

import (
	"fmt"

	"github.com/rafaeldepontes/comp/ast"
)

type SemanticError struct {
	Message string
	// add column and lines from lexer...
}

type Analyser struct {
	Scp    *Scope
	Errors []SemanticError
}

func (a *Analyser) Error(msg string) {
	a.Errors = append(a.Errors, SemanticError{Message: msg})
}

func (a *Analyser) Errorf(format string, args ...any) {
	a.Errors = append(a.Errors, SemanticError{
		Message: fmt.Sprintf(format, args...),
	})
}

func (a *Analyser) WalkStmt(stmt ast.Stmt) {
	switch node := stmt.(type) {
	case ast.VarDeclStmt:
		a.checkVarDecl(node)
	case ast.ExpressionStmt:
		a.TypeCheckExpr(node.Expression)
	// case ast.BlockStmt:
	// 	a.checkBlock(node)
	case ast.IfStmt:
		a.checkIf(node)
	// case ast.WhileStmt:
	// 	a.checkWhile(node)
	// case ast.ForStmt:
	// 	a.checkFor(node)
	// case ast.ForEachStmt:
	// 	a.checkForEach(node)
	// case ast.FuncStmt:
	// 	a.checkFunc(node)
	// case ast.ImplStmt:
	// 	a.checkImpl(node)
	// case ast.ClassStmt:
	// 	a.checkClass(node)
	default:
		a.Error("unsupported statement type")
	}
}

func (a *Analyser) TypeCheckExpr(expr ast.Expr) ast.Type {
	switch node := expr.(type) {
	case ast.NumberExpr:
		return ast.PrimitiveType{Type: ast.Number}
	case ast.StringExpr:
		return ast.PrimitiveType{Type: ast.String}
	case ast.BooleanExpr:
		return ast.PrimitiveType{Type: ast.Boolean}
	case ast.SymbolExpr:
		return a.checkSymbol(node)
	case ast.BinaryExpr:
		return a.checkBinary(node)
	case ast.AssignExpr:
		return a.checkAssign(node)
	// case ast.CallExpr:
	// 	return a.checkCall(node)
	// case ast.MemberExpr:
	// 	return a.checkMember(node)
	// case ast.NewExpr:
	// 	return a.checkNew(node)
	default:
		a.Error("unsupported expression type")
		return ast.PrimitiveType{Type: ast.Invalid}
	}
}
