package compiler

import "fmt"

type InterpretResult int

const (
	InterpretOk InterpretResult = iota
	InterpretCompileError
	InterpretRuntimeError
)

var vm *VM

type VM struct {
	chunk    *Chunk
	ip       int
	stack    [256]Value
	stackTop int
}

func NewVM() *VM {
	return &VM{}
}

func (v *VM) Interpret(source []byte) InterpretResult {
	v.compile(source)
	return InterpretOk
}

func (v *VM) compile(source []byte) {
	
}

func (v *VM) run() InterpretResult {
	for {
		operation := v.readByte()
		switch operation {
		case OpConstant:
			constant := v.readConstant()
			v.push(constant)
			fmt.Printf("%g\n", constant)
		case OpAdd, OpSubtract, OpDivide, OpMultiply:
			v.binaryOp(operation)
			break
		case OpNegate:
			v.push(-v.pop())
			break
		case OpReturn:
			printValue(v.pop())
			fmt.Println()
			return InterpretOk
		}
	}
}

func (v *VM) resetStack() {
	v.stackTop = 0
}

func (v *VM) push(value Value) {
	v.stack[v.stackTop] = value
	v.stackTop++
}

func (v *VM) pop() Value {
	v.stackTop--
	return v.stack[v.stackTop]
}

func (v *VM) readByte() byte {
	v.ip++
	return v.chunk.Code[v.ip-1]
}

func (v *VM) readConstant() Value {
	return v.chunk.Constants.values[v.readByte()]
}

func (v *VM) DebugTrace() {
	fmt.Print("          ")

	for i := 0; i < v.stackTop; i++ {
		fmt.Print("[ ")
		printValue(v.stack[i])
		fmt.Print(" ]")
	}
	fmt.Println()

	v.chunk.disassembleInstruction(v.ip)
}

func (v *VM) binaryOp(op byte) {
	a := v.pop()
	b := v.pop()

	switch op {

	case OpAdd:
		v.push(a + b)
		break
	case OpSubtract:
		v.push(a - b)
		break
	case OpMultiply:
		v.push(a * b)
		break
	case OpDivide:
		v.push(a / b)
		break

	}
}

// Standalone Functions

func printValue(value Value) {
	fmt.Print(value)
}
