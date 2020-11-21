package main

import (
	"github.com/andrewesterhuizen/penpal/assembler"
	"github.com/andrewesterhuizen/penpal/vm"
)

func main() {
	a := assembler.New()

	source := `
		CALL 0x6
		MOVB 0xaa
		HALT
		
		MOVA 0xff
		RET
	
		HALT
	`

	i := a.GetInstructions(source)

	vm := vm.New()

	vm.Load(i)
	vm.Run()
	vm.PrintReg()
	vm.PrintMem(0xffff-10, 10)
}
