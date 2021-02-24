package chunk

import "fmt"

type ValueType int

const (
	TypeBool ValueType = iota
	TypeInteger
	TypeFloat
	TypeString
	TypeFunction
	TypeNative
	TypeNull
)

type Value struct {
	Type  ValueType
	Value interface{}
}

type GString struct {
	String string
	Hash   uint32
}

type GFunction struct {
	Arity int
	Chunk Chunk
	Name  *GString
}

type NativeFn func(argCount byte, args []Value) Value
type GNative struct {
	Function NativeFn
}

func NewGString(value string) *GString {
	return &GString{
		value,
		hashString(value),
	}
}

func NewGFunction() *GFunction {
	return &GFunction{
		Arity: 0,
		Chunk: *NewChunk(),
		Name:  nil,
	}
}

func NewGNative(function NativeFn) *GNative {
	return &GNative{Function: function}
}

func (v *Value) Bool() bool {
	value, _ := v.Value.(bool)
	return value
}

func (v *Value) Integer() int {
	value, _ := v.Value.(int)
	return value
}

func (v *Value) Float() float64 {
	value, _ := v.Value.(float64)
	return value
}

func (v *Value) String() string {
	value, _ := v.Value.(*GString)
	return value.String
}

func (v *Value) FunctionName() string {
	value, _ := v.Value.(*GFunction)
	return value.Name.String
}

func PrintValue(value Value) {
	switch value.Type {
	case TypeBool:
		if value.Bool() {
			fmt.Println("true")
		} else {
			fmt.Println("false")
		}
		break
	case TypeNull:
		fmt.Println("null")
		break
	case TypeInteger:
		fmt.Printf("%g", value.Integer())
		break
	case TypeFloat:
		fmt.Printf("%g", value.Float())
		break
	case TypeString:
		fmt.Print(value.String())
		break
	case TypeFunction:
		if name := value.FunctionName(); name == "" {
			fmt.Printf("<script>")
		} else {
			fmt.Printf("<fn %s>", name)
		}
		break
	case TypeNative:
		fmt.Printf("<native fn>")
		break
	default:
		fmt.Printf("%g", value.Value)
		break
	}

}
func hashString(key string) uint32 {
	var hash uint32 = 2166136261
	for i := 0; i < len(key); i++ {
		hash ^= uint32(key[i])
		hash *= 16777619
	}
	return hash
}
