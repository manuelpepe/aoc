package vmio

import (
	"bytes"
	"context"
	"fmt"
	"io"
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
	outs := GetOutput(buf)
	return outs[len(outs)-1]
}

func GetOutput(buf *bytes.Buffer) []int {
	rawOut := buf.String()
	rawOut = strings.TrimSuffix(rawOut, "\n")
	lines := strings.Split(rawOut, "\n")

	buf.Reset() // clear buffer

	out := make([]int, len(lines))

	for ix, line := range lines {
		prefix := "OUT: "

		if !strings.HasPrefix(line, prefix) {
			panic(fmt.Sprintf("expected 'OUT: ' suffix. Got: '%s'", line))
		}

		num, err := strconv.ParseInt(line[len(prefix):], 10, 64)
		if err != nil {
			panic(err)
		}

		out[ix] = int(num)
	}

	return out
}

type SendInputFunc func(inp []byte) bool

func StartInputSender(inbuf io.Writer) (SendInputFunc, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	inputs := make(chan []byte)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case v := <-inputs:
				fmt.Fprintf(inbuf, "%s\n", v)
			}
		}
	}()

	return func(inp []byte) bool {
		select {
		case inputs <- inp:
			return true
		case <-ctx.Done():
			return false
		}
	}, cancel
}