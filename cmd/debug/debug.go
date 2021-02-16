package main

import "GNBS/old"

func main() {
	chunk := old.NewChunk()
	vm := old.NewVM()

	constant := chunk.Constants.AddConstant(1.2)
	chunk.WriteChunk(old.OpConstant, 123)
	chunk.WriteChunk(constant, 123)

	constant = chunk.Constants.AddConstant(3.4)
	chunk.WriteChunk(old.OpConstant, 123)
	chunk.WriteChunk(constant, 123)

	chunk.WriteChunk(old.OpAdd, 123)

	constant = chunk.Constants.AddConstant(5.6)
	chunk.WriteChunk(old.OpConstant, 123)
	chunk.WriteChunk(constant, 123)

	chunk.WriteChunk(old.OpDivide, 123)
	chunk.WriteChunk(old.OpNegate, 123)
	chunk.WriteChunk(old.OpReturn, 123)
	chunk.DisassembleChunk("test chunk")

	vm.Interpret(chunk)
}
