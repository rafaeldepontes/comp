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

	type_ := node.ExplicitType
	if type_ != nil {
		if !type_.Equals(assignType) {
			a.Errorf(
				"[ERROR] cannot assign value of type %s to variable %s of type %s",
				assignType.String(),
				node.VariableName,
				node.ExplicitType.String(),
			)
			return
		}
	} else {
		type_ = assignType
	}

	sym := &Symbol{
		Name:       node.VariableName,
		Type:       type_,
		IsConstant: node.IsConstant,
		IsGlobal:   a.Scp.IsGlobal(),
		Value:      assignType.GetValue(),
	}

	a.Scp.Define(node.VariableName, sym)
}

func (a *Analyser) checkIf(node ast.IfStmt) {
	// cond := a.TypeCheckExpr()
}

func (a *Analyser) checkReturn(node ast.ReturnStmt) {
	expr := a.TypeCheckExpr(node.ReturnValue)

	println(expr)
}

func (a *Analyser) checkBlock(node ast.BlockStmt) {
	a.WalkStmt(node)
}

func (a *Analyser) checkFunc(node ast.FuncStmt) {
	if _, ok := a.Scp.Lookup(node.Name); ok {
		a.Errorf(
			"[ERROR] %s is already defined",
			node.Name,
		)
		return
	}

	a.Scp.Define(node.Name, &Symbol{
		Name:     node.Name,
		IsGlobal: a.Scp.IsGlobal(),
		Type:     ast.FunctionType{},
	})

	scp := NewScope(a.Scp)
	a.Scp = scp

	params := make([]ast.Type, 0)
	for i := range node.Params {
		params = append(params, a.TypeCheckExpr(node.Params[i]))
		a.Scp.Define(node.Params[i].Name, &Symbol{
			Name:       node.Params[i].Name,
			Type:       node.Params[i].Type,
			IsConstant: false,
			IsGlobal:   a.Scp.IsGlobal(),
		})
	}

	for i := range node.Body.Body {
		a.WalkStmt(node.Body.Body[i])
	}

	symb, _ := a.Scp.Lookup(node.Name)
	symb.Type = ast.FunctionType{
		Name:       node.Name,
		Params:     params,
		ReturnType: node.ReturnType,
	}
}
