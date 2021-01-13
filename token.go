package GNBS

import (
	"fmt"
)

type Token struct {
	tokenType TokenType
	lexeme    string
	literal   interface{}
	line      uint
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line uint) *Token {
	return &Token{tokenType: tokenType, lexeme: lexeme, literal: literal, line: line}
}

func (t *Token) ToString() string {
	return fmt.Sprintf("%s %s %v", tokens[t.tokenType], t.lexeme, t.literal)
}

type TokenType int

const (
	// Single character tokens
	LParentheses TokenType = iota
	RParentheses
	LBrace
	RBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star

	// One or two character tokens
	Not
	NotEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	// Identifier
	Identifier
	String
	Integer
	Float

	//Keywords
	And
	Or
	If
	Else
	Var
	Func
	Struct
	Return
	True
	False

	Eof
)

var tokens = [...]string{
	// Single character tokens
	LParentheses: "(",
	RParentheses: ")",
	LBrace:       "{",
	RBrace:       "}",
	Comma:        ",",
	Dot:          ".",
	Minus:        "-",
	Plus:         "+",
	Semicolon:    ";",
	Slash:        "/",
	Star:         "*",

	// One or two character tokens
	Not:          "!",
	NotEqual:     "!=",
	Equal:        "=",
	EqualEqual:   "==",
	Greater:      ">",
	GreaterEqual: ">=",
	Less:         "<",
	LessEqual:    "<=",

	// Keywords
	And:    "et",
	Or:     "ou",
	If:     "si",
	Else:   "autre",
	Var:    "var",
	Func:   "fonction",
	Struct: "struct",
	Return: "revenir",
	True:   "vrai",
	False:  "faux",

	Eof: "EOF",
}
