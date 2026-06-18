package ast

// { [...]Stmt }
type BlockStmt struct {
	Body []Stmt
}

func (bs BlockStmt) stmt() {

}

// any shit...*;* <--
type ExpressionStmt struct {
	Expression Expr
}

func (es ExpressionStmt) stmt() {

}
