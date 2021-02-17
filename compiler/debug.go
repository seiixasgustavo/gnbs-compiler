package compiler

import (
	"GNBS/chunk"
	"fmt"
)

func DisassembleChunk(c *chunk.Chunk, name string) {
	fmt.Printf("== %s ==\n", name)
	for i := 0; i < len(c.Code); {
		i = DisassembleInstruction(c, i)
	}
}

func DisassembleInstruction(c *chunk.Chunk, offset int) int {
	fmt.Printf("%04d ", offset)
	fmt.Printf("%s ", c.Pos[offset].String())

	instruction := c.Code[offset]
	switch instruction {
	case OpReturn:
		return simpleInstruction("OP_RETURN", offset)
	case OpConstant:
		return constantInstruction("OP_CONSTANT", c, offset)
	case OpNegate:
		return simpleInstruction("OP_NEGATE", offset)
	case OpNot:
		return simpleInstruction("OP_NOT", offset)
	case OpAdd:
		return simpleInstruction("OP_ADD", offset)
	case OpSubtract:
		return simpleInstruction("OP_SUBTRACT", offset)
	case OpDivide:
		return simpleInstruction("OP_DIVIDE", offset)
	case OpMultiply:
		return simpleInstruction("OP_MULTIPLY", offset)
	case OpNull:
		return simpleInstruction("OP_NULL", offset)
	case OpTrue:
		return simpleInstruction("OP_TRUE", offset)
	case OpFalse:
		return simpleInstruction("OP_FALSE", offset)
	case OpEqual:
		return simpleInstruction("OP_EQUAL", offset)
	case OpGreater:
		return simpleInstruction("OP_GREATER", offset)
	case OpLess:
		return simpleInstruction("OP_LESS", offset)
	default:
		fmt.Printf("Unknown OpCode %d\n", instruction)
		return offset + 1
	}
}

func simpleInstruction(name string, offset int) int {
	fmt.Printf("%s\n", name)
	return offset + 1
}
func constantInstruction(name string, c *chunk.Chunk, offset int) int {
	constant := c.Code[offset+1]

	fmt.Printf("%-16s %4d '", name, constant)
	chunk.PrintValue(c.Values[constant])
	fmt.Printf("'\n")
	return offset + 2
}
