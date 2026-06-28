package parser

import (
	"github.com/rafaeldepontes/comp/ast"
	"github.com/rafaeldepontes/comp/lexer"
)

type BindingPower int

const (
	DefaultBP BindingPower = iota
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
	nudLT[type_] = nudFn
}

func stmt(type_ lexer.TokenType, stmtFn stmtHandler) {
	bpLT[type_] = DefaultBP
	stmtLT[type_] = stmtFn
}

func createTokensLookups() {
	stmt(lexer.Fn, parseFuncStmt)
	stmt(lexer.Let, parseValDeclStmt)
	stmt(lexer.Const, parseValDeclStmt)
	stmt(lexer.Import, parseImportStmt)
	stmt(lexer.Struct, parseStructStmt)
	// stmt(lexer.Class, parseClassStmt)
	stmt(lexer.Impl, parseImplStmt)
	stmt(lexer.If, parseIfStmt)
	stmt(lexer.While, parseWhileStmt)
	stmt(lexer.Foreach, parseForEachStmt)
	stmt(lexer.For, parseForStmt)
	stmt(lexer.Return, parseReturnStmt)

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

	led(lexer.Assignment, Assignment, parseAssignExpr)
	led(lexer.PlusEquals, Assignment, parseAssignExpr)
	led(lexer.MinusEquals, Assignment, parseAssignExpr)
	led(lexer.StarEquals, Assignment, parseAssignExpr)
	led(lexer.SlashEquals, Assignment, parseAssignExpr)
	led(lexer.PercentEquals, Assignment, parseAssignExpr)
	led(lexer.NullishAssignment, Assignment, parseAssignExpr)

	led(lexer.PlusPlus, Unary, parsePostfixExpr)
	led(lexer.MinusMinus, Unary, parsePostfixExpr)

	led(lexer.OpenParen, Call, parseCallExpr)
	led(lexer.Dot, Member, parseMemberExpr)

	led(lexer.OpenBracket, Primary, parseIndexExpr)

	nud(lexer.Number, parsePrimaryExpr)
	nud(lexer.String, parsePrimaryExpr)
	nud(lexer.Identifier, parsePrimaryExpr)

	nud(lexer.True, parsePrimaryExpr)
	nud(lexer.False, parsePrimaryExpr)
	nud(lexer.Null, parsePrimaryExpr)

	nud(lexer.PlusPlus, parsePrefixExpr)
	nud(lexer.MinusMinus, parsePrefixExpr)

	nud(lexer.Plus, parsePrefixExpr)
	nud(lexer.Dash, parsePrefixExpr)
	nud(lexer.Not, parsePrefixExpr)

	nud(lexer.OpenParen, parseGroupingExpr)
	nud(lexer.OpenBracket, parseArrayExpr)
	nud(lexer.New, parseNewExpr)
	nud(lexer.This, parseThisExpr)
	nud(lexer.Fn, parseFuncExpr)
}
