package penpal

import (
	"bytes"
	"text/template"
)

// TODO: these are just random addresses, these should be changed to labels to DBs in the include
const (
	AddressBPM              = 0xaa00
	AddressPPQN             = 0xaa01
	AddressMidiMessageStart = 0xaa02
	AddressStatus           = 0xaa02
	AddressData1            = 0xaa03
	AddressData2            = 0xaa04
	AddressSendMessage      = 0xaa05
)

var midiNoteIncludeTemplateText = `
{{ range $key, $value := . }}
#define MIDI_NOTE_{{ $key }} {{ $value | printf "0x%02x" -}}
{{ end }}
`
var midiNoteIncludeTemplateData = map[string]int{
	"C1":  24,
	"C#1": 25,
	"D1":  26,
	"D#1": 27,
	"E":   28,
	"F":   29,
	"F#":  30,
	"G":   31,
	"G#":  32,
	"A1":  33,
	"A#1": 34,
	"B1":  35,
}

var midiIncludeTemplateText = `
#define MIDI_ADDRESS_BPM {{.AddressBPM | printf "0x%04x"}} 
#define MIDI_ADDRESS_PPQN {{.AddressPPQN | printf "0x%04x"}} 
#define MIDI_ADDRESS_STATUS {{.AddressStatus | printf "0x%04x"}} 
#define MIDI_ADDRESS_DATA1 {{.AddressData1 | printf "0x%04x"}} 
#define MIDI_ADDRESS_DATA2 {{.AddressData2 | printf "0x%04x"}} 
#define MIDI_ADDRESS_SEND_MESSAGE {{.AddressSendMessage | printf "0x%04x"}} 

// args: (status, data1, data2)
midi_send_message:
	// status
	STORE +7(fp) MIDI_ADDRESS_STATUS
	// data1
	STORE +8(fp) MIDI_ADDRESS_DATA1
	// data2
	STORE +9(fp) MIDI_ADDRESS_DATA2

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
	PUSH +7(fp)
	// status (0x90/note on)
	PUSH 0x90
	// number of args
	PUSH 0x3
	CALL midi_send_message

	
	// send note off

	// data2 (velocity)
	PUSH 0x63
	// data1 (note)
    PUSH +7(fp)
	// status (0x80/note off)
	PUSH 0x80
	// number of args
    PUSH 0x3
    CALL midi_send_message

    RET
`

var midiIncludeIncludeData = map[string]int{
	"AddressBPM":         AddressBPM,
	"AddressPPQN":        AddressPPQN,
	"AddressStatus":      AddressStatus,
	"AddressData1":       AddressData1,
	"AddressData2":       AddressData2,
	"AddressSendMessage": AddressSendMessage,
}

func getIncludeTemplate(name string, templateText string, data interface{}) (string, error) {
	buf := bytes.Buffer{}

	midiIncludeTemplate, err := template.New(name).Parse(templateText)
	if err != nil {
		return "", err
	}

	err = midiIncludeTemplate.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GetSystemIncludes() (map[string]string, error) {
	var includes = map[string]string{}

	midiInlude, err := getIncludeTemplate("<midi>", midiIncludeTemplateText, midiIncludeIncludeData)
	if err != nil {
		return nil, err
	}

	includes["midi"] = midiInlude

	notesInclude, err := getIncludeTemplate("<notes>", midiNoteIncludeTemplateText, midiNoteIncludeTemplateData)
	if err != nil {
		return nil, err
	}

	includes["notes"] = notesInclude

	return includes, nil
}
