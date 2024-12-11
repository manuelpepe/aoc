package main

import (
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	parts := strings.Split(string(data), " ")
	nums := make([]int64, len(parts))
	for ix, n := range parts {
		nums[ix], err = strconv.ParseInt(n, 10, 64)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("Total stones after %d blinks: %d\n", 25, sol1(slices.Clone(nums), 25))
	clear(cache)
	fmt.Printf("Total stones after %d blinks: %d\n", 75, sol1(nums, 75))
}

// If the stone is engraved with the number 0, it is replaced by a stone engraved with the number 1.
// If the stone is engraved with a number that has an even number of digits, it is replaced by two stones. The left half of the digits are engraved on the new left stone, and the right half of the digits are engraved on the new right stone. (The new numbers don't keep extra leading zeroes: 1000 would become stones 10 and 0.)
// If none of the other rules apply, the stone is replaced by a new stone; the old stone's number multiplied by 2024 is engraved on the new stone.
func sol1(nums []int64, blinks int) int {
	acc := 0
	for ix := 0; ix < len(nums); ix++ {
		acc += int(recSol(nums[ix], 0, blinks))
	}
	return acc
}

type call struct {
	n     int64
	depth int
}

var cache = make(map[call]int64)

func recSol(n int64, depth int, maxDepth int) int64 {
	if depth == maxDepth {
		return 1
	}

	if v, ok := cache[call{n, depth}]; ok {
		return v
	}

	if n == 0 {
		x := recSol(1, depth+1, maxDepth)
		cache[call{n, depth}] = x
		return x
	}

	numLen := int(math.Log10(float64(n))) + 1
	if numLen%2 == 0 {
		div := int64(math.Pow(10, math.Floor(float64(numLen)/2)))
		x := recSol(n/div, depth+1, maxDepth) + recSol(n%div, depth+1, maxDepth)
		cache[call{n, depth}] = x
		return x
	}

	x := recSol(n*2024, depth+1, maxDepth)
	return x
}
