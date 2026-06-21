package analyser

import (
	"fmt"

	"github.com/rafaeldepontes/comp/ast"
	"github.com/rafaeldepontes/comp/lexer"
)

func (a *Analyser) checkBinary(expr ast.BinaryExpr) ast.Type {
	leftT := a.TypeCheckExpr(expr.Left)
	rightT := a.TypeCheckExpr(expr.Right)

	switch expr.Opr.Type {
	case lexer.Plus:
		if leftT.GetType() == ast.String && rightT.GetType() == ast.String {
			return ast.PrimitiveType{Type: ast.String}
		}
		fallthrough

	case lexer.Dash, lexer.Percent, lexer.Slash, lexer.Star:
		if leftT.GetType() == ast.Number && rightT.GetType() == ast.Number {
			return ast.PrimitiveType{Type: ast.Number}
		}

		a.Error(
			fmt.Sprintf(
				"[ERROR] invalid operation: %s %s %s",
				leftT.String(),
				lexer.TokenTypeString(expr.Opr.Type),
				rightT.String(),
			),
		)
		return ast.PrimitiveType{Type: ast.Invalid}

	case lexer.Equals, lexer.NotEquals:
		if leftT.Equals(rightT) {
			return ast.PrimitiveType{Type: ast.Boolean}
		}

		a.Error(
			fmt.Sprintf(
				"[ERROR] invalid operation: %s %s %s",
				leftT.String(),
				lexer.TokenTypeString(expr.Opr.Type),
				rightT.String(),
			),
		)
		return ast.PrimitiveType{Type: ast.Invalid}

	default:
		return ast.PrimitiveType{Type: ast.Invalid}
	}
}

func (a *Analyser) checkSymbol(expr ast.SymbolExpr) ast.Type {
	sym, has := a.Scp.Lookup(expr.Val)
	if !has {
		a.Error(
			fmt.Sprintf(
				"[ERROR] undefined variable: %s",
				expr.Val,
			),
		)
		return ast.PrimitiveType{Type: ast.Invalid}
	}

	return ast.NamedType{
		Name: sym.Name,
	}
}

func (a *Analyser) checkAssign(expr ast.AssignExpr) ast.Type {
	value := a.TypeCheckExpr(expr.Value)
	lhs, ok := expr.Assigne.(ast.SymbolExpr)
	if !ok {
		a.Error("[ERROR] left side of assignment must be assignable")
		return ast.PrimitiveType{Type: ast.Invalid}
	}

	sym, has := a.Scp.Lookup(lhs.Val)
	if !has {
		a.Error(
			fmt.Sprintf(
				"[ERROR] undefined variable: %s",
				lhs.Val,
			),
		)
		return ast.PrimitiveType{Type: ast.Invalid}
	}

	if sym.IsConstant {
		a.Error(
			fmt.Sprintf(
				"[ERROR] cannot assign to constant %s",
				sym.Name,
			),
		)
		return ast.PrimitiveType{Type: ast.Invalid}
	}

	if !sym.Type.Equals(value) && value.GetType() != ast.Null {
		a.Error(
			fmt.Sprintf(
				"[ERROR] cannot assign value of type %s to variable %s of type %s",
				value.String(),
				sym.Name,
				sym.Type.String(),
			),
		)
		return ast.PrimitiveType{Type: ast.Invalid}
	}

	// change this...
	return sym.Type
}
