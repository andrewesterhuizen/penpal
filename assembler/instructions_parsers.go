package assembler

import (
	"fmt"

	"github.com/andrewesterhuizen/penpal/instructions"
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

func (p *parser) parseNoOperandInstruction(instruction byte) error {
	p.addByte(instruction)
	p.skipIf(tokenTypeNewLine)
	return nil
}

func (p *parser) parseDB() error {
	t, err := p.expect(tokenTypeInteger)
	if err != nil {
		return err
	}

	n, err := parseIntegerToken(t)
	if err != nil {
		return err
	}

	p.addByte(byte(n))

	p.skipIf(tokenTypeNewLine)
	return nil
}

func (p *parser) parseAddressInstruction(instruction byte) error {
	p.addByte(instruction)

	t := p.nextToken()

	var addr uint16

	switch t.tokenType {
	case tokenTypeInteger:
		i, err := parseIntegerToken(t)
		if err != nil {
			return err
		}

		addr = uint16(i)

	case tokenTypeText:
		i, exists := p.labels[t.value]
		if !exists {
			return fmt.Errorf("no definitions found for label %s", t.value)
		}

		addr = i

	default:
		return fmt.Errorf("unexpected token %s", t.value)
	}

	h := (addr & 0xff00) >> 8
	l := addr & 0xff

	p.addByte(byte(h))
	p.addByte(byte(l))

	p.skipIf(tokenTypeNewLine)
	return nil
}

// func (p *parser) parseArithmeticLogicInstruction(instruction byte) error {
// 	p.addByte(instruction)

// 	t := p.nextToken()

// 	switch t.tokenType {
// 	// no operand = implied B register
// 	case tokenTypeNewLine:
// 		p.addByte(instructions.Register)
// 		p.addByte(instructions.RegisterB)

// 		return nil
// 	case tokenTypeInteger:
// 		p.addByte(instructions.Immediate)

// 		n, err := parseIntegerToken(t)
// 		if err != nil {
// 			return err
// 		}

// 		p.addByte(byte(n))

// 	default:
// 		return fmt.Errorf("unexpected operand \"%s\"", t.value)
// 	}

// 	_, err := p.expect(tokenTypeNewLine)
// 	return err
// }

func (p *parser) parsePushInstruction() error {
	p.addByte(instructions.Push)

	t := p.nextToken()

	switch t.tokenType {
	// no operand = implied A register
	case tokenTypeEndOfFile:
		p.addByte(instructions.Register)
		p.addByte(instructions.RegisterA)

		return nil
	case tokenTypeNewLine:
		p.addByte(instructions.Register)
		p.addByte(instructions.RegisterA)

		return nil
	case tokenTypeInteger:
		p.addByte(instructions.Immediate)

		n, err := parseIntegerToken(t)
		if err != nil {
			return err
		}

		p.addByte(byte(n))

	default:
		fmt.Println(t)
		return fmt.Errorf("unexpected operand \"%s\"", t.value)
	}

	p.skipIf(tokenTypeNewLine)
	return nil
}

func (p *parser) parseIndex() (bool, byte, error) {
	_, err := p.expect(tokenTypeLeftBracket)
	if err != nil {
		return false, 0, err
	}

	expectedtokenTypes := []tokenType{tokenTypeInteger, tokenTypeText}
	t, err := p.expectRange(expectedtokenTypes)
	if err != nil {
		return false, 0, err
	}

	_, err = p.expect(tokenTypeRightBracket)
	if err != nil {
		return false, 0, err
	}

	switch t.tokenType {
	case tokenTypeInteger:
		n, err := parseIntegerToken(t)
		if err != nil {
			return false, 0, err
		}

		return false, byte(n), nil

	case tokenTypeText:
		r, err := getRegister(t.value)
		if err != nil {
			return false, 0, err
		}

		return true, r, nil

	default:
		return false, 0, fmt.Errorf("unxpected tokend %s", t.value)
	}
}

// this function needs some tests and cleaning up
func (p *parser) parseOffsetAddress() (byte, byte, uint16, error) {
	_, err := p.expect(tokenTypeLeftParen)
	if err != nil {
		return 0, 0, 0, err
	}

	t, err := p.expect(tokenTypeText)
	if err != nil {
		return 0, 0, 0, err
	}

	mode := byte(instructions.Immediate)
	modeArg := byte(0)

	isLabel := false
	var labelAddress uint16

	if t.value == "fp" {
		mode = instructions.FramePointerWithOffset
	} else {
		addr, exists := p.labels[t.value]
		if !exists {
			return 0, 0, 0, fmt.Errorf("no definitions found for label %s", t.value)
		}

		isLabel = true
		labelAddress = addr
	}

	next := p.nextToken()

	// get the offset
	switch next.tokenType {
	case tokenTypePlus:
		t := p.nextToken()

		switch t.tokenType {
		case tokenTypeInteger:
			n, err := parseIntegerToken(t)
			if err != nil {
				return 0, 0, 0, err
			}

			modeArg = byte(n)
		case tokenTypeText:
			reg, err := getRegister(t.value)
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
			return 0, 0, 0, fmt.Errorf("unexpected token \"%s\"", t.value)
		}

	case tokenTypeMinus:
		t := p.nextToken()

		switch t.tokenType {
		case tokenTypeInteger:
			n, err := parseIntegerToken(t)
			if err != nil {
				return 0, 0, 0, err
			}

			nsigned := int8(n) * int8(-1)
			modeArg = byte(nsigned)

		case tokenTypeText:
			reg, err := getRegister(t.value)
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
			return 0, 0, 0, fmt.Errorf("unexpected token \"%s\"", t.value)
		}

	// handle offset from [n] or [reg]
	case tokenTypeLeftBracket:
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
		return 0, 0, 0, fmt.Errorf("unexpected token \"%s\"", next.value)
	}

	p.expect(tokenTypeRightParen)

	return mode, modeArg, labelAddress, nil
}

func (p *parser) parseMemoryAddress() (byte, byte, byte, byte, error) {
	t := p.nextToken()

	switch t.tokenType {
	case tokenTypeInteger:
		n, err := parseIntegerToken(t)
		if err != nil {
			return 0, 0, 0, 0, err
		}

		h := byte((n & 0xff00) >> 8)
		l := byte(n & 0xff)

		return instructions.Immediate, 0, h, l, nil

	case tokenTypeLeftParen:
		p.backup()
		mode, offset, addr, err := p.parseOffsetAddress()

		if err != nil {
			return 0, 0, 0, 0, err
		}

		h := byte((addr & 0xff00) >> 8)
		l := byte(addr & 0xff)

		return mode, offset, h, l, nil

	case tokenTypeText:
		if t.value == "fp" {
			return instructions.FramePointerWithOffset, 0, 0, 0, nil
		}

		addr, err := p.getLabelAddress(t.value)
		if err != nil {
			return 0, 0, 0, 0, err
		}

		h := byte((addr & 0xff00) >> 8)
		l := byte(addr & 0xff)

		return instructions.Immediate, 0, h, l, nil

	default:
		return 0, 0, 0, 0, fmt.Errorf("unexpected token \"%s\"", t.value)
	}
}

func (p *parser) parseMov() error {
	p.addByte(instructions.Mov)

	dest := p.nextToken()

	if dest.tokenType != tokenTypeText {
		return fmt.Errorf("unexpected token %s", dest.tokenType)
	}

	reg, err := getRegister(dest.value)
	if err != nil {
		return err
	}

	p.addByte(reg)

	_, err = p.expect(tokenTypeComma)
	if err != nil {
		return err
	}

	t, err := p.expect(tokenTypeInteger)
	if err != nil {
		return err
	}

	n, err := parseIntegerToken(t)
	if err != nil {
		return err
	}

	p.addByte(byte(n))

	p.skipIf(tokenTypeNewLine)
	return nil
}

func (p *parser) parseLoad() error {
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
	_, err = p.expect(tokenTypeComma)
	if err != nil {
		return err
	}

	// register
	dest := p.nextToken()
	if dest.tokenType != tokenTypeText {
		return fmt.Errorf("expected register, got %s", dest.tokenType)
	}

	reg, err := getRegister(dest.value)
	if err != nil {
		return err
	}

	p.addByte(reg)

	p.skipIf(tokenTypeNewLine)
	return nil
}

func (p *parser) parseStore() error {
	p.addByte(instructions.Store)

	// register
	dest := p.nextToken()

	if dest.tokenType != tokenTypeText {
		return fmt.Errorf("expected register, got %s", dest.tokenType)
	}

	reg, err := getRegister(dest.value)
	if err != nil {
		return err
	}

	p.addByte(reg)

	// ,
	_, err = p.expect(tokenTypeComma)
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

	p.skipIf(tokenTypeNewLine)
	return nil
}
