package main

import (
	"fmt"
	"os"

	"github.com/rafaeldepontes/comp/lexer"
	"github.com/rafaeldepontes/comp/parser"
	"github.com/sanity-io/litter"
)

var TestFilePaths = []string{
	"./examples/test_case_01.rcs",
	"./examples/test_case_02.rcs",
	"./examples/test_case_03.rcs",
	"./examples/test_case_04.rcs",
	"./examples/test_case_05.rcs",
	"./examples/test_case_06.rcs",
	"./examples/test_case_07.rcs",
	"./examples/test_case_08.rcs",
	"./examples/test_case_09.rcs",
}

func main() {
	type_ := ""
	// text := "Choose your lexer type (e.g.: 1): "
	// fmt.Println("\n> 1. Regex\n> 2. State Machine")
	// fmt.Println("\033[5A\033")

	// fmt.Printf("\033[%dC", len(text))
	// fmt.Println("\r")
	// print(text)

	// reader := bufio.NewReader(os.Stdin)
	// opt, err := reader.ReadString('\n')
	// if err != nil {
	// 	panic(err)
	// }

	// type_ := strings.TrimSpace(opt)
	// fmt.Println("\033c")

	for i := range TestFilePaths {
		b, err := os.ReadFile(TestFilePaths[i])
		if err != nil {
			panic("[ERROR] missing example file")
		}
		src := string(b)

		var tokens []lexer.Token
		switch type_ {
		case "1":
			tokens = lexer.TokenizeRegex(TestFilePaths[i], src)

		case "2":
			tokens = lexer.TokenizeStateMachine(TestFilePaths[i], src)

		default:
			tokens = lexer.TokenizeStateMachine(TestFilePaths[i], src)
		}

		// Tokens are alright I guess...
		// for j := range tokens {
		// 	tokens[j].Debbug()
		// }

		fmt.Printf("\n\ncode snippet inside examples file: %s\n\n%s\n", TestFilePaths[i], src)
		println("=============\n\n")

		ast := parser.Parse(tokens)

		// This is ok also... I think I can finally start writting the AST
		litter.Dump(ast)
	}
}
