package bytesvalue

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/hyprstereo/go-dao/encoding/json"
	"github.com/hyprstereo/go-dao/encoding/msg"
	"github.com/hyprstereo/go-dao/utils"
)

type EncoderType uint8

const (
	AUTO EncoderType = iota
	JSON
	MSGPACK
)

type Bytes []byte

func NewBytesBuffer(size int) (b Bytes) {
	b = make([]byte, size)
	return
}

func NewBytes(src any, encoder ...EncoderType) Bytes {
	encodingType := AUTO
	if len(encoder) > 0 {
		encodingType = encoder[0]
	}
	var data []byte
	switch val := src.(type) {

	case string:
		data = []byte(val)
	case int:
		data = Int(val).Bytes()
	case float64:
		data = Float(val).Bytes()
	case bool: //TODO: is it needed?
		if val {
			data = []byte{uint8(1)}
		} else {
			data = []byte{uint8(0)}
		}
	default:
		if encodingType == MSGPACK {

			data, _ = msg.Encode(val)
		} else {
			data = json.Encode(val)
		}
	}
	return Bytes(data)
}

func (b Bytes) Empty() bool {
	return len(b) == 0
}

func (b Bytes) String() string {
	return string(b)
}

func (b Bytes) GoString() string {
	return string(b)
}
func (b Bytes) Contains(f []byte) (ok bool) {
	ok = bytes.Contains(b, f)
	return
}

func (b Bytes) Encode(v any) (r Bytes, err error) {

	switch val := v.(type) {
	case map[string]any:
		b = Bytes(json.Encode(val))
	case string:
		b = Bytes(val)
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

func (b Bytes) Join(src []Bytes, delim ...rune) (r Bytes) {
	if len(src) > 0 {
		r = Bytes{}
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

func (b Bytes) Split(delim ...rune) (r []Bytes) {
	res := bytes.FieldsFunc(b, func(r rune) bool {
		for _, v := range delim {
			if ok := r == v; ok {
				return ok
			}
		}
		return false
	})

	r = make([]Bytes, 0)
	for _, v := range res {
		r = append(r, v)
	}
	return
}

func (b Bytes) UnsafeString() string {
	return utils.UnsafeString(b)
}

func (b Bytes) Size() int {
	return len(b)
}

func (b Bytes) Humanize() string {
	return utils.ByteSize(uint64(b.Size()))
}

func (b Bytes) Write(value []byte) int {
	return copy(b[0:], value)
}

func (b Bytes) WriteString(s string) int {

	return b.Write(utils.UnsafeBytes(s))
}

func (b Bytes) WriteRune(r []rune) int {

	pos := b.Size()
	for _, rv := range r {
		b = append(b, byte(rv))
	}
	return pos
}

func (b Bytes) WriteTo(w io.Writer) (int, error) {
	return w.Write(b)
}

func (b Bytes) ReadFrom(r io.Reader) (int, error) {
	return r.Read(b)
}

func (b Bytes) Clear() {
	b = []byte{}
}

func (b Bytes) Append(data []byte) Bytes {
	b = append(b, data...)

	return b
}

func (b Bytes) Slice(pos int, length ...int) (data []byte) {
	if len(length) > 0 {
		return b[pos:length[0]]
	}
	return b[pos:]
}

func (b Bytes) Find(value []byte) (start, end int, extract []byte) {
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

func (b Bytes) IsJSON() bool {
	return json.Valid(b)
}

func (v Bytes) UiInt64(i ...uint64) (val uint64) {
	if len(i) > 0 {
		binary.LittleEndian.PutUint64(v, i[0])
		val = i[0]
	} else {
		ui := binary.BigEndian.Uint64(v)
		val = ui
	}
	return
}

func (v Bytes) Uint32(i ...uint32) (val uint32) {
	if len(i) > 0 {
		binary.LittleEndian.PutUint32(v, i[0])
		val = i[0]
	} else {
		ui := binary.BigEndian.Uint32(v)
		val = ui
	}
	return
}

func (v Bytes) Uint16(i ...uint16) (val uint16) {
	if len(i) > 0 {
		binary.LittleEndian.PutUint16(v, i[0])
		val = i[0]
	} else {
		ui := binary.BigEndian.Uint16(v)
		val = ui
	}
	return
}

func (v Bytes) Float(i ...float64) (val Float) {
	if len(i) > 0 {
		v = Float(i[0]).Bytes()
	} else {
		val = Float.FromBytes(0, v)
	}
	return
}

func (v Bytes) Int(i ...int) (val Int) {
	if len(i) > 0 {
		v = Int(i[0]).Bytes()
	} else {
		val = Int.FromBytes(0, v)
	}
	return
}
