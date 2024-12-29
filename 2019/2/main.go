package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	nums := parse(os.Args[1])
	fmt.Printf("solution 1: %d\n", sol1(nums))
	fmt.Printf("solution 2: %d\n", sol2(nums))
}

func sol1(nums []int) int {
	nums[1] = 12
	nums[2] = 2
	m := NewMachine(nums)
	return m.Run()
}

func sol2(nums []int) int {
	// could be BFS'd, just did it manually
	// relize that N at [2] only increments output by N
	// with [2] = 0 find max [1] that produces a max value M below target T, then [2] = T - M for answer
	nums[1] = 31
	nums[2] = 46
	m := NewMachine(nums)
	return m.Run()
}

type Machine struct {
	memory     []int
	intrsPoint int
}

const MAX_MEM = 2048

func NewMachine(program []int) *Machine {
	if len(program) > MAX_MEM {
		panic("memory limit exceeded")
	}

	mem := make([]int, MAX_MEM)
	copy(mem, program)

	return &Machine{
		memory:     mem,
		intrsPoint: 0,
	}
}

func (m *Machine) Run() int {
OUT:
	for {

		opcode := m.curOpcode()

		switch opcode {
		case 99:
			break OUT // HALT
		case 1:
			outloc := m.operand(3)
			m.memory[outloc] = m.memory[m.operand(1)] + m.memory[m.operand(2)]
		case 2:
			outloc := m.operand(3)
			m.memory[outloc] = m.memory[m.operand(1)] * m.memory[m.operand(2)]
		}

		m.advanceInstrPointer()
	}

	return m.memory[0]
}

func (m *Machine) curOpcode() int {
	opcode := m.memory[m.intrsPoint]
	return opcode
}

func (m *Machine) operand(n int) int {
	return m.memory[m.intrsPoint+n]
}

func (m *Machine) advanceInstrPointer() {
	m.intrsPoint += 4
}

func parse(fn string) []int {
	rawdata, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	data, _ := strings.CutSuffix(string(rawdata), "\n")
	items := strings.Split(data, ",")

	out := make([]int, 0)

	for _, i := range items {
		n, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			panic(err)
		}
		out = append(out, int(n))
	}

	return out
}
