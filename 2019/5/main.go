package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"2019/utils/machine/parser"
	"2019/utils/machine/vm"
)

func main() {
	nums := parser.Parse(os.Args[1])
	sol1(nums)
	sol2(nums)

}

func sol1(nums []int) {
	inBuf := createInBuffer(1)
	outBuf := bytes.NewBuffer([]byte{})

	m := vm.NewVM(nums, inBuf, outBuf)
	m.Run()

	fmt.Printf("solution 1: %d\n", getLastOutput(outBuf))
}

func sol2(nums []int) {
	inBuf := createInBuffer(5)
	outBuf := bytes.NewBuffer([]byte{})

	m := vm.NewVM(nums, inBuf, outBuf)
	m.Run()

	fmt.Printf("solution 2: %d\n", getLastOutput(outBuf))
}

func createInBuffer(inputs ...int) io.Reader {
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

		if num != 0 {
			panic(fmt.Sprintf("ERROR: expected 0 but got %d", num))
		}
	}

	return -1
}
