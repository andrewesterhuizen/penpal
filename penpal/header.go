package penpal

import "bytes"

// HeaderSize is the size of the binary header in bytes
const HeaderSize = 8

// GetHeaderBytes returns bytes of the header to be added to the start of a penpal program
func GetHeaderBytes() []byte {
	buf := bytes.Buffer{}

	buf.WriteString("PENPAL") // program

	buf.WriteByte(0) // version major
	buf.WriteByte(1) // version minor

	return buf.Bytes()
}
