package vm

import (
	"fmt"
	"log"

	"a.com/mvm/instructions"
)

type VM struct {
	ip           uint16
	sp           uint16
	fp           uint16
	memory       [0xf]uint8
	instructions []uint8
}

func New() VM {
	vm := VM{}
	vm.ip = 0
	vm.sp = 0xf - 1
	vm.fp = 0xf - 1

	vm.memory[0] = 0xae

	return vm
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
	return vm.instructions[vm.ip]
}

func (vm *VM) fetch16() uint16 {
	vm.ip++
	h := uint16(vm.instructions[vm.ip])
	vm.ip++
	l := uint16(vm.instructions[vm.ip])
	return (h << 8) | l
}

func (vm *VM) execute(instruction uint8) {
	name, exists := instructions.Names[instruction]
	if !exists {
		log.Fatalf("encountered unknown instruction %d", instruction)
	}

	fmt.Printf("Executing %s\n", name)

	switch instruction {
	case instructions.STORE:
		addr := vm.pop16()
		val := vm.pop()
		vm.memory[addr] = val
		vm.ip++

	case instructions.LOAD:
		addr := vm.pop16()
		vm.push(vm.memory[addr])
		vm.ip++

	case instructions.ADD:
		vm.push(vm.pop() + vm.pop())
		vm.ip++

	case instructions.SUB:
		vm.push(vm.pop() - vm.pop())
		vm.ip++

	case instructions.MUL:
		vm.push(vm.pop() * vm.pop())
		vm.ip++

	case instructions.DIV:
		vm.push(vm.pop() / vm.pop())
		vm.ip++

	case instructions.SHL:
		vm.push(vm.pop() << vm.pop())
		vm.ip++

	case instructions.SHR:
		vm.push(vm.pop() >> vm.pop())
		vm.ip++

	case instructions.AND:
		vm.push(vm.pop() & vm.pop())
		vm.ip++

	case instructions.OR:
		vm.push(vm.pop() | vm.pop())
		vm.ip++

	case instructions.JUMP:
		addr := vm.fetch16()
		vm.ip = addr

	case instructions.JUMPZ:
		addr := vm.fetch16()
		n := vm.pop()
		if n == 0 {
			vm.ip = addr
		} else {
			vm.ip++
		}

	case instructions.JUMPNZ:
		addr := vm.fetch16()
		if vm.pop() != 0 {
			vm.ip = addr
		} else {
			vm.ip++
		}

	case instructions.PUSH:
		vm.push(vm.fetch())
		vm.ip++

	case instructions.POP:
		vm.pop()
		vm.ip++

	case instructions.CALL:
		addr := vm.fetch16()
		vm.push16(vm.ip + 1)
		vm.fp = vm.sp
		vm.ip = addr

	case instructions.RET:
		vm.sp = vm.fp
		addr := vm.pop16()
		vm.ip = addr

	default:
		log.Fatalf("encountered unknown instruction %d, name=%s", instruction, instructions.Names[instruction])
	}

}

func (vm *VM) Load(instructions []uint8) {
	vm.instructions = instructions
}

func (vm *VM) Run() uint8 {
	ins := vm.instructions

	for {
		if ins[vm.ip] == instructions.HALT {
			break
		}

		vm.execute(ins[vm.ip])
	}

	return vm.memory[0xf-1]
}
