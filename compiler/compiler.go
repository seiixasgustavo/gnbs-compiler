package compiler

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

func Compile(source []byte, chunk *Chunk) bool {
	parser := NewParser(NewScanner(source), false, false, chunk)
	parser.advance()
	parser.consume(Eof, []byte("Expect end of expression."))
	return !parser.HadError
}

type Parser struct {
	Scanner        *Scanner
	current        *Token
	previous       *Token
	HadError       bool
	PanicMode      bool
	CompilingChunk *Chunk
	rules          []ParseRule
}

func NewParser(scanner *Scanner, hadError bool, panicMode bool, compilingChunk *Chunk) *Parser {
	parser := &Parser{Scanner: scanner, HadError: hadError, PanicMode: panicMode, CompilingChunk: compilingChunk}
	parser.rules = parser.initRules()
	return parser
}

func (p *Parser) advance() {
	p.previous = p.current
	for {
		p.current = p.Scanner.ScanToken()
		if p.current.TokenType != Error {
			break
		}
		p.errorAtCurrent(p.current.source)
	}
}

func (p *Parser) consume(tokenType int, message []byte) {
	if p.current.TokenType == tokenType {
		p.advance()
		return
	}
	p.errorAtCurrent(message)
}

func (p *Parser) expression() {
	p.parsePrecedence(PrecAssignment)
}

// Emit Functions

func (p *Parser) emitByte(by byte) {
	p.CompilingChunk.WriteChunk(by, p.previous.line)
}

func (p *Parser) emitBytes(by1, by2 byte) {
	p.emitByte(by1)
	p.emitByte(by2)
}

func (p *Parser) emitReturn() {
	p.emitByte(OpReturn)
}

func (p *Parser) emitConstant(value Value) {
	p.emitBytes(OpConstant, p.makeConstant(value))
}

// Operation Functions

func (p *Parser) number() {
	uvl := binary.LittleEndian.Uint64(p.previous.source)
	value := math.Float64frombits(uvl)
	p.emitConstant(Value(value))
}

func (p *Parser) grouping() {
	p.expression()
	p.consume(RParentheses, []byte("Expect ')' after expression."))
}

func (p *Parser) unary() {
	operatorType := p.previous.TokenType

	p.expression()

	p.parsePrecedence(PrecUnary)

	switch operatorType {
	case Minus:
		p.emitByte(OpNegate)
		break
	default:
		return
	}
}

func (p *Parser) binary() {
	operatorType := p.previous.TokenType
	rule := p.getRule(operatorType)

	p.parsePrecedence(rule.precedence + 1)

	switch operatorType {
	case Plus:
		p.emitByte(OpAdd)
		break
	case Minus:
		p.emitByte(OpSubtract)
		break
	case Star:
		p.emitByte(OpMultiply)
		break
	case Slash:
		p.emitByte(OpDivide)
		break
	default:
		return
	}
}

// Util Functions

func (p *Parser) makeConstant(value Value) byte {
	constant := p.CompilingChunk.Constants.AddConstant(value)
	if constant > math.MaxUint8 {
		p.error([]byte("Too many constants in one chunk."))
		return 0
	}
	return constant
}

// Error Handling Functions

func (p *Parser) errorAtCurrent(message []byte) {
	p.errorAt(p.current, message)
}

func (p *Parser) error(message []byte) {
	p.errorAt(p.previous, message)
}

func (p *Parser) errorAt(token *Token, message []byte) {
	if p.PanicMode {
		return
	}
	p.PanicMode = true

	fmt.Fprintf(os.Stderr, "[line %d] Error", token.line)

	if token.TokenType == Eof {
		fmt.Fprintf(os.Stderr, " at the end")
	} else if token.TokenType == Error {

	} else {
		fmt.Fprintf(os.Stderr, " at '%s'", token.source)
	}

	fmt.Fprintf(os.Stderr, ": %s\n", message)
}

func (p *Parser) endCompiler() {
	p.emitReturn()
	if !p.HadError {
		p.CompilingChunk.DisassembleChunk("code")
	}
}
