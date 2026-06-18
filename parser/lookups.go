package parser

import (
	"github.com/rafaeldepontes/comp/ast"
	"github.com/rafaeldepontes/comp/lexer"
)

type BindingPower int

const (
	DefaltBP BindingPower = iota
	Comma
	Assignment
	Logical
	Relational
	Additive
	Multiplicative
	Unary
	Call
	Member
	Primary
)

type stmtHandler func(p *parser) ast.Stmt
type nudHandler func(p *parser) ast.Expr
type ledHandler func(p *parser, left ast.Expr, bp BindingPower) ast.Expr

type stmtLookup map[lexer.TokenType]stmtHandler
type nudLookup map[lexer.TokenType]nudHandler
type ledLookup map[lexer.TokenType]ledHandler
type bpLookup map[lexer.TokenType]BindingPower

var bpLT = bpLookup{}
var nudLT = nudLookup{}
var ledLT = ledLookup{}
var stmtLT = stmtLookup{}

func led(type_ lexer.TokenType, bp BindingPower, ledFn ledHandler) {
	bpLT[type_] = bp
	ledLT[type_] = ledFn
}

func nud(type_ lexer.TokenType, nudFn nudHandler) {
	bpLT[type_] = Primary
	nudLT[type_] = nudFn
}

func stmt(type_ lexer.TokenType, stmtFn stmtHandler) {
	bpLT[type_] = DefaltBP
	stmtLT[type_] = stmtFn
}

func createTokensLookups() {
	led(lexer.And, Logical, parseBinaryExpr)
	led(lexer.Or, Logical, parseBinaryExpr)
	led(lexer.DotDot, Logical, parseBinaryExpr)

	led(lexer.Less, Relational, parseBinaryExpr)
	led(lexer.LessEquals, Relational, parseBinaryExpr)
	led(lexer.Greater, Relational, parseBinaryExpr)
	led(lexer.GreaterEquals, Relational, parseBinaryExpr)
	led(lexer.Equals, Relational, parseBinaryExpr)
	led(lexer.NotEquals, Relational, parseBinaryExpr)

	led(lexer.Plus, Additive, parseBinaryExpr)
	led(lexer.Dash, Additive, parseBinaryExpr)

	led(lexer.Star, Multiplicative, parseBinaryExpr)
	led(lexer.Slash, Multiplicative, parseBinaryExpr)
	led(lexer.Percent, Multiplicative, parseBinaryExpr)

	nud(lexer.Number, parsePrimaryExpr)
	nud(lexer.String, parsePrimaryExpr)
	nud(lexer.Identifier, parsePrimaryExpr)
}
