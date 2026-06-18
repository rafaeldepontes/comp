package lexer

import (
	"fmt"
	"regexp"
)

type regexHandler func(l *lexer, reg *regexp.Regexp)

type regexpPattern struct {
	reg *regexp.Regexp
	h   regexHandler
}

type lexer struct {
	patterns []regexpPattern
	Tokens   []Token
	src      string
	pos      int
	line     int
}

func (l *lexer) advanceN(n int) {
	l.pos += n
}

func (l *lexer) push(t Token) {
	l.Tokens = append(l.Tokens, t)
}

func (l *lexer) remainder() string {
	return l.src[l.pos:]
}

func (l *lexer) atEOF() bool {
	return l.pos >= len(l.src)
}

func defaultHandler(t TokenType, val string) regexHandler {
	return func(l *lexer, reg *regexp.Regexp) {
		l.push(NewToken(t, val))
		l.advanceN(len(val))
	}
}

func numberHandler(l *lexer, reg *regexp.Regexp) {
	match := reg.FindString(l.remainder())
	l.push(NewToken(Number, match))
	l.advanceN(len(match))
}

func stringHandler(l *lexer, reg *regexp.Regexp) {
	match := reg.FindStringIndex(l.remainder())
	sl := l.remainder()[match[0]:match[1]]

	l.push(NewToken(String, sl))
	l.advanceN(len(sl))
}

func symbolHandler(l *lexer, reg *regexp.Regexp) {
	match := reg.FindString(l.remainder())

	if t, has := KwdMem[match]; has {
		l.push(NewToken(t, match))
	} else {
		l.push(NewToken(Identifier, match))
	}

	l.advanceN(len(match))
}

func skipHandler(l *lexer, reg *regexp.Regexp) {
	match := reg.FindStringIndex(l.remainder())
	l.advanceN(match[1])
}

func commentHandler(l *lexer, reg *regexp.Regexp) {
	match := reg.FindStringIndex(l.remainder())

	if match != nil {
		l.advanceN(match[1])
		l.line++
	}
}

func newLexer(src string) *lexer {
	return &lexer{
		pos:    0,
		line:   1,
		src:    src,
		Tokens: make([]Token, 0),
		patterns: []regexpPattern{
			{
				reg: regexp.MustCompile(`\s+`),
				h:   skipHandler,
			},
			{
				reg: regexp.MustCompile(`\/\/.*`),
				h:   commentHandler,
			},
			{
				reg: regexp.MustCompile(`"[^"]*"`),
				h:   stringHandler,
			},
			{
				reg: regexp.MustCompile(`[0-9]+(\.[0-9]+)?`),
				h:   numberHandler,
			},
			{
				reg: regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`),
				h:   symbolHandler,
			},
			{
				reg: regexp.MustCompile(`\(`),
				h:   defaultHandler(OpenParen, "("),
			},
			{
				reg: regexp.MustCompile(`\)`),
				h:   defaultHandler(CloseParen, ")"),
			},
			{
				reg: regexp.MustCompile(`\{`),
				h:   defaultHandler(OpenCurly, "{"),
			},
			{
				reg: regexp.MustCompile(`\}`),
				h:   defaultHandler(CloseCurly, "}"),
			},
			{
				reg: regexp.MustCompile(`==`),
				h:   defaultHandler(Equals, "=="),
			},
			{
				reg: regexp.MustCompile(`=`),
				h:   defaultHandler(Assignment, "="),
			},
			{
				reg: regexp.MustCompile(`!=`),
				h:   defaultHandler(NotEquals, "!="),
			},
			{
				reg: regexp.MustCompile(`!`),
				h:   defaultHandler(Not, "!"),
			},
			{
				reg: regexp.MustCompile(`<=`),
				h:   defaultHandler(LessEquals, "<="),
			},
			{
				reg: regexp.MustCompile(`<`),
				h:   defaultHandler(Less, "<"),
			},
			{
				reg: regexp.MustCompile(`>=`),
				h:   defaultHandler(GreaterEquals, ">="),
			},
			{
				reg: regexp.MustCompile(`>`),
				h:   defaultHandler(Greater, ">"),
			},
			{
				reg: regexp.MustCompile(`\|\|`),
				h:   defaultHandler(Or, "||"),
			},
			{
				reg: regexp.MustCompile(`&&`),
				h:   defaultHandler(And, "&&"),
			},
			{
				reg: regexp.MustCompile(`\.\.`),
				h:   defaultHandler(DotDot, ".."),
			},
			{
				reg: regexp.MustCompile(`\.`),
				h:   defaultHandler(Dot, "."),
			},
			{
				reg: regexp.MustCompile(`\[`),
				h:   defaultHandler(OpenBracket, "["),
			},
			{
				reg: regexp.MustCompile(`\]`),
				h:   defaultHandler(CloseBracket, "]"),
			},
			{
				reg: regexp.MustCompile(`;`),
				h:   defaultHandler(SemiColon, ";"),
			},
			{
				reg: regexp.MustCompile(`:`),
				h:   defaultHandler(Colon, ":"),
			},
			{
				reg: regexp.MustCompile(`\?\?=`),
				h:   defaultHandler(NullishAssignment, "??="),
			},
			{
				reg: regexp.MustCompile(`\?`),
				h:   defaultHandler(Question, "?"),
			},
			{
				reg: regexp.MustCompile(`,`),
				h:   defaultHandler(Comma, ","),
			},
			{
				reg: regexp.MustCompile(`\+\+`),
				h:   defaultHandler(PlusPlus, "++"),
			},
			{
				reg: regexp.MustCompile(`--`),
				h:   defaultHandler(MinusMinus, "--"),
			},
			{
				reg: regexp.MustCompile(`\+=`),
				h:   defaultHandler(PlusEquals, "+="),
			},
			{
				reg: regexp.MustCompile(`-=`),
				h:   defaultHandler(MinusEquals, "-="),
			},
			{
				reg: regexp.MustCompile(`\+`),
				h:   defaultHandler(Plus, "+"),
			},
			{
				reg: regexp.MustCompile(`-`),
				h:   defaultHandler(Dash, "-"),
			},
			{
				reg: regexp.MustCompile(`/`),
				h:   defaultHandler(Slash, "/"),
			},
			{
				reg: regexp.MustCompile(`\*`),
				h:   defaultHandler(Star, "*"),
			},
			{
				reg: regexp.MustCompile(`%`),
				h:   defaultHandler(Percent, "%"),
			},
		},
	}
}

func Tokenize(src string) []Token {
	l := newLexer(src)

	for !l.atEOF() {
		matched := false

		for i := range l.patterns {
			loc := l.patterns[i].reg.FindStringIndex(l.remainder())

			if loc != nil && loc[0] == 0 {
				l.patterns[i].h(l, l.patterns[i].reg)
				matched = true
				break
			}
		}

		// TODO: improve this to show the location where the error occurred
		if !matched {
			test := l.remainder()
			println(test)
			panic(fmt.Sprintf("[ERROR] unrecognized token near %s\n", l.remainder()))
		}
	}

	l.push(NewToken(EOF, "EOF"))
	return l.Tokens
}
