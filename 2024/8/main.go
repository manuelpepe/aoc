package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

type Pos struct{ x, y int }

func (p Pos) Sub(p2 Pos) Pos {
	return Pos{
		p.x - p2.x,
		p.y - p2.y,
	}
}

func (p Pos) Add(p2 Pos) Pos {
	return Pos{
		p.x + p2.x,
		p.y + p2.y,
	}
}

func main() {
	fh, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	positions := make(map[byte][]Pos, 0)
	table := make([][]byte, 0)

	r := bufio.NewScanner(fh)
	y := 0
	for r.Scan() {
		line := []byte(r.Text())
		for x, v := range line {
			if v == '.' {
				continue
			}

			positions[v] = append(positions[v], Pos{x, y})
		}

		table = append(table, line)
		y++
	}

	table2 := make([][]byte, len(table))
	for ix, line := range table {
		table2[ix] = slices.Clone(line)
	}
	fmt.Printf("antinode count is: %d\n", sol1(table, positions))
	fmt.Printf("resonant harmonious antinode count is: %d\n", sol2(table2, positions))

}

func sol1(table [][]byte, positions map[byte][]Pos) int {
	acc := 0
	for _, antennas := range positions {
		for pos1, ant1 := range antennas {
			for pos2, ant2 := range antennas {
				if pos1 == pos2 {
					continue
				}

				distance := ant1.Sub(ant2)
				antinode := ant1.Add(distance)

				if inBounds(table, antinode) && table[antinode.y][antinode.x] != '#' {
					table[antinode.y][antinode.x] = '#'
					acc++
				}
			}
		}
	}
	return acc
}

func inBounds(table [][]byte, pos Pos) bool {
	if pos.y < 0 || pos.y >= len(table) {
		return false
	}

	if pos.x < 0 || pos.x >= len(table[0]) {
		return false
	}

	return true
}

func sol2(table [][]byte, positions map[byte][]Pos) int {
	acc := 0
	for _, antennas := range positions {
		for pos1, ant1 := range antennas {
			for pos2, ant2 := range antennas {
				if pos1 == pos2 {
					continue
				}

				distance := ant1.Sub(ant2)
				antinode := ant1

				for inBounds(table, antinode) {
					if table[antinode.y][antinode.x] != '#' {
						table[antinode.y][antinode.x] = '#'
						acc++
					}
					antinode = antinode.Add(distance)
				}
			}
		}
	}
	return acc
}
