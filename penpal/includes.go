package penpal

import (
	"bytes"
	"text/template"
)

var midiNoteIncludeTemplateText = `
{{ range $key, $value := . }}
midi_note_{{ $key }}: db {{ $value -}}
{{ end }}
`
var midiNoteIncludeTemplateData = map[string]int{
	"C1":   24,
	"Csh1": 25,
	"D1":   26,
	"Dsh1": 27,
	"E":    28,
	"F":    29,
	"Fsh":  30,
	"G":    31,
	"Gsh":  32,
	"A1":   33,
	"Ash1": 34,
	"B1":   35,
}

var midiIncludeTemplateText = `
midi_clock_enable: db 1
midi_bpm: db 120
midi_ppqn: db 2
midi_status: db 0
midi_data1: db 0
midi_data2: db 0
midi_send_bit: db 0

// args: (status, data1, data2)
midi_send_message:
	load (fp+7), A
	store A, midi_status
	load (fp+8), A
	store A, midi_data1
	load (fp+9), A
	store A, midi_data2
	mov A, 1
	store A, midi_send_bit
	ret

// args: (note, velocity)
midi_note_on:
	load (fp+8), A
	push
	load (fp+7), A
	push
	push 0x90
	push 0x3
	call midi_send_message
	ret

// args: (note, velocity)
midi_note_off:
	load (fp+8), A
	push
	load (fp+7), A
	push 0x80
	push 0x3
	call midi_send_message
	ret

// args: (note)
midi_trig:
	load (fp+7), A

	push 0x7f
	push
	push 2
	call midi_note_on

	push 0x7f
	push
	push 2
	call midi_note_off

    ret
`

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

	midiInlude, err := getIncludeTemplate("<midi>", midiIncludeTemplateText, map[string]int{})
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
