package penpal

import (
	"bytes"
	"html/template"
)

const (
	AddressBPM              = 0x0
	AddressPPQN             = 0x1
	AddressMidiMessageStart = 0x2
	AddressStatus           = 0x2
	AddressData1            = 0x3
	AddressData2            = 0x4
	AddressSendMessage      = 0x5
)

var midiIncludeTemplateText = `
#define MIDI_ADDRESS_BPM {{.AddressBPM | printf "0x%02x"}} 
#define MIDI_ADDRESS_PPQN {{.AddressPPQN | printf "0x%02x"}} 
#define MIDI_ADDRESS_STATUS {{.AddressStatus | printf "0x%02x"}} 
#define MIDI_ADDRESS_DATA1 {{.AddressData1 | printf "0x%02x"}} 
#define MIDI_ADDRESS_DATA2 {{.AddressData2 | printf "0x%02x"}} 
#define MIDI_ADDRESS_SEND_MESSAGE {{.AddressSendMessage | printf "0x%02x"}} 

// args: (status, data1, data2)
midi_send_message:
	// status
	STORE +5(fp) MIDI_ADDRESS_STATUS
	// data1
	STORE +6(fp) MIDI_ADDRESS_DATA1
	// data2
	STORE +7(fp) MIDI_ADDRESS_DATA2

	// set send byte
	MOV A 0x1
	STORE A MIDI_ADDRESS_SEND_MESSAGE

	RET

// args: (note)
midi_trig:
	// send note on

	// data2 (velocity)
	PUSH 0x7F
	// data1 (note)
	PUSH +5(fp)
	// status (0x90/note on)
	PUSH 0x90
	// number of args
	PUSH 0x3
	CALL midi_send_message

	
	// send note off

	// data2 (velocity)
	PUSH 0x63
	// data1 (note)
    PUSH +5(fp)
	// status (0x80/note off)
	PUSH 0x80
	// number of args
    PUSH 0x3
    CALL midi_send_message

    RET
`

func GetSystemIncludes() (map[string]string, error) {
	var includes = map[string]string{}

	buf := bytes.Buffer{}

	midiTemplate, err := template.New("midi_include").Parse(midiIncludeTemplateText)
	if err != nil {
		return nil, err
	}

	data := map[string]int{
		"AddressBPM":         AddressBPM,
		"AddressPPQN":        AddressPPQN,
		"AddressStatus":      AddressStatus,
		"AddressData1":       AddressData1,
		"AddressData2":       AddressData2,
		"AddressSendMessage": AddressSendMessage,
	}
	err = midiTemplate.Execute(&buf, data)
	if err != nil {
		return nil, err
	}

	midiInclude := buf.String()

	includes["midi"] = midiInclude

	return includes, nil
}
