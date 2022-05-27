package bytesvalue

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"
)

type Int int64

func (i Int) From(value int) (r Int) {
	i = Int(value)
	return
}

func (i Int) String() string {
	return fmt.Sprint(i.I64())
}

func (i Int) FromBytes(v []byte) Int {
	return Int(new(big.Int).SetBytes(v).Int64())
}
func (i Int) I64() int64 {
	return int64(i)
}

func (i Int) I32() int32 {
	return int32(i)
}

func (i Int) I16() int16 {
	return int16(i)
}

func (i Int) I8() int8 {
	return int8(i)
}

func (i Int) Bytes() []byte {
	v := new(big.Int).SetInt64(i.I64()).Bytes()
	return v
}

type Float float64

func (i Float) F32() float32 {
	return float32(i)
}

func (i Float) FromBytes(v []byte) Float {
	bits := binary.LittleEndian.Uint64(v)
	i = Float(math.Float64frombits(bits))
	return i
}

func (i Float) Bytes() []byte {
	bits := math.Float64bits(float64(i))
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func (i Float) String() string {
	return fmt.Sprint(float64(i))
}
