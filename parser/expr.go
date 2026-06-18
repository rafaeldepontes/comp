package parser

import (
	"fmt"
	"strconv"

	"github.com/rafaeldepontes/comp/ast"
	"github.com/rafaeldepontes/comp/lexer"
)

func parseExpr(p *parser, bp BindingPower) ast.Expr {
	type_ := p.currentTokenType()

	nudFn, ok := nudLT[type_]
	if !ok {
		panic("[ERROR] somehow I've coded poorly 1.0... Jesus man... Maybe is the age??")
		// NUD handler missing...
	}

	left := nudFn(p)
	for bpLT[p.currentTokenType()] > bp {
		type_ = p.currentTokenType()

		ledFn, has := ledLT[type_]
		if !has {
			panic("[ERROR] somehow I've coded poorly 2.0... Jesus man... Maybe is the age??")
			// LED handler missing...
		}

		left = ledFn(p, left, bpLT[p.currentTokenType()])
	}

	return left
}

func parsePrimaryExpr(p *parser) ast.Expr {
	switch p.currentTokenType() {
	case lexer.Number:
		number, _ := strconv.ParseFloat(p.advance().Val, 64)
		return ast.NumberExpr{Val: number}

	case lexer.String:
		return ast.StringExpr{Val: p.advance().Val}

	case lexer.Identifier:
		return ast.SymbolExpr{Val: p.advance().Val}

	default:
		panic(
			fmt.Sprintf(
				"[ERROR] cannot create primary expression from %s\n",
				lexer.TokenTypeString(p.currentTokenType()),
			),
		)
	}
}

func parseBinaryExpr(p *parser, left ast.Expr, bp BindingPower) ast.Expr {
	optToken := p.advance()
	right := parseExpr(p, bp)

	return ast.BinaryExpr{
		Left:  left,
		Opr:   optToken,
		Right: right,
	}
}
