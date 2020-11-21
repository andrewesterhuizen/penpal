package assembler

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/andrewesterhuizen/vm/instructions"
)

// TODO: labels

type Assembler struct {
	index        int
	source       string
	instructions []uint8
	defines      map[string]string
	labels       map[string]uint16
}

func New() Assembler {
	a := Assembler{}
	a.defines = make(map[string]string)
	a.labels = make(map[string]uint16)
	return a
}

var instructionRegex = regexp.MustCompile(`(\$*\w+)`)
var defineRegex = regexp.MustCompile(`\$(\w+)\s?=\s?(\w+)`)

var wordRegex = regexp.MustCompile(`(\w+)`)

func getInstructionAndArg(line string) (string, string, bool, bool) {
	matches := instructionRegex.FindAllString(line, -1)

	if matches == nil || len(matches) <= 0 {
		return "", "", false, false
	}

	if len(matches) == 2 {
		if matches[1][0] == '$' {
			return matches[0], matches[1], true, false
		} else if matches[1][0] == '0' && matches[1][1] == 'x' {
			return matches[0], matches[1], false, false
		} else if wordRegex.MatchString(matches[1]) {
			return matches[0], matches[1], false, true
		} else {
			return matches[0], matches[1], false, false
		}

	}

	return matches[0], "", false, false
}

func parseInt(s string, base int, instruction string) uint64 {
	n, err := strconv.ParseUint(s, 0, base)
	if err != nil {
		log.Fatalf("failed at instruction %s: unable to parse int from value %s", instruction, s)
	}

	return n
}

func (a *Assembler) ParseLine(line string) {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "$") {
		a.parseDefine(line)

	} else if strings.Contains(line, ":") {
		// skip
	} else {
		a.parseInstruction(line)
	}

}

func (a *Assembler) parseDefine(line string) {
	matches := instructionRegex.FindAllString(line, -1)
	a.defines[matches[0]] = matches[1]
}

func (a *Assembler) getDefine(d string) string {
	v, exists := a.defines[d]
	if !exists {
		log.Fatalf("no definition found for %s", d)
	}

	return v
}

func (a *Assembler) getLabel(d string) uint16 {
	v, exists := a.labels[d]
	if !exists {
		log.Fatalf("no label definition found for %s", d)
	}

	return v
}

func (a *Assembler) parseInstruction(line string) {
	line = strings.TrimSpace(line)

	if line == "" {
		return
	}

	instruction, arg, argIsDefine, argIsLabel := getInstructionAndArg(line)

	switch instruction {
	case "MOVA":
		a.instructions = append(a.instructions, instructions.MOVA)

		if argIsDefine {
			value := a.getDefine(arg)
			a.instructions = append(a.instructions, uint8(parseInt(value, 16, instruction)))
		} else if argIsLabel {
			value := a.getLabel(arg)
			a.instructions = append(a.instructions, uint8((value&0xff00)>>8))
			a.instructions = append(a.instructions, uint8((value & 0xff)))
		} else {
			a.instructions = append(a.instructions, uint8(parseInt(arg, 8, instruction)))
		}

	case "MOVB":
		a.instructions = append(a.instructions, instructions.MOVB)

		if argIsDefine {
			value := a.getDefine(arg)
			a.instructions = append(a.instructions, uint8(parseInt(value, 16, instruction)))
		} else if argIsLabel {
			value := a.getLabel(arg)
			a.instructions = append(a.instructions, uint8((value&0xff00)>>8))
			a.instructions = append(a.instructions, uint8((value & 0xff)))
		} else {
			a.instructions = append(a.instructions, uint8(parseInt(arg, 8, instruction)))
		}

	case "SWAP":
		a.instructions = append(a.instructions, instructions.SWAP)

	case "LOAD":
		a.instructions = append(a.instructions, instructions.LOAD)

		if argIsDefine {
			value := a.getDefine(arg)
			n := parseInt(value, 16, instruction)
			a.instructions = append(a.instructions, uint8((n&0xff00)>>8))
			a.instructions = append(a.instructions, uint8(n&0xff))
		} else if argIsLabel {
			value := a.getLabel(arg)
			a.instructions = append(a.instructions, uint8((value&0xff00)>>8))
			a.instructions = append(a.instructions, uint8((value & 0xff)))
		} else {
			n := parseInt(arg, 16, instruction)
			a.instructions = append(a.instructions, uint8((n&0xff00)>>8))
			a.instructions = append(a.instructions, uint8(n&0xff))
		}

	case "POP":
		a.instructions = append(a.instructions, instructions.POP)

	case "PUSH":
		a.instructions = append(a.instructions, instructions.PUSH)

	case "STORE":
		a.instructions = append(a.instructions, instructions.STORE)

		if argIsDefine {
			value := a.getDefine(arg)
			n := parseInt(value, 16, instruction)
			a.instructions = append(a.instructions, uint8((n&0xff00)>>8))
			a.instructions = append(a.instructions, uint8(n&0xff))
		} else if argIsLabel {
			value := a.getLabel(arg)
			a.instructions = append(a.instructions, uint8((value&0xff00)>>8))
			a.instructions = append(a.instructions, uint8((value & 0xff)))
		} else {
			n := parseInt(arg, 16, instruction)
			a.instructions = append(a.instructions, uint8((n&0xff00)>>8))
			a.instructions = append(a.instructions, uint8(n&0xff))
		}

	case "ADD":
		a.instructions = append(a.instructions, instructions.ADD)

	case "SUB":
		a.instructions = append(a.instructions, instructions.SUB)

	case "MUL":
		a.instructions = append(a.instructions, instructions.MUL)

	case "DIV":
		a.instructions = append(a.instructions, instructions.DIV)

	case "SHL":
		a.instructions = append(a.instructions, instructions.SHL)

	case "SHR":
		a.instructions = append(a.instructions, instructions.SHR)

	case "AND":
		a.instructions = append(a.instructions, instructions.AND)

	case "OR":
		a.instructions = append(a.instructions, instructions.OR)

	case "HALT":
		a.instructions = append(a.instructions, instructions.HALT)

	case "JUMP":
		a.instructions = append(a.instructions, instructions.JUMP)

		if argIsDefine {
			value := a.getDefine(arg)

			n := parseInt(value, 16, instruction)
			a.instructions = append(a.instructions, uint8((n&0xff00)>>8))
			a.instructions = append(a.instructions, uint8(n&0xff))
		} else if argIsLabel {
			value := a.getLabel(arg)
			a.instructions = append(a.instructions, uint8((value&0xff00)>>8))
			a.instructions = append(a.instructions, uint8((value & 0xff)))
		} else {
			n := parseInt(arg, 16, instruction)
			a.instructions = append(a.instructions, uint8((n&0xff00)>>8))
			a.instructions = append(a.instructions, uint8(n&0xff))
		}

	case "JUMPZ":
		a.instructions = append(a.instructions, instructions.JUMPZ)

		if argIsDefine {
			value := a.getDefine(arg)

			n := parseInt(value, 16, instruction)
			a.instructions = append(a.instructions, uint8((n&0xff00)>>8))
			a.instructions = append(a.instructions, uint8(n&0xff))
		} else if argIsLabel {
			value := a.getLabel(arg)
			a.instructions = append(a.instructions, uint8((value&0xff00)>>8))
			a.instructions = append(a.instructions, uint8((value & 0xff)))
		} else {
			n := parseInt(arg, 16, instruction)
			a.instructions = append(a.instructions, uint8((n&0xff00)>>8))
			a.instructions = append(a.instructions, uint8(n&0xff))
		}

	case "JUMPNZ":
		a.instructions = append(a.instructions, instructions.JUMPNZ)

		if argIsDefine {
			value := a.getDefine(arg)

			n := parseInt(value, 16, instruction)
			a.instructions = append(a.instructions, uint8((n&0xff00)>>8))
			a.instructions = append(a.instructions, uint8(n&0xff))
		} else if argIsLabel {
			value := a.getLabel(arg)
			a.instructions = append(a.instructions, uint8((value&0xff00)>>8))
			a.instructions = append(a.instructions, uint8((value & 0xff)))
		} else {
			n := parseInt(arg, 16, instruction)
			a.instructions = append(a.instructions, uint8((n&0xff00)>>8))
			a.instructions = append(a.instructions, uint8(n&0xff))
		}

	case "CALL":
		a.instructions = append(a.instructions, instructions.CALL)

		if argIsDefine {
			value := a.getDefine(arg)

			n := parseInt(value, 16, instruction)
			a.instructions = append(a.instructions, uint8((n&0xff00)>>8))
			a.instructions = append(a.instructions, uint8(n&0xff))
		} else if argIsLabel {
			value := a.getLabel(arg)
			a.instructions = append(a.instructions, uint8((value&0xff00)>>8))
			a.instructions = append(a.instructions, uint8((value & 0xff)))
		} else {
			n := parseInt(arg, 16, instruction)
			a.instructions = append(a.instructions, uint8((n&0xff00)>>8))
			a.instructions = append(a.instructions, uint8(n&0xff))
		}

	case "RET":
		a.instructions = append(a.instructions, instructions.RET)

	default:
		log.Fatalf("encountered unknown instruction %s", instruction)
	}

}

var labelRegex = regexp.MustCompile(`^(\w+):`)

func (a *Assembler) getLabels(lines []string) {
	for i, l := range lines {

		l = strings.TrimSpace(l)
		if label := labelRegex.FindString(l); len(label) > 0 {
			fmt.Printf("%d: %s\n", i, label)

			a.labels[label] = uint16(i)
		}
	}

	// fmt.Println(a.labels)
}

func (a *Assembler) GetInstructions(source string) []uint8 {
	a.source = strings.TrimSpace(source)

	lines := strings.Split(source, "\n")

	a.getLabels(lines)

	for _, l := range lines {
		a.ParseLine(l)
	}

	return a.instructions
}
