package parser

import (
	"errors"
	"fmt"

	"github.com/rafaeldepontes/comp/ast"
	"github.com/rafaeldepontes/comp/lexer"
)

func parseType(p *parser, msg string) string {
	if msg == "" {
		msg = "expected type name"
	}

	if p.currentTokenType() == lexer.OpenBracket {
		p.advance()
		_ = p.expect(lexer.CloseBracket)
		return "[]" + parseType(p, msg)
	}
	return p.expectError(lexer.Identifier, errors.New(msg)).Val
}

func parseStmt(p *parser) ast.Stmt {
	stmtFn, ok := stmtLT[p.currentTokenType()]
	if ok {
		return stmtFn(p)
	}

	expr := parseExpr(p, DefaultBP)
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
	varName := parseType(p, "inside variable declaration expected to find variable name")

	var type_ ast.Type
	if p.currentTokenType() == lexer.Colon {
		p.advance()
		type_ = parseType(p, "")
	}

	var assignVal ast.Expr
	switch p.currentTokenType() {
	case lexer.PlusEquals, lexer.MinusEquals, lexer.SlashEquals, lexer.StarEquals, lexer.PercentEquals, lexer.NullishAssignment:
		assignVal = parseAssignExpr(p, ast.SymbolExpr{Val: varName}, Assignment)
	case lexer.Assignment:
		_ = p.expect(lexer.Assignment)
		assignVal = parseExpr(p, DefaultBP)
	case lexer.SemiColon:
		// No assignment, assignVal remains nil
	default:
		_ = p.expect(lexer.Assignment)
	}
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

	nameOrPkg := parseType(p, "inside import declaration expected to find package name or import name")

	if p.currentTokenType() == lexer.From {
		p.advance()

		pkgName := parseType(p, "expected package name after from")
		_ = p.expect(lexer.SemiColon)
		return ast.FromImportStmt{
			PackageName: pkgName,
			ImportName:  nameOrPkg,
		}
	}

	_ = p.expect(lexer.SemiColon)

	return ast.ImportStmt{
		PackageName: nameOrPkg,
	}
}

func parseStructStmt(p *parser) ast.Stmt {
	p.advance()
	name := parseType(p, "inside struct declaration expected to find struct name")

	_ = p.expect(lexer.OpenCurly)

	fields := make([]ast.StructFields, 0)
	for p.currentTokenType() != lexer.CloseCurly && p.hasTokens() {
		fn := parseType(p, "expected field name in struct declaration")

		_ = p.expect(lexer.Colon)

		ft := parseType(p, "")

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

	name := parseType(p, "inside fn declaration expected to find function name")

	function := parseFuncGeneric(p)
	function.Name = name

	return ast.FuncStmt{
		Function: function,
	}
}

func parseImplStmt(p *parser) ast.Stmt {
	p.advance() // consume 'impl'

	name := parseType(p, "inside impl declaration expected to find struct name")

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

	condition := parseExpr(p, DefaultBP)
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

	condition := parseExpr(p, DefaultBP)
	body := parseBlockStmt(p)

	return ast.WhileStmt{
		Condition: condition,
		Body:      body,
	}
}

func parseForEachStmt(p *parser) ast.Stmt {
	_ = p.expect(lexer.Foreach)

	itemName := parseType(p, "expected variable name in foreach loop")

	var index string
	if p.currentTokenType() == lexer.In {
		_ = p.expect(lexer.In)
	} else if p.currentTokenType() == lexer.Comma {
		p.advance() // consume ','
		index = parseType(p, "expected variable name in foreach loop")
		_ = p.expect(lexer.In)
	}

	iterable := parseExpr(p, DefaultBP)
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
			expr := parseExpr(p, DefaultBP)
			_ = p.expect(lexer.SemiColon)
			init = ast.ExpressionStmt{Expression: expr}
		}
	} else {
		_ = p.expect(lexer.SemiColon)
	}

	var cond ast.Expr
	if p.currentTokenType() != lexer.SemiColon {
		cond = parseExpr(p, DefaultBP)
	}
	_ = p.expect(lexer.SemiColon)

	var post ast.Expr
	if p.currentTokenType() != lexer.OpenCurly {
		post = parseExpr(p, DefaultBP)
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

	name := parseType(p, "inside class declaration expected to find class name")

	_ = p.expect(lexer.OpenCurly)

	fields := make([]ast.StructFields, 0)
	methods := make([]ast.FuncStmt, 0)
	for p.currentTokenType() != lexer.CloseCurly && p.hasTokens() {
		switch p.currentTokenType() {
		case lexer.Fn:
			methods = append(methods, parseFuncStmt(p).(ast.FuncStmt))

		case lexer.Let, lexer.Const:
			p.advance()
			fn := parseType(p, "expected field name in class declaration")

			_ = p.expect(lexer.Colon)

			ft := parseType(p, "")

			var defaultVal ast.Expr
			if p.currentTokenType() == lexer.Assignment {
				p.advance() // consume '='
				defaultVal = parseExpr(p, DefaultBP)
			}

			fields = append(fields, ast.StructFields{
				Name:         fn,
				Type:         ft,
				DefaultValue: defaultVal,
			})

			if p.currentTokenType() == lexer.SemiColon {
				p.advance()
			} else if p.currentTokenType() != lexer.CloseCurly {
				panic(
					fmt.Sprintf(
						"expected ';' or '}' after field declaration, but got %s",
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
