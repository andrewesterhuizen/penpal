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
	SystemIncludes map[string]string
}

type Assembler struct {
	config       Config
	instructions []uint8
	defines      map[string]string
	labels       map[string]uint16
	getFile      FileGetterFunc

	currentLableAddress uint16

	systemIncludeSources map[string]string

	includeTokens       map[string][]lexer.Token
	systemIncludeTokens map[string][]lexer.Token
}

func New(config Config) Assembler {
	a := Assembler{config: config}

	if config.FileGetterFunc != nil {
		a.getFile = config.FileGetterFunc
	} else {
		a.getFile = FileSystemFileGetterFunc
	}

	if config.SystemIncludes != nil {
		a.systemIncludeSources = config.SystemIncludes
	} else {
		a.systemIncludeSources = make(map[string]string)
	}

	a.defines = make(map[string]string)
	a.labels = make(map[string]uint16)
	a.includeTokens = make(map[string][]lexer.Token)
	a.systemIncludeTokens = make(map[string][]lexer.Token)

	a.currentLableAddress = 0

	if !config.disableHeader {
		a.currentLableAddress = HeaderSize
	}

	return a
}

func (a *Assembler) getDefine(define string) (string, bool) {
	value, exists := a.defines[define]
	return value, exists
}

func (a *Assembler) getLabel(label string) (uint16, bool) {
	value, exists := a.labels[label]
	return value, exists
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
	name := t.Value
	tokens := a.includeTokens[name]
	return a.processTokens(tokens)
}

func (a *Assembler) processSystemIncludeToken(t lexer.Token) error {
	name := t.Value
	tokens := a.systemIncludeTokens[name]
	return a.processTokens(tokens)
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
	if arg.IsIdentifier {
		// first try get define
		value, exists := a.getDefine(arg.Value)
		if exists {
			n, err := strconv.ParseUint(value, 0, 16)
			if err != nil {
				return err
			}

			a.appendInstruction(uint8((n & 0xff00) >> 8))
			a.appendInstruction(uint8(n & 0xff))
		} else {
			// try get label
			value, exists := a.getLabel(arg.Value)
			if !exists {
				return fmt.Errorf("no definition found for identifier %s", arg.Value)
			}
			a.appendInstruction(uint8((value & 0xff00) >> 8))
			a.appendInstruction(uint8((value & 0xff)))
		}

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
	} else if arg.IsIdentifier {
		defineValue, exists := a.getDefine(arg.Value)
		if exists {
			n, err := strconv.ParseUint(defineValue, 0, 16)
			if err != nil {
				return 0, err
			}
			return uint8(n), nil
		}

		labelValue, exists := a.getLabel(arg.Value)
		if !exists {
			return 0, fmt.Errorf("no definition found for identifier %s", arg.Value)
		}

		return uint8((labelValue & 0xff)), nil

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

func (a *Assembler) getFileInclude(name string) error {
	source, err := a.getFile(name)
	if err != nil {
		return err
	}

	l := lexer.New()
	tokens, err := l.GetTokens(name, source)
	if err != nil {
		return err
	}

	a.includeTokens[name] = tokens

	err = a.getIncludes(name, tokens)
	if err != nil {
		return err
	}

	return nil
}

func (a *Assembler) getSystemInclude(name string) error {
	source, exists := a.systemIncludeSources[name]

	if !exists {
		return fmt.Errorf("no source available for system include <%s>", name)
	}

	l := lexer.New()
	tokens, err := l.GetTokens(name, source)
	if err != nil {
		return err
	}

	a.systemIncludeTokens[name] = tokens

	err = a.getIncludes(name, tokens)
	if err != nil {
		return err
	}

	return nil
}

func (a *Assembler) getIncludes(filename string, tokens []lexer.Token) error {
	for _, t := range tokens {
		switch t.Type {
		case lexer.TokenFileInclude:
			err := a.getFileInclude(t.Value)
			if err != nil {
				return err
			}

		case lexer.TokenSystemInclude:
			err := a.getSystemInclude(t.Value)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (a *Assembler) getLabels(tokens []lexer.Token) error {
	for _, t := range tokens {
		switch t.Type {
		case lexer.TokenFileInclude:
			filename := t.Value
			tokens := a.includeTokens[filename]
			a.getLabels(tokens)

		case lexer.TokenSystemInclude:
			name := t.Value
			tokens := a.systemIncludeTokens[name]
			a.getLabels(tokens)

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
	err = a.getIncludes(filename, tokens)
	if err != nil {
		return nil, err
	}

	// recursively gets labels from all included files
	err = a.getLabels(tokens)
	if err != nil {
		return nil, err
	}

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
