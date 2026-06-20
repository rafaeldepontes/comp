package lexer

import (
	"fmt"
	"strings"
)

const (
	ColorReset    = "\033[0m"
	ColorBold     = "\033[1m"
	ColorBoldRed  = "\033[1;31m"
	ColorBoldCyan = "\033[1;36m"
)

func ErrorHandler(l lexer, src, path string) {
	actualLine := 1
	actualCol := 1
	for i := 0; i < l.pos; i++ {
		if src[i] == '\n' {
			actualLine++
			actualCol = 1
		} else {
			actualCol++
		}
	}

	start := l.pos
	for start > 0 && src[start-1] != '\n' {
		start--
	}

	end := l.pos
	for end < len(src) && src[end] != '\n' && src[end] != '\r' {
		end++
	}

	lineStr := src[start:end]

	var caretLine strings.Builder
	for _, r := range src[start:l.pos] {
		if r == '\t' {
			caretLine.WriteString("\t")
		} else {
			caretLine.WriteString(" ")
		}
	}
	caretLine.WriteString("^")

	lineNumStr := fmt.Sprintf("%d", actualLine)
	paddingWidth := len(lineNumStr)

	fmt.Printf("%serror%s: %sunrecognized token%s\n", ColorBoldRed, ColorReset, ColorBold, ColorReset)
	fmt.Printf("  --> %s:%d:%d\n", path, actualLine, actualCol)
	fmt.Printf("%s%*s|%s\n", ColorBoldCyan, paddingWidth+1, "", ColorReset)
	fmt.Printf("%s%d |%s %s\n", ColorBoldCyan, actualLine, ColorReset, lineStr)
	fmt.Printf("%s%*s|%s %s%s%s\n", ColorBoldCyan, paddingWidth+1, "", ColorReset, ColorBoldRed, caretLine.String(), ColorReset)
}
