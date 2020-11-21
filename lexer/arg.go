package lexer

import (
	"log"
	"strconv"
	"unicode"
)

type Arg struct {
	rawValue string
	Value    string
	IsDefine bool
	IsLabel  bool
}

func parseInt(s string, base int) uint64 {
	n, err := strconv.ParseUint(s, 0, base)
	if err != nil {
		log.Fatalf("failed to convert string %s to integer", s)
	}

	return n
}

func (a *Arg) AsUint() uint64 {
	return parseInt(a.Value, 16)
}

func (a *Arg) AsUint8() uint8 {
	return uint8(parseInt(a.Value, 16))
}

func (a *Arg) init() {
	rv := a.rawValue
	a.IsDefine = rv[0] == '$'
	a.IsLabel = unicode.IsLetter(rune(rv[0]))

	switch {
	case a.IsDefine:
		a.Value = rv[1:len(rv)]
	case a.IsLabel:
		a.Value = rv
	default:
		a.Value = rv
	}
}

func newArg(value string) Arg {
	a := Arg{rawValue: value}
	a.init()
	return a
}

func newArgs(values ...string) []Arg {
	args := []Arg{}

	for _, v := range values {
		args = append(args, newArg(v))
	}

	return args
}
