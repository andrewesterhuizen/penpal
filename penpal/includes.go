package penpal

import (
	"bytes"
	"fmt"
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
// args: (status, data1, data2)
midi_send_message:
	// status
	STORE +5(fp) {{.AddressStatus | printf "0x%02x"}} 
	// data1
	STORE +6(fp) {{.AddressData1 | printf "0x%02x"}}  
	// data2
	STORE +7(fp) {{.AddressData2 | printf "0x%02x"}} 

	// set send byte
	MOV A 0x1
	STORE A {{.AddressSendMessage | printf "0x%02x"}} 

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

	fmt.Println(midiInclude)

	includes["midi"] = midiInclude

	return includes, nil
}
