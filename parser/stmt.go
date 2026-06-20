package parser

import (
	"errors"
	"fmt"

	"github.com/rafaeldepontes/comp/ast"
	"github.com/rafaeldepontes/comp/lexer"
)

func parseStmt(p *parser) ast.Stmt {
	stmtFn, ok := stmtLT[p.currentTokenType()]
	if ok {
		return stmtFn(p)
	}

	expr := parseExpr(p, DefaltBP)
	_ = p.expect(lexer.SemiColon)

	return ast.ExpressionStmt{
		Expression: expr,
	}
}

func parseBlockStmt(p *parser) ast.BlockStmt {
	_ = p.expect(lexer.OpenCurly)

	body := make([]ast.Stmt, 0)
	for p.currentTokenType() != lexer.CloseCurly && p.hasTokens() {
		body = append(body, parseStmt(p))
	}

	_ = p.expect(lexer.CloseCurly)
	return ast.BlockStmt{
		Body: body,
	}
}

func parseValDeclStmt(p *parser) ast.Stmt {
	isConstant := p.advance().Type == lexer.Const
	varName := p.expectError(
		lexer.Identifier,
		errors.New("inside variable declaration expected to find variable name"),
	).Val

	var type_ ast.Type
	if p.currentTokenType() == lexer.Colon {
		p.advance()
		type_ = p.expectError(
			lexer.Identifier,
			errors.New("expected return type after ':' in variable declaration"),
		).Val
	}

	_ = p.expect(lexer.Assignment)
	assignVal := parseExpr(p, Assignment)
	_ = p.expect(lexer.SemiColon)

	return ast.VarDeclStmt{
		VariableName:  varName,
		IsConstant:    isConstant,
		AssignedValue: assignVal,
		ExplicitType:  type_,
	}
}

func parseImportStmt(p *parser) ast.Stmt {
	_ = p.expect(lexer.Import)

	pkgName := p.expectError(
		lexer.Identifier,
		errors.New("inside import declaration expected to find package name"),
	).Val
	_ = p.expect(lexer.SemiColon)

	return ast.ImportStmt{
		PackageName: pkgName,
	}
}

func parseStructStmt(p *parser) ast.Stmt {
	p.advance()
	name := p.expectError(
		lexer.Identifier,
		errors.New("inside struct declaration expected to find struct name"),
	).Val

	_ = p.expect(lexer.OpenCurly)

	fields := make([]ast.StructFields, 0)
	for p.currentTokenType() != lexer.CloseCurly && p.hasTokens() {
		fn := p.expectError(
			lexer.Identifier,
			errors.New("expected field name in struct declaration"),
		).Val

		_ = p.expect(lexer.Colon)

		ft := p.expectError(
			lexer.Identifier,
			errors.New("expected field type in struct declaration"),
		).Val

		fields = append(fields, ast.StructFields{
			Name: fn,
			Type: ft,
		})

		if p.currentTokenType() == lexer.Comma {
			p.advance()
		} else if p.currentTokenType() != lexer.CloseCurly {
			panic(
				fmt.Sprintf(
					"expected ',' or '}' after field declaration, but got %s",
					lexer.TokenTypeString(p.currentTokenType()),
				),
			)
		}
	}

	_ = p.expect(lexer.CloseCurly)

	return ast.StructStmt{
		Name:   name,
		Fields: fields,
	}
}

func parseFuncStmt(p *parser) ast.Stmt {
	_ = p.expect(lexer.Fn)

	name := p.expectError(
		lexer.Identifier,
		errors.New("inside fn declaration expected to find function name"),
	).Val

	_ = p.expect(lexer.OpenParen)

	params := make([]ast.FuncParam, 0)
	for p.currentTokenType() != lexer.CloseParen && p.hasTokens() {
		paramName := p.expectError(
			lexer.Identifier,
			errors.New("expected parameter name in function declaration"),
		).Val

		_ = p.expect(lexer.Colon)

		paramType := p.expectError(
			lexer.Identifier,
			errors.New("expected parameter type in function declaration"),
		).Val

		params = append(params, ast.FuncParam{
			Name: paramName,
			Type: paramType,
		})

		if p.currentTokenType() == lexer.Comma {
			p.advance()
		} else if p.currentTokenType() != lexer.CloseParen {
			panic(
				fmt.Sprintf(
					"expected ',' or ')' after parameter, but got %s",
					lexer.TokenTypeString(p.currentTokenType()),
				),
			)
		}
	}
	_ = p.expect(lexer.CloseParen)

	var rt ast.Type
	if p.currentTokenType() == lexer.Colon {
		p.advance() // consume ':'
		rt = p.expectError(
			lexer.Identifier,
			errors.New("expected return type after ':' in function declaration"),
		).Val
	}

	body := parseBlockStmt(p)
	return ast.FuncStmt{
		Name:       name,
		Parameters: params,
		ReturnType: rt,
		Body:       body,
	}
}

func parseImplStmt(p *parser) ast.Stmt {
	p.advance() // consume 'impl'

	name := p.expectError(
		lexer.Identifier,
		errors.New("inside impl declaration expected to find struct name"),
	).Val

	_ = p.expect(lexer.OpenCurly)

	methods := make([]ast.FuncStmt, 0)
	for p.currentTokenType() != lexer.CloseCurly && p.hasTokens() {
		if p.currentTokenType() != lexer.Fn {
			panic(
				fmt.Sprintf(
					"expected method declaration ('fn') inside impl block, but got %s",
					lexer.TokenTypeString(p.currentTokenType()),
				),
			)
		}

		methods = append(methods, parseFuncStmt(p).(ast.FuncStmt))
	}

	_ = p.expect(lexer.CloseCurly)

	return ast.ImplStmt{
		Name:    name,
		Methods: methods,
	}
}

func parseIfStmt(p *parser) ast.Stmt {
	_ = p.expect(lexer.If)

	condition := parseExpr(p, DefaltBP)
	thenBlock := parseBlockStmt(p)

	var elseBlock ast.Stmt
	if p.currentTokenType() == lexer.Else {
		p.advance()
		if p.currentTokenType() == lexer.If {
			elseBlock = parseIfStmt(p)
		} else {
			elseBlock = parseBlockStmt(p)
		}
	}

	return ast.IfStmt{
		Condition: condition,
		Then:      thenBlock,
		Else:      elseBlock,
	}
}

func parseWhileStmt(p *parser) ast.Stmt {
	_ = p.expect(lexer.While)

	condition := parseExpr(p, DefaltBP)
	body := parseBlockStmt(p)

	return ast.WhileStmt{
		Condition: condition,
		Body:      body,
	}
}

func parseForEachStmt(p *parser) ast.Stmt {
	_ = p.expect(lexer.Foreach)

	itemName := p.expectError(
		lexer.Identifier,
		errors.New("expected variable name in foreach loop"),
	).Val

	var index string
	if p.currentTokenType() == lexer.In {
		_ = p.expect(lexer.In)
	} else if p.currentTokenType() == lexer.Comma {
		p.advance() // consume ','
		index = p.expectError(
			lexer.Identifier,
			errors.New("expected variable name in foreach loop"),
		).Val
		_ = p.expect(lexer.In)
	}

	iterable := parseExpr(p, DefaltBP)
	body := parseBlockStmt(p)

	return ast.ForEachStmt{
		Item:     itemName,
		Index:    index,
		Iterable: iterable,
		Body:     body,
	}
}

func parseForStmt(p *parser) ast.Stmt {
	_ = p.expect(lexer.For)

	var init ast.Stmt
	if p.currentTokenType() != lexer.SemiColon {
		if p.currentTokenType() == lexer.Let || p.currentTokenType() == lexer.Const {
			init = parseValDeclStmt(p)
		} else {
			expr := parseExpr(p, DefaltBP)
			_ = p.expect(lexer.SemiColon)
			init = ast.ExpressionStmt{Expression: expr}
		}
	} else {
		_ = p.expect(lexer.SemiColon)
	}

	var cond ast.Expr
	if p.currentTokenType() != lexer.SemiColon {
		cond = parseExpr(p, DefaltBP)
	}
	_ = p.expect(lexer.SemiColon)

	var post ast.Expr
	if p.currentTokenType() != lexer.OpenCurly {
		post = parseExpr(p, DefaltBP)
	}

	body := parseBlockStmt(p)

	return ast.ForStmt{
		Init: init,
		Cond: cond,
		Post: post,
		Body: body,
	}
}

func parseClassStmt(p *parser) ast.Stmt {
	p.advance() // consume 'class'

	name := p.expectError(
		lexer.Identifier,
		errors.New("inside class declaration expected to find class name"),
	).Val

	_ = p.expect(lexer.OpenCurly)

	fields := make([]ast.StructFields, 0)
	methods := make([]ast.FuncStmt, 0)
	for p.currentTokenType() != lexer.CloseCurly && p.hasTokens() {
		switch p.currentTokenType() {
		case lexer.Fn:
			methods = append(methods, parseFuncStmt(p).(ast.FuncStmt))

		case lexer.Let, lexer.Const:
			p.advance()
			fn := p.expectError(
				lexer.Identifier,
				errors.New("expected field name in class declaration"),
			).Val

			_ = p.expect(lexer.Colon)

			ft := p.expectError(
				lexer.Identifier,
				errors.New("expected field type in class declaration"),
			).Val

			fields = append(fields, ast.StructFields{
				Name: fn,
				Type: ft,
			})

			if p.currentTokenType() == lexer.SemiColon {
				p.advance()
			} else if p.currentTokenType() != lexer.CloseCurly {
				panic(
					fmt.Sprintf(
						"expected ',' or '}' after field declaration, but got %s",
						lexer.TokenTypeString(p.currentTokenType()),
					),
				)
			}

		default:
			panic(
				fmt.Sprintf(
					"unexpected token %s inside class declaration",
					lexer.TokenTypeString(p.currentTokenType()),
				),
			)
		}
	}

	_ = p.expect(lexer.CloseCurly)

	return ast.ClassStmt{
		Name:    name,
		Fields:  fields,
		Methods: methods,
	}
}
