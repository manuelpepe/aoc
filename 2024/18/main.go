package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/manuelpepe/aoc/utils/astar"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	size, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		panic(err)
	}

	limit, err := strconv.ParseInt(os.Args[3], 10, 64)
	if err != nil {
		panic(err)
	}

	scanner, grid := parse(file, int(size), int(limit))

	fmt.Printf("shortest path is %d steps\n", sol1(grid))

	x, y := sol2(grid, scanner)
	fmt.Printf("first byte to block exit is %d,%d\n", x, y)
}

func parse(file *os.File, size, readto int) (*bufio.Scanner, [][]byte) {
	grid := make([][]byte, size)
	for i := 0; i < size; i++ {
		grid[i] = make([]byte, size)
		for j := range grid[i] {
			grid[i][j] = '.' // for debug
		}
	}

	r := bufio.NewScanner(file)
	cur := 0
	for r.Scan() {
		x, y := parseLine(r)
		grid[y][x] = '#'

		cur++
		if cur == readto {
			break
		}
	}

	return r, grid
}

func parseLine(scanner *bufio.Scanner) (int, int) {
	line := strings.Split(scanner.Text(), ",")
	if len(line) != 2 {
		panic("expected 2 values")
	}

	x, err := strconv.ParseInt(line[0], 10, 64)
	if err != nil {
		panic(err)
	}

	y, err := strconv.ParseInt(line[1], 10, 64)
	if err != nil {
		panic(err)
	}

	return int(x), int(y)
}

func sol1(grid [][]byte) int {
	from := astar.Pos{0, 0}
	to := astar.Pos{len(grid) - 1, len(grid) - 1}

	path, _ := astar.Find(grid, from, to, func(_, _, _ astar.Pos) int {
		return 1
	})

	return len(path) - 1

}

func sol2(grid [][]byte, scanner *bufio.Scanner) (int, int) {
	for {
		scanner.Scan()
		x, y := parseLine(scanner)
		grid[y][x] = '#'

		cost := sol1(grid)
		if cost == 1 {
			return x, y
		}
	}
}
