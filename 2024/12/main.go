package main

import (
	"bufio"
	"fmt"
	"os"
)

type Plot struct {
	Plant   byte
	Visited bool
}

func main() {
	fh, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	table := make([][]*Plot, 0)

	r := bufio.NewScanner(fh)

	for r.Scan() {
		line := []byte(r.Text())

		row := make([]*Plot, len(line))
		for ix, p := range line {
			row[ix] = &Plot{Plant: p, Visited: false}
		}

		table = append(table, row)
	}

	fmt.Printf("total price of land is: %d\n", sol1(table))
}

func sol1(table [][]*Plot) int {
	acc := 0

	for y := range table {
		for x := range table[0] {
			if table[y][x].Visited {
				continue
			}

			area, perim := checkRegion(table, x, y)

			acc += area * perim
		}
	}

	return acc
}

var DIRS = [][]int{
	{-1, 0},
	{1, 0},
	{0, 1},
	{0, -1},
}

func checkRegion(table [][]*Plot, x, y int) (int, int) {
	plot := table[y][x]
	plot.Visited = true

	area, perim := 1, 0

	for _, dir := range DIRS {
		nx, ny := x+dir[0], y+dir[1]

		if !inBounds(table, nx, ny) {
			perim++
			continue
		}

		nplot := table[ny][nx]

		if nplot.Plant != plot.Plant {
			perim++
			continue
		}

		if nplot.Visited {
			continue
		}

		narea, nperim := checkRegion(table, nx, ny)
		area += narea
		perim += nperim
	}

	return area, perim
}

func inBounds(table [][]*Plot, x, y int) bool {
	if y < 0 || y >= len(table) {
		return false
	}

	if x < 0 || x >= len(table[0]) {
		return false
	}

	return true
}
