package opcodes

type Opcode byte

const (
	OP_NOP = 0

	OP_ADD = 1
	OP_MUL = 2

	OP_SET = 3
	OP_OUT = 4

	OP_JOT   = 5 // Jump On True
	OP_JOF   = 6 // Jump On False
	OP_LESS  = 7
	OP_EQUAL = 8

	OP_HALT = 99
)

type Instruction struct {
	Opcode Opcode
	Params []Param
}

type ParamMode byte

const (
	ParamModePosition  ParamMode = 0
	ParamModeImmediate ParamMode = 1
)

type Param struct {
	ParamMode ParamMode
	Value     int
}

var OPCODE_PARAMS = map[Opcode]int{
	OP_NOP: 0,

	OP_ADD: 3,
	OP_MUL: 3,

	OP_SET: 1,
	OP_OUT: 1,

	OP_JOT:   2,
	OP_JOF:   2,
	OP_LESS:  3,
	OP_EQUAL: 3,

	OP_HALT: 0,
}
