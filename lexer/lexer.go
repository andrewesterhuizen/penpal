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

var includeRegex = regexp.MustCompile(`#include\s+(<|")(\w+.\w+)["|>]\s*`)
var defineRegex = regexp.MustCompile(`#define\s+(\w+)\s+(\w+)`)

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
		// skip
		return nil
	case defineRegex.MatchString(line):
		s := defineRegex.FindStringSubmatch(line)

		name := s[1]
		value := s[2]

		t.Type = TokenDefine
		t.Value = name

		t.Args = []Arg{newArg(value)}
	case includeRegex.MatchString(line):
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
