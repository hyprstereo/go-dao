package dao

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"math"
	"math/big"
	"reflect"

	"github.com/hyprstereo/go-dao/encoding/json"
	"github.com/hyprstereo/go-dao/utils"

	"github.com/vmihailenco/msgpack/v5"
)

type Int int

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

type BytesValue []byte

func (b BytesValue) Empty() bool {
	return len(b) == 0
}

func (b BytesValue) String() string {
	return string(b)
}

func (b BytesValue) GoString() string {
	return string(b)
}
func (b BytesValue) Contains(f []byte) (ok bool) {
	ok = bytes.Contains(b, f)
	return
}

func (b BytesValue) Encode(v any) (r BytesValue, err error) {

	switch val := v.(type) {
	case map[string]any:
		b = BytesValue(json.Encode(val))
	case Map:
		b = BytesValue(val.Bytes())
	case string:
		b = BytesValue(val)
	case []byte:
		b = val
	case int:
		b = Int(val).Bytes()
	case float64:
		b = Float(val).Bytes()
	default:
		err = fmt.Errorf("cannot convert %v", val)
	}
	r = b
	return
}

func (b BytesValue) Join(src []BytesValue, delim ...rune) (r BytesValue) {
	if len(src) > 0 {
		r = BytesValue{}
		var del byte
		for x, d := range src {
			r = append(r, d...)
			if x < len(delim)-1 {
				del = byte(delim[x])
			} else {
				del = byte(delim[len(delim)-1])
			}
			r = append(r, del)
		}
	}
	return
}

func (b BytesValue) Split(delim ...rune) (r []BytesValue) {
	res := bytes.FieldsFunc(b, func(r rune) bool {
		for _, v := range delim {
			if ok := r == v; ok {
				return ok
			}
		}
		return false
	})

	r = make([]BytesValue, 0)
	for _, v := range res {
		r = append(r, v)
	}
	return
}

func (b BytesValue) UnsafeString() string {
	return UnsafeString(b)
}

func (b BytesValue) Size() int {
	return len(b)
}

func (b BytesValue) Humanize() string {
	return utils.ByteSize(uint64(b.Size()))
}

func (b BytesValue) Write(value []byte) int {
	return copy(b[0:], value)
}

func (b BytesValue) WriteString(s string) int {

	return b.Write(UnsafeBytes(s))
}

func (b BytesValue) WriteRune(r []rune) int {

	pos := b.Size()
	for _, rv := range r {
		b = append(b, byte(rv))
	}
	return pos
}

func (b BytesValue) WriteTo(w io.Writer) (int, error) {
	return w.Write(b)
}

func (b BytesValue) ReadFrom(r io.Reader) (int, error) {
	return r.Read(b)
}

func (b BytesValue) Clear() {
	b = []byte{}
}

func (b BytesValue) Append(data []byte) BytesValue {
	b = append(b, data...)

	return b
}

func (b BytesValue) Slice(pos int, length ...int) (data []byte) {
	if len(length) > 0 {
		return b[pos:length[0]]
	}
	return b[pos:]
}

func (b BytesValue) Find(value []byte) (start, end int, extract []byte) {
	pos := bytes.IndexByte(b, value[0])
	if pos > -1 {
		for v := 1; v < len(value); v++ {
			fmt.Println(value[v], b[pos+v])
			if value[v] != b[pos+v] {
				break
			}
		}
	}
	start = pos
	end = start + len(value)
	extract = b[start:end]
	return
}

func (b BytesValue) IsJSON() bool {
	return json.Valid(b)
}

func (v BytesValue) UiInt64(i ...uint64) (val uint64) {
	if len(i) > 0 {
		binary.LittleEndian.PutUint64(v, i[0])
		val = i[0]
	} else {
		ui := binary.BigEndian.Uint64(v)
		val = ui
	}
	return
}

func (v BytesValue) Uint32(i ...uint32) (val uint32) {
	if len(i) > 0 {
		binary.LittleEndian.PutUint32(v, i[0])
		val = i[0]
	} else {
		ui := binary.BigEndian.Uint32(v)
		val = ui
	}
	return
}

func (v BytesValue) Uint16(i ...uint16) (val uint16) {
	if len(i) > 0 {
		binary.LittleEndian.PutUint16(v, i[0])
		val = i[0]
	} else {
		ui := binary.BigEndian.Uint16(v)
		val = ui
	}
	return
}

func (v BytesValue) Float(i ...float64) (val Float) {
	if len(i) > 0 {
		v = Float(i[0]).Bytes()
	} else {
		val = Float.FromBytes(0, v)
	}
	return
}

func (v BytesValue) Int(i ...int) (val Int) {
	if len(i) > 0 {
		v = Int(i[0]).Bytes()
	} else {
		val = Int.FromBytes(0, v)
	}
	return
}

func GOBEncode(v any, useValues ...bool) (data []byte, err error) {
	buff := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buff)
	if len(useValues) > 0 && useValues[0] {
		err = enc.EncodeValue(reflect.ValueOf(v))
	} else {
		err = enc.Encode(v)
	}
	if err == nil {
		data = buff.Bytes()
	}
	return
}

func GOBDecode(data []byte, src any, useValues ...bool) (err error) {
	buff := bytes.NewReader(data)
	dec := gob.NewDecoder(buff)
	if len(useValues) > 0 && useValues[0] {
		err = dec.DecodeValue(reflect.ValueOf(src))
	} else {
		err = dec.Decode(src)
	}
	return
}

func MSGPackEncode(v any) (data []byte, err error) {
	data, err = msgpack.Marshal(v)
	return
}

func MSGPackDecode(data []byte, v any) (err error) {
	err = msgpack.Unmarshal(data, v)
	return
}
