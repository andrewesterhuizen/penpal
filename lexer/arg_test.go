package lexer

import "testing"

type ArgTestCase struct {
	arg               Arg
	isDefine          bool
	isLabel           bool
	isFPOffsetAddress bool
	value             string
}

var argTestCases = []ArgTestCase{
	{newArg("0x12"), false, false, false, "0x12"},
	{newArg("$test"), true, false, false, "test"},
	{newArg("test"), false, true, false, "test"},
	{newArg("+1(fp)"), false, false, true, "0x1"},
	{newArg("+9(fp)"), false, false, true, "0x9"},
	{newArg("+10(fp)"), false, false, true, "0xa"},
	{newArg("+15(fp)"), false, false, true, "0xf"},
	{newArg("+16(fp)"), false, false, true, "0x10"},
	{newArg("-1(fp)"), false, false, true, "0xff"},
	{newArg("-9(fp)"), false, false, true, "0xf7"},
	{newArg("-10(fp)"), false, false, true, "0xf6"},
	{newArg("-15(fp)"), false, false, true, "0xf1"},
	{newArg("-16(fp)"), false, false, true, "0xf0"},
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

		isFPOffsetAddress := a.IsFPOffsetAddress
		if isFPOffsetAddress != tc.isFPOffsetAddress {
			t.Errorf("expected isFPOffsetAddress to be %v and got %v", tc.isFPOffsetAddress, isFPOffsetAddress)
		}

		value := a.Value
		if value != tc.value {
			t.Errorf("expected value to be %v and got %v", tc.value, value)
		}
	}
}
