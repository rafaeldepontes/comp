package analyser

import (
	"fmt"

	"github.com/rafaeldepontes/comp/ast"
	"github.com/sanity-io/litter"
)

func Analyses(p ast.BlockStmt) any {
	a := Analyser{
		Scp:    NewScope(nil), // Global scope...
		Errors: make([]SemanticError, 0),
	}

	for i := range p.Body {
		a.WalkStmt(p.Body[i])
	}

	fmt.Println("[INFO] semantic analyser errors:")
	for i := range a.Errors {
		fmt.Println(a.Errors[i].Message)
	}

	litter.Dump(a)
	return nil
}
