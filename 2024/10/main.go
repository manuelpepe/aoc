package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

func main() {
	fh, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	table := make([][]byte, 0)

	r := bufio.NewScanner(fh)
	for r.Scan() {
		line := []byte(r.Text())
		for ix := range line {
			line[ix] -= '0'
		}
		table = append(table, line)
	}

	r1, r2 := sol1and2(table)

	fmt.Printf("found %d trailheads with a total rating of %d\n", r1, r2)
}

func sol1and2(table [][]byte) (int, int) {
	trailheads := 0
	rating := 0

	for y := range table {
		for x := range table[y] {
			if table[y][x] == 0 {
				score := findFrom(table, x, y, 0, x, y)
				rating += len(score)

				slices.SortFunc(score, func(a, b []int) int {
					if a[0] > b[0] {
						return 1
					} else if a[0] < b[0] {
						return -1
					} else {
						if a[1] > b[1] {
							return 1
						} else if a[1] < b[1] {
							return -1
						} else {
							return 0
						}
					}

				})

				lastPos := []int{-1, -1}
				for _, end := range score {
					if end[0] == lastPos[0] && end[1] == lastPos[1] {
						continue
					}
					trailheads++
					lastPos = end
				}
			}
		}
	}

	return trailheads, rating
}

var DIRS = [][]int{
	{-1, 0},
	{1, 0},
	{0, 1},
	{0, -1},
}

func findFrom(table [][]byte, x, y int, next byte, lastX int, lastY int) [][]int {
	if !inBounds(table, x, y) {
		return [][]int{}
	}
	if table[y][x] != next {
		return [][]int{}
	}
	if table[y][x] == 9 && next == 9 {
		return [][]int{{x, y}}
	}

	acc := make([][]int, 0)
	for _, dir := range DIRS {
		nx, ny := x+dir[0], y+dir[1]
		if nx == lastX && ny == lastY {
			continue
		}
		acc = append(acc, findFrom(table, nx, ny, next+1, x, y)...)
	}
	return acc
}

func inBounds(table [][]byte, x, y int) bool {
	if y < 0 || y >= len(table) {
		return false
	}

	if x < 0 || x >= len(table[0]) {
		return false
	}

	return true
}
