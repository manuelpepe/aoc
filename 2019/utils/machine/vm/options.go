package vm

import (
	"fmt"
	"io"
)

type IOMode int

const (
	IOInteger IOMode = 1
	IOAscii   IOMode = 2
)

type vmconfig struct {
	reader io.Reader
	writer io.Writer
	mode   IOMode

	id         string
	debugLevel int
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

func WithNumericID(id int) vmOptionFunc {
	return func(v *vmconfig) {
		v.id = fmt.Sprintf("%03d", id)
	}
}

func WithDebugLevel(lvl int) vmOptionFunc {
	return func(v *vmconfig) {
		v.debugLevel = lvl
	}
}
