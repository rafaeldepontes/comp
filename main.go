package main

import (
	"os"

	"github.com/rafaeldepontes/comp/lexer"
)

var TestFilePaths = []string{
	"./examples/test_case_01.rcs",
	"./examples/test_case_02.rcs",
}

func main() {
	for i := range TestFilePaths {
		b, err := os.ReadFile(TestFilePaths[i])
		if err != nil {
			panic("[ERROR] missing example file")
		}
		src := string(b)

		tokens := lexer.Tokenize(src)
		for j := range tokens {
			tokens[j].Debbug()
		}

		println("\n\ncode snippet inside examples file:", src)
		println("=============\n\n")
	}
}
