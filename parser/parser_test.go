package parser

import (
	"fmt"
	"testing"

	"github.com/andrewesterhuizen/penpal/instructions"
	"github.com/andrewesterhuizen/penpal/lexer_rewrite"
)

type parserTestCase struct {
	input  []lexer_rewrite.Token
	output []byte
}

var parserTestCases = []parserTestCase{
	{input: []lexer_rewrite.Token{ // add
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "add"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Add}},
	{input: []lexer_rewrite.Token{ // move 0xab, A
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "move"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInteger, Value: "0xab"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeComma},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "A"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Mov, instructions.AddressingModeImmediate, instructions.RegisterA, 0xab}},
	{input: []lexer_rewrite.Token{ // move (fp+1), A
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "move"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeLeftParen, Value: "("},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "fp"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypePlus, Value: "+"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInteger, Value: "1"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeRightParen, Value: ")"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeComma, Value: ","},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "A"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine, Value: "\n"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Mov, instructions.AddressingModeFPRelative, instructions.RegisterA, 0x1}},
	{input: []lexer_rewrite.Token{ // move (fp-1), A
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "move"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeLeftParen, Value: "("},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "fp"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeMinus, Value: "-"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInteger, Value: "1"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeRightParen, Value: ")"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeComma, Value: ","},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "A"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine, Value: "\n"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Mov, instructions.AddressingModeFPRelative, instructions.RegisterA, 0xff}},
}

func TestParser(t *testing.T) {

	for _, tc := range parserTestCases {
		p := NewParser()
		p.Load(tc.input)
		out, err := p.Run()
		if err != nil {
			t.Error(err)
			return
		}

		fmt.Println(out)

		if len(tc.output) != len(out) {
			t.Errorf("expected %v bytes and got %v", len(tc.output), len(out))
			return
		}

		for i, ins := range tc.output {
			if ins != out[i] {
				t.Errorf("expected %02x and got %02x", ins, out[i])
			}

		}
	}
}
