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

func (v *ValueArray) AddConstant(value Value) byte {
	v.WriteValueArray(value)
	return byte(len(v.values) - 1)
}
