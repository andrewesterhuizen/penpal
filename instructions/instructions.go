package instructions

const (
	HALT = iota
	MOVA
	MOVB
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
)

var Names = map[uint8]string{
	MOVA:   "MOVA",
	MOVB:   "MOVB",
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
}
