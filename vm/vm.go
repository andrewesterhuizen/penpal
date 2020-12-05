package vm

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/andrewesterhuizen/penpal/instructions"
)

const memorySize = 0xffff

type VM struct {
	Halted bool

	ip     uint16
	sp     uint16
	fp     uint16
	a      uint8
	b      uint8
	memory [memorySize]uint8

	// TODO: make nested interupts work
	inInterupt bool
}

func New() *VM {
	rand.Seed(time.Now().UnixNano())

	vm := VM{}
	vm.init()
	return &vm
}

func (vm *VM) init() {
	vm.ip = 0
	vm.sp = memorySize - 1
	vm.fp = memorySize - 1
}

func (vm *VM) push(value uint8) {
	vm.memory[vm.sp] = value
	vm.sp--
}

func (vm *VM) pop() uint8 {
	vm.sp++
	return vm.memory[vm.sp]
}

func (vm *VM) push16(n uint16) {
	h := uint8((n & 0xff00) >> 8)
	l := uint8(n & 0xff)

	vm.push(l)
	vm.push(h)
}

func (vm *VM) pop16() uint16 {
	var h, l uint16
	h = uint16(vm.pop())
	l = uint16(vm.pop())
	return (h << 8) | l
}

func (vm *VM) fetch() uint8 {
	vm.ip++
	return vm.memory[vm.ip]
}

func (vm *VM) fetch16() uint16 {
	vm.ip++
	h := uint16(vm.memory[vm.ip])
	vm.ip++
	l := uint16(vm.memory[vm.ip])
	return (h << 8) | l
}

func (vm *VM) getFramePointerRelativeAddress(offset int8) uint16 {
	addr := vm.fp

	if offset >= 0 {
		addr += uint16(offset)
	} else {
		offset *= -1
		addr -= uint16(offset)
	}

	return addr
}

func (vm *VM) saveState(interupt bool) {
	// a register is used for return value in subroutines so we don't save it for non interupts
	if interupt {
		vm.push(vm.a)
	}

	vm.push(vm.b)
	vm.push16(vm.fp)
	vm.push16(vm.ip)
	vm.fp = vm.sp
}

func (vm *VM) restoreState(interupt bool) {
	vm.sp = vm.fp
	vm.ip = vm.pop16()
	prevfp := vm.pop16()
	vm.b = vm.pop()

	if interupt {
		vm.a = vm.pop()
	}

	vm.sp = vm.fp
	vm.fp = prevfp
}

func (vm *VM) call(addr uint16) {
	vm.saveState(false)
	vm.ip = addr
}

func (vm *VM) ret() {
	vm.restoreState(false)

	// remove args from stack
	nArgs := vm.pop()
	for i := 0; i < int(nArgs); i++ {
		vm.pop()
	}
}

func (vm *VM) callInterupt(addr uint16) {
	vm.inInterupt = true
	vm.saveState(true)
	vm.ip = addr
}

func (vm *VM) retFromInterupt() {
	vm.restoreState(true)
	vm.inInterupt = false
}

func (vm *VM) Interupt(n int) {
	if vm.inInterupt {
		return
	}

	// only 3 interupts for now
	if n >= 3 {
		return
	}

	// each jump instruction is 3 bytes wide
	// address of interupt jump location = entry point + (interupt number * 3 bytes)
	addr := uint16(3 + (n * 3))

	// if interupt has been set
	if vm.memory[addr] > 0 {
		vm.callInterupt(addr)
	}
}

func (vm *VM) execute(instruction uint8) {
	switch instruction {
	case instructions.Swap:
		vm.a, vm.b = vm.b, vm.a
		vm.ip++

	case instructions.Mov:
		mode := vm.fetch()
		register := vm.fetch()

		var dest *uint8

		switch register {
		case instructions.RegisterA:
			dest = &vm.a
		case instructions.RegisterB:
			dest = &vm.b
		default:
			log.Fatalf("encountered unknown destination for mov, 0x%02x", register)
		}

		switch mode {
		case instructions.AddressingModeImmediate:
			*dest = vm.fetch()
		case instructions.AddressingModeFPRelative:
			offset := int8(vm.fetch())
			addr := vm.getFramePointerRelativeAddress(offset)
			v := vm.memory[addr]
			*dest = v
		default:
			log.Fatalf("encountered unknown addressing mode for mov, 0x%02x", mode)
		}

		vm.ip++

	case instructions.Store:
		mode := vm.fetch()
		modeArg := vm.fetch()

		var value byte

		switch mode {
		case instructions.FramePointerWithOffset:
			addr := vm.getFramePointerRelativeAddress(int8(modeArg))
			value = vm.memory[addr]
		case instructions.Register:
			switch modeArg {
			case instructions.RegisterA:
				value = vm.a
			case instructions.RegisterB:
				value = vm.b
			default:
				log.Fatalf("store: encountered unknown register source 0x%02x", mode)
			}

		default:
			log.Fatalf("store: encountered unknown addressing mode 0x%02x", mode)
		}

		addr := vm.fetch16()
		vm.memory[addr] = value
		vm.ip++

	case instructions.Load:
		mode := vm.fetch()
		modeArg := vm.fetch()

		srcAddr := vm.fetch16()

		switch mode {
		case instructions.FramePointerWithOffset:
			addr := vm.getFramePointerRelativeAddress(int8(modeArg))
			vm.memory[addr] = vm.memory[srcAddr]
		case instructions.Register:
			switch modeArg {
			case instructions.RegisterA:
				vm.a = vm.memory[srcAddr]
			case instructions.RegisterB:
				vm.b = vm.memory[srcAddr]
			default:
				log.Fatalf("load: encountered unknown register source 0x%02x", mode)
			}

		default:
			log.Fatalf("load: encountered unknown addressing mode 0x%02x", mode)
		}

		vm.ip++

	case instructions.Add:
		vm.a += vm.b
		vm.ip++

	case instructions.Sub:
		vm.a -= vm.b
		vm.ip++

	case instructions.Mul:
		vm.a *= vm.b
		vm.ip++

	case instructions.Div:
		vm.a /= vm.b
		vm.ip++

	case instructions.Shl:
		vm.a = vm.a << vm.b
		vm.ip++

	case instructions.Shr:
		vm.a = vm.a >> vm.b
		vm.ip++

	case instructions.And:
		vm.a = vm.a & vm.b
		vm.ip++

	case instructions.Or:
		vm.a = vm.a | vm.b
		vm.ip++

	case instructions.Jump:
		addr := vm.fetch16()
		vm.ip = addr

	case instructions.Jumpz:
		addr := vm.fetch16()
		if vm.a == 0 {
			vm.ip = addr
		} else {
			vm.ip++
		}

	case instructions.Jumpnz:
		addr := vm.fetch16()
		if vm.a != 0 {
			vm.ip = addr
		} else {
			vm.ip++
		}

	case instructions.Push:
		mode := vm.fetch()
		value := vm.fetch()

		switch mode {
		case instructions.Register:
			switch value {
			case instructions.RegisterA:
				vm.push(vm.a)
			case instructions.RegisterB:
				vm.push(vm.b)
			default:
				log.Fatalf("push: encountered unknown register 0x%02x\n", value)
			}
		case instructions.FramePointerWithOffset:
			addr := vm.getFramePointerRelativeAddress(int8(value))
			vm.push(vm.memory[addr])
		case instructions.Value:
			vm.push(value)
		default:
			log.Fatalf("push: encountered unknown mode 0x%02x\n", mode)
		}

		vm.ip++

	case instructions.Pop:
		vm.a = vm.pop()
		vm.ip++

	case instructions.Call:
		addr := vm.fetch16()
		vm.call(addr)

	case instructions.Ret:
		vm.ret()
		vm.ip++

	case instructions.Reti:
		vm.retFromInterupt()

	case instructions.Rand:
		vm.a = uint8(rand.Intn(255))
		vm.ip++

	default:
		log.Fatalf("encountered unknown instruction 0x%02x, name=%s", instruction, instructions.Names[instruction])
	}

}

func (vm *VM) Load(instructions []uint8) {
	vm.init()
	copy(vm.memory[:], instructions)
}

func (vm *VM) GetMemorySection(start uint16, n uint16) []byte {
	return vm.memory[start : start+n]
}

func (vm *VM) GetMemory(addr uint16) uint8 {
	return vm.memory[addr]
}

func (vm *VM) SetMemory(addr uint16, value uint8) {
	vm.memory[addr] = value
}

func (vm *VM) PrintReg() {
	fmt.Printf("a: 0x%02x | b: 0x%02x\n", vm.a, vm.b)
}

func (vm *VM) PrintMem(start uint16, n uint16) {
	for i := start; i < start+n; i++ {
		if vm.sp == i {
			fmt.Printf("sp ->")
		} else {
			fmt.Print("     ")
		}

		fmt.Printf("%04x: 0x%02x", i, vm.memory[i])

		if vm.fp == i {
			fmt.Printf("<- fp\n")
		} else {
			fmt.Print("\n")
		}
	}
}

func (vm *VM) Tick() {
	if vm.memory[vm.ip] == instructions.Halt {
		vm.Halted = true
		return
	}

	vm.execute(vm.memory[vm.ip])
}
