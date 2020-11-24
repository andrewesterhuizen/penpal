package main

import (
	"github.com/andrewesterhuizen/penpal/assembler"
	"github.com/andrewesterhuizen/penpal/midi"
	"github.com/andrewesterhuizen/penpal/vm"
)

func main() {
	a := assembler.New()

	source := `
	MOV A 0x40
	PUSH
	MOV A 0x1
	PUSH
	CALL trig

	HALT

	trig:
		// note on
		MOV A 0x63
		PUSH
		MOV A +5(fp)
		PUSH
		MOV A 0x90
		PUSH
		MOV A 0x3
		PUSH
		CALL send_midi

		// note off
		MOV A 0x63
		PUSH
		MOV A +5(fp)
		PUSH
		MOV A 0x80
		PUSH
		MOV A 0x3
		PUSH
		CALL send_midi

		RET
	
	send_midi:
		// status
		MOV A +5(fp) 
		STORE 0x0 
	
		// data1
		MOV A +6(fp) 
		STORE 0x1 
	
		// data2
		MOV A +7(fp) 
		STORE 0x2
	
		SEND
	
		RET
	`

	i := a.GetInstructions(source)

	vm := vm.New(midi.NewPortMidiMidiHandler())

	vm.Load(i)
	vm.Run()
	vm.PrintReg()
	vm.PrintMem(0, 0xf)
	vm.Close()
}
