package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/andrewesterhuizen/penpal/assembler"
	"github.com/andrewesterhuizen/penpal/midi"
	"github.com/andrewesterhuizen/penpal/vm"
)

func main() {
	midiHandler := midi.NewPortMidiMidiHandler()
	source := ""

	args := os.Args[1:]
	if len(args) > 0 {
		arg0 := args[0]
		switch arg0 {
		case "devices":
			inputs, outputs := midiHandler.GetDevices()

			fmt.Println("inputs:")
			for _, d := range inputs {
				fmt.Printf("[%v] %s\n", d.Id, d.Name)
			}

			fmt.Println("outputs:")
			for _, d := range outputs {
				fmt.Printf("[%v] %s\n", d.Id, d.Name)
			}

			return

		case "compile":
			if len(args) < 2 {
				log.Fatal("no input file")
			}

			file := args[1]

			f, err := ioutil.ReadFile(file)
			if err != nil {
				log.Fatal(err)
			}

			a := assembler.New(assembler.Config{})

			instructions, err := a.GetInstructions(string(f))
			if err != nil {
				log.Fatal(err)
			}

			for _, i := range instructions {
				fmt.Printf("%c", i)
			}

			return

		default:
			// TODO: handle binary files too

			f, err := ioutil.ReadFile(arg0)
			if err != nil {
				log.Fatal(err)
			}

			source = string(f)
		}

	}

	a := assembler.New(assembler.Config{})

	i, err := a.GetInstructions(source)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(i)

	vm := vm.New(midiHandler)

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
