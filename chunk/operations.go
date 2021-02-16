package chunk

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
)
