package chunk

import (
	"GNBS/token"
	"encoding/binary"
)

type Chunk struct {
	Code   []byte
	Values []Value
	Pos    []token.Position
}

func NewChunk() *Chunk {
	return &Chunk{}
}

func (c *Chunk) WriteChunk(by byte, position token.Position) {
	c.Code = append(c.Code, by)
	c.Pos = append(c.Pos, position)
}

func (c *Chunk) AddConstant(value Value) byte {
	c.Values = append(c.Values, value)
	length := len(c.Values) - 1

	bs := make([]byte, 4)
	binary.LittleEndian.PutUint16(bs, uint16(length))
	return bs[0]
}
