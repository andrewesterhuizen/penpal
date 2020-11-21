package main

import (
	"github.com/andrewesterhuizen/vm/assembler"
	"github.com/andrewesterhuizen/vm/vm"
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
