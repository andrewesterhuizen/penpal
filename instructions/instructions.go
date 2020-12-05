package instructions

const (
	Halt = iota
	Mov
	Swap
	Load
	Store
	Add
	Sub
	Mul
	Div
	Shl
	Shr
	And
	Or
	Jump
	Jumpz
	Jumpnz
	Push
	Pop
	Call
	Ret
	Reti
	Rand
	Db

	Value                     = 0x0 // TODO: remove this one when it is no longer used
	Immediate                 = 0x0
	Register                  = 0x1
	FramePointerWithOffset    = 0x2
	FramePointerPlusRegister  = 0x3
	FramePointerMinusRegister = 0x4

	RegisterA = 0x0
	RegisterB = 0x1

	AddressingModeImmediate  = 0x0
	AddressingModeFPRelative = 0x1
)

var Names = map[uint8]string{
	Mov:    "mov",
	Swap:   "swap",
	Halt:   "halt",
	Load:   "load",
	Store:  "store",
	Add:    "add",
	Sub:    "sub",
	Mul:    "mul",
	Div:    "div",
	Shl:    "shl",
	Shr:    "shr",
	And:    "and",
	Or:     "or",
	Jump:   "jump",
	Jumpz:  "jumpz",
	Jumpnz: "jumpnz",
	Push:   "push",
	Pop:    "pop",
	Call:   "call",
	Ret:    "ret",
	Reti:   "reti",
	Rand:   "rand",
	Db:     "db",
}

var InstructionByName = map[string]uint8{
	"mov":    Mov,
	"swap":   Swap,
	"halt":   Halt,
	"load":   Load,
	"store":  Store,
	"add":    Add,
	"sub":    Sub,
	"mul":    Mul,
	"div":    Div,
	"shl":    Shl,
	"shr":    Shr,
	"and":    And,
	"or":     Or,
	"jump":   Jump,
	"jumpz":  Jumpz,
	"jumpnz": Jumpnz,
	"push":   Push,
	"pop":    Pop,
	"call":   Call,
	"ret":    Ret,
	"reti":   Reti,
	"rand":   Rand,
	"db":     Db,
}

var Width = map[uint8]int{
	Mov:    4,
	Swap:   1,
	Halt:   1,
	Load:   5,
	Store:  5,
	Add:    1,
	Sub:    1,
	Mul:    1,
	Div:    1,
	Shl:    1,
	Shr:    1,
	And:    1,
	Or:     1,
	Jump:   3,
	Jumpz:  3,
	Jumpnz: 3,
	Push:   3,
	Pop:    1,
	Call:   3,
	Ret:    1,
	Reti:   1,
	Rand:   1,
	Db:     1,
}

var RegistersByName = map[string]uint8{
	"A": RegisterA,
	"B": RegisterB,
}
