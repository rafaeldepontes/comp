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

	case lexer.String:
		return ast.StringExpr{
			Val: p.advance().Val,
		}

	case lexer.Identifier:
		return ast.SymbolExpr{
			Val: p.advance().Val,
		}

	case lexer.True:
		p.advance()
		return ast.BooleanExpr{
			Val: true,
		}

	case lexer.False:
		p.advance()
		return ast.BooleanExpr{
			Val: false,
		}

	case lexer.Null:
		p.advance()
		return ast.NullExpr{}

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

func parseThisExpr(p *parser) ast.Expr {
	return ast.ThisExpr{Token: p.advance()}
}

func parseGroupingExpr(p *parser) ast.Expr {
	_ = p.expect(lexer.OpenParen)
	expr := parseExpr(p, DefaultBP)
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
		args = append(args, parseExpr(p, DefaultBP))
		if p.currentTokenType() == lexer.Comma {
			p.advance()
		} else if p.currentTokenType() != lexer.CloseParen {
			panic(
				fmt.Sprintf(
					"expected ',' or ')' after argument, but got %s",
					lexer.TokenTypeString(p.currentTokenType()),
				),
			)
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
		args = append(args, parseExpr(p, DefaultBP))
		if p.currentTokenType() == lexer.Comma {
			p.advance()

		} else if p.currentTokenType() != lexer.CloseParen {
			panic(
				fmt.Sprintf("expected ',' or ')' after argument, but got %s",
					lexer.TokenTypeString(p.currentTokenType()),
				),
			)
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
	rhs := parseExpr(p, DefaultBP)

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
	rhs := parseExpr(p, bp-1)

	return ast.AssignExpr{
		Assigne:  left,
		Operator: opt,
		Value:    rhs,
	}
}

func parseArrayExpr(p *parser) ast.Expr {
	p.advance() // consume '['
	elements := make([]ast.Expr, 0)
	for p.currentTokenType() != lexer.CloseBracket && p.hasTokens() {
		elements = append(elements, parseExpr(p, DefaultBP))
		if p.currentTokenType() == lexer.Comma {
			p.advance()
		} else if p.currentTokenType() != lexer.CloseBracket {
			p.errors = append(
				p.errors,
				fmt.Errorf(
					"expected ',' or ']' in array literal, but got %s",
					lexer.TokenTypeString(p.currentTokenType()),
				),
			)
			break
		}
	}
	_ = p.expect(lexer.CloseBracket)
	return ast.ArrayLiteralExpr{
		Elements: elements,
	}
}

func parseIndexExpr(p *parser, left ast.Expr, bp BindingPower) ast.Expr {
	open := p.advance() // consume '['
	idx := parseExpr(p, DefaultBP)

	p.expect(lexer.CloseBracket)

	return ast.IndexExpr{
		Bracket:    open,
		Collection: left,
		Index:      idx,
	}
}

func parseFuncExpr(p *parser) ast.Expr {
	_ = p.advance() // consume 'fn'

	return ast.FuncDeclExpr{
		Function: parseFuncGeneric(p),
	}
}
