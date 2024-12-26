package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	grid, moves := parse(os.Args[1])
	// printGrid(grid)
	// fmt.Printf("Grid: %+v\nMoves: %+v\n", grid, moves)

	res := sol1(grid, moves)
	printGrid(grid)
	fmt.Printf("sum of GPS for all boxes is %d\n", res)
}

func sol1(grid Grid, moves Moves) int {
	robotPos := findRobot(grid)

	for _, move := range moves {
		// fmt.Printf("trying move #%d : %c\n", ix, move)
		robotPos = tryMove(grid, robotPos, move)
		// printGrid(grid)
	}

	acc := 0

	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid); x++ {
			if grid[y][x] == Box {
				acc += 100*y + x
			}
		}
	}

	return acc
}

func tryMove(grid Grid, robotPos Pos, move Move) Pos {
	dir := DIRS[move]

	nextEmpty, found := emptyInDirection(grid, robotPos, dir)
	if !found {
		// fmt.Printf("  no empty move found\n")
		return robotPos // cant make move
	}

	// fmt.Printf("  making move to %v\n", nextEmpty)
	return makeMove(grid, robotPos, dir, nextEmpty)
}

func makeMove(grid Grid, robotPos Pos, dir Direction, upTo Pos) Pos {
	curPosition := upTo
	var nextPosition Pos
	for {
		nextPosition = curPosition.Sub(dir)

		tmp := grid[curPosition.Y][curPosition.X]
		grid[curPosition.Y][curPosition.X] = grid[nextPosition.Y][nextPosition.X]
		grid[nextPosition.Y][nextPosition.X] = tmp
		if nextPosition == robotPos {
			break
		}

		curPosition = nextPosition
	}

	return curPosition
}

func emptyInDirection(grid Grid, from Pos, dir Direction) (Pos, bool) {
	cur := from

	for {
		cur = cur.Add(dir)

		if !grid.InBounds(cur) {
			break
		}

		if grid[cur.Y][cur.X] == Wall {
			break
		}

		if grid[cur.Y][cur.X] == Empty {
			return cur, true
		}
	}

	return Pos{}, false
}

func findRobot(grid Grid) Pos {
	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] == Robot {
				return Pos{x, y}
			}
		}
	}
	panic("robot not found")
}

func parse(fn string) (Grid, Moves) {
	fh, err := os.Open(fn)
	if err != nil {
		panic(err)
	}

	grid := make(Grid, 0)
	moves := make(Moves, 0)

	parsingGrid := true

	s := bufio.NewScanner(fh)
	for s.Scan() {
		line := []byte(s.Text())

		if len(line) == 0 {
			parsingGrid = false
			continue
		}

		if parsingGrid {
			grid = append(grid, line)
			continue
		}

		moves = append(moves, line...)
	}

	return grid, moves
}

type Item = byte

const (
	Robot Item = '@'
	Box   Item = 'O'
	Wall  Item = '#'
	Empty Item = '.'
)

type Move = byte

const (
	Up    Move = '^'
	Down  Move = 'v'
	Left  Move = '<'
	Right Move = '>'
)

type Grid [][]byte

func (g Grid) InBounds(pos Pos) bool {
	if pos.Y < 0 || pos.Y >= len(g) {
		return false
	}

	if pos.X < 0 || pos.X >= len(g[0]) {
		return false
	}

	return true
}

type Moves []Move

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

type Direction = Pos

var DIRS = map[Move]Direction{
	Down:  {0, 1},
	Left:  {-1, 0},
	Up:    {0, -1},
	Right: {1, 0},
}

func printGrid(table Grid) {
	for _, row := range table {
		for _, cell := range row {
			fmt.Printf("%c ", cell)
		}
		fmt.Printf("\n")
	}

	// reader := bufio.NewReader(os.Stdin)
	// reader.ReadString('\n')
}
