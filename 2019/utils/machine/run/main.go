package main

import (
	"flag"
	"fmt"
	"os"

	"2019/utils/machine/parser"
	"2019/utils/machine/vm"
)

func main() {
	runPath := flag.String("run", "", "file to run")
	disPath := flag.String("dis", "", "file to dissasemble")

	flag.Parse()

	if runPath != nil && *runPath != "" {
		program := parser.Parse(*runPath)
		m := vm.NewVM(program, os.Stdin, os.Stdout)
		m.Run()
	}

	if disPath != nil && *disPath != "" {
		program := parser.Parse(*disPath)
		out := vm.Dissasemble(program)

		vix := 0
		for _, instr := range out {
			fmt.Printf("%03d:\t%s\n", vix, instr)
			vix += 1 + len(instr.Params)
		}
	}
}
