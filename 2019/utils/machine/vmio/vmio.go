package vmio

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

func CreateInBuffer(inputs ...int) *bytes.Buffer {
	buf := bytes.NewBuffer([]byte{})
	for _, n := range inputs {
		fmt.Fprintf(buf, "%d\n", n)
	}
	return buf
}

func GetLastOutput(buf *bytes.Buffer) int {
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
