package parser

import (
	"GNBS/token"
	"testing"
)

func TestPrint(t *testing.T) {
	var astPrinter AstPrinter

	expr := binary{
		left: unary{
			operator: token.NewToken(token.Minus, token.Lexeme[token.Minus], nil, 1),
			right:    literal{value: 123},
		},
		operator: token.NewToken(token.Star, token.Lexeme[token.Star], nil, 1),
		right:    grouping{expressions: literal{value: 45.67}},
	}

	if astPrinter.print(expr) != `(* (- 123) (group 45.67))` {
		t.Error("string didn't match")
	}
}

func TestPrint2(t *testing.T) {
	var astPrinter AstPrinter

	expr := binary{
		left:     grouping{expressions: literal{value: "asd"}},
		operator: token.NewToken(token.Star, token.Lexeme[token.Star], nil, 1),
		right:    grouping{expressions: literal{value: 11}},
	}

	if astPrinter.print(expr) != `(* (group "asd") (group 11))` {
		t.Error("string didn't match")
	}
}
