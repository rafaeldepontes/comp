package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rafaeldepontes/comp/lexer"
	"github.com/rafaeldepontes/comp/parser"
	semanticAnalyser "github.com/rafaeldepontes/comp/semantic/analyser"
)

var TestFilePaths = []string{
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

// func main() {
// 	type_ := ""
// 	TUI(&type_, false)

// 	for i := range TestFilePaths {
// 		b, err := os.ReadFile(TestFilePaths[i])
// 		if err != nil {
// 			panic("[ERROR] missing example file")
// 		}
// 		src := string(b)

// 		var tokens []lexer.Token
// 		chooseTokenizer(type_, TestFilePaths[i], src, &tokens)

// 		// Tokens are alright I guess...
// 		// for j := range tokens {
// 		// 	tokens[j].Debbug()
// 		// }

// 		// AST seems to have little problems, but I need
// 		// to test my interpreter to be sure... So more tests
// 		// are needed in order to decide if this is really
// 		// correct or not.
// 		ast := parser.Parse(tokens)

// 		if len(ast.Errors) > 0 {
// 			printLogs(ast, TestFilePaths[i], src)
// 		} else {
// 			fmt.Printf("%sFile: %s is OK\n%s", lexer.ColorBoldCyan, TestFilePaths[i], lexer.ColorReset)
// 		}

// 		semanticAnalyser.Analyses(ast)
// 	}
// }

func main() {
	type_ := ""
	TUI(&type_, false)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Comp running...")
	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break // EOF (Ctrl+D)
		}

		src := scanner.Text()

		if strings.TrimSpace(src) == "exit()" {
			fmt.Println("Goodbye!")
			break
		}

		var tokens []lexer.Token
		chooseTokenizer(type_, "<stdin>", src, &tokens)

		ast := parser.Parse(tokens)

		if len(ast.Errors) > 0 {
			printLogs(ast, "<stdin>", src)
		} else {
			fmt.Printf("%sInput is OK\n%s", lexer.ColorBoldCyan, lexer.ColorReset)
		}

		semanticAnalyser.Analyses(ast)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
