package parser

import (
	"fmt"
	"strconv"

	"github.com/andrewesterhuizen/penpal/lexer_rewrite"
)

type Parser struct {
	index        int
	tokens       []lexer_rewrite.Token
	instructions []byte
}

func NewParser() *Parser {
	p := Parser{}
	return &p
}

func (p *Parser) Load(tokens []lexer_rewrite.Token) {
	p.index = 0
	p.tokens = tokens
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

func (p *Parser) appendInstruction(i byte) {
	p.instructions = append(p.instructions, i)
}

func (p *Parser) nextToken() lexer_rewrite.Token {
	p.index++
	return p.tokens[p.index]
}

func (p *Parser) accept(t lexer_rewrite.TokenType) error {
	n := p.nextToken()
	if n.Type != t {
		return fmt.Errorf("expected %v and got %v", t, n.Type)
	}

	return nil
}

func (p *Parser) expect(t lexer_rewrite.TokenType) (lexer_rewrite.Token, error) {
	n := p.nextToken()
	if n.Type != t {
		return n, fmt.Errorf("expected %v and got %v", t, n.Type)
	}

	return n, nil
}

func (p *Parser) parseInstruction(t lexer_rewrite.Token) error {
	switch t.Value {
	case "add":
		return p.parseAdd()
	case "mov":
		return p.parseMov()
	default:
		return fmt.Errorf("unexpected instruction %v", t.Value)
	}

}

func (p *Parser) parseToken(t lexer_rewrite.Token) error {
	var err error
	switch t.Type {
	case lexer_rewrite.TokenTypeText:
		err = p.parseInstruction(t)
	default:
		return fmt.Errorf("encountered unknown token %v", t.Type)
	}

	if err != nil {
		return fmt.Errorf("error parsing (%s) instruction: %w", t.Value, err)
	}

	return nil
}

func (p *Parser) Run() ([]byte, error) {
	for t := p.tokens[p.index]; t.Type != lexer_rewrite.TokenTypeEndOfFile; t = p.nextToken() {
		if err := p.parseToken(t); err != nil {
			return nil, err
		}
	}

	return p.instructions, nil
}
