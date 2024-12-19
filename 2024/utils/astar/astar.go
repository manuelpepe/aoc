package astar

import (
	"bufio"
	"fmt"
	"os"
	"slices"

	"github.com/manuelpepe/aoc/utils/pq"
)

var DEBUG = false

// Basic (X,Y) position
type Pos struct {
	X, Y int
}

func (p Pos) Sub(p2 Pos) Pos {
	return Pos{
		p.X - p2.X,
		p.Y - p2.Y,
	}
}

func (p Pos) Add(p2 Pos) Pos {
	return Pos{
		p.X + p2.X,
		p.Y + p2.Y,
	}
}

type CostFunction func(last, from, to Pos) int

// An example cost function (from adv 2024 ex.16) that adds 1 per step,
// but 1000 for every turn (clockwise or counter-clockwise)
func Straight1Turns1000Bothwise(last, from, to Pos) int {
	if last == (Pos{-1, -1}) {
		last = from.Sub(Pos{1, 0}) // starts pointing east
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

// Find the shortest path from start to goal on grid using A*
func Find(grid [][]byte, start Pos, goal Pos, costFn CostFunction) ([]Pos, int) {
	frontier := pq.New[Pos](pq.MIN_HEAP)
	frontier.Push(start, 0)

	cameFrom := make(map[Pos]Pos)
	costTo := make(map[Pos]int)

	cameFrom[start] = Pos{-1, -1}
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
			newCostToNext := costTo[current] + costFn(cameFrom[current], current, next)
			curCostToNext, ok := costTo[next]
			if !ok || newCostToNext < curCostToNext {
				costTo[next] = newCostToNext
				cameFrom[next] = current

				priority := newCostToNext + manhattanDistance(goal, next)
				frontier.Push(next, priority)
				if DEBUG {
					printTable(grid, frontier)
				}
			}
		}

	}

	path := make([]Pos, 0)
	current := goal
	for current != start {
		path = append(path, current)
		current = cameFrom[current]
	}
	path = append(path, start)

	return path, costTo[goal]
}

var DIRS = []Pos{
	{0, -1},
	{1, 0},
	{0, 1},
	{-1, 0},
}

// TODO: could be parametrized
func neighs(grid [][]byte, p Pos) []Pos {
	out := make([]Pos, 0)
	for _, dir := range DIRS {
		new := p.Add(dir)
		if inBounds(grid, new) && grid[new.Y][new.X] != '#' {
			out = append(out, new)
		}
	}
	return out
}

func inBounds(matrix [][]byte, p Pos) bool {
	if p.Y < 0 || p.Y >= len(matrix) {
		return false
	}

	if p.X < 0 || p.X >= len(matrix[0]) {
		return false
	}

	return true
}

// TODO: could be parametrized
func manhattanDistance(a, b Pos) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y)
}

func abs[T ~int](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func printTable(table [][]byte, frontier *pq.PriorityQueue[Pos]) {
	itemsInFrontier := frontier.UnordereredItems()

	for y, row := range table {
		for x, cell := range row {
			if slices.Contains(itemsInFrontier, Pos{x, y}) {
				cell = 'O'
			}
			fmt.Printf("%c ", cell)
		}
		fmt.Printf("\n")
	}

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
