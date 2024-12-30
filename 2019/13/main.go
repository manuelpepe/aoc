package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"2019/utils/machine/parser"
	"2019/utils/machine/vm"
	"2019/utils/machine/vmio"
)

func main() {
	program := parser.Parse(os.Args[1])
	// play(program)
	sol1(program)
	sol2(program)
}

func sol2(program []int) {
	program[0] = 2
	m, tiles, inbuf, outbuf := prepareArcade(program)
	score := runArcade(m, tiles, inbuf, outbuf, false)
	fmt.Printf("solution 2: %d\n", score)
}

func sol1(program []int) {
	_, tiles, _, _ := prepareArcade(program)

	blocks := 0
	for y := range tiles {
		for x := range tiles[y] {
			if tiles[y][x] == 2 {
				blocks++
			}
		}
	}

	fmt.Printf("solution 1: %d\n", blocks)
}

func prepareArcade(program []int) (*vm.VM, [][]byte, io.Writer, *bytes.Buffer) {
	inReader, inWriter := io.Pipe()
	outbuf := bytes.NewBuffer([]byte{})

	m := vm.NewVM(program, inReader, outbuf)

	grid := make([][]byte, 24)
	for y := range grid {
		grid[y] = make([]byte, 36)
	}

	for {
		tile, halt := readTile(m, outbuf)
		if halt {
			break
		}

		if tile.X == -1 && tile.Y == 0 {
			break
		}

		grid[tile.Y][tile.X] = byte(tile.Value)
	}

	return m, grid, inWriter, outbuf

}

func runArcade(m *vm.VM, grid [][]byte, instream io.Writer, outstream *bytes.Buffer, display bool) int {
	// locate ball and player
	var ballPos, playerPos Pos
	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] == Ball {
				ballPos = Pos{x, y}
			}
			if grid[y][x] == Paddle {
				playerPos = Pos{x, y}
			}
		}
	}

	// start input goroutine
	ctx, stopInputGoroutine := context.WithCancel(context.Background())
	inputs := make(chan string, 10)

	go func(ch chan string) {
		for {
			select {
			case v := <-ch:
				if _, err := instream.Write([]byte(v)); err != nil {
					panic(err)
				}
			case <-ctx.Done():
				return
			}
		}
	}(inputs)

	// play game
	score := 0

	nextInput := "0\n"
	inputs <- nextInput // send first input

	for !m.Halted() {
		nextInput = "0\n"

		tile, halt := readTile(m, outstream)
		if halt {
			break
		}

		if tile.X == -1 && tile.Y == 0 {
			score = tile.Value
		} else {
			grid[tile.Y][tile.X] = byte(tile.Value)

			switch tile.Value {
			case int(Paddle):
				playerPos = Pos{tile.X, tile.Y}

			case int(Ball):
				ballPos = Pos{tile.X, tile.Y}
				if ballPos.X > playerPos.X {
					nextInput = "1\n"
				} else if ballPos.X < playerPos.X {
					nextInput = "-1\n"
				}
				inputs <- nextInput
			}

			if display {
				printGrid(grid, score)
			}
		}

	}

	stopInputGoroutine()
	return score
}

func readTile(m *vm.VM, outbuf *bytes.Buffer) (Tile, bool) {
	m.RunForOutput()
	if m.Halted() {
		return Tile{}, true
	}
	m.RunForOutput()
	m.RunForOutput()

	outs := vmio.GetOutput(outbuf)
	if len(outs) != 3 {
		panic(fmt.Sprintf("expected two outputs (color and rotation) got: %+v", outs))
	}

	return Tile{X: outs[0], Y: outs[1], Value: outs[2]}, false
}

const (
	Empty  byte = 0
	Wall   byte = 1
	Block  byte = 2
	Paddle byte = 3
	Ball   byte = 4
)

type Pos struct {
	X, Y int
}

type Tile struct {
	X, Y  int
	Value int
}

func printGrid(grid [][]byte, score int) {
	for _, row := range grid {
		for _, cell := range row {
			fmt.Printf("%c ", getShowValue(cell))
		}
		fmt.Printf("\n")
	}
	fmt.Printf("Score: %d\n", score)
}

var showvalues = map[byte]byte{
	0: ' ',
	1: '#',
	2: 'B',
	3: '-',
	4: 'o',
}

func getShowValue(v byte) byte {
	return showvalues[v]
}

//// below this is outside of the exercise scope but I wanted to try it nonetheless.. ////

func play(program []int) {
	program[0] = 2
	m, tiles, inbuf, outbuf := prepareArcade(program)

	ctx, cancel := context.WithCancel(context.Background())
	watch(ctx, inbuf)

	score := runArcade2(m, tiles, outbuf)
	fmt.Printf("solution 2: %d\n", score)

	cancel()
}

func runArcade2(m *vm.VM, grid [][]byte, outstream *bytes.Buffer) int {
	score := 0

	for !m.Halted() {
		tile, halt := readTile(m, outstream)
		if halt {
			break
		}

		if tile.X == -1 && tile.Y == 0 {
			score = tile.Value
		} else {
			grid[tile.Y][tile.X] = byte(tile.Value)
			printGrid(grid, score)
		}

	}

	return score
}

// used this function while testing to play the game myself
// not the best experience
func watch(ctx context.Context, w io.Writer) {
	t := time.NewTicker(1 * time.Second)
	inps := make(chan byte)

	// read goroutine reads 1 byte at a time from stdin and sends it to an inputs channel
	go func() {
		// disable input buffering
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()

		data := make([]byte, 1)
		for {
			_, err := os.Stdin.Read(data)
			if err != nil {
				err = fmt.Errorf("error reading stdin: %w", err)
				panic(err)
			}
			inps <- data[0]
		}
	}()

	// process goroutine reads the inputs channel and sends its data to the writer
	// if no input is provided it sends '0' once per second
	go func() {
		for {
			select {
			case <-t.C:
				if _, err := w.Write([]byte{'0', '\n'}); err != nil {
					panic(err)
				}

			case b := <-inps:
				// -1 is bothersome to input
				// shifted inputs +2 so [1, 2, 3] = [-1, 0, 1]
				actual := fmt.Sprintf("%d\n", int(b-'0')-2)
				if _, err := w.Write([]byte(actual)); err != nil {
					panic(err)
				}
				t.Reset(1000 * time.Millisecond) // reset for smoothness

			case <-ctx.Done():
				fmt.Printf("done reading stdin\n")
				return
			}
		}
	}()
}
