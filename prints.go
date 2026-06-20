package main

import (
	"fmt"

	"github.com/rafaeldepontes/comp/ast"
	"github.com/rafaeldepontes/comp/lexer"
	"github.com/sanity-io/litter"
)

func printLogs(ast ast.BlockStmt, path, src string) {
	fmt.Printf("%s%s%s\n", lexer.ColorBoldCyan, "=================================================================\n", lexer.ColorReset)
	fmt.Printf("%sFile: %s\n\n* Code Sample:%s\n%s\n", lexer.ColorBoldCyan, path, lexer.ColorReset, src)
	fmt.Printf("%s%s%s\n", lexer.ColorBoldCyan, "=================================================================\n", lexer.ColorReset)

	fmt.Printf("%s%s%s\n", lexer.ColorBoldCyan, "AST:", lexer.ColorReset)
	litter.Dump(ast.Body)
	fmt.Printf("%s%s%s\n", lexer.ColorBoldCyan, "=================================================================\n", lexer.ColorReset)

	fmt.Printf("%s[INFO] Parser errors:%s\n", lexer.ColorBoldCyan, lexer.ColorReset)
	fmt.Printf("%s", lexer.ColorBoldRed)
	for i := range ast.Errors {
		println(ast.Errors[i].Error())
	}
	fmt.Printf("%s", lexer.ColorReset)
}
