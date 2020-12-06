package lexer_rewrite

import (
	"fmt"
	"testing"
)

func newToken(t TokenType, v string) Token {
	return Token{Type: t, Value: v}
}

type testCase struct {
	input  string
	output []Token
}

var testCases = []testCase{
	{"db 5\n", []Token{newToken(TokenTypeInstruction, "db"), newToken(TokenTypeInteger, "5"), newToken(TokenTypeNewLine, "\n"), newToken(TokenTypeEndOfFile, "")}},
	{"add\n", []Token{newToken(TokenTypeInstruction, "add"), newToken(TokenTypeNewLine, "\n"), newToken(TokenTypeEndOfFile, "")}},
	{"add 5\n", []Token{newToken(TokenTypeInstruction, "add"), newToken(TokenTypeInteger, "5"), newToken(TokenTypeNewLine, "\n"), newToken(TokenTypeEndOfFile, "")}},
	{"add 0xbb\n", []Token{newToken(TokenTypeInstruction, "add"), newToken(TokenTypeInteger, "0xbb"), newToken(TokenTypeNewLine, "\n"), newToken(TokenTypeEndOfFile, "")}},
	{"add 0x5f\n", []Token{newToken(TokenTypeInstruction, "add"), newToken(TokenTypeInteger, "0x5f"), newToken(TokenTypeNewLine, "\n"), newToken(TokenTypeEndOfFile, "")}},
	{"add 0b101\n", []Token{newToken(TokenTypeInstruction, "add"), newToken(TokenTypeInteger, "0b101"), newToken(TokenTypeNewLine, "\n"), newToken(TokenTypeEndOfFile, "")}},
	{"mov 23\n", []Token{newToken(TokenTypeInstruction, "mov"), newToken(TokenTypeInteger, "23"), newToken(TokenTypeNewLine, "\n"), newToken(TokenTypeEndOfFile, "")}},
	{"mov label\n", []Token{newToken(TokenTypeInstruction, "mov"), newToken(TokenTypeText, "label"), newToken(TokenTypeNewLine, "\n"), newToken(TokenTypeEndOfFile, "")}},
	{
		"mov 11, A\n",
		[]Token{
			newToken(TokenTypeInstruction, "mov"),
			newToken(TokenTypeInteger, "11"),
			newToken(TokenTypeComma, ","),
			newToken(TokenTypeText, "A"),
			newToken(TokenTypeNewLine, "\n"),
			newToken(TokenTypeEndOfFile, ""),
		},
	},
	{"label:\n", []Token{newToken(TokenTypeLabel, "label"), newToken(TokenTypeNewLine, "\n"), newToken(TokenTypeEndOfFile, "")}},
	{
		"mov (fp)\n",
		[]Token{
			newToken(TokenTypeInstruction, "mov"),
			newToken(TokenTypeLeftParen, "("),
			newToken(TokenTypeText, "fp"),
			newToken(TokenTypeRightParen, ")"),
			newToken(TokenTypeNewLine, "\n"),
			newToken(TokenTypeEndOfFile, ""),
		},
	},
	{
		"mov (label[3])\n",
		[]Token{
			newToken(TokenTypeInstruction, "mov"),
			newToken(TokenTypeLeftParen, "("),
			newToken(TokenTypeText, "label"),
			newToken(TokenTypeLeftBracket, "["),
			newToken(TokenTypeInteger, "3"),
			newToken(TokenTypeRightBracket, "]"),
			newToken(TokenTypeRightParen, ")"),
			newToken(TokenTypeNewLine, "\n"),
			newToken(TokenTypeEndOfFile, ""),
		},
	},
	{
		"mov (fp+1)\n",
		[]Token{
			newToken(TokenTypeInstruction, "mov"),
			newToken(TokenTypeLeftParen, "("),
			newToken(TokenTypeText, "fp"),
			newToken(TokenTypePlus, "+"),
			newToken(TokenTypeInteger, "1"),
			newToken(TokenTypeRightParen, ")"),
			newToken(TokenTypeNewLine, "\n"),
			newToken(TokenTypeEndOfFile, ""),
		},
	},
	{
		"mov (fp+1), B\n",
		[]Token{
			newToken(TokenTypeInstruction, "mov"),
			newToken(TokenTypeLeftParen, "("),
			newToken(TokenTypeText, "fp"),
			newToken(TokenTypePlus, "+"),
			newToken(TokenTypeInteger, "1"),
			newToken(TokenTypeRightParen, ")"),
			newToken(TokenTypeComma, ","),
			newToken(TokenTypeText, "B"),
			newToken(TokenTypeNewLine, "\n"),
			newToken(TokenTypeEndOfFile, ""),
		},
	},
	{
		"mov A, 67\nadd 13\n",
		[]Token{
			newToken(TokenTypeInstruction, "mov"),
			newToken(TokenTypeText, "A"),
			newToken(TokenTypeComma, ","),
			newToken(TokenTypeInteger, "67"),
			newToken(TokenTypeNewLine, "\n"),
			newToken(TokenTypeInstruction, "add"),
			newToken(TokenTypeInteger, "13"),
			newToken(TokenTypeNewLine, "\n"),
			newToken(TokenTypeEndOfFile, ""),
		},
	},
	{
		"#include \"test.asm\"\n",
		[]Token{
			newToken(TokenTypeFileInclude, "test.asm"),
			newToken(TokenTypeNewLine, "\n"),
			newToken(TokenTypeEndOfFile, ""),
		},
	},
	{
		"#include <test>\n",
		[]Token{
			newToken(TokenTypeSystemInclude, "test"),
			newToken(TokenTypeNewLine, "\n"),
			newToken(TokenTypeEndOfFile, ""),
		},
	},
}

func TestLexer(t *testing.T) {
	for _, tc := range testCases {
		l := NewLexer()
		l.Load(tc.input)
		tokens, err := l.Run()
		if err != nil {
			t.Error(err)
			return
		}

		fmt.Println(tokens)
		fmt.Println()

		if len(tc.output) != len(tokens) {
			t.Errorf("expected %v tokens and got %v", len(tc.output), len(tokens))
			return
		}

		for i, token := range tc.output {
			if token.Type != tokens[i].Type {
				t.Errorf("expected type %v and got %v", token.Type, tokens[i].Type)
			}

			if token.Value != tokens[i].Value {
				t.Errorf("expected value %s and got %s for token type %v", token.Value, tokens[i].Value, token.Type)
			}
		}
	}
}
