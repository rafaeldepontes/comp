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
		p.errors = append(
			p.errors,
			fmt.Errorf("[ERROR] nud handler missing for %s",
				lexer.TokenTypeString(p.currentTokenType()),
			),
		)
		p.synchronize()
		return ast.NumberExpr{} // dummy to avoid crash
	}

	left := nudFn(p)
	for bpLT[p.currentTokenType()] > bp {
		type_ = p.currentTokenType()

		ledFn, has := ledLT[type_]
		if !has {
			p.errors = append(
				p.errors,
				fmt.Errorf(
					"[ERROR] led handler missing for %s",
					lexer.TokenTypeString(p.currentTokenType()),
				),
			)
			return left
		}

		left = ledFn(p, left, bpLT[p.currentTokenType()])
	}

	return left
}

func parsePrimaryExpr(p *parser) ast.Expr {
	switch p.currentTokenType() {
	case lexer.Number:
		number, _ := strconv.ParseFloat(p.advance().Val, 64)
		return ast.NumberExpr{
			Val: number,
		}

	case lexer.Identifier:
		return ast.SymbolExpr{
			Val: p.advance().Val,
		}

	default:
		p.errors = append(
			p.errors,
			fmt.Errorf(
				"[ERROR] cannot create primary expression from %s\n",
				lexer.TokenTypeString(p.currentTokenType()),
			),
		)
		return ast.NumberExpr{}
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
