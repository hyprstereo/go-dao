package term

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

// stolen from gorilla color:
//    https://code.google.com/p/gorilla/source/browse/?r=ef489f63418265a7249b1d53bdc358b09a4a2ea0#hg%2Fcolor

func Stylize(in []byte, format ...[]byte) (res []byte) {
	res = packMutilbytes(format...)
	res = append(append(res, in...), ResetFG...)

	return
}

func packMutilbytes(byt ...[]byte) []byte {
	res := []byte{}
	for _, by := range byt {
		res = append(res, by...)
	}
	return res
}

func HexToRGB(hex string) (r, g, b uint8) {
	if strings.Contains(hex, "#") {
		hex = hex[strings.Index(hex, "#")+1:]
	}
	tmphex := hex + "000000"

	ur, _ := strconv.ParseUint(tmphex[0:2], 16, 64)
	ug, _ := strconv.ParseUint(tmphex[2:4], 16, 64)
	ub, _ := strconv.ParseUint(tmphex[4:6], 16, 64)
	r = uint8(ur)
	g = uint8(ug)
	b = uint8(ub)

	return
}

func Hex(hex string) []byte {
	r, g, b := HexToRGB(hex)
	return rgb(r, g, b, uint8(0), uint8(0), uint8(0))
}

func RandomRGB() (r, g, b uint8) {
	r = uint8(rand.Uint32() * 255)
	g = uint8(rand.Uint32() * 255)
	b = uint8(rand.Uint32() * 255)
	return
}

func RandomHex() string {
	r := uint8(rand.Uint32() * 255)
	g := uint8(rand.Uint32() * 255)
	b := uint8(rand.Uint32() * 255)
	return fmt.Sprintf("%02x%02x%02x", r, g, b)
}

func RandomHSL() (h, s, l float64) {
	h = (rand.Float64() * 255) / 255
	s = (rand.Float64() * 255) / 255
	l = (rand.Float64() * 255) / 255
	return
}

// RGBtoHSL is a helper that transforms 8bits HSL colors to
// 8bits RGB colors.
func RGBtoHSL(r, g, b uint8) (h, s, l float64) {
	fR := float64(r) / 255
	fG := float64(g) / 255
	fB := float64(b) / 255
	max := math.Max(math.Max(fR, fG), fB)
	min := math.Min(math.Min(fR, fG), fB)
	l = (max + min) / 2
	if max == min {
		// Achromatic.
		h, s = 0, 0
	} else {
		// Chromatic.
		d := max - min
		if l > 0.5 {
			s = d / (2.0 - max - min)
		} else {
			s = d / (max + min)
		}
		switch max {
		case fR:
			h = (fG - fB) / d
			if fG < fB {
				h += 6
			}
		case fG:
			h = (fB-fR)/d + 2
		case fB:
			h = (fR-fG)/d + 4
		}
		h /= 6
	}
	return
}

// HSLtoRGB is a helper that transforms 8bits RGB colors to
// 8bits HSL colors.
func HSLtoRGB(h, s, l float64) (r, g, b uint8) {
	var fR, fG, fB float64
	if s == 0 {
		fR, fG, fB = l, l, l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - s*l
		}
		p := 2*l - q
		fR = hueToRGB(p, q, h+1.0/3)
		fG = hueToRGB(p, q, h)
		fB = hueToRGB(p, q, h-1.0/3)
	}
	r = uint8((fR * 255) + 0.5)
	g = uint8((fG * 255) + 0.5)
	b = uint8((fB * 255) + 0.5)
	return
}

// hueToRGB is a helper function for HSLtoRGB.
func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	if t < 1.0/6 {
		return p + (q-p)*6*t
	}
	if t < 0.5 {
		return q
	}
	if t < 2.0/3 {
		return p + (q-p)*(2.0/3-t)*6
	}
	return p
}
