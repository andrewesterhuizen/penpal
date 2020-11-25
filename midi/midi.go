package midi

import (
	"fmt"
	"log"

	"github.com/rakyll/portmidi"
)

type Device struct {
	Name string
	Id   int
}

type MidiHandler interface {
	Send(status byte, data1 byte, data2 byte)
	Close()
	GetDevices() (inputs []Device, ouputs []Device)
}

type PortMidiMidiHandler struct {
	midi             *portmidi.Stream
	bpm              int
	ppqn             int
	clockRunning     bool
	getMidiClockData *func() (uint8, uint8)
	tick             *func()
}

func NewPortMidiMidiHandler() MidiHandler {
	portmidi.Initialize()

	out, err := portmidi.NewOutputStream(1, 1024, 0)
	if err != nil {
		log.Fatal(err)
	}

	return &PortMidiMidiHandler{midi: out}
}

func (m *PortMidiMidiHandler) GetDevices() (inputs []Device, outputs []Device) {
	n := portmidi.CountDevices()

	inputs = []Device{}
	outputs = []Device{}

	for i := 0; i < n; i++ {
		device := portmidi.Info(portmidi.DeviceID(i))
		if device.IsInputAvailable {
			inputs = append(inputs, Device{Id: i, Name: device.Name})
		}
		if device.IsOutputAvailable {
			outputs = append(outputs, Device{Id: i, Name: device.Name})
		}
	}

	return inputs, outputs
}

func (m *PortMidiMidiHandler) Send(status byte, data1 byte, data2 byte) {
	m.midi.WriteShort(int64(status), int64(data1), int64(data2))
	fmt.Printf("SEND %02x|%02x|%02x\n", status, data1, data2)
}

func (m *PortMidiMidiHandler) Close() {
	m.midi.Close()
}
