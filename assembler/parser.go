package assembler

import (
	"fmt"
	"strconv"

	"github.com/andrewesterhuizen/penpal/instructions"
)

type parser struct {
	index               int
	tokens              []token
	instructions        []byte
	currentLableAddress uint16
	labels              map[string]uint16
}

func newParser() *parser {
	p := parser{}
	p.labels = map[string]uint16{}
	return &p
}

func parseIntegerToken(t token) (uint64, error) {
	if t.tokenType != tokenTypeInteger {
		return 0, fmt.Errorf("expected token to be integer and got %s", t.tokenType)
	}

	n, err := strconv.ParseUint(t.value, 0, 64)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (p *parser) addByte(i byte) {
	p.instructions = append(p.instructions, i)
}

func (p *parser) nextToken() token {
	p.index++

	if p.index >= len(p.tokens) {
		return token{tokenType: tokenTypeEndOfFile}
	}

	return p.tokens[p.index]
}

func (p *parser) backup() {
	p.index--
	return
}

func (p *parser) expect(t tokenType) (token, error) {
	n := p.nextToken()
	if n.tokenType != t {
		return n, fmt.Errorf("expected %v and got %v", t, n.tokenType)
	}

	return n, nil
}

func (p *parser) expectRange(types []tokenType) (token, error) {
	n := p.nextToken()

	for _, t := range types {
		if n.tokenType == t {
			return n, nil
		}
	}

	return n, fmt.Errorf("unexpected token %s", n.value)
}

func (p *parser) skipIf(t tokenType) {
	n := p.nextToken()
	if n.tokenType != t {
		p.backup()
	}
}

func (p *parser) peek() token {
	nextIndex := p.index + 1

	if nextIndex >= len(p.tokens) {
		return token{tokenType: tokenTypeEndOfFile}
	}

	return p.tokens[nextIndex]
}

func (p *parser) getLabelAddress(l string) (uint16, error) {
	addr, exists := p.labels[l]
	if !exists {
		return 0, fmt.Errorf("no definition found for label %s", l)
	}

	return addr, nil
}

func (p *parser) parseInstruction(t token) error {
	switch t.value {
	case "swap":
		return p.parseNoOperandInstruction(instructions.Swap)
	case "pop":
		return p.parseNoOperandInstruction(instructions.Pop)
	case "ret":
		return p.parseNoOperandInstruction(instructions.Ret)
	case "reti":
		return p.parseNoOperandInstruction(instructions.Reti)
	case "push":
		return p.parsePushInstruction()
	case "add":
		return p.parseNoOperandInstruction(instructions.Add)
	case "sub":
		return p.parseNoOperandInstruction(instructions.Sub)
	case "mul":
		return p.parseNoOperandInstruction(instructions.Mul)
	case "div":
		return p.parseNoOperandInstruction(instructions.Div)
	case "shl":
		return p.parseNoOperandInstruction(instructions.Shl)
	case "shr":
		return p.parseNoOperandInstruction(instructions.Shr)
	case "and":
		return p.parseNoOperandInstruction(instructions.And)
	case "or":
		return p.parseNoOperandInstruction(instructions.Or)
	case "gt":
		return p.parseNoOperandInstruction(instructions.GT)
	case "gte":
		return p.parseNoOperandInstruction(instructions.GTE)
	case "lt":
		return p.parseNoOperandInstruction(instructions.LT)
	case "lte":
		return p.parseNoOperandInstruction(instructions.LTE)
	case "eq":
		return p.parseNoOperandInstruction(instructions.Eq)
	case "neq":
		return p.parseNoOperandInstruction(instructions.Neq)
	case "rand":
		return p.parseNoOperandInstruction(instructions.Rand)
	case "halt":
		return p.parseNoOperandInstruction(instructions.Halt)
	case "call":
		return p.parseAddressInstruction(instructions.Call)
	case "jump":
		return p.parseAddressInstruction(instructions.Jump)
	case "jumpz":
		return p.parseAddressInstruction(instructions.Jumpz)
	case "jumpnz":
		return p.parseAddressInstruction(instructions.Jumpnz)
	case "load":
		return p.parseLoad()
	case "store":
		return p.parseStore()
	case "mov":
		return p.parseMov()
	case "db":
		return p.parseDB()
	default:
		return fmt.Errorf("unexpected instruction %v", t.value)
	}

}

func (p *parser) parseToken(t token) error {
	var err error
	switch t.tokenType {
	case tokenTypeNewLine:
		// skip whitespace
	case tokenTypeInstruction:
		err = p.parseInstruction(t)
	case tokenTypeLabel:
		n := p.peek()

		// newline is optional
		if n.tokenType == tokenTypeNewLine {
			p.index++
		}

	default:
		return fmt.Errorf("unexpected token %v", t.value)
	}

	if err != nil {
		return fmt.Errorf("error parsing \"%s\" instruction: %w", t.value, err)
	}

	return nil
}

func (p *parser) parseTokens() error {
	for t := p.tokens[p.index]; t.tokenType != tokenTypeEndOfFile; t = p.nextToken() {
		if err := p.parseToken(t); err != nil {
			return fmt.Errorf("[%s:%d:%d] %s", t.fileName, t.line, t.column, err)
		}
	}

	return nil
}

func (p *parser) getLabels() error {
	for _, t := range p.tokens {
		switch t.tokenType {
		case tokenTypeLabel:
			p.labels[t.value] = uint16(p.currentLableAddress)

		case tokenTypeInstruction:
			ins := instructions.InstructionByName[t.value]
			w := instructions.Width[ins]
			p.currentLableAddress += uint16(w)
		}
	}

	return nil
}

func (p *parser) Run(tokens []token) ([]byte, error) {
	p.tokens = tokens

	err := p.getLabels()
	if err != nil {
		return nil, err
	}

	err = p.parseTokens()
	if err != nil {
		return nil, err
	}

	return p.instructions, nil
}
