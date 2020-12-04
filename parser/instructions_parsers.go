package parser

import (
	"fmt"

	"github.com/andrewesterhuizen/penpal/instructions"
	"github.com/andrewesterhuizen/penpal/lexer_rewrite"
)

func (p *Parser) parseNoOperandInstruction(instruction byte) error {
	p.addByte(instruction)
	return p.expect(lexer_rewrite.TokenTypeNewLine)
}

// encoding = instruction, (immediate mode|register mode), (immediate value|register number)
func (p *Parser) parseArithmeticLogicInstruction(instruction byte) error {
	p.addByte(instruction)

	op1 := p.nextToken()

	switch op1.Type {
	// no operand = implied B register
	case lexer_rewrite.TokenTypeNewLine:
		p.addByte(instructions.Register)
		p.addByte(instructions.RegisterB)

		return nil
	case lexer_rewrite.TokenTypeInteger:
		p.addByte(instructions.Immediate)

		n, err := parseIntegerToken(op1)
		if err != nil {
			return err
		}

		p.addByte(byte(n))

	default:
		return fmt.Errorf("unexpected operand \"%s\"", op1.Value)
	}

	return p.expect(lexer_rewrite.TokenTypeNewLine)
}

// encoding = mov, value, mode, dest
func (p *Parser) parseMov() error {
	p.addByte(instructions.Mov)

	t := p.nextToken()

	switch t.Type {
	case lexer_rewrite.TokenTypeInteger:
		n, err := parseIntegerToken(t)
		if err != nil {
			return err
		}

		p.addByte(byte(n))
		p.addByte(instructions.AddressingModeImmediate)

	case lexer_rewrite.TokenTypeLeftParen:
		register, err := p.expectAndGet(lexer_rewrite.TokenTypeText)
		if err != nil {
			return err
		}

		mode := byte(instructions.Immediate)

		switch register.Value {
		case "fp":
			mode = instructions.AddressingModeFPRelative
		default:
			return fmt.Errorf("unexpected register \"%s\"", t.Value)
		}

		next := p.nextToken()

		switch next.Type {
		case lexer_rewrite.TokenTypePlus:
			value, err := p.expectAndGet(lexer_rewrite.TokenTypeInteger)
			if err != nil {
				return err
			}

			n, err := parseIntegerToken(value)
			if err != nil {
				return err
			}

			p.addByte(byte(n))

		case lexer_rewrite.TokenTypeMinus:
			value, err := p.expectAndGet(lexer_rewrite.TokenTypeInteger)
			if err != nil {
				return err
			}

			n, err := parseIntegerToken(value)
			if err != nil {
				return err
			}
			nsigned := int8(n) * int8(-1)
			p.addByte(byte(nsigned))

		case lexer_rewrite.TokenTypeLeftBracket:
			t := p.nextToken()

			// default to +
			operator := lexer_rewrite.TokenTypePlus

			// valid next tokens here are, +, - or the integer

			// check if token is - for inverting later
			if t.Type == lexer_rewrite.TokenTypeMinus {
				operator = lexer_rewrite.TokenTypeMinus
			}

			intToken := t

			// if current token isn't the integer then it must be the next token
			if t.Type != lexer_rewrite.TokenTypeInteger {
				intToken, err = p.expectAndGet(lexer_rewrite.TokenTypeInteger)
				if err != nil {
					return err
				}
			}

			err = p.expect(lexer_rewrite.TokenTypeRightBracket)
			if err != nil {
				return err
			}

			n, err := parseIntegerToken(intToken)
			if err != nil {
				return err
			}

			if operator == lexer_rewrite.TokenTypeMinus {
				nsigned := int8(n) * int8(-1)
				p.addByte(byte(nsigned))
			} else {
				p.addByte(byte(n))
			}

		default:
			return fmt.Errorf("unexpected token \"%s\"", next.Value)
		}

		p.addByte(mode)

		p.expect(lexer_rewrite.TokenTypeRightParen)
	case lexer_rewrite.TokenTypeText:
		switch t.Value {
		case "fp":
			// mov fp, A => (A = (fp+0))
			p.addByte(0)
			p.addByte(instructions.AddressingModeFPRelative)

		default:
			return fmt.Errorf("unexpected register \"%s\"", t.Value)
		}

	default:
		return fmt.Errorf("unexpected token \"%s\"", t.Value)
	}

	err := p.expect(lexer_rewrite.TokenTypeComma)
	if err != nil {
		return err
	}

	dest := p.nextToken()

	if dest.Type != lexer_rewrite.TokenTypeText {
		return fmt.Errorf("unexpected token %s", dest.Type)
	}

	switch dest.Value {
	case "A":
		p.addByte(instructions.RegisterA)
	case "B":
		p.addByte(instructions.RegisterB)
	default:
		return fmt.Errorf("expected register and got %s", dest.Value)
	}

	return p.expect(lexer_rewrite.TokenTypeNewLine)
}

// load 0xae, A
// instructions.Load, instructions.Register, instructions.RegisterA, 0x00, 0xae
func (p *Parser) parseLoad() error {
	p.addByte(instructions.Load)

	var mode uint8
	var valueH byte
	var valueL byte

	t := p.nextToken()

	switch t.Type {
	case lexer_rewrite.TokenTypeInteger:
		mode = instructions.Register
		p.addByte(instructions.Register)

		n, err := parseIntegerToken(t)
		if err != nil {
			return err
		}

		value := uint16(n)

		valueH = byte((value & 0xff00) >> 8)
		valueL = byte(value & 0xff)

	default:
		return fmt.Errorf("unexpected token \"%s\"", t.Value)
	}

	err := p.expect(lexer_rewrite.TokenTypeComma)
	if err != nil {
		return err
	}

	dest := p.nextToken()

	if dest.Type != lexer_rewrite.TokenTypeText {
		return fmt.Errorf("unexpected token %s", dest.Type)
	}

	switch dest.Value {
	case "A":
		p.addByte(instructions.RegisterA)
	case "B":
		p.addByte(instructions.RegisterB)
	default:
		return fmt.Errorf("expected register and got %s", dest.Value)
	}

	switch mode {
	case instructions.AddressingModeImmediate:
		p.addByte(valueH)
		p.addByte(valueL)
	case instructions.AddressingModeFPRelative:
		p.addByte(valueH)
		p.addByte(valueL)
	}

	return p.expect(lexer_rewrite.TokenTypeNewLine)
}
