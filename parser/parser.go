package parser

import "GNBS/token"

type Parser struct {
	tokens  []token.Type
	current uint
}
