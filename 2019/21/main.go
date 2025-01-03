package main

import (
	"fmt"
	"os"

	"2019/utils/machine/parser"
	"2019/utils/machine/vm"
	"2019/utils/machine/vmio"
)

func main() {
	program := parser.Parse(os.Args[1])
	sol1(program)
}

func sol1(program []int) {
	m, inbuf, outbuf := vm.NewVMPiped(program, vm.WithIOMode(vm.IOAscii))
	sendInput, stop := vmio.StartInputSender(inbuf)
	defer stop()

	sendInput([]byte("NOT A J\nNOT B T\nOR T J\nNOT C T\nOR T J\nAND D J\nWALK\n"))
	m.Run()

	lines := vmio.GetOutputRaw(outbuf)

	fmt.Printf("solution 1: %s\n", lines[len(lines)-1])
}
