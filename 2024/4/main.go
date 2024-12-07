package main

import (
	"bufio"
	"fmt"
	"os"
)

var DIRS = [][]int{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1} /*{0, 0},*/, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

func main() {
	fh, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	r := bufio.NewScanner(fh)
	matrix := make([][]byte, 0)

	for r.Scan() {
		line := r.Text()
		matrix = append(matrix, []byte(line))
	}

	res := countExactSequence(matrix, []byte("XMAS"))
	fmt.Printf("XMAS appears: %d times\n", res)
}

func countExactSequence(matrix [][]byte, seq []byte) int {
	acc := 0

	for y := range matrix {
		for x := range matrix[y] {
			if matrix[y][x] == seq[0] {
				seqs := recCountSeq(matrix, x, y, seq)
				acc += seqs
			}
		}
	}

	return acc
}

func recCountSeq(matrix [][]byte, x, y int, seq []byte) int {
	if matrix[y][x] != seq[0] {
		return 0
	}

	total := 0

	dirs := findPosibleDirection(matrix, x, y, seq[1])
	for _, posibleDir := range dirs {
		total += recCountSeqWithDir(matrix, x, y, posibleDir, seq, 1)
	}

	return total
}

func recCountSeqWithDir(matrix [][]byte, x, y int, dir []int, seq []byte, depth int) int {
	aX, aY := x+dir[0], y+dir[1]

	// base: not in matrix bounds
	if !inBounds(matrix, aX, aY) {
		return 0
	}

	// base: overflowed seq
	if depth == len(seq) {
		return 0
	}

	// base: last char of seq
	if depth == len(seq)-1 {
		if seq[depth] == matrix[aY][aX] {
			return 1
		} else {
			return 0
		}
	}

	// base: next not valid
	if matrix[aY][aX] != seq[depth] {
		return 0
	}

	return recCountSeqWithDir(matrix, aX, aY, dir, seq, depth+1)
}

func findPosibleDirection(matrix [][]byte, x, y int, ch byte) [][]int {
	out := make([][]int, 0)

	for _, dir := range DIRS {
		aX, aY := x+dir[0], y+dir[1]

		if !inBounds(matrix, aX, aY) {
			continue
		}

		if matrix[aY][aX] == ch {
			out = append(out, dir)
		}
	}

	return out
}

func inBounds(matrix [][]byte, x, y int) bool {
	if y < 0 || y >= len(matrix) {
		return false
	}

	if x < 0 || x >= len(matrix[0]) {
		return false
	}

	return true
}
