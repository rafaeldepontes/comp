package ast

type BlockStmt struct {
	Statements []Stmt
}

func (BlockStmt) stmt() {}

type ExitStmt struct {
	Expr Expr
}

func (ExitStmt) stmt() {}

type ExprStmt struct {
	Expr Expr
}

func (ExprStmt) stmt() {}
