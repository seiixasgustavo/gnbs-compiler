package compiler

const (
	LParentheses = iota
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

	Not
	NotEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	String
	Number
	Identifier

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
	While

	Error
	Eof
)

var Lexeme = [...]string{
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
}
