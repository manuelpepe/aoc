package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/manuelpepe/aoc/utils/astar"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	limit, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		panic(err)
	}

	grid, start, end := parse(file)
	fmt.Printf("there are %d cheats that would save you %d picoseconds\n", sol1(grid, start, end, int(limit)), limit)
}

func parse(file *os.File) ([][]byte, astar.Pos, astar.Pos) {
	r := bufio.NewScanner(file)

	grid := make([][]byte, 0)
	var start, end astar.Pos

	y := 0
	for r.Scan() {
		line := []byte(r.Text())
		for x := range line {
			if line[x] == 'S' {
				start = astar.Pos{x, y}
			} else if line[x] == 'E' {
				end = astar.Pos{x, y}
			}
		}
		y++
		grid = append(grid, line)
	}

	return grid, start, end
}

func sol1(grid [][]byte, start, end astar.Pos, threshold int) int {
	_, bestCost := astar.Find(grid, start, end, func(last, from, to astar.Pos) int {
		return 1
	})

	acc := 0

	for y := 1; y < len(grid)-1; y++ {
		for x := 1; x < len(grid[0])-1; x++ {
			if grid[y][x] != '#' {
				continue
			}

			grid[y][x] = 0

			_, newBest := astar.Find(grid, start, end, func(last, from, to astar.Pos) int {
				return 1
			})

			grid[y][x] = '#'

			if bestCost-newBest >= threshold {
				acc += 1
			}

		}
	}

	return acc
}
