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

type Compiler struct {
	Enclosing *Compiler
	Function  *chunk.GFunction
	Type      FunctionType

	Locals     []Local
	LocalCount int
	ScoreDepth int
}

type Local struct {
	Name  *scanner.Token
	Depth int
}

var (
	parser  *Parser
	current *Compiler = nil
)

func InitParser() {

}

func InitCompiler(comp *Compiler, Type FunctionType) {
	comp = &Compiler{Enclosing: current, Function: nil, Type: Type, LocalCount: 0, ScoreDepth: 0, Locals: make([]Local, math.MaxUint8+1)}
	comp.Function = chunk.NewGFunction()

	current = comp

	if Type == TypeScript {
		current.Function.Name = chunk.NewGString(parser.previous.LitName)
	}

	local := &current.Locals[current.LocalCount]
	current.LocalCount++
	local.Depth = 0
	local.Name.LitName = ""
}

// Parser

func GetParser() *Parser {
	return parser
}

func Compile(source []byte) *chunk.GFunction {
	parser.scanner = scanner.NewScanner(source, nil)
	parser.hadError = false

	var compiler Compiler
	InitCompiler(&compiler, TypeScript)

	advance()

	for !match(token.Eof) {
		declaration()
	}

	function := endCompiler()
	if parser.hadError {
		return nil
	}
	return function
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

func declaration() {
	if match(token.Function) {
		funcDeclaration()
	} else if match(token.Var) {
		varDeclaration()
	} else {
		statement()
	}

	if parser.panicMode {
		synchronize()
	}
}

func synchronize() {
	parser.panicMode = false

	for parser.current.Token != token.Eof {
		if parser.previous.Token != token.Semicolon {
			return
		}

		switch parser.current.Token {
		case token.Class, token.Function, token.Var, token.For, token.If, token.Print, token.Return:
			return
		default:

		}
		advance()
	}
}

func match(tk token.TokenType) bool {
	if !check(tk) {
		return false
	}
	advance()
	return true
}

func check(tk token.TokenType) bool {
	return parser.current.Token == tk
}

func beginScope() {
	current.ScoreDepth++
}

func endScope() {
	current.ScoreDepth--

	for current.LocalCount > 0 && current.Locals[current.LocalCount-1].Depth > current.ScoreDepth {
		emitByte(OpPop)
		current.LocalCount--
	}
}

func endCompiler() *chunk.GFunction {
	emitReturn()
	fn := current.Function

	current = current.Enclosing
	return fn
}

// Statement Handlers

func statement() {
	if match(token.Print) {
		printStatement()
	} else if match(token.For) {
		forStatement()
	} else if match(token.If) {
		ifStatement()
	} else if match(token.Return) {
		returnStatement()
	}else {
		expressionStatement()
	}
}

func printStatement() {
	expression()
	consume(token.Semicolon, "Expect ';' after value.")
	emitByte(OpPrint)
}

func expressionStatement() {
	expression()
	consume(token.Semicolon, "Expect ';' after expression.")
	emitByte(OpPop)
}

func ifStatement() {
	consume(token.LParentheses, "Expect '(' after 'if'.")
	expression()
	consume(token.RParentheses, "Expect ')' after condition.")

	thenJump := emitJump(OpJumpIfFalse)
	emitByte(OpPop)
	statement()

	elseJump := emitJump(OpJump)

	patchJump(thenJump)
	emitByte(OpPop)

	if match(token.Else) {
		statement()
	}
	patchJump(elseJump)
}

func forStatement() {
	beginScope()

	consume(token.LParentheses, "Expect '(' after 'for'.")
	if match(token.Semicolon) {

	} else if match(token.Var) {
		varDeclaration()
	} else {
		expressionStatement()
	}

	loopStart := uint16(len(currentChunk().Code))

	exitJump := -1
	if !match(token.Semicolon) {
		expression()
		consume(token.Semicolon, "Expect ';' after loop condition.")

		exitJump = emitJump(OpJumpIfFalse)
		emitByte(OpPop)
	}

	consume(token.RParentheses, "Expect ')' after for clauses.")

	if !match(token.RParentheses) {
		bodyJump := emitJump(OpJump)

		incrementStart := len(currentChunk().Code)
		expression()
		emitByte(OpPop)
		consume(token.RParentheses, "Expect ')' after for clauses.")

		emitLoop(loopStart)
		loopStart = uint16(incrementStart)
		patchJump(bodyJump)
	}

	statement()
	emitLoop(loopStart)

	if exitJump != -1 {
		patchJump(exitJump)
		emitByte(OpPop)
	}

	endScope()
}

func returnStatement() {
	if current.Type == TypeScript {
		error("Can't return from top-level code.")
	}

	if match(token.Semicolon) {
		emitReturn()
	} else {
		expression()
		consume(token.Semicolon, "Expect ';' after return value.")
		emitByte(OpReturn)
	}
}

// Expression handlers

func expression() {
	parsePrecedence(Assignment)
}

func block() {
	for !check(token.RBrace) && !check(token.Eof) {
		declaration()
	}

	consume(token.RBrace, "Expect '}' after block.")
}

func function(functionType FunctionType) {
	var compiler Compiler
	InitCompiler(&compiler, functionType)
	beginScope()

	consume(token.LParentheses, "Expect '(' after function name.")
	if !check(token.RParentheses) {
		firstRun := true
		for firstRun || match(token.Comma) {
			firstRun = false

			current.Function.Arity++
			if current.Function.Arity > 255 {
				errorAtCurrent("Can't have more than 255 parameters")
			}

			paramConstant := parseVariable("Expect parameter name.")
			defineVariable(paramConstant)
		}
	}
	consume(token.RParentheses, "Expect ')' after parameters.")

	consume(token.LBrace, "Expect '{' before function body.")

	block()

	fun := endCompiler()
	emitBytes(OpConstant, makeConstant(chunk.Value{
		Type:  chunk.TypeFunction,
		Value: fun,
	}))
}

func varDeclaration() {
	global := parseVariable("Expect variable name.")

	if match(token.Equal) {
		expression()
	} else {
		emitByte(OpNull)
	}
	consume(token.Semicolon, "Expect ';' after variable declaration.")

	defineVariable(global)
}

func funcDeclaration() {
	global := parseVariable("Expect function name.")
	markInitialized()
	function(TypeFunction)
	defineVariable(global)
}

func grouping(canAssign bool) {
	expression()
	consume(token.RParentheses, "Expect ')' after expression.")
}

func unary(canAssign bool) {
	operatorType := parser.previous.Token

	parsePrecedence(Unary)

	switch operatorType {
	case token.Not:
		emitByte(OpNot)
		break
	case token.Minus:
		emitByte(OpNegate)
		break
	default:
		return
	}
}

func binary(canAssign bool) {
	operatorType := parser.previous.Token
	rule := getRule(operatorType)
	parsePrecedence(rule.precedence + 1)

	switch operatorType {
	case token.Plus:
		emitByte(OpAdd)
		break
	case token.Minus:
		emitByte(OpSubtract)
		break
	case token.Star:
		emitByte(OpMultiply)
		break
	case token.Slash:
		emitByte(OpDivide)
		break
	case token.NotEqual:
		emitBytes(OpEqual, OpNot)
		break
	case token.EqualEqual:
		emitByte(OpEqual)
		break
	case token.Greater:
		emitByte(OpGreater)
		break
	case token.GreaterEqual:
		emitBytes(OpLess, OpNot)
		break
	case token.Less:
		emitByte(OpLess)
		break
	case token.LessEqual:
		emitBytes(OpGreater, OpNot)
	default:
		break
	}
}

func literal(canAssign bool) {
	switch parser.previous.Token {
	case token.False:
		emitByte(OpFalse)
		break
	case token.Null:
		emitByte(OpNull)
		break
	case token.True:
		emitByte(OpTrue)
		break
	default:
		return
	}
}

func callFn(canAssign bool) {
	argCount := argumentList()
	emitBytes(OpCall, argCount)
}

// Values functions

func floatnumber(canAssign bool) {
	value, _ := strconv.ParseFloat(parser.previous.LitName, 64)
	valueDS := chunk.Value{
		Type:  chunk.TypeFloat,
		Value: value,
	}
	emitConstant(valueDS)
}

func intnumber(canAssign bool) {
	value, _ := strconv.ParseInt(parser.previous.LitName, 10, 64)
	valueDS := chunk.Value{
		Type:  chunk.TypeInteger,
		Value: value,
	}
	emitConstant(valueDS)
}

func stringvalue(canAssign bool) {
	emitConstant(chunk.Value{
		Type:  chunk.TypeString,
		Value: chunk.NewGString(parser.previous.LitName),
	})
}

func variable(canAssign bool) {
	namedVariable(parser.previous, canAssign)
}

func makeConstant(value chunk.Value) byte {
	constant := currentChunk().AddConstant(value)
	if constant > math.MaxUint8 {
		error("Too many constants in one chunk.")
		return 0
	}
	return constant
}

func namedVariable(tk *scanner.Token, canAssign bool) {
	arg := resolveLocal(current, tk)
	var getOp, setOp byte

	if arg != -1 {
		getOp = OpGetLocal
		setOp = OpSetLocal
	} else {
		arg = int(identifierConstant(tk))
		getOp = OpGetGlobal
		setOp = OpSetGlobal
	}

	if canAssign && match(token.Equal) {
		expression()
		emitBytes(setOp, byte(arg))
	} else {
		emitBytes(getOp, byte(arg))
	}
}

func and_(canAssign bool) {
	endJump := emitJump(OpJumpIfFalse)
	emitByte(OpPop)
	parsePrecedence(And)

	patchJump(endJump)
}

func or_(canAssign bool) {
	elseJump := emitJump(OpJumpIfFalse)
	endJump := emitJump(OpJump)

	patchJump(elseJump)
	emitByte(OpPop)

	parsePrecedence(Or)
	patchJump(endJump)
}

// Emit Bytes

func emitByte(by byte) {
	currentChunk().WriteChunk(by, *parser.scanner.GetPosition(parser.previous.Position))
}

func emitBytes(by, by2 byte) {
	emitByte(by)
	emitByte(by2)
}

func emitJump(instruction byte) int {
	emitByte(instruction)
	emitByte(0xff)
	emitByte(0xff)
	return len(currentChunk().Code) - 2
}

func emitConstant(value chunk.Value) {
	emitBytes(OpConstant, makeConstant(value))
}

func emitLoop(loopStart uint16) {
	emitByte(OpLoop)

	offset := uint16(len(currentChunk().Code)) - loopStart + 2
	if offset > math.MaxUint16 {
		error("Loop body too large.")
	}

	emitByte(byte((offset >> 8) & 0xff))
	emitByte(byte(offset & 0xff))
}

func patchJump(offset int) {
	jump := len(currentChunk().Code) - offset - 2
	if jump > math.MaxUint16 {
		error("Too much code to jump over.")
	}

	currentChunk().Code[offset] = byte((jump >> 8) & 0xff)
	currentChunk().Code[offset+1] = byte(jump & 0xff)
}

func emitReturn() {
	emitByte(OpNull)
	emitByte(OpReturn)
}

// Compiling Chunk

func currentChunk() *chunk.Chunk {
	return &current.Function.Chunk
}

// Precedence

func parsePrecedence(precedence Precedence) {
	advance()
	prefixRule := getRule(parser.previous.Token).prefix
	if prefixRule == nil {
		error("Expect expression.")
		return
	}

	canAssign := precedence <= token.PrecAssignment
	prefixRule(canAssign)

	for precedence <= getRule(parser.current.Token).precedence {
		advance()
		infixRule := getRule(parser.previous.Token).infix

		infixRule(canAssign)
	}

	if canAssign && match(token.Equal) {
		error("Invalid assignment target.")
	}
}

func identifierConstant(tk *scanner.Token) byte {
	return makeConstant(chunk.Value{
		Type:  chunk.TypeString,
		Value: chunk.NewGString(tk.LitName),
	})
}

func identifiersEqual(a, b *scanner.Token) bool {
	return a.LitName == b.LitName
}

func resolveLocal(compiler *Compiler, tk *scanner.Token) int {
	for i := compiler.LocalCount - 1; i >= 0; i-- {
		local := &compiler.Locals[i]
		if identifiersEqual(tk, local.Name) {
			return i
		}
	}

	return -1
}

func addLocal(tk *scanner.Token) {
	if current.LocalCount == math.MaxUint8+1 {
		error("Too many local variables in function.")
		return
	}
	local := &current.Locals[current.LocalCount]
	current.LocalCount++
	local.Name = tk
	local.Depth = -1
}

func parseVariable(errorMessage string) byte {
	consume(token.Identifier, errorMessage)

	declareVariable()
	if current.ScoreDepth > 0 {
		return 0
	}
	return identifierConstant(parser.previous)
}

func markInitialized() {
	if current.ScoreDepth == 0 {
		return
	}
	current.Locals[current.LocalCount-1].Depth = current.ScoreDepth
}

func defineVariable(global byte) {
	if current.ScoreDepth > 0 {
		markInitialized()
		return
	}
	emitBytes(OpDefineGlobal, global)
}

func argumentList() byte {
	var argCount byte = 0
	var firstRound = true
	if !check(token.RParentheses) {
		for firstRound || match(token.Comma) {
			firstRound = false
			expression()
			if argCount == 255 {
				error("Can't have more than 255 arguments.")
			}
			argCount++
		}
	}

	consume(token.RParentheses, "Expect ')' after arguments.")
	return argCount
}

func declareVariable() {
	if current.ScoreDepth == 0 {
		return
	}
	tk := parser.previous

	for i := current.LocalCount - 1; i >= 0; i-- {
		local := &current.Locals[i]
		if local.Depth != -1 && local.Depth < current.ScoreDepth {
			break
		}

		if identifiersEqual(tk, local.Name) {
			error("Already variable with this name in this scope.")
		}
	}

	addLocal(tk)
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
