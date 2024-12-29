package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	nums := parse(os.Args[1])
	fmt.Printf("solution 1: %d\n", sol1(nums))
}

func sol1(nums []int) int {
	acc := 0

	for _, n := range nums {
		last := (n / 3) - 2
		subacc := last
		for {
			next := (last / 3) - 2
			if next < 0 {
				break
			}
			subacc += next
			last = next
		}
		acc += subacc
	}

	return acc
}

func parse(fn string) []int {
	fh, err := os.Open(fn)
	if err != nil {
		panic(err)
	}

	out := make([]int, 0)

	s := bufio.NewScanner(fh)
	for s.Scan() {
		line := s.Text()
		n, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			panic(err)
		}

		out = append(out, int(n))
	}

	return out
}
