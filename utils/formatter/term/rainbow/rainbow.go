package rainbow

import (
	"bytes"
	"io"

	"github.com/hyprstereo/go-dao/utils/formatter/term"
)

// Rainbow prints all the colors of the rainbow! All 232 of them colors!
type Rainbow struct {
	wrap    io.Writer
	h, s, l float64
}

func PrintRainbow(word string) string {
	rw := bytes.NewBuffer([]byte(word))
	fn := New(rw, 247, 147, 30)
	fn.Write([]byte(word))
	return rw.String()
}

// New creates a Rainbow writer, wrapping each byte with a new color
// of the rainbow. This is clearly a good idea! This is especially
// useful when w is something like a net.Conn or the output of a log.
func New(w io.Writer, r, g, b uint8) *Rainbow {
	h, s, l := term.RGBtoHSL(r, g, b)
	return &Rainbow{wrap: w, h: h, s: s, l: l}
}

// Write wraps each byte with an RGB color. Weeee!
func (r *Rainbow) Write(p []byte) (int, error) {

	buf := bytes.NewBuffer(nil)
	for i := range p {
		r.h += (0.5 / 360)
		if r.h > 1.0 {
			r.h = 0
		}

		r, g, b := term.HSLtoRGB(r.h, r.s, r.l)
		if p[i] != '\n' {
			_, _ = buf.Write(term.FgByte(p[i], r, g, b))
		} else {
			buf.Write([]byte{p[i]})
		}

	}

	_, err := buf.WriteTo(r.wrap)
	return len(p), err
}
