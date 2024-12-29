package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	min, max := parse(os.Args[1])

	fmt.Printf("solution 1: %d\n", sol(min, max, fulfills))
	fmt.Printf("solution 2: %d\n", sol(min, max, fulfills2))

}

func sol(min, max int, fn func(int) bool) int {
	acc := 0
	for pw := min; pw < max; pw++ {
		if fn(pw) {
			acc++
		}
	}
	return acc
}

func fulfills2(pw int) bool {
	digits := []int{
		int(pw % 10),
		int(pw % 100 / 10),
		int(pw % 1000 / 100),
		int(pw % 10000 / 1000),
		int(pw % 100000 / 10000),
		int(pw % 1000000 / 100000),
	}

	groups := make(map[int]int)

	alwaysIncreases := true
	for i := 0; i < len(digits)-1; i++ {
		if digits[i] == digits[i+1] {
			groups[digits[i]]++
		}

		if digits[i] < digits[i+1] {
			alwaysIncreases = false
			break
		}
	}

	doubleFound := false
	for _, c := range groups {
		if c == 1 {
			doubleFound = true
			break
		}
	}

	return doubleFound && alwaysIncreases
}

func fulfills(pw int) bool {
	digits := []int{
		int(pw % 10),
		int(pw % 100 / 10),
		int(pw % 1000 / 100),
		int(pw % 10000 / 1000),
		int(pw % 100000 / 10000),
		int(pw % 1000000 / 100000),
	}

	doubleFound := false
	alwaysIncreases := true
	for i := 0; i < len(digits)-1; i++ {
		if digits[i] == digits[i+1] {
			doubleFound = true
		}

		if digits[i] < digits[i+1] {
			alwaysIncreases = false
			break
		}
	}
	return doubleFound && alwaysIncreases
}

func parse(file string) (int, int) {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	rang := string(data)

	parts := strings.Split(rang, "-")
	min, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		panic(err)
	}
	max, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		panic(err)
	}
	return int(min), int(max)
}
