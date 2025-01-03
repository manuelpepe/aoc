package vm

import (
	"2019/utils/machine/parser"
	"2019/utils/machine/vmio"
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Day2(t *testing.T) {
	nums := parser.Parse("res/day2-input-0.txt")

	// example
	out := NewVM(nums).Run()
	assert.Equal(t, 3500, out)

	nums = parser.Parse("res/day2-input-1.txt")

	// part 1
	nums[1], nums[2] = 12, 2
	out = NewVM(nums).Run()
	assert.Equal(t, 8017076, out)

	// part 2
	nums[1], nums[2] = 31, 46
	out = NewVM(nums).Run()
	assert.Equal(t, 19690720, out)
}

func Test_Day7(t *testing.T) {
	nums := parser.Parse("res/day7-input-1.txt")

	outbuf := bytes.NewBuffer([]byte{})
	NewVM(
		nums,
		WithInput(vmio.CreateInBuffer(0, 1)),
		WithOutput(outbuf),
	).Run()
	last := vmio.GetLastOutput(outbuf)
	assert.Equal(t, 11, last)

	outbuf = bytes.NewBuffer([]byte{})
	NewVM(
		nums,
		WithInput(vmio.CreateInBuffer(0, 2)),
		WithOutput(outbuf),
	).Run()
	last = vmio.GetLastOutput(outbuf)
	assert.Equal(t, 13, last)
}

func Test_Day9(t *testing.T) {
	nums := parser.Parse("res/day9-input-1.txt")

	outbuf := bytes.NewBuffer([]byte{})
	NewVM(nums,
		WithInput(vmio.CreateInBuffer(1)),
		WithOutput(outbuf),
	).Run()
	last := vmio.GetLastOutput(outbuf)
	assert.Equal(t, 2662308295, last)

	outbuf = bytes.NewBuffer([]byte{})
	NewVM(nums,
		WithInput(vmio.CreateInBuffer(2)),
		WithOutput(outbuf),
	).Run()
	last = vmio.GetLastOutput(outbuf)
	assert.Equal(t, 63441, last)
}

func Test_ParseParameterModes(t *testing.T) {
	out := parseModes(10102)
	assert.Equal(t, out, []byte{1, 0, 1})

	out = parseModes(02)
	assert.Equal(t, out, []byte{})
}

func Test_SetWithRelativeParam(t *testing.T) {
	program := []int{203, 0, 99}
	m := NewVM(
		program,
		WithInput(vmio.CreateInBuffer(12345)),
		WithOutput(os.Stdout),
	)
	out := m.Run()
	assert.Equal(t, 12345, out)
}
