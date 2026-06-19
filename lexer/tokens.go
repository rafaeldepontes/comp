package lexer

import (
	"fmt"
	"slices"
)

type TokenType int

const (
	EOF TokenType = iota
	Null
	True
	False
	Number
	String
	Boolean
	Return
	This
	Struct
	Array
	Impl
	Identifier

	OpenBracket
	CloseBracket
	OpenCurly
	CloseCurly
	OpenParen
	CloseParen

	Assignment

	Equals
	Not
	NotEquals

	Less
	LessEquals
	Greater
	GreaterEquals

	Or
	And

	Dot
	DotDot
	SemiColon
	Colon
	Question
	Comma

	// Inc or Dec Ops
	PlusPlus
	MinusMinus
	PlusEquals
	MinusEquals
	NullishAssignment

	// Math Ops
	Plus
	Dash
	Slash
	Star
	Percent

	// KWD (Key Words)
	Let
	Const
	Class
	New
	Import
	From
	Fn
	If
	Else
	Foreach
	While
	For
	Export
	Typeof
	In
)

var KwdMem = map[string]TokenType{
	"let":     Let,
	"const":   Const,
	"class":   Class,
	"new":     New,
	"import":  Import,
	"from":    From,
	"fn":      Fn,
	"if":      If,
	"else":    Else,
	"foreach": Foreach,
	"while":   While,
	"for":     For,
	"export":  Export,
	"typeof":  Typeof,
	"in":      In,
	"return":  Return,
	"struct":  Struct,
	"this":    This,
	"impl":    Impl,
}

type Token struct {
	Val  string
	Type TokenType
}

func NewToken(type_ TokenType, val string) Token {
	return Token{
		Val:  val,
		Type: type_,
	}
}

func (t Token) isOneOfMany(exp ...TokenType) bool {
	return slices.Contains(exp, t.Type)
}

func (t Token) Debbug() {
	if t.isOneOfMany(Identifier, String, Number) {
		fmt.Printf("type: %s, \t\t\t      (%s)\n", TokenTypeString(t.Type), t.Val)
	} else {
		fmt.Printf("type: %s, \t\t\tValue: %s \n", TokenTypeString(t.Type), t.Val)
	}
}

func TokenTypeString(tt TokenType) string {
	switch tt {
	case EOF:
		return "eof"
	case Null:
		return "null"
	case Number:
		return "number"
	case String:
		return "string"
	case True:
		return "true"
	case False:
		return "false"
	case Identifier:
		return "identifier"
	case OpenBracket:
		return "openBracket"
	case CloseBracket:
		return "closeBracket"
	case OpenCurly:
		return "openCurly"
	case CloseCurly:
		return "closeCurly"
	case OpenParen:
		return "openParen"
	case CloseParen:
		return "closeParen"
	case Assignment:
		return "assignment"
	case Equals:
		return "equals"
	case NotEquals:
		return "not_equals"
	case Not:
		return "not"
	case Less:
		return "less"
	case LessEquals:
		return "lessEquals"
	case Greater:
		return "greater"
	case GreaterEquals:
		return "greaterEquals"
	case Or:
		return "or"
	case And:
		return "and"
	case Dot:
		return "dot"
	case DotDot:
		return "dotDot"
	case SemiColon:
		return "semiColon"
	case Colon:
		return "colon"
	case Question:
		return "question"
	case Comma:
		return "comma"
	case PlusPlus:
		return "plusPlus"
	case MinusMinus:
		return "minusMinus"
	case PlusEquals:
		return "plusEquals"
	case MinusEquals:
		return "minusEquals"
	case NullishAssignment:
		return "nullishAssignment"
	case Plus:
		return "plus"
	case Dash:
		return "dash"
	case Slash:
		return "slash"
	case Star:
		return "star"
	case Percent:
		return "percent"
	case Let:
		return "let"
	case Const:
		return "const"
	case Class:
		return "class"
	case New:
		return "new"
	case Import:
		return "import"
	case From:
		return "from"
	case Fn:
		return "fn"
	case If:
		return "if"
	case Else:
		return "else"
	case Foreach:
		return "foreach"
	case For:
		return "for"
	case While:
		return "while"
	case Export:
		return "export"
	case In:
		return "in"
	case Return:
		return "return"
	case This:
		return "this"
	case Impl:
		return "impl"
	case Struct:
		return "struct"
	case Boolean:
		return "boolean"
	default:
		return fmt.Sprintf("unknown(%d)", tt)
	}
}
