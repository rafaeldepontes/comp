package analyser

import (
	"github.com/rafaeldepontes/comp/ast"
	"github.com/rafaeldepontes/comp/lexer"
)

func (a *Analyser) checkBinary(expr ast.BinaryExpr) ast.Type {
	leftT := a.TypeCheckExpr(expr.Left)
	rightT := a.TypeCheckExpr(expr.Right)

	switch expr.Opr.Type {
	case lexer.Plus:
		if leftT.GetType() == ast.String && rightT.GetType() == ast.String {
			return ast.PrimitiveType{Type: ast.String, Val: ""}
		}
		fallthrough

	case lexer.Dash, lexer.Percent, lexer.Slash, lexer.Star:
		if leftT.GetType() == ast.Number && rightT.GetType() == ast.Number {
			return ast.PrimitiveType{Type: ast.Number, Val: 0}
		}

		a.Errorf(
			"[ERROR] invalid operation: %s %s %s",
			leftT.String(),
			lexer.TokenTypeString(expr.Opr.Type),
			rightT.String(),
		)
		return ast.PrimitiveType{Type: ast.Invalid}

	case lexer.Equals, lexer.NotEquals:
		if leftT.Equals(rightT) {
			return ast.PrimitiveType{Type: ast.Boolean, Val: leftT}
		}

		a.Errorf(
			"[ERROR] invalid operation: %s %s %s",
			leftT.String(),
			lexer.TokenTypeString(expr.Opr.Type),
			rightT.String(),
		)
		return ast.PrimitiveType{Type: ast.Invalid}

	default:
		return ast.PrimitiveType{Type: ast.Invalid}
	}
}

func (a *Analyser) checkSymbol(expr ast.SymbolExpr) ast.Type {
	sym, has := a.Scp.Lookup(expr.Val)
	if !has {
		a.Errorf(
			"[ERROR] undefined variable: %s",
			expr.Val,
		)
		return ast.PrimitiveType{Type: ast.Invalid}
	}

	return sym.Type
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
		a.Errorf(
			"[ERROR] undefined variable: %s",
			lhs.Val,
		)
		return ast.PrimitiveType{Type: ast.Invalid}
	}

	if sym.IsConstant {
		a.Errorf(
			"[ERROR] cannot assign to constant %s",
			sym.Name,
		)
		return ast.PrimitiveType{Type: ast.Invalid}
	}

	if !sym.Type.Equals(value) && value.GetType() != ast.Null {
		a.Errorf(
			"[ERROR] cannot assign value of type %s to variable %s of type %s",
			value.String(),
			sym.Name,
			sym.Type.String(),
		)
		return ast.PrimitiveType{Type: ast.Invalid}
	}

	// change this...
	return sym.Type
}

// TODO: Test to see if this is actually right
// and if it's the correct behaviour.
func (a *Analyser) checkCall(expr ast.CallExpr) ast.Type {
	fn := a.TypeCheckExpr(expr.Callee)

	if fn.GetType() != ast.Fn {
		a.Errorf(
			"[ERROR] %s is not a function",
			fn.String(),
		)
		return ast.PrimitiveType{Type: ast.Invalid}
	}

	// var _ ast.NamedType
	// if fn.GetType() == ast.Struct {
	// st := fn.(ast.StructFields)

	// TODO: check structs methods???
	// }

	val, ok := fn.(ast.FunctionType)
	if !ok {
		a.Errorf(
			"[ERROR] %s is not a function",
			fn.String(),
		)
		return ast.PrimitiveType{Type: ast.Invalid}
	}

	args := make([]ast.Type, 0)
	for i := range expr.Args {
		args = append(args, a.TypeCheckExpr(expr.Args[i]))
	}

	if len(val.Params) != len(args) {
		a.Errorf(
			"[ERROR] %s expected %d but received %d parameters",
			val.Name,
			len(val.Params),
			len(args),
		)
		return ast.PrimitiveType{Type: ast.Invalid}
	}

	for i := range val.Params {
		if !val.Params[i].Equals(args[i]) {
			a.Errorf(
				"[ERROR] expected %s but got %s",
				val.Params[i].String(),
				args[i].String(),
			)
			return ast.PrimitiveType{Type: ast.Invalid}
		}
	}

	return val.ReturnType
}
