package vm

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Day2(t *testing.T) {
	nums := parseDay2("res/day2-input-0.txt")

	// part 0
	m := NewVM(nums, os.Stdin, os.Stdout)
	out := m.Run()
	assert.Equal(t, 3500, out)

	nums = parseDay2("res/day2-input-1.txt")

	// part 1
	nums[1] = 12
	nums[2] = 2
	m = NewVM(nums, os.Stdin, os.Stdout)
	out = m.Run()
	assert.Equal(t, 8017076, out)

	// part 2
	nums[1] = 31
	nums[2] = 46
	m = NewVM(nums, os.Stdin, os.Stdout)
	out = m.Run()
	assert.Equal(t, 19690720, out)
}

func parseDay2(fn string) []int {
	rawdata, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	data, _ := strings.CutSuffix(string(rawdata), "\n")
	items := strings.Split(data, ",")

	out := make([]int, 0)

	for _, i := range items {
		n, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			panic(err)
		}
		out = append(out, int(n))
	}

	return out
}

func Test_ParseParameterModes(t *testing.T) {
	out := parseModes(10102)
	assert.Equal(t, out, []byte{1, 0, 1})

	out = parseModes(02)
	assert.Equal(t, out, []byte{})
}
