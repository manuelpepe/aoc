package vm

import "io"

type IOMode int

const (
	IOInteger IOMode = 1
	IOAscii   IOMode = 2
)

type vmconfig struct {
	reader io.Reader
	writer io.Writer
	mode   IOMode
}

type vmOptionFunc func(*vmconfig)

func WithInput(reader io.Reader) vmOptionFunc {
	return func(v *vmconfig) {
		v.reader = reader
	}
}

func WithOutput(writer io.Writer) vmOptionFunc {
	return func(v *vmconfig) {
		v.writer = writer
	}
}

func WithIOMode(mode IOMode) vmOptionFunc {
	return func(v *vmconfig) {
		v.mode = mode
	}
}
