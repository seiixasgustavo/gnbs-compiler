package compiler

type Value float64

type ValueArray struct {
	values []Value
}

func NewValueArray() *ValueArray {
	return &ValueArray{}
}

func (v *ValueArray) WriteValueArray(value Value) {
	v.values = append(v.values, value)
}

func (v *ValueArray) AddConstant(value Value) OpCode {
	v.WriteValueArray(value)
	return OpCode(len(v.values) - 1)
}
