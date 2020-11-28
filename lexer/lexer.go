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

var lineCommentRegex = regexp.MustCompile(`^\s*[//|;]`)
var labelRegex = regexp.MustCompile(`(\w+):`)
var includeRegex = regexp.MustCompile(`#include\s+(<|")(\w+.\w+)["|>]\s*`)
var defineRegex = regexp.MustCompile(`#define\s+(\w+)\s+(\w+)`)
var instructionRegex = regexp.MustCompile(`([\w\(\)\+\-]+)`)

func (l *Lexer) parseLine(filename string, lineNumber int, line string) error {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	t := Token{FileName: filename, LineNumber: lineNumber}

	switch {
	case lineCommentRegex.MatchString(line):
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

	case labelRegex.MatchString(line):
		s := labelRegex.FindStringSubmatch(line)
		t.Type = TokenLabel
		t.Value = s[1]

	default:
		matches := instructionRegex.FindAllString(line, -1)

		t.Type = TokenInstruction
		t.Value = matches[0]

		operands := matches[1:]

		args := []Arg{}

		for _, o := range operands {
			if o == " " {
				continue
			}

			a := newArg(o)

			args = append(args, a)
		}

		t.Args = args
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
