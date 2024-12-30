package main

import (
	"2019/utils/machine/parser"
	"2019/utils/machine/vm"
	"2019/utils/machine/vmio"
	"bytes"
	"fmt"
	"maps"
	"os"
	"slices"
)

func main() {
	program := parser.Parse(os.Args[1])
	sol1(program)
	sol2(program)
}

func sol2(program []int) {
	colors := runEHPR(program, 1)
	fmt.Printf("solution 2:\n")
	printGrid(colors)
}

func printGrid(colors map[Pos]byte) {
	positions := slices.Collect(maps.Keys(colors))

	maxX := slices.MaxFunc(positions, func(a, b Pos) int { return a.X - b.X })
	maxY := slices.MaxFunc(positions, func(a, b Pos) int { return a.Y - b.Y })

	grid := make([][]byte, maxY.Y+1)

	for y := range grid {
		grid[y] = make([]byte, maxX.X)
		for x := range grid[y] {
			if colors[Pos{x, y}] == 1 {
				grid[y][x] = '#'
			} else {
				grid[y][x] = '.'
			}
		}
	}

	for _, row := range grid {
		for _, cell := range row {
			fmt.Printf("%c ", cell)
		}
		fmt.Printf("\n")
	}
}

func sol1(program []int) {
	colors := runEHPR(program, 0)
	fmt.Printf("solution 1: %d\n", len(colors))
}

func runEHPR(program []int, startColor int) map[Pos]byte {
	inbuf := vmio.CreateInBuffer(startColor)
	outbuf := bytes.NewBuffer([]byte{})
	m := vm.NewVM(program, inbuf, outbuf)

	colors := make(map[Pos]byte)
	curpos := Pos{0, 0}
	curdir := 0

	for {
		m.RunForOutput()
		if m.Halted() {
			break
		}

		m.RunForOutput()

		outs := vmio.GetOutput(outbuf)
		if len(outs) != 2 {
			panic(fmt.Sprintf("expected two outputs (color and rotation) got: %+v", outs))
		}

		color, rotation := outs[0], outs[1]
		if rotation == 0 {
			curdir = mod(curdir-1, len(DIRS)) // left turn
		} else if rotation == 1 {
			curdir = (curdir + 1) % len(DIRS) // right turn
		} else {
			panic("unexpected next dir")
		}

		colors[curpos] = byte(color)               // store tile color
		curpos = curpos.Add(DIRS[curdir])          // move forward
		fmt.Fprintf(inbuf, "%d\n", colors[curpos]) // send next input
	}

	return colors
}

type Pos struct {
	X, Y int
}

func (p Pos) Add(p2 Pos) Pos {
	return Pos{
		p.X + p2.X,
		p.Y + p2.Y,
	}
}

var DIRS = []Pos{
	{0, -1},
	{1, 0},
	{0, 1},
	{-1, 0},
}

// https://stackoverflow.com/questions/43018206/modulo-of-negative-integers-in-go
// https://torstencurdt.com/tech/posts/modulo-of-negative-numbers/
// https://github.com/golang/go/issues/448
func mod(a, b int) int {
	return (a%b + b) % b
}
