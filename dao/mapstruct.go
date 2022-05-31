package dao

import (
	"reflect"
	"time"

	"github.com/iancoleman/strcase"
)

var (
	preType = map[string]reflect.Type{
		"BOOL":      reflect.TypeOf(bool(true)),
		"INT":       reflect.TypeOf(int32(0)),
		"FLOAT":     reflect.TypeOf(float64(0.1)),
		"STRING":    reflect.TypeOf(string("")),
		"TIME":      reflect.TypeOf(time.Now()),
		"TIMESTAMP": reflect.TypeOf(int64(0)),
		"MAP":       reflect.TypeOf(map[string]any{}),
		"ARRAY":     reflect.TypeOf([]any{}),
		"BLOB":      reflect.TypeOf([]byte{}),
		"RUNE":      reflect.TypeOf(rune('0')),
	}
)

type MStructConfig struct {
	Tag string
}

func MapToStruct(src Map, tags string) (newStruct reflect.Type) {
	sFields := []reflect.StructField{}
	for name, field := range src {
		st := reflect.StructField{
			Name: strcase.ToCamel(SafeString(name)),
			Type: reflect.TypeOf(field),
			Tag:  reflect.StructTag(`json:"` + strcase.ToLowerCamel(name) + `"`),
		}
		sFields = append(sFields, st)
	}
	newStruct = reflect.StructOf(sFields)
	return
}
