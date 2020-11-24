package instructions

import "testing"

type FlagsTestCase struct {
	register uint8
	mode     uint8
	flag     uint8
}

var flagsTestCase = []FlagsTestCase{
	{register: RegisterA, mode: AddressingModeImmediate, flag: 0},
	{register: RegisterB, mode: AddressingModeImmediate, flag: 0x10},
	{register: RegisterA, mode: AddressingModeFPRelative, flag: 0x1},
	{register: RegisterB, mode: AddressingModeFPRelative, flag: 0x11},
}

func TestEncodeFlags(t *testing.T) {
	for _, tc := range flagsTestCase {
		flag := EncodeFlags(tc.register, tc.mode)

		if flag != tc.flag {
			t.Errorf("expected flag 0x%02x and got 0x%02x with input register=0x%02x, mode=0x%02x", tc.flag, flag, tc.register, tc.mode)
			return
		}

		register, mode := DecodeFlags(flag)

		if register != tc.register {
			t.Errorf("expected register 0x%02x and got 0x%02x with input 0x%02x", tc.register, register, tc.flag)
			return
		}

		if mode != tc.mode {
			t.Errorf("expected mode 0x%02x and got 0x%02x with input 0x%02x", tc.mode, mode, tc.flag)
			return
		}

	}
}
