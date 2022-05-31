package term

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/hyprstereo/go-dao/utils/template/ft"
)

var (
	//Reset     = []byte("\x1b[0m")
	ResetFG   = []byte("\033[0;00m")
	Bold      = []byte("\x1b[1m")
	Dim       = []byte("\x1b[2m")
	Italic    = []byte("\x1b[3m")
	Underline = []byte("\x1b[4m")
	Blink     = []byte("\x1b[5m")
	Reverse   = []byte("\x1b[7m")
	Hidden    = []byte("\x1b[8m")
)

type Fmt struct {
}

func Printf(w io.Writer, format string, v ...interface{}) {
	entryStr := fmt.Sprintf(format, v...)
	str := ft.ExecuteFuncString(entryStr, "<%", "%>", termFuncHandler)
	w.Write([]byte(str))
}

func Sprintf(format string, v ...interface{}) string {
	entryStr := fmt.Sprintf(format, v...)
	str := ft.ExecuteFuncString(entryStr, "<%", "%>", termFuncHandler)
	return InterpretStr(str)
}

func Println(format string, v ...interface{}) {
	entryStr := fmt.Sprintf(format, v...)
	str := InterpretStr(ft.ExecuteFuncString(entryStr, "<%", "%>", termFuncHandler))
	fmt.Println(str)
}

func termFuncHandler(w io.Writer, tag string) (p int, err error) {
	v := []byte(strings.TrimSpace(tag))
	if bytes.HasPrefix(v, []byte("c:[")) {
		v = bytes.TrimPrefix(v, []byte("c:["))
		v = bytes.TrimSuffix(v, []byte("]"))
		toks := bytes.Split(v, []byte(","))
		val := FgString(string(toks[0]), strTouint8(toks[1][:]), strTouint8(toks[2]), strTouint8(toks[3]))
		return w.Write([]byte(val))
	}
	return
}

func strTouint8(s []byte) uint8 {
	v, _ := strconv.ParseUint(string(s), 10, 8)
	return uint8(v)
}
