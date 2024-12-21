package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	s1 := parse("logs")
	s2 := parse("log1")

	compare(s1, s2)
}

func parse(p string) []Pos {
	file, err := os.Open(p)
	if err != nil {
		panic(err)
	}

	out := make([]Pos, 0)

	s := bufio.NewScanner(file)
	for s.Scan() {
		line := s.Text()
		parts := strings.Split(line, ",")

		v1, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			panic(err)
		}

		v2, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			panic(err)
		}

		out = append(out, Pos{int(v1), int(v2)})
	}

	return out

}

type Pos struct {
	X, Y int
}

func compare(s1 []Pos, s2 []Pos) {

	m := make(map[Pos]bool)

	for _, v := range s1 {
		m[v] = true
	}

	for _, v := range s2 {
		if !m[v] {
			fmt.Printf("Pos{%d, %d}\n", v.X, v.Y)
		}
	}
}
