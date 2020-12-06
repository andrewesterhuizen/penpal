package parser

import (
	"fmt"
	"strconv"

	"github.com/andrewesterhuizen/penpal/instructions"
	"github.com/andrewesterhuizen/penpal/lexer_rewrite"
)

type Parser struct {
	index        int
	tokens       []lexer_rewrite.Token
	instructions []byte

	entryFileName string
	files         map[string][]lexer_rewrite.Token

	currentLableAddress uint16
	labels              map[string]uint16
}

func NewParser() *Parser {
	p := Parser{}
	p.labels = map[string]uint16{}
	return &p
}

func (p *Parser) Load(entryFileName string, files map[string][]lexer_rewrite.Token) {
	p.entryFileName = entryFileName
	p.files = files
}

func parseIntegerToken(t lexer_rewrite.Token) (uint64, error) {
	if t.Type != lexer_rewrite.TokenTypeInteger {
		return 0, fmt.Errorf("expected token to be integer and got %s", t.Type)
	}

	n, err := strconv.ParseUint(t.Value, 0, 64)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (p *Parser) addByte(i byte) {
	p.instructions = append(p.instructions, i)
}

func (p *Parser) nextToken() lexer_rewrite.Token {
	p.index++

	if p.index >= len(p.tokens) {
		return lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile}
	}

	return p.tokens[p.index]
}

func (p *Parser) backup() {
	p.index--
	return
}

func (p *Parser) expect(t lexer_rewrite.TokenType) (lexer_rewrite.Token, error) {
	n := p.nextToken()
	if n.Type != t {
		return n, fmt.Errorf("expected %v and got %v", t, n.Type)
	}

	return n, nil
}
func (p *Parser) peek() lexer_rewrite.Token {
	nextIndex := p.index + 1

	if nextIndex >= len(p.tokens) {
		return lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile}
	}

	return p.tokens[nextIndex]
}

func (p *Parser) parseInstruction(t lexer_rewrite.Token) error {
	switch t.Value {
	case "swap":
		return p.parseNoOperandInstruction(instructions.Swap)
	case "pop":
		return p.parseNoOperandInstruction(instructions.Pop)
	case "ret":
		return p.parseNoOperandInstruction(instructions.Ret)
	case "reti":
		return p.parseNoOperandInstruction(instructions.Reti)
	case "push":
		return p.parseArithmeticLogicInstruction(instructions.Push)
	case "add":
		return p.parseArithmeticLogicInstruction(instructions.Add)
	case "sub":
		return p.parseArithmeticLogicInstruction(instructions.Sub)
	case "mul":
		return p.parseArithmeticLogicInstruction(instructions.Mul)
	case "div":
		return p.parseArithmeticLogicInstruction(instructions.Div)
	case "shl":
		return p.parseArithmeticLogicInstruction(instructions.Shl)
	case "shr":
		return p.parseArithmeticLogicInstruction(instructions.Shr)
	case "and":
		return p.parseArithmeticLogicInstruction(instructions.And)
	case "or":
		return p.parseArithmeticLogicInstruction(instructions.Or)
	case "gt":
		return p.parseArithmeticLogicInstruction(instructions.GT)
	case "gte":
		return p.parseArithmeticLogicInstruction(instructions.GTE)
	case "lt":
		return p.parseArithmeticLogicInstruction(instructions.LT)
	case "lte":
		return p.parseArithmeticLogicInstruction(instructions.LTE)
	case "eq":
		return p.parseArithmeticLogicInstruction(instructions.Eq)
	case "neq":
		return p.parseArithmeticLogicInstruction(instructions.Neq)
	case "rand":
		return p.parseArithmeticLogicInstruction(instructions.Rand)
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
		return fmt.Errorf("unexpected instruction %v", t.Value)
	}

}

func (p *Parser) parseToken(t lexer_rewrite.Token) error {
	var err error
	switch t.Type {
	case lexer_rewrite.TokenTypeInstruction:
		err = p.parseInstruction(t)
	case lexer_rewrite.TokenTypeLabel:
		n := p.peek()

		// newline is optional
		if n.Type == lexer_rewrite.TokenTypeNewLine {
			p.index++
		}

	default:
		return fmt.Errorf("encountered unknown token %v", t.Type)
	}

	if err != nil {
		return fmt.Errorf("error parsing \"%s\" instruction: %w", t.Value, err)
	}

	return nil
}

func (p *Parser) parseTokens() error {
	for t := p.tokens[p.index]; t.Type != lexer_rewrite.TokenTypeEndOfFile; t = p.nextToken() {
		if err := p.parseToken(t); err != nil {
			return fmt.Errorf("[%d:%d] %s", t.Line, t.Column, err)
		}
	}

	return nil
}

func (p *Parser) processIncludes(filename string) ([]lexer_rewrite.Token, error) {
	tokens := p.files[filename]
	out := []lexer_rewrite.Token{}

	for _, t := range tokens {
		switch t.Type {
		case lexer_rewrite.TokenTypeFileInclude:
			includeTokens, err := p.processIncludes(t.Value)
			if err != nil {
				return nil, err
			}

			out = append(out, includeTokens...)
		case lexer_rewrite.TokenTypeSystemInclude:

		default:
			out = append(out, t)
		}
	}

	return out, nil
}

func (p *Parser) getLabels() error {
	for _, t := range p.tokens {
		switch t.Type {
		case lexer_rewrite.TokenTypeLabel:
			p.labels[t.Value] = uint16(p.currentLableAddress)

		case lexer_rewrite.TokenTypeInstruction:
			ins := instructions.InstructionByName[t.Value]
			w := instructions.WidthNew[ins]
			p.currentLableAddress += uint16(w)
		}
	}

	return nil
}

func (p *Parser) Run() ([]byte, error) {
	tokens, err := p.processIncludes(p.entryFileName)
	if err != nil {
		return nil, err
	}

	p.tokens = tokens

	err = p.getLabels()
	if err != nil {
		return nil, err
	}

	err = p.parseTokens()
	if err != nil {
		return nil, err
	}

	return p.instructions, nil
}
