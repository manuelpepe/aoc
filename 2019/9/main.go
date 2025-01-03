package main

import (
	"2019/utils/machine/parser"
	"2019/utils/machine/vm"
	"2019/utils/machine/vmio"
	"bytes"
	"fmt"
	"os"
)

func main() {
	program := parser.Parse(os.Args[1])
	run(program, 1)
	run(program, 2)
}

func run(program []int, n int) {
	outbuf := bytes.NewBuffer([]byte{})
	m := vm.NewVM(
		program,
		vm.WithInput(vmio.CreateInBuffer(n)),
		vm.WithOutput(outbuf),
	)
	m.Run()
	res := vmio.GetLastOutput(outbuf)
	fmt.Printf("solution %d: %d\n", n, res)
}
