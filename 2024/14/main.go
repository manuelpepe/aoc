package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
)

var PAT = regexp.MustCompile(`p=(\d+),(\d+) v=(-?\d+),(-?\d+)`)

type Robot struct {
	X, Y   int
	VX, VY int
}

func main() {
	fh, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	bots := make([]Robot, 0)

	r := bufio.NewScanner(fh)

	for r.Scan() {
		data := PAT.FindStringSubmatch(r.Text())
		parsed := make([]int, 4)
		for i := 0; i < 4; i++ {
			val, err := strconv.ParseInt(data[i+1], 10, 64)
			if err != nil {
				panic(err)
			}
			parsed[i] = int(val)
		}

		bots = append(bots, Robot{
			X:  parsed[0],
			Y:  parsed[1],
			VX: parsed[2],
			VY: parsed[3],
		})
	}

	// fmt.Printf("score is: %d\n", sol1(bots, 11, 7, 100))
	fmt.Printf("score is: %d\n", sol1(slices.Clone(bots), 101, 103, 100))
	fmt.Printf("score is: %d\n", sol1(bots, 101, 103, 100000))
}

// https://stackoverflow.com/questions/43018206/modulo-of-negative-integers-in-go
// https://torstencurdt.com/tech/posts/modulo-of-negative-numbers/
// https://github.com/golang/go/issues/448
func mod(a, b int) int {
	return (a%b + b) % b
}

func advance(bots []Robot, width, height, secs int) []Robot {
	for i := 0; i < secs; i++ {
		for ix := range bots {
			bots[ix].X = mod((bots[ix].X + bots[ix].VX), width)
			bots[ix].Y = mod((bots[ix].Y + bots[ix].VY), height)
		}
		if hasLine(bots, width, height) {
			printTable(bots, width, height, i)

		}
	}
	return bots
}

type pos struct {
	x, y int
}

func sol1(bots []Robot, width, height, secs int) int {
	advance(bots, width, height, secs)

	counts := make(map[pos]int)

	for _, b := range bots {
		counts[pos{b.X, b.Y}] += 1
	}

	QUADS := [][]pos{
		{{0, 0}, {width/2 - 1, height/2 - 1}},
		{{width/2 + 1, 0}, {width - 1, height/2 - 1}},
		{{0, height/2 + 1}, {width/2 - 1, height - 1}},
		{{width/2 + 1, height/2 + 1}, {width - 1, height - 1}},
	}

	score := 1

	for _, q := range QUADS {
		acc := 0
		for i := q[0].x; i <= q[1].x; i++ {
			for j := q[0].y; j <= q[1].y; j++ {
				acc += counts[pos{i, j}]
			}
		}
		score *= acc
	}

	return score
}

func hasLine(bots []Robot, w, h int) bool {
	counts := make(map[pos]bool)

	for _, b := range bots {
		counts[pos{b.X, b.Y}] = true
	}

	for y := 0; y < w; y++ {
		inLineCount := 0

		for x := 0; x < h; x++ {
			if counts[pos{x, y}] {
				inLineCount++
			} else {
				inLineCount = 0
			}

			if inLineCount == 10 {
				return true
			}
		}
	}

	return false
}

func printTable(bots []Robot, w, h int, secs int) {
	counts := make(map[pos]bool)

	for _, b := range bots {
		counts[pos{b.X, b.Y}] = true
	}

	table := make([][]byte, w)
	for y := range table {
		table[y] = make([]byte, h)
		for x := range table[y] {
			if counts[pos{x, y}] {
				table[y][x] = '#'
			} else {
				table[y][x] = '.'
			}
		}
	}

	for _, row := range table {
		for _, cell := range row {
			fmt.Printf("%c ", cell)
		}
		fmt.Printf("\n")
	}

	fmt.Printf("\nsecs: %d", secs+1)

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
