package analyser

import (
	"fmt"

	"github.com/rafaeldepontes/comp/ast"
)

func (a *Analyser) checkVarDecl(node ast.VarDeclStmt) {
	if _, has := a.Scp.Lookup(node.VariableName); has {
		a.Error(
			fmt.Sprintf("[ERROR] %s is already declared in the current scope",
				node.VariableName,
			),
		)
	}

	assignType := a.TypeCheckExpr(node.AssignedValue)
	if !node.ExplicitType.Equals(assignType) {
		a.Error(
			fmt.Sprintf(
				"[ERROR] cannot assign value of type %s to variable %s of type %s",
				assignType.String(),
				node.VariableName,
				node.ExplicitType.String(),
			),
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
