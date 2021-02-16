package compiler

import (
	"GNBS/chunk"
	"GNBS/scanner"
	"GNBS/token"
	"fmt"
	"math"
	"os"
	"strconv"
)

type Parser struct {
	current  *scanner.Token
	previous *scanner.Token
	scanner  *scanner.Scanner

	hadError  bool
	panicMode bool
}

var (
	parser         *Parser
	compilingChunk *chunk.Chunk
)

func InitParser() {

}

// Parser

func GetParser() *Parser {
	return parser
}

func Compile(source []byte, chunk *chunk.Chunk) bool {
	parser.scanner = scanner.NewScanner(source, nil)
	compilingChunk = chunk

	advance()
	expression()
	consume(token.Eof, "Expect end of expression.")
	endCompiler()
	return parser.hadError
}

func advance() {
	parser.previous = parser.current

	for {
		parser.current = parser.scanner.Scan()
		if parser.current.Token != token.Error {
			break
		}

		errorAtCurrent(parser.current.LitName)
	}
}

func consume(tp token.TokenType, message string) {
	if parser.current.Token == tp {
		advance()
		return
	}

	errorAtCurrent(message)
}

func expression() {
	parsePrecedence(Assignment)
}

func grouping() {
	expression()
	consume(token.RParentheses, "Expect ')' after expression.")
}

func unary() {
	operatorType := parser.previous.Token

	parsePrecedence(Unary)

	switch operatorType {
	case token.Not:
		emitByte(chunk.OpNot)
		break
	case token.Minus:
		emitByte(chunk.OpNegate)
		break
	default:
		return
	}
}

func binary() {
	operatorType := parser.previous.Token
	rule := getRule(operatorType)
	parsePrecedence(rule.precedence + 1)

	switch operatorType {
	case token.Plus:
		emitByte(chunk.OpAdd)
		break
	case token.Minus:
		emitByte(chunk.OpSubtract)
		break
	case token.Star:
		emitByte(chunk.OpMultiply)
		break
	case token.Slash:
		emitByte(chunk.OpDivide)
		break
	case token.NotEqual:
		emitBytes(chunk.OpEqual, chunk.OpNot)
		break
	case token.EqualEqual:
		emitByte(chunk.OpEqual)
		break
	case token.Greater:
		emitByte(chunk.OpGreater)
		break
	case token.GreaterEqual:
		emitBytes(chunk.OpLess, chunk.OpNot)
		break
	case token.Less:
		emitByte(chunk.OpLess)
		break
	case token.LessEqual:
		emitBytes(chunk.OpGreater, chunk.OpNot)
	default:
		break
	}
}

func literal() {
	switch parser.previous.Token {
	case token.False:
		emitByte(chunk.OpFalse)
		break
	case token.Null:
		emitByte(chunk.OpNull)
		break
	case token.True:
		emitByte(chunk.OpTrue)
		break
	default:
		return
	}
}

func number() {
	value, _ := strconv.ParseFloat(parser.previous.LitName, 64)
	valueDS := chunk.Value{
		Type:  chunk.TypeFloat,
		Value: value,
	}
	emitConstant(valueDS)
}

func makeConstant(value chunk.Value) byte {
	constant := currentChunk().AddConstant(value)
	if constant > math.MaxUint8 {
		error("Too many constants in one chunk.")
		return 0
	}
	return constant
}

func endCompiler() {
	emitReturn()
}

// Emit Bytes

func emitByte(by byte) {
	currentChunk().WriteChunk(by, *parser.scanner.GetPosition(parser.previous.Position))
}

func emitBytes(by, by2 byte) {
	emitByte(by)
	emitByte(by2)
}

func emitConstant(value chunk.Value) {
	emitBytes(chunk.OpConstant, makeConstant(value))
}

func emitReturn() {
	emitByte(chunk.OpReturn)
}

// Compiling Chunk

func currentChunk() *chunk.Chunk {
	return compilingChunk
}

// Precedence

func parsePrecedence(precedence Precedence) {
	advance()
	prefixRule := getRule(parser.previous.Token).prefix
	if prefixRule == nil {
		error("Expect expression.")
		return
	}

	prefixRule()

	for precedence <= getRule(parser.current.Token).precedence {
		advance()
		infixRule := getRule(parser.previous.Token).infix

		infixRule()
	}
}

// Error Handlers

func errorAt(tk *scanner.Token, message string) {
	if parser.panicMode {
		return
	}
	parser.panicMode = true

	pos := parser.scanner.GetPosition(tk.Position)
	fmt.Fprintf(os.Stderr, "[line %d:%d] Error", pos.Line, pos.Column)

	if tk.Token == token.Eof {
		fmt.Fprintf(os.Stderr, " at end")
	} else if tk.Token == token.Error {

	} else {
		fmt.Fprintf(os.Stderr, " at '%s'", tk.LitName)
	}

	fmt.Fprintf(os.Stderr, ": %s\n", message)
	parser.hadError = true
}

func error(message string) {
	errorAt(parser.previous, message)
}

func errorAtCurrent(message string) {
	errorAt(parser.current, message)
}
