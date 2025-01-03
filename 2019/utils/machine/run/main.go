package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"2019/utils/machine/parser"
	"2019/utils/machine/vm"
)

func main() {
	runPath := flag.String("run", "", "file to run")
	disPath := flag.String("dis", "", "file to dissasemble")
	asciiMode := flag.Bool("ascii", false, "treat output as ascii codes")

	flag.Parse()

	if runPath != nil && *runPath != "" {
		if *asciiMode {
			runAsciiMode(*runPath)
		} else {
			runIntegerMode(*runPath)
		}
	}

	if disPath != nil && *disPath != "" {
		runDissasembler(*disPath)
	}
}

func runDissasembler(srcfile string) {
	program := parser.Parse(srcfile)
	out := vm.DissasembleIter(program)

	vix := 0
	for instr := range out {
		fmt.Printf("%03d:\t%s\n", vix, instr)
		vix += 1 + len(instr.Params)
	}
}

func runIntegerMode(srcfile string) {
	program := parser.Parse(srcfile)
	m := vm.NewVM(
		program,
		vm.WithInput(os.Stdin),
		vm.WithOutput(os.Stdout),
		vm.WithIOMode(vm.IOInteger),
	)
	m.Run()
}

func runAsciiMode(srcfile string) {
	program := parser.Parse(srcfile)

	inReader, inWriter := io.Pipe()
	cancel := asciiReader(inWriter)
	defer cancel()

	m := vm.NewVM(
		program,
		vm.WithInput(inReader),
		vm.WithOutput(os.Stdout),
		vm.WithIOMode(vm.IOAscii),
	)
	m.Run()
}

func asciiReader(w io.Writer) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	inps := make(chan []byte)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			data, err := reader.ReadBytes('\n')
			if err != nil {
				err = fmt.Errorf("error reading stdin: %w", err)
				panic(err)
			}

			inps <- data
		}
	}()

	// process goroutine reads the inputs channel and sends its data to the writer
	go func() {
		for {
			select {
			case data := <-inps:
				for _, b := range data {
					if _, err := w.Write([]byte{b}); err != nil {
						panic(err)
					}
				}
			case <-ctx.Done():
				fmt.Printf("done reading stdin\n")
				return
			}
		}
	}()

	return cancel
}
