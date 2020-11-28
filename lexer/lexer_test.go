package lexer

import (
	"testing"
)

type TestCase struct {
	input  string
	output []Token
}

var testCases = []TestCase{
	TestCase{input: "MOV A 0xab", output: []Token{newInstructionToken("MOV", newArgs("A", "0xab"))}},
	TestCase{input: "JUMP 0xabcd", output: []Token{newInstructionToken("JUMP", newArgs("0xabcd"))}},
	TestCase{input: "#define TEST 0xbc", output: []Token{newDefineToken("TEST", newArgs("0xbc"))}},
	TestCase{input: "#include \"test.asm\"", output: []Token{newFileIncludeToken("test.asm")}},
	TestCase{input: "#include <test>", output: []Token{newSystemIncludeToken("test")}},
	TestCase{input: "test:", output: []Token{newLabelToken("test")}},
	TestCase{
		input: `
			MOVA 0xa
			MOVB 0xfa`,
		output: []Token{
			newInstructionToken("MOVA", newArgs("0xa")),
			newInstructionToken("MOVB", newArgs("0xfa")),
		},
	},
	TestCase{
		input: `
			#define TEST 0xfa
			MOVA 0xa
			MOVB TEST`,
		output: []Token{
			newDefineToken("TEST", newArgs("0xfa")),
			newInstructionToken("MOVA", newArgs("0xa")),
			newInstructionToken("MOVB", newArgs("TEST")),
		},
	},
	TestCase{
		input: `
			test:
				MOVA 0ec
				JUMP test`,
		output: []Token{
			newLabelToken("test"),
			newInstructionToken("MOVA", newArgs("0ec")),
			newInstructionToken("JUMP", newArgs("test")),
		},
	},
}

func TestLexer(t *testing.T) {
	for _, tc := range testCases {
		l := New()

		tokens, err := l.GetTokens("", tc.input)
		if err != nil {
			t.Errorf("lexer returned error: %w", err)
		}

		if len(tokens) != len(tc.output) {
			t.Errorf("expected %d instuctions and got %d", len(tc.output), len(tokens))
			return
		}

		for i, expectedToken := range tc.output {
			token := tokens[i]

			if token.Type != expectedToken.Type {
				t.Errorf("expected Type: %s and got %s", expectedToken.Type, token.Type)
			}

			if expectedToken.Value != token.Value {
				t.Errorf("expected Value: %s and got %s", expectedToken.Value, token.Value)
			}

			expectedArgsLength := len(expectedToken.Args)
			argsLength := len(token.Args)

			if argsLength != expectedArgsLength {
				t.Errorf("expected %d args and got %d", expectedArgsLength, argsLength)
				return
			}

			for j, expectedArg := range expectedToken.Args {
				arg := token.Args[j]

				if arg.Value != expectedArg.Value {
					t.Errorf("expected arg Value: %s and got %s", arg.Value, expectedArg.Value)
				}
			}
		}
	}
}
