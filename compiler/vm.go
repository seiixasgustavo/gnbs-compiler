package compiler

import (
	"GNBS/chunk"
	"fmt"
	"math"
	"os"
)

const (
	FrameMax = 64
	StackMax = (math.MaxUint8 + 1) * FrameMax
)

type VM struct {
	Frames     []CallFrame
	FrameCount int

	stack    []chunk.Value
	stackTop int

	stringsTable *chunk.Table
	globals      *chunk.Table
}

type CallFrame struct {
	Function *chunk.GFunction
	Ip       uint16
	Code     []byte
	Slots    []chunk.Value
}

var vm *VM

func InitVM(ck *chunk.Chunk) {
	vm = &VM{
		FrameCount:   0,
		stackTop:     0,
		stack:        make([]chunk.Value, StackMax),
		stringsTable: chunk.NewTable(),
		globals:      chunk.NewTable(),
	}
}

type InterpretResult int

const (
	InterpretOk InterpretResult = iota
	InterpretCompileError
	InterpretRuntimeError
)

func (v *VM) interpret(source []byte) InterpretResult {
	fn := Compile(source)
	if fn == nil {
		return InterpretCompileError
	}
	push(chunk.Value{
		Type:  chunk.TypeFunction,
		Value: fn,
	})
	frame := &vm.Frames[vm.FrameCount]
	vm.FrameCount++
	frame.Ip = 0
	frame.Code = fn.Chunk.Code
	frame.Slots = vm.stack

	return run()
}

func run() InterpretResult {
	frame := &vm.Frames[vm.FrameCount-1]

	for {
		instruction := readByte(frame)

		switch instruction {
		case OpReturn:
			result := pop()
			vm.FrameCount--
			if vm.FrameCount == 0 {
				pop()
				return InterpretOk
			}
			vm.stackTop -= len(frame.Slots)
			push(result)
			frame = &vm.Frames[vm.FrameCount-1]
			break
		case OpConstant:
			constant := readConstant(frame)
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

		case OpPrint:
			chunk.PrintValue(pop())
			fmt.Println()
			break

		case OpPop:
			pop()
			break

		case OpDefineGlobal:
			name := readString(frame)
			vm.globals.TableSet(name, peek(0))
			pop()
			break

		case OpGetLocal:
			slot := readByte(frame)
			push(frame.Slots[slot])
			break

		case OpGetGlobal:
			name := readString(frame)
			var value chunk.Value

			if vm.globals.TableGet(name, &value) {
				runtimeError("Undefined variable '%s'.", name.String)
				return InterpretRuntimeError
			}
			push(value)
			break

		case OpSetLocal:
			slot := readByte(frame)
			frame.Slots[slot] = peek(0)
			break

		case OpSetGlobal:
			name := readString(frame)
			if vm.globals.TableSet(name, peek(0)) {
				vm.globals.TableDelete(name)
				runtimeError("Undefined variable '%s'.", name.String)
				return InterpretRuntimeError
			}
			break

		case OpJump:
			offset := readShort(frame)
			frame.Ip += offset
			break

		case OpJumpIfFalse:
			offset := readShort(frame)
			if isFalsey(peek(0)) {
				frame.Ip += offset
			}
			break

		case OpLoop:
			offset := readShort(frame)
			frame.Ip -= offset
			break

		case OpCall:
			argCount := readByte(frame)
			if !callValue(peek(argCount), argCount) {
				return InterpretRuntimeError
			}
			frame = &vm.Frames[vm.FrameCount-1]
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

func readByte(frame *CallFrame) byte {
	frame.Ip++
	return frame.Code[frame.Ip-1]
}

func readConstant(frame *CallFrame) chunk.Value {
	index := readByte(frame)
	return frame.Slots[index]
}

func readString(frame *CallFrame) *chunk.GString {
	value := readConstant(frame)
	return value.Value.(*chunk.GString)
}

func readShort(frame *CallFrame) uint16 {
	frame.Ip += 2
	return uint16(frame.Code[frame.Ip-2])<<8 | uint16(frame.Code[frame.Ip-1])
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

func peek(distance byte) chunk.Value {
	return vm.stack[vm.stackTop-1-int(distance)]
}

func call(fn *chunk.GFunction, argCount byte) bool {
	if int(argCount) != fn.Arity {
		runtimeError("Expected %d arguments but got %d.", fn.Arity, argCount)
		return false
	}

	if vm.FrameCount == FrameMax {
		runtimeError("Stack overflow.")
		return false
	}

	frame := &vm.Frames[vm.FrameCount]
	vm.FrameCount++
	frame.Ip, frame.Function, frame.Code = 0, fn, fn.Chunk.Code

	frame.Slots = vm.stack[argCount-1:vm.stackTop]
	return true
}

func callValue(callee chunk.Value, argCount byte) bool {
	switch callee.Type {
	case chunk.TypeFunction:
		fn, _ := callee.Value.(*chunk.GFunction)
		return call(fn, argCount)
	case chunk.TypeNative:
		fn, _ := callee.Value.(*chunk.GNative)
		result := fn.Function(argCount, vm.stack[vm.stackTop-int(argCount):])
		vm.stackTop -= int(argCount) + 1
		push(result)
		return true
	default:
		break
	}

	runtimeError("Can only call functions and classes.")
	return false
}

func isFalsey(value chunk.Value) bool {
	return value.Type == chunk.TypeNull || (value.Type == chunk.TypeBool && value.Bool())
}

// Errors

func runtimeError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args)
	fmt.Fprintln(os.Stderr)

	for i := vm.FrameCount - 1; i >= 0; i-- {
		frame := &vm.Frames[i]
		fn := frame.Function

		instruction := frame.Code[frame.Ip-1]
		pos := frame.Function.Chunk.Pos[instruction]
		fmt.Fprintf(os.Stderr, "[line %4d:%3d] in %s\n", pos.Line, pos.Column, pos.Filename)

		if fn.Name == nil {
			fmt.Fprintf(os.Stderr, "script\n")
		} else {
			fmt.Fprintf(os.Stderr, "%s()\n", fn.Name.String)
		}

	}



	resetStack()
}

// Native Values

func defineNative(name string, function chunk.NativeFn) {
	push(chunk.Value{
		Type:  chunk.TypeString,
		Value: chunk.NewGString(name),
	})
	push(chunk.Value{
		Type:  chunk.TypeNative,
		Value: chunk.NewGNative(function),
	})
	str, _ := vm.stack[0].Value.(*chunk.GString)
	vm.globals.TableSet(str, vm.stack[1])
	pop()
	pop()
}
