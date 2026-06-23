package analyser

import (
	"github.com/rafaeldepontes/comp/ast"
)

func (a *Analyser) checkVarDecl(node ast.VarDeclStmt) {
	if _, has := a.Scp.Lookup(node.VariableName); has {
		a.Errorf("[ERROR] %s is already declared in the current scope",
			node.VariableName,
		)
	}

	assignType := a.TypeCheckExpr(node.AssignedValue)
	if !node.ExplicitType.Equals(assignType) {
		a.Errorf(
			"[ERROR] cannot assign value of type %s to variable %s of type %s",
			assignType.String(),
			node.VariableName,
			node.ExplicitType.String(),
		)
		return
	}

	sym := &Symbol{
		Name:       node.VariableName,
		Type:       node.ExplicitType,
		IsConstant: node.IsConstant,
		IsGlobal:   a.Scp.IsGlobal(),
	}

	a.Scp.Define(node.VariableName, sym)
}

func (a *Analyser) checkIf(node ast.IfStmt) {
	// cond := a.TypeCheckExpr()
}

func (a *Analyser) checkFunc(node ast.FuncStmt) {
	if _, ok := a.Scp.Lookup(node.Name); ok {
		a.Errorf(
			"[ERROR] %s is already defined",
			node.Name,
		)
		return
	}

	params := make([]ast.Type, 0)
	for i := range node.Params {
		params = append(params, a.TypeCheckExpr(node.Params[i]))
	}

	a.Scp.Define(node.Name, &Symbol{
		Name:     node.Name,
		IsGlobal: a.Scp.IsGlobal(),
		Type: ast.FunctionType{
			Name:       node.Name,
			Params:     params,
			ReturnType: node.ReturnType,
		},
	})

	scp := NewScope(a.Scp)
	a.Scp = scp
}
