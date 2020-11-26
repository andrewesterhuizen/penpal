package assembler

import (
	"fmt"
	"strings"
	"testing"
)

type ParseIncludeTestCase struct {
	input  string
	isFile bool
	name   string
}

var parseIncludeTestCases = []ParseIncludeTestCase{
	{"#include <midi>", false, "midi"},
	{`#include "midi.asm"`, true, "midi.asm"},
}

func TestParseInclude(t *testing.T) {
	for _, tc := range parseIncludeTestCases {
		isFile, name, err := parseInclude(tc.input)

		if err != nil {
			t.Errorf("error parsing include: %v", err)
		}

		if isFile != tc.isFile {
			t.Errorf("expect isFile to be %v and got %v with input %s", tc.isFile, isFile, tc.input)
		}

		if name != tc.name {
			t.Errorf("expect name to be %v and got %v with input %s", tc.name, name, tc.input)
		}

	}
}

func newMockFileGetter(files map[string]string) FileGetter {
	return func(name string) (string, error) {
		file, exists := files[name]
		if !exists {
			return "", fmt.Errorf("failed to get file %s", name)
		}

		return file, nil
	}
}

func newEmptyMockFileGetter() FileGetter {
	return newMockFileGetter(map[string]string{})
}

func TestPreprocessor_ReturnsErrorIfFileGetterFails(t *testing.T) {
	fileGetter := newEmptyMockFileGetter()

	p := NewPreprocessor(fileGetter)

	input := `#include "test"`
	_, err := p.Process(input)

	if err == nil {
		t.Errorf("expected preprocessor to return error if fileGetter failed")
	}
}

func TestPreprocessor_DoesNotChangeInstructions(t *testing.T) {
	fileGetter := newEmptyMockFileGetter()

	p := NewPreprocessor(fileGetter)

	input := `
		MOV A 0x1
		MOV B 0x2
		ADD
		HALT`

	output, _ := p.Process(input)

	if strings.TrimSpace(output) != strings.TrimSpace(input) {
		t.Errorf("expected output to match input and got %s", output)
	}
}

func TestPreprocessor_IncludesFile(t *testing.T) {
	files := map[string]string{}

	files["testfile"] = `
	PUSH 0xae
	SWAP`

	fileGetter := newMockFileGetter(files)
	p := NewPreprocessor(fileGetter)

	input := `
		#include "testfile"

		MOV A 0x1
		MOV B 0x2
		ADD
		HALT`

	output, _ := p.Process(input)

	for file := range files {
		if !strings.Contains(output, files[file]) {
			t.Errorf("expected output to include file %s", file)
		}
	}
}

func TestPreprocessor_IncludesFiles(t *testing.T) {
	files := map[string]string{}

	files["testfile"] = `
	PUSH 0xae
	SWAP`

	files["testfile2"] = `
	MOV A 0xff
	`

	fileGetter := newMockFileGetter(files)
	p := NewPreprocessor(fileGetter)

	input := `
		#include "testfile"
		#include "testfile2"

		MOV A 0x1
		MOV B 0x2
		ADD
		HALT`

	output, _ := p.Process(input)

	for file := range files {
		if !strings.Contains(output, files[file]) {
			t.Errorf("expected output to include file %s", file)
		}
	}
}

func TestPreprocessor_IncludesFilesWithIncludes(t *testing.T) {
	files := map[string]string{}

	files["testfile"] = `
	#include "testfile2"

	PUSH 0xae
	SWAP`

	files["testfile2"] = `
	MOV A 0xff
	`

	fileGetter := newMockFileGetter(files)
	p := NewPreprocessor(fileGetter)

	input := `
		#include "testfile"

		MOV A 0x1
		MOV B 0x2
		ADD
		HALT`

	output, _ := p.Process(input)

	if !strings.Contains(output, files["testfile2"]) {
		t.Errorf("expected output to include file testfile2")
	}
}
