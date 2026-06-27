package ast

import "github.com/rafaeldepontes/comp/lexer"

type NodeExpr struct {
	Value string
	Kind  lexer.TokenKind
}

func (n NodeExpr) expr() {}
