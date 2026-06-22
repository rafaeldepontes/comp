package main

import (
	"fmt"
	"os"

	"github.com/rafaeldepontes/comp/lexer"
	"github.com/rafaeldepontes/comp/parser"
)

var TestFilePaths = []string{
	"./examples/test_case_01.rcs",
}

func main() {
	type_ := ""
	TUI(&type_, false)

	for i := range TestFilePaths {
		b, err := os.ReadFile(TestFilePaths[i])
		if err != nil {
			panic("[ERROR] missing example file")
		}
		src := string(b)

		var tokens []lexer.Token
		chooseTokenizer(type_, TestFilePaths[i], src, &tokens)

		ast := parser.Parse(tokens)

		if len(ast.Errors) == 0 {
			printLogs(ast, TestFilePaths[i], src)
		} else {
			fmt.Printf("%sFile: %s is OK\n%s", lexer.ColorBoldCyan, TestFilePaths[i], lexer.ColorReset)
		}
	}
}
