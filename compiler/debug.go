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
	case OpPrint:
		return simpleInstruction("OP_PRINT", offset)
	case OpPop:
		return simpleInstruction("OP_POP", offset)
	case OpDefineGlobal:
		return constantInstruction("OP_DEFINE_GLOBAL", c, offset)
	case OpGetLocal:
		return byteInstruction("OP_GET_LOCAL", c, offset)
	case OpSetLocal:
		return byteInstruction("OP_SET_LOCAL", c, offset)
	case OpGetGlobal:
		return constantInstruction("OP_GET_GLOBAL", c, offset)
	case OpJump:
		return jumpInstruction("OP_JUMP", 1, c, offset)
	case OpJumpIfFalse:
		return jumpInstruction("OP_JUMP_IF_FALSE", 1, c, offset)
	case OpLoop:
		return jumpInstruction("OP_LOOP", -1, c, offset)
	case OpCall:
		return byteInstruction("OP_CALL", c, offset)
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

func byteInstruction(name string, c *chunk.Chunk, offset int) int {
	slot := c.Code[offset+1]
	fmt.Printf("%-16s %4d\n", name, slot)
	return offset + 2
}

func jumpInstruction(name string, sign uint16, c *chunk.Chunk, offset int) int {
	jump := uint16(c.Code[offset + 1]) << 8
	jump |= uint16(c.Code[offset + 2])

	fmt.Printf("%-16s %4d -> %d\n", name, offset, uint16(offset) + 3 + sign * jump)
	return offset + 3
}
