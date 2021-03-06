package template

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	godao "github.com/hyprstereo/go-dao"

	"github.com/hyprstereo/go-dao/utils/template/ft"
)

type FuncMaps = godao.Map
type Funcs = func(args ...any) (result any, err error)

type Writer struct {
	*bytes.Buffer
	linePos   int
	cursorPos int
	lines     []godao.Bytes
	lastSize  int
}

func (w *Writer) Init(size int) (err error) {
	sizeBuf := make([]byte, size)
	w.Buffer = bytes.NewBuffer(sizeBuf)
	return
}

func (w *Writer) Lines() []godao.Bytes {
	if len(w.lines) < 1 || w.lastSize != len(w.Bytes()) {

		res := []godao.Bytes{}
		rdr := bufio.NewReader(w)
		isEof := false
		for !isEof {
			if line, _, err := rdr.ReadLine(); err != nil {
				isEof = true
				break
			} else {
				res = append(res, godao.Bytes(line))
			}
		}
		w.lines = res
		w.lastSize = len(w.Bytes())
	}
	return w.lines
}

func (w *Writer) Clear() {
	w.lines = []godao.Bytes{}
	w.Reset()
	w.cursorPos = 0
	w.linePos = 0
}

func (w *Writer) LastLine() godao.Bytes {
	return w.lines[len(w.lines)-1]
}

func (w *Writer) ReadLine(lineIndex int) godao.Bytes {
	if lineIndex < len(w.Lines()) {
		return w.lines[lineIndex]
	}
	return nil
}

func (w *Writer) Output() godao.Bytes {
	return godao.Bytes(w.Bytes())
}

func (w *Writer) Render(dst io.Writer, mapData ...FuncMaps) (p int, err error) {
	if bytes.ContainsAny(w.Bytes(), "${}") {
		dataMap := FuncMaps{}.Merge(mapData...)
		result := ft.ExecuteFuncString(w.Output().String(), "${", "}", func(wr io.Writer, tag string) (int, error) {
			tag = strings.TrimSpace(tag)
			var tokens []string
			if strings.Contains(tag, " ") {
				tokens = strings.Split(tag, " ")
			} else {
				tokens = []string{tag}
			}
			if value, ok := dataMap[tokens[0]]; ok {
				switch fn := value.(type) {
				case Funcs:
					if res, er := fn(extractArgs(tokens[1:]...)...); er == nil {
						return wr.Write([]byte(fmt.Sprint(res)))
					}
				case ft.TagFunc:
					return fn(wr, tag)
				default:
					return wr.Write([]byte(fmt.Sprint(fn)))
				}
			}
			return -1, nil
		})
		//p, err = dst.Write(w.Bytes())
		return dst.Write([]byte(result))
	} else {
		return dst.Write(w.Bytes())
	}
}

func extractArgs(str ...string) []any {
	res := []any{}
	for _, v := range str {
		res = append(res, v)
	}
	return res
}
