package parser

import (
	"fmt"

	"github.com/andrewesterhuizen/penpal/instructions"
	"github.com/andrewesterhuizen/penpal/lexer_rewrite"
)

func getRegister(r string) (byte, error) {
	switch r {
	case "A":
		return instructions.RegisterA, nil
	case "B":
		return instructions.RegisterB, nil
	default:
		return 0, fmt.Errorf("expected register and got %s", r)
	}
}

func (p *Parser) parseNoOperandInstruction(instruction byte) error {
	p.addByte(instruction)
	_, err := p.expect(lexer_rewrite.TokenTypeNewLine)
	return err
}

func (p *Parser) parseAddressInstruction(instruction byte) error {
	p.addByte(instruction)

	t := p.nextToken()

	var addr uint16

	switch t.Type {
	case lexer_rewrite.TokenTypeInteger:
		i, err := parseIntegerToken(t)
		if err != nil {
			return err
		}

		addr = uint16(i)

	case lexer_rewrite.TokenTypeText:
		i, exists := p.labels[t.Value]
		if !exists {
			return fmt.Errorf("no definitions found for label %s", t.Value)
		}

		addr = i

	default:
		return fmt.Errorf("unexpected token %s", t.Value)
	}

	h := (addr & 0xff00) >> 8
	l := addr & 0xff

	p.addByte(byte(h))
	p.addByte(byte(l))

	_, err := p.expect(lexer_rewrite.TokenTypeNewLine)
	return err
}

// encoding = instruction, (immediate mode|register mode), (immediate value|register number)
func (p *Parser) parseArithmeticLogicInstruction(instruction byte) error {
	p.addByte(instruction)

	t := p.nextToken()

	switch t.Type {
	// no operand = implied B register
	case lexer_rewrite.TokenTypeNewLine:
		p.addByte(instructions.Register)
		p.addByte(instructions.RegisterB)

		return nil
	case lexer_rewrite.TokenTypeInteger:
		p.addByte(instructions.Immediate)

		n, err := parseIntegerToken(t)
		if err != nil {
			return err
		}

		p.addByte(byte(n))

	default:
		return fmt.Errorf("unexpected operand \"%s\"", t.Value)
	}

	_, err := p.expect(lexer_rewrite.TokenTypeNewLine)
	return err
}

func (p *Parser) parseOffsetAddress() (byte, byte, error) {
	_, err := p.expect(lexer_rewrite.TokenTypeLeftParen)
	if err != nil {
		return 0, 0, err
	}

	register, err := p.expect(lexer_rewrite.TokenTypeText)
	if err != nil {
		return 0, 0, err
	}

	mode := byte(instructions.Immediate)
	modeArg := byte(0)

	switch register.Value {
	case "fp":
		mode = instructions.FramePointerWithOffset

	default:
		return 0, 0, fmt.Errorf("unexpected register \"%s\"", register.Value)
	}

	next := p.nextToken()

	switch next.Type {
	case lexer_rewrite.TokenTypePlus:
		t := p.nextToken()

		switch t.Type {
		case lexer_rewrite.TokenTypeInteger:
			n, err := parseIntegerToken(t)
			if err != nil {
				return 0, 0, err
			}

			modeArg = byte(n)
		case lexer_rewrite.TokenTypeText:
			reg, err := getRegister(t.Value)
			if err != nil {
				return 0, 0, err
			}

			mode = instructions.FramePointerPlusRegister
			modeArg = byte(reg)

		default:
			return 0, 0, fmt.Errorf("unexpected token \"%s\"", t.Value)
		}

	case lexer_rewrite.TokenTypeMinus:
		t := p.nextToken()

		switch t.Type {
		case lexer_rewrite.TokenTypeInteger:
			n, err := parseIntegerToken(t)
			if err != nil {
				return 0, 0, err
			}

			nsigned := int8(n) * int8(-1)
			modeArg = byte(nsigned)

		case lexer_rewrite.TokenTypeText:
			reg, err := getRegister(t.Value)
			if err != nil {
				return 0, 0, err
			}

			mode = instructions.FramePointerMinusRegister
			modeArg = byte(reg)

		default:
			return 0, 0, fmt.Errorf("unexpected token \"%s\"", t.Value)
		}

	case lexer_rewrite.TokenTypeLeftBracket:
		t := p.nextToken()

		// default to +
		operator := lexer_rewrite.TokenTypePlus

		// valid next tokens are, +, - or the integer

		// check if token is - for inverting later
		if t.Type == lexer_rewrite.TokenTypeMinus {
			operator = lexer_rewrite.TokenTypeMinus
		}

		// if current token isn't the integer then it must be the next token
		if t.Type != lexer_rewrite.TokenTypeInteger {
			t, err = p.expect(lexer_rewrite.TokenTypeInteger)
			if err != nil {
				return 0, 0, err
			}
		}

		_, err = p.expect(lexer_rewrite.TokenTypeRightBracket)
		if err != nil {
			return 0, 0, err
		}

		n, err := parseIntegerToken(t)
		if err != nil {
			return 0, 0, err
		}

		if operator == lexer_rewrite.TokenTypeMinus {
			nsigned := int8(n) * int8(-1)
			modeArg = byte(nsigned)

		} else {
			modeArg = byte(n)
		}

	default:
		return 0, 0, fmt.Errorf("unexpected token \"%s\"", next.Value)
	}

	p.expect(lexer_rewrite.TokenTypeRightParen)

	return mode, modeArg, nil
}

func (p *Parser) parseMemoryAddress() (byte, byte, byte, byte, error) {
	t := p.nextToken()

	switch t.Type {
	case lexer_rewrite.TokenTypeInteger:
		n, err := parseIntegerToken(t)
		if err != nil {
			return 0, 0, 0, 0, err
		}

		h := byte((n & 0xff00) >> 8)
		l := byte(n & 0xff)

		return instructions.Immediate, 0, h, l, nil

	case lexer_rewrite.TokenTypeLeftParen:
		p.backup()
		mode, offset, err := p.parseOffsetAddress()
		if err != nil {
			return 0, 0, 0, 0, err
		}

		return mode, offset, 0, 0, nil

	case lexer_rewrite.TokenTypeText:
		switch t.Value {
		case "fp":
			return instructions.FramePointerWithOffset, 0, 0, 0, nil

		default:
			return 0, 0, 0, 0, fmt.Errorf("unknown register \"%s\"", t.Value)
		}
	default:
		return 0, 0, 0, 0, fmt.Errorf("unexpected token \"%s\"", t.Value)
	}
}

// encoding =  mov, mode, value, dest
func (p *Parser) parseMov() error {
	p.addByte(instructions.Mov)

	t := p.nextToken()

	switch t.Type {
	case lexer_rewrite.TokenTypeInteger:
		n, err := parseIntegerToken(t)
		if err != nil {
			return err
		}

		p.addByte(instructions.Immediate)
		p.addByte(byte(n))

	case lexer_rewrite.TokenTypeLeftParen:
		p.backup()
		mode, offset, err := p.parseOffsetAddress()
		if err != nil {
			return err
		}

		p.addByte(mode)
		p.addByte(offset)

	case lexer_rewrite.TokenTypeText:
		switch t.Value {
		case "fp":
			// mov fp, A => (A = (fp+0))
			p.addByte(instructions.FramePointerWithOffset)
			p.addByte(0)

		default:
			return fmt.Errorf("unknown register \"%s\"", t.Value)
		}

	default:
		return fmt.Errorf("unexpected token \"%s\"", t.Value)
	}

	_, err := p.expect(lexer_rewrite.TokenTypeComma)
	if err != nil {
		return err
	}

	dest := p.nextToken()

	if dest.Type != lexer_rewrite.TokenTypeText {
		return fmt.Errorf("unexpected token %s", dest.Type)
	}

	reg, err := getRegister(dest.Value)
	if err != nil {
		return err
	}

	p.addByte(reg)

	_, err = p.expect(lexer_rewrite.TokenTypeNewLine)
	return err
}

// instructions.Load, addressh, addressl, mode, offset, dest reg
func (p *Parser) parseLoad() error {
	p.addByte(instructions.Load)

	// address bytes
	mode, modeArg, h, l, err := p.parseMemoryAddress()
	if err != nil {
		return err
	}

	p.addByte(h)
	p.addByte(l)
	p.addByte(mode)
	p.addByte(modeArg)

	// ,
	_, err = p.expect(lexer_rewrite.TokenTypeComma)
	if err != nil {
		return err
	}

	// register
	dest := p.nextToken()
	if dest.Type != lexer_rewrite.TokenTypeText {
		return fmt.Errorf("expected register, got %s", dest.Type)
	}

	reg, err := getRegister(dest.Value)
	if err != nil {
		return err
	}

	p.addByte(reg)

	_, err = p.expect(lexer_rewrite.TokenTypeNewLine)
	return err
}

// instructions.Store, dest reg, mode, offset, addressh, addressl
func (p *Parser) parseStore() error {
	p.addByte(instructions.Store)

	// register
	dest := p.nextToken()
	if dest.Type != lexer_rewrite.TokenTypeText {
		return fmt.Errorf("expected register, got %s", dest.Type)
	}

	reg, err := getRegister(dest.Value)
	if err != nil {
		return err
	}

	p.addByte(reg)

	// ,
	_, err = p.expect(lexer_rewrite.TokenTypeComma)
	if err != nil {
		return err
	}

	// address bytes
	mode, modeArg, h, l, err := p.parseMemoryAddress()
	if err != nil {
		return err
	}

	p.addByte(mode)
	p.addByte(modeArg)
	p.addByte(h)
	p.addByte(l)

	_, err = p.expect(lexer_rewrite.TokenTypeNewLine)
	return err
}
