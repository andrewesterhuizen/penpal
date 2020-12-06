package assembler2

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

func tokenInstruction(v string) lexer_rewrite.Token {
	return lexer_rewrite.Token{Type: lexer_rewrite.TokenTypeInstruction, Value: v}
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
	return token(lexer_rewrite.TokenTypeRightParen, "(")
}

func newNoOperandInstructionTest(text string, ins byte) parserTestCase {
	return parserTestCase{
		input:  tokens(tokenInstruction(text), tokenNL(), tokenEOF()),
		output: []byte{ins},
	}
}

func newAritmeticLogicInstructionTest(text string, ins byte) []parserTestCase {
	return []parserTestCase{
		// ins
		{
			input:  tokens(tokenInstruction(text), tokenNL(), tokenEOF()),
			output: []byte{ins, instructions.Register, instructions.RegisterB},
		},
		// ins 0xbb
		{
			input:  tokens(tokenInstruction(text), tokenInt("0xbb"), tokenNL(), tokenEOF()),
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
	// mov A, 0xab
	{
		input:  tokens(tokenInstruction("mov"), tokenText("A"), tokenComma(), tokenInt("0xab"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Mov, instructions.RegisterA, 0xab},
	},
	// mov B, 0xcd
	{
		input:  tokens(tokenInstruction("mov"), tokenText("B"), tokenComma(), tokenInt("0xcd"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Mov, instructions.RegisterB, 0xcd},
	},
}

var loadTestCases = []parserTestCase{
	// load 0xae, A
	{
		input:  tokens(tokenInstruction("load"), tokenInt("0xae"), tokenComma(), tokenText("A"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Load, 0x00, 0xae, instructions.Immediate, 0x0, instructions.RegisterA},
	},
	// load 0xaecd, B
	{
		input: tokens(
			tokenInstruction("load"),
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
			tokenInstruction("load"),
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
			tokenInstruction("load"),
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
			tokenInstruction("load"),
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
			tokenInstruction("load"),
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
			tokenInstruction("load"),
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
			tokenInstruction("load"),
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
			tokenInstruction("load"),
			tokenText("fp"),
			tokenComma(),
			tokenText("A"),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Load, 0x00, 0x00, instructions.FramePointerWithOffset, 0x0, instructions.RegisterA},
	},
}

var storeTestCases = []parserTestCase{
	// store A, 0xabcd
	{
		input:  tokens(tokenInstruction("store"), tokenText("A"), tokenComma(), tokenInt("0xabcd"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Store, instructions.RegisterA, instructions.Immediate, 0x0, 0xab, 0xcd},
	},
	// store B, fp
	{
		input:  tokens(tokenInstruction("store"), tokenText("B"), tokenComma(), tokenText("fp"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Store, instructions.RegisterB, instructions.FramePointerWithOffset, 0x0, 0x0, 0x0},
	},
	// store A, (fp + B)
	{
		input: tokens(
			tokenInstruction("store"),
			tokenText("A"),
			tokenComma(),
			tokenLParen(),
			tokenText("fp"),
			token(lexer_rewrite.TokenTypePlus, "+"),
			tokenText("B"),
			tokenRParen(),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Store, instructions.RegisterA, instructions.FramePointerPlusRegister, instructions.RegisterB, 0x0, 0x0},
	},
	// store B, (fp - A)
	{
		input: tokens(
			tokenInstruction("store"),
			tokenText("B"),
			tokenComma(),
			tokenLParen(),
			tokenText("fp"),
			token(lexer_rewrite.TokenTypeMinus, "-"),
			tokenText("A"),
			tokenRParen(),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Store, instructions.RegisterB, instructions.FramePointerMinusRegister, instructions.RegisterA, 0x0, 0x0},
	},
	// store B, (fp + 3)
	{
		input: tokens(
			tokenInstruction("store"),
			tokenText("B"),
			tokenComma(),
			tokenLParen(),
			tokenText("fp"),
			token(lexer_rewrite.TokenTypePlus, "+"),
			tokenInt("3"),
			tokenRParen(),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Store, instructions.RegisterB, instructions.FramePointerWithOffset, 0x3, 0x0, 0x0},
	},
	// store A, (fp - 1)
	{
		input: tokens(
			tokenInstruction("store"),
			tokenText("A"),
			tokenComma(),
			tokenLParen(),
			tokenText("fp"),
			token(lexer_rewrite.TokenTypeMinus, "-"),
			tokenInt("1"),
			tokenRParen(),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Store, instructions.RegisterA, instructions.FramePointerWithOffset, 0xff, 0x0, 0x0},
	},
	// store A, (fp[3])
	{
		input: tokens(
			tokenInstruction("store"),
			tokenText("A"),
			tokenComma(),
			tokenLParen(),
			tokenText("fp"),
			token(lexer_rewrite.TokenTypeLeftBracket, "["),
			tokenInt("3"),
			token(lexer_rewrite.TokenTypeRightBracket, "]"),
			tokenRParen(),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Store, instructions.RegisterA, instructions.FramePointerWithOffset, 0x3, 0x0, 0x0},
	},
	// store A, (fp[+3])
	{
		input: tokens(
			tokenInstruction("store"),
			tokenText("A"),
			tokenComma(),
			tokenLParen(),
			tokenText("fp"),
			token(lexer_rewrite.TokenTypeLeftBracket, "["),
			token(lexer_rewrite.TokenTypePlus, "+"),
			tokenInt("3"),
			token(lexer_rewrite.TokenTypeRightBracket, "]"),
			tokenRParen(),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Store, instructions.RegisterA, instructions.FramePointerWithOffset, 0x3, 0x0, 0x0},
	},
	// store A, (fp[-1])
	{
		input: tokens(
			tokenInstruction("store"),
			tokenText("A"),
			tokenComma(),
			tokenLParen(),
			tokenText("fp"),
			token(lexer_rewrite.TokenTypeLeftBracket, "["),
			token(lexer_rewrite.TokenTypeMinus, "-"),
			tokenInt("1"),
			token(lexer_rewrite.TokenTypeRightBracket, "]"),
			tokenRParen(),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Store, instructions.RegisterA, instructions.FramePointerWithOffset, 0xff, 0x0, 0x0},
	},
}

var callTestCases = []parserTestCase{
	// call 0xba
	{
		input:  tokens(tokenInstruction("call"), tokenInt("0xba"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Call, 0x00, 0xba},
	},
	// call 0xcdba
	{
		input:  tokens(tokenInstruction("call"), tokenInt("0xcdba"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Call, 0xcd, 0xba},
	},
}

var jumpTestCases = []parserTestCase{
	// jump 0xba
	{
		input:  tokens(tokenInstruction("jump"), tokenInt("0xba"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Jump, 0x00, 0xba},
	},
	// jump 0xcdba
	{
		input:  tokens(tokenInstruction("jump"), tokenInt("0xcdba"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Jump, 0xcd, 0xba},
	},
	// jump test_label
	{
		input: tokens(
			token(lexer_rewrite.TokenTypeLabel, "test_label"),
			tokenInstruction("add"),
			tokenNL(),
			tokenInstruction("jump"),
			tokenText("test_label"),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Add, instructions.Jump, 0, 0},
	},
	// jump test_label, reference before definition
	{
		input: tokens(
			tokenInstruction("jump"),
			tokenText("test_label"),
			tokenNL(),
			// label
			token(lexer_rewrite.TokenTypeLabel, "test_label"),
			tokenInstruction("add"),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Jump, 0, 3, instructions.Add},
	},
	// jump test_label, label def has newline
	{
		input: tokens(
			tokenInstruction("jump"),
			tokenText("test_label"),
			tokenNL(),
			// label
			token(lexer_rewrite.TokenTypeLabel, "test_label"),
			tokenNL(),
			tokenInstruction("add"),
			tokenNL(),
			tokenEOF(),
		),
		output: []byte{instructions.Jump, 0, 3, instructions.Add},
	},
	// jumpz 0xba
	{
		input:  tokens(tokenInstruction("jumpz"), tokenInt("0xba"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Jumpz, 0x00, 0xba},
	},
	// jumpz 0xcdba
	{
		input:  tokens(tokenInstruction("jumpz"), tokenInt("0xcdba"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Jumpz, 0xcd, 0xba},
	},
	// jumpnz 0xba
	{
		input:  tokens(tokenInstruction("jumpnz"), tokenInt("0xba"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Jumpnz, 0x00, 0xba},
	},
	// jumpnz 0xcdba
	{
		input:  tokens(tokenInstruction("jumpnz"), tokenInt("0xcdba"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Jumpnz, 0xcd, 0xba},
	},
}

var labelTestCases = []parserTestCase{
	{input: tokens(token(lexer_rewrite.TokenTypeLabel, "label"), tokenNL()), output: []byte{}},
	{
		input:  tokens(token(lexer_rewrite.TokenTypeLabel, "label"), tokenNL(), tokenInstruction("db"), tokenInt("55"), tokenNL()),
		output: []byte{55},
	},
}

var dbTestCases = []parserTestCase{
	{
		input: tokens(tokenInstruction("db"), tokenInt("33"), tokenNL()), output: []byte{33},
	},
	{
		input: tokens(
			tokenInstruction("db"),
			tokenInt("33"),
			tokenNL(),
			tokenInstruction("db"),
			tokenInt("44"),
			tokenNL(),
			tokenInstruction("db"),
			tokenInt("22"),
			tokenNL(),
		), output: []byte{33, 44, 22},
	},
}

var testCases = []parserTestCase{
	{
		// check parser ignores whitespace
		input: tokens(
			tokenNL(),
			tokenInstruction("db"),
			tokenInt("1"),
			tokenNL(),
			tokenInstruction("db"),
			tokenInt("2"),
			tokenNL(),
			tokenNL(),
			tokenInstruction("db"),
			tokenInt("3"),
			tokenNL(),
			tokenNL(),
			tokenNL(),
			tokenNL(),
			tokenEOF(),
		), output: []byte{1, 2, 3},
	},
}

var pushTestCases = []parserTestCase{
	{
		input:  tokens(tokenInstruction("push"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Push, instructions.Register, instructions.RegisterA},
	},
	{
		input:  tokens(tokenInstruction("push"), tokenInt("0xae"), tokenNL(), tokenEOF()),
		output: []byte{instructions.Push, instructions.Immediate, 0xae},
	},
}

func TestParser(t *testing.T) {
	var parserTestCases = []parserTestCase{}
	// arithmetic and logic
	// no operand = implied B register as value
	// instruction, (immediate mode|register mode), (immediate value|register number)
	parserTestCases = append(parserTestCases, newNoOperandInstructionTest("add", instructions.Add))
	parserTestCases = append(parserTestCases, newNoOperandInstructionTest("sub", instructions.Sub))
	parserTestCases = append(parserTestCases, newNoOperandInstructionTest("mul", instructions.Mul))
	parserTestCases = append(parserTestCases, newNoOperandInstructionTest("div", instructions.Div))
	parserTestCases = append(parserTestCases, newNoOperandInstructionTest("or", instructions.Or))
	parserTestCases = append(parserTestCases, newNoOperandInstructionTest("shl", instructions.Shl))
	parserTestCases = append(parserTestCases, newNoOperandInstructionTest("shr", instructions.Shr))
	parserTestCases = append(parserTestCases, newNoOperandInstructionTest("rand", instructions.Rand))
	parserTestCases = append(parserTestCases, newNoOperandInstructionTest("gt", instructions.GT))
	parserTestCases = append(parserTestCases, newNoOperandInstructionTest("gte", instructions.GTE))
	parserTestCases = append(parserTestCases, newNoOperandInstructionTest("lt", instructions.LT))
	parserTestCases = append(parserTestCases, newNoOperandInstructionTest("lte", instructions.LTE))
	parserTestCases = append(parserTestCases, newNoOperandInstructionTest("eq", instructions.Eq))
	parserTestCases = append(parserTestCases, newNoOperandInstructionTest("neq", instructions.Neq))

	parserTestCases = append(parserTestCases, singleOperandInstructionTestCases...)
	parserTestCases = append(parserTestCases, movTestCases...)
	parserTestCases = append(parserTestCases, callTestCases...)
	parserTestCases = append(parserTestCases, jumpTestCases...)
	parserTestCases = append(parserTestCases, loadTestCases...)
	parserTestCases = append(parserTestCases, storeTestCases...)
	parserTestCases = append(parserTestCases, labelTestCases...)
	parserTestCases = append(parserTestCases, dbTestCases...)
	parserTestCases = append(parserTestCases, testCases...)
	parserTestCases = append(parserTestCases, pushTestCases...)

	for _, tc := range parserTestCases {
		p := NewParser()

		out, err := p.Run(tc.input)
		if err != nil {
			t.Errorf("%s\nwith input:\n %v", err, tc.input)
			return
		}

		fmt.Println(out)

		if len(tc.output) != len(out) {
			t.Errorf("expected %v bytes and got %v, with input:\n%s", len(tc.output), len(out), tc.input)
			return
		}

		for i, ins := range tc.output {
			if ins != out[i] {
				t.Errorf("expected byte 0x%02x and got 0x%02x", ins, out[i])
			}

		}
	}
}
