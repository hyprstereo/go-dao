package dao

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Type = reflect.Type
type Kind = reflect.Kind
type Value[T any] struct {
	Val any
}

func (v Value[T]) As(v2 any) {
	v.Val = v2
}

func (v Value[T]) Type() (r Type) {
	return reflect.TypeOf(v)
}

func (v Value[T]) ToInt() int {
	return v.Val.(int)
}

func (v Value[T]) ToString() string {
	return v.Val.(string)
}

func (v Value[T]) ToBool() bool {
	return v.Val.(bool)
}

func (v Value[T]) ToFloat() float64 {
	return v.Val.(float64)
}

func (v Value[T]) GetType() string {
	return reflect.TypeOf(v.Val).String()
}

func (v Value[T]) ToBytes() BytesValue {
	return BytesValue([]byte(fmt.Sprint(v.Val)))
}

func GetFieldValue(key string, exprs string, defaultValue ...interface{}) (r any) {
	if strings.Contains(exprs, key) {
		pat := fmt.Sprintf(`%s=[\"\'](.*?)[\"\']`, key)
		reg := regexp.MustCompile(pat)
		m := reg.FindStringSubmatch(exprs)
		if len(m) == 2 {
			r = m[1]
			return
		}
		if len(defaultValue) > 0 {
			r = defaultValue[0]
			return
		}
	}
	return
}
