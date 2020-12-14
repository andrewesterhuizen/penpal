package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/andrewesterhuizen/penpal/assembler"
	"github.com/andrewesterhuizen/penpal/instructions"
	"github.com/andrewesterhuizen/penpal/penpal"

	"github.com/andrewesterhuizen/penpal/midi"
	"github.com/andrewesterhuizen/penpal/vm"
)

func printDisasm(program []byte) {
	w := 0

	for i, b := range program {
		if i < 19 {
			fmt.Printf("%02d: %03d\n", i, b)
			continue
		}

		if w == 0 {
			n, exists := instructions.Names[b]
			if exists {
				fmt.Printf("%02d: (%s)\n", i, n)
				w = instructions.Width[b]
			} else {
				fmt.Printf("%02d: %03d\n", i, b)
			}
		} else {
			fmt.Printf("%02d: %03d\n", i, b)
		}

		if w > 0 {
			w--
		}
	}

	fmt.Println()
}

func printMidiDevices() {
	midiHandler := midi.NewPortMidiMidiHandler()
	inputs, outputs := midiHandler.GetDevices()

	fmt.Println("inputs:")
	for _, d := range inputs {
		fmt.Printf("[%v] %s\n", d.Id, d.Name)
	}

	fmt.Println("outputs:")
	for _, d := range outputs {
		fmt.Printf("[%v] %s\n", d.Id, d.Name)
	}
}

func compileFromFile(filename string) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	systemIncludes, err := penpal.GetSystemIncludes()
	if err != nil {
		log.Fatal(err)
	}

	a := assembler.New(assembler.Config{SystemIncludes: systemIncludes})

	program, err := a.GetProgram(filename, string(f))
	if err != nil {
		log.Fatal(err)
	}

	header := penpal.GetHeaderBytes()

	binary.Write(os.Stdout, binary.LittleEndian, header)
	binary.Write(os.Stdout, binary.LittleEndian, program)
}

func loadProgramFromFile(filename string) []byte {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	// determine if file is compiled binary by checking header
	header := []byte("PENPAL")
	binary := true

	for i, c := range header {
		if f[i] != c {
			binary = false
			break
		}
	}

	if binary {
		return f[penpal.HeaderSize:]
	}

	systemIncludes, err := penpal.GetSystemIncludes()
	if err != nil {
		log.Fatal(err)
	}

	a := assembler.New(assembler.Config{
		SystemIncludes: systemIncludes,
		InteruptLabels: [3]string{"on_tick"},
	})

	program, err := a.GetProgram(filename, string(f))
	if err != nil {
		log.Fatal(err)
	}

	return program
}

func executeProgramFromFile(filename string) {
	program := loadProgramFromFile(filename)

	midiHandler := midi.NewPortMidiMidiHandler()
	defer midiHandler.Close()

	vm := vm.New()

	msPerMinute := 60 * 1000

	// TODO: clock should be enabled according to a flag
	go func() {
		for {
			bpm := vm.GetMemory(0xd)
			ppqn := vm.GetMemory(0xe)

			if bpm == 0 || ppqn == 0 {
				continue
			}

			interval := (msPerMinute / int(bpm)) / int(ppqn)
			vm.Interupt(0)
			time.Sleep(time.Duration(interval) * time.Millisecond)

		}
	}()

	vm.Load(program)

	clockSpeedMHz := 1
	clockInterval := (1000 / time.Duration(clockSpeedMHz)) * time.Nanosecond

	ticker := time.NewTicker(clockInterval)
	defer ticker.Stop()

	done := make(chan bool)
	messages := make(chan midi.MidiMessage)

	go func() {
		for range ticker.C {
			if vm.Halted {
				vm.PrintReg()
				vm.PrintMem(0, 24)
				done <- true
				return
			}

			vm.Tick()
			m := vm.GetMemorySection(0x000f, 4)

			if m[3] > 0 {
				messages <- midi.MidiMessage{m[0], m[1], m[2]}
			}
		}
	}()

	go func() {
		for m := range messages {
			midiHandler.Send(m[0], m[1], m[2])
			vm.SetMemory(18, 0x0)
		}
	}()

	<-done
}

func main() {
	args := os.Args[1:]
	if len(args) > 0 {
		switch args[0] {
		case "devices":
			printMidiDevices()
			return

		case "compile":
			if len(args) < 2 {
				log.Fatal("no input file")
			}

			compileFromFile(args[1])
			return

		default:
			executeProgramFromFile(args[0])
		}

		return
	}

	// TODO: print help info if no args supplied
}
