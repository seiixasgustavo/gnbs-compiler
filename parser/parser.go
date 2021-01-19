package parser

import (
	"GNBS"
	"GNBS/token"
	"errors"
	"fmt"
)

type Parser struct {
	tokens  []token.Token
	current uint
}

func (p *Parser) parse() (Expr, error) {
	return p.expression()
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) check(tp token.Type) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == tp
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == token.Eof
}

func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(token.NotEqual, token.EqualEqual) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = binary{
			left:     expr,
			operator: operator,
			right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) match(types ...token.Type) bool {
	for _, tk := range types {
		if p.check(tk) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(token.Greater, token.GreaterEqual, token.Less, token.LessEqual) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = binary{expr, operator, right}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}
	for p.match(token.Minus, token.Plus) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = binary{expr, operator, right}
	}
	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.Slash, token.Star) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = binary{expr, operator, right}
	}
	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	for p.match(token.Not, token.Minus) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return unary{
			operator: operator,
			right:    right,
		}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(token.False) {
		return literal{false}, nil
	}
	if p.match(token.True) {
		return literal{true}, nil
	}
	if p.match(token.Null) {
		return literal{nil}, nil
	}

	if p.match(token.Integer, token.Float, token.String) {
		return literal{p.previous().Literal}, nil
	}

	if p.match(token.LParentheses) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.RParentheses, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return grouping{expr}, nil
	}
	return nil, errors.New("expect expression")
}

func (p *Parser) consume(tokenType token.Type, message string) (token.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	p.error(p.peek(), message)
	return token.Token{}, errors.New(fmt.Sprintf("%s %s", p.peek(), message))
}

func (p *Parser) error(tk token.Token, message string) {
	GNBS.Compiler.Error(tk.Line, message)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == token.Semicolon {
			return
		}
		switch p.peek().TokenType {
		case token.Struct:
		case token.Func:
		case token.Var:
		case token.For:
		case token.If:
		case token.Return:
		}
		p.advance()
	}
}
