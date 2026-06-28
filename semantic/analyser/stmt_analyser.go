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
	a.TypeCheckExpr(node.Condition)

	scope := NewScope(a.Scp)
	prev := a.Scp
	a.Scp = scope
	for _, stmt := range node.Then.Body {
		a.WalkStmt(stmt)
	}
	a.Scp = prev

	if node.Else != nil {
		scopeElse := NewScope(a.Scp)
		a.Scp = scopeElse
		a.WalkStmt(node.Else)
		a.Scp = prev
	}
}

func (a *Analyser) checkWhile(node ast.WhileStmt) {
	a.TypeCheckExpr(node.Condition)

	scope := NewScope(a.Scp)
	prev := a.Scp
	a.Scp = scope
	for _, stmt := range node.Body.Body {
		a.WalkStmt(stmt)
	}
	a.Scp = prev
}

func (a *Analyser) checkFor(node ast.ForStmt) {
	scope := NewScope(a.Scp)
	prev := a.Scp
	a.Scp = scope

	if node.Init != nil {
		a.WalkStmt(node.Init)
	}
	if node.Cond != nil {
		a.TypeCheckExpr(node.Cond)
	}
	if node.Post != nil {
		a.TypeCheckExpr(node.Post)
	}

	for _, stmt := range node.Body.Body {
		a.WalkStmt(stmt)
	}
	a.Scp = prev
}

func (a *Analyser) checkForEach(node ast.ForEachStmt) {
	a.TypeCheckExpr(node.Iterable)

	scope := NewScope(a.Scp)
	prev := a.Scp
	a.Scp = scope

	// For basic struct-only language, fallback to double type for loop items
	a.Scp.Define(node.Item, &Symbol{
		Name:       node.Item,
		Type:       ast.PrimitiveType{Type: ast.Number},
		IsConstant: false,
		IsGlobal:   false,
	})

	if node.Index != "" {
		a.Scp.Define(node.Index, &Symbol{
			Name:       node.Index,
			Type:       ast.PrimitiveType{Type: ast.Number},
			IsConstant: false,
			IsGlobal:   false,
		})
	}

	for _, stmt := range node.Body.Body {
		a.WalkStmt(stmt)
	}
	a.Scp = prev
}

func (a *Analyser) checkReturn(node ast.ReturnStmt) {
	expr := a.TypeCheckExpr(node.ReturnValue)

	println(expr)
}

func (a *Analyser) checkBlock(node ast.BlockStmt) {
	scope := NewScope(a.Scp)
	prev := a.Scp
	a.Scp = scope
	for _, stmt := range node.Body {
		a.WalkStmt(stmt)
	}
	a.Scp = prev
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

func (a *Analyser) checkStruct(node ast.StructStmt) {
	if _, ok := a.Scp.Lookup(node.Name); ok {
		a.Errorf("[ERROR] %s is already defined", node.Name)
		return
	}

	structType := ast.StructType{
		Name:    node.Name,
		Fields:  make(map[string]ast.Type),
		Methods: make(map[string]ast.FunctionType),
	}

	for i := range node.Fields {
		structType.Fields[node.Fields[i].Name] = node.Fields[i].Type
	}

	a.Scp.Define(node.Name, &Symbol{
		Name:     node.Name,
		Type:     structType,
		IsGlobal: a.Scp.IsGlobal(),
	})
}

func (a *Analyser) checkImpl(node ast.ImplStmt) {
	sym, ok := a.Scp.Lookup(node.Name)
	if !ok {
		a.Errorf("[ERROR] struct %s is not defined", node.Name)
		return
	}

	structType, ok := sym.Type.(ast.StructType)
	if !ok {
		a.Errorf("[ERROR] %s is not a struct", node.Name)
		return
	}

	if structType.Methods == nil {
		structType.Methods = make(map[string]ast.FunctionType)
	}

	for i := range node.Methods {
		if _, exists := structType.Methods[node.Methods[i].Name]; exists {
			a.Errorf("[ERROR] method %s is already defined in struct %s", node.Methods[i].Name, node.Name)
			continue
		}

		params := make([]ast.Type, 0)
		for j := range node.Methods[i].Params {
			params = append(params, node.Methods[i].Params[j].Type)
		}

		fnType := ast.FunctionType{
			Name:       node.Methods[i].Name,
			Params:     params,
			ReturnType: node.Methods[i].ReturnType,
		}

		structType.Methods[node.Methods[i].Name] = fnType

		methodScope := NewScope(a.Scp)
		prevScope := a.Scp
		a.Scp = methodScope

		a.Scp.Define("this", &Symbol{
			Name:     "this",
			Type:     structType,
			IsGlobal: false,
		})

		for j := range node.Methods[i].Params {
			a.Scp.Define(node.Methods[i].Params[j].Name, &Symbol{
				Name:       node.Methods[i].Params[j].Name,
				Type:       node.Methods[i].Params[j].Type,
				IsConstant: false,
				IsGlobal:   false,
			})
		}

		for j := range node.Methods[i].Body.Body {
			a.WalkStmt(node.Methods[i].Body.Body[j])
		}

		a.Scp = prevScope
	}

	sym.Type = structType
}
