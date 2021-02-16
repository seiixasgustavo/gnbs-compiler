package chunk

import (
	"GNBS/token"
	"fmt"
	"testing"
)

func TestChunk_DisassembleChunk(t *testing.T) {
	fmt.Printf("\n\n")
	chunk := NewChunk()
	pos := token.Position{"test", 1, 123, 5}
	chunk.writeChunk(OpReturn, pos)

	constant := chunk.AddConstant(1.2)
	chunk.writeChunk(OpConstant, pos)
	chunk.writeChunk(constant, pos)

	chunk.writeChunk(OpNegate, pos)
	chunk.DisassembleChunk("test")
	fmt.Printf("\n\n")
}
