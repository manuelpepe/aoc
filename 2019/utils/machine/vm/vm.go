package vm

import (
	"bufio"
	"bytes"
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
	id         string
	debugLevel int

	memory     []int
	intrsPoint int
	lastOpcode opcodes.Opcode

	// for reload
	originalMemory []int

	halted bool

	// isStdin is used to show the '>> ' legend when waiting for input
	isStdin bool

	ioMode IOMode

	inReader  *bufio.Reader
	outWriter *bufio.Writer

	relativeBase int
}

const MAX_MEM = 4096 * 2

func NewVM(program []int, opts ...vmOptionFunc) *VM {
	if len(program) > MAX_MEM {
		panic("memory limit exceeded")
	}

	mem := make([]int, MAX_MEM)
	copy(mem, program)

	mem2 := make([]int, MAX_MEM)
	copy(mem2, program)

	cfg := &vmconfig{
		reader:     os.Stdin,
		writer:     os.Stdout,
		mode:       IOInteger,
		id:         "???",
		debugLevel: 0,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	fileIn, isFile := cfg.reader.(*os.File)
	isStdin := isFile && fileIn.Fd() == uintptr(syscall.Stdin)

	return &VM{
		id:         cfg.id,
		debugLevel: cfg.debugLevel,

		memory:     mem,
		intrsPoint: 0,

		originalMemory: mem2,

		halted: false,

		isStdin: isStdin,
		ioMode:  cfg.mode,

		inReader:  bufio.NewReader(cfg.reader),
		outWriter: bufio.NewWriter(cfg.writer),

		relativeBase: 0,
	}
}

func NewVMPiped(program []int, options ...vmOptionFunc) (*VM, io.Writer, *bytes.Buffer) {
	inReader, inWriter := io.Pipe()
	outbuf := bytes.NewBuffer([]byte{})
	opts := []vmOptionFunc{
		WithInput(inReader),
		WithOutput(outbuf),
	}
	opts = append(opts, options...)
	m := NewVM(program, opts...)
	return m, inWriter, outbuf
}

func (m *VM) Restart() {
	copy(m.memory, m.originalMemory)
	m.intrsPoint = 0
	m.halted = false
	m.relativeBase = 0
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
		m.evalNext()
		if m.halted || m.lastOpcode == opcodes.OP_OUT {
			break
		}
	}
}

func (m *VM) RunForInput() {
	if m.halted {
		panic("executing on halted vm")
	}

	for {
		m.evalNext()
		if m.halted || m.lastOpcode == opcodes.OP_INP {
			break
		}
	}
}

func (m *VM) evalNext() {
	instr := m.curInstruction()
	m.lastOpcode = instr.Opcode

	m.debug("INSTR: %+v", instr)

	switch instr.Opcode {
	case opcodes.OP_HALT:
		m.halted = true
		m.debug("HALTED")

	case opcodes.OP_ADD:
		out := m.getWritingParamValue(instr.Params[2])
		m.memory[out] = m.getParamValue(instr.Params[0]) + m.getParamValue(instr.Params[1])

	case opcodes.OP_MUL:
		out := m.getWritingParamValue(instr.Params[2])
		m.memory[out] = m.getParamValue(instr.Params[0]) * m.getParamValue(instr.Params[1])

	case opcodes.OP_INP:
		if m.isStdin {
			fmt.Print(">> ")
			m.outWriter.Flush()
		}

		switch m.ioMode {
		case IOInteger:
			raw, err := m.inReader.ReadString('\n')
			if err != nil {
				panic(err)
			}

			raw, _ = strings.CutSuffix(raw, "\n")
			v, err := strconv.ParseInt(raw, 10, 64)
			if err != nil {
				panic(err)
			}

			m.memory[m.getWritingParamValue(instr.Params[0])] = int(v)

		case IOAscii:
			b, err := m.inReader.ReadByte()
			if err != nil {
				panic(err)
			}

			m.memory[m.getWritingParamValue(instr.Params[0])] = int(b)
		}

		m.debug("READ: %d", m.memory[m.getWritingParamValue(instr.Params[0])])

	case opcodes.OP_OUT:
		switch m.ioMode {
		case IOInteger:
			fmt.Fprintf(m.outWriter, "OUT: %d\n", m.getParamValue(instr.Params[0]))
			m.debug("WROTE: %d", m.getParamValue(instr.Params[0]))
		case IOAscii:
			v := m.getParamValue(instr.Params[0])
			if v < 0 || v > 255 {
				fmt.Fprintf(m.outWriter, "%d", v)
			} else {
				fmt.Fprintf(m.outWriter, "%c", v)
			}
		}

		m.outWriter.Flush()

	case opcodes.OP_JNZ:
		if m.getParamValue(instr.Params[0]) != 0 {
			m.intrsPoint = m.getParamValue(instr.Params[1])
			return // skip adv instr pointer
		}

	case opcodes.OP_JEZ:
		if m.getParamValue(instr.Params[0]) == 0 {
			m.intrsPoint = m.getParamValue(instr.Params[1])
			return // skip adv instr pointer
		}

	case opcodes.OP_LESS:
		if m.getParamValue(instr.Params[0]) < m.getParamValue(instr.Params[1]) {
			m.memory[m.getWritingParamValue(instr.Params[2])] = 1
		} else {
			m.memory[m.getWritingParamValue(instr.Params[2])] = 0
		}

	case opcodes.OP_EQUAL:
		if m.getParamValue(instr.Params[0]) == m.getParamValue(instr.Params[1]) {
			m.memory[m.getWritingParamValue(instr.Params[2])] = 1
		} else {
			m.memory[m.getWritingParamValue(instr.Params[2])] = 0
		}

	case opcodes.OP_RBO:
		m.relativeBase += m.getParamValue(instr.Params[0])

	default:
		panic(fmt.Sprintf("unexpected opcode at eval: '%b' at $%d", instr.Opcode, m.intrsPoint))
	}

	m.advanceInstrPointer(len(instr.Params))

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

	case opcodes.ParamModeRelative:
		return m.memory[param.Value+m.relativeBase]

	default:
		panic(fmt.Sprintf("unexpected param mode: %b", param.ParamMode))
	}
}

func (m *VM) getWritingParamValue(param opcodes.Param) int {
	switch param.ParamMode {
	case opcodes.ParamModePosition:
		return param.Value

	case opcodes.ParamModeImmediate:
		panic("writing with immediate mode, breaks Day 5 constraint")

	case opcodes.ParamModeRelative:
		return param.Value + m.relativeBase

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

func (m *VM) debug(msg string, args ...any) {
	if m.debugLevel > 0 {
		formatted := fmt.Sprintf(msg, args...)
		fmt.Printf("%s - %s\n", m.id, formatted)
	}
}
