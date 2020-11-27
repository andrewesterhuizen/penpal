package lexer

import (
	"fmt"
	"regexp"
	"strings"
)

type Lexer struct {
	source string
	tokens []Token
}

func New() Lexer {
	l := Lexer{}
	return l
}

func getInstructionArgs(parts []string) []Arg {
	args := []Arg{}
	for _, s := range parts {
		if s == " " {
			continue
		}

		a := newArg(s)

		args = append(args, a)
	}

	return args
}

func getDefineArgs(parts []string) ([]Arg, error) {
	partsNoWhiteSpace := []string{}

	for _, s := range parts {
		if s == " " {
			continue
		}

		partsNoWhiteSpace = append(partsNoWhiteSpace, s)
	}

	parts = partsNoWhiteSpace

	if parts[0] != "=" {
		return nil, fmt.Errorf("expected '=' and got %s", parts[0])
	}

	return []Arg{newArg(parts[1])}, nil
}

var includeRegex = regexp.MustCompile(`#include\s+(<|")(\w+.\w+)["|>]\s*`)

func (l *Lexer) parseLine(filename string, lineNumber int, line string) error {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	t := Token{FileName: filename, LineNumber: lineNumber}

	lineSplit := strings.Split(line, " ")

	v := strings.TrimSpace(lineSplit[0])
	args := lineSplit[1:]

	switch {
	case v[0:2] == "//": // comment
		// skip, might need to be smarter about this in the future
		// if we want to give error messages with line numbers
		return nil
	case v[0] == '$': // define
		t.Type = TokenDefine
		t.Value = v[1:]
		args, err := getDefineArgs(args)
		if err != nil {
			return fmt.Errorf("failed to parse define %s: %w", v, err)
		}

		t.Args = args
	case v[0] == '#':
		s := includeRegex.FindStringSubmatch(line)

		startChar := s[1]
		filename := s[2]

		if startChar == "<" {
			t.Value = filename
			t.Type = TokenSystemInclude
		} else if startChar == "\"" {
			t.Value = filename
			t.Type = TokenFileInclude
		} else {
			return fmt.Errorf("expected next char to be \" or < and got %v", startChar)
		}

	case v[len(v)-1] == ':': // label
		t.Type = TokenLabel
		t.Value = v[0 : len(v)-1]

	default: // instruction
		t.Type = TokenInstruction
		t.Value = v
		t.Args = getInstructionArgs(args)
	}

	l.tokens = append(l.tokens, t)

	return nil
}

func (l *Lexer) GetTokens(filename string, source string) ([]Token, error) {
	l.source = strings.TrimSpace(source)
	lines := strings.Split(l.source, "\n")

	for i, line := range lines {
		err := l.parseLine(filename, i, line)
		if err != nil {
			return nil, fmt.Errorf("lexing failed at line %d, (%s): %w", i, line, err)
		}
	}

	return l.tokens, nil
}
