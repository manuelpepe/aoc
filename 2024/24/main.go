package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

type Wire string
type GateType string

type Gate struct {
	Inputs []Wire
	Type   GateType
	Output Wire
}

func (g Gate) Eval(a, b bool) bool {
	switch g.Type {
	case "XOR":
		return a != b
	case "OR":
		return a || b
	case "AND":
		return a && b
	}
	panic(fmt.Sprintf("unexpected gate type: %s", g.Type))
}

type System struct {
	Wires map[Wire]bool
	Gates []Gate
}

func (s *System) Run() {
	pendingGates := make([]int, len(s.Gates))
	for ix := range pendingGates {
		pendingGates[ix] = ix
	}

	for len(pendingGates) > 0 {
		// pop
		curGateIx := pendingGates[0]
		pendingGates = pendingGates[1:]

		gate := s.Gates[curGateIx]

		if s.IsWritten(gate.Output) {
			continue // done
		}

		if !s.IsReadyForEval(gate) {
			// push back
			pendingGates = append(pendingGates, curGateIx)
			continue
		}

		inA := s.Wires[gate.Inputs[0]]
		inB := s.Wires[gate.Inputs[1]]

		s.Wires[gate.Output] = gate.Eval(inA, inB)
	}
}

func (s *System) GetOutput() int {
	zeds := make([]Wire, 0)
	for k, _ := range s.Wires {
		if k[0:1] != "z" {
			continue
		}

		zeds = append(zeds, k)
	}

	slices.Sort(zeds)
	slices.Reverse(zeds)

	val := 0
	for _, w := range zeds {
		fmt.Printf("%s: %v\n", w, s.Wires[w])
		if s.Wires[w] {
			val++
		}
		val = val << 1
		fmt.Printf("%016b\n", val)
	}

	val = val >> 1

	return val
}

func (s *System) IsWritten(wire Wire) bool {
	_, ok := s.Wires[wire]
	return ok
}

func (s *System) IsReadyForEval(gate Gate) bool {
	return s.IsWritten(gate.Inputs[0]) && s.IsWritten(gate.Inputs[1])
}

func main() {
	sys := parse(os.Args[1])
	sys.Run()

	fmt.Printf("%+v\n", sys.GetOutput())
}

func parse(fn string) *System {
	fh, err := os.Open(fn)
	if err != nil {
		panic(err)
	}

	sys := System{
		Wires: make(map[Wire]bool),
		Gates: make([]Gate, 0),
	}

	s := bufio.NewScanner(fh)

	parsingGates := false
	for s.Scan() {
		line := s.Text()
		if line == "" {
			parsingGates = true
			continue
		}

		if parsingGates {
			parts := strings.Split(line, " ")
			sys.Gates = append(sys.Gates, Gate{
				Inputs: []Wire{Wire(parts[0]), Wire(parts[2])},
				Type:   GateType(parts[1]),
				Output: Wire(parts[4]),
			})
			continue
		}

		// parsing default values
		parts := strings.Split(line, ": ")
		if len(parts) != 2 {
			panic("expected two parts parsing wire")
		}

		sys.Wires[Wire(parts[0])] = parts[1] == "1"
	}

	return &sys
}
