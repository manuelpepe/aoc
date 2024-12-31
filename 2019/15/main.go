package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"2019/utils/machine/parser"
	"2019/utils/machine/vm"
	"2019/utils/machine/vmio"
)

var SHOW_GRID = false

const SHOW_GRID_DELAY = 100 * time.Millisecond

func main() {
	playflag := flag.Bool("play", false, "play the game")
	gridflag := flag.Bool("show", false, "show the grid on the simulation")
	flag.CommandLine.Parse(os.Args[2:])

	SHOW_GRID = *gridflag

	program := parser.Parse(os.Args[1])

	if *playflag {
		play(program)
	} else {
		sol1(program)
		sol2(program)
	}
}

func sol1(program []int) {
	m, inbuf, outbuf := vm.NewVMPiped(program)
	sendInput, stop := vmio.StartInputSender(inbuf)
	defer stop()

	sendSingleByte := func(b byte) {
		sendInput([]byte{'0' + b})
	}

	playerPos := Pos{30, 30}

	grid := NewGrid(60, 60)
	grid.Set(playerPos, PLAYER)

	minSteps := DFS(m, grid, playerPos, 0, sendSingleByte, outbuf)
	fmt.Printf("solution 1: %d\n", minSteps)
}

func DFS(m *vm.VM, grid *Grid, playerPos Pos, curSteps int, sendInput func(b byte), outbuf *bytes.Buffer) int {
	curSteps++
	minSteps := -1

	reverseDirs := []byte{2, 1, 4, 3}

	for dir := range 4 {
		checkingPos := playerPos.Add(DIRS[dir])

		if grid.Visited(checkingPos) || grid.Wall(checkingPos) {
			continue // skip visited tiles
		}

		inputDir := byte(dir + 1)
		sendInput(inputDir)
		m.RunForOutput()
		res := vmio.GetLastOutput(outbuf)

		switch res {
		case HIT_WALL:
			grid.Set(checkingPos, WALL)
			showgrid(grid, inputDir, playerPos, checkingPos, curSteps, "HIT WALL")

		case MOVED:
			grid.Set(playerPos, STEPPED)
			grid.Set(checkingPos, PLAYER)

			showgrid(grid, inputDir, playerPos, checkingPos, curSteps, "MOVED")

			steps := DFS(m, grid, checkingPos, curSteps, sendInput, outbuf)
			if steps > -1 { // found goal down this path
				if minSteps == -1 || minSteps > steps {
					minSteps = steps
				}
			}

			// undo movement on backtrack
			inputDir := reverseDirs[dir]
			sendInput(inputDir)
			m.RunForOutput()
			res := vmio.GetLastOutput(outbuf)

			if res != MOVED {
				panic(fmt.Sprintf("unable to undo move. expected moved, got %d", res))
			}

			grid.Set(checkingPos, STEPPED)
			grid.Set(playerPos, PLAYER)
			showgrid(grid, inputDir, checkingPos, playerPos, curSteps, "BACKTRACKED")

		case FOUND:
			grid.Set(playerPos, STEPPED)
			grid.Set(checkingPos, TANK)
			showgrid(grid, inputDir, playerPos, checkingPos, curSteps, "OXYGEN TANK")

			// undo movement on backtrack
			inputDir := reverseDirs[dir]
			sendInput(inputDir)
			m.RunForOutput()
			res := vmio.GetLastOutput(outbuf)

			if res != MOVED {
				panic(fmt.Sprintf("unable to undo move. expected moved, got %d", res))
			}

			grid.Set(playerPos, PLAYER)
			showgrid(grid, inputDir, checkingPos, playerPos, curSteps, "BACKTRACKED")

			return curSteps // steps to goal
		}
	}

	return minSteps
}

func sol2(program []int) {
	m, inbuf, outbuf := vm.NewVMPiped(program)
	sendInput, stop := vmio.StartInputSender(inbuf)
	defer stop()

	sendSingleByte := func(b byte) {
		sendInput([]byte{'0' + b})
	}

	playerPos := Pos{30, 30}

	grid := NewGrid(60, 60)
	grid.Set(playerPos, PLAYER)

	DFS(m, grid, playerPos, 0, sendSingleByte, outbuf)

	parsedGrid := grid.ExtractMap()
	printGrid(parsedGrid)

	// TODO: Find longest path to closed section starting from TANK
}

func play(program []int) {
	m, inbuf, outbuf := vm.NewVMPiped(program)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	inputs := watch(ctx, inbuf)

	playerPos := Pos{20, 20}
	tankPos := Pos{-1, -1}

	grid := NewGrid(40, 40)
	grid.Set(playerPos, PLAYER)

	_ = []int{0, 2, 1, 3}

	for {
		m.RunForOutput()
		res := vmio.GetLastOutput(outbuf)

		currentDirection := <-inputs - '0'

		checkingPos := playerPos.Add(DIRS[currentDirection-1])

		switch res {
		case HIT_WALL:
			grid.Set(checkingPos, WALL)
		case MOVED:
			grid.Set(checkingPos, PLAYER)
			if playerPos == tankPos {
				grid.Set(playerPos, TANK)
			} else {
				grid.Set(playerPos, STEPPED)
			}
			playerPos = checkingPos
		case FOUND:
			grid.Set(checkingPos, TANK)
			grid.Set(playerPos, STEPPED)
			playerPos = checkingPos
			tankPos = checkingPos
		}

		clearScreen()
		printGrid(grid.data)
		time.Sleep(100 * time.Millisecond)
	}

}

var DIRS = []Pos{
	{0, -1},
	{0, 1},
	{-1, 0},
	{1, 0},
}

const (
	HIT_WALL = 0
	MOVED    = 1
	FOUND    = 2
)

const (
	NORTH = 1
	SOUTH = 2
	WEST  = 3
	EAST  = 4
)

const (
	EMPTY   = ' '
	WALL    = '#'
	PLAYER  = 'D'
	TANK    = '$'
	STEPPED = '.'
)

// An X-Y grid that grows automatically.
// Also keeps track of written bounds to extract the relevant parts of the grid.
type Grid struct {
	data [][]byte

	minWritten Pos
	maxWritten Pos
}

func NewGrid(x, y int) *Grid {
	g := &Grid{
		data:       make([][]byte, y),
		minWritten: Pos{9999, 9999},
	}
	for i := range g.data {
		g.data[i] = make([]byte, x)
		for h := range g.data[i] {
			g.data[i][h] = EMPTY
		}
	}
	return g
}

func (g *Grid) Set(pos Pos, val byte) {
	if pos.X >= len(g.data[0]) {
		g.growX(pos.X)
	}
	if pos.Y >= len(g.data) {
		g.growY(pos.Y)
	}

	g.minWritten.X = min(g.minWritten.X, pos.X)
	g.maxWritten.X = max(g.maxWritten.X, pos.X)

	g.minWritten.Y = min(g.minWritten.Y, pos.Y)
	g.maxWritten.Y = max(g.maxWritten.Y, pos.Y)

	g.data[pos.Y][pos.X] = val
}

func (g *Grid) growY(upto int) {
	for len(g.data) <= upto {
		newRow := make([]byte, len(g.data[0]))
		g.data = append(g.data, newRow)
	}
}

func (g *Grid) growX(upto int) {
	for y := range g.data {
		for len(g.data[y]) < upto+1 {
			g.data[y] = append(g.data[y], EMPTY)
		}
	}
}

func (g *Grid) Visited(p Pos) bool {
	return p.Y >= 0 && p.X >= 0 && p.Y < len(g.data) && p.X < len(g.data[p.Y]) && g.data[p.Y][p.X] == STEPPED
}

func (g *Grid) Wall(p Pos) bool {
	return p.Y >= 0 && p.X >= 0 && p.Y < len(g.data) && p.X < len(g.data[p.Y]) && g.data[p.Y][p.X] == WALL
}

func (g *Grid) ExtractMap() [][]byte {
	out := make([][]byte, 0)
	for y := g.minWritten.Y; y <= g.maxWritten.Y; y++ {
		row := make([]byte, 0)
		for x := g.minWritten.X; x <= g.maxWritten.X; x++ {
			row = append(row, g.data[y][x])
		}
		out = append(out, row)
	}
	return out
}

// watch reads data from stdin and sends it to the given input buffer in real time
// it only watches for the arrow buttons (up, down, left, right) and translates the key codes
// to the corresponding direction inputs for the arcade machine (1,2,3,4).
//
// It returns a read-only channel used to read the machine inputs.
func watch(ctx context.Context, w io.Writer) <-chan byte {
	inps := make(chan []byte)
	out := make(chan byte)

	// read goroutine reads 3 byte at a time from stdin and sends it to an inputs channel
	go func() {
		// disable input buffering
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()

		data := make([]byte, 3)
		for {
			_, err := os.Stdin.Read(data)
			if err != nil {
				err = fmt.Errorf("error reading stdin: %w", err)
				panic(err)
			}
			inps <- data
		}
	}()

	// process goroutine reads the inputs channel and sends its data to the writer
	go func() {
		ARROW_CODES_TO_MOVES := map[byte]byte{
			65: 1,
			66: 2,
			67: 4,
			68: 3,
		}

		for {
			select {
			case b := <-inps:
				if b[0] != 27 || b[1] != 91 {
					continue
				}

				val := '0' + ARROW_CODES_TO_MOVES[b[2]]
				if _, err := w.Write([]byte{val, '\n'}); err != nil {
					panic(err)
				}

				select {
				case out <- val:
				case <-ctx.Done():
					return
				}

			case <-ctx.Done():
				fmt.Printf("done reading stdin\n")
				return
			}
		}
	}()

	return out
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

func showgrid(grid *Grid, dir byte, playerPos Pos, checkPos Pos, steps int, msg string) {
	if !SHOW_GRID {
		return
	}

	// fmt.Printf("Step: %d\n", steps)
	// fmt.Printf("Direction: %c - [ (%d, %d) -> (%d, %d) ]\n", dir, playerPos.X, playerPos.Y, checkPos.X, checkPos.Y)
	// fmt.Printf("%s\n", msg)
	printGrid(grid.data)
	// fmt.Println("==================")

	time.Sleep(SHOW_GRID_DELAY)
}

func printGrid(grid [][]byte) {
	for _, row := range grid {
		for _, cell := range row {
			fmt.Printf("%c ", cell)
		}
		fmt.Printf("\n")
	}
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
