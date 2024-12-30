package opcodes

import (
	"bytes"
	"fmt"
)

type Opcode byte

const (
	OP_NOP Opcode = 0

	OP_ADD Opcode = 1
	OP_MUL Opcode = 2

	OP_INP Opcode = 3
	OP_OUT Opcode = 4

	OP_JNZ   Opcode = 5 // Jump Not Zero
	OP_JEZ   Opcode = 6 // Jump Equals Zero
	OP_LESS  Opcode = 7
	OP_EQUAL Opcode = 8

	OP_RBO Opcode = 9 // Relative Base Offset

	OP_HALT Opcode = 99
)

type Instruction struct {
	Opcode Opcode
	Params []Param
}

func (i Instruction) String() string {
	out := bytes.NewBuffer([]byte{})
	out.WriteString(OPCODE_NAMES[i.Opcode])
	out.WriteString("  \t")
	for _, p := range i.Params {
		tag := "?"
		switch p.ParamMode {
		case ParamModePosition:
			tag = "P"
		case ParamModeImmediate:
			tag = "I"
		case ParamModeRelative:
			tag = "R"
		}
		fmt.Fprintf(out, "%s %d\t", tag, p.Value)
	}
	return out.String()
}

type ParamMode byte

const (
	ParamModePosition  ParamMode = 0
	ParamModeImmediate ParamMode = 1
	ParamModeRelative  ParamMode = 2
)

type Param struct {
	ParamMode ParamMode
	Value     int
}

var OPCODE_PARAMS = map[Opcode]int{
	OP_NOP: 0,

	OP_ADD: 3,
	OP_MUL: 3,

	OP_INP: 1,
	OP_OUT: 1,

	OP_JNZ:   2,
	OP_JEZ:   2,
	OP_LESS:  3,
	OP_EQUAL: 3,

	OP_RBO: 1,

	OP_HALT: 0,
}

var OPCODE_NAMES = map[Opcode]string{
	OP_NOP: "NOP",

	OP_ADD: "ADD",
	OP_MUL: "MUL",

	OP_INP: "INP",
	OP_OUT: "OUT",

	OP_JNZ:   "JNZ",
	OP_JEZ:   "JEZ",
	OP_LESS:  "LESS",
	OP_EQUAL: "EQUAL",

	OP_RBO: "RBO",

	OP_HALT: "HALT",
}
