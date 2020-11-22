package main

import (
	"github.com/andrewesterhuizen/penpal/assembler"
	"github.com/andrewesterhuizen/penpal/vm"
)

func main() {
	a := assembler.New()

	source := `
	MOV A 0x4
	PUSH
	MOV A 0x1
	PUSH
	CALL square 
	HALT

	square:
		MOV A +5(fp)
		MOV B +5(fp)
		MUL
		RET
	`

	i := a.GetInstructions(source)

	vm := vm.New()

	vm.Load(i)
	vm.Run()
	vm.PrintReg()
	vm.PrintMem(0, 0xf)
}
