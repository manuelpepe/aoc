package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

var DIRS = [][]int{
	{0, -1},
	{1, 0},
	{0, 1},
	{-1, 0},
}

func main() {
	fh, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	r := bufio.NewScanner(fh)

	table := make([][]byte, 0)
	playerPos := []int{0, 0}
	playerDirIx := 0

	rix := 0
	for r.Scan() {
		line := r.Text()

		cells := []byte(line)
		for cix, c := range cells {
			if c == '^' {
				playerPos = []int{cix, rix}
				cells[cix] = '.'
			}
		}

		table = append(table, cells)
		rix++
	}

	fmt.Printf("Guard visited %d diferent positions\n", sol1(table, playerPos, playerDirIx))

}

func sol1(table [][]byte, playerPos []int, playerDirIx int) int {
	acc := 0

	for {
		playerPos[0] += DIRS[playerDirIx][0]
		playerPos[1] += DIRS[playerDirIx][1]

		if !inBounds(table, playerPos) {
			break
		}

		if table[playerPos[1]][playerPos[0]] == '#' {
			playerPos[0] -= DIRS[playerDirIx][0]
			playerPos[1] -= DIRS[playerDirIx][1]
			playerDirIx = (playerDirIx + 1) % len(DIRS)
			continue
		}

		if table[playerPos[1]][playerPos[0]] == '.' {
			table[playerPos[1]][playerPos[0]] = 'X'
			acc++
		}

		// printTable(table, playerPos, playerDirIx, acc)
	}

	return acc
}

func inBounds(matrix [][]byte, pos []int) bool {
	x, y := pos[0], pos[1]

	if y < 0 || y >= len(matrix) {
		return false
	}

	if x < 0 || x >= len(matrix[0]) {
		return false
	}

	return true
}

func printTable(table [][]byte, playerPos []int, playerDirIx int, acc int) {
	playerDirs := []byte{
		'^',
		'>',
		'v',
		'<',
	}

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	tmp := table[playerPos[1]][playerPos[0]]
	table[playerPos[1]][playerPos[0]] = playerDirs[playerDirIx]

	for _, row := range table {
		for _, cell := range row {
			fmt.Printf("%c ", cell)
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\nacc: %d", acc)

	table[playerPos[1]][playerPos[0]] = tmp

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}
