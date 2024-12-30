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
	sol1(program)
	sol2(program)

}

func sol1(program []int) {
	outbuf := bytes.NewBuffer([]byte{})
	m := vm.NewVM(program, vmio.CreateInBuffer(1), outbuf)
	m.Run()
	res := vmio.GetLastOutput(outbuf)
	fmt.Printf("solution 1: %d\n", res)
}

func sol2(program []int) {
	outbuf := bytes.NewBuffer([]byte{})
	m := vm.NewVM(program, vmio.CreateInBuffer(2), outbuf)
	m.Run()
	res := vmio.GetLastOutput(outbuf)
	fmt.Printf("solution 2: %d\n", res)
}
