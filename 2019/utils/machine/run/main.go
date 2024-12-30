package main

import (
	"2019/utils/machine/parser"
	"2019/utils/machine/vm"
	"os"
)

func main() {
	nums := parser.Parse(os.Args[1])
	m := vm.NewVM(nums, os.Stdin, os.Stdout)
	m.Run()
}
