package main

import "testing"

func TestAdvance(t *testing.T) {
	bots := []Robot{
		{
			X:  2,
			Y:  4,
			VX: 2,
			VY: -3,
		},
	}

	res := advance(bots, 11, 7, 5)
	if res[0].X != 1 {
		t.Fatalf("expected X to be 1, got %d", res[0].X)
	}
	if res[0].Y != 3 {
		t.Fatalf("expected Y to be 3, got %d", res[0].Y)
	}
}
