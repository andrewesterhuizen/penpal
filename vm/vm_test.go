package vm

// import (
// 	"fmt"
// 	"testing"
// )

// func TestVM_pushpop(t *testing.T) {
// 	vm := New()

// 	value := byte(0xab)

// 	vm.push(value)

// 	poppedValue := vm.pop()

// 	if value != poppedValue {
// 		t.Errorf("expected popped value to be %04x, got %04x", value, poppedValue)
// 	}
// }

// func TestVM_push16pop16(t *testing.T) {
// 	vm := New()

// 	value := uint16(0xabcd)

// 	vm.push16(value)

// 	poppedValue := vm.pop16()

// 	if value != poppedValue {
// 		t.Errorf("expected popped value to be %04x, got %04x", value, poppedValue)
// 	}
// }

// func uint16hl(n uint16) (byte, byte) {
// 	h := (n & 0xff00) >> 8
// 	l := n & 0xff

// 	return byte(h), byte(l)
// }

// func TestVM_StateSaveRestore(t *testing.T) {
// 	vm := New()

// 	initialSP := uint16(0xffff - 1)
// 	initialBRegister := byte(0x56)
// 	initialFP := uint16(0x1234)
// 	initialIP := uint16(0xabcd)

// 	// set VM state
// 	vm.b = initialBRegister
// 	vm.ip = initialIP
// 	vm.sp = initialSP
// 	vm.fp = initialFP

// 	// save state to stack
// 	vm.saveState(false)

// 	iph, ipl := uint16hl(initialIP)
// 	fph, fpl := uint16hl(initialFP)

// 	// assert values are saved on stack
// 	if vm.memory[initialSP] != vm.b {
// 		t.Errorf("VM did not save B register")
// 	}

// 	if vm.memory[initialSP-1] != fpl {
// 		t.Errorf("VM did not save fp low byte")
// 	}

// 	if vm.memory[initialSP-2] != fph {
// 		t.Errorf("VM did not save fp high byte")
// 	}

// 	if vm.memory[initialSP-3] != ipl {
// 		t.Errorf("VM did not save ip low byte")
// 	}

// 	if vm.memory[initialSP-4] != iph {
// 		t.Errorf("VM did not save ip high byte")
// 	}

// 	// restore state from stack
// 	vm.restoreState(false)

// 	// assert state has been restored
// 	if vm.b != initialBRegister {
// 		t.Errorf("VM did not restore B register")
// 	}

// 	if vm.ip != initialIP {
// 		t.Errorf("expected ip to be 0x%04x and got 0x%04x", initialIP, vm.ip)
// 	}

// 	if vm.sp != initialSP {
// 		t.Errorf("expected sp to be 0x%04x and got 0x%04x", initialSP, vm.sp)
// 	}
// }

// func TestVM_ret_RemovesArgsFromStack(t *testing.T) {
// 	vm := New()

// 	initialStackPointer := vm.sp

// 	arg := byte(0xbb)

// 	vm.push(arg) // push arg
// 	vm.push(arg) // push arg
// 	vm.push(arg) // push arg
// 	vm.push(3)   // number of args

// 	// save state to stack
// 	vm.saveState(false)

// 	// restore state and remove args from stack
// 	vm.ret()

// 	if vm.sp != initialStackPointer {
// 		t.Errorf("VM did not restore stack pointer")
// 	}
// }

// func TestVM_ret_RemovesArgsFromStackNested(t *testing.T) {
// 	vm := New()

// 	initialStackPointer := vm.sp

// 	// sub 1
// 	fmt.Println("sub 1")

// 	vm.ip = 0xabcd
// 	vm.b = 0xbb

// 	// push args and number of args
// 	vm.push(0x12)
// 	vm.push(1)

// 	// save state to stack
// 	vm.saveState(false)

// 	// sub 2
// 	fmt.Println("sub 2")

// 	vm.ip = 0xef56
// 	vm.b = 0xcc

// 	// push args and number of args
// 	vm.push(0x34)
// 	vm.push(1)

// 	// save state to stack
// 	vm.saveState(false)

// 	// restore state from sub 1
// 	vm.ret()

// 	// restore state from sub 2
// 	vm.ret()

// 	if vm.sp != initialStackPointer {
// 		t.Errorf("expected sp to be 0x%04x and got 0x%04x", initialStackPointer, vm.sp)
// 	}
// }
