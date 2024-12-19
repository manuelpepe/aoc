package main

import (
	"16/pq"
	"bufio"
	"fmt"
	"os"
	"slices"
)

func main() {

	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	grid, start, end := parse(file)
	fmt.Printf("cost of shortest path is: %d\n", pathCost(grid, start, end))
}

func parse(file *os.File) ([][]byte, pos, pos) {
	r := bufio.NewScanner(file)

	grid := make([][]byte, 0)
	var start, end pos

	y := 0
	for r.Scan() {
		line := []byte(r.Text())
		for x := range line {
			if line[x] == 'S' {
				start = pos{x, y}
			} else if line[x] == 'E' {
				end = pos{x, y}
			}
		}
		y++
		grid = append(grid, line)
	}

	return grid, start, end
}

type pos struct {
	x, y int
}

func (p pos) Sub(p2 pos) pos {
	return pos{
		p.x - p2.x,
		p.y - p2.y,
	}
}

func (p pos) Add(p2 pos) pos {
	return pos{
		p.x + p2.x,
		p.y + p2.y,
	}
}

// find shortest path from start to goal on grid using A*
func pathCost(grid [][]byte, start pos, goal pos) int {
	frontier := pq.New[pos](pq.MIN_HEAP)
	frontier.Push(start, 0)

	cameFrom := make(map[pos]pos)
	costTo := make(map[pos]int)

	cameFrom[start] = pos{-1, -1}
	costTo[start] = 0

	for {
		current, ok := frontier.Pop()
		if !ok {
			break
		}
		if current == goal {
			break
		}

		for _, next := range neighs(grid, current) {
			newCostToNext := costTo[current] + costForMove(cameFrom[current], current, next)
			curCostToNext, ok := costTo[next]
			if !ok || newCostToNext < curCostToNext {
				costTo[next] = newCostToNext
				cameFrom[next] = current

				priority := newCostToNext + manhattanDistance(goal, next)
				frontier.Push(next, priority)
			}
		}
	}

	return costTo[goal]
}

func costForMove(last, from, to pos) int {
	if last == (pos{-1, -1}) {
		last = from.Sub(pos{1, 0}) // starts pointing east
	}
	last_dir := from.Sub(last)
	cur_dir := to.Sub(from)

	if last_dir == cur_dir {
		return 1
	}

	lastDirIx := slices.Index(DIRS, last_dir)
	nextDirIx := slices.Index(DIRS, cur_dir)

	return (abs(lastDirIx-nextDirIx)%2)*1000 + 1
}

var DIRS = []pos{
	{0, -1},
	{1, 0},
	{0, 1},
	{-1, 0},
}

func neighs(grid [][]byte, p pos) []pos {
	out := make([]pos, 0)
	for _, dir := range DIRS {
		new := p.Add(dir)
		if inBounds(grid, new) && grid[new.y][new.x] != '#' {
			out = append(out, new)
		}
	}
	return out
}

func inBounds(matrix [][]byte, p pos) bool {
	if p.y < 0 || p.y >= len(matrix) {
		return false
	}

	if p.x < 0 || p.x >= len(matrix[0]) {
		return false
	}

	return true
}

func manhattanDistance(a, b pos) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func abs[T ~int](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func printTable(table [][]byte, frontier *pq.PriorityQueue[pos]) {
	itemsInFrontier := frontier.UnordereredItems()

	for y, row := range table {
		for x, cell := range row {
			if slices.Contains(itemsInFrontier, pos{x, y}) {
				cell = 'O'
			}
			fmt.Printf("%c ", cell)
		}
		fmt.Printf("\n")
	}

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
