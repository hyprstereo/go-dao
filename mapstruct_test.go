package godao

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestMapToStruct(t *testing.T) {
	var m Map = Map{
		"id":        int(0),
		"name":      string(""),
		"createdAt": time.Time(time.Now()),
		"is_active": bool(true),
	}
	v := MapToStruct(m, "json")
	newStruct := reflect.New(v).Elem()

	newStruct.Field(0).SetInt(1)
	newStruct.Field(1).SetString("Name")
	newStruct.Field(2).Set(reflect.ValueOf(time.Now()))
	newStruct.Field(3).SetBool(true)

	d := newStruct.Addr().Interface()
	str := strings.ReplaceAll(reflect.TypeOf(d).String(), `"j`, "`j")
	str = strings.ReplaceAll(str, `" `, "` ")
	str = strings.ReplaceAll(str, `\"`, `"`)
	println(str)
}
