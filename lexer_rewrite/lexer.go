package lexer_rewrite

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/andrewesterhuizen/penpal/instructions"
)

const (
	TokenTypeEndOfFile TokenType = iota
	TokenTypeNewLine
	TokenTypeText
	TokenTypeInstruction
	TokenTypeDefine
	TokenTypeFileInclude
	TokenTypeSystemInclude
	TokenTypeInteger
	TokenTypePlus
	TokenTypeMinus
	TokenTypeLeftParen
	TokenTypeRightParen
	TokenTypeLeftBracket
	TokenTypeRightBracket
	TokenTypeComma
	TokenTypeColon
	TokenTypeDoubleQuote
	TokenTypeDot
	TokenTypeAngleBracketLeft
	TokenTypeAngleBracketRight
	TokenTypeLabel
)

const eof = -1

type TokenType int

func (t TokenType) String() string {
	switch t {
	case TokenTypeEndOfFile:
		return "EndOfFile"
	case TokenTypeNewLine:
		return "NewLine"
	case TokenTypeText:
		return "Text"
	case TokenTypeInstruction:
		return "Instruction"
	case TokenTypeInteger:
		return "Integer"
	case TokenTypePlus:
		return "Plus"
	case TokenTypeMinus:
		return "Minus"
	case TokenTypeLeftParen:
		return "LeftParen"
	case TokenTypeRightParen:
		return "RightParen"
	case TokenTypeLeftBracket:
		return "LeftBracket"
	case TokenTypeRightBracket:
		return "RightBracket"
	case TokenTypeComma:
		return "Comma"
	case TokenTypeColon:
		return "Colon"
	case TokenTypeDoubleQuote:
		return "DoubleQuote"
	case TokenTypeDot:
		return "Dot"
	case TokenTypeAngleBracketLeft:
		return "AngleBracketLeft"
	case TokenTypeAngleBracketRight:
		return "AngleBracketRight"
	case TokenTypeLabel:
		return "Label"
	case TokenTypeDefine:
		return "Define"
	case TokenTypeFileInclude:
		return "FileInclude"
	case TokenTypeSystemInclude:
		return "SystemInclude"
	default:
		return "Unknown"
	}
}

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

func (t Token) String() string {
	if t.Value == "\n" {
		t.Value = "(newline)"
	} else if t.Value == "" {
		t.Value = "(empty)"
	} else if strings.TrimSpace(t.Value) == "" {
		t.Value = "(whitespace)"
	}

	return fmt.Sprintf("{ Type: %v, Value: \"%s\", Line: %v, Column: %v }\n", t.Type, t.Value, t.Line, t.Column)
}

type Lexer struct {
	input       string
	start       int
	pos         int
	startOfLine int
	line        int
	tokens      []Token
}

func NewLexer() *Lexer {
	return &Lexer{}
}

func (l *Lexer) Load(input string) {
	l.input = input
	l.start = 0
	l.pos = 0
}

func (l *Lexer) Run() ([]Token, error) {
	for {
		l.start = l.pos
		r := rune(l.input[l.pos])

		switch {
		case unicode.IsDigit(r):
			l.lexInteger()
		case isAlphaNumeric(r):
			l.lexText()

		case r == '#':
			nr := l.peek()
			l.pos++

			switch nr {
			case 'd':
				// TODO: return single token instead of define + text tokens
				// #define test = 0xaa => {type: Define, value: "test" }
				l.lexDefine()
			case 'i':
				err := l.lexInclude()
				if err != nil {
					return nil, err
				}
			default:
				return nil, fmt.Errorf("unexpected next rune %v", string(r))
			}
		case r == '\n':
			l.pos++
			l.addToken(TokenTypeNewLine)
			l.line++

			l.startOfLine = l.pos
		case r == ',':
			l.pos++
			l.addToken(TokenTypeComma)
		case r == '[':
			l.pos++
			l.addToken(TokenTypeLeftBracket)
		case r == ']':
			l.pos++
			l.addToken(TokenTypeRightBracket)
		case r == '(':
			l.pos++
			l.addToken(TokenTypeLeftParen)
		case r == ')':
			l.pos++
			l.addToken(TokenTypeRightParen)
		case r == '+':
			l.pos++
			l.addToken(TokenTypePlus)
		case r == '-':
			l.pos++
			l.addToken(TokenTypeMinus)
		case r == '"':
			l.pos++
			l.addToken(TokenTypeDoubleQuote)
		case r == '.':
			l.pos++
			l.addToken(TokenTypeDot)
		case r == '<':
			l.pos++
			l.addToken(TokenTypeAngleBracketLeft)
		case r == '>':
			l.pos++
			l.addToken(TokenTypeAngleBracketRight)
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
			l.addToken(TokenTypeEndOfFile)
			break
		}
	}

	return l.tokens, nil
}

func (l *Lexer) next() rune {
	l.pos++

	if l.pos >= len(l.input) {
		return eof
	}

	return rune(l.input[l.pos])
}

func (l *Lexer) peek() rune {
	nextPos := l.pos + 1

	if nextPos >= len(l.input) {
		return eof
	}

	return rune(l.input[nextPos])
}

func (l *Lexer) getText() string {
	return l.input[l.start:l.pos]
}

func (l *Lexer) addToken(tokenType TokenType) {
	v := l.input[l.start:l.pos]
	t := Token{Type: tokenType, Value: v, Line: l.line, Column: l.pos - l.startOfLine - len(v)}
	l.tokens = append(l.tokens, t)
	l.start = l.pos
}

func (l *Lexer) lexText() {
	r := rune(l.input[l.pos])

	for isAlphaNumeric(r) {
		r = l.next()
	}

	if r == ':' {
		l.addToken(TokenTypeLabel)
		l.pos++
		return
	}

	text := l.getText()

	_, isInstruction := instructions.InstructionByName[text]
	if isInstruction {
		l.addToken(TokenTypeInstruction)
		return
	}

	l.addToken(TokenTypeText)
	return
}

func (l *Lexer) lexDefine() {
	r := rune(l.input[l.pos])

	for isAlphaNumeric(r) {
		r = l.next()
	}

	l.addToken(TokenTypeDefine)
	return
}

func (l *Lexer) lexInclude() error {
	r := rune(l.input[l.pos])

	// skip "include" text
	for isAlphaNumeric(r) {
		r = l.next()
	}

	r = l.next()

	tt := TokenTypeFileInclude

	switch r {
	case '"':
		tt = TokenTypeFileInclude
	case '<':
		tt = TokenTypeSystemInclude
	default:
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

func (l *Lexer) lexInteger() {
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

	l.addToken(TokenTypeInteger)
	return
}

func isAlphaNumeric(r rune) bool {
	return r == '_' || unicode.IsLetter(r) || unicode.IsDigit(r)
}

func isHex(r rune) bool {
	return unicode.IsDigit(r) || strings.IndexRune("abcdef", r) > 0
}
