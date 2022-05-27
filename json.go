package godao

import "github.com/hyprstereo/go-dao/internal/encoding/json"

func JsonEncode(v any, pretty ...bool) (data json.RawValue) {
	return json.Encode(v, pretty...)
}

func JsonDecode(data []byte, v any) (err error) {
	return json.Decode(data, v)
}
