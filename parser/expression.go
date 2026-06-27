package parser

import (
	"errors"

	"github.com/rafaeldepontes/comp/ast"
	"github.com/rafaeldepontes/comp/lexer"
)

func (p *Parser) parseExpr() ast.Expr {
	switch node := p.advance(); node.Kind {
	case lexer.Int:
		return ast.NodeExpr{
			Value: node.Value,
			Kind:  lexer.Int,
		}

	default:
		p.Errors = append(p.Errors, errors.New("[ERROR] invalid expression"))
		return nil
	}
}
