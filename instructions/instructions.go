package instructions

const (
	HALT = iota
	// memory access
	LOAD
	STORE
	// maths / logic
	ADD
	SUB
	MUL
	DIV
	SHL
	SHR
	AND
	OR
	// conditional
	JUMP
	JUMPZ
	JUMPNZ
	// stack
	PUSH
	POP
	CALL
	RET
)

var Names = map[uint8]string{
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
}
