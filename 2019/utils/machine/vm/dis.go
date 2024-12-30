package vm

import "2019/utils/machine/opcodes"

func Dissasemble(program []int) []opcodes.Instruction {
	m := NewVM(program, nil, nil)

	out := make([]opcodes.Instruction, 0)

	for m.intrsPoint < len(program) {
		instr := m.curInstruction()
		out = append(out, instr)
		m.advanceInstrPointer(len(instr.Params))
	}

	return out
}
