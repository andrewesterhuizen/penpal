package main

import (
	"fmt"
	"time"

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

	msPerMinute := 60 * 1000

	// TODO: clock should be enabled according to a flag
	go func() {
		for {
			bpm, ppqn := vm.GetMidiClockData()
			fmt.Println(bpm)
			if bpm == 0 || ppqn == 0 {
				time.Sleep(10 * time.Millisecond)
				continue
			}

			interval := (msPerMinute / int(bpm)) / int(ppqn)
			vm.Tick()
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}()

	vm.Load(i)
	vm.Run()
	vm.PrintReg()
	vm.PrintMem(0, 0xf)
	vm.Close()
}
