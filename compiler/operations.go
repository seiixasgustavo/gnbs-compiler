package compiler

import "GNBS/chunk"

func binaryOperation(operation byte) InterpretResult {
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
	case chunk.OpAdd, chunk.OpSubtract, chunk.OpMultiply, chunk.OpDivide:
		if val.Type == chunk.TypeInteger {
			return binaryIntegerOperation(operation, val, val2)
		} else {
			return binaryFloatOperation(operation, val, val2)
		}
	case chunk.OpGreater, chunk.OpLess:
		return binaryComparison(operation, val, val2)
	}
	return InterpretRuntimeError
}

func binaryIntegerOperation(operation byte, val, val2 chunk.Value) InterpretResult {
	switch operation {
	case chunk.OpAdd:
		push(chunk.Value{
			Type:  chunk.TypeInteger,
			Value: val.Integer() + val2.Integer(),
		})
		break
	case chunk.OpNegate:
		push(chunk.Value{
			Type:  chunk.TypeInteger,
			Value: val.Integer() - val2.Integer(),
		})
		break
	case chunk.OpMultiply:
		push(chunk.Value{
			Type:  chunk.TypeInteger,
			Value: val.Integer() * val2.Integer(),
		})
		break
	case chunk.OpDivide:
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
	case chunk.OpAdd:
		push(chunk.Value{
			Type:  chunk.TypeInteger,
			Value: val.Float() + val2.Float(),
		})
		break
	case chunk.OpNegate:
		push(chunk.Value{
			Type:  chunk.TypeInteger,
			Value: val.Float() - val2.Float(),
		})
		break
	case chunk.OpMultiply:
		push(chunk.Value{
			Type:  chunk.TypeInteger,
			Value: val.Float() * val2.Float(),
		})
		break
	case chunk.OpDivide:
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
		case chunk.OpGreater:
			push(chunk.Value{
				Type:  chunk.TypeBool,
				Value: val.Float() > val2.Float(),
			})
			break
		case chunk.OpLess:
			push(chunk.Value{
				Type:  chunk.TypeBool,
				Value: val.Float() < val2.Float(),
			})
		}
		return InterpretOk
	} else {
		switch operation {
		case chunk.OpGreater:
			push(chunk.Value{
				Type:  chunk.TypeBool,
				Value: val.Integer() > val2.Integer(),
			})
			break
		case chunk.OpLess:
			push(chunk.Value{
				Type:  chunk.TypeBool,
				Value: val.Integer() < val2.Integer(),
			})
		}
		return InterpretOk
	}
}
