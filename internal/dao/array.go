package dao

import (
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
	"sort"

	"github.com/hyprstereo/go-dao/internal/encoding/json"
)

type Arr[T string | any] struct {
	Values []T `json:"data,omitempty"`
}

func (a *Arr[T]) GoString() string {
	return string(json.Encode(a))
}

func (a *Arr[T]) ForEach(cb func(int, T, []T)) {
	for x, v := range a.Values {
		cb(x, v, a.Values)
	}

}

func (a *Arr[T]) Filter(cb func(int, T) bool) (r []T) {
	r = make([]T, 0)
	for x, v := range a.Values {
		if cb(x, v) {
			r = append(r, v)
		}
	}
	return
}

func (a *Arr[T]) Map(cb func(int, T) T) (r []T) {
	r = make([]T, 0)
	for x, v := range a.Values {
		r = append(r, cb(x, v))
	}
	return
}

func (a *Arr[T]) MatchPattern(matcher string, result func(int, T)) (r bool) {
	for x, v := range a.Values {
		if Match(fmt.Sprint(v), matcher) {
			result(x, v)
		}
	}
	return
}

func (a *Arr[T]) Contains(matcher ...T) (r bool) {
	if len(a.Values) > 0 {
		for _, v := range a.Values {
			for _, u := range matcher {
				if reflect.DeepEqual(v, u) {
					r = true
					return
				}
			}
		}
	}
	return
}

func (a *Arr[T]) IndexOf(matcher T) (r int) {
	for x, v := range a.Values {
		if reflect.DeepEqual(v, matcher) {
			r = x
			return
		}
	}
	r = -1
	return
}

func (a *Arr[T]) Len() int {
	return len(a.Values)
}

func (a *Arr[T]) First() T {
	return a.Values[0]
}

func (a *Arr[T]) Last() T {
	return a.Values[a.Len()-1]
}

func (a *Arr[T]) Get(index int) (r T) {
	if len(a.Values) > 0 {
		r = a.Values[index]
	}
	return
}

func (a *Arr[T]) Remove(index int, count ...int) []T {
	c := 1
	if len(count) > 0 {
		c = count[0]
	}

	spliced := []T{}

	if index >= 0 && index+c < a.Len() {
		spliced = a.Values[index : index+c]
	}
	nV := []T{}
	for x, v := range a.Values {
		if x < index || x >= index+c {
			nV = append(nV, v)
		}
	}
	a.Values = nV
	return spliced
}

func (a *Arr[T]) Pop() (v T) {
	if len(a.Values) > 0 {
		spliced := a.Values[0]
		a.Values = a.Values[1:]
		v = spliced
	}
	return
}

func (a *Arr[T]) Shift() (v T) {

	if len(a.Values) > 0 {
		spliced := a.Last()
		a.Values = a.Values[0 : a.Len()-1]
		v = spliced
	}
	return
}

func (a *Arr[T]) Insert(value T, index int) *Arr[T] {
	nv := []T{}
	if index > 0 {
		if index < a.Len()-1 {
			nv = append(nv, a.Values[0:index]...)
			nv = append(nv, value)
			nv = append(nv, a.Values[index:]...)
		} else {
			nv = append(nv, a.Values...)
			nv = append(nv, value)
		}
	} else if index == 0 {
		nv = append(nv, value)
		nv = append(nv, a.Values...)
	} else {
		nv = a.Values
	}

	a.Values = nv
	return a
}

func (a *Arr[T]) Push(value T) (index int) {
	index = a.Len() - 1
	a.Values = append(a.Values, value)
	return
}

func (a *Arr[T]) Sort(order string) *Arr[T] {
	sort.Slice(a.Values, func(i, j int) bool {
		a0 := []byte(fmt.Sprint(a.Values[i]))
		a1 := []byte(fmt.Sprint(a.Values[j]))
		if order == "asc" {
			return bytes.Compare(a0, a1) < 0
		} else {
			return bytes.Compare(a0, a1) >= 0
		}
	})

	return a
}

func (a *Arr[T]) Randomise() *Arr[T] {
	cloned := &Arr[T]{}
	cloned.Values = append(cloned.Values, a.Values...)

	v := []T{}
	counter := cloned.Len() - 1
	for {
		if counter > 0 {
			i := rand.Intn(counter)
			c := cloned.Remove(i)
			v = append(v, c[0])
			counter--
		} else {
			v = append(v, cloned.First())
			break
		}

	}
	a.Values = v
	return a
}

func (a *Arr[T]) String() string {
	v := ""
	for _, val := range a.Values {
		v += fmt.Sprint(val) + ","
	}
	return v[:len(v)-1]
}

func PopulateIntArray(count int) (r []int) {
	r = make([]int, 0)
	for i := 0; i < count; i++ {
		r = append(r, i)
	}
	return
}
