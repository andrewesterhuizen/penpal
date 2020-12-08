package assembler2

import (
	"fmt"
	"testing"

	"github.com/andrewesterhuizen/penpal/instructions"
	"github.com/andrewesterhuizen/penpal/lexer_rewrite"
)

type parserTestCase struct {
	input  string
	output []byte
}

func newAritmeticLogicInstructionTest(text string, ins byte) []parserTestCase {
	return []parserTestCase{
		// ins
		{
			input:  text,
			output: []byte{ins, instructions.Register, instructions.RegisterB},
		},
		// ins 0xbb
		{
			input:  fmt.Sprintf("%s 0xbb", text),
			output: []byte{ins, instructions.Immediate, 0xbb},
		},
	}
}

var movTestCases = []parserTestCase{
	{
		input:  "mov A, 1",
		output: []byte{instructions.Mov, instructions.RegisterA, 1},
	},
	{
		input:  "mov A, 0xab",
		output: []byte{instructions.Mov, instructions.RegisterA, 0xab},
	},
	{
		input:  "mov B, 0xcd",
		output: []byte{instructions.Mov, instructions.RegisterB, 0xcd},
	},
}

var loadTestCases = []parserTestCase{
	{
		input:  "load 1, A",
		output: []byte{instructions.Load, 0x00, 0x01, instructions.Immediate, 0x0, instructions.RegisterA},
	},
	{
		input:  "load 0xae, A",
		output: []byte{instructions.Load, 0x00, 0xae, instructions.Immediate, 0x0, instructions.RegisterA},
	},
	{
		input:  "load 0xaecd, B",
		output: []byte{instructions.Load, 0xae, 0xcd, instructions.Immediate, 0x0, instructions.RegisterB},
	},
	{
		input:  "load (fp + 5), A",
		output: []byte{instructions.Load, 0x00, 0x00, instructions.FramePointerWithOffset, 0x5, instructions.RegisterA},
	},
	{
		input:  "load (fp[5]), A",
		output: []byte{instructions.Load, 0x00, 0x00, instructions.FramePointerWithOffset, 0x5, instructions.RegisterA},
	},
	{
		input:  "load (fp - A), A",
		output: []byte{instructions.Load, 0x00, 0x00, instructions.FramePointerMinusRegister, instructions.RegisterA, instructions.RegisterA},
	},
	{
		input:  "load fp, A",
		output: []byte{instructions.Load, 0x00, 0x00, instructions.FramePointerWithOffset, 0x0, instructions.RegisterA},
	},
	{
		input: `
		test_label: db 1

		load (test_label + 6), A`,
		output: []byte{1, instructions.Load, 0x00, 0x00, instructions.Immediate, 6, instructions.RegisterA},
	},
	{
		input: `
		test_label: db 1

		load (test_label - 1), A`,
		output: []byte{1, instructions.Load, 0x00, 0x00, instructions.Immediate, 0xff, instructions.RegisterA},
	},
	{
		input: `
		test_label: db 1

		load (test_label[3]), A`,
		output: []byte{1, instructions.Load, 0x00, 0x00, instructions.Immediate, 3, instructions.RegisterA},
	},
	{
		input: `
		test_label: db 1

		load (test_label + A), B`,
		output: []byte{1, instructions.Load, 0x00, 0x00, instructions.ImmediatePlusRegister, instructions.RegisterA, instructions.RegisterB},
	},
	{
		input: `
		test_label: db 1

		load (test_label - A), B`,
		output: []byte{1, instructions.Load, 0x00, 0x00, instructions.ImmediateMinusRegister, instructions.RegisterA, instructions.RegisterB},
	},
	{
		input: `
		test_label: db 1

		load (test_label[A]), B`,
		output: []byte{1, instructions.Load, 0x00, 0x00, instructions.ImmediatePlusRegister, instructions.RegisterA, instructions.RegisterB},
	},
	{
		input: `
		test_label: db 1

		load test_label, B`,
		output: []byte{1, instructions.Load, 0x00, 0x00, instructions.Immediate, 0x0, instructions.RegisterB},
	},
}

var storeTestCases = []parserTestCase{
	{
		input:  "store A, 0xabcd",
		output: []byte{instructions.Store, instructions.RegisterA, instructions.Immediate, 0x0, 0xab, 0xcd},
	},
	{
		input:  "store B, fp",
		output: []byte{instructions.Store, instructions.RegisterB, instructions.FramePointerWithOffset, 0x0, 0x0, 0x0},
	},
	{
		input:  "store A, (fp + B)",
		output: []byte{instructions.Store, instructions.RegisterA, instructions.FramePointerPlusRegister, instructions.RegisterB, 0x0, 0x0},
	},
	{
		input:  "store B, (fp - A)",
		output: []byte{instructions.Store, instructions.RegisterB, instructions.FramePointerMinusRegister, instructions.RegisterA, 0x0, 0x0},
	},
	{
		input:  "store B, (fp + 3)",
		output: []byte{instructions.Store, instructions.RegisterB, instructions.FramePointerWithOffset, 0x3, 0x0, 0x0},
	},
	{
		input:  "store A, (fp - 1)",
		output: []byte{instructions.Store, instructions.RegisterA, instructions.FramePointerWithOffset, 0xff, 0x0, 0x0},
	},
	{
		input:  "store A, (fp[3])",
		output: []byte{instructions.Store, instructions.RegisterA, instructions.FramePointerWithOffset, 0x3, 0x0, 0x0},
	},
	{
		input: `
		test_label: db 1

		store B, test_label`,
		output: []byte{1, instructions.Store, instructions.RegisterB, instructions.Immediate, 0x0, 0x0, 0x0},
	},
	{
		input: `
		test_label: db 1

		store B, (test_label + 3)`,
		output: []byte{1, instructions.Store, instructions.RegisterB, instructions.Immediate, 0x3, 0x0, 0x0},
	},
	{
		input: `
		test_label: db 1

		store B, (test_label - 2)`,
		output: []byte{1, instructions.Store, instructions.RegisterB, instructions.Immediate, 0xfe, 0x0, 0x0},
	},
	{
		input: `
		test_label: db 1

		store B, (test_label[5])`,
		output: []byte{1, instructions.Store, instructions.RegisterB, instructions.Immediate, 0x5, 0x0, 0x0},
	},
	{
		input: `
		test_label: db 1

		store B, (test_label[B])`,
		output: []byte{1, instructions.Store, instructions.RegisterB, instructions.ImmediatePlusRegister, instructions.RegisterB, 0x0, 0x0},
	},
}

var callTestCases = []parserTestCase{
	{
		input:  "call 0xba",
		output: []byte{instructions.Call, 0x00, 0xba},
	},
	{
		input:  "call 0xcdba",
		output: []byte{instructions.Call, 0xcd, 0xba},
	},
}

var jumpTestCases = []parserTestCase{
	{
		input:  "jump 0xba",
		output: []byte{instructions.Jump, 0x00, 0xba},
	},
	{
		input:  "jump 0xcdba",
		output: []byte{instructions.Jump, 0xcd, 0xba},
	},
	{
		input: `test_label:
		add

		jump test_label
		`,
		output: []byte{instructions.Add, instructions.Jump, 0, 0},
	},
	{
		input: `jump test_label
		
		test_label: add
		`,
		output: []byte{instructions.Jump, 0, 3, instructions.Add},
	},
	{
		input:  "jumpz 0xba",
		output: []byte{instructions.Jumpz, 0x00, 0xba},
	},
	{
		input:  "jumpz 0xcdba",
		output: []byte{instructions.Jumpz, 0xcd, 0xba},
	},
	{
		input:  "jumpnz 0xba",
		output: []byte{instructions.Jumpnz, 0x00, 0xba},
	},
	{
		input:  "jumpnz 0xcdba",
		output: []byte{instructions.Jumpnz, 0xcd, 0xba},
	},
}

var labelTestCases = []parserTestCase{
	{
		input: "label:", output: []byte{},
	},
	{
		input:  "label: db 55",
		output: []byte{55},
	},
}

var dbTestCases = []parserTestCase{
	{
		input: "db 33", output: []byte{33},
	},
	{
		input: `db 33
		db 44
		db 22`,
		output: []byte{33, 44, 22},
	},
}

var testCases = []parserTestCase{
	// check parser ignores whitespace
	{
		input: `

		db 1

		db 2


		db 3



		`,
		output: []byte{1, 2, 3},
	},
}

var pushTestCases = []parserTestCase{
	{
		input:  "push",
		output: []byte{instructions.Push, instructions.Register, instructions.RegisterA},
	},
	{
		input:  "push 0xae",
		output: []byte{instructions.Push, instructions.Immediate, 0xae},
	},
}

func TestParser(t *testing.T) {
	var parserTestCases = []parserTestCase{}

	parserTestCases = append(parserTestCases, parserTestCase{"swap", []byte{instructions.Swap}})
	parserTestCases = append(parserTestCases, parserTestCase{"pop", []byte{instructions.Pop}})
	parserTestCases = append(parserTestCases, parserTestCase{"ret", []byte{instructions.Ret}})
	parserTestCases = append(parserTestCases, parserTestCase{"reti", []byte{instructions.Reti}})
	parserTestCases = append(parserTestCases, parserTestCase{"halt", []byte{instructions.Halt}})
	parserTestCases = append(parserTestCases, parserTestCase{"add", []byte{instructions.Add}})
	parserTestCases = append(parserTestCases, parserTestCase{"sub", []byte{instructions.Sub}})
	parserTestCases = append(parserTestCases, parserTestCase{"mul", []byte{instructions.Mul}})
	parserTestCases = append(parserTestCases, parserTestCase{"div", []byte{instructions.Div}})
	parserTestCases = append(parserTestCases, parserTestCase{"or", []byte{instructions.Or}})
	parserTestCases = append(parserTestCases, parserTestCase{"shl", []byte{instructions.Shl}})
	parserTestCases = append(parserTestCases, parserTestCase{"shr", []byte{instructions.Shr}})
	parserTestCases = append(parserTestCases, parserTestCase{"rand", []byte{instructions.Rand}})
	parserTestCases = append(parserTestCases, parserTestCase{"gt", []byte{instructions.GT}})
	parserTestCases = append(parserTestCases, parserTestCase{"gte", []byte{instructions.GTE}})
	parserTestCases = append(parserTestCases, parserTestCase{"lt", []byte{instructions.LT}})
	parserTestCases = append(parserTestCases, parserTestCase{"lte", []byte{instructions.LTE}})
	parserTestCases = append(parserTestCases, parserTestCase{"eq", []byte{instructions.Eq}})
	parserTestCases = append(parserTestCases, parserTestCase{"neq", []byte{instructions.Neq}})

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
		l := lexer_rewrite.NewLexer()
		tokens, err := l.Run("", tc.input)
		if err != nil {
			t.Errorf("%s\nwith input:\n %v", err, tc.input)
			return
		}

		p := NewParser()

		out, err := p.Run(tokens)
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
				t.Errorf("expected byte 0x%02x and got 0x%02x in %d position\nwith input:\n %v", ins, out[i], i, tc.input)
			}

		}
	}
}
