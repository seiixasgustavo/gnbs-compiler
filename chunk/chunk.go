package chunk

import (
	"GNBS/token"
	"encoding/binary"
	"fmt"
)

type Chunk struct {
	Code   []byte
	Values []Value
	Pos    []token.Position
}

func NewChunk() *Chunk {
	return &Chunk{}
}

func (c *Chunk) WriteChunk(by byte, position token.Position) {
	c.Code = append(c.Code, by)
	c.Pos = append(c.Pos, position)
}

func (c *Chunk) AddConstant(value Value) byte {
	c.Values = append(c.Values, value)
	length := len(c.Values) - 1

	bs := make([]byte, 4)
	binary.LittleEndian.PutUint16(bs, uint16(length))
	return bs[0]
}

func (c *Chunk) DisassembleChunk(name string) {
	fmt.Printf("== %s ==\n", name)
	for i := 0; i < len(c.Code); {
		i = c.disassembleInstruction(i)
	}
}

func (c *Chunk) disassembleInstruction(offset int) int {
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
func constantInstruction(name string, chunk *Chunk, offset int) int {
	constant := chunk.Code[offset+1]

	fmt.Printf("%-16s %4d '", name, constant)
	PrintValue(chunk.Values[constant])
	fmt.Printf("'\n")
	return offset + 2
}

func PrintValue(value Value) {
	switch value.Type {
	case TypeBool:
		if value.Bool() {
			fmt.Println("true")
		} else {
			fmt.Println("false")
		}
		break
	case TypeNull:
		fmt.Println("null")
		break
	case TypeInteger:
		fmt.Printf("%g", value.Integer())
		break
	case TypeFloat:
		fmt.Printf("%g", value.Float())
		break
	default:
		fmt.Printf("%g", value.Value)
		break
	}

}
