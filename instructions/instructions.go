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
	GT
	GTE
	LT
	LTE
	Eq
	Neq
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

	Immediate                 = 0x0
	ImmediatePlusRegister     = 0x1
	ImmediateMinusRegister    = 0x2
	Register                  = 0x3
	FramePointerWithOffset    = 0x4
	FramePointerPlusRegister  = 0x5
	FramePointerMinusRegister = 0x6

	AddressingModeImmediate  = 0x7
	AddressingModeFPRelative = 0x8

	RegisterA = 0xa
	RegisterB = 0xb
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
	GT:     "gt",
	GTE:    "gte",
	LT:     "lt",
	LTE:    "lte",
	Eq:     "eq",
	Neq:    "neq",
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
	"gt":     GT,
	"gte":    GTE,
	"lt":     LT,
	"lte":    LTE,
	"eq":     Eq,
	"neq":    Neq,
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
	Mov:    3,
	Swap:   1,
	Halt:   1,
	Load:   6,
	Store:  6,
	Add:    1,
	Sub:    1,
	Mul:    1,
	Div:    1,
	Shl:    1,
	Shr:    1,
	And:    1,
	Or:     1,
	GT:     1,
	GTE:    1,
	LT:     1,
	LTE:    1,
	Eq:     1,
	Neq:    1,
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
