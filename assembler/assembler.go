package assembler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/andrewesterhuizen/penpal/instructions"
	"github.com/andrewesterhuizen/penpal/lexer"
)

const HeaderSize = 10

type FileGetterFunc func(path string) (string, error)

func FileSystemFileGetterFunc(path string) (string, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(f), nil
}

type Config struct {
	disableHeader  bool
	FileGetterFunc FileGetterFunc
}

type Assembler struct {
	config       Config
	instructions []uint8
	defines      map[string]string
	labels       map[string]uint16
	getFile      FileGetterFunc

	currentLableAddress uint16

	sourceTokens map[string][]lexer.Token
}

func New(config Config) Assembler {
	a := Assembler{config: config}

	if config.FileGetterFunc != nil {
		a.getFile = config.FileGetterFunc
	} else {
		a.getFile = FileSystemFileGetterFunc
	}

	a.defines = make(map[string]string)
	a.labels = make(map[string]uint16)
	a.sourceTokens = make(map[string][]lexer.Token)

	a.currentLableAddress = 0

	if !config.disableHeader {
		a.currentLableAddress = HeaderSize
	}

	return a
}

func (a *Assembler) getDefine(define string) (string, error) {
	value, exists := a.defines[define]
	if !exists {
		return "", fmt.Errorf("no definition found for %s", define)
	}

	return value, nil
}

func (a *Assembler) getLabel(label string) (uint16, error) {
	value, exists := a.labels[label]
	if !exists {
		return 0, fmt.Errorf("no label definition found for %s", label)
	}

	return value, nil
}

func (a *Assembler) processTokens(tokens []lexer.Token) error {
	for _, t := range tokens {
		err := a.processToken(t)
		if err != nil {
			return fmt.Errorf("%s:%v: %w", t.FileName, t.LineNumber, err)
		}
	}

	return nil
}

func (a *Assembler) processFileIncludeToken(t lexer.Token) error {
	filename := t.Value
	tokens := a.sourceTokens[filename]
	return a.processTokens(tokens)
}

func (a *Assembler) processSystemIncludeToken(t lexer.Token) error {
	// TODO: process lexer.TokenSystemInclude here
	return nil
}

func (a *Assembler) processDefineToken(t lexer.Token) error {
	a.defines[t.Value] = t.Args[0].Value
	return nil
}

func (a *Assembler) processInstructionToken(t lexer.Token) error {
	i, exists := instructions.InstructionByName[t.Value]
	if !exists {
		return fmt.Errorf("encountered unknown instruction when getting labels: %s", t.Value)
	}

	a.currentLableAddress += uint16(instructions.Width[i])

	err := a.addInstruction(t)
	if err != nil {
		return err
	}

	return nil
}

func (a *Assembler) processToken(t lexer.Token) error {
	switch t.Type {
	case lexer.TokenFileInclude:
		return a.processFileIncludeToken(t)
	case lexer.TokenSystemInclude:
		return a.processSystemIncludeToken(t)
	case lexer.TokenDefine:
		return a.processDefineToken(t)
	case lexer.TokenInstruction:
		return a.processInstructionToken(t)
	}

	return nil
}

func (a *Assembler) appendInstruction(b uint8) {
	a.instructions = append(a.instructions, b)
}

func (a *Assembler) appendInstructions(bytes ...uint8) {
	for _, b := range bytes {
		a.instructions = append(a.instructions, b)
	}
}

func (a *Assembler) addInstruction(t lexer.Token) error {
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
			return fmt.Errorf("encountered unknown destination for MOV instruction: %s", register.Value)
		}

		arg2 := t.Args[1]

		var addressingMode uint8 = instructions.AddressingModeImmediate

		if arg2.IsFPOffsetAddress {
			addressingMode = instructions.AddressingModeFPRelative
		}

		v, err := a.getInstructionArgs8(arg2, instruction)
		if err != nil {
			return err
		}

		a.appendInstructions(instructions.MOV, addressingMode, dest, v)

	case "SWAP":
		a.appendInstruction(instructions.SWAP)

	case "LOAD":
		a.appendInstruction(instructions.LOAD)

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
				return fmt.Errorf("LOAD: encountered unknown operand type %v", arg0)
			}

			a.addInstructionArgs16(t.Args[1], instruction)
		} else {
			a.appendInstruction(instructions.Register)
			a.appendInstruction(instructions.RegisterA)
			a.addInstructionArgs16(t.Args[0], instruction)
		}

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
				return fmt.Errorf("STORE: encountered unknown operand type %v", arg0)
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
		if len(t.Args) < 1 {
			return fmt.Errorf("expected 1 operand for instruction")
		}
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
		return fmt.Errorf("encountered unknown instruction %s", instruction)
	}

	return nil
}

func (a *Assembler) addInstructionArgs16(arg lexer.Arg, instruction string) error {
	if arg.IsDefine {
		value, err := a.getDefine(arg.Value)
		if err != nil {
			return err
		}

		n, err := strconv.ParseUint(value, 0, 16)
		if err != nil {
			return err
		}

		a.appendInstruction(uint8((n & 0xff00) >> 8))
		a.appendInstruction(uint8(n & 0xff))
	} else if arg.IsLabel {
		value, err := a.getLabel(arg.Value)
		if err != nil {
			return err
		}
		a.appendInstruction(uint8((value & 0xff00) >> 8))
		a.appendInstruction(uint8((value & 0xff)))
	} else {
		n := arg.AsUint()
		a.appendInstruction(uint8((n & 0xff00) >> 8))
		a.appendInstruction(uint8(n & 0xff))
	}

	return nil
}

func (a *Assembler) getInstructionArgs8(arg lexer.Arg, instruction string) (uint8, error) {
	if arg.IsFPOffsetAddress {
		return arg.AsUint8(), nil
	} else if arg.IsDefine {
		value, err := a.getDefine(arg.Value)
		if err != nil {
			return 0, err
		}

		n, err := strconv.ParseUint(value, 0, 16)
		if err != nil {
			return 0, err
		}

		return uint8(n), nil
	} else if arg.IsLabel {
		value, err := a.getLabel(arg.Value)
		if err != nil {
			return 0, err
		}
		return uint8((value & 0xff)), nil
	} else {
		return arg.AsUint8(), nil
	}

}

func (a *Assembler) addInstructionArgs8(arg lexer.Arg, instruction string) error {
	v, err := a.getInstructionArgs8(arg, instruction)
	if err != nil {
		return err
	}

	a.appendInstruction(v)

	return nil
}

func (a *Assembler) getHeaderBytes() ([]byte, error) {
	buf := bytes.Buffer{}

	entryPointAddress, exists := a.labels["__start"]
	// check for entry point
	if !exists {
		return nil, fmt.Errorf("source has no entry point")
	}

	for _, b := range []byte("PENPAL") {
		buf.WriteByte(b)
	}

	// version
	buf.WriteByte(0)
	buf.WriteByte(1)

	// entry point
	buf.WriteByte(byte((entryPointAddress & 0xff00) >> 8))
	buf.WriteByte(byte(entryPointAddress & 0xff))

	return buf.Bytes(), nil
}

func (a *Assembler) getIncludes(filename string, tokens []lexer.Token) error {
	for _, t := range tokens {
		if t.Type == lexer.TokenFileInclude {

			filename := t.Value
			file, err := a.getFile(filename)
			if err != nil {
				return err
			}

			l := lexer.New()
			tokens, err := l.GetTokens(filename, file)
			if err != nil {
				return err
			}

			a.sourceTokens[filename] = tokens

			a.getIncludes(filename, tokens)
		} else if t.Type == lexer.TokenSystemInclude {
			// TODO: handle lexer.TokenSystemInclude here
		}
	}

	return nil
}

func (a *Assembler) getLabels(tokens []lexer.Token) error {
	for _, t := range tokens {
		switch t.Type {
		case lexer.TokenFileInclude:
			filename := t.Value
			tokens := a.sourceTokens[filename]
			a.getLabels(tokens)

		case lexer.TokenSystemInclude:
		// TODO: handle lexer.TokenSystemInclude here

		case lexer.TokenLabel:
			if t.Type == lexer.TokenLabel {
				a.labels[t.Value] = uint16(a.currentLableAddress)
			}

		case lexer.TokenInstruction:
			i, exists := instructions.InstructionByName[t.Value]
			if !exists {
				return fmt.Errorf("encountered unknown instruction %s", t.Value)
			}

			a.currentLableAddress += uint16(instructions.Width[i])
		}
	}

	return nil
}

func (a *Assembler) GetProgram(filename string, source string) ([]uint8, error) {
	l := lexer.New()

	// get tokens for entry point file
	tokens, err := l.GetTokens(filename, source)
	if err != nil {
		return nil, err
	}

	// recursively gets tokens for each included file
	a.getIncludes(filename, tokens)

	// recursively gets labels from all included files
	a.getLabels(tokens)

	out := bytes.Buffer{}

	if !a.config.disableHeader {
		h, err := a.getHeaderBytes()
		if err != nil {
			return nil, err
		}

		out.Write(h)
	}

	a.processTokens(tokens)

	out.Write(a.instructions)

	return out.Bytes(), nil
}
