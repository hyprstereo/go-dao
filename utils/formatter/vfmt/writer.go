package vfmt

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/chzyer/readline"
	"github.com/hyprstereo/go-dao/utils"
	"github.com/hyprstereo/go-dao/utils/formatter/term"
	"github.com/hyprstereo/go-dao/utils/template/ft"
)

func NewWriter(format string, t ...string) *Writer {
	w := &Writer{
		format:     format,
		styleClass: make(map[string]any),
	}

	w.SetFormat(format)
	return w
}

type Writer struct {
	io.Writer
	format     string
	wg         sync.WaitGroup
	styleClass map[string]any
}

func (w *Writer) SetStyle(clsName, style string) {
	w.styleClass[clsName] = style
}

func (w *Writer) SetFormat(format string) {
	w.format = format
}

func (w *Writer) Write(data []byte) (n int, err error) {
	res, _ := ParseAsync(utils.UnsafeString(data), "{{", "}}")
	<-res
	n, err = w.Writer.Write(data)
	return
}

func (w *Writer) HTML(input string) (r string) {
	r = Parse(input)
	return
}

func Println(v ...interface{}) (int, error) {
	for x, t := range v {
		v[x] = Parse(fmt.Sprint(t))
	}
	return fmt.Println(v...)
}

func APrintln(v ...interface{}) (int, error) {
	for x, t := range v {
		v[x] = fmt.Sprint(t)
	}
	return fmt.Println(v...)
}

func Print(v ...interface{}) (int, error) {
	for x, t := range v {
		v[x] = Parse(fmt.Sprint(t))
	}

	return fmt.Print(v...)
}

func Sprintf(format string, v ...interface{}) string {
	return Parse(fmt.Sprintf(format, v...))
}

func Printf(format string, v ...interface{}) (int, error) {
	return fmt.Printf(Sprintf(format, v...))
}

type TagFunc = func(w io.Writer, tag string) (int, error)

func tagFn(val string) TagFunc {
	return func(w io.Writer, tag string) (int, error) {
		return w.Write([]byte(val))
	}
}

type Formatter struct {
	buffer *bytes.Buffer
	lt     string
	rt     string
}

func (c *Formatter) Parse(in string) string {
	res, _ := ParseAsync(in, c.lt, c.rt)
	r := <-res
	return r
}

func (c *Formatter) Sprintf(f string, v ...any) string {
	return c.Parse(fmt.Sprintf(f, v...))
}
func (c *Formatter) Sprint(v ...any) string {
	return c.Parse(fmt.Sprint(v...))
}
func (c *Formatter) Println(v ...any) (int, error) {

	return fmt.Println(c.Sprint(v...))
}

func (c *Formatter) Print(v ...any) (int, error) {
	return fmt.Print(c.Sprint(v...))
}

func (c *Formatter) Printf(f string, v ...any) (int, error) {
	return c.Print(c.Sprintf(f, v...))
}

func (c *Formatter) fprint(f string, v ...any) (n int, err error) {
	return fmt.Fprint(c.buffer, c.Sprintf(f, v...))
}

func (c *Formatter) Output() []byte {
	return c.buffer.Bytes()
}

func (c *Formatter) Reset() {
	c.buffer.Reset()
}

func (*Formatter) Exec(in, lt, rt string, d map[string]any) string {
	return ft.ExecuteString(in, lt, rt, d)
}

var ResetAll = "\x1b[0m"

var resetFg = "\x1b[39m"
var resetBg = "\x1b[49m"

var styleMap = map[string][]string{
	"0": []string{`<\.(.*):(.*)+>`, `<b><fg#random>$1</fg></b>`},
}

var (
	Reset     = []byte("\x1b[0m")
	ResetFG   = []byte("\033[0;00m")
	Bold      = []byte("\x1b[1m")
	Dim       = []byte("\x1b[2m")
	Italic    = []byte("\x1b[3m")
	Underline = []byte("\x1b[4m")
	Blink     = []byte("\x1b[5m")
	Reverse   = []byte("\x1b[7m")
	Hidden    = []byte("\x1b[8m")
)

var (
	wg sync.WaitGroup
)

func Parse(data string, ctags ...string) string {
	ltag := "<"
	rtag := ">"
	if len(ctags) == 2 {
		ltag = ctags[0]
		rtag = ctags[1]
	}

	var res = ""
	res = parse(data, ltag, rtag)
	return res
}

func ParseAsync(data string, ctags ...string) (result chan string, err error) {
	ltag := "<"
	rtag := ">"
	if len(ctags) == 2 {
		ltag = ctags[0]
		rtag = ctags[1]
	}
	result = make(chan string)
	var res = ""
	wg.Add(1)
	go func(w *sync.WaitGroup) {

		res = parse(data, ltag, rtag)
		w.Done()

	}(&wg)
	wg.Wait()
	result <- res
	return
}

func parse(input, ltag, rtag string, styleM ...map[string]string) (res string) {

	styleMap = map[string][]string{
		"0": {`<\.(.*):(.*)+>`, `<b><fg#random>$1</fg></b>`},
	}

	lastFg := []byte("")
	lastBg := []byte("")
	isBold := false
	isItalic := false
	isUnderline := false
	isReverse := false
	listOrder := ""
	lastSeq := 0
	res = ft.ExecuteFuncString(input, ltag, rtag, func(w io.Writer, tag string) (int, error) {
		//tag = strings.TrimSpace(tag)
		newText := ""

		switch tag {
		case "cl", "clearline":
			return w.Write(clearLine)
		case "md", "movedown":
			return w.Write(moveDown)
		case "mu", "moveup":
			return w.Write(moveUp)
		case "bold", "b", "strong":
			isBold = true
			return w.Write(Bold)
		case "italic", "i", "em":
			isItalic = true
			return w.Write(Italic)
		case "underline", "u":
			isUnderline = true
			return w.Write(Underline)
		case "blink":

			return w.Write(Blink)
		case "reverse":
			isReverse = true
			return w.Write(Reverse)
		case "dim":
			return w.Write(Dim)
		case "hidden":
			return w.Write(Hidden)
		case "ul", "ol":
			listOrder = tag
			if tag == "ol" {
				lastSeq = 1
			}
		case "/ul", "/ol":
			listOrder = ""
			lastSeq = 0
			return w.Write(Reset)
		case "li":
			if listOrder == "ul" {
				newText = "  • "
			} else if listOrder == "ol" {
				newText = fmt.Sprintf("  % 2d. ", lastSeq)
				lastSeq++
			}
			return w.Write([]byte(newText))
		case "br", "br/":
			newText = "\n"
		default:

			var next []byte = Reset
			if strings.HasPrefix(tag, "/fg") {
				lastFg = []byte{}
				next = append(next, []byte(resetFg)...)
			} else if strings.HasPrefix(tag, "/bg") {
				lastBg = []byte{}
				next = append(next, []byte(resetBg)...)
			} else if strings.HasPrefix(tag, "/") {
				chk := tag
				switch chk {
				case "/b", "/strong", "/bold":
					isBold = false
				case "/u", "/underline":
					isUnderline = false
				case "/i", "/italic", "/em":
					isItalic = false
				case "/reverse":
					isReverse = false
				}

				//next = append(next, Reset...)

			} else if strings.HasPrefix(tag, "fg#") {

				r, g, b := uint8(0), uint8(0), uint8(0)
				if tag == "random" {
					r, g, b = term.RandomRGB()
				} else {
					r, g, b = term.HexToRGB(tag)
				}
				lastFg = term.GetColor(r, g, b, true)
				//return w.Write(lastFg)
			} else if strings.HasPrefix(tag, "bg#") {
				r, g, b := uint8(0), uint8(0), uint8(0)
				if tag == "random" {
					r, g, b = term.RandomRGB()
				} else {
					r, g, b = term.HexToRGB(tag)
				}
				lastFg = term.GetColor(r, g, b, true)
				lastBg = term.GetColor(r, g, b, false)
				//return w.Write(lastBg)
			}

			if isBold {
				next = append(next, Bold...)
			}
			if isItalic {
				next = append(next, Italic...)
			}
			if isUnderline {
				next = append(next, Underline...)
			}
			if isReverse {
				next = append(next, Reverse...)
			}

			if len(lastBg) > 0 {
				next = append(next, lastBg...)
			}

			if len(lastFg) > 0 {
				next = append(next, lastFg...)
			}

			return w.Write(next)
		}

		return w.Write([]byte(newText))
	})
	res += string(Reset)
	return
}

func Block(widthSize int, bg, fg, in string) (r string) {
	if widthSize < 0 {
		widthSize = readline.GetScreenWidth()
	}
	blk := make([]byte, widthSize)
	buffer := bytes.NewBuffer(blk)
	bal := widthSize - len(in)
	buf := strings.Repeat(" ", bal+1)
	buffer.WriteString(fmt.Sprintf("\r<bg#%s><fg#%s>%s</fg>%s</bg>", bg, fg, in, buf))
	r = buffer.String()
	return
}
