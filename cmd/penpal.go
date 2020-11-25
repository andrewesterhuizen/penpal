package main

import (
	"fmt"
	"log"
	"time"

	"github.com/andrewesterhuizen/penpal/assembler"
	"github.com/andrewesterhuizen/penpal/midi"
	"github.com/andrewesterhuizen/penpal/vm"
)

func main() {
	a := assembler.New(assembler.Config{})

	source := `
	// this should be skipped
	HALT 

	__start:
		PUSH 0x40
		PUSH 0x1
		CALL trig

		HALT

	trig:
		// note on
		PUSH 0x63
		MOV A +5(fp)
		PUSH
		PUSH 0x90
		PUSH 0x3
		CALL send_midi

		// note off
		PUSH 0x63
		MOV A +5(fp)
		PUSH
		PUSH 0x80
		PUSH 0x3
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

	i, err := a.GetInstructions(source)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(i)

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
