package ast

// { [...]Stmt }
type BlockStmt struct {
	Body   []Stmt
	Errors []error
}

func (bs BlockStmt) stmt() {}

// any shit...*;* <--
type ExpressionStmt struct {
	Expression Expr
}

func (es ExpressionStmt) stmt() {}

type VarDeclStmt struct {
	VariableName  string
	IsConstant    bool
	AssignedValue Expr
	ExplicitType  Type
}

func (vd VarDeclStmt) stmt() {}
