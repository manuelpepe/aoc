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
	lines := GetOutputRaw(buf)

	out := make([]int, len(lines))

	for ix, line := range lines {
		prefix := "OUT: "

		if !strings.HasPrefix(line, prefix) {
			continue
		}

		num, err := strconv.ParseInt(line[len(prefix):], 10, 64)
		if err != nil {
			panic(err)
		}

		out[ix] = int(num)
	}

	return out
}

func GetOutputRaw(buf *bytes.Buffer) []string {
	rawOut := buf.String()
	buf.Reset() // clear buffer
	rawOut = strings.TrimSuffix(rawOut, "\n")
	return strings.Split(rawOut, "\n")
}

type SendInputFunc func(inp []byte) bool

func StartInputSender(inbuf io.Writer) (SendInputFunc, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	inputs := make(chan []byte)

	go func() {
		for {
			select {
			case v := <-inputs:
				if _, err := inbuf.Write(v); err != nil {
					panic(err)
				}
			case <-ctx.Done():
				return
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
