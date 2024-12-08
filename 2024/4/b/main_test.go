package main

import "testing"

func TestFindAdjacent(t *testing.T) {
	inp := [][]byte{
		{0, 0, 0, 0, 0, 0},
		{0, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 9, 0, 0},
		{0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0},
	}

	var outs [][]int

	outs = findPosibleDirection(inp, 0, 0, 1)
	if len(outs) != 1 {
		t.Fail()
	}

	outs = findPosibleDirection(inp, 0, 0, 0)
	if len(outs) != 2 {
		t.Fail()
	}

	outs = findPosibleDirection(inp, 4, 1, 0)
	if len(outs) != 8 {
		t.Fail()
	}
}

func TestCountSeq(t *testing.T) {
	inp := [][]byte{
		{3, 0, 3, 0, 1, 0},
		{0, 2, 0, 0, 0, 0},
		{1, 0, 1, 1, 1, 1},
		{0, 0, 0, 2, 0, 0},
		{0, 0, 3, 3, 3, 0},
		{3, 2, 1, 0, 3, 0},
	}

	out := countExactSequence(inp, []byte{1, 2, 3})
	if out != 2 {
		t.Fail()
	}
}
