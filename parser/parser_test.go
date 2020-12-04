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

func newBasicInstructionTestCase(text string, ins byte) parserTestCase {
	input := []lexer_rewrite.Token{
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: text},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}

	return parserTestCase{
		input:  input,
		output: []byte{ins},
	}
}

func newALTests(text string, ins byte) []parserTestCase {
	return []parserTestCase{
		{input: []lexer_rewrite.Token{
			// ins
			lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: text},
			lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine},
			lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
		}, output: []byte{ins, instructions.Register, instructions.RegisterB}},
		{input: []lexer_rewrite.Token{
			// ins 0xbb
			lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: text},
			lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInteger, Value: "0xbb"},
			lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine},
			lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
		}, output: []byte{ins, instructions.Immediate, 0xbb}},
	}
}

var singleOperandInstructionTestCases = []parserTestCase{
	newBasicInstructionTestCase("swap", instructions.Swap),
	newBasicInstructionTestCase("pop", instructions.Pop),
	newBasicInstructionTestCase("ret", instructions.Ret),
	newBasicInstructionTestCase("reti", instructions.Reti),
	newBasicInstructionTestCase("halt", instructions.Halt),
}

var movTestCases = []parserTestCase{
	{input: []lexer_rewrite.Token{
		// mov 0xab, A
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "mov"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInteger, Value: "0xab"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeComma},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "A"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Mov, 0xab, instructions.AddressingModeImmediate, instructions.RegisterA}},
	{input: []lexer_rewrite.Token{
		// mov fp, A
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "mov"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "fp"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeComma, Value: ","},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "A"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine, Value: "\n"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Mov, 0x0, instructions.AddressingModeFPRelative, instructions.RegisterA}},
	{input: []lexer_rewrite.Token{
		// mov (fp+1), A
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "mov"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeLeftParen, Value: "("},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "fp"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypePlus, Value: "+"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInteger, Value: "1"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeRightParen, Value: ")"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeComma, Value: ","},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "A"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine, Value: "\n"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Mov, 0x1, instructions.AddressingModeFPRelative, instructions.RegisterA}},
	{input: []lexer_rewrite.Token{
		// mov (fp-1), B
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "mov"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeLeftParen, Value: "("},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "fp"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeMinus, Value: "-"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInteger, Value: "1"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeRightParen, Value: ")"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeComma, Value: ","},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "B"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine, Value: "\n"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Mov, 0xff, instructions.AddressingModeFPRelative, instructions.RegisterB}},
	{input: []lexer_rewrite.Token{
		// mov (fp[3]), B
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "mov"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeLeftParen, Value: "("},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "fp"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeLeftBracket, Value: "["},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInteger, Value: "3"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeRightBracket, Value: "]"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeRightParen, Value: ")"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeComma, Value: ","},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "B"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine, Value: "\n"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Mov, 0x3, instructions.AddressingModeFPRelative, instructions.RegisterB}},
	{input: []lexer_rewrite.Token{
		// mov (fp[+3]), A
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "mov"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeLeftParen, Value: "("},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "fp"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeLeftBracket, Value: "["},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypePlus, Value: "+"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInteger, Value: "3"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeRightBracket, Value: "]"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeRightParen, Value: ")"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeComma, Value: ","},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "A"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine, Value: "\n"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Mov, 0x3, instructions.AddressingModeFPRelative, instructions.RegisterA}},
	{input: []lexer_rewrite.Token{
		// mov (fp[-3]), A
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "mov"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeLeftParen, Value: "("},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "fp"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeLeftBracket, Value: "["},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeMinus, Value: "-"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInteger, Value: "3"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeRightBracket, Value: "]"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeRightParen, Value: ")"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeComma, Value: ","},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "A"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine, Value: "\n"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Mov, 0xfd, instructions.AddressingModeFPRelative, instructions.RegisterA}},
}

var loadTestCases = []parserTestCase{
	// load src (address|label), dest (register)
	{input: []lexer_rewrite.Token{
		// load 0xae, A
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "load"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInteger, Value: "0xae"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeComma},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "A"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Load, instructions.Register, instructions.RegisterA, 0x00, 0xae}},
	{input: []lexer_rewrite.Token{
		// load 0xaecd, B
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "load"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInteger, Value: "0xaecd"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeComma},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "B"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Load, instructions.Register, instructions.RegisterB, 0xae, 0xcd}},
}

func TestParser(t *testing.T) {
	var parserTestCases = []parserTestCase{}
	// arithmetic and logic
	// no operand = implied B register as value
	// instruction, (immediate mode|register mode), (immediate value|register number)
	parserTestCases = append(parserTestCases, newALTests("add", instructions.Add)...)
	parserTestCases = append(parserTestCases, newALTests("sub", instructions.Sub)...)
	parserTestCases = append(parserTestCases, newALTests("mul", instructions.Mul)...)
	parserTestCases = append(parserTestCases, newALTests("div", instructions.Div)...)
	parserTestCases = append(parserTestCases, newALTests("or", instructions.Or)...)
	parserTestCases = append(parserTestCases, newALTests("shl", instructions.Shl)...)
	parserTestCases = append(parserTestCases, newALTests("shr", instructions.Shr)...)
	parserTestCases = append(parserTestCases, newALTests("rand", instructions.Rand)...)

	parserTestCases = append(parserTestCases, movTestCases...)

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
				t.Errorf("expected byte 0x%02x and got 0x%02x", ins, out[i])
			}

		}
	}
}
