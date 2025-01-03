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
	sol2(program)
}

func sol1(program []int) {
	m, inbuf, outbuf := vm.NewVMPiped(program)
	sendInput, stop := vmio.StartInputSender(inbuf)
	defer stop()

	grid := make([][]byte, 0)

	acc := 0
	for y := 0; y < 50; y++ {
		row := make([]byte, 0)
		for x := 0; x < 50; x++ {
			res := tryCoordinates(m, sendInput, outbuf, Pos{x, y})
			acc += res
			row = append(row, byte(res))
		}
		grid = append(grid, row)
	}

	// printGrid(grid)

	fmt.Printf("solution 1: %d\n", acc)
}

var cache = make(map[Pos]int)

func tryCoordinates(m *vm.VM, sendInput vmio.SendInputFunc, outbuf *bytes.Buffer, p Pos) int {
	if _, ok := cache[p]; !ok {
		m.Restart()
		data := fmt.Sprintf("%d\n%d\n", p.X, p.Y)
		sendInput([]byte(data))
		m.RunForOutput()
		cache[p] = vmio.GetLastOutput(outbuf)
	}
	return cache[p]
}

func sol2(program []int) {
	m, inbuf, outbuf := vm.NewVMPiped(program)
	sendInput, stop := vmio.StartInputSender(inbuf)
	defer stop()

	acc := 0

	queued := make([]Pos, 0)
	queued = append(queued, Pos{0, 0})

	visited := make(map[Pos]bool)
	visited[queued[0]] = true

	for len(queued) > 0 {
		cur := queued[0]
		queued = queued[1:]

		if squareAt(m, sendInput, outbuf, cur.X, cur.Y) {
			acc = cur.X*10000 + cur.Y
			break
		}

		for _, neigh := range neighbors(cur) {
			if visited[neigh] {
				continue
			}

			if tryCoordinates(m, sendInput, outbuf, neigh) == 0 {
				continue
			}

			queued = append(queued, neigh)
			visited[neigh] = true
		}
	}

	fmt.Printf("solution 2: %d\n", acc)
}

func squareAt(m *vm.VM, sendInput vmio.SendInputFunc, outbuf *bytes.Buffer, x, y int) bool {
	for cy := y; cy < y+100; cy++ {
		if tryCoordinates(m, sendInput, outbuf, Pos{x, cy}) == 0 {
			return false
		}
	}
	for cx := x; cx < x+100; cx++ {
		if tryCoordinates(m, sendInput, outbuf, Pos{cx, y}) == 0 {
			return false
		}
	}
	return true
}

var DIRS = []Pos{
	{1, 0},
	{1, 1},
	{0, 1},
}

func neighbors(p Pos) []Pos {
	out := make([]Pos, 0)

	for _, dir := range DIRS {
		new := p.Add(dir)
		if !inBounds(new) {
			continue
		}
		out = append(out, new)
	}

	return out
}

func inBounds(pos Pos) bool {
	return pos.Y >= 0 && pos.X >= 0
}

type Pos struct{ X, Y int }

func (p Pos) Add(p2 Pos) Pos {
	return Pos{
		p.X + p2.X,
		p.Y + p2.Y,
	}
}

func printGrid(grid [][]byte) {
	for _, row := range grid {
		for _, cell := range row {
			fmt.Printf("%d ", cell)
		}
		fmt.Printf("\n")
	}
}
