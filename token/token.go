package token

import (
	"fmt"
	"go/token"
	"strconv"
)

type TokenType int

const (
	LParentheses TokenType = iota
	RParentheses
	LBrace
	RBrace
	LBracket
	RBracket
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Star

	ddig_begin
	Slash
	Not
	Equal
	Greater
	Less
	ddig_end

	NotEqual
	EqualEqual
	GreaterEqual
	LessEqual

	String
	Integer
	Float
	Identifier

	keywords_begin
	And
	Class
	Else
	False
	For
	Function
	If
	Null
	Or
	Print
	Return
	Super
	This
	True
	Var
	Break
	keywords_end

	Illegal
	Error
	Eof
)

var tokens = [...]string{
	LParentheses: "(",
	RParentheses: ")",
	LBrace:       "{",
	RBrace:       "}",
	LBracket:     "[",
	RBracket:     "]",
	Comma:        ",",
	Dot:          ".",
	Minus:        "-",
	Plus:         "+",
	Semicolon:    ";",
	Star:         "*",

	Slash:   "/",
	Not:     "!",
	Equal:   "=",
	Less:    "<",
	Greater: ">",

	LessEqual:    "<=",
	GreaterEqual: ">=",
	NotEqual:     "!=",
	EqualEqual:   "==",


	String:     "STRING",
	Integer:    "INTEGER",
	Float:      "FLOAT",
	Identifier: "IDENTIFIER",

	And:      "AND",
	Class:    "CLASS",
	Else:     "ELSE",
	False:    "FALSE",
	For:      "FOR",
	Function: "FUNCTION",
	If:       "IF",
	Null:     "NULL",
	Or:       "OR",
	Print:    "PRINT",
	Return:   "RETURN",
	Super:    "SUPER",
	This:     "THIS",
	True:     "TRUE",
	Var:      "VAR",
	Break:    "BREAK",

	Eof: "EOF",
}

func (t TokenType) String() string {
	s := ""
	if 0 <= t && t < TokenType(len(tokens)) {
		s = tokens[t]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(t)) + ")"
	}
	return s
}

var keywords map[string]TokenType

func init() {
	keywords = make(map[string]TokenType)
	for i := keywords_begin + 1; i < keywords_end; i++ {
		keywords[tokens[i]] = i
	}
}

func Lookup(ident string) TokenType {
	if tok, isKeyword := keywords[ident]; isKeyword {
		return tok
	}
	return Identifier
}

func IsKeyword(name string) bool {
	_, ok := keywords[name]
	return ok
}

// Precedence

const (
	PrecNone = iota
	PrecAssignment
	PrecOr
	PrecAnd
	PrecEquality
	PrecComparison
	PrecTerm
	PrecFactor
	PrecUnary
	PrecCall
	PrecPrimary
)

type Position struct {
	Filename string
	Offset   int
	Line     int
	Column   int
}

func GoTokenPosToPos(position token.Position) Position {
	return Position{
		Filename: position.Filename,
		Offset:   position.Offset,
		Line:     position.Line,
		Column:   position.Column,
	}
}

func (p *Position) IsValid() bool { return p.Line > 0 }
func (p *Position) String() string {
	s := ""
	if p.IsValid() {
		if s != "" {
			s += ":"
		}
		s += fmt.Sprintf("%4d", p.Line)
		if p.Column != 0 {
			s += fmt.Sprintf(":%-3d", p.Column)
		}
	}
	if s == "" {
		s = "-"
	}
	return s
}
