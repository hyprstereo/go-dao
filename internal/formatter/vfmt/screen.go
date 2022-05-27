package vfmt

import (
	"bytes"
	"fmt"
	"os"

	"github.com/chzyer/readline"
)

type Screen struct {
	Width   int
	Height  int
	bgColor []byte
	color   []byte
	//screenBuffer *screenbuf.ScreenBuf
	//bufs    [][]byte
	writer  *bytes.Buffer
	rb      *readline.RuneBuffer
	rl      *readline.Config
	linePos int
}

func (c *Screen) Init() {
	//c.bufs = make([][]byte, c.Height)
	c.writer = bytes.NewBuffer([]byte{})
	c.Clear()
}
func (c *Screen) MoveDown() {
	if c.linePos > 0 {
		c.linePos++
	}
}
func (c *Screen) MoveUp() {
	if c.linePos > 0 {
		c.linePos--

	}
}

func (c *Screen) Clear() {
	for l := 0; l < c.Height; l++ {
		//c.bufs[l] = c.blockBuffer()
		c.writer.Write(c.blockBuffer())
		c.writer.Write(moveUp)
	}
	c.linePos = 0
	fmt.Fprintf(os.Stdout, "%s", c.writer.String())
}

func (c *Screen) Write(str string, startPos ...int) (int, error) {
	blk := c.blockBuffer(str)
	return c.writer.Write(blk)
}

func (c *Screen) blockBuffer(txt ...string) []byte {
	line := []byte{}
	line = append(line, c.bgColor...)
	line = append(line, c.color...)
	if len(txt) > 0 {
		line = append(line, []byte(txt[0])...)
	}
	buf := bytes.Repeat([]byte(" "), c.Width-len(line))
	line = append(line, buf...)
	line = append(line, Reset...)
	return buf
}

func NewScreen(width, height int, fgColor, bgColor []byte) *Screen {
	w := bytes.NewBuffer([]byte{})
	s := &Screen{
		Width:   width,
		Height:  height,
		bgColor: bgColor,
		color:   fgColor,
		//screenBuffer: screenbuf.New(w),
		writer: w,
	}

	s.Init()
	return s
}
