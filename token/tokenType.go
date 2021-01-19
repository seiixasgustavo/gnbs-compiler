package token

type Type int

const (
	LParentheses Type = iota
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

	Identifier
	String
	Integer
	Float

	And
	Or
	If
	Else
	For
	Var
	Func
	Struct
	Return
	True
	False
	Null

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

	// Keywords
	And:    "et",
	Or:     "ou",
	If:     "si",
	Else:   "autre",
	For:    "pendant",
	Var:    "var",
	Func:   "fonction",
	Struct: "struct",
	Return: "revenir",
	True:   "vrai",
	False:  "faux",
	Null:   "null",

	Eof: "EOF",
}
