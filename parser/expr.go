package parser

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/rafaeldepontes/comp/ast"
	"github.com/rafaeldepontes/comp/lexer"
)

func parseExpr(p *parser, bp BindingPower) ast.Expr {
	type_ := p.currentTokenType()

	nudFn, ok := nudLT[type_]
	if !ok {
		// panic("[ERROR] somehow I've coded poorly 1.0... Jesus man... Maybe is the age??")
		panic(
			fmt.Errorf(
				"[ERROR] nud handler missing for %s",
				lexer.TokenTypeString(p.currentTokenType()),
			),
		)
		// NUD handler missing...
	}

	left := nudFn(p)
	for bpLT[p.currentTokenType()] > bp {
		type_ = p.currentTokenType()

		ledFn, has := ledLT[type_]
		if !has {
			// panic("[ERROR] somehow I've coded poorly 2.0... Jesus man... Maybe is the age??")
			panic(
				fmt.Errorf(
					"[ERROR] led handler missing for %s",
					lexer.TokenTypeString(p.currentTokenType()),
				),
			)
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

func parseThisExpr(p *parser) ast.Expr {
	return ast.ThisExpr{Token: p.advance()}
}

func parseGroupingExpr(p *parser) ast.Expr {
	_ = p.expect(lexer.OpenParen)
	expr := parseExpr(p, DefaltBP)
	_ = p.expect(lexer.CloseParen)
	return expr
}

func parseNewExpr(p *parser) ast.Expr {
	newToken := p.advance() // consume 'new'
	className := p.expectError(
		lexer.Identifier,
		errors.New("expected class or struct name after 'new'"),
	).Val

	_ = p.expect(lexer.OpenParen)
	args := make([]ast.Expr, 0)
	for p.currentTokenType() != lexer.CloseParen && p.hasTokens() {
		args = append(args, parseExpr(p, DefaltBP))
		if p.currentTokenType() == lexer.Comma {
			p.advance()
		} else if p.currentTokenType() != lexer.CloseParen {
			panic(fmt.Sprintf("expected ',' or ')' after argument, but got %s", lexer.TokenTypeString(p.currentTokenType())))
		}
	}
	_ = p.expect(lexer.CloseParen)

	return ast.NewExpr{
		Token:     newToken,
		ClassName: className,
		Args:      args,
	}
}

func parseCallExpr(p *parser, left ast.Expr, bp BindingPower) ast.Expr {
	_ = p.expect(lexer.OpenParen)
	args := make([]ast.Expr, 0)
	for p.currentTokenType() != lexer.CloseParen && p.hasTokens() {
		args = append(args, parseExpr(p, DefaltBP))
		if p.currentTokenType() == lexer.Comma {
			p.advance()
		} else if p.currentTokenType() != lexer.CloseParen {
			panic(fmt.Sprintf("expected ',' or ')' after argument, but got %s", lexer.TokenTypeString(p.currentTokenType())))
		}
	}
	_ = p.expect(lexer.CloseParen)

	return ast.CallExpr{
		Callee: left,
		Args:   args,
	}
}

func parseMemberExpr(p *parser, left ast.Expr, bp BindingPower) ast.Expr {
	opToken := p.advance() // consume '.'
	propName := p.expectError(
		lexer.Identifier,
		errors.New("expected property name after '.'"),
	).Val

	return ast.MemberExpr{
		Object:   left,
		Operator: opToken,
		Property: ast.SymbolExpr{Val: propName},
	}
}

func parsePrefixExpr(p *parser) ast.Expr {
	opt := p.advance()
	rhs := parseExpr(p, DefaltBP)

	return ast.UpdateExpr{
		Opr:      opt,
		Operand:  rhs,
		IsPrefix: true,
	}
}

func parsePostfixExpr(p *parser, left ast.Expr, bp BindingPower) ast.Expr {
	opt := p.advance()

	return ast.UpdateExpr{
		Opr:     opt,
		Operand: left,
	}
}

func parseAssignExpr(p *parser, left ast.Expr, bp BindingPower) ast.Expr {
	opt := p.advance()
	rhs := parseExpr(p, bp)

	return ast.AssignExpr{
		Assigne:  left,
		Operator: opt,
		Value:    rhs,
	}
}
