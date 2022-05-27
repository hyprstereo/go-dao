package godao

import (
	"github.com/hyprstereo/go-dao/pkg/encoding/json"
	"github.com/hyprstereo/go-dao/pkg/encoding/msg"
)

type RawValue = json.RawValue

func MsgEncode(v any, pretty ...bool) (data RawValue, err error) {
	return msg.Encode(v)
}

func MsgDecode(data []byte, v any) (err error) {
	return msg.Decode(data, v)
}
