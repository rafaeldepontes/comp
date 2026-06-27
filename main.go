package main

import (
	"fmt"
	"os"

	"github.com/rafaeldepontes/comp/builder"
	"github.com/rafaeldepontes/comp/lexer"
	"github.com/rafaeldepontes/comp/parser"
)

var TestFilePaths = []string{
	"./examples/test_case_01.dr",
}

func main() {
	println("dragged is running...")
	for i := range len(TestFilePaths) {
		b, err := os.ReadFile(TestFilePaths[i])
		if err != nil {
			panic(fmt.Errorf("[ERROR] incorrect path or missing dragged file..."))
		}
		src := string(b)

		tokens := lexer.Tokenize(src)
		if len(tokens) < 1 {
			panic("[ERROR] unable to parse")
		}

		body := parser.Parse(tokens)
		if len(body) < 1 {
			panic("[ERROR] unable to compile")
		}

		builder.Compile(body, i)
	}

	println("compilation completed!")
}
