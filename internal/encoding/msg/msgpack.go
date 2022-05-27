package msg

import (
	"github.com/vmihailenco/msgpack/v5"
)

func Encode(v any) (data []byte, err error) {
	data, err = msgpack.Marshal(v)
	return
}

func Decode(data []byte, v any) (err error) {
	err = msgpack.Unmarshal(data, v)
	return
}
