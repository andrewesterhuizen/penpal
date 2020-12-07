package assembler2

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
	p.skipIf(lexer_rewrite.TokenTypeNewLine)
	return nil
}

func (p *Parser) parseDB() error {
	t, err := p.expect(lexer_rewrite.TokenTypeInteger)
	if err != nil {
		return err
	}

	n, err := parseIntegerToken(t)
	if err != nil {
		return err
	}

	p.addByte(byte(n))

	p.skipIf(lexer_rewrite.TokenTypeNewLine)
	return nil
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

	p.skipIf(lexer_rewrite.TokenTypeNewLine)
	return nil
}

// func (p *Parser) parseArithmeticLogicInstruction(instruction byte) error {
// 	p.addByte(instruction)

// 	t := p.nextToken()

// 	switch t.Type {
// 	// no operand = implied B register
// 	case lexer_rewrite.TokenTypeNewLine:
// 		p.addByte(instructions.Register)
// 		p.addByte(instructions.RegisterB)

// 		return nil
// 	case lexer_rewrite.TokenTypeInteger:
// 		p.addByte(instructions.Immediate)

// 		n, err := parseIntegerToken(t)
// 		if err != nil {
// 			return err
// 		}

// 		p.addByte(byte(n))

// 	default:
// 		return fmt.Errorf("unexpected operand \"%s\"", t.Value)
// 	}

// 	_, err := p.expect(lexer_rewrite.TokenTypeNewLine)
// 	return err
// }

func (p *Parser) parsePushInstruction() error {
	p.addByte(instructions.Push)

	t := p.nextToken()

	switch t.Type {
	// no operand = implied A register
	case lexer_rewrite.TokenTypeEndOfFile:
		p.addByte(instructions.Register)
		p.addByte(instructions.RegisterA)

		return nil
	case lexer_rewrite.TokenTypeNewLine:
		p.addByte(instructions.Register)
		p.addByte(instructions.RegisterA)

		return nil
	case lexer_rewrite.TokenTypeInteger:
		p.addByte(instructions.Immediate)

		n, err := parseIntegerToken(t)
		if err != nil {
			return err
		}

		p.addByte(byte(n))

	default:
		fmt.Println(t)
		return fmt.Errorf("unexpected operand \"%s\"", t.Value)
	}

	p.skipIf(lexer_rewrite.TokenTypeNewLine)
	return nil
}

func (p *Parser) parseIndex() (bool, byte, error) {
	_, err := p.expect(lexer_rewrite.TokenTypeLeftBracket)
	if err != nil {
		return false, 0, err
	}

	expectedTokenTypes := []lexer_rewrite.TokenType{lexer_rewrite.TokenTypeInteger, lexer_rewrite.TokenTypeText}
	t, err := p.expectRange(expectedTokenTypes)
	if err != nil {
		return false, 0, err
	}

	_, err = p.expect(lexer_rewrite.TokenTypeRightBracket)
	if err != nil {
		return false, 0, err
	}

	switch t.Type {
	case lexer_rewrite.TokenTypeInteger:
		n, err := parseIntegerToken(t)
		if err != nil {
			return false, 0, err
		}

		return false, byte(n), nil

	case lexer_rewrite.TokenTypeText:
		r, err := getRegister(t.Value)
		if err != nil {
			return false, 0, err
		}

		return true, r, nil

	default:
		return false, 0, fmt.Errorf("unxpected tokend %s", t.Value)
	}
}

// this function needs some tests and cleaning up
func (p *Parser) parseOffsetAddress() (byte, byte, uint16, error) {
	_, err := p.expect(lexer_rewrite.TokenTypeLeftParen)
	if err != nil {
		return 0, 0, 0, err
	}

	t, err := p.expect(lexer_rewrite.TokenTypeText)
	if err != nil {
		return 0, 0, 0, err
	}

	mode := byte(instructions.Immediate)
	modeArg := byte(0)

	isLabel := false
	var labelAddress uint16

	if t.Value == "fp" {
		mode = instructions.FramePointerWithOffset
	} else {
		addr, exists := p.labels[t.Value]
		if !exists {
			return 0, 0, 0, fmt.Errorf("no definitions found for label %s", t.Value)
		}

		isLabel = true
		labelAddress = addr
	}

	next := p.nextToken()

	// get the offset
	switch next.Type {
	case lexer_rewrite.TokenTypePlus:
		t := p.nextToken()

		switch t.Type {
		case lexer_rewrite.TokenTypeInteger:
			n, err := parseIntegerToken(t)
			if err != nil {
				return 0, 0, 0, err
			}

			modeArg = byte(n)
		case lexer_rewrite.TokenTypeText:
			reg, err := getRegister(t.Value)
			if err != nil {
				return 0, 0, 0, err
			}

			if isLabel {
				mode = instructions.ImmediatePlusRegister
			} else {
				mode = instructions.FramePointerPlusRegister
			}

			modeArg = byte(reg)

		default:
			return 0, 0, 0, fmt.Errorf("unexpected token \"%s\"", t.Value)
		}

	case lexer_rewrite.TokenTypeMinus:
		t := p.nextToken()

		switch t.Type {
		case lexer_rewrite.TokenTypeInteger:
			n, err := parseIntegerToken(t)
			if err != nil {
				return 0, 0, 0, err
			}

			nsigned := int8(n) * int8(-1)
			modeArg = byte(nsigned)

		case lexer_rewrite.TokenTypeText:
			reg, err := getRegister(t.Value)
			if err != nil {
				return 0, 0, 0, err
			}

			if isLabel {
				mode = instructions.ImmediateMinusRegister
			} else {
				mode = instructions.FramePointerMinusRegister
			}

			modeArg = byte(reg)

		default:
			return 0, 0, 0, fmt.Errorf("unexpected token \"%s\"", t.Value)
		}

	// handle offset from [n] or [reg]
	case lexer_rewrite.TokenTypeLeftBracket:
		p.backup()

		isRegister, n, err := p.parseIndex()
		if err != nil {
			return 0, 0, 0, err
		}

		if isRegister {
			mode = instructions.ImmediatePlusRegister
		}

		modeArg = n
	default:
		return 0, 0, 0, fmt.Errorf("unexpected token \"%s\"", next.Value)
	}

	p.expect(lexer_rewrite.TokenTypeRightParen)

	return mode, modeArg, labelAddress, nil
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
		mode, offset, addr, err := p.parseOffsetAddress()

		if err != nil {
			return 0, 0, 0, 0, err
		}

		h := byte((addr & 0xff00) >> 8)
		l := byte(addr & 0xff)

		return mode, offset, h, l, nil

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

func (p *Parser) parseMov() error {
	p.addByte(instructions.Mov)

	dest := p.nextToken()

	if dest.Type != lexer_rewrite.TokenTypeText {
		return fmt.Errorf("unexpected token %s", dest.Type)
	}

	reg, err := getRegister(dest.Value)
	if err != nil {
		return err
	}

	p.addByte(reg)

	_, err = p.expect(lexer_rewrite.TokenTypeComma)
	if err != nil {
		return err
	}

	t, err := p.expect(lexer_rewrite.TokenTypeInteger)
	if err != nil {
		return err
	}

	n, err := parseIntegerToken(t)
	if err != nil {
		return err
	}

	p.addByte(byte(n))

	p.skipIf(lexer_rewrite.TokenTypeNewLine)
	return nil
}

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

	p.skipIf(lexer_rewrite.TokenTypeNewLine)
	return nil
}

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

	p.skipIf(lexer_rewrite.TokenTypeNewLine)
	return nil
}
