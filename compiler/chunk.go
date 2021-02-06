package compiler

import "fmt"

type Chunk struct {
	Code      []byte
	Constants ValueArray
	Lines     []int
}

func NewChunk() *Chunk {
	return &Chunk{}
}

func (c *Chunk) WriteChunk(bt byte, line int) {
	c.Code = append(c.Code, bt)
	c.Lines = append(c.Lines, line)
}

func (c *Chunk) DisassembleChunk(name string) {
	fmt.Printf("== %s ==\n", name)
	for i := 0; i < len(c.Code); {
		i = c.disassembleInstruction(i)
	}
}

func (c *Chunk) disassembleInstruction(offset int) int {
	fmt.Printf("%04d ", offset)

	if offset > 0 && c.Lines[offset] == c.Lines[offset-1] {
		fmt.Print("   | ")
	} else {
		fmt.Printf("%4d ", c.Lines[offset])
	}

	opCode := c.Code[offset]
	switch opCode {
	case OpConstant:
		return c.constantInstruction("OP_CONSTANT", offset)
	case OpAdd:
		return c.simpleInstruction("OP_ADD", offset)
	case OpSubtract:
		return c.simpleInstruction("OP_SUBTRACT", offset)
	case OpMultiply:
		return c.simpleInstruction("OP_MULTIPLY", offset)
	case OpDivide:
		return c.simpleInstruction("OP_DIVIDE", offset)
	case OpNegate:
		return c.simpleInstruction("OP_NEGATE", offset)
	case OpReturn:
		return c.simpleInstruction("OP_RETURN", offset)
	default:
		fmt.Printf("Unknown opcode %d\n", opCode)
		return offset + 1
	}
}

func (c *Chunk) simpleInstruction(name string, offset int) int {
	fmt.Println(name)
	return offset + 1
}

func (c *Chunk) constantInstruction(name string, offset int) int {
	constant := c.Code[offset+1]

	fmt.Printf("%-16s %4d '", name, constant)
	fmt.Printf("%g'\n", c.Constants.values[constant])
	return offset + 2
}