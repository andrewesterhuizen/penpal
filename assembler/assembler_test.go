package assembler

import (
	"fmt"
	"testing"

	"github.com/andrewesterhuizen/penpal/instructions"
)

type TestCase struct {
	input          string
	output         []uint8
	files          map[string]string
	systemIncludes map[string]string
}

var instructionTestCases = []TestCase{
	{input: "MOV A 0xab", output: []uint8{instructions.MOV, instructions.AddressingModeImmediate, instructions.RegisterA, 0xab}},
	{input: "MOV A +1(fp)", output: []uint8{instructions.MOV, instructions.AddressingModeFPRelative, instructions.RegisterA, 0x1}},
	{input: "MOV A -1(fp)", output: []uint8{instructions.MOV, instructions.AddressingModeFPRelative, instructions.RegisterA, 0xff}},
	{input: "MOV B 0xcd", output: []uint8{instructions.MOV, instructions.AddressingModeImmediate, instructions.RegisterB, 0xcd}},
	{input: "SWAP", output: []uint8{instructions.SWAP}},

	// LOAD dest src
	{input: "LOAD 0xae", output: []uint8{instructions.LOAD, instructions.Register, instructions.RegisterA, 0x00, 0xae}},
	{input: "LOAD 0xaecd", output: []uint8{instructions.LOAD, instructions.Register, instructions.RegisterA, 0xae, 0xcd}},
	{input: "LOAD A 0xae", output: []uint8{instructions.LOAD, instructions.Register, instructions.RegisterA, 0x00, 0xae}},
	{input: "LOAD A 0xaecd", output: []uint8{instructions.LOAD, instructions.Register, instructions.RegisterA, 0xae, 0xcd}},
	{input: "LOAD B 0xae", output: []uint8{instructions.LOAD, instructions.Register, instructions.RegisterB, 0x00, 0xae}},
	{input: "LOAD B 0xaecd", output: []uint8{instructions.LOAD, instructions.Register, instructions.RegisterB, 0xae, 0xcd}},
	{input: "LOAD +1(fp) 0xae", output: []uint8{instructions.LOAD, instructions.FramePointerRelativeAddress, 0x1, 0x00, 0xae}},
	{input: "LOAD +1(fp) 0xaecd", output: []uint8{instructions.LOAD, instructions.FramePointerRelativeAddress, 0x1, 0xae, 0xcd}},
	{input: "LOAD -1(fp) 0xae", output: []uint8{instructions.LOAD, instructions.FramePointerRelativeAddress, 0xff, 0x00, 0xae}},
	{input: "LOAD -1(fp) 0xaecd", output: []uint8{instructions.LOAD, instructions.FramePointerRelativeAddress, 0xff, 0xae, 0xcd}},

	// STORE src dest
	{input: "STORE 0xae", output: []uint8{instructions.STORE, instructions.Register, instructions.RegisterA, 0x00, 0xae}},
	{input: "STORE 0xaecd", output: []uint8{instructions.STORE, instructions.Register, instructions.RegisterA, 0xae, 0xcd}},
	{input: "STORE A 0xae", output: []uint8{instructions.STORE, instructions.Register, instructions.RegisterA, 0x00, 0xae}},
	{input: "STORE B 0xaecd", output: []uint8{instructions.STORE, instructions.Register, instructions.RegisterB, 0xae, 0xcd}},
	{input: "STORE +1(fp) 0xae", output: []uint8{instructions.STORE, instructions.FramePointerRelativeAddress, 0x1, 0x00, 0xae}},
	{input: "STORE +1(fp) 0xaecd", output: []uint8{instructions.STORE, instructions.FramePointerRelativeAddress, 0x1, 0xae, 0xcd}},
	{input: "STORE -1(fp) 0xae", output: []uint8{instructions.STORE, instructions.FramePointerRelativeAddress, 0xff, 0x00, 0xae}},
	{input: "STORE -1(fp) 0xaecd", output: []uint8{instructions.STORE, instructions.FramePointerRelativeAddress, 0xff, 0xae, 0xcd}},

	{input: "POP", output: []uint8{instructions.POP}},
	{input: "PUSH", output: []uint8{instructions.PUSH, instructions.Register, instructions.RegisterA}},
	{input: "PUSH A", output: []uint8{instructions.PUSH, instructions.Register, instructions.RegisterA}},
	{input: "PUSH B", output: []uint8{instructions.PUSH, instructions.Register, instructions.RegisterB}},
	{input: "PUSH 0xbd", output: []uint8{instructions.PUSH, instructions.Value, 0xbd}},
	{input: "PUSH +1(fp)", output: []uint8{instructions.PUSH, instructions.FramePointerRelativeAddress, 0x1}},
	{input: "PUSH -1(fp)", output: []uint8{instructions.PUSH, instructions.FramePointerRelativeAddress, 0xff}},
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
	{input: "RAND", output: []uint8{instructions.RAND}},
	{input: "JUMP 0xcd", output: []uint8{instructions.JUMP, 0x00, 0xcd}},
	{input: "JUMP 0xaecd", output: []uint8{instructions.JUMP, 0xae, 0xcd}},
	{input: "JUMPZ 0xcd", output: []uint8{instructions.JUMPZ, 0x00, 0xcd}},
	{input: "JUMPZ 0xaecd", output: []uint8{instructions.JUMPZ, 0xae, 0xcd}},
	{input: "JUMPNZ 0xcd", output: []uint8{instructions.JUMPNZ, 0x00, 0xcd}},
	{input: "JUMPNZ 0xaecd", output: []uint8{instructions.JUMPNZ, 0xae, 0xcd}},
	{input: "#define TEST 0xbc\nMOV A TEST\n", output: []uint8{instructions.MOV, instructions.AddressingModeImmediate, instructions.RegisterA, 0xbc}},
	{input: "#define TEST 0xabcd\nMOV A TEST\n", output: []uint8{instructions.MOV, instructions.AddressingModeImmediate, instructions.RegisterA, 0xcd}},
	{input: "#define TEST 0xbc\nMOV A TEST\n", output: []uint8{instructions.MOV, instructions.AddressingModeImmediate, instructions.RegisterA, 0xbc}},
	{input: "#define TEST 0xabcd\nMOV A TEST\n", output: []uint8{instructions.MOV, instructions.AddressingModeImmediate, instructions.RegisterA, 0xcd}},
	{input: "#define TEST 0xbc\nCALL TEST\n", output: []uint8{instructions.CALL, 0x0, 0xbc}},
	{input: "#define TEST 0xabcd\nCALL TEST\n", output: []uint8{instructions.CALL, 0xab, 0xcd}},
	{input: "#define TEST 0xbc\nJUMP TEST\n", output: []uint8{instructions.JUMP, 0x0, 0xbc}},
	{input: "#define TEST 0xabcd\nJUMP TEST\n", output: []uint8{instructions.JUMP, 0xab, 0xcd}},
	{input: "#define TEST 0xbc\nJUMPZ TEST\n", output: []uint8{instructions.JUMPZ, 0x0, 0xbc}},
	{input: "#define TEST 0xabcd\nJUMPZ TEST\n", output: []uint8{instructions.JUMPZ, 0xab, 0xcd}},
	{input: "#define TEST 0xbc\nJUMPNZ TEST\n", output: []uint8{instructions.JUMPNZ, 0x0, 0xbc}},
	{input: "#define TEST 0xabcd\nJUMPNZ TEST\n", output: []uint8{instructions.JUMPNZ, 0xab, 0xcd}},
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
	{
		input: `
			#include "testfile.asm"

			HALT`,
		output: []uint8{instructions.SWAP, instructions.HALT},
		files:  map[string]string{"testfile.asm": `SWAP`},
	},
	{
		input: `
			#include "testfile.asm"

			HALT`,
		output: []uint8{instructions.SHL, instructions.SHR, instructions.ADD, instructions.HALT},
		files: map[string]string{
			"testfile.asm":  "#include \"testfile2.asm\"\nADD\n",
			"testfile2.asm": "SHL\nSHR\n",
		},
	},
	{
		input: `
			#include <test>

			HALT`,
		output:         []uint8{instructions.SUB, instructions.HALT},
		systemIncludes: map[string]string{"test": `SUB`},
	},
	{
		input: `
			#include <test>

			HALT`,
		output: []uint8{instructions.SHL, instructions.SHR, instructions.ADD, instructions.HALT},
		systemIncludes: map[string]string{
			"test":  "#include <test2>\nADD\n",
			"test2": "SHL\nSHR\n",
		},
	},
	{
		input: `
			#include <test>

			HALT`,
		output: []uint8{instructions.SHL, instructions.SHR, instructions.ADD, instructions.HALT},
		files: map[string]string{
			"testfile.asm": "SHL\nSHR\n",
		},
		systemIncludes: map[string]string{
			"test": "#include \"testfile.asm\"\nADD\n",
		},
	},
	{
		input: `
				#include "testfile.asm"
	
				HALT`,
		output: []uint8{instructions.MUL, instructions.SUB, instructions.HALT},
		files: map[string]string{
			"testfile.asm": "#include <test>\nSUB",
		},
		systemIncludes: map[string]string{
			"test": "MUL\n",
		},
	},
}

func newMockFileGetterFunc(files map[string]string) FileGetterFunc {
	return func(name string) (string, error) {
		file, exists := files[name]
		if !exists {
			return "", fmt.Errorf("failed to get file %s", name)
		}

		return file, nil
	}
}

func TestInstructions(t *testing.T) {
	for _, tc := range instructionTestCases {
		files := tc.files

		if files == nil {
			files = map[string]string{}
		}

		cfg := Config{
			disableHeader:  true,
			FileGetterFunc: newMockFileGetterFunc(files),
			SystemIncludes: tc.systemIncludes,
		}

		a := New(cfg)

		ins, err := a.GetProgram("", tc.input)
		if err != nil {
			t.Errorf("failed to get instructions due to error: %s, with input %s", err, tc.input)
			return
		}

		if len(ins) != len(tc.output) {
			t.Errorf("expected %d instuctions and got %d, with input %s", len(tc.output), len(ins), tc.input)
			return
		}

		for i, in := range tc.output {
			if in != ins[i] {
				t.Errorf("expected 0x%02x and got 0x%02x at pos %d, with input %s", in, ins[i], i, tc.input)
			}
		}
	}
}

func TestAssembler_NoEntryPoint_ReturnsError(t *testing.T) {
	a := New(Config{})
	_, err := a.GetProgram("", "SWAP\n")

	if err == nil {
		t.Errorf("expected assembler to return error for source with no entry point")
	}
}

// this test should fail if I forget to update the HeaderSize constant
func TestAssembler_EntryPointSkipsHeader(t *testing.T) {
	a := New(Config{})
	ins, _ := a.GetProgram("", "__start:\nSWAP\nOR\nADD\n")

	if (ins[HeaderSize] != instructions.SWAP) || (ins[HeaderSize+1] != instructions.OR) || (ins[HeaderSize+2] != instructions.ADD) {
		t.Errorf("expected assembler to return error for source with no entry point")
	}
}

func TestAssembler_MissingSystemInclude_ReturnsError(t *testing.T) {
	a := New(Config{disableHeader: true})

	source := `#include <missing>`

	_, err := a.GetProgram("", source)

	if err == nil {
		t.Errorf("expected assembler to return error for missing system include")
	}
}

func TestAssembler_UndefinedLabel_ReturnsError(t *testing.T) {
	a := New(Config{disableHeader: true})

	source := `JUMP undefined`

	_, err := a.GetProgram("", source)

	if err == nil {
		t.Errorf("expected assembler to return error for undefined label")
	}
}

func TestAssembler_EntryPointCorrectWithSystemInclude(t *testing.T) {
	systemIncludes := map[string]string{"test": "SWAP\n"}
	a := New(Config{disableHeader: false, SystemIncludes: systemIncludes})

	source := `
	#include <test>

	__start:
    JUMP __start
`

	program, err := a.GetProgram("", source)
	if err != nil {
		t.Errorf("test failed with error %s", err)
	}

	startIndex := HeaderSize + 1 // header + include with 1 instruction

	if program[startIndex] != instructions.JUMP {
		t.Errorf("incorrect entry point in header")
	}

	if program[startIndex+2] != uint8(startIndex) {
		t.Errorf("incorrect entry point in header")
	}

}
