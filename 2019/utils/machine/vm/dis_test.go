package vm

import (
	"2019/utils/machine/opcodes"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDis(t *testing.T) {
	instrs := slices.Collect(DissasembleIter([]int{203, 0, 99}))
	assert.Len(t, instrs, 2)
	assert.Equal(t, opcodes.OP_INP, instrs[0].Opcode)
	assert.Equal(t, opcodes.OP_HALT, instrs[1].Opcode)
}
