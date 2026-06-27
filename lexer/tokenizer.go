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

func (l *Lexer) isWhiteSpaces(ch byte) bool {
	if ch == ' ' || ch == '\t' || ch == '\r' {
		l.advance()
		return true
	}
	return false
}

func (l *Lexer) isNewLine(ch byte) bool {
	if ch == '\n' {
		l.advance()
		l.lin++
		l.col = 1
		return true
	}
	return false
}

func (l *Lexer) isExit() bool {
	if isAlpha(l.src[l.pos]) || l.src[l.pos] == '_' {
		start := l.pos

		for l.hasData() && (isAlphaNumeric(l.src[l.pos]) || l.src[l.pos] == '_') {
			l.advance()
		}

		val := l.src[start:l.pos]
		if val, has := kwd[val]; has {
			l.Tokens = append(l.Tokens, newToken("", val))
		} else {
			l.Tokens = append(l.Tokens, newToken(l.src[start:l.pos], Unknown))
		}

		return true
	}
	return false
}

func (l *Lexer) isInteger() bool {
	if isDigit(l.src[l.pos]) {
		start := l.pos

		for l.hasData() && isDigit(l.src[l.pos]) {
			l.advance()
		}

		if l.hasData() && l.src[l.pos] == '.' {
			if l.pos+1 < len(l.src) && isDigit(l.src[l.pos+1]) {
				l.advance() // skips the '.'

				for l.hasData() && isDigit(l.src[l.pos]) {
					l.advance()
				}
			}
		}
		l.Tokens = append(l.Tokens, newToken(l.src[start:l.pos], Int))
		return true
	}
	return false
}
