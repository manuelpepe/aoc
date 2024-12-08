package main

import (
	"bufio"
	"fmt"
	"os"
)

// {-1, -1},
//		-x 0
//		0 -y
// {-1, 1},
//		-x 0
//		0 +y
// {1, -1},
//		+x 0
//		0 -y
// {1, 1},
//		+x 0
//		0 +y

var DIRS = [][]int{
	{-1, -1}, {-1, 1},
	{1, -1}, {1, 1},
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

	res := countExactSequence(matrix, []byte("MAS"))

	fmt.Printf("X-MAS appears: %d times\n", res)
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

	// divide by two because each X is counted twice
	// I.e.
	// 	1 . 1
	// 	. 2 .
	// 	3 . 3
	// is counted as two crosses but it's only one
	return acc / 2
}

func recCountSeq(matrix [][]byte, x, y int, seq []byte) int {
	if matrix[y][x] != seq[0] {
		return 0
	}

	total := 0

	dirs := findPosibleDirection(matrix, x, y, seq[1])
	for _, posibleDir := range dirs {
		ok := recCountSeqWithDir(matrix, x, y, posibleDir, seq, 1)
		if !ok {
			continue
		}

		// check first possible crossing seq
		deltaDiagX := (len(seq) - 1) * posibleDir[0]
		diagX, diagY := x+deltaDiagX, y
		diagDir := []int{-posibleDir[0], posibleDir[1]}

		if inBounds(matrix, diagX, diagY) && matrix[diagY][diagX] == seq[0] {
			if ok := recCountSeqWithDir(matrix, diagX, diagY, diagDir, seq, 1); ok {
				total++
				continue
			}
		}

		// check second possible crossing seq
		deltaDiagY2 := (len(seq) - 1) * posibleDir[1]
		diagX2, diagY2 := x, y+deltaDiagY2
		diagDir2 := []int{posibleDir[0], -posibleDir[1]}

		if inBounds(matrix, diagX2, diagY2) && matrix[diagY2][diagX2] == seq[0] {
			if ok := recCountSeqWithDir(matrix, diagX2, diagY2, diagDir2, seq, 1); ok {
				total++
				continue
			}
		}

	}

	return total
}

func recCountSeqWithDir(matrix [][]byte, x, y int, dir []int, seq []byte, depth int) bool {
	aX, aY := x+dir[0], y+dir[1]

	// base: not in matrix bounds
	if !inBounds(matrix, aX, aY) {
		return false
	}

	// base: overflowed seq
	if depth == len(seq) {
		return false
	}

	// base: last char of seq
	if depth == len(seq)-1 {
		if seq[depth] == matrix[aY][aX] {
			return true
		} else {
			return false
		}
	}

	// base: next not valid
	if matrix[aY][aX] != seq[depth] {
		return false
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
