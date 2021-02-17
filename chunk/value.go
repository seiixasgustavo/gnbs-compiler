package chunk

import "fmt"

type ValueType int

const (
	TypeBool ValueType = iota
	TypeInteger
	TypeFloat
	TypeString
	TypeNull
)

type Value struct {
	Type  ValueType
	Value interface{}
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
	value, _ := v.Value.(string)
	return value
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
	default:
		fmt.Printf("%g", value.Value)
		break
	}

}
