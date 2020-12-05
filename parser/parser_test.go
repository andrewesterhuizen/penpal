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

// TODO: move these token constructers to lexer once packages are merged to assembler
func token(t lexer_rewrite.TokenType, v string) lexer_rewrite.Token {
	return lexer_rewrite.Token{Type: t, Value: v}
}

func tokens(t ...lexer_rewrite.Token) []lexer_rewrite.Token {
	return t
}

func tokenText(v string) lexer_rewrite.Token {
	return lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeText, Value: v}
}

func tokenInt(v string) lexer_rewrite.Token {
	return lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInteger, Value: v}
}

func tokenNL() lexer_rewrite.Token {
	return token(lexer_rewrite.TokenTypeNewLine, "\n")
}

func tokenEOF() lexer_rewrite.Token {
	return token(lexer_rewrite.TokenTypeEndOfFile, "")
}

func tokenComma() lexer_rewrite.Token {
	return token(lexer_rewrite.TokenTypeComma, ",")
}

func tokenLParen() lexer_rewrite.Token {
	return token(lexer_rewrite.TokenTypeLeftParen, "(")
}

func tokenRParen() lexer_rewrite.Token {
	return token(lexer_rewrite.TokenTypeLeftParen, "(")
}

func newNoOperandInstructionTest(text string, ins byte) parserTestCase {
	return parserTestCase{
		tokens(tokenText(text), tokenNL(), tokenEOF()),
		[]byte{ins},
	}
}

func newAritmeticLogicInstructionTest(text string, ins byte) []parserTestCase {
	return []parserTestCase{
		// ins
		{
			tokens(tokenText(text), tokenNL(), tokenEOF()),
			[]byte{ins, instructions.Register, instructions.RegisterB},
		},
		// ins 0xbb
		{
			input:  tokens(tokenText(text), tokenInt("0xbb"), tokenNL(), tokenEOF()),
			output: []byte{ins, instructions.Immediate, 0xbb},
		},
	}
}

var singleOperandInstructionTestCases = []parserTestCase{
	newNoOperandInstructionTest("swap", instructions.Swap),
	newNoOperandInstructionTest("pop", instructions.Pop),
	newNoOperandInstructionTest("ret", instructions.Ret),
	newNoOperandInstructionTest("reti", instructions.Reti),
	newNoOperandInstructionTest("halt", instructions.Halt),
}

var movTestCases = []parserTestCase{
	// mov 0xab, A
	{
		input:  tokens(tokenText("mov"), tokenInt("0xab"), tokenComma(), tokenText("A"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Mov, instructions.AddressingModeImmediate, 0xab, instructions.RegisterA},
	},
	// mov fp, A
	{
		input:  tokens(tokenText("mov"), tokenText("fp"), tokenComma(), tokenText("A"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Mov, instructions.FramePointerWithOffset, 0x0, instructions.RegisterA},
	},
	// mov (fp+1), A
	{
		input: tokens(
			tokenText("mov"),
			tokenLParen(),
			tokenText("fp"),
			token(lexer_rewrite.TokenTypePlus, "+"),
			tokenInt("1"),
			tokenRParen(),
			tokenComma(),
			tokenText("A"),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Mov, instructions.FramePointerWithOffset, 0x1, instructions.RegisterA},
	},
	// mov (fp-1), B
	{
		input: tokens(
			tokenText("mov"),
			tokenLParen(),
			tokenText("fp"),
			token(lexer_rewrite.TokenTypeMinus, "-"),
			tokenInt("1"),
			tokenRParen(),
			tokenComma(),
			tokenText("B"),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Mov, instructions.FramePointerWithOffset, 0xff, instructions.RegisterB},
	},
	// mov (fp[3]), B
	{
		input: tokens(
			tokenText("mov"),
			tokenLParen(),
			tokenText("fp"),
			token(lexer_rewrite.TokenTypeLeftBracket, "["),
			tokenInt("3"),
			token(lexer_rewrite.TokenTypeRightBracket, "]"),
			tokenRParen(),
			tokenComma(),
			tokenText("B"),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Mov, instructions.FramePointerWithOffset, 0x3, instructions.RegisterB},
	},
	// mov (fp[+3]), A
	{
		input: tokens(tokenText("mov"),
			tokenLParen(),
			tokenText("fp"),
			token(lexer_rewrite.TokenTypeLeftBracket, "["),
			token(lexer_rewrite.TokenTypePlus, "+"),
			tokenInt("3"),
			token(lexer_rewrite.TokenTypeRightBracket, "]"),
			tokenRParen(),
			tokenComma(),
			tokenText("A"),
			tokenNL(),
			tokenEOF()),
		output: []byte{instructions.Mov, instructions.FramePointerWithOffset, 0x3, instructions.RegisterA},
	},
	// mov (fp[-3]), A
	{
		input: tokens(tokenText("mov"),
			tokenLParen(),
			tokenText("fp"),
			token(lexer_rewrite.TokenTypeLeftBracket, "["),
			token(lexer_rewrite.TokenTypeMinus, "-"),
			tokenInt("3"),
			token(lexer_rewrite.TokenTypeRightBracket, "]"),
			tokenRParen(), tokenComma(),
			tokenText("A"),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Mov, instructions.FramePointerWithOffset, 0xfd, instructions.RegisterA},
	},
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
	// load 0xae, A
	{
		input:  tokens(tokenText("load"), tokenInt("0xae"), tokenComma(), tokenText("A"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Load, 0x00, 0xae, instructions.Immediate, 0x0, instructions.RegisterA},
	},
	// load 0xaecd, B
	{
		input: tokens(
			tokenText("load"),
			tokenInt("0xaecd"),
			tokenComma(),
			tokenText("B"),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Load, 0xae, 0xcd, instructions.Immediate, 0x0, instructions.RegisterB},
	},
	// load (fp + 5), A
	{
		input: tokens(
			tokenText("load"),
			tokenLParen(),
			tokenText("fp"),
			token(lexer_rewrite.TokenTypePlus, "+"),
			tokenInt("5"),
			tokenRParen(),
			tokenComma(),
			tokenText("A"),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Load, 0x00, 0x00, instructions.FramePointerWithOffset, 0x5, instructions.RegisterA},
	},
	// load (fp + 5), A
	{
		input: tokens(
			tokenText("load"),
			tokenLParen(),
			tokenText("fp"),
			token(lexer_rewrite.TokenTypePlus, "+"),
			tokenInt("5"),
			tokenRParen(),
			tokenComma(),
			tokenText("A"),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Load, 0x00, 0x00, instructions.FramePointerWithOffset, 0x5, instructions.RegisterA},
	},
	// load (fp[5]), A
	{
		input: tokens(
			tokenText("load"),
			tokenLParen(),
			tokenText("fp"),
			token(lexer_rewrite.TokenTypeLeftBracket, "["),
			tokenInt("5"),
			token(lexer_rewrite.TokenTypeRightBracket, "]"),
			tokenRParen(),
			tokenComma(),
			tokenText("A"),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Load, 0x00, 0x00, instructions.FramePointerWithOffset, 0x5, instructions.RegisterA},
	},
	// load (fp[-1]), B
	{
		input: tokens(
			tokenText("load"),
			tokenLParen(),
			tokenText("fp"),
			token(lexer_rewrite.TokenTypeLeftBracket, "["),
			token(lexer_rewrite.TokenTypeMinus, "-"),
			tokenInt("1"),
			token(lexer_rewrite.TokenTypeRightBracket, "]"),
			tokenRParen(),
			tokenComma(),
			tokenText("B"),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Load, 0x00, 0x00, instructions.FramePointerWithOffset, 0xff, instructions.RegisterB},
	},
	// load (fp - A), A
	{
		input: tokens(
			tokenText("load"),
			tokenLParen(),
			tokenText("fp"),
			token(lexer_rewrite.TokenTypeMinus, "-"),
			tokenText("A"),
			tokenRParen(),
			tokenComma(),
			tokenText("A"),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Load, 0x00, 0x00, instructions.FramePointerMinusRegister, instructions.RegisterA, instructions.RegisterA},
	},
	//  load (fp - 1), A
	{
		input: tokens(
			tokenText("load"),
			tokenLParen(),
			tokenText("fp"),
			token(lexer_rewrite.TokenTypeMinus, "-"),
			tokenInt("1"),
			tokenRParen(),
			tokenComma(),
			tokenText("A"),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Load, 0x00, 0x00, instructions.FramePointerWithOffset, 0xff, instructions.RegisterA},
	},
	//  load fp, A
	{
		input: tokens(
			tokenText("load"),
			tokenText("fp"),
			tokenComma(),
			tokenText("A"),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Load, 0x00, 0x00, instructions.FramePointerWithOffset, 0x0, instructions.RegisterA},
	},
}

// TODO: labels
var callTestCases = []parserTestCase{
	// call 0xba
	{
		input:  tokens(tokenText("call"), tokenInt("0xba"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Call, 0x00, 0xba},
	},
	// call 0xcdba
	{
		input:  tokens(tokenText("call"), tokenInt("0xcdba"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Call, 0xcd, 0xba},
	},
}

var jumpTestCases = []parserTestCase{
	// jump 0xba
	{
		input:  tokens(tokenText("jump"), tokenInt("0xba"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Jump, 0x00, 0xba},
	},
	// jump 0xcdba
	{
		input:  tokens(tokenText("jump"), tokenInt("0xcdba"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Jump, 0xcd, 0xba},
	},
	// jump 0xba
	{
		input:  tokens(tokenText("jumpz"), tokenInt("0xba"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Jumpz, 0x00, 0xba},
	},
	// jump 0xcdba
	{
		input:  tokens(tokenText("jumpz"), tokenInt("0xcdba"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Jumpz, 0xcd, 0xba},
	},
	// jump 0xba
	{
		input:  tokens(tokenText("jumpnz"), tokenInt("0xba"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Jumpnz, 0x00, 0xba},
	},
	// jump 0xcdba
	{
		input:  tokens(tokenText("jumpnz"), tokenInt("0xcdba"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Jumpnz, 0xcd, 0xba},
	},
}

func TestParser(t *testing.T) {
	var parserTestCases = []parserTestCase{}
	// arithmetic and logic
	// no operand = implied B register as value
	// instruction, (immediate mode|register mode), (immediate value|register number)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTest("add", instructions.Add)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTest("sub", instructions.Sub)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTest("mul", instructions.Mul)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTest("div", instructions.Div)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTest("or", instructions.Or)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTest("shl", instructions.Shl)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTest("shr", instructions.Shr)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTest("rand", instructions.Rand)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTest("gt", instructions.GT)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTest("gte", instructions.GTE)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTest("lt", instructions.LT)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTest("lte", instructions.LTE)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTest("eq", instructions.Eq)...)
	parserTestCases = append(parserTestCases, newAritmeticLogicInstructionTest("neq", instructions.Neq)...)

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
