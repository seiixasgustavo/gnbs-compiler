package main

import "GNBS/compiler"

func main() {
	chunk := compiler.NewChunk()
	chunk.WriteChunk(compiler.OpConstant, 123)

	constant := chunk.Constants.AddConstant(1.2)
	chunk.WriteChunk(constant, 123)

	chunk.WriteChunk(compiler.OpReturn, 123)
	chunk.DisassembleChunk("test chunk")
}
