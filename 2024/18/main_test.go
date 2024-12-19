package main

import (
	"os"
	"testing"
)

func TestBasic(t *testing.T) {
	file, err := os.Open("input-0.txt")
	if err != nil {
		panic(err)
	}

	size := 7
	limit := 12

	scanner, g := parse(file, size, limit)

	sol2(g, scanner)
}
