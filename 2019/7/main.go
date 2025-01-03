package main

import (
	"2019/utils/machine/parser"
	"2019/utils/machine/vm"
	"bytes"
	"fmt"
	"iter"
	"os"
	"strconv"
	"strings"
)

func main() {
	program := parser.Parse(os.Args[1])
	sol1(program)
	sol2(program)
}

func sol2(program []int) {
	lastMax := 0

PHASES:
	for phase := range alts(5, 9) {
		inbuffers := []*bytes.Buffer{
			createInBuffer(phase[0]),
			createInBuffer(phase[1]),
			createInBuffer(phase[2]),
			createInBuffer(phase[3]),
			createInBuffer(phase[4]),
		}

		outbuffers := []*bytes.Buffer{
			bytes.NewBuffer([]byte{}),
			bytes.NewBuffer([]byte{}),
			bytes.NewBuffer([]byte{}),
			bytes.NewBuffer([]byte{}),
			bytes.NewBuffer([]byte{}),
		}

		amps := []*vm.VM{
			vm.NewVM(program, vm.WithInput(inbuffers[0]), vm.WithOutput(outbuffers[0])),
			vm.NewVM(program, vm.WithInput(inbuffers[1]), vm.WithOutput(outbuffers[1])),
			vm.NewVM(program, vm.WithInput(inbuffers[2]), vm.WithOutput(outbuffers[2])),
			vm.NewVM(program, vm.WithInput(inbuffers[3]), vm.WithOutput(outbuffers[3])),
			vm.NewVM(program, vm.WithInput(inbuffers[4]), vm.WithOutput(outbuffers[4])),
		}

		lastOutput := 0
		loopCount := 0

		for {
			for ix, vm := range amps {
				in := lastOutput
				fmt.Fprintf(inbuffers[ix], "%d\n", in)
				vm.RunForOutput()

				if vm.Halted() {
					// fmt.Printf("Halted on VM %d at loop %d\n", ix, loopCount)
					// fmt.Printf("Result: %d\n", lastOutput)
					lastMax = max(lastMax, lastOutput)
					continue PHASES
				}

				lastOutput = getLastOutput(outbuffers[ix])
				// fmt.Printf("-> Loop #%d - VM #%d - IN:%d OUT:%d\n", loopCount, ix, in, lastOutput)

			}
			loopCount++
		}
	}

	fmt.Printf("solution 2: %d\n", lastMax)
}

func sol1(program []int) {
	maxval := 0

	for phases := range alts(0, 4) {
		o0 := runAmp(program, phases[0], 0)
		o1 := runAmp(program, phases[1], o0)
		o2 := runAmp(program, phases[2], o1)
		o3 := runAmp(program, phases[3], o2)
		o4 := runAmp(program, phases[4], o3)
		maxval = max(o4, maxval)
	}

	fmt.Printf("solution 1: %d\n", maxval)
}

func runAmp(program []int, phase, val int) int {
	inBuf := createInBuffer(phase, val)
	outBuf := bytes.NewBuffer([]byte{})

	m := vm.NewVM(program, vm.WithInput(inBuf), vm.WithOutput(outBuf))
	m.Run()

	return getLastOutput(outBuf)
}

func alts(min, max int) iter.Seq[[]int] {
	valid := func(p []int) bool {
		kv := make(map[int]int)
		for _, v := range p {
			kv[v]++
			if kv[v] > 1 {
				return false
			}
		}
		return true
	}

	return func(yield func([]int) bool) {
		for a := min; a <= max; a++ {
			for b := min; b <= max; b++ {
				for c := min; c <= max; c++ {
					for d := min; d <= max; d++ {
						for e := min; e <= max; e++ {
							phase := []int{a, b, c, d, e}
							if valid(phase) {
								if !yield(phase) {
									return
								}
							}
						}
					}
				}
			}
		}
	}
}

// TODO: move these funcs to utils

func createInBuffer(inputs ...int) *bytes.Buffer {
	buf := bytes.NewBuffer([]byte{})
	for _, n := range inputs {
		fmt.Fprintf(buf, "%d\n", n)
	}
	return buf
}

func getLastOutput(buf *bytes.Buffer) int {
	rawOut := buf.String()
	rawOut = strings.TrimSuffix(rawOut, "\n")
	lines := strings.Split(rawOut, "\n")

	buf.Reset() // clear buffer

	for ix, line := range lines {
		prefix := "OUT: "

		if !strings.HasPrefix(line, prefix) {
			panic("expected 'OUT: ' suffix")
		}

		num, err := strconv.ParseInt(line[len(prefix):], 10, 64)
		if err != nil {
			panic(err)
		}

		if ix == len(lines)-1 {
			return int(num)
		}
	}

	return -1
}
