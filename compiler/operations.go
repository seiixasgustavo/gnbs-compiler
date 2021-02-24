package compiler

import "GNBS/chunk"

const (
	OpReturn byte = iota
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
	OpConstant
	OpNegate
	OpNull
	OpTrue
	OpFalse
	OpNot
	OpEqual
	OpGreater
	OpLess
	OpPrint
	OpPop
	OpDefineGlobal

	OpGetGlobal
	OpGetLocal
	OpSetGlobal
	OpSetLocal

	OpJump
	OpJumpIfFalse

	OpLoop

	OpCall
)

func binaryOperation(operation byte) InterpretResult {
	if operation == OpAdd {
		val, val2 := peek(0), peek(1)
		if val.Type == chunk.TypeString && val2.Type == chunk.TypeString {
			val, val2 = pop(), pop()
			push(chunk.Value{
				Type:  chunk.TypeString,
				Value: chunk.NewGString(val.String() + val2.String()),
			})
		}
	}

	if val, val2 := peek(0), peek(1);
		(val.Type != chunk.TypeFloat && val2.Type != chunk.TypeFloat) ||
			(val.Type != chunk.TypeInteger && val2.Type != chunk.TypeInteger) {
		runtimeError("Operands must be numbers.")
		return InterpretRuntimeError
	}

	if val, val2 := peek(0), peek(1);
		(val.Type != chunk.TypeInteger && val2.Type != chunk.TypeFloat) ||
			(val.Type != chunk.TypeFloat && val2.Type != chunk.TypeInteger) {
		runtimeError("Operands must have the same type")
		return InterpretRuntimeError
	}

	val := pop()
	val2 := pop()

	switch operation {
	case OpAdd, OpSubtract, OpMultiply, OpDivide:
		if val.Type == chunk.TypeInteger {
			return binaryIntegerOperation(operation, val, val2)
		} else {
			return binaryFloatOperation(operation, val, val2)
		}
	case OpGreater, OpLess:
		return binaryComparison(operation, val, val2)
	}
	return InterpretRuntimeError
}

func binaryIntegerOperation(operation byte, val, val2 chunk.Value) InterpretResult {
	switch operation {
	case OpAdd:
		push(chunk.Value{
			Type:  chunk.TypeInteger,
			Value: val.Integer() + val2.Integer(),
		})
		break
	case OpNegate:
		push(chunk.Value{
			Type:  chunk.TypeInteger,
			Value: val.Integer() - val2.Integer(),
		})
		break
	case OpMultiply:
		push(chunk.Value{
			Type:  chunk.TypeInteger,
			Value: val.Integer() * val2.Integer(),
		})
		break
	case OpDivide:
		push(chunk.Value{
			Type:  chunk.TypeInteger,
			Value: val.Integer() / val2.Integer(),
		})
		break
	}
	return InterpretOk
}
func binaryFloatOperation(operation byte, val, val2 chunk.Value) InterpretResult {
	switch operation {
	case OpAdd:
		push(chunk.Value{
			Type:  chunk.TypeInteger,
			Value: val.Float() + val2.Float(),
		})
		break
	case OpNegate:
		push(chunk.Value{
			Type:  chunk.TypeInteger,
			Value: val.Float() - val2.Float(),
		})
		break
	case OpMultiply:
		push(chunk.Value{
			Type:  chunk.TypeInteger,
			Value: val.Float() * val2.Float(),
		})
		break
	case OpDivide:
		push(chunk.Value{
			Type:  chunk.TypeInteger,
			Value: val.Float() / val2.Float(),
		})
	}
	return InterpretOk
}
func binaryComparison(operation byte, val, val2 chunk.Value) InterpretResult {
	if val.Type == chunk.TypeFloat {
		switch operation {
		case OpGreater:
			push(chunk.Value{
				Type:  chunk.TypeBool,
				Value: val.Float() > val2.Float(),
			})
			break
		case OpLess:
			push(chunk.Value{
				Type:  chunk.TypeBool,
				Value: val.Float() < val2.Float(),
			})
		}
		return InterpretOk
	} else {
		switch operation {
		case OpGreater:
			push(chunk.Value{
				Type:  chunk.TypeBool,
				Value: val.Integer() > val2.Integer(),
			})
			break
		case OpLess:
			push(chunk.Value{
				Type:  chunk.TypeBool,
				Value: val.Integer() < val2.Integer(),
			})
		}
		return InterpretOk
	}
}
