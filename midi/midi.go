package midi

import (
	"fmt"
	"log"
	"time"

	"github.com/rakyll/portmidi"
)

type MidiHandler interface {
	Send(status byte, data1 byte, data2 byte)
	Close()
	StartClock()
	OnRequestMidiClockData(getMidiClockData func() (uint8, uint8))
	OnTick(onTick func())
}

type PortMidiMidiHandler struct {
	midi             *portmidi.Stream
	bpm              int
	ppqn             int
	clockRunning     bool
	getMidiClockData func() (uint8, uint8)
	tick             func()
}

func NewPortMidiMidiHandler() MidiHandler {
	portmidi.Initialize()

	out, err := portmidi.NewOutputStream(1, 1024, 0)
	if err != nil {
		log.Fatal(err)
	}

	msPerMinute := 60 * 1000

	handler := &PortMidiMidiHandler{midi: out}

	// this will need to work with an external clock at some point
	go func() {
		for {
			if !handler.clockRunning {
				continue
			}

			bpm, ppqn := handler.getMidiClockData()
			interval := (msPerMinute / int(bpm)) / int(ppqn)
			handler.tick()
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}
	}()

	return handler
}

func (m *PortMidiMidiHandler) Send(status byte, data1 byte, data2 byte) {
	m.midi.WriteShort(int64(status), int64(data1), int64(data2))
	fmt.Printf("SEND %02x|%02x|%02x\n", status, data1, data2)
}

func (m *PortMidiMidiHandler) Close() {
	m.midi.Close()
}

func (m *PortMidiMidiHandler) StartClock() {
	m.clockRunning = true
}

func (m *PortMidiMidiHandler) OnRequestMidiClockData(getMidiClockData func() (uint8, uint8)) {
	m.getMidiClockData = getMidiClockData
}

func (m *PortMidiMidiHandler) OnTick(onTick func()) {
	m.tick = onTick
}
