package scanner

import (
	"strconv"
	"unicode"
)

type TokenType int

const (
	LParentheses TokenType = iota
	RParentheses
	LBrace
	RBrace
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

func IsIdentifier(name string) bool {
	for i, c := range name {
		if !unicode.IsLetter(c) && c != '_' && (i == 0 || !unicode.IsDigit(c)) {
			return false
		}
	}
	return name != "" && !IsKeyword(name)
}
func IsPossibleDoubleDigit(name string) bool {
	tk, ok := keywords[name]
	if !ok {
		return false
	}
	if tk > ddig_begin && tk < ddig_end {
		return true
	}
	return false
}
