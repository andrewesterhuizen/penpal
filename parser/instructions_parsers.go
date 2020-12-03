package parser

import (
	"fmt"

	"github.com/andrewesterhuizen/penpal/instructions"
	"github.com/andrewesterhuizen/penpal/lexer_rewrite"
)

func (p *Parser) parseAdd() error {
	p.appendInstruction(instructions.Add)
	return p.accept(lexer_rewrite.TokenTypeNewLine)
}

func (p *Parser) parseMove() error {
	p.appendInstruction(instructions.Mov)

	var mode uint8
	var value byte

	t := p.nextToken()

	switch t.Type {
	case lexer_rewrite.TokenTypeInteger:
		mode = instructions.AddressingModeImmediate
		p.appendInstruction(instructions.AddressingModeImmediate)

		n, err := parseIntegerToken(t)
		if err != nil {
			return err
		}

		value = byte(n)

	case lexer_rewrite.TokenTypeLeftParen:
		register, err := p.expect(lexer_rewrite.TokenTypeText)
		if err != nil {
			return err
		}

		switch register.Value {
		case "fp":
			mode = instructions.AddressingModeFPRelative
			p.appendInstruction(instructions.AddressingModeFPRelative)
		default:
			return fmt.Errorf("unexpected register \"%s\"", t.Value)
		}

		operator := p.nextToken()

		v, err := p.expect(lexer_rewrite.TokenTypeInteger)
		if err != nil {
			return err
		}

		switch operator.Type {
		case lexer_rewrite.TokenTypePlus:
			n, err := parseIntegerToken(v)
			if err != nil {
				return err
			}

			value = byte(n)

		case lexer_rewrite.TokenTypeMinus:
			n, err := parseIntegerToken(v)
			if err != nil {
				return err
			}
			nsigned := int8(n) * int8(-1)
			value = byte(nsigned)
		default:
			return fmt.Errorf("unexpected operator \"%s\"", t.Value)
		}

		p.accept(lexer_rewrite.TokenTypeRightParen)

	default:
		return fmt.Errorf("unexpected token \"%s\"", t.Value)
	}

	err := p.accept(lexer_rewrite.TokenTypeComma)
	if err != nil {
		return err
	}

	dest := p.nextToken()

	if dest.Type != lexer_rewrite.TokenTypeText {
		return fmt.Errorf("unexpected token %s", dest.Type)
	}

	switch dest.Value {
	case "A":
		p.appendInstruction(instructions.RegisterA)
	case "B":
		p.appendInstruction(instructions.RegisterA)
	default:
		return fmt.Errorf("expected register and got %s", dest.Value)
	}

	switch mode {
	case instructions.AddressingModeImmediate:
		p.appendInstruction(value)
	case instructions.AddressingModeFPRelative:
		p.appendInstruction(value)
	}

	return p.accept(lexer_rewrite.TokenTypeNewLine)
}
