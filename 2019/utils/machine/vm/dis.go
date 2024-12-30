package vm

import (
	"iter"

	"2019/utils/machine/opcodes"
)

func DissasembleIter(program []int) iter.Seq[opcodes.Instruction] {
	return func(yield func(opcodes.Instruction) bool) {
		m := NewVM(program, nil, nil)

		for m.intrsPoint < len(program) {
			instr := m.curInstruction()
			if !yield(instr) {
				return
			}
			m.advanceInstrPointer(len(instr.Params))
		}
	}
}
