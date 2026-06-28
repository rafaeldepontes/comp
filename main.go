package main

import (
	"fmt"
	"os"

	"github.com/rafaeldepontes/comp/builder"
	"github.com/rafaeldepontes/comp/lexer"
	"github.com/rafaeldepontes/comp/parser"
	semanticAnalyser "github.com/rafaeldepontes/comp/semantic/analyser"
)

var TestFilePaths = []string{
	"./examples/control_flow_test.rcs",
	"./examples/structs_test.rcs",
	"./examples/test_case_01.rcs",
	// "./examples/test_case_02.rcs",
	// "./examples/test_case_03.rcs",
	// "./examples/test_case_04.rcs",
	// "./examples/test_case_05.rcs",
	// "./examples/test_case_06.rcs",
	// "./examples/test_case_07.rcs",
	// "./examples/test_case_08.rcs",
	// "./examples/test_case_09.rcs",
	// "./examples/test_case_10.rcs",
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
		if len(ast.Errors) > 0 {
			printLogs(ast, TestFilePaths[i], src)
		} else {
			fmt.Printf("%sFile: %s is OK\n%s", lexer.ColorBoldCyan, TestFilePaths[i], lexer.ColorReset)
		}

		semanticAnalyser.Analyses(ast)

		if len(ast.Errors) == 0 {
			compBuilder := builder.NewBuilder()
			err := compBuilder.Build(ast, TestFilePaths[i])
			if err != nil {
				fmt.Printf("[ERROR] Builder failed: %v\n", err)
			}
		}
	}
}
