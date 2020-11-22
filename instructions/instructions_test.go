package instructions

import "testing"

type FlagsTestCase struct {
	dest uint8
	mode uint8
	flag uint8
}

var flagsTestCase = []FlagsTestCase{
	{dest: DestRegisterA, mode: AddressingModeImmediate, flag: 0},
	{dest: DestRegisterB, mode: AddressingModeImmediate, flag: 0x10},
	{dest: DestRegisterA, mode: AddressingModeFPRelative, flag: 0x1},
	{dest: DestRegisterB, mode: AddressingModeFPRelative, flag: 0x11},
}

func TestEncodeFlags(t *testing.T) {
	for _, tc := range flagsTestCase {
		flag := EncodeFlags(tc.dest, tc.mode)

		if flag != tc.flag {
			t.Errorf("expected flag 0x%02x and got 0x%02x with input dest=0x%02x, mode=0x%02x", tc.flag, flag, tc.dest, tc.mode)
			return
		}

		dest, mode := DecodeFlags(flag)

		if dest != tc.dest {
			t.Errorf("expected dest 0x%02x and got 0x%02x with input 0x%02x", tc.dest, dest, tc.flag)
			return
		}

		if mode != tc.mode {
			t.Errorf("expected mode 0x%02x and got 0x%02x with input 0x%02x", tc.mode, mode, tc.flag)
			return
		}

	}
}
