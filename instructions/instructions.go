package instructions

const (
	HALT = iota
	MOV
	SWAP
	LOAD
	STORE
	ADD
	SUB
	MUL
	DIV
	SHL
	SHR
	AND
	OR
	JUMP
	JUMPZ
	JUMPNZ
	PUSH
	POP
	CALL
	RET
	RAND
	SEND

	RegisterA = 0x0
	RegisterB = 0x1

	AddressingModeImmediate  = 0x0
	AddressingModeFPRelative = 0x1
)

var Names = map[uint8]string{
	MOV:    "MOV",
	SWAP:   "SWAP",
	HALT:   "HALT",
	LOAD:   "LOAD",
	STORE:  "STORE",
	ADD:    "ADD",
	SUB:    "SUB",
	MUL:    "MUL",
	DIV:    "DIV",
	SHL:    "SHL",
	SHR:    "SHR",
	AND:    "AND",
	OR:     "OR",
	JUMP:   "JUMP",
	JUMPZ:  "JUMPZ",
	JUMPNZ: "JUMPNZ",
	PUSH:   "PUSH",
	POP:    "POP",
	CALL:   "CALL",
	RET:    "RET",
	SEND:   "SEND",
	RAND:   "RAND",
}

var InstructionByName = map[string]uint8{
	"MOV":    MOV,
	"SWAP":   SWAP,
	"HALT":   HALT,
	"LOAD":   LOAD,
	"STORE":  STORE,
	"ADD":    ADD,
	"SUB":    SUB,
	"MUL":    MUL,
	"DIV":    DIV,
	"SHL":    SHL,
	"SHR":    SHR,
	"AND":    AND,
	"OR":     OR,
	"JUMP":   JUMP,
	"JUMPZ":  JUMPZ,
	"JUMPNZ": JUMPNZ,
	"PUSH":   PUSH,
	"POP":    POP,
	"CALL":   CALL,
	"RET":    RET,
	"RAND":   RAND,
	"SEND":   SEND,
}

var Width = map[uint8]int{
	MOV:    3,
	SWAP:   1,
	HALT:   1,
	LOAD:   3,
	STORE:  3,
	ADD:    1,
	SUB:    1,
	MUL:    1,
	DIV:    1,
	SHL:    1,
	SHR:    1,
	AND:    1,
	OR:     1,
	JUMP:   3,
	JUMPZ:  3,
	JUMPNZ: 3,
	PUSH:   1,
	POP:    1,
	CALL:   3,
	RET:    1,
	RAND:   1,
	SEND:   1,
}

func EncodeFlags(register uint8, addressingMode uint8) uint8 {
	return (register << 4) | addressingMode
}

func DecodeFlags(flags uint8) (addressingMode uint8, register uint8) {
	register = (flags & 0xf0) >> 4
	addressingMode = (flags & 0xf)
	return register, addressingMode
}
