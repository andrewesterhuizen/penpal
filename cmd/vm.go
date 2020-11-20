package main

import (
	"fmt"

	"a.com/mvm/assembler"
	"a.com/mvm/vm"
)

func main() {
	a := assembler.New()

	source := `
	CALL 0x6
	PUSH 0xae
	HALT

	PUSH 0x1
	RET
	`

	i := a.GetInstructions(source)

	vm := vm.New()

	vm.Load(i)
	f := vm.Run()
	fmt.Printf("finished with 0x%02x\n", f)
}
