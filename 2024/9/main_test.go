package main

import "testing"

func Test2(t *testing.T) {
	data := []byte{2, 3, 3, 3, 1, 3, 3, 1, 2, 1, 4, 1, 4, 1, 3, 1, 4, 0, 2}
	res := sol2(data)
	if res != 2858 {
		t.Fatalf("expected 2858 got %d", res)
	}
}
