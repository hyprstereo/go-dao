package vfmt

import (
	"fmt"
	"strings"
)

const esc = "\033["

var (
	clearLine = []byte(esc + "2K\r")
	moveUp    = []byte(esc + "1A")
	moveDown  = []byte(esc + "1B")
)

type FormatSymbol uint8

const (
	ResetSymbol FormatSymbol = iota
	DimSymbol
	ResetFGSymbol
	ResetBGSymbol
	BoldSymbol
	UnderlineSymbol
	ItalicSymbol
	ReverseSymbol
	StrikeOutSymbol
	NewlineSymbol
	ColorSymbol
	BGColorSymbol
	StyleSymbol
	BlinkSymbol
	HiddenSymbol
	ClearLineSymbol
	MovedownSymbol
	MoveupSymbol
	ListOrderedSymbol
	ListUnorderedSymbol
	ListItemSymbol
)

type SymValue struct {
	Sym   FormatSymbol
	Value []byte
}

type SymbolsMap map[FormatSymbol]string

func (s SymbolsMap) GetValueOf(sym FormatSymbol) (result any, ok bool) {
	res, ok0 := s[sym]
	if ok0 && strings.Contains(res, ",") {
		result = strings.Split(res, ",")[0]
	} else {
		result = res
	}
	ok = ok0
	return
}

type IFormatter interface {
	Init() error
	GetSymbols() SymbolsMap
	IsSymbolOf(string) (FormatSymbol, bool)
}

var htmlSet = SymbolsMap{
	DimSymbol:           "dim",
	ColorSymbol:         "fg",
	ResetFGSymbol:       "/fg",
	BGColorSymbol:       "bg",
	ResetBGSymbol:       "/bg",
	BoldSymbol:          "b,strong,bold",
	UnderlineSymbol:     "u,underline",
	ItalicSymbol:        "i,italic,em",
	NewlineSymbol:       "br",
	StyleSymbol:         "style",
	ReverseSymbol:       "reverse",
	ClearLineSymbol:     "cl,clearline",
	MovedownSymbol:      "md,movedown",
	MoveupSymbol:        "mu,moveup",
	BlinkSymbol:         "blink",
	ListOrderedSymbol:   "ol",
	ListUnorderedSymbol: "ul",
	ListItemSymbol:      "li",
}

type htmlFormatter struct {
	IFormatter
	symbols *SymbolsMap
}

func (h htmlFormatter) Init() (err error) {
	if *h.symbols != nil || len(*h.symbols) > 0 {
		err = fmt.Errorf("already_init")
		return
	}
	h.symbols = &htmlSet
	return
}

func (h htmlFormatter) GetSymbols() SymbolsMap {
	return *h.symbols
}

func (h htmlFormatter) IsSymbolOf(tags string) (sym FormatSymbol, ok bool) {
	for s, els := range *h.symbols {
		tags = strings.TrimPrefix(tags, "/")
		if strings.ContainsAny(els, tags) {
			sym = s
			ok = true
			return
		}
	}
	return
}
