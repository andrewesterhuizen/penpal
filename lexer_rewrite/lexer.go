package lexer_rewrite

import (
	"fmt"
	"strings"
	"unicode"
)

const (
	TokenTypeEndOfFile TokenType = iota
	TokenTypeNewLine
	TokenTypeText
	TokenTypeDefine
	TokenTypeInclude
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
	case TokenTypeAngleBracketLeft:
		return "AngleBracketLeft"
	case TokenTypeAngleBracketRight:
		return "AngleBracketRight"
	case TokenTypeLabel:
		return "Label"
	case TokenTypeDefine:
		return "Define"
	case TokenTypeInclude:
		return "Include"
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
				// TODO: return single token instead of define + text tokens
				// #include "test" / <test> => {type: Include/SystemInclude, value: "test" }
				l.lexInclude()
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
		case r == '<':
			l.pos++
			l.addToken(TokenTypeAngleBracketLeft)
		case r == '>':
			l.pos++
			l.addToken(TokenTypeAngleBracketRight)
		case r == ' ':
			// skip
			l.pos++
		default:
			return nil, fmt.Errorf("encountered unexpected rune %v", string(r))
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

func (l *Lexer) lexInclude() {
	r := rune(l.input[l.pos])

	for isAlphaNumeric(r) {
		r = l.next()
	}

	l.addToken(TokenTypeInclude)
	return
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
