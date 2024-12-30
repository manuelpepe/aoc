package vm

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"syscall"

	"2019/utils/machine/opcodes"
)

type VM struct {
	memory     []int
	intrsPoint int

	halted bool

	// isStdin is used to show the '>> ' legend when waiting for input
	isStdin bool

	inReader  *bufio.Reader
	outWriter *bufio.Writer
}

const MAX_MEM = 2048

func NewVM(program []int, in io.Reader, out io.Writer) *VM {
	if len(program) > MAX_MEM {
		panic("memory limit exceeded")
	}

	mem := make([]int, MAX_MEM)
	copy(mem, program)

	fileIn, isFile := in.(*os.File)
	isStdin := isFile && fileIn.Fd() == uintptr(syscall.Stdin)

	return &VM{
		memory:     mem,
		intrsPoint: 0,

		halted: false,

		isStdin: isStdin,

		inReader:  bufio.NewReader(in),
		outWriter: bufio.NewWriter(out),
	}
}

func (m *VM) Halted() bool {
	return m.halted
}

func (m *VM) Run() int {
	if m.halted {
		panic("executing on halted vm")
	}

	for {
		m.evalNext()
		if m.halted {
			break
		}
	}

	return m.memory[0]
}

func (m *VM) RunForOutput() {
	if m.halted {
		panic("executing on halted vm")
	}

	for {
		ret := m.evalNext()
		if ret || m.halted {
			break
		}
	}
}

func (m *VM) evalNext() bool {
	instr := m.curInstruction()

	outputted := false

	switch instr.Opcode {
	case opcodes.OP_HALT:
		m.halted = true

	case opcodes.OP_ADD:
		out := instr.Params[2].Value
		m.memory[out] = m.getParamValue(instr.Params[0]) + m.getParamValue(instr.Params[1])

	case opcodes.OP_MUL:
		out := instr.Params[2].Value
		m.memory[out] = m.getParamValue(instr.Params[0]) * m.getParamValue(instr.Params[1])

	case opcodes.OP_SET:
		if m.isStdin {
			fmt.Fprint(m.outWriter, ">> ")
			m.outWriter.Flush()
		}

		raw, err := m.inReader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		raw, _ = strings.CutSuffix(raw, "\n")
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			panic(err)
		}

		m.memory[instr.Params[0].Value] = int(v)

	case opcodes.OP_OUT:
		fmt.Fprintf(m.outWriter, "OUT: %d\n", m.getParamValue(instr.Params[0]))
		m.outWriter.Flush()
		outputted = true

	case opcodes.OP_JOT:
		if m.getParamValue(instr.Params[0]) != 0 {
			m.intrsPoint = m.getParamValue(instr.Params[1])
			return false // skip adv instr pointer
		}

	case opcodes.OP_JOF:
		if m.getParamValue(instr.Params[0]) == 0 {
			m.intrsPoint = m.getParamValue(instr.Params[1])
			return false // skip adv instr pointer
		}

	case opcodes.OP_LESS:
		if m.getParamValue(instr.Params[0]) < m.getParamValue(instr.Params[1]) {
			m.memory[instr.Params[2].Value] = 1
		} else {
			m.memory[instr.Params[2].Value] = 0
		}

	case opcodes.OP_EQUAL:
		if m.getParamValue(instr.Params[0]) == m.getParamValue(instr.Params[1]) {
			m.memory[instr.Params[2].Value] = 1
		} else {
			m.memory[instr.Params[2].Value] = 0
		}

	default:
		panic(fmt.Sprintf("unexpected opcode at eval: %b", instr.Opcode))
	}

	m.advanceInstrPointer(len(instr.Params))

	return outputted
}

func (m *VM) curInstruction() opcodes.Instruction {
	return opcodes.Instruction{
		Opcode: m.curOpcode(),
		Params: m.curParams(),
	}
}

func (m *VM) curOpcode() opcodes.Opcode {
	opcode := m.memory[m.intrsPoint] % 100
	return opcodes.Opcode(opcode)
}

func (m *VM) curParams() []opcodes.Param {
	opcode := m.curOpcode()

	params := make([]opcodes.Param, opcodes.OPCODE_PARAMS[opcode])
	for i := 0; i < opcodes.OPCODE_PARAMS[opcode]; i++ {
		params[i] = opcodes.Param{Value: m.operand(i + 1)}
	}

	modes := parseModes(m.memory[m.intrsPoint])
	for ix, mode := range modes {
		params[ix].ParamMode = opcodes.ParamMode(mode)
	}

	return params
}

func parseModes(encoded int) []byte {
	modesRaw := int(encoded / 100)
	modesCount := int(math.Log10(float64(modesRaw))) + 1

	out := make([]byte, 0)

	for i := 0; i < modesCount; i++ {
		div := int(math.Pow(10, float64(i+1)))
		mode := modesRaw % div / (div / 10)
		out = append(out, byte(mode))
	}

	return out
}

func (m *VM) getParamValue(param opcodes.Param) int {
	switch param.ParamMode {
	case opcodes.ParamModePosition:
		return m.memory[param.Value]

	case opcodes.ParamModeImmediate:
		return param.Value

	default:
		panic(fmt.Sprintf("unexpected param mode: %b", param.ParamMode))
	}
}

func (m *VM) operand(n int) int {
	return m.memory[m.intrsPoint+n]
}

func (m *VM) advanceInstrPointer(n int) {
	m.intrsPoint += n + 1
}
