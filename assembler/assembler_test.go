package assembler

import (
	"testing"

	"github.com/andrewesterhuizen/vm/instructions"
)

type TestCase struct {
	input  string
	output []uint8
}

var instructionTestCases = []TestCase{
	TestCase{input: "MOVA 0xab", output: []uint8{instructions.MOVA, 0xab}},
	TestCase{input: "MOVB 0xcd", output: []uint8{instructions.MOVB, 0xcd}},
	TestCase{input: "SWAP", output: []uint8{instructions.SWAP}},
	TestCase{input: "LOAD 0xae", output: []uint8{instructions.LOAD, 0x00, 0xae}},
	TestCase{input: "LOAD 0xaecd", output: []uint8{instructions.LOAD, 0xae, 0xcd}},
	TestCase{input: "STORE 0xae", output: []uint8{instructions.STORE, 0x00, 0xae}},
	TestCase{input: "STORE 0xaecd", output: []uint8{instructions.STORE, 0xae, 0xcd}},
	TestCase{input: "POP", output: []uint8{instructions.POP}},
	TestCase{input: "PUSH", output: []uint8{instructions.PUSH}},
	TestCase{input: "CALL 0xbeaf", output: []uint8{instructions.CALL, 0xbe, 0xaf}},
	TestCase{input: "RET", output: []uint8{instructions.RET}},
	TestCase{input: "ADD", output: []uint8{instructions.ADD}},
	TestCase{input: "SUB", output: []uint8{instructions.SUB}},
	TestCase{input: "MUL", output: []uint8{instructions.MUL}},
	TestCase{input: "DIV", output: []uint8{instructions.DIV}},
	TestCase{input: "SHL", output: []uint8{instructions.SHL}},
	TestCase{input: "SHR", output: []uint8{instructions.SHR}},
	TestCase{input: "AND", output: []uint8{instructions.AND}},
	TestCase{input: "OR", output: []uint8{instructions.OR}},
	TestCase{input: "HALT", output: []uint8{instructions.HALT}},
	TestCase{input: "JUMP 0xcd", output: []uint8{instructions.JUMP, 0x00, 0xcd}},
	TestCase{input: "JUMP 0xaecd", output: []uint8{instructions.JUMP, 0xae, 0xcd}},
	TestCase{input: "JUMPZ 0xcd", output: []uint8{instructions.JUMPZ, 0x00, 0xcd}},
	TestCase{input: "JUMPZ 0xaecd", output: []uint8{instructions.JUMPZ, 0xae, 0xcd}},
	TestCase{input: "JUMPNZ 0xcd", output: []uint8{instructions.JUMPNZ, 0x00, 0xcd}},
	TestCase{input: "JUMPNZ 0xaecd", output: []uint8{instructions.JUMPNZ, 0xae, 0xcd}},
	TestCase{input: "$test = 0xbc\nMOVA $test\n", output: []uint8{instructions.MOVA, 0xbc}},
	TestCase{input: "$test = 0xabcd\nMOVA $test\n", output: []uint8{instructions.MOVA, 0xcd}},
	TestCase{input: "$test = 0xbc\nMOVB $test\n", output: []uint8{instructions.MOVB, 0xbc}},
	TestCase{input: "$test = 0xabcd\nMOVB $test\n", output: []uint8{instructions.MOVB, 0xcd}},
	TestCase{input: "$test = 0xbc\nCALL $test\n", output: []uint8{instructions.CALL, 0x0, 0xbc}},
	TestCase{input: "$test = 0xabcd\nCALL $test\n", output: []uint8{instructions.CALL, 0xab, 0xcd}},
	TestCase{input: "$test = 0xbc\nJUMP $test\n", output: []uint8{instructions.JUMP, 0x0, 0xbc}},
	TestCase{input: "$test = 0xabcd\nJUMP $test\n", output: []uint8{instructions.JUMP, 0xab, 0xcd}},
	TestCase{input: "$test = 0xbc\nJUMPZ $test\n", output: []uint8{instructions.JUMPZ, 0x0, 0xbc}},
	TestCase{input: "$test = 0xabcd\nJUMPZ $test\n", output: []uint8{instructions.JUMPZ, 0xab, 0xcd}},
	TestCase{input: "$test = 0xbc\nJUMPNZ $test\n", output: []uint8{instructions.JUMPNZ, 0x0, 0xbc}},
	TestCase{input: "$test = 0xabcd\nJUMPNZ $test\n", output: []uint8{instructions.JUMPNZ, 0xab, 0xcd}},
	TestCase{input: "test:", output: []uint8{}},
	TestCase{input: `
	test:
		SWAP

	JUMP test
	`,
		output: []uint8{instructions.SWAP, instructions.JUMP, 0x0, 0x0},
	},
	TestCase{input: `
	JUMP test

	test:
		SWAP
	`,
		output: []uint8{instructions.JUMP, 0x0, 0x3, instructions.SWAP},
	},
}

func TestInstructions(t *testing.T) {
	for _, tc := range instructionTestCases {
		a := New()

		ins := a.GetInstructions(tc.input)

		if len(ins) != len(tc.output) {
			t.Errorf("expected %d instuctions and got %d", len(tc.output), len(ins))
			return
		}

		for i, in := range tc.output {
			if in != ins[i] {
				t.Errorf("expected 0x%02x and got 0x%02x", in, ins[i])
			}
		}

	}

}
