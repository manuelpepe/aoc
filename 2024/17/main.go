package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

type OPCODE int

const (
	ADV OPCODE = 0
	BXL OPCODE = 1
	BST OPCODE = 2
	JNZ OPCODE = 3
	BXC OPCODE = 4
	OUT OPCODE = 5
	BDV OPCODE = 6
	CDV OPCODE = 7
)

type Instruction struct {
	opcode  OPCODE
	operand byte
}

type Machine struct {
	registers         []int
	startingRegisters []int

	instrPointer int
	intructions  []Instruction
}

func NewMachine(instructions []Instruction, A, B, C int) Machine {
	return Machine{
		registers:         []int{A, B, C},
		startingRegisters: []int{A, B, C},

		instrPointer: 0,
		intructions:  instructions,
	}
}

func (m *Machine) Reset() {
	for i := range m.registers {
		m.registers[i] = m.startingRegisters[i]
	}
}

func (m *Machine) RawInstructions() string {
	bytes := make([]string, len(m.intructions)*2)
	for ix := range m.intructions {
		bytes[ix*2] = fmt.Sprintf("%d", m.intructions[ix].opcode)
		bytes[ix*2+1] = fmt.Sprintf("%d", m.intructions[ix].operand)
	}
	return strings.Join(bytes, ",")
}

func (m *Machine) Run() ([]int, error) {
	allouts := make([]int, 0)
	for {
		out, err := m.Advance()
		if errors.Is(err, ErrNoMoreInstructions) {
			break
		}
		if err != nil {
			return nil, err
		}
		allouts = append(allouts, out)
	}
	return allouts, nil
}

var ErrNoMoreInstructions = errors.New("out of instructions")

// Advance the machine to the next OUT instruction
// Returns ErrNoMoreInstructions
func (m *Machine) Advance() (int, error) {
	for {
		if m.instrPointer >= len(m.intructions) {
			return 0, ErrNoMoreInstructions
		}

		instr := m.intructions[m.instrPointer]

		// fmt.Printf("(%d) %+v\n", m.instrPointer, instr)

		switch instr.opcode {
		case ADV:
			m.registers[0] /= int(math.Pow(2, float64(m.valueForCombo(instr.operand))))
		case BXL:
			m.registers[1] ^= int(instr.operand)
		case BST:
			m.registers[1] = m.valueForCombo(instr.operand) & 0b111
		case JNZ:
			if m.registers[0] != 0 {
				m.instrPointer = int(instr.operand)
				continue
			}
		case BXC:
			m.registers[1] ^= m.registers[2]
		case OUT:
			m.instrPointer++
			return m.valueForCombo(instr.operand) & 0b111, nil
		case BDV:
			m.registers[1] = m.registers[0] / int(math.Pow(2, float64(m.valueForCombo(instr.operand))))
		case CDV:
			m.registers[2] = m.registers[0] / int(math.Pow(2, float64(m.valueForCombo(instr.operand))))
		}

		m.instrPointer++
	}

}

func (m *Machine) valueForCombo(operand byte) int {
	switch operand {
	case 0:
		return int(operand)
	case 1:
		return int(operand)
	case 2:
		return int(operand)
	case 3:
		return int(operand)
	case 4:
		return m.registers[operand-4]
	case 5:
		return m.registers[operand-4]
	case 6:
		return m.registers[operand-4]
	}

	err := fmt.Sprintf("invalid combo operand %b", operand)
	panic(err)
}

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	machine := parse(file)

	if len(os.Args) == 2 {
		machine.Reset()
		fmt.Printf("solution 1: %s\n", sol1(machine))
		machine.Reset()
		fmt.Printf("solution 2: %d\n", sol2(machine))
	}

	if len(os.Args) == 3 {
		exp := os.Args[2]
		asint, err := strconv.ParseInt(exp, 10, 64)

		if err == nil {
			res := with_a(machine, int(asint))
			fmt.Printf("output for A=%d (o%o) (b%b): %s\n", asint, asint, asint, res)
		} else {
			fmt.Printf("finding '%s'\n", exp)
			res := smallest_for(machine, exp)
			fmt.Printf("output for A=%d (o%o) (b%b): %s\n", res, res, res, exp)
		}
	}

	if len(os.Args) == 4 {
		from, err := strconv.ParseInt(os.Args[2], 10, 64)
		if err != nil {
			panic(err)
		}

		to, err := strconv.ParseInt(os.Args[3], 10, 64)
		if err != nil {
			panic(err)
		}

		// v := 5
		// for i := from; i < to; i++ {
		// 	machine.Reset()
		// 	fmt.Printf("output for A=%d (o%o) (b%b): %s\n", v, v, v, with_a(machine, int(v)))
		// 	v = v << 3
		// 	v += 1
		// }

		for i := from; i < to; i++ {
			fmt.Printf("output for A=%d (o%o) (b%b): %s\n", i, i, i, with_a(machine, int(i)))
		}
	}

}

func sol1(m Machine) string {
	outs, err := m.Run()
	if err != nil {
		panic(err)
	}
	return join(outs)
}

func with_a(m Machine, A int) string {
	m.Reset()
	m.registers[0] = A
	outs, err := m.Run()
	if err != nil {
		panic(err)
	}
	return join(outs)
}

func smallest_for(m Machine, exp string) int {
	return smallest_from(m, exp, 0)
}

func smallest_from(m Machine, exp string, from int) int {
	n := from
	for {
		res := with_a(m, n)
		if res == exp {
			return n
		}
		if len(res) > len(exp)+3 {
			return 0
		}
		n++
	}
}

func sol2(m Machine) int {
	// get all instructions
	nums := make([]byte, 0)
	for i := range m.intructions {
		nums = append(nums, byte(m.intructions[i].opcode))
		nums = append(nums, byte(m.intructions[i].operand))
	}

	octIx := 15            // hardcoded for the 16 instruction input
	n := 01000000000000000 // start with leftmost octet at 1

	for {
		if octIx < 0 {
			break
		}

		exp := join(nums[octIx:])

		wa := with_a(m, n)
		spl := strings.Split(wa, ",")[octIx:]
		res := strings.Join(spl, ",")

		// fmt.Printf("0o%o (%d)-> %+v == %+v\n", n, octIx, res, exp)

		if res == exp {
			octIx -= 1
			continue
		}

		n += 01 << (3 * octIx) // add 1 to the leftmost unmatching octal
	}

	return n
}

func join[T int | byte](elems []T) string {
	outs := make([]string, len(elems))
	for ix := range elems {
		outs[ix] = fmt.Sprintf("%d", elems[ix])
	}
	return strings.Join(outs, ",")
}

func parse(file *os.File) Machine {
	r := bufio.NewReader(file)

	skipTo(r, ':')
	rA := readIntTo(r, '\n')
	skipTo(r, ':')
	rB := readIntTo(r, '\n')
	skipTo(r, ':')
	rC := readIntTo(r, '\n')

	skipTo(r, '\n')
	skipTo(r, ' ')

	instr := make([]Instruction, 0)
	for {
		i := Instruction{}

		op, err := r.ReadByte()
		if err != nil {
			panic(err)
		}
		i.opcode = OPCODE(op - 48)

		sep, err := r.ReadByte()
		if err != nil {
			panic(err)
		}
		if sep != ',' {
			panic(fmt.Sprintf("wrong separator '%c'", sep))
		}

		op, err = r.ReadByte()
		if err != nil {
			panic(err)
		}
		i.operand = op - 48

		instr = append(instr, i)

		sep, err = r.ReadByte()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}
		if sep != ',' {
			panic(fmt.Sprintf("wrong separator %c", sep))
		}
	}

	return NewMachine(instr, rA, rB, rC)
}

func skipTo(r *bufio.Reader, delim byte) error {
	if _, err := r.ReadBytes(delim); err != nil {
		if errors.Is(err, io.EOF) {
			return err
		}
		panic(err)
	}
	return nil
}

func readIntTo(r *bufio.Reader, delim byte) int {
	val, err := r.ReadString(delim)
	if err != nil {
		panic(err)
	}

	parsed, err := strconv.ParseInt(val[1:len(val)-1], 10, 64)
	if err != nil {
		panic(err)
	}

	return int(parsed)
}
