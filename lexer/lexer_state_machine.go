package lexer

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func isAlphaNumeric(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}

func isWhiteSpaces(l *lexer, ch byte) bool {
	if ch == ' ' || ch == '\t' || ch == '\r' {
		l.advanceN(1)
		return true
	}
	return false
}

func isNewLine(l *lexer, ch byte) bool {
	if ch == '\n' {
		l.pos++
		l.line++
		l.col = 1
		return true
	}
	return false
}

func isComments(l *lexer, ch byte) bool {
	if ch == '/' && l.pos+1 < len(l.src) && l.src[l.pos+1] == '/' {
		l.advanceN(2)
		for l.pos < len(l.src) && l.src[l.pos] != '\n' {
			l.advanceN(1)
		}
		return true
	}
	return false
}

func isStringLiteral(l *lexer, ch byte, path string) bool {
	if ch == '"' {
		start := l.pos
		l.advanceN(1)
		for l.pos < len(l.src) && l.src[l.pos] != '"' {
			if l.src[l.pos] == '\n' {
				l.line++
				l.col = 1
				l.pos++
			} else {
				l.advanceN(1)
			}
		}

		if l.pos < len(l.src) && l.src[l.pos] == '"' {
			l.advanceN(1)
			val := l.src[start:l.pos]
			l.push(NewToken(String, val))
			return true
		} else {
			ErrorHandler(*l, l.src, path)
			panic(0)
		}
	}
	return false
}

func isNumberLiteral(l *lexer, ch byte) bool {
	if isDigit(ch) {
		start := l.pos
		for l.pos < len(l.src) && isDigit(l.src[l.pos]) {
			l.advanceN(1)
		}

		if l.pos < len(l.src) && l.src[l.pos] == '.' {
			if l.pos+1 < len(l.src) && isDigit(l.src[l.pos+1]) {
				l.advanceN(1) // consume '.'
				for l.pos < len(l.src) && isDigit(l.src[l.pos]) {
					l.advanceN(1)
				}
			}
		}

		val := l.src[start:l.pos]
		l.push(NewToken(Number, val))
		return true
	}
	return false
}

func isIdentifierOrKeyword(l *lexer, ch byte) bool {
	if isAlpha(ch) || ch == '_' {
		start := l.pos
		for l.pos < len(l.src) && (isAlphaNumeric(l.src[l.pos]) || l.src[l.pos] == '_') {
			l.advanceN(1)
		}

		val := l.src[start:l.pos]
		if t, has := KwdMem[val]; has {
			l.push(NewToken(t, val))
		} else {
			l.push(NewToken(Identifier, val))
		}
		return true
	}
	return false
}

func TokenizeStateMachine(path, src string) []Token {
	l := &lexer{
		src:    src,
		pos:    0,
		line:   1,
		col:    1,
		Tokens: make([]Token, 0),
	}

	for !l.atEOF() {
		ch := l.src[l.pos]

		if isWhiteSpaces(l, ch) ||
			isNewLine(l, ch) ||
			isComments(l, ch) ||
			isStringLiteral(l, ch, path) ||
			isNumberLiteral(l, ch) ||
			isIdentifierOrKeyword(l, ch) {
			continue
		}

		// 6. Operators and Punctuation symbols
		switch ch {
		case '(':
			l.push(NewToken(OpenParen, "("))
			l.advanceN(1)
		case ')':
			l.push(NewToken(CloseParen, ")"))
			l.advanceN(1)
		case '{':
			l.push(NewToken(OpenCurly, "{"))
			l.advanceN(1)
		case '}':
			l.push(NewToken(CloseCurly, "}"))
			l.advanceN(1)
		case '[':
			l.push(NewToken(OpenBracket, "["))
			l.advanceN(1)
		case ']':
			l.push(NewToken(CloseBracket, "]"))
			l.advanceN(1)
		case ';':
			l.push(NewToken(SemiColon, ";"))
			l.advanceN(1)
		case ':':
			l.push(NewToken(Colon, ":"))
			l.advanceN(1)
		case ',':
			l.push(NewToken(Comma, ","))
			l.advanceN(1)
		case '%':
			l.push(NewToken(Percent, "%"))
			l.advanceN(1)
		case '*':
			l.push(NewToken(Star, "*"))
			l.advanceN(1)
		case '/':
			l.push(NewToken(Slash, "/"))
			l.advanceN(1)
		case '=':
			if l.pos+1 < len(l.src) && l.src[l.pos+1] == '=' {
				l.push(NewToken(Equals, "=="))
				l.advanceN(2)
			} else {
				l.push(NewToken(Assignment, "="))
				l.advanceN(1)
			}
		case '!':
			if l.pos+1 < len(l.src) && l.src[l.pos+1] == '=' {
				l.push(NewToken(NotEquals, "!="))
				l.advanceN(2)
			} else {
				l.push(NewToken(Not, "!"))
				l.advanceN(1)
			}
		case '<':
			if l.pos+1 < len(l.src) && l.src[l.pos+1] == '=' {
				l.push(NewToken(LessEquals, "<="))
				l.advanceN(2)
			} else {
				l.push(NewToken(Less, "<"))
				l.advanceN(1)
			}
		case '>':
			if l.pos+1 < len(l.src) && l.src[l.pos+1] == '=' {
				l.push(NewToken(GreaterEquals, ">="))
				l.advanceN(2)
			} else {
				l.push(NewToken(Greater, ">"))
				l.advanceN(1)
			}
		case '|':
			if l.pos+1 < len(l.src) && l.src[l.pos+1] == '|' {
				l.push(NewToken(Or, "||"))
				l.advanceN(2)
			} else {
				ErrorHandler(*l, l.src, path)
				panic(0)
			}
		case '&':
			if l.pos+1 < len(l.src) && l.src[l.pos+1] == '&' {
				l.push(NewToken(And, "&&"))
				l.advanceN(2)
			} else {
				ErrorHandler(*l, l.src, path)
				panic(0)
			}
		case '.':
			if l.pos+1 < len(l.src) && l.src[l.pos+1] == '.' {
				l.push(NewToken(DotDot, ".."))
				l.advanceN(2)
			} else {
				l.push(NewToken(Dot, "."))
				l.advanceN(1)
			}
		case '?':
			if l.pos+2 < len(l.src) && l.src[l.pos+1] == '?' && l.src[l.pos+2] == '=' {
				l.push(NewToken(NullishAssignment, "??="))
				l.advanceN(3)
			} else {
				l.push(NewToken(Question, "?"))
				l.advanceN(1)
			}
		case '+':
			if l.pos+1 < len(l.src) && l.src[l.pos+1] == '+' {
				l.push(NewToken(PlusPlus, "++"))
				l.advanceN(2)
			} else if l.pos+1 < len(l.src) && l.src[l.pos+1] == '=' {
				l.push(NewToken(PlusEquals, "+="))
				l.advanceN(2)
			} else {
				l.push(NewToken(Plus, "+"))
				l.advanceN(1)
			}
		case '-':
			if l.pos+1 < len(l.src) && l.src[l.pos+1] == '-' {
				l.push(NewToken(MinusMinus, "--"))
				l.advanceN(2)
			} else if l.pos+1 < len(l.src) && l.src[l.pos+1] == '=' {
				l.push(NewToken(MinusEquals, "-="))
				l.advanceN(2)
			} else {
				l.push(NewToken(Dash, "-"))
				l.advanceN(1)
			}
		default:
			ErrorHandler(*l, l.src, path)
			panic(0)
		}
	}

	l.push(NewToken(EOF, "EOF"))
	return l.Tokens
}
