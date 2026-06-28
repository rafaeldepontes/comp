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

	switch lhs := expr.Assigne.(type) {
	case ast.SymbolExpr:
		sym, has := a.Scp.Lookup(lhs.Val)
		if !has {
			a.Errorf("[ERROR] undefined variable: %s", lhs.Val)
			return ast.PrimitiveType{Type: ast.Invalid}
		}

		if sym.IsConstant {
			a.Errorf("[ERROR] cannot assign to constant %s", sym.Name)
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
		return sym.Type

	case ast.MemberExpr:
		memType := a.checkMember(lhs)
		if memType.GetType() == ast.Invalid {
			return ast.PrimitiveType{Type: ast.Invalid}
		}
		if !memType.Equals(value) && value.GetType() != ast.Null {
			a.Errorf(
				"[ERROR] cannot assign value of type %s to member of type %s",
				value.String(),
				memType.String(),
			)
			return ast.PrimitiveType{Type: ast.Invalid}
		}
		return memType

	default:
		a.Error("[ERROR] left side of assignment must be assignable")
		return ast.PrimitiveType{Type: ast.Invalid}
	}
}

func (a *Analyser) checkCall(expr ast.CallExpr) ast.Type {
	fn := a.TypeCheckExpr(expr.Callee)

	if fn.GetType() != ast.Fn {
		a.Errorf(
			"[ERROR] %s is not a function",
			fn.String(),
		)
		return ast.PrimitiveType{Type: ast.Invalid}
	}

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

func (a *Analyser) checkNew(expr ast.NewExpr) ast.Type {
	sym, ok := a.Scp.Lookup(expr.ClassName)
	if !ok {
		a.Errorf("[ERROR] struct %s is not defined", expr.ClassName)
		return ast.PrimitiveType{Type: ast.Invalid}
	}

	structType, ok := sym.Type.(ast.StructType)
	if !ok {
		a.Errorf("[ERROR] %s is not a struct", expr.ClassName)
		return ast.PrimitiveType{Type: ast.Invalid}
	}

	return structType
}

func (a *Analyser) checkMember(expr ast.MemberExpr) ast.Type {
	objType := a.TypeCheckExpr(expr.Object)

	if objType.GetType() != ast.Struct {
		a.Errorf("[ERROR] cannot access member of non-struct type %s", objType.String())
		return ast.PrimitiveType{Type: ast.Invalid}
	}

	var structType ast.StructType
	if st, ok := objType.(ast.StructType); ok {
		structType = st
	} else if nt, ok := objType.(ast.NamedType); ok {
		sym, has := a.Scp.Lookup(nt.Name)
		if !has {
			return ast.PrimitiveType{Type: ast.Invalid}
		}
		structType = sym.Type.(ast.StructType)
	} else {
		return ast.PrimitiveType{Type: ast.Invalid}
	}

	propExpr, ok := expr.Property.(ast.SymbolExpr)
	if !ok {
		a.Errorf("[ERROR] expected identifier after '.'")
		return ast.PrimitiveType{Type: ast.Invalid}
	}

	propName := propExpr.Val

	if fieldType, exists := structType.Fields[propName]; exists {
		return fieldType
	}

	if methodType, exists := structType.Methods[propName]; exists {
		return methodType
	}

	a.Errorf("[ERROR] struct %s has no field or method %s", structType.Name, propName)
	return ast.PrimitiveType{Type: ast.Invalid}
}

func (a *Analyser) checkThis() ast.Type {
	sym, ok := a.Scp.Lookup("this")
	if !ok {
		a.Errorf("[ERROR] 'this' can only be used inside a method")
		return ast.PrimitiveType{Type: ast.Invalid}
	}
	return sym.Type
}
