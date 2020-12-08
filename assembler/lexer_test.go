package assembler

import (
	"fmt"
	"testing"
)

func newToken(t tokenType, v string) token {
	return token{tokenType: t, value: v}
}

type lexerTestCase struct {
	input  string
	output []token
}

var lexerTestCases = []lexerTestCase{
	{"// args: (status, data1, data2)\n", []token{newToken(tokenTypeNewLine, "\n"), newToken(tokenTypeEndOfFile, "")}},
	{"// db 5\n", []token{newToken(tokenTypeNewLine, "\n"), newToken(tokenTypeEndOfFile, "")}},
	{"; db 5\n", []token{newToken(tokenTypeNewLine, "\n"), newToken(tokenTypeEndOfFile, "")}},
	{"db 5\n", []token{newToken(tokenTypeInstruction, "db"), newToken(tokenTypeInteger, "5"), newToken(tokenTypeNewLine, "\n"), newToken(tokenTypeEndOfFile, "")}},
	{"add\n", []token{newToken(tokenTypeInstruction, "add"), newToken(tokenTypeNewLine, "\n"), newToken(tokenTypeEndOfFile, "")}},
	{"add 5\n", []token{newToken(tokenTypeInstruction, "add"), newToken(tokenTypeInteger, "5"), newToken(tokenTypeNewLine, "\n"), newToken(tokenTypeEndOfFile, "")}},
	{"add 0xbb\n", []token{newToken(tokenTypeInstruction, "add"), newToken(tokenTypeInteger, "0xbb"), newToken(tokenTypeNewLine, "\n"), newToken(tokenTypeEndOfFile, "")}},
	{"add 0x5f\n", []token{newToken(tokenTypeInstruction, "add"), newToken(tokenTypeInteger, "0x5f"), newToken(tokenTypeNewLine, "\n"), newToken(tokenTypeEndOfFile, "")}},
	{"add 0b101\n", []token{newToken(tokenTypeInstruction, "add"), newToken(tokenTypeInteger, "0b101"), newToken(tokenTypeNewLine, "\n"), newToken(tokenTypeEndOfFile, "")}},
	{"mov 23\n", []token{newToken(tokenTypeInstruction, "mov"), newToken(tokenTypeInteger, "23"), newToken(tokenTypeNewLine, "\n"), newToken(tokenTypeEndOfFile, "")}},
	{"mov label\n", []token{newToken(tokenTypeInstruction, "mov"), newToken(tokenTypeText, "label"), newToken(tokenTypeNewLine, "\n"), newToken(tokenTypeEndOfFile, "")}},
	{
		"mov 11, A\n",
		[]token{
			newToken(tokenTypeInstruction, "mov"),
			newToken(tokenTypeInteger, "11"),
			newToken(tokenTypeComma, ","),
			newToken(tokenTypeText, "A"),
			newToken(tokenTypeNewLine, "\n"),
			newToken(tokenTypeEndOfFile, ""),
		},
	},
	{"label:\n", []token{newToken(tokenTypeLabel, "label"), newToken(tokenTypeNewLine, "\n"), newToken(tokenTypeEndOfFile, "")}},
	{
		"mov (fp)\n",
		[]token{
			newToken(tokenTypeInstruction, "mov"),
			newToken(tokenTypeLeftParen, "("),
			newToken(tokenTypeText, "fp"),
			newToken(tokenTypeRightParen, ")"),
			newToken(tokenTypeNewLine, "\n"),
			newToken(tokenTypeEndOfFile, ""),
		},
	},
	{
		"mov (label[3])\n",
		[]token{
			newToken(tokenTypeInstruction, "mov"),
			newToken(tokenTypeLeftParen, "("),
			newToken(tokenTypeText, "label"),
			newToken(tokenTypeLeftBracket, "["),
			newToken(tokenTypeInteger, "3"),
			newToken(tokenTypeRightBracket, "]"),
			newToken(tokenTypeRightParen, ")"),
			newToken(tokenTypeNewLine, "\n"),
			newToken(tokenTypeEndOfFile, ""),
		},
	},
	{
		"mov (fp+1)\n",
		[]token{
			newToken(tokenTypeInstruction, "mov"),
			newToken(tokenTypeLeftParen, "("),
			newToken(tokenTypeText, "fp"),
			newToken(tokenTypePlus, "+"),
			newToken(tokenTypeInteger, "1"),
			newToken(tokenTypeRightParen, ")"),
			newToken(tokenTypeNewLine, "\n"),
			newToken(tokenTypeEndOfFile, ""),
		},
	},
	{
		"mov (fp+1), B\n",
		[]token{
			newToken(tokenTypeInstruction, "mov"),
			newToken(tokenTypeLeftParen, "("),
			newToken(tokenTypeText, "fp"),
			newToken(tokenTypePlus, "+"),
			newToken(tokenTypeInteger, "1"),
			newToken(tokenTypeRightParen, ")"),
			newToken(tokenTypeComma, ","),
			newToken(tokenTypeText, "B"),
			newToken(tokenTypeNewLine, "\n"),
			newToken(tokenTypeEndOfFile, ""),
		},
	},
	{
		"mov A, 67\nadd 13\n",
		[]token{
			newToken(tokenTypeInstruction, "mov"),
			newToken(tokenTypeText, "A"),
			newToken(tokenTypeComma, ","),
			newToken(tokenTypeInteger, "67"),
			newToken(tokenTypeNewLine, "\n"),
			newToken(tokenTypeInstruction, "add"),
			newToken(tokenTypeInteger, "13"),
			newToken(tokenTypeNewLine, "\n"),
			newToken(tokenTypeEndOfFile, ""),
		},
	},
	{
		"#include \"test.asm\"\n",
		[]token{
			newToken(tokenTypeFileInclude, "test.asm"),
			newToken(tokenTypeNewLine, "\n"),
			newToken(tokenTypeEndOfFile, ""),
		},
	},
	{
		"#include <test>\n",
		[]token{
			newToken(tokenTypeSystemInclude, "test"),
			newToken(tokenTypeNewLine, "\n"),
			newToken(tokenTypeEndOfFile, ""),
		},
	},
	{
		"push 0xae\n",
		[]token{
			newToken(tokenTypeInstruction, "push"),
			newToken(tokenTypeInteger, "0xae"),
			newToken(tokenTypeNewLine, "\n"),
			newToken(tokenTypeEndOfFile, ""),
		},
	},
}

func TestLexer(t *testing.T) {
	l := newLexer()

	for _, tc := range lexerTestCases {
		tokens, err := l.Run("", tc.input)
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
			if token.tokenType != tokens[i].tokenType {
				t.Errorf("expected type %v and got %v", token.tokenType, tokens[i].tokenType)
			}

			if token.value != tokens[i].value {
				t.Errorf("expected value %s and got %s for token type %v", token.value, tokens[i].value, token.tokenType)
			}
		}
	}
}
