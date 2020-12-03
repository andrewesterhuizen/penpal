package lexer

import (
	"testing"
)

type TestCase struct {
	input  string
	output []Token
}

var testCases = []TestCase{
	TestCase{input: "mov A 0xab", output: []Token{newInstructionToken("mov", newArgs("A", "0xab"))}},
	TestCase{input: "jump 0xabcd", output: []Token{newInstructionToken("jump", newArgs("0xabcd"))}},
	TestCase{input: "#define TEST 0xbc", output: []Token{newDefineToken("TEST", newArgs("0xbc"))}},
	TestCase{input: "#include \"test.asm\"", output: []Token{newFileIncludeToken("test.asm")}},
	TestCase{input: "#include <test>", output: []Token{newSystemIncludeToken("test")}},
	TestCase{input: "test:", output: []Token{newLabelToken("test")}},
	TestCase{input: "__test:", output: []Token{newLabelToken("__test")}},
	TestCase{
		input: `
			movA 0xa
			movB 0xfa`,
		output: []Token{
			newInstructionToken("movA", newArgs("0xa")),
			newInstructionToken("movB", newArgs("0xfa")),
		},
	},
	TestCase{
		input: `
			#define TEST 0xfa
			movA 0xa
			movB TEST`,
		output: []Token{
			newDefineToken("TEST", newArgs("0xfa")),
			newInstructionToken("movA", newArgs("0xa")),
			newInstructionToken("movB", newArgs("TEST")),
		},
	},
	TestCase{
		input: `
			test:
				movA 0ec
				jump test`,
		output: []Token{
			newLabelToken("test"),
			newInstructionToken("movA", newArgs("0ec")),
			newInstructionToken("jump", newArgs("test")),
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
