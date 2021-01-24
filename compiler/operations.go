package compiler

type OpCode int

func (o *OpCode) Validate() bool {
	if *o < OpConstant && *o > OpReturn {
		return false
	}
	return true
}

const (
	OpConstant OpCode = iota
	OpReturn
)
