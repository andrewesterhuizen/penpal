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

func tok(t lexer_rewrite.TokenType, v string) lexer_rewrite.Token {
	return lexer_rewrite.Token{Type: t, Value: v}
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

func newAritmeticLogicInstructionTests(text string, ins byte) []parserTestCase {
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
	}, output: []byte{instructions.Mov, instructions.AddressingModeImmediate, 0xab, instructions.RegisterA}},
	{input: []lexer_rewrite.Token{
		// mov fp, A
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "mov"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "fp"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeComma, Value: ","},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "A"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine, Value: "\n"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Mov, instructions.FramePointerWithOffset, 0x0, instructions.RegisterA}},
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
	}, output: []byte{instructions.Mov, instructions.FramePointerWithOffset, 0x1, instructions.RegisterA}},
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
	}, output: []byte{instructions.Mov, instructions.FramePointerWithOffset, 0xff, instructions.RegisterB}},
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
	}, output: []byte{instructions.Mov, instructions.FramePointerWithOffset, 0x3, instructions.RegisterB}},
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
	}, output: []byte{instructions.Mov, instructions.FramePointerWithOffset, 0x3, instructions.RegisterA}},
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
	}, output: []byte{instructions.Mov, instructions.FramePointerWithOffset, 0xfd, instructions.RegisterA}},
}

// TODO:
// load label, A           (A = memory[label])
// load (label + A), A     (A = memory[label + A])
// load (label - A), A     (A = memory[label - A])
// load (label + 5), A     (A = memory[label + 5])
// load (label - 5), A     (A = memory[label - 5])
// load (label[5]), A      (A = memory[label + 5])

// load src (address|label), dest (register)
var loadTestCases = []parserTestCase{
	{input: []lexer_rewrite.Token{
		// load 0xae, A
		tok(lexer_rewrite.TokenTypeText, "load"),
		tok(lexer_rewrite.TokenTypeInteger, "0xae"),
		tok(lexer_rewrite.TokenTypeComma, ","),
		tok(lexer_rewrite.TokenTypeText, "A"),
		tok(lexer_rewrite.TokenTypeNewLine, "\n"),
		tok(lexer_rewrite.TokenTypeEndOfFile, ""),
	}, output: []byte{instructions.Load, 0x00, 0xae, instructions.Immediate, 0x0, instructions.RegisterA}},
	{input: []lexer_rewrite.Token{
		// load 0xaecd, B
		tok(lexer_rewrite.TokenTypeText, "load"),
		tok(lexer_rewrite.TokenTypeInteger, "0xaecd"),
		tok(lexer_rewrite.TokenTypeComma, ""),
		tok(lexer_rewrite.TokenTypeText, "B"),
		tok(lexer_rewrite.TokenTypeNewLine, ""),
		tok(lexer_rewrite.TokenTypeEndOfFile, ""),
	}, output: []byte{instructions.Load, 0xae, 0xcd, instructions.Immediate, 0x0, instructions.RegisterB}},
	{input: []lexer_rewrite.Token{
		// load (fp + 5), A
		tok(lexer_rewrite.TokenTypeText, "load"),
		tok(lexer_rewrite.TokenTypeLeftParen, "("),
		tok(lexer_rewrite.TokenTypeText, "fp"),
		tok(lexer_rewrite.TokenTypePlus, "+"),
		tok(lexer_rewrite.TokenTypeInteger, "5"),
		tok(lexer_rewrite.TokenTypeRightParen, ")"),
		tok(lexer_rewrite.TokenTypeComma, ","),
		tok(lexer_rewrite.TokenTypeText, "A"),
		tok(lexer_rewrite.TokenTypeNewLine, "\n"),
		tok(lexer_rewrite.TokenTypeEndOfFile, ""),
	}, output: []byte{instructions.Load, 0x00, 0x00, instructions.FramePointerWithOffset, 0x5, instructions.RegisterA}},
	{input: []lexer_rewrite.Token{
		// load (fp + 5), A
		tok(lexer_rewrite.TokenTypeText, "load"),
		tok(lexer_rewrite.TokenTypeLeftParen, "("),
		tok(lexer_rewrite.TokenTypeText, "fp"),
		tok(lexer_rewrite.TokenTypePlus, "+"),
		tok(lexer_rewrite.TokenTypeInteger, "5"),
		tok(lexer_rewrite.TokenTypeRightParen, ")"),
		tok(lexer_rewrite.TokenTypeComma, ","),
		tok(lexer_rewrite.TokenTypeText, "A"),
		tok(lexer_rewrite.TokenTypeNewLine, "\n"),
		tok(lexer_rewrite.TokenTypeEndOfFile, ""),
	}, output: []byte{instructions.Load, 0x00, 0x00, instructions.FramePointerWithOffset, 0x5, instructions.RegisterA}},
	{input: []lexer_rewrite.Token{
		// load (fp[5]), A
		tok(lexer_rewrite.TokenTypeText, "load"),
		tok(lexer_rewrite.TokenTypeLeftParen, "("),
		tok(lexer_rewrite.TokenTypeText, "fp"),
		tok(lexer_rewrite.TokenTypeLeftBracket, "["),
		tok(lexer_rewrite.TokenTypeInteger, "5"),
		tok(lexer_rewrite.TokenTypeRightBracket, "]"),
		tok(lexer_rewrite.TokenTypeRightParen, ")"),
		tok(lexer_rewrite.TokenTypeComma, ","),
		tok(lexer_rewrite.TokenTypeText, "A"),
		tok(lexer_rewrite.TokenTypeNewLine, "\n"),
		tok(lexer_rewrite.TokenTypeEndOfFile, ""),
	}, output: []byte{instructions.Load, 0x00, 0x00, instructions.FramePointerWithOffset, 0x5, instructions.RegisterA}},
	{input: []lexer_rewrite.Token{
		// load (fp[-1]), B
		tok(lexer_rewrite.TokenTypeText, "load"),
		tok(lexer_rewrite.TokenTypeLeftParen, "("),
		tok(lexer_rewrite.TokenTypeText, "fp"),
		tok(lexer_rewrite.TokenTypeLeftBracket, "["),
		tok(lexer_rewrite.TokenTypeMinus, "-"),
		tok(lexer_rewrite.TokenTypeInteger, "1"),
		tok(lexer_rewrite.TokenTypeRightBracket, "]"),
		tok(lexer_rewrite.TokenTypeRightParen, ")"),
		tok(lexer_rewrite.TokenTypeComma, ","),
		tok(lexer_rewrite.TokenTypeText, "B"),
		tok(lexer_rewrite.TokenTypeNewLine, "\n"),
		tok(lexer_rewrite.TokenTypeEndOfFile, ""),
	}, output: []byte{instructions.Load, 0x00, 0x00, instructions.FramePointerWithOffset, 0xff, instructions.RegisterB}},
	{input: []lexer_rewrite.Token{
		// load (fp - A), A
		tok(lexer_rewrite.TokenTypeText, "load"),
		tok(lexer_rewrite.TokenTypeLeftParen, "("),
		tok(lexer_rewrite.TokenTypeText, "fp"),
		tok(lexer_rewrite.TokenTypeMinus, "-"),
		tok(lexer_rewrite.TokenTypeText, "A"),
		tok(lexer_rewrite.TokenTypeRightParen, ")"),
		tok(lexer_rewrite.TokenTypeComma, ","),
		tok(lexer_rewrite.TokenTypeText, "A"),
		tok(lexer_rewrite.TokenTypeNewLine, "\n"),
		tok(lexer_rewrite.TokenTypeEndOfFile, ""),
	}, output: []byte{instructions.Load, 0x00, 0x00, instructions.FramePointerMinusRegister, instructions.RegisterA, instructions.RegisterA}},
	{input: []lexer_rewrite.Token{
		//  load (fp - 1), A
		tok(lexer_rewrite.TokenTypeText, "load"),
		tok(lexer_rewrite.TokenTypeLeftParen, "("),
		tok(lexer_rewrite.TokenTypeText, "fp"),
		tok(lexer_rewrite.TokenTypeMinus, "-"),
		tok(lexer_rewrite.TokenTypeInteger, "1"),
		tok(lexer_rewrite.TokenTypeRightParen, ")"),
		tok(lexer_rewrite.TokenTypeComma, ","),
		tok(lexer_rewrite.TokenTypeText, "A"),
		tok(lexer_rewrite.TokenTypeNewLine, "\n"),
		tok(lexer_rewrite.TokenTypeEndOfFile, ""),
	}, output: []byte{instructions.Load, 0x00, 0x00, instructions.FramePointerWithOffset, 0xff, instructions.RegisterA}},
	{input: []lexer_rewrite.Token{
		//  load fp, A
		tok(lexer_rewrite.TokenTypeText, "load"),
		tok(lexer_rewrite.TokenTypeText, "fp"),
		tok(lexer_rewrite.TokenTypeComma, ","),
		tok(lexer_rewrite.TokenTypeText, "A"),
		tok(lexer_rewrite.TokenTypeNewLine, "\n"),
		tok(lexer_rewrite.TokenTypeEndOfFile, ""),
	}, output: []byte{instructions.Load, 0x00, 0x00, instructions.FramePointerWithOffset, 0x0, instructions.RegisterA}},
}

// TODO: labels
var callTestCases = []parserTestCase{
	{input: []lexer_rewrite.Token{
		// call 0xba
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "call"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInteger, Value: "0xba"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Call, 0x00, 0xba}},
	{input: []lexer_rewrite.Token{
		// call 0xcdba
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "call"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInteger, Value: "0xcdba"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Call, 0xcd, 0xba}},
}

var jumpTestCases = []parserTestCase{
	{input: []lexer_rewrite.Token{
		// jump 0xba
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "jump"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInteger, Value: "0xba"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Jump, 0x00, 0xba}},
	{input: []lexer_rewrite.Token{
		// jump 0xcdba
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "jump"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInteger, Value: "0xcdba"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Jump, 0xcd, 0xba}},
	{input: []lexer_rewrite.Token{
		// jump 0xba
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "jumpz"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInteger, Value: "0xba"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Jumpz, 0x00, 0xba}},
	{input: []lexer_rewrite.Token{
		// jump 0xcdba
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "jumpz"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInteger, Value: "0xcdba"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Jumpz, 0xcd, 0xba}},
	{input: []lexer_rewrite.Token{
		// jump 0xba
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "jumpnz"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInteger, Value: "0xba"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Jumpnz, 0x00, 0xba}},
	{input: []lexer_rewrite.Token{
		// jump 0xcdba
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: "jumpnz"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInteger, Value: "0xcdba"},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeNewLine},
		lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeEndOfFile},
	}, output: []byte{instructions.Jumpnz, 0xcd, 0xba}},
}

func TestParser(t *testing.T) {
	var parserTestCases = []parserTestCase{}
	// arithmetic and logic
	// no operand = implied B register as value
	// instruction, (immediate mode|register mode), (immediate value|register number)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTests("add", instructions.Add)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTests("sub", instructions.Sub)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTests("mul", instructions.Mul)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTests("div", instructions.Div)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTests("or", instructions.Or)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTests("shl", instructions.Shl)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTests("shr", instructions.Shr)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTests("rand", instructions.Rand)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTests("gt", instructions.GT)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTests("gte", instructions.GTE)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTests("lt", instructions.LT)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTests("lte", instructions.LTE)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTests("eq", instructions.Eq)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTests("neq", instructions.Neq)...)

	parserTestCases = append(parserTestCases, singleOperandInstructionTestCases...)
	parserTestCases = append(parserTestCases, movTestCases...)
	parserTestCases = append(parserTestCases, callTestCases...)
	parserTestCases = append(parserTestCases, jumpTestCases...)
	parserTestCases = append(parserTestCases, loadTestCases...)

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
