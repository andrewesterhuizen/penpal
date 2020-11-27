package penpal

var SystemIncludes = map[string]string{
	"midi": `
		send_midi:
		// status
		STORE +5(fp) 0x0 
		// data1
		STORE +6(fp) 0x1 
		// data2
		STORE +7(fp) 0x2

		SEND
		RET
`,
}
