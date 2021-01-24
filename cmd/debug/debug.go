package main

import "GNBS/compiler"

func main() {
	chunk := compiler.NewChunk()
	vm := compiler.NewVM()

	constant := chunk.Constants.AddConstant(1.2)
	chunk.WriteChunk(compiler.OpConstant, 123)
	chunk.WriteChunk(constant, 123)

	constant = chunk.Constants.AddConstant(3.4)
	chunk.WriteChunk(compiler.OpConstant, 123)
	chunk.WriteChunk(constant, 123)

	chunk.WriteChunk(compiler.OpAdd, 123)

	constant = chunk.Constants.AddConstant(5.6)
	chunk.WriteChunk(compiler.OpConstant, 123)
	chunk.WriteChunk(constant, 123)

	chunk.WriteChunk(compiler.OpDivide, 123)
	chunk.WriteChunk(compiler.OpNegate, 123)
	chunk.WriteChunk(compiler.OpReturn, 123)
	chunk.DisassembleChunk("test chunk")

	vm.Interpret(chunk)
}
