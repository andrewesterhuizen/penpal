package assembler

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/andrewesterhuizen/penpal/instructions"
	"github.com/andrewesterhuizen/penpal/lexer"
)

const HeaderSize = 10

type Config struct {
	disableHeader bool
}

type Assembler struct {
	config       Config
	index        int
	source       string
	instructions []uint8
	defines      map[string]string
	labels       map[string]uint16
}

func New(config Config) Assembler {
	a := Assembler{config: config}
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

func (a *Assembler) appendInstruction(b uint8) {
	a.instructions = append(a.instructions, b)
}

func (a *Assembler) appendInstructions(bytes ...uint8) {
	for _, b := range bytes {
		a.instructions = append(a.instructions, b)
	}
}

func (a *Assembler) addInstruction(t lexer.Token) {
	instruction := t.Value

	switch instruction {
	case "MOV":
		register := t.Args[0]
		var dest uint8

		switch register.Value {
		case "A":
			dest = instructions.RegisterA
		case "B":
			dest = instructions.RegisterB
		default:
			log.Fatalf("encountered unknown destination for MOV instruction: %s", register.Value)
		}

		arg2 := t.Args[1]

		var addressingMode uint8 = instructions.AddressingModeImmediate

		if arg2.IsFPOffsetAddress {
			addressingMode = instructions.AddressingModeFPRelative
		}

		a.appendInstructions(instructions.MOV, addressingMode, dest, a.getInstructionArgs8(arg2, instruction))

	case "SWAP":
		a.appendInstruction(instructions.SWAP)

	case "LOAD":
		a.appendInstruction(instructions.LOAD)
		a.addInstructionArgs16(t.Args[0], instruction)

	case "POP":
		a.appendInstruction(instructions.POP)

	case "PUSH":
		if len(t.Args) == 1 {
			a.appendInstruction(instructions.PUSH)
			arg := t.Args[0]

			if arg.IsRegister {
				a.appendInstruction(instructions.Register)
				a.appendInstruction(instructions.RegistersByName[arg.Value])
			} else if arg.IsFPOffsetAddress {
				a.appendInstruction(instructions.FramePointerRelativeAddress)
				a.appendInstruction(arg.AsUint8())
			} else {
				a.appendInstruction(instructions.Value)
				a.appendInstruction(arg.AsUint8())
			}

		} else {
			a.appendInstruction(instructions.PUSH)
			a.appendInstruction(instructions.Register)
			a.appendInstruction(instructions.RegisterA)
		}

	case "STORE":
		a.appendInstruction(instructions.STORE)

		if len(t.Args) == 2 {
			arg0 := t.Args[0]
			if arg0.IsFPOffsetAddress {
				a.appendInstruction(instructions.FramePointerRelativeAddress)
				a.appendInstruction(t.Args[0].AsUint8())
			} else if arg0.IsRegister {
				a.appendInstruction(instructions.Register)
				register := instructions.RegistersByName[t.Args[0].Value]
				a.appendInstruction(register)

			} else {
				log.Fatalf("STORE: encountered unknown operand type %v", arg0)
			}

			a.addInstructionArgs16(t.Args[1], instruction)
		} else {
			a.appendInstruction(instructions.Register)
			a.appendInstruction(instructions.RegisterA)
			a.addInstructionArgs16(t.Args[0], instruction)
		}

	case "ADD":
		a.appendInstruction(instructions.ADD)

	case "SUB":
		a.appendInstruction(instructions.SUB)

	case "MUL":
		a.appendInstruction(instructions.MUL)

	case "DIV":
		a.appendInstruction(instructions.DIV)

	case "SHL":
		a.appendInstruction(instructions.SHL)

	case "SHR":
		a.appendInstruction(instructions.SHR)

	case "AND":
		a.appendInstruction(instructions.AND)

	case "OR":
		a.appendInstruction(instructions.OR)

	case "HALT":
		a.appendInstruction(instructions.HALT)

	case "SEND":
		a.appendInstruction(instructions.SEND)

	case "RAND":
		a.appendInstruction(instructions.RAND)

	case "JUMP":
		a.appendInstruction(instructions.JUMP)
		a.addInstructionArgs16(t.Args[0], instruction)

	case "JUMPZ":
		a.appendInstruction(instructions.JUMPZ)
		a.addInstructionArgs16(t.Args[0], instruction)

	case "JUMPNZ":
		a.appendInstruction(instructions.JUMPNZ)
		a.addInstructionArgs16(t.Args[0], instruction)

	case "CALL":
		a.appendInstruction(instructions.CALL)
		a.addInstructionArgs16(t.Args[0], instruction)

	case "RET":
		a.appendInstruction(instructions.RET)

	default:
		log.Fatalf("encountered unknown instruction %s", instruction)
	}

}

func (a *Assembler) addInstructionArgs16(arg lexer.Arg, instruction string) {
	if arg.IsDefine {
		value := a.getDefine(arg.Value)
		n := parseInt(value, 16, instruction)
		a.appendInstruction(uint8((n & 0xff00) >> 8))
		a.appendInstruction(uint8(n & 0xff))
	} else if arg.IsLabel {
		value := a.getLabel(arg.Value)
		a.appendInstruction(uint8((value & 0xff00) >> 8))
		a.appendInstruction(uint8((value & 0xff)))
	} else {
		n := arg.AsUint()
		a.appendInstruction(uint8((n & 0xff00) >> 8))
		a.appendInstruction(uint8(n & 0xff))
	}
}

func (a *Assembler) getInstructionArgs8(arg lexer.Arg, instruction string) uint8 {
	if arg.IsFPOffsetAddress {
		return arg.AsUint8()
	} else if arg.IsDefine {
		value := a.getDefine(arg.Value)
		return uint8(parseInt(value, 16, instruction))
	} else if arg.IsLabel {
		value := a.getLabel(arg.Value)
		return uint8((value & 0xff))
	} else {
		return arg.AsUint8()
	}
}

func (a *Assembler) addInstructionArgs8(arg lexer.Arg, instruction string) {
	a.appendInstruction(a.getInstructionArgs8(arg, instruction))
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

func (a *Assembler) addHeader() error {
	entryPointAddress, exists := a.labels["__start"]

	// check for entry point
	if !exists {
		return fmt.Errorf("source has no entry point")
	}

	for _, b := range []byte("PENPAL") {
		a.appendInstruction(b)
	}

	// version
	a.appendInstruction(0)
	a.appendInstruction(1)

	// entry point
	entryPointAddress += HeaderSize
	a.appendInstruction(byte((entryPointAddress & 0xff00) >> 8))
	a.appendInstruction(byte(entryPointAddress & 0xff))

	return nil
}

func (a *Assembler) GetInstructions(source string) ([]uint8, error) {
	a.source = strings.TrimSpace(source)

	l := lexer.New()
	tokens, err := l.GetTokens(source)
	if err != nil {
		log.Fatalf("assembler failed: %v\n", err)
	}

	a.getLabels(tokens)

	if !a.config.disableHeader {
		err = a.addHeader()
		if err != nil {
			return nil, err
		}

		// update labels to account for header size
		for l := range a.labels {
			a.labels[l] += uint16(HeaderSize)
		}
	}

	for _, t := range tokens {
		a.ParseToken(t)
	}

	return a.instructions, nil
}
