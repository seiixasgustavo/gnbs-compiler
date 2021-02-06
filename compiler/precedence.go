package compiler

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

func (p *Parser) parsePrecedence(precedence int) {
	p.advance()

	prefixRule := p.getRule(p.previous.TokenType).preffix
	if prefixRule == nil {
		p.error([]byte("Expect expression."))
		return
	}
	prefixRule()

	for precedence <= p.getRule(p.current.TokenType).precedence {
		p.advance()
		infixRule := p.getRule(p.previous.TokenType).infix
		infixRule()
	}
}
