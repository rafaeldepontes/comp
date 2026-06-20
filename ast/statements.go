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

type ImportStmt struct {
	PackageName string
}

func (i ImportStmt) stmt() {}

type FromImportStmt struct {
	PackageName string
	ImportName  string
}

func (f FromImportStmt) stmt() {}

type StructFields struct {
	Name         string
	Type         Type
	DefaultValue Expr
}
type StructStmt struct {
	Name   string
	Fields []StructFields
}

func (s StructStmt) stmt() {}

type FuncParam struct {
	Name string
	Type Type
}

type FuncStmt struct {
	Function
}

func (f FuncStmt) stmt() {}

type ImplStmt struct {
	Name    string
	Methods []FuncStmt
}

func (i ImplStmt) stmt() {}

type ClassStmt struct {
	Name    string
	Fields  []StructFields
	Methods []FuncStmt
}

func (c ClassStmt) stmt() {}

type IfStmt struct {
	Condition Expr
	Then      BlockStmt
	Else      Stmt
}

func (i IfStmt) stmt() {}

type WhileStmt struct {
	Condition Expr
	Body      BlockStmt
}

func (w WhileStmt) stmt() {}

type ForEachStmt struct {
	Item     string
	Index    string
	Iterable Expr
	Body     BlockStmt
}

func (f ForEachStmt) stmt() {}

type ForStmt struct {
	Init Stmt
	Cond Expr
	Post Expr
	Body BlockStmt
}

func (f ForStmt) stmt() {}
