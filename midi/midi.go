package midi

import (
	"fmt"
	"log"

	"github.com/rakyll/portmidi"
)

type MidiHandler interface {
	Send(status byte, data1 byte, data2 byte)
	Close()
}

type PortMidiMidiHandler struct {
	midi *portmidi.Stream
}

func NewPortMidiMidiHandler() MidiHandler {
	portmidi.Initialize()

	out, err := portmidi.NewOutputStream(1, 1024, 0)
	if err != nil {
		log.Fatal(err)
	}

	return &PortMidiMidiHandler{midi: out}
}

func (m *PortMidiMidiHandler) Send(status byte, data1 byte, data2 byte) {
	m.midi.WriteShort(int64(status), int64(data1), int64(data2))
	fmt.Printf("SEND %02x|%02x|%02x\n", status, data1, data2)
}

func (m *PortMidiMidiHandler) Close() {
	m.midi.Close()
}
