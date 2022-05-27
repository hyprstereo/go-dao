package godao

import "github.com/hyprstereo/go-dao/internal/dao"

// example of speficic type of array
type StringArray = dao.Arr[string]

// Array is a dynamic type of array
type Array = dao.Arr[any]

// create a new Array.
func NewArray(initialValue ...[]any) *Array {
	initValue := []any{}
	if len(initialValue) > 0 {
		initValue = initialValue[0]
	}
	return &Array{Values: initValue}
}

// generic array
func NewArrayType[T comparable](initialValue ...T) (arr *dao.Arr[T]) {
	arr = &dao.Arr[T]{}
	if len(initialValue) > 0 {
		arr.Values = initialValue
	}
	return
}
