package lexer

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Arg struct {
	rawValue          string
	Value             string
	IsIdentifier      bool
	IsFPOffsetAddress bool
	IsRegister        bool
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

var identifierRegex = regexp.MustCompile(`\w+`)

func (a *Arg) init() {
	rv := a.rawValue

	// TODO: This won't work for other offset addressing modes
	a.IsFPOffsetAddress = rv[0] == '+' || rv[0] == '-'

	a.IsRegister = len(rv) == 1 && (rv == "A" || rv == "B")

	if !a.IsRegister && !a.IsFPOffsetAddress && len(rv) >= 2 && rv[0] != '0' && rv[1] != 'x' {
		a.IsIdentifier = identifierRegex.MatchString(rv)
	}

	switch {
	case a.IsIdentifier:
		a.Value = rv
	case a.IsFPOffsetAddress:
		// hack to get this working :(
		sign := rv[0]
		value := strings.Split(rv[1:], "(fp)")[0] // decimal number string
		valueInt := parseInt(value, 10)

		if sign == '-' {
			n := int8(valueInt) * int8(-1)
			a.Value = fmt.Sprintf("0x%x", uint8(n))
		} else {
			a.Value = fmt.Sprintf("0x%x", valueInt)

		}

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
