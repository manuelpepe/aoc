package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type XY struct {
	X int
	Y int
}

func (xy XY) Plus(xy2 XY) XY {
	return XY{xy.X + xy2.X, xy.Y + xy2.Y}
}

type Button XY

func (b Button) Times(n int) XY {
	return XY{n * b.X, n * b.Y}
}

func (b Button) For(xy XY) int {
	return min(xy.X/b.X, xy.Y/b.Y)
}

type Machine struct {
	A     Button
	B     Button
	Prize XY
}

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	machines := parse(file)

	fmt.Printf("solution 1 is: %d\n", sol1(machines))
}

func parse(file *os.File) []Machine {
	machines := make([]Machine, 0)
	r := bufio.NewReader(file)

	for {
		machine := Machine{}

		skipTo(r, '+')
		machine.A.X = readIntTo(r, ',')
		skipTo(r, '+')
		machine.A.Y = readIntTo(r, '\n')

		skipTo(r, '+')
		machine.B.X = readIntTo(r, ',')
		skipTo(r, '+')
		machine.B.Y = readIntTo(r, '\n')

		skipTo(r, '=')
		machine.Prize.X = readIntTo(r, ',')
		skipTo(r, '=')
		machine.Prize.Y = readIntTo(r, '\n')

		// fmt.Printf("%+v\n", machine)
		machines = append(machines, machine)

		if err := skipTo(r, '\n'); err != nil {
			break
		}
	}

	return machines
}

func skipTo(r *bufio.Reader, delim byte) error {
	if _, err := r.ReadBytes(delim); err != nil {
		if errors.Is(err, io.EOF) {
			return err
		}
		panic(err)
	}
	return nil
}

func readIntTo(r *bufio.Reader, delim byte) int {
	val, err := r.ReadString(delim)
	if err != nil {
		panic(err)
	}

	parsed, err := strconv.ParseInt(val[:len(val)-1], 10, 64)
	if err != nil {
		panic(err)
	}

	return int(parsed)
}

func sol1(machines []Machine) int {
	acc := 0

	for _, m := range machines {
		tokens := solve(m, 0, 0)
		if tokens == 99999999 {
			continue
		}
		acc += tokens
	}

	return acc
}

type call struct {
	m    Machine
	a, b int
}

var cache = make(map[call]int)

func solve(m Machine, a, b int) int {
	if val, ok := cache[call{m, a, b}]; ok {
		return val
	}

	axy := m.A.Times(a)
	bxy := m.B.Times(b)

	pos := axy.Plus(bxy)

	if pos.X == m.Prize.X && pos.Y == m.Prize.Y {
		return a*3 + b
	}

	if pos.X > m.Prize.X || pos.Y > m.Prize.Y {
		return 99999999
	}

	res := min(solve(m, a+1, b), solve(m, a, b+1))
	cache[call{m, a, b}] = res

	return res
}
