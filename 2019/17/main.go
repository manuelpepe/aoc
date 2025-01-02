package main

import (
	"bytes"
	"fmt"
	"os"

	"2019/utils/machine/parser"
	"2019/utils/machine/vm"
	"2019/utils/machine/vmio"
)

func main() {
	program := parser.Parse(os.Args[1])
	sol1(program)
}

func sol1(program []int) {
	m, _, outbuf := vm.NewVMPiped(program)
	m.Run()

	grid := parseGrid(outbuf)
	intersections := make([]Pos, 0)

	for y := range grid {
	ROW:
		for x := range grid[y] {
			if grid[y][x] != '#' {
				continue
			}

			neighs := neighbors(grid, Pos{x, y})
			if len(neighs) != 4 {
				continue
			}

			for _, n := range neighs {
				if grid[n.Y][n.X] != '#' {
					continue ROW
				}
			}

			intersections = append(intersections, Pos{x, y})
		}
	}

	for y := range grid {
		fmt.Printf("> %s\n", string(grid[y]))
	}
	fmt.Printf("intersections: %+v\n", intersections)

	acc := 0
	for _, inter := range intersections {
		acc += inter.X * inter.Y
	}

	fmt.Printf("solution 1: %d\n", acc)
}

func parseGrid(outbuf *bytes.Buffer) [][]byte {
	grid := make([][]byte, 0)
	row := make([]byte, 0)
	for _, c := range vmio.GetOutput(outbuf) {
		if c == '\n' {
			grid = append(grid, row)
			row = make([]byte, 0)
		} else {
			row = append(row, byte(c))
		}
	}
	return grid[:len(grid)-1]
}

var DIRS = []Pos{
	{0, -1},
	{1, 0},
	{0, 1},
	{-1, 0},
}

func neighbors(grid [][]byte, p Pos) []Pos {
	out := make([]Pos, 0)

	for _, dir := range DIRS {
		new := p.Add(dir)

		if !inBounds(grid, new) {
			continue
		}

		out = append(out, new)
	}

	return out
}

func neighPositions(p Pos) []Pos {
	out := make([]Pos, len(DIRS))
	for _, dir := range DIRS {
		out = append(out, p.Add(dir))
	}
	return out
}

func inBounds(table [][]byte, pos Pos) bool {
	if pos.Y < 0 || pos.Y >= len(table) {
		return false
	}

	if pos.X < 0 || pos.X >= len(table[0]) {
		return false
	}

	return true
}

type Pos struct{ X, Y int }

func (p Pos) Add(p2 Pos) Pos {
	return Pos{
		p.X + p2.X,
		p.Y + p2.Y,
	}
}
