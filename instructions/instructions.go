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
	RETI
	RAND
	DB

	Value                       = 0x0
	Register                    = 0x1
	FramePointerRelativeAddress = 0x2

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
	RETI:   "RETI",
	RAND:   "RAND",
	DB:     "DB",
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
	"RETI":   RETI,
	"RAND":   RAND,
	"DB":     DB,
}

var Width = map[uint8]int{
	MOV:    4,
	SWAP:   1,
	HALT:   1,
	LOAD:   5,
	STORE:  5,
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
	PUSH:   3,
	POP:    1,
	CALL:   3,
	RET:    1,
	RETI:   1,
	RAND:   1,
	DB:     1,
}

var RegistersByName = map[string]uint8{
	"A": RegisterA,
	"B": RegisterB,
}
