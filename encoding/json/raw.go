package json

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/hyprstereo/go-dao/utils"
	"github.com/hyprstereo/go-dao/utils/formatter/term"
	"github.com/hyprstereo/go-dao/utils/template/ft"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type RawValue []byte

func get(data RawValue, p string) (r Result) {
	r = Result{Result: gjson.GetBytes(data, p)}
	return
}

func getMany(data RawValue, p ...string) (r []Result) {
	res := gjson.GetManyBytes(data, p...)
	r = make([]Result, 0)
	for _, j := range res {
		r = append(r, Result{Result: j})
	}
	return
}

func set(data RawValue, p string, value any) RawValue {
	if buff, err := sjson.SetBytes(data, p, value); err != nil {
		panic(err)
	} else {
		return RawValue(buff)
	}
}

func delete(data RawValue, p string) (r RawValue) {
	if res, er := sjson.DeleteBytes(data, p); er != nil {
		panic(er)
	} else {
		r = res
	}
	return
}

func (o RawValue) String() string {
	return utils.UnsafeString(o)
}

func (o RawValue) Value() any {
	return o.Get("@this").Value()
}

func (o RawValue) Object() Result {
	return o.Get("@this")
}

// get a node by path, borrowed from gjson.Get
func (o RawValue) Get(p string) (result Result) {
	result = get(o, p)
	return
}

func (o RawValue) Keys(p string) (result Result) {
	return o.Get(p + "|@keys")
}

// get a node by path, borrowed from gjson.Get
func (o RawValue) Render(template string, key string) (r RawValue) {
	root := get(o, key)
	//fmt.Println(root.Raw)
	if root.Exists() {
		if root.IsArray() {
			result := []byte{}
			values := root.Array()
			for _, v := range values {
				tmp := ft.ExecuteFuncString(template, "{{", "}}", func(w io.Writer, tag string) (int, error) {
					t := v.Get(tag).Raw
					return w.Write([]byte(t))
				})
				result = append(result, []byte(tmp)...)
				result = append(result, []byte("\n")...)

			}
			r = result
		} else if root.IsObject() {
			result := []byte{}
			values := root.Map()
			for key, v := range values {
				tmp := ft.ExecuteFuncString(template, "{{", "}}", func(w io.Writer, tag string) (int, error) {
					if tag == "key" {
						return w.Write([]byte(key))

					}

					if strings.Contains(tag, ".") {
						tags := strings.Split(tag, ".")
						t := v.Get(tags[0])
						switch tags[1] {
						case "raw":
							return w.Write([]byte(t.Raw))
						case "str":
							return w.Write([]byte(t.Str))
						case "string":
							return w.Write([]byte(t.String()))
						case "int":
							return w.Write([]byte(fmt.Sprint(t.Int())))
						case "value":
							return w.Write([]byte(fmt.Sprint(t.Value())))
						case "type":
							return w.Write([]byte(t.Type.String()))
						case "bool":
							return w.Write([]byte(fmt.Sprint(t.Bool())))
						default:
							return w.Write([]byte(tag))
						}

					} else {
						t := v.Get(tag)
						return w.Write([]byte(fmt.Sprint(t.Value())))
					}

				})

				result = append(result, []byte(tmp)...)
				//result = append(result, []byte("\n")...)

			}
			r = result
		} else {
			result := ft.ExecuteFuncString(template, "{{", "}}", func(w io.Writer, tag string) (int, error) {
				t := get(o, tag)
				return w.Write([]byte(t.String()))
			})
			r = []byte(result)
		}
	}

	return
}

func (o RawValue) Clean() (r RawValue) {
	r = bytes.ReplaceAll(o, []byte(`"`), []byte(""))
	r = bytes.ReplaceAll(r, []byte("{"), []byte(""))
	r = bytes.ReplaceAll(r, []byte("}"), []byte(""))
	r = bytes.ReplaceAll(r, []byte(","), []byte(""))
	r = bytes.TrimSpace(r)
	return
}

func (o RawValue) Highlight(query string) (r RawValue) {
	res := get(o, query)
	if res.Exists() {
		start := res.Index
		last := start + len(res.Raw)
		chunks := o[0 : start-1]
		chunks = append(chunks, term.Hex("#ffff00")...)
		chunks = append(chunks, o[start:last]...)
		chunks = append(chunks, term.Reset...)
		chunks = append(chunks, o[last+1:]...)
		r = chunks
	}
	return
}

// get more paths, similiar to Get but with multiple path
func (o RawValue) GetMany(p ...string) (result []Result) {
	result = getMany(o, p...)
	return
}

//set new node by path
//it wont overwrite the original value unless
//it set to itself, ie: data = data.Set(path, 1234)
func (o RawValue) Set(p string, value any) RawValue {
	o = set(o, p, value)
	return o
}

func (o RawValue) Map() (m map[string]any) {
	if e := Decode(o, &m); e != nil {
		panic(e)
	}
	return
}

//delete node by path
//it wont delete the original value unless
//it set to itself, ie: data = data.Delete(path)
func (o RawValue) Delete(p string) RawValue {
	o = delete(o, p)
	return o
}

func (o RawValue) Size() (size string) {
	return utils.ByteSize([]byte(o))
}

func (o RawValue) Format(f fmt.State, verb rune) {
	fmt.Println(f, verb)
}
