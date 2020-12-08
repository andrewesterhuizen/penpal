package assembler

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/andrewesterhuizen/penpal/instructions"
)

const (
	tokenTypeEndOfFile tokenType = iota
	tokenTypeNewLine
	tokenTypeText
	tokenTypeInstruction
	tokenTypeFileInclude
	tokenTypeSystemInclude
	tokenTypeInteger
	tokenTypePlus
	tokenTypeMinus
	tokenTypeLeftParen
	tokenTypeRightParen
	tokenTypeLeftBracket
	tokenTypeRightBracket
	tokenTypeComma
	tokenTypeColon
	tokenTypeDoubleQuote
	tokenTypeDot
	tokenTypeAngleBracketLeft
	tokenTypeAngleBracketRight
	tokenTypeLabel
)

const eof = -1

type tokenType int

func (t tokenType) String() string {
	switch t {
	case tokenTypeEndOfFile:
		return "EndOfFile"
	case tokenTypeNewLine:
		return "NewLine"
	case tokenTypeText:
		return "Text"
	case tokenTypeInstruction:
		return "Instruction"
	case tokenTypeInteger:
		return "Integer"
	case tokenTypePlus:
		return "Plus"
	case tokenTypeMinus:
		return "Minus"
	case tokenTypeLeftParen:
		return "LeftParen"
	case tokenTypeRightParen:
		return "RightParen"
	case tokenTypeLeftBracket:
		return "LeftBracket"
	case tokenTypeRightBracket:
		return "RightBracket"
	case tokenTypeComma:
		return "Comma"
	case tokenTypeColon:
		return "Colon"
	case tokenTypeDoubleQuote:
		return "DoubleQuote"
	case tokenTypeDot:
		return "Dot"
	case tokenTypeAngleBracketLeft:
		return "AngleBracketLeft"
	case tokenTypeAngleBracketRight:
		return "AngleBracketRight"
	case tokenTypeLabel:
		return "Label"
	case tokenTypeFileInclude:
		return "FileInclude"
	case tokenTypeSystemInclude:
		return "SystemInclude"
	default:
		return "Unknown"
	}
}

type token struct {
	tokenType tokenType
	value     string
	fileName  string
	line      int
	column    int
}

func (t token) String() string {
	if t.value == "\n" {
		t.value = "(newline)"
	} else if t.value == "" {
		t.value = "(empty)"
	} else if strings.TrimSpace(t.value) == "" {
		t.value = "(whitespace)"
	}

	return fmt.Sprintf("{ tokenType: %v, value: \"%s\", line: %v, column: %v }\n", t.tokenType, t.value, t.line, t.column)
}

type lexer struct {
	input       string
	start       int
	pos         int
	startOfLine int
	line        int
	filename    string
	tokens      []token
}

func newLexer() *lexer {
	return &lexer{}
}

func (l *lexer) reset(filename string, input string) {
	l.filename = filename
	l.input = input
	l.start = 0
	l.pos = 0
	l.line = 1
	l.startOfLine = 0
	l.tokens = []token{}
}

func (l *lexer) errWithPos(err error) error {
	return fmt.Errorf("[%d:%d] %s", l.line, l.getColumn(), err)
}

func (l *lexer) Run(filename string, input string) ([]token, error) {
	l.reset(filename, input)

	for {
		l.start = l.pos
		r := rune(l.input[l.pos])

		switch {
		case unicode.IsDigit(r):
			l.lexInteger()
		case isAlphaNumeric(r):
			l.lexText()

		case r == ';':
			l.skipUntil('\n')

		case r == '/':
			if l.next() != '/' {
				return nil, fmt.Errorf("unexpected character /")
			}

			l.skipUntil('\n')

		case r == '#':
			l.next() // skip #

			err := l.lexInclude()
			if err != nil {
				return nil, l.errWithPos(err)
			}

		case r == '\n':
			l.pos++
			l.addToken(tokenTypeNewLine)
			l.line++

			l.startOfLine = l.pos
		case r == ',':
			l.pos++
			l.addToken(tokenTypeComma)
		case r == '[':
			l.pos++
			l.addToken(tokenTypeLeftBracket)
		case r == ']':
			l.pos++
			l.addToken(tokenTypeRightBracket)
		case r == '(':
			l.pos++
			l.addToken(tokenTypeLeftParen)
		case r == ')':
			l.pos++
			l.addToken(tokenTypeRightParen)
		case r == '+':
			l.pos++
			l.addToken(tokenTypePlus)
		case r == '-':
			l.pos++
			l.addToken(tokenTypeMinus)
		case r == '"':
			l.pos++
			l.addToken(tokenTypeDoubleQuote)
		case r == '.':
			l.pos++
			l.addToken(tokenTypeDot)
		case r == '<':
			l.pos++
			l.addToken(tokenTypeAngleBracketLeft)
		case r == '>':
			l.pos++
			l.addToken(tokenTypeAngleBracketRight)
		case r == ' ':
			// skip
			l.pos++
		case r == '\t':
			// skip
			l.pos++
		default:
			return nil, fmt.Errorf("encountered unexpected rune '%v' (%d)", string(r), r)
		}

		if l.pos >= len(l.input) {
			l.addToken(tokenTypeEndOfFile)
			break
		}
	}

	return l.tokens, nil
}

func (l *lexer) next() rune {
	l.pos++

	if l.pos >= len(l.input) {
		return eof
	}

	return rune(l.input[l.pos])
}

func (l *lexer) peek() rune {
	nextPos := l.pos + 1

	if nextPos >= len(l.input) {
		return eof
	}

	return rune(l.input[nextPos])
}

func (l *lexer) skipUntil(r rune) {
	n := l.next()
	for n != r && n != eof {
		n = l.next()
	}
}

func (l *lexer) getText() string {
	return l.input[l.start:l.pos]
}

func (l *lexer) getColumn() int {
	v := l.input[l.start:l.pos]
	return l.pos - l.startOfLine - len(v)
}

func (l *lexer) addToken(tt tokenType) {
	t := token{
		tokenType: tt,
		value:     l.input[l.start:l.pos],
		fileName:  l.filename,
		line:      l.line,
		column:    l.getColumn(),
	}
	l.tokens = append(l.tokens, t)
	l.start = l.pos
}

func (l *lexer) lexText() {
	r := rune(l.input[l.pos])

	for isAlphaNumeric(r) {
		r = l.next()
	}

	if r == ':' {
		l.addToken(tokenTypeLabel)
		l.pos++
		return
	}

	text := l.getText()

	_, isInstruction := instructions.InstructionByName[text]
	if isInstruction {
		l.addToken(tokenTypeInstruction)
		return
	}

	l.addToken(tokenTypeText)
	return
}

func (l *lexer) lexInclude() error {
	r := rune(l.input[l.pos])

	// skip "include" text
	for isAlphaNumeric(r) {
		r = l.next()
	}

	r = l.next()

	tt := tokenTypeFileInclude

	switch r {
	case '"':
		tt = tokenTypeFileInclude
	case '<':
		tt = tokenTypeSystemInclude
	default:
		fmt.Println(l.input)
		return fmt.Errorf("expected '<' or '\"', got %s", string(r))
	}

	l.pos++
	l.start = l.pos

	r = rune(l.input[l.pos])

	for isAlphaNumeric(r) {
		r = l.next()
	}

	if r == '.' {
		r = l.next()

		// get extension
		for isAlphaNumeric(r) {
			r = l.next()
		}

		l.addToken(tt)
		l.pos++ // skip '>' or '"'
		return nil
	}

	l.addToken(tt)
	l.pos++ // skip '>' or '"'
	return nil
}

func (l *lexer) lexInteger() {
	r := rune(l.input[l.pos])

	// first rune needs to be digit
	if unicode.IsDigit(r) {
		r = l.next()
	}

	if r == 'x' { // hex literal
		r = l.next()
		for isHex(r) {
			r = l.next()
		}

	} else if r == 'b' { // bin literal
		r = l.next()
		for r == '0' || r == '1' {
			r = l.next()
		}

	} else { // decimal
		for unicode.IsDigit(r) {
			r = l.next()
		}

	}

	l.addToken(tokenTypeInteger)
	return
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func isHex(r rune) bool {
	return unicode.IsDigit(r) || strings.IndexRune("abcdefABCDEF", r) >= 0
}
