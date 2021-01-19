package token

import (
	"fmt"
)

type Token struct {
	TokenType Type
	Lexeme    string
	Literal   interface{}
	Line      uint
}

func NewToken(tokenType Type, lexeme string, literal interface{}, line uint) Token {
	return Token{TokenType: tokenType, Lexeme: lexeme, Literal: literal, Line: line}
}

func (t *Token) ToString() string {
	return fmt.Sprintf("%s", Lexeme[t.TokenType])
}
