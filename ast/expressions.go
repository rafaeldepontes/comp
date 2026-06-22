package ast

import (
	"github.com/rafaeldepontes/comp/lexer"
)

type NumberExpr struct {
	Val float64
}

func (n NumberExpr) expr() {}

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
