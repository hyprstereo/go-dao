package dao

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/tidwall/gjson"

	"github.com/hyprstereo/go-dao/encoding/hjson"
	"github.com/hyprstereo/go-dao/encoding/json"
	"github.com/hyprstereo/go-dao/internal/utils"

	"github.com/tidwall/sjson"
)

var mu sync.RWMutex

type Result struct {
	gjson.Result
}

func (r Result) fromBytes(data []byte) (rs Result) {
	r = Result{Result: gjson.ParseBytes(data)}
	rs = r
	return
}

func (r Result) Get(key string) Result {
	return Result{Result: r.Result.Get(key)}
}

func (r Result) Bytes() []byte {
	return []byte(utils.UnsafeBytes(r.Raw))
}

func (r Result) Decode(v any) (err error) {
	err = json.Decode([]byte(r.Raw), v)
	return
}

func (r Result) Map() Map {
	return Map(r.Value().(map[string]any))
}

func (r Result) Set(p string, value any) (err error) {
	if d, e := sjson.SetBytes(r.Bytes(), p, value); e != nil {
		err = e
		return
	} else {
		r.fromBytes(d)
	}

	return
}

func (r Result) Delete(path string) (err error) {
	if d, e := sjson.DeleteBytes(r.Bytes(), path); e != nil {
		err = e
		return
	} else {
		r.fromBytes(d)
	}
	return
}

func (r Result) Each(cb func(key, value any) bool) {
	r.ForEach(func(key, value gjson.Result) bool {
		return cb(key.Value(), value.Value())
	})
}

type SMap[T string, V any] struct {
	V map[T]V
}

func (s *SMap[T, V]) String(pretty ...bool) string {
	return string(s.Bytes(pretty...))
}

func (s *SMap[T, V]) Bytes(pretty ...bool) []byte {
	return json.Encode(&s.V, pretty...)
}

func (s *SMap[T, V]) Parse(val []byte) *SMap[T, V] {
	json.Decode(val, &s.V)
	return s
}

func (s *SMap[T, V]) Get(path string) (val Result) {
	res := gjson.GetBytes(s.Bytes(), path)
	val = Result{Result: res}
	return
}

func (s *SMap[T, V]) GetMany(path ...string) (val []Result) {
	res := gjson.GetManyBytes(s.Bytes(), path...)
	val = make([]Result, 0)
	for _, v := range res {
		val = append(val, Result{Result: v})
	}
	return
}

func (s *SMap[T, V]) Set(path string, val V) {
	buffer, _ := sjson.SetBytesOptions(s.Bytes(), path, val, &sjson.Options{Optimistic: true})
	json.Decode(buffer, &s)
}

func (s *SMap[T, V]) Delete(path string) {
	buffer, _ := sjson.DeleteBytes(s.Bytes(), path)
	json.Decode(buffer, &s.V)
}

type Map map[string]any

func (m Map) FromBytes(data []byte) (r Result) {
	r = Result{Result: gjson.ParseBytes(data)}
	return
}

func (m Map) FromString(data string) (r Result) {
	return m.FromBytes([]byte(utils.UnsafeString([]byte(data))))
}

func (m Map) EncodeBinary(additional ...map[string][]rune) (data []byte, err error) {
	buf := bytes.NewBuffer([]byte{})
	enc := gob.NewEncoder(buf)
	err = enc.Encode(&m)
	if err == nil {
		data = buf.Bytes()
	}
	return
}

func (m Map) DecodeBinary(data []byte) (mp Map, err error) {
	buf := bytes.NewReader(data)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&m)
	if err == nil {
		mp = m
	}
	return
}

func (m Map) Length() int {
	mu.RLock()
	defer mu.RUnlock()
	return len(m)
}
func (m Map) Each(fn func(int, string, any) bool) (a []any) {
	mu.Lock()

	index := 0
	for k, v := range m {
		if fn(index, k, v) {
			a = append(a, v)
		}
	}
	mu.Unlock()
	return
}

func (m Map) Selective(keys ...string) (o Map) {
	mu.Lock()
	defer mu.Unlock()
	o = make(Map)
	for k, val := range m {
		o[k] = val
	}
	return
}

func (m Map) Keys() []string {
	mu.Lock()
	ky := []string{}
	for k, _ := range m {
		ky = append(ky, k)
	}
	mu.Unlock()
	return ky
}

func (m Map) Values() *Arr[any] {
	mu.Lock()
	ky := &Arr[any]{Values: make([]any, 0)}
	for _, v := range m {
		ky.Push(v)
	}
	mu.Unlock()
	return ky
}

func (m Map) ValuesString() []string {
	mu.Lock()
	ky := []string{}
	for _, v := range m {
		ky = append(ky, fmt.Sprint(v))
	}
	mu.Unlock()
	return ky
}

func (m Map) CanCall(key string) (r bool) {
	r = reflect.TypeOf(m.Get(key).Value()).Kind() == reflect.Func
	return
}

//special func
func (m Map) Call(key string, args ...any) (out any, err error) {
	r := reflect.ValueOf(m.Get(key).Value())
	if r.Kind() == reflect.Func {
		ins := []reflect.Value{}
		for _, e := range args {
			ins = append(ins, reflect.ValueOf(e))
		}
		res := r.Call(ins)
		if len(res) > 0 {
			if len(res) == 1 {
				out = res[0].Interface()
			} else {
				o := []any{}
				for _, v := range res {
					o = append(o, v.Interface())
				}
				out = o
			}
			return
		}
	} else {
		err = fmt.Errorf("%s is not callable", key)
	}
	return
}

func (m Map) GetByIndex(i int) (k string, v any) {
	mu.RLock()
	c := 0
	for _, v := range m {
		if c == i {
			return k, v
		}
	}
	mu.RUnlock()
	return "", nil
}

func (m Map) MatchesKey(pattern string, res func(k string, v any)) (ok bool) {
	mu.RLock()
	defer mu.RUnlock()
	for f, v := range m {
		if Match(f, pattern) {
			res(f, v)
			ok = true
			return
		}
	}
	return
}

func (m Map) KeyMatch(pattern string) (key string, value any, ok bool) {
	mu.RLock()
	defer mu.RUnlock()
	for f, v := range m {
		if Match(f, pattern) {
			key = f
			value = v
			ok = true
			return
		}
	}
	return
}

func (m Map) ForEach(cb func(int, string, any)) {
	mu.RLock()
	defer mu.RUnlock()

	cnt := 0
	for k, n := range m {
		cb(cnt, k, n)
		cnt++
	}
}

func (m Map) ByKeyOrder(sel []string, cb func(int, string, any)) {
	mu.RLock()
	defer mu.RUnlock()
	arrs := &Arr[string]{Values: sel}
	arrs.Map(func(i int, s string) string {
		if m.Has(s) {
			cb(i, s, m.Get(s).Value())
		}
		return s
	})
}

func (m Map) EqualKeys(src Map) (r bool) {
	mu.RLock()
	defer mu.RUnlock()

	k := m.Keys()
	//k2 := src.Keys()
	for _, v := range k {
		if !src.Has(v) {
			return
		}
	}
	r = true
	return
}

func (m Map) Diffs(src Map) (diff []string) {
	mu.RLock()
	defer mu.RUnlock()

	diff = make([]string, 0)
	k := m.Keys()
	k2 := src.Keys()

	sort.Strings(k)
	sort.Strings(k2)

	diffKeys := []string{}
	for x, v := range k {
		if v != k2[x] {
			diffKeys = append(diffKeys, v)
		}
	}
	diff = diffKeys
	return
}

func (m Map) Merge(s ...Map) Map {
	for _, m2 := range s {
		for n, v := range m2 {
			m[n] = v
		}
	}
	return m
}

func (m Map) MergeStrict(s Map, excludes []string) Map {
	ex := strings.Join(excludes, " ")
	for n, v := range s {
		if !strings.Contains(ex, n) {
			m[n] = v
		} else {
			fmt.Printf("%s field is restricted", n)
		}
	}
	return m
}

func (m Map) Consume(v map[string]any) Map {
	m = Map(v)
	return m
}

func (m Map) Has(keys ...string) bool {
	mu.RLock()
	defer mu.RUnlock()
	for _, k := range keys {
		if _, ok := m[k]; !ok {
			return false
		}
	}
	return true
}

func (m Map) Map(key string) Map {
	if m.Get(key).Exists() {
		//fmt.Println(m[key])
		return Map(m[key].(Map))
	}
	return nil
}

func (m Map) String(key string) string {
	mu.RLock()
	defer mu.RUnlock()
	r := m.Get(key)
	if r.Type == gjson.String {
		return r.String()
	}
	return fmt.Sprint(r.Value())
}

func (m Map) Int(key string) int {
	mu.RLock()
	defer mu.RUnlock()
	r := m.Get(key)
	if r.Type == gjson.Number {
		return int(r.Int())
	}
	return -1
}

func (m Map) Bool(key string) bool {
	mu.RLock()
	defer mu.RUnlock()
	r := m.Get(key)
	if r.Type == gjson.True || r.Type == gjson.False {
		return r.Bool()
	}
	return false
}

func (m Map) Dom(tag string, inner ...string) (res string) {
	mu.Lock()
	defer mu.Unlock()
	attrs := []string{}
	for key, val := range m {
		attrs = append(attrs, key+`="`+fmt.Sprint(val)+`"`)
	}
	dom := fmt.Sprintf("<%[1]s %[2]s>%[3]s</%[1]s>", tag, strings.Join(attrs, " "), strings.Join(inner, " "))
	res = dom
	return
}

func (m Map) JSON(keys ...string) any {
	mu.RLock()
	defer mu.RUnlock()
	if len(keys) > 0 {
		dat := gjson.GetManyBytes(m.Bytes(), keys...)
		res := []string{}
		for _, d := range dat {
			res = append(res, fmt.Sprint(d))
		}
		return res
	}

	return string(m.Bytes())
}

func (m Map) JSON_Indent(keys ...string) any {
	mu.RLock()
	defer mu.RUnlock()
	if len(keys) > 0 {
		dat := gjson.GetManyBytes(m.Bytes(), keys...)
		res := []string{}
		for _, d := range dat {
			res = append(res, fmt.Sprint(d))
		}
		return res
	}

	return string(json.Encode(&m, true))
}

func (m Map) HJSON(keys ...string) interface{} {
	if len(keys) > 0 {
		dat := gjson.GetManyBytes(m.Bytes(), keys...)
		res := []string{}
		for _, d := range dat {
			buf, _ := hjson.Marshal(d.Value())
			res = append(res, string(buf))
		}
		return res
	}
	byt, _ := hjson.Marshal(&m)
	return string(byt)
}

func (m Map) Get(p string, defaultValue ...any) (res Result) {
	mu.RLock()
	defer mu.RUnlock()
	if len(m) > 0 {
		if p == "" {
			p = "@this"
		}
		res = Result{Result: gjson.GetBytes(m.Bytes(), p)}
	}
	if len(defaultValue) > 0 && !res.Exists() {
		data, _ := sjson.Set("{}", p, defaultValue[0])
		res.Value()
		res = Result{Result: gjson.Parse(data)}
	}
	return
}

func (m Map) Interface(p string) (res any) {
	mu.RLock()
	defer mu.RUnlock()
	res = m[p]
	return
}

func (m Map) Set(p string, value interface{}) (res any) {
	mu.Lock()
	defer mu.Unlock()
	buf, er := sjson.SetBytes(m.Bytes(), p, value)
	if er != nil {
		//Result(gjson.ParseBytes(json.Encode(Map{"Error": er.Error()})))
	} else {
		json.Decode(buf, &m)
		//res = Result(gjson.ParseBytes(buf))
		res = string(m.Bytes())
	}
	return
}

func (m Map) Del(p string) (res Result) {
	mu.Lock()
	defer mu.Unlock()
	buf, er := sjson.DeleteBytes(m.Bytes(), p)
	if er != nil {
		res = Result{Result: gjson.ParseBytes(json.Encode(Map{"Error": er.Error()}))}
	} else {
		json.Decode(buf, &m)
		res = Result{Result: gjson.ParseBytes(buf)}
	}
	return
}

func (m Map) GetAsArray(jpath string) (res []Map) {
	mu.RLock()
	defer mu.RUnlock()

	j := gjson.GetBytes(m.Bytes(), jpath)
	fmt.Println(j.Raw)
	if j.IsArray() {
		for _, r := range j.Array() {
			nmap := Map{}
			jMap := r.Map()
			for k, f := range jMap {
				nmap[k] = f.Value()
			}
			res = append(res, nmap)
		}
	}
	return
}

func (m Map) Bytes(keys ...string) json.RawValue {
	buf := json.Encode(&m)
	if len(keys) > 0 {
		if len(keys) > 1 {
			res := gjson.GetBytes(buf, keys[0])
			buf = json.Encode(&res)
		} else {
			res := gjson.GetManyBytes(buf, keys...)
			buf = json.Encode(&res)
		}

	}
	return buf
}

func (m Map) Iterate(from string, cb func(string, Result) bool) (result Map) {
	mu.RLock()
	defer mu.RUnlock()
	result = make(Map)
	if i := m.Get(from); i.Exists() {
		i.ForEach(func(key, value gjson.Result) bool {
			ret := cb(key.Str, Result{Result: value})
			if ret {
				result[key.Str] = value.Value()
			}
			return ret
		})
	}

	return
}

// decode the values into any variable
func (m Map) Decode(key string, dst any) (err error) {
	d := m.Get(key)
	if d.Exists() {
		switch d.Type {
		case gjson.Number:
			dst = d.Num
		case gjson.String:
			dst = d.String()
		case gjson.False, gjson.True:
			dst = d.Bool()
		case gjson.JSON:
			err = json.Decode([]byte(d.Raw), dst)
		default:
			err = json.Decode([]byte(d.Raw), dst)
		}

	}
	return
}

func MapKeys[K comparable, V any](m map[K]V) (s []K) {
	s = make([]K, 0)
	for k, _ := range m {
		s = append(s, k)
	}
	return
}

func MapValues[K comparable, V comparable](m map[K]V) (s []V) {
	s = make([]V, 0)
	for _, v := range m {
		s = append(s, v)
	}
	return
}
