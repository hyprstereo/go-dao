package json

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"strings"

	gojson "github.com/goccy/go-json"

	"github.com/spf13/afero"
	"github.com/tidwall/gjson"
)

var (
	fs = afero.NewOsFs()
)

type Raw = gojson.RawMessage
type Result struct {
	gjson.Result
}
type MarshalJSON = gojson.Marshaler
type UnmarshalJSON = gojson.Unmarshaler

func Marshal(v interface{}) ([]byte, error) {
	return gojson.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return gojson.Unmarshal(data, v)
}

func Compact(dst *bytes.Buffer, src RawValue) error {
	return gojson.Compact(dst, src)
}

func EncodeWithColor(v any, pretty ...bool) (r RawValue) {

	if len(pretty) > 0 && pretty[0] {
		r, _ = gojson.MarshalIndentWithOption(v, "", "  ", gojson.Colorize(gojson.DefaultColorScheme))
		return
	}
	r, _ = gojson.MarshalWithOption(v, gojson.Colorize(gojson.DefaultColorScheme))
	return
}

func EncodeNoEscape(v any) (r RawValue) {
	r, _ = gojson.MarshalNoEscape(v)

	return
}

func MapToStruct(m any, target any) (err error) {
	err = Decode(Encode(m), target)
	return
}

func Encode(v interface{}, pretty ...bool) (r RawValue) {
	if len(pretty) > 0 {
		if r, e := gojson.MarshalIndent(v, "", "  "); e != nil {
			//fmt.Println(e.Error())
			return nil
		} else {
			return r
		}
	}

	r, _ = gojson.Marshal(v)
	return
}

func EncodePrefix(v interface{}, prefix string, pretty ...bool) (r []byte) {
	if len(pretty) > 0 {
		r, _ = gojson.MarshalIndent(v, prefix, "  ")
		return
	}
	r, _ = gojson.Marshal(v)
	return
}

func EncodeAndStore(v interface{}, w interface{}, pretty ...bool) (r []byte, err error) {
	if len(pretty) > 0 {
		r, _ = gojson.MarshalIndent(v, "", "  ")
		return
	}
	r, _ = gojson.Marshal(v)
	switch va := w.(type) {
	case io.Writer:
		_, err = va.Write(r)
	default:
		err = Save(va.(string), r)
	}
	return
}

func Decode(data []byte, v interface{}) (err error) {
	err = gojson.Unmarshal(data, v)
	return
}

func Load(src string) (data RawValue) {

	if strings.HasPrefix(src, "file:") {
		src = strings.TrimPrefix(src, "file:")
	}
	if d, e := afero.ReadFile(fs, src); e != nil {
		return
	} else {
		data = d
	}

	return
}

func LoadWithErr(src string) (data RawValue, err error) {

	if strings.HasPrefix(src, "file:") {
		src = strings.TrimPrefix(src, "file:")
	}
	if d, e := afero.ReadFile(fs, src); e != nil {
		err = e
		return
	} else {
		data = d
	}

	return
}

func Save(src string, data []byte) (err error) {
	fmt.Println("saving json", src)
	err = afero.WriteFile(fs, src, data, os.FileMode(0755))

	return
}

func Valid(src []byte) bool {
	return gojson.Valid(src)
}

func Read(src string, v interface{}) (err error) {
	if d, e := afero.ReadFile(fs, src); e != nil {
		err = e
		return
	} else {
		Decode(d, v)
	}
	return
}

func GetFromFile(src string, path string) (g gjson.Result) {
	buf := Load(src)
	g = gjson.GetBytes(buf, path)
	return
}
