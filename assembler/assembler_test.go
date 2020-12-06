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
	{input: "mov A 0xab", output: []uint8{instructions.Mov, instructions.AddressingModeImmediate, instructions.RegisterA, 0xab}},
	{input: "mov A +1(fp)", output: []uint8{instructions.Mov, instructions.AddressingModeFPRelative, instructions.RegisterA, 0x1}},
	{input: "mov A -1(fp)", output: []uint8{instructions.Mov, instructions.AddressingModeFPRelative, instructions.RegisterA, 0xff}},
	{input: "mov B 0xcd", output: []uint8{instructions.Mov, instructions.AddressingModeImmediate, instructions.RegisterB, 0xcd}},
	{input: "swap", output: []uint8{instructions.Swap}},

	// load dest src
	{input: "load 0xae", output: []uint8{instructions.Load, instructions.Register, instructions.RegisterA, 0x00, 0xae}},
	{input: "load 0xaecd", output: []uint8{instructions.Load, instructions.Register, instructions.RegisterA, 0xae, 0xcd}},
	{input: "load A 0xae", output: []uint8{instructions.Load, instructions.Register, instructions.RegisterA, 0x00, 0xae}},
	{input: "load A 0xaecd", output: []uint8{instructions.Load, instructions.Register, instructions.RegisterA, 0xae, 0xcd}},
	{input: "load B 0xae", output: []uint8{instructions.Load, instructions.Register, instructions.RegisterB, 0x00, 0xae}},
	{input: "load B 0xaecd", output: []uint8{instructions.Load, instructions.Register, instructions.RegisterB, 0xae, 0xcd}},
	{input: "load +1(fp) 0xae", output: []uint8{instructions.Load, instructions.FramePointerWithOffset, 0x1, 0x00, 0xae}},
	{input: "load +1(fp) 0xaecd", output: []uint8{instructions.Load, instructions.FramePointerWithOffset, 0x1, 0xae, 0xcd}},
	{input: "load -1(fp) 0xae", output: []uint8{instructions.Load, instructions.FramePointerWithOffset, 0xff, 0x00, 0xae}},
	{input: "load -1(fp) 0xaecd", output: []uint8{instructions.Load, instructions.FramePointerWithOffset, 0xff, 0xae, 0xcd}},

	// store src dest
	{input: "store 0xae", output: []uint8{instructions.Store, instructions.Register, instructions.RegisterA, 0x00, 0xae}},
	{input: "store 0xaecd", output: []uint8{instructions.Store, instructions.Register, instructions.RegisterA, 0xae, 0xcd}},
	{input: "store A 0xae", output: []uint8{instructions.Store, instructions.Register, instructions.RegisterA, 0x00, 0xae}},
	{input: "store B 0xaecd", output: []uint8{instructions.Store, instructions.Register, instructions.RegisterB, 0xae, 0xcd}},
	{input: "store +1(fp) 0xae", output: []uint8{instructions.Store, instructions.FramePointerWithOffset, 0x1, 0x00, 0xae}},
	{input: "store +1(fp) 0xaecd", output: []uint8{instructions.Store, instructions.FramePointerWithOffset, 0x1, 0xae, 0xcd}},
	{input: "store -1(fp) 0xae", output: []uint8{instructions.Store, instructions.FramePointerWithOffset, 0xff, 0x00, 0xae}},
	{input: "store -1(fp) 0xaecd", output: []uint8{instructions.Store, instructions.FramePointerWithOffset, 0xff, 0xae, 0xcd}},

	{input: "pop", output: []uint8{instructions.Pop}},
	{input: "push", output: []uint8{instructions.Push, instructions.Register, instructions.RegisterA}},
	{input: "push A", output: []uint8{instructions.Push, instructions.Register, instructions.RegisterA}},
	{input: "push B", output: []uint8{instructions.Push, instructions.Register, instructions.RegisterB}},
	{input: "push 0xbd", output: []uint8{instructions.Push, instructions.Value, 0xbd}},
	{input: "push +1(fp)", output: []uint8{instructions.Push, instructions.FramePointerWithOffset, 0x1}},
	{input: "push -1(fp)", output: []uint8{instructions.Push, instructions.FramePointerWithOffset, 0xff}},
	{input: "call 0xbeaf", output: []uint8{instructions.Call, 0xbe, 0xaf}},
	{input: "ret", output: []uint8{instructions.Ret}},
	{input: "reti", output: []uint8{instructions.Reti}},
	{input: "add", output: []uint8{instructions.Add}},
	{input: "sub", output: []uint8{instructions.Sub}},
	{input: "mul", output: []uint8{instructions.Mul}},
	{input: "div", output: []uint8{instructions.Div}},
	{input: "shl", output: []uint8{instructions.Shl}},
	{input: "shr", output: []uint8{instructions.Shr}},
	{input: "and", output: []uint8{instructions.And}},
	{input: "or", output: []uint8{instructions.Or}},
	{input: "halt", output: []uint8{instructions.Halt}},
	{input: "rand", output: []uint8{instructions.Rand}},
	{input: "db 0xeb", output: []uint8{0xeb}},
	{input: "jump 0xcd", output: []uint8{instructions.Jump, 0x00, 0xcd}},
	{input: "jump 0xaecd", output: []uint8{instructions.Jump, 0xae, 0xcd}},
	{input: "jumpz 0xcd", output: []uint8{instructions.Jumpz, 0x00, 0xcd}},
	{input: "jumpz 0xaecd", output: []uint8{instructions.Jumpz, 0xae, 0xcd}},
	{input: "jumpnz 0xcd", output: []uint8{instructions.Jumpnz, 0x00, 0xcd}},
	{input: "jumpnz 0xaecd", output: []uint8{instructions.Jumpnz, 0xae, 0xcd}},
	{input: "#define TEST 0xbc\nmov A TEST\n", output: []uint8{instructions.Mov, instructions.AddressingModeImmediate, instructions.RegisterA, 0xbc}},
	{input: "#define TEST 0xabcd\nmov A TEST\n", output: []uint8{instructions.Mov, instructions.AddressingModeImmediate, instructions.RegisterA, 0xcd}},
	{input: "#define TEST 0xbc\nmov A TEST\n", output: []uint8{instructions.Mov, instructions.AddressingModeImmediate, instructions.RegisterA, 0xbc}},
	{input: "#define TEST 0xabcd\nmov A TEST\n", output: []uint8{instructions.Mov, instructions.AddressingModeImmediate, instructions.RegisterA, 0xcd}},
	{input: "#define TEST 0xbc\ncall TEST\n", output: []uint8{instructions.Call, 0x0, 0xbc}},
	{input: "#define TEST 0xabcd\ncall TEST\n", output: []uint8{instructions.Call, 0xab, 0xcd}},
	{input: "#define TEST 0xbc\njump TEST\n", output: []uint8{instructions.Jump, 0x0, 0xbc}},
	{input: "#define TEST 0xabcd\njump TEST\n", output: []uint8{instructions.Jump, 0xab, 0xcd}},
	{input: "#define TEST 0xbc\njumpz TEST\n", output: []uint8{instructions.Jumpz, 0x0, 0xbc}},
	{input: "#define TEST 0xabcd\njumpz TEST\n", output: []uint8{instructions.Jumpz, 0xab, 0xcd}},
	{input: "#define TEST 0xbc\njumpnz TEST\n", output: []uint8{instructions.Jumpnz, 0x0, 0xbc}},
	{input: "#define TEST 0xabcd\njumpnz TEST\n", output: []uint8{instructions.Jumpnz, 0xab, 0xcd}},
	{input: "test:", output: []uint8{}},
	{input: `
	test:
		swap

	jump test
	`,
		output: []uint8{instructions.Swap, instructions.Jump, 0x0, 0x0},
	},
	{input: `
	jump test

	test:
		swap
	`,
		output: []uint8{instructions.Jump, 0x0, 0x3, instructions.Swap},
	},
	{input: `
	call label
	halt

	label:
		ret
	`,
		output: []uint8{instructions.Call, 0x0, 0x4, instructions.Halt, instructions.Ret},
	},
	{
		input: `
			#include "testfile.asm"

			halt`,
		output: []uint8{instructions.Swap, instructions.Halt},
		files:  map[string]string{"testfile.asm": `swap`},
	},
	{
		input: `
			#include "testfile.asm"

			halt`,
		output: []uint8{instructions.Shl, instructions.Shr, instructions.Add, instructions.Halt},
		files: map[string]string{
			"testfile.asm":  "#include \"testfile2.asm\"\nadd\n",
			"testfile2.asm": "shl\nshr\n",
		},
	},
	{
		input: `
			#include <test>

			halt`,
		output:         []uint8{instructions.Sub, instructions.Halt},
		systemIncludes: map[string]string{"test": `sub`},
	},
	{
		input: `
			#include <test>

			halt`,
		output: []uint8{instructions.Shl, instructions.Shr, instructions.Add, instructions.Halt},
		systemIncludes: map[string]string{
			"test":  "#include <test2>\nadd\n",
			"test2": "shl\nshr\n",
		},
	},
	{
		input: `
			#include <test>

			halt`,
		output: []uint8{instructions.Shl, instructions.Shr, instructions.Add, instructions.Halt},
		files: map[string]string{
			"testfile.asm": "shl\nshr\n",
		},
		systemIncludes: map[string]string{
			"test": "#include \"testfile.asm\"\nadd\n",
		},
	},
	{
		input: `
				#include "testfile.asm"
	
				halt`,
		output: []uint8{instructions.Mul, instructions.Sub, instructions.Halt},
		files: map[string]string{
			"testfile.asm": "#include <test>\nsub",
		},
		systemIncludes: map[string]string{
			"test": "mul\n",
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
			disableEntryPointsTable: true,
			FileGetterFunc:          newMockFileGetterFunc(files),
			SystemIncludes:          tc.systemIncludes,
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
	_, err := a.GetProgram("", "swap\n")

	if err == nil {
		t.Errorf("expected assembler to return error for source with no entry point")
	}
}

func TestAssembler_getEntryPointTableBytes_returnsExpectedBytes(t *testing.T) {
	interuptLabels := [3]string{"test", "", "test2"}
	a := New(Config{InteruptLabels: interuptLabels})

	a.labels["__start"] = 0x1234
	a.labels["test"] = 0xabcd
	a.labels["test2"] = 0xa1b2

	entryPointBytes, err := a.getEntryPointTableBytes()
	if err != nil {
		t.Errorf("unexpected error while getting entry point bytes: %s", err)
	}

	expected := []byte{
		instructions.Jump, 0x12, 0x34,
		instructions.Jump, 0xab, 0xcd,
		0x0, 0x0, 0x0,
		instructions.Jump, 0xa1, 0xb2,
	}

	for i, b := range expected {
		if b != entryPointBytes[i] {
			t.Errorf("expected byte 0x%02x and got 0x%02x, at position %d", b, entryPointBytes[i], i)

		}
	}
}

func TestAssembler_AddsEntryPointBytesAtExpectedLocation(t *testing.T) {
	a := New(Config{})
	program, _ := a.GetProgram("", "__start:halt\n")

	entryPoint := EntryPointsTableSize
	eph := byte((entryPoint & 0xff00) >> 8)
	epl := byte(entryPoint & 0xff)

	if program[0] != instructions.Jump || program[1] != eph || program[2] != epl {
		t.Errorf("expected assembler to add entry point bytes")
	}
}

func TestAssembler_MissingSystemInclude_ReturnsError(t *testing.T) {
	a := New(Config{})

	source := `#include <missing>`

	_, err := a.GetProgram("", source)

	if err == nil {
		t.Errorf("expected assembler to return error for missing system include")
	}
}

func TestAssembler_UndefinedLabel_ReturnsError(t *testing.T) {
	a := New(Config{})

	source := `jump undefined`

	_, err := a.GetProgram("", source)

	if err == nil {
		t.Errorf("expected assembler to return error for undefined label")
	}
}

func TestAssembler_EntryPointCorrectWithSystemInclude(t *testing.T) {
	systemIncludes := map[string]string{"test": "swap\n"}
	a := New(Config{
		disableEntryPointsTable: false,
		SystemIncludes:          systemIncludes,
	})

	source := `
	#include <test>

	__start:
    jump __start
`

	program, err := a.GetProgram("", source)
	if err != nil {
		t.Errorf("test failed with error %s", err)
	}

	startIndex := EntryPointsTableSize + 1

	if program[startIndex] != instructions.Jump {
		t.Errorf("incorrect entry point in header")
	}

	if program[startIndex+2] != uint8(startIndex) {
		t.Errorf("incorrect entry point in header")
	}
}
