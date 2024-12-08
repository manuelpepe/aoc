package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	fh, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	r := bufio.NewScanner(fh)

	updates := make([][]int64, 0)
	rules := make(map[int64][]int64)

	parsingRules := true
	for r.Scan() {
		line := r.Text()
		if line == "" {
			parsingRules = false
			continue
		}

		if parsingRules {
			parts := strings.Split(line, "|")

			x, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				panic(err)
			}

			y, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				panic(err)
			}

			rules[y] = append(rules[y], x)
		} else {
			parts := strings.Split(line, ",")

			update := make([]int64, len(parts))
			for i, v := range parts {
				nv, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					panic(err)
				}
				update[i] = nv
			}

			updates = append(updates, update)
		}

	}

	fmt.Printf("mid-page sum of ordered updates: %d\n", sol1(rules, updates))
	fmt.Printf("mid-page sum of unordered updates after fix: %d\n", sol2(rules, updates))

}

func sol1(rules map[int64][]int64, updates [][]int64) int64 {
	var acc int64

	for _, u := range updates {
		if !isOrdered(rules, u) {
			continue
		}

		acc += u[len(u)/2]
	}

	return acc
}

func sol2(rules map[int64][]int64, updates [][]int64) int64 {
	var acc int64
	for _, u := range updates {
		if isOrdered(rules, u) {
			continue
		}

		slices.SortFunc(u, func(a, b int64) int {
			mustBeBefore := rules[a]
			if slices.Contains(mustBeBefore, b) {
				return 1
			} else {
				return -1
			}
		})

		acc += u[len(u)/2]
	}

	return acc

}

func isOrdered(rules map[int64][]int64, update []int64) bool {
	for ix, n := range update {
		if ix == len(update) {
			continue
		}

		mustBeBefore := rules[n]

		for _, n2 := range update[ix+1:] {
			if slices.Contains(mustBeBefore, n2) {
				return false
			}
		}
	}

	return true
}
