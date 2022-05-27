package godao

import "github.com/hyprstereo/go-dao/encoding/json"

func JsonEncode(v any, pretty ...bool) (data []byte) {
	return json.Encode(v, pretty...)
}

func JsonDecode(data []byte, v any) (err error) {
	return json.Decode(data, v)
}
