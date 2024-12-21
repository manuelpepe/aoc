package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"sync"
	"time"
)

var DIRS = [][]int{
	{0, -1},
	{1, 0},
	{0, 1},
	{-1, 0},
}

type Pos struct {
	X, Y int
}

func main() {
	fh, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	r := bufio.NewScanner(fh)

	table := make([][]byte, 0)
	playerPos := Pos{0, 0}
	playerDirIx := 0

	rix := 0
	for r.Scan() {
		line := r.Text()

		cells := []byte(line)
		for cix, c := range cells {
			if c == '^' {
				playerPos = Pos{cix, rix}
				cells[cix] = '.'
			}
		}

		table = append(table, cells)
		rix++
	}

	fmt.Printf("Guard visited %d diferent positions\n", sol1(clone(table), playerPos, playerDirIx))
	fmt.Printf("Obstructable positions to make loops: %d\n", sol2(table, playerPos, playerDirIx))

}

func sol1(table [][]byte, playerPos Pos, playerDirIx int) int {
	acc := 0

	for {
		playerPos.X += DIRS[playerDirIx][0]
		playerPos.Y += DIRS[playerDirIx][1]

		if !inBounds(table, playerPos) {
			break
		}

		if table[playerPos.Y][playerPos.X] == '#' {
			playerPos.X -= DIRS[playerDirIx][0]
			playerPos.Y -= DIRS[playerDirIx][1]
			playerDirIx = (playerDirIx + 1) % len(DIRS)
			continue
		}

		if table[playerPos.Y][playerPos.X] == '.' {
			table[playerPos.Y][playerPos.X] = 'X'
			acc++
		}

		// printTable(table, playerPos, playerDirIx, acc)
	}

	return acc
}

// I **WILDLY** overoptimized this before realizing I had a bug with cerain loops like:
//
//	.#....#......#......#.
//	#...................^#
//	..........#..#......##
//
// I would never reach the starting square with the same direction twice.
// Fixed adding a simple counter and if I go many times through the same position
// it's probably a loop ¯\_(ツ)_/¯
//
// Then I also missed that the loop must start from... the *beginning!*
// I was checking the loop starting from the middle of the path, at the move right before adding the obstacle,
// failing to realize that adding it might make the player never get to the same square to begin with.
//
// I'll leave the optimizations just because they are already done and part of the journey.
func sol2(table [][]byte, playerPos Pos, playerDirIx int) int {
	type traveledPosition struct {
		pos Pos
		dir int
	}

	traveledPositions := make(chan traveledPosition)

	go func(t [][]byte, playerPos Pos, dirIx int) {
		visited := make(map[traveledPosition]bool)

		for {
			playerPos.X += DIRS[dirIx][0]
			playerPos.Y += DIRS[dirIx][1]

			if !inBounds(t, playerPos) {
				break
			}

			if t[playerPos.Y][playerPos.X] == '#' {
				playerPos.X -= DIRS[dirIx][0]
				playerPos.Y -= DIRS[dirIx][1]
				dirIx = (dirIx + 1) % len(DIRS)

				tpos := traveledPosition{playerPos, dirIx}
				if !visited[tpos] {
					visited[tpos] = true
					traveledPositions <- tpos
				}
				continue
			}

			tpos := traveledPosition{playerPos, dirIx}
			if !visited[tpos] {
				visited[tpos] = true
				traveledPositions <- tpos
			}

			if t[playerPos.Y][playerPos.X] == '.' {
				t[playerPos.Y][playerPos.X] = 'X'
			}

			// printTable(t, playerPos, dirIx, acc)
		}

		close(traveledPositions)
	}(clone(table), playerPos, playerDirIx)

	routines := 12
	wg := sync.WaitGroup{}
	wg.Add(routines)

	validObstacles := make(chan Pos)

	for i := 0; i < routines; i++ {
		go func(ix int, t [][]byte) {
			// fmt.Printf("routine %d: starting\n", ix)

			for {
				p, open := <-traveledPositions
				if !open {
					break
				}

				newObstaclePos := Pos{
					X: p.pos.X + DIRS[p.dir][0],
					Y: p.pos.Y + DIRS[p.dir][1],
				}

				if wouldLoop(t, playerPos, playerDirIx, p.pos, p.dir, newObstaclePos) {
					validObstacles <- newObstaclePos
				}
			}

			wg.Done()
			// fmt.Printf("routine %d: done\n", ix)
		}(i, clone(table))
	}

	go func() {
		wg.Wait()
		close(validObstacles)
	}()

	seen := make(map[Pos]bool)
	acc := 0

	for {
		obs, open := <-validObstacles
		if !open {
			break
		}

		// grid := clone(table)
		// grid[obs.Y][obs.X] = '#'
		// showPath(grid, obs.playerPos, obs.playerDir)

		if !seen[obs] {
			seen[obs] = true
			acc++
			// fmt.Printf("loop found at Pos{%d, %d}\n", obs.X, obs.Y)
		}
	}

	return acc
}

func wouldLoop(table [][]byte, start Pos, startDir int, targetPos Pos, targetDir int, obstacle Pos) bool {
	if cached, ok := CACHE.Get(obstacle); ok {
		return cached
	}

	found := false

	// count the number of times a cell has been visited in case
	// a loop forms without ever reaching the starting square in the same direction
	counters := make(map[Pos]int)

	targetHit := 0

	if !inBounds(table, obstacle) {
		return false
	}

	tmp := table[obstacle.Y][obstacle.X]
	table[obstacle.Y][obstacle.X] = '#'

	playerPos := start
	dirIx := startDir

	for {
		playerPos.X += DIRS[dirIx][0]
		playerPos.Y += DIRS[dirIx][1]

		if !inBounds(table, playerPos) {
			break
		}

		if table[playerPos.Y][playerPos.X] == '#' {
			playerPos.X -= DIRS[dirIx][0]
			playerPos.Y -= DIRS[dirIx][1]
			dirIx = (dirIx + 1) % len(DIRS)
			continue
		}

		// early exit if target position and direction before adding obstacle
		// gets hit twice.
		if playerPos == targetPos && targetDir == dirIx {
			targetHit++
			if targetHit > 1 {
				found = true
				break
			}
		}

		table[playerPos.Y][playerPos.X] = 'X'
		counters[playerPos] += 1

		if counters[playerPos] > 10 {
			found = true // likely a loop...
			break
		}

		// printTable(table, playerPos, dirIx, -1)
	}

	table[obstacle.Y][obstacle.X] = tmp

	return CACHE.Set(obstacle, found)
}

func clone(t [][]byte) [][]byte {
	o := make([][]byte, len(t))
	for i := range t {
		o[i] = slices.Clone(t[i])
	}
	return o
}

func inBounds(matrix [][]byte, pos Pos) bool {
	x, y := pos.X, pos.Y

	if y < 0 || y >= len(matrix) {
		return false
	}

	if x < 0 || x >= len(matrix[0]) {
		return false
	}

	return true
}

// debugging...

func printTable(table [][]byte, playerPos Pos, playerDirIx int, acc int) {
	playerDirs := []byte{
		'^',
		'>',
		'v',
		'<',
	}

	tmp := table[playerPos.Y][playerPos.X]
	table[playerPos.Y][playerPos.X] = playerDirs[playerDirIx]

	for _, row := range table {
		for _, cell := range row {
			fmt.Printf("%c ", cell)
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\nacc: %d", acc)

	table[playerPos.Y][playerPos.X] = tmp

	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')
}

func printTable2(table [][]byte, playerPos Pos, playerDirIx int) {

	playerDirs := []byte{
		'^',
		'>',
		'v',
		'<',
	}

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	tmp := table[playerPos.Y][playerPos.X]
	table[playerPos.Y][playerPos.X] = playerDirs[playerDirIx]

	for _, row := range table {
		for _, cell := range row {
			fmt.Printf("%c ", cell)
		}
		fmt.Printf("\n")
	}

	table[playerPos.Y][playerPos.X] = tmp
}

var ticker = time.NewTicker(10 * time.Millisecond)

func showPath(table [][]byte, from Pos, dirIx int) {

	pos := from
	dir := dirIx

	exit := make(chan struct{})

	go func() {
		reader := bufio.NewReader(os.Stdin)
		reader.ReadString('\n')
		exit <- struct{}{}
	}()

	for {
		pos.X += DIRS[dir][0]
		pos.Y += DIRS[dir][1]

		if !inBounds(table, pos) {
			break
		}

		if table[pos.Y][pos.X] == '#' {
			pos.X -= DIRS[dir][0]
			pos.Y -= DIRS[dir][1]
			dir = (dir + 1) % len(DIRS)
			continue
		}

		if table[pos.Y][pos.X] == '.' {
			table[pos.Y][pos.X] = 'X'
		}

		printTable2(table, pos, dir)

		select {
		case <-ticker.C:
			continue
		case <-exit:
			return
		}
	}

}

// memoization

type cache struct {
	c map[Pos]bool
	m sync.RWMutex
}

func (c *cache) Get(k Pos) (bool, bool) {
	c.m.RLock()
	defer c.m.RUnlock()
	v, ok := c.c[k]
	return v, ok
}

func (c *cache) Set(k Pos, v bool) bool {
	c.m.Lock()
	defer c.m.Unlock()
	c.c[k] = v
	return v
}

var CACHE = cache{
	c: make(map[Pos]bool),
	m: sync.RWMutex{},
}
