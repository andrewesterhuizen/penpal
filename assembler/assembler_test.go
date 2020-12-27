package assembler

import (
	"fmt"
	"testing"
)

type assemblerTestCase struct {
	input          string
	output         []uint8
	files          map[string]string
	systemIncludes map[string]string
}

var fileIncludeTestCases = []assemblerTestCase{
	{
		input: `
			#include "testfile.asm"

			db 2
			`,
		output: []uint8{1, 2},
		files:  map[string]string{"testfile.asm": `db 1`},
	},
	{
		input: `
			#include "testfile.asm"

			db 3
			`,
		output: []uint8{1, 2, 3},
		files: map[string]string{
			"testfile.asm":  "#include \"testfile2.asm\"\ndb 2\n",
			"testfile2.asm": "db 1\n",
		},
	},
	{
		input: `
			#include <test>

			db 2
			`,
		output:         []uint8{1, 2},
		systemIncludes: map[string]string{"test": `db 1`},
	},
	{
		input: `
			#include <test>

			db 3
			`,
		output: []uint8{1, 2, 3},
		systemIncludes: map[string]string{
			"test":  "#include <test2>\ndb 2\n",
			"test2": "db 1",
		},
	},
	{
		input: `
			#include <test>

			db 3
			`,
		output: []uint8{1, 2, 3},
		files: map[string]string{
			"testfile.asm": "db 1",
		},
		systemIncludes: map[string]string{
			"test": "#include \"testfile.asm\"\ndb 2\n",
		},
	},
	{
		input: `
				#include "testfile.asm"

				db 3
				`,
		output: []uint8{1, 2, 3},
		files: map[string]string{
			"testfile.asm": "#include <test>\ndb 2\n",
		},
		systemIncludes: map[string]string{
			"test": "db 1\n",
		},
	},
}

func newMockFileGetterFunc(files map[string]string) fileGetterFunc {
	return func(name string) (string, error) {
		file, exists := files[name]
		if !exists {
			return "", fmt.Errorf("failed to get file %s", name)
		}

		return file, nil
	}
}

func TestInstructions(t *testing.T) {
	testCases := []assemblerTestCase{}

	testCases = append(testCases, fileIncludeTestCases...)

	for _, tc := range testCases {
		files := tc.files

		if files == nil {
			files = map[string]string{}
		}

		cfg := Config{
			disableEntryPointsTable: true,
			fileGetterFunc:          newMockFileGetterFunc(files),
			SystemIncludes:          tc.systemIncludes,
		}

		a := New(cfg)

		ins, err := a.GetProgram("", fmt.Sprintf("start:\n%s", tc.input))
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

	source := `start:
	jump undefined
	`

	_, err := a.GetProgram("", source)

	if err == nil {
		t.Errorf("expected assembler to return error for undefined label")
	}
}
