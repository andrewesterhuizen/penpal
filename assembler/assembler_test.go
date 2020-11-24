package assembler

import (
	"testing"

	"github.com/andrewesterhuizen/penpal/instructions"
)

type TestCase struct {
	input  string
	output []uint8
}

var instructionTestCases = []TestCase{
	{input: "MOV A 0xab", output: []uint8{instructions.MOV, instructions.EncodeFlags(instructions.AddressingModeImmediate, instructions.RegisterA), 0xab}},
	{input: "MOV A +1(fp)", output: []uint8{instructions.MOV, instructions.EncodeFlags(instructions.AddressingModeFPRelative, instructions.RegisterA), 0x1}},
	{input: "MOV A -1(fp)", output: []uint8{instructions.MOV, instructions.EncodeFlags(instructions.AddressingModeFPRelative, instructions.RegisterA), 0xff}},
	{input: "MOV B 0xcd", output: []uint8{instructions.MOV, instructions.EncodeFlags(instructions.AddressingModeImmediate, instructions.RegisterB), 0xcd}},
	{input: "SWAP", output: []uint8{instructions.SWAP}},
	{input: "LOAD 0xae", output: []uint8{instructions.LOAD, 0x00, 0xae}},
	{input: "LOAD 0xaecd", output: []uint8{instructions.LOAD, 0xae, 0xcd}},
	{input: "STORE 0xae", output: []uint8{instructions.STORE, 0x00, 0xae}},
	{input: "STORE 0xaecd", output: []uint8{instructions.STORE, 0xae, 0xcd}},
	{input: "POP", output: []uint8{instructions.POP}},
	{input: "PUSH", output: []uint8{instructions.PUSH}},
	{input: "CALL 0xbeaf", output: []uint8{instructions.CALL, 0xbe, 0xaf}},
	{input: "RET", output: []uint8{instructions.RET}},
	{input: "ADD", output: []uint8{instructions.ADD}},
	{input: "SUB", output: []uint8{instructions.SUB}},
	{input: "MUL", output: []uint8{instructions.MUL}},
	{input: "DIV", output: []uint8{instructions.DIV}},
	{input: "SHL", output: []uint8{instructions.SHL}},
	{input: "SHR", output: []uint8{instructions.SHR}},
	{input: "AND", output: []uint8{instructions.AND}},
	{input: "OR", output: []uint8{instructions.OR}},
	{input: "HALT", output: []uint8{instructions.HALT}},
	{input: "SEND", output: []uint8{instructions.SEND}},
	{input: "RAND", output: []uint8{instructions.RAND}},
	{input: "JUMP 0xcd", output: []uint8{instructions.JUMP, 0x00, 0xcd}},
	{input: "JUMP 0xaecd", output: []uint8{instructions.JUMP, 0xae, 0xcd}},
	{input: "JUMPZ 0xcd", output: []uint8{instructions.JUMPZ, 0x00, 0xcd}},
	{input: "JUMPZ 0xaecd", output: []uint8{instructions.JUMPZ, 0xae, 0xcd}},
	{input: "JUMPNZ 0xcd", output: []uint8{instructions.JUMPNZ, 0x00, 0xcd}},
	{input: "JUMPNZ 0xaecd", output: []uint8{instructions.JUMPNZ, 0xae, 0xcd}},
	{input: "$test = 0xbc\nMOV A $test\n", output: []uint8{instructions.MOV, instructions.EncodeFlags(instructions.AddressingModeImmediate, instructions.RegisterA), 0xbc}},
	{input: "$test = 0xabcd\nMOV A $test\n", output: []uint8{instructions.MOV, instructions.EncodeFlags(instructions.AddressingModeImmediate, instructions.RegisterA), 0xcd}},
	{input: "$test = 0xbc\nMOV A $test\n", output: []uint8{instructions.MOV, instructions.EncodeFlags(instructions.AddressingModeImmediate, instructions.RegisterA), 0xbc}},
	{input: "$test = 0xabcd\nMOV A $test\n", output: []uint8{instructions.MOV, instructions.EncodeFlags(instructions.AddressingModeImmediate, instructions.RegisterA), 0xcd}},
	{input: "$test = 0xbc\nCALL $test\n", output: []uint8{instructions.CALL, 0x0, 0xbc}},
	{input: "$test = 0xabcd\nCALL $test\n", output: []uint8{instructions.CALL, 0xab, 0xcd}},
	{input: "$test = 0xbc\nJUMP $test\n", output: []uint8{instructions.JUMP, 0x0, 0xbc}},
	{input: "$test = 0xabcd\nJUMP $test\n", output: []uint8{instructions.JUMP, 0xab, 0xcd}},
	{input: "$test = 0xbc\nJUMPZ $test\n", output: []uint8{instructions.JUMPZ, 0x0, 0xbc}},
	{input: "$test = 0xabcd\nJUMPZ $test\n", output: []uint8{instructions.JUMPZ, 0xab, 0xcd}},
	{input: "$test = 0xbc\nJUMPNZ $test\n", output: []uint8{instructions.JUMPNZ, 0x0, 0xbc}},
	{input: "$test = 0xabcd\nJUMPNZ $test\n", output: []uint8{instructions.JUMPNZ, 0xab, 0xcd}},
	{input: "test:", output: []uint8{}},
	{input: `
	test:
		SWAP

	JUMP test
	`,
		output: []uint8{instructions.SWAP, instructions.JUMP, 0x0, 0x0},
	},
	{input: `
	JUMP test

	test:
		SWAP
	`,
		output: []uint8{instructions.JUMP, 0x0, 0x3, instructions.SWAP},
	},
	{input: `
	CALL label 
	HALT

	label:
		RET
	`,
		output: []uint8{instructions.CALL, 0x0, 0x4, instructions.HALT, instructions.RET},
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
