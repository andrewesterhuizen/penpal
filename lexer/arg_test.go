package lexer

import "testing"

type ArgTestCase struct {
	arg      Arg
	isDefine bool
	isLabel  bool
	value    string
}

var argTestCases = []ArgTestCase{
	ArgTestCase{newArg("0x12"), false, false, "0x12"},
	ArgTestCase{newArg("$test"), true, false, "test"},
	ArgTestCase{newArg("test"), false, true, "test"},
}

func TestArg(t *testing.T) {
	for _, tc := range argTestCases {
		a := tc.arg

		isDefine := a.IsDefine
		if isDefine != tc.isDefine {
			t.Errorf("expected IsDefine to be %v and got %v", tc.isDefine, isDefine)
		}

		isLabel := a.IsLabel
		if isLabel != tc.isLabel {
			t.Errorf("expected isLabel to be %v and got %v", tc.isLabel, isLabel)
		}

		value := a.Value
		if value != tc.value {
			t.Errorf("expected value to be %v and got %v", tc.value, value)
		}
	}
}
