package ast

import (
	"github.com/rafaeldepontes/comp/lexer"
)

type NumberExpr struct {
	Val float64
}

func (n NumberExpr) expr() {}

type StringExpr struct {
	Val string
}

func (s StringExpr) expr() {}

type SymbolExpr struct {
	Val string
}

func (n SymbolExpr) expr() {}

type BinaryExpr struct {
	Left  Expr
	Opr   lexer.Token
	Right Expr
}

func (b BinaryExpr) expr() {}

type NewExpr struct {
	Token     lexer.Token // 'new'
	ClassName string
	Args      []Expr
}

func (n NewExpr) expr() {}

type ThisExpr struct {
	Token lexer.Token // 'this'
}

func (t ThisExpr) expr() {}

type CallExpr struct {
	Callee Expr
	Args   []Expr
}

func (c CallExpr) expr() {}

type MemberExpr struct {
	Object   Expr
	Operator lexer.Token // '.'
	Property Expr        // Identifier
}

func (m MemberExpr) expr() {}

type UpdateExpr struct {
	Opr      lexer.Token
	Operand  Expr
	IsPrefix bool
}

func (p UpdateExpr) expr() {}

type AssignExpr struct {
	Assigne  Expr
	Operator lexer.Token
	Value    Expr
}

func (a AssignExpr) expr() {}

type BooleanExpr struct {
	Val bool
}

func (b BooleanExpr) expr() {}

type NullExpr struct {}

func (n NullExpr) expr() {}

type ArrayLiteralExpr struct {
	Elements []Expr
}

func (a ArrayLiteralExpr) expr() {}
