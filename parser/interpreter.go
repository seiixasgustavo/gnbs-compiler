package parser

import (
	"GNBS/token"
	"reflect"
)

type Interpreter struct{}

func (i Interpreter) visitLiteralExpr(expr literal) interface{} {
	return expr.value
}

func (i Interpreter) visitGroupingExpr(expr grouping) interface{} {
	return i.evaluate(expr)
}

func (i Interpreter) visitUnaryExpr(expr unary) interface{} {
	right := i.evaluate(expr.right)

	switch expr.operator.TokenType {
	case token.Not:
		return !i.isTruthy(right)
	case token.Minus:
		if reflect.TypeOf(right) == reflect.TypeOf(10.0) {
			return -(right.(float64))
		}
		return -(right.(int))
	}
	return nil
}

func (i Interpreter) visitBinaryExpr(expr binary) interface{} {
	left := i.evaluate(expr.left)
	right := i.evaluate(expr.right)
	return i.operation(left, right, expr.operator.TokenType)
}

func (i Interpreter) evaluate(expr Expr) interface{} {
	return expr.accept(i)
}

func (i Interpreter) isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	if value, ok := object.(bool); ok {
		return value
	}
	return true
}

func (i Interpreter) isEqual(left interface{}, right interface{}) bool {
	return left == right
}

func (i Interpreter) operation(left interface{}, right interface{}, operator token.Type) interface{} {
	leftDouble, leftDoubleBool := left.(float64)
	rightDouble, rightDoubleBool := right.(float64)

	switch operator {
	case token.Minus:
		if leftDoubleBool && rightDoubleBool {
			return leftDouble - rightDouble
		} else if leftDoubleBool {
			return leftDouble - float64(right.(int))
		}
		return left.(int) - right.(int)
	case token.Star:
		if leftDoubleBool && rightDoubleBool {
			return leftDouble * rightDouble
		} else if leftDoubleBool {
			return leftDouble * float64(right.(int))
		}
		return left.(int) * right.(int)
	case token.Slash:
		if leftDoubleBool && rightDoubleBool {
			return leftDouble / rightDouble
		} else if leftDoubleBool {
			return leftDouble / float64(right.(int))
		}
		return left.(int) / right.(int)
	case token.Plus:
		if !leftDoubleBool && !rightDoubleBool {
			leftStr, leftStrBool := left.(string)
			rightStr, rightStrBool := right.(string)
			if leftStrBool && rightStrBool {
				return leftStr + rightStr
			}
		} else if leftDoubleBool && rightDoubleBool {
			return leftDouble + rightDouble
		} else if leftDoubleBool {
			return leftDouble + float64(right.(int))
		}
		return left.(int) + right.(int)
	case token.Greater:
		if leftDoubleBool && rightDoubleBool {
			return leftDouble > rightDouble
		} else if leftDoubleBool {
			return leftDouble > float64(right.(int))
		}
		return left.(int) > right.(int)
	case token.GreaterEqual:
		if leftDoubleBool && rightDoubleBool {
			return leftDouble >= rightDouble
		} else if leftDoubleBool {
			return leftDouble >= float64(right.(int))
		}
		return left.(int) >= right.(int)
	case token.Less:
		if leftDoubleBool && rightDoubleBool {
			return leftDouble < rightDouble
		} else if leftDoubleBool {
			return leftDouble < float64(right.(int))
		}
		return left.(int) < right.(int)
	case token.LessEqual:
		if leftDoubleBool && rightDoubleBool {
			return leftDouble <= rightDouble
		} else if leftDoubleBool {
			return leftDouble <= float64(right.(int))
		}
		return left.(int) <= right.(int)
	case token.NotEqual:
		return !i.isEqual(left, right)
	case token.Equal:
		return i.isEqual(left, right)
	}
	return nil
}
