package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rafaeldepontes/comp/lexer"
)

func chooseTokenizer(type_, path, src string, tokens *[]lexer.Token) {
	switch type_ {
	case "1":
		*tokens = lexer.TokenizeRegex(path, src)

	case "2":
		*tokens = lexer.TokenizeStateMachine(path, src)

	default:
		*tokens = lexer.TokenizeStateMachine(path, src)
	}
}

func TUI(type_ *string, enabled bool) {
	if enabled {
		text := "Choose your lexer type (e.g.: 1): "
		fmt.Println("\n> 1. Regex\n> 2. State Machine")
		fmt.Println("\033[5A\033")

		fmt.Printf("\033[%dC", len(text))
		fmt.Println("\r")
		print(text)

		reader := bufio.NewReader(os.Stdin)
		opt, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		*type_ = strings.TrimSpace(opt)
		fmt.Println("\033c")
	}
}
