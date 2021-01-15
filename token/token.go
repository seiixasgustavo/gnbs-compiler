package token

import (
	"fmt"
)

type Token struct {
	tokenType Type
	lexeme    string
	literal   interface{}
	line      uint
}

func NewToken(tokenType Type, lexeme string, literal interface{}, line uint) Token {
	return Token{tokenType: tokenType, lexeme: lexeme, literal: literal, line: line}
}

func (t *Token) ToString() string {
	return fmt.Sprintf("%s", Lexeme[t.tokenType])
}
