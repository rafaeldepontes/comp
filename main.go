package main

import (
	"fmt"
	"os"

	"github.com/rafaeldepontes/comp/lexer"
)

var TestFilePaths = []string{
	"./examples/test_case_01.rcs",
	"./examples/test_case_02.rcs",
	"./examples/test_case_03.rcs",
}

func main() {
	for i := range TestFilePaths {
		b, err := os.ReadFile(TestFilePaths[i])
		if err != nil {
			panic("[ERROR] missing example file")
		}
		src := string(b)

		tokens := lexer.Tokenize(TestFilePaths[i], src)
		for j := range tokens {
			tokens[j].Debbug()
		}

		fmt.Printf("\n\ncode snippet inside examples file: %s\n\n%s\n", TestFilePaths[i], src)
		println("=============\n\n")
	}
}
