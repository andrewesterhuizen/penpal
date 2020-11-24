package assembler

import (
	"log"
	"strconv"
	"strings"

	"github.com/andrewesterhuizen/penpal/instructions"
	"github.com/andrewesterhuizen/penpal/lexer"
)

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

func parseInt(s string, base int, instruction string) uint64 {
	n, err := strconv.ParseUint(s, 0, base)
	if err != nil {
		log.Fatalf("failed at instruction %s: unable to parse int from value %s", instruction, s)
	}

	return n
}

func (a *Assembler) ParseToken(t lexer.Token) {
	switch t.Type {
	case lexer.TokenDefine:
		a.defines[t.Value] = t.Args[0].Value
	case lexer.TokenLabel:
	case lexer.TokenInstruction:
		a.addInstruction(t)

	default:
		log.Fatalf("encountered unexpected token type %s\n", t.Type)
	}
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

func (a *Assembler) addInstruction(t lexer.Token) {
	instruction := t.Value

	switch instruction {
	case "MOV":

		a.instructions = append(a.instructions, instructions.MOV)

		dest := t.Args[0]

		var destination uint8

		switch dest.Value {
		case "A":
			destination = instructions.DestRegisterA
		case "B":
			destination = instructions.DestRegisterB
		default:
			log.Fatalf("encountered unknown destination for MOV instruction: %s", dest.Value)
		}

		arg2 := t.Args[1]

		var addressingMode uint8 = instructions.AddressingModeImmediate

		if arg2.IsFPOffsetAddress {
			addressingMode = instructions.AddressingModeFPRelative
		}

		a.instructions = append(a.instructions, instructions.EncodeFlags(addressingMode, destination))
		a.addInstructionArgs8(arg2, instruction)

	case "SWAP":
		a.instructions = append(a.instructions, instructions.SWAP)

	case "LOAD":
		a.instructions = append(a.instructions, instructions.LOAD)
		a.addInstructionArgs16(t.Args[0], instruction)

	case "POP":
		a.instructions = append(a.instructions, instructions.POP)

	case "PUSH":
		a.instructions = append(a.instructions, instructions.PUSH)

	case "STORE":
		a.instructions = append(a.instructions, instructions.STORE)
		a.addInstructionArgs16(t.Args[0], instruction)

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

	case "SEND":
		a.instructions = append(a.instructions, instructions.SEND)

	case "RAND":
		a.instructions = append(a.instructions, instructions.RAND)

	case "JUMP":
		a.instructions = append(a.instructions, instructions.JUMP)
		a.addInstructionArgs16(t.Args[0], instruction)

	case "JUMPZ":
		a.instructions = append(a.instructions, instructions.JUMPZ)
		a.addInstructionArgs16(t.Args[0], instruction)

	case "JUMPNZ":
		a.instructions = append(a.instructions, instructions.JUMPNZ)
		a.addInstructionArgs16(t.Args[0], instruction)

	case "CALL":
		a.instructions = append(a.instructions, instructions.CALL)
		a.addInstructionArgs16(t.Args[0], instruction)

	case "RET":
		a.instructions = append(a.instructions, instructions.RET)

	default:
		log.Fatalf("encountered unknown instruction %s", instruction)
	}

}

func (a *Assembler) addInstructionArgs16(arg lexer.Arg, instruction string) {
	if arg.IsDefine {
		value := a.getDefine(arg.Value)
		n := parseInt(value, 16, instruction)
		a.instructions = append(a.instructions, uint8((n&0xff00)>>8))
		a.instructions = append(a.instructions, uint8(n&0xff))
	} else if arg.IsLabel {
		value := a.getLabel(arg.Value)
		a.instructions = append(a.instructions, uint8((value&0xff00)>>8))
		a.instructions = append(a.instructions, uint8((value & 0xff)))
	} else {
		n := arg.AsUint()
		a.instructions = append(a.instructions, uint8((n&0xff00)>>8))
		a.instructions = append(a.instructions, uint8(n&0xff))
	}
}

func (a *Assembler) addInstructionArgs8(arg lexer.Arg, instruction string) {
	if arg.IsFPOffsetAddress {
		a.instructions = append(a.instructions, arg.AsUint8())
	} else if arg.IsDefine {
		value := a.getDefine(arg.Value)
		a.instructions = append(a.instructions, uint8(parseInt(value, 16, instruction)))
	} else if arg.IsLabel {
		value := a.getLabel(arg.Value)
		// a.instructions = append(a.instructions, uint8((value&0xff00)>>8))
		a.instructions = append(a.instructions, uint8((value & 0xff)))
	} else {
		a.instructions = append(a.instructions, arg.AsUint8())
	}
}

func (a *Assembler) getLabels(tokens []lexer.Token) {
	instructionNumber := 0

	for _, t := range tokens {
		switch t.Type {
		case lexer.TokenDefine:
			// skip
		case lexer.TokenLabel:
			// add label
			if t.Type == lexer.TokenLabel {
				a.labels[t.Value] = uint16(instructionNumber)
			}

		case lexer.TokenInstruction:
			i, exists := instructions.InstructionByName[t.Value]
			if !exists {
				log.Fatalf("encountered unknown instruction when getting labels: %s", t.Value)
			}

			instructionNumber += instructions.Width[i]

		default:
			log.Fatalf("encountered unexpected token type when getting labels: %s\n", t.Type)
		}

	}
}

func (a *Assembler) GetInstructions(source string) []uint8 {
	a.source = strings.TrimSpace(source)

	l := lexer.New()
	tokens, err := l.GetTokens(source)
	if err != nil {
		log.Fatalf("assembler failed: %v\n", err)
	}

	a.getLabels(tokens)

	for _, t := range tokens {
		a.ParseToken(t)
	}

	return a.instructions
}
