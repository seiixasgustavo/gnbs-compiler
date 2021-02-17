package compiler

import (
	"GNBS/chunk"
	"fmt"
	"os"
)

const StackMax = 256

type VM struct {
	Chunk *chunk.Chunk
	Ip    int

	stack    []chunk.Value
	stackTop int
}

var vm *VM

func InitVM(ck *chunk.Chunk) {
	vm = &VM{Chunk: ck, Ip: 0, stackTop: 0, stack: make([]chunk.Value, StackMax)}
}

type InterpretResult int

const (
	InterpretOk InterpretResult = iota
	InterpretCompileError
	InterpretRuntimeError
)

func (v *VM) interpret(source []byte) InterpretResult {
	chunk := chunk.NewChunk()
	if !Compile(source, chunk) {
		return InterpretCompileError
	}

	v.Chunk = chunk
	v.Ip = 0

	return v.run()
}

func (v *VM) run() InterpretResult {
	for {
		instruction := readByte()

		switch instruction {
		case OpReturn:
			chunk.PrintValue(pop())
			fmt.Println()
			return InterpretOk
		case OpConstant:
			constant := readConstant()
			push(constant)
			break
		case OpNull:
			push(chunk.Value{
				Type:  chunk.TypeNull,
				Value: nil,
			})
			break
		case OpNot:
			val := pop()
			push(chunk.Value{
				Type:  chunk.TypeBool,
				Value: !val.Bool(),
			})
			break
		case OpTrue:
			push(chunk.Value{
				Type:  chunk.TypeNull,
				Value: true,
			})
			break
		case OpFalse:
			push(chunk.Value{
				Type:  chunk.TypeNull,
				Value: false,
			})
			break
		case OpEqual:
			val, val2 := pop(), pop()
			push(chunk.Value{
				Type:  chunk.TypeBool,
				Value: val.Value == val2.Value,
			})
			break
		case OpNegate:
			if typ := peek(0); typ.Type != chunk.TypeInteger && typ.Type != chunk.TypeFloat {
				runtimeError("Operand must be a number.")
				return InterpretRuntimeError
			}

			val := pop()
			if val.Type == chunk.TypeInteger {
				val.Value = -val.Integer()
			} else {
				val.Value = -val.Float()
			}
			push(val)
			break

		case OpAdd, OpSubtract, OpMultiply, OpDivide:
			result := binaryOperation(instruction)
			if result != InterpretOk {
				runtimeError("Operands must be numbers of the same type")
				return InterpretRuntimeError
			}
			break
		}
	}
}

func readByte() byte {
	vm.Ip++
	return vm.Chunk.Code[vm.Ip-1]
}

func readConstant() chunk.Value {
	index := readByte()
	return vm.Chunk.Values[index]
}

func resetStack() {
	vm.stackTop = 0
	vm.stack = []chunk.Value{}
}

func pop() chunk.Value {
	vm.stackTop--
	return vm.stack[vm.stackTop]
}

func push(value chunk.Value) {
	vm.stackTop++
	vm.stack[vm.stackTop-1] = value
}

func peek(distance int) chunk.Value {
	return vm.stack[vm.stackTop-1-distance]
}

// Errors

func runtimeError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args)
	fmt.Fprintln(os.Stderr)

	pos := vm.Chunk.Pos[vm.Ip]

	fmt.Fprintf(os.Stderr, "[line %4d:%3d] in %s\n", pos.Line, pos.Column, pos.Filename)

	resetStack()
}
