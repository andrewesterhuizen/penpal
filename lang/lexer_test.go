package lang

import (
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
	{
		input: `var x = 2;
		var y = 3;
		var z;
		if (x < y) {
			z = x + y;
		}`,
		output: []token{
			newToken(tokenTypeKeyword, "var"),
			newToken(tokenTypeIdentifier, "x"),
			newToken(tokenTypeEquals, "="),
			newToken(tokenTypeInteger, "2"),
			newToken(tokenTypeSemiColon, ";"),
			newToken(tokenTypeNewLine, "\n"),

			newToken(tokenTypeKeyword, "var"),
			newToken(tokenTypeIdentifier, "y"),
			newToken(tokenTypeEquals, "="),
			newToken(tokenTypeInteger, "3"),
			newToken(tokenTypeSemiColon, ";"),
			newToken(tokenTypeNewLine, "\n"),

			newToken(tokenTypeKeyword, "var"),
			newToken(tokenTypeIdentifier, "z"),
			newToken(tokenTypeSemiColon, ";"),
			newToken(tokenTypeNewLine, "\n"),

			newToken(tokenTypeKeyword, "if"),
			newToken(tokenTypeLeftParen, "("),
			newToken(tokenTypeIdentifier, "x"),
			newToken(tokenTypeLessThan, "<"),
			newToken(tokenTypeIdentifier, "y"),
			newToken(tokenTypeRightParen, ")"),
			newToken(tokenTypeLeftBrace, "{"),
			newToken(tokenTypeNewLine, "\n"),

			newToken(tokenTypeIdentifier, "z"),
			newToken(tokenTypeEquals, "="),
			newToken(tokenTypeIdentifier, "x"),
			newToken(tokenTypePlus, "+"),
			newToken(tokenTypeIdentifier, "y"),
			newToken(tokenTypeSemiColon, ";"),
			newToken(tokenTypeNewLine, "\n"),

			newToken(tokenTypeRightBrace, "}"),
			newToken(tokenTypeEndOfFile, ""),
		},
	},
	{
		input: `func test() {}`,
		output: []token{
			newToken(tokenTypeKeyword, "func"),
			newToken(tokenTypeIdentifier, "test"),
			newToken(tokenTypeLeftParen, "("),
			newToken(tokenTypeRightParen, ")"),
			newToken(tokenTypeLeftBrace, "{"),
			newToken(tokenTypeRightBrace, "}"),
			newToken(tokenTypeEndOfFile, ""),
		},
	},
	{
		input: `asm { "mov A, 12" }`,
		output: []token{
			newToken(tokenTypeKeyword, "asm"),
			newToken(tokenTypeLeftBrace, "{"),
			newToken(tokenTypeString, "mov A, 12"),
			newToken(tokenTypeRightBrace, "}"),
			newToken(tokenTypeEndOfFile, ""),
		},
	},
	{
		input: `asm { "mov A, 12
push A
" }`,
		output: []token{
			newToken(tokenTypeKeyword, "asm"),
			newToken(tokenTypeLeftBrace, "{"),
			newToken(tokenTypeString, "mov A, 12\npush A\n"),
			newToken(tokenTypeRightBrace, "}"),
			newToken(tokenTypeEndOfFile, ""),
		},
	},
}

func TestLexer(t *testing.T) {

	for _, tc := range lexerTestCases {
		input := tc.input
		expectedOut := tc.output

		l := newLexer()

		out, err := l.Run("", input)
		if err != nil {
			t.Error(err)
			return
		}

		lenOut := len(out)
		lenExpected := len(expectedOut)

		if lenOut != lenExpected {
			t.Errorf("expected len %d and got %d", lenExpected, lenOut)
			return
		}

		for i, expected := range expectedOut {
			token := out[i]

			if expected.tokenType != token.tokenType {
				t.Errorf("expected type %v and got %v, at index %d", expected.tokenType, token.tokenType, i)
			}

			if expected.value != token.value {
				t.Errorf("expected value %s and got %s, for token type %v, at index %d", expected.value, token.value, token.tokenType, i)
			}
		}
	}
}
