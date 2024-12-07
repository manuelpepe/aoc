package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var debug = false

func log(msg string, args ...any) {
	if debug {
		fmt.Printf(msg, args...)
	}
}

func main() {

	fh, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	r := bufio.NewScanner(fh)

	var total int

	lines := make([][]int, 0)

	for r.Scan() {
		total++
		line := r.Text()

		numberStrs := strings.Split(line, " ")
		var numbers []int
		for _, v := range numberStrs {
			n, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				panic(err)
			}
			numbers = append(numbers, int(n))
		}

		lines = append(lines, numbers)
	}

	sol1(lines)
	sol2_bruteforce(lines)
	sol2(lines)

}

func sol1(lines [][]int) {
	var acc int

	for _, numbers := range lines {
		if isSafe(numbers) {
			acc++
		}
	}

	fmt.Printf("Total of safe reports: %d of %d\n", acc, len(lines))
}

func sol2_bruteforce(lines [][]int) {

	alternatives := func(line []int) [][]int {
		alts := make([][]int, len(line))
		for i := range line {
			alts[i] = append(line[:i:i], line[i+1:]...)
		}
		return alts
	}

	hasSafeAlt := func(line []int) bool {
		for _, alt := range alternatives(line) {
			if ok := isSafe(alt); ok {
				return true
			}
		}
		return false
	}

	var acc int

	for _, line := range lines {
		if isSafe(line) || hasSafeAlt(line) {
			acc++
		}
	}

	fmt.Printf("Total of safe reports (damp+brute): %d of %d\n", acc, len(lines))
}

func isSafe(numbers []int) bool {
	n := numbers[0]
	n2 := numbers[1]

	var cmp func(int, int) bool
	if n > n2 {
		cmp = gt
	} else {
		cmp = lt
	}

	i := 2
	for {
		if !fulfils(cmp, n, n2) {
			return false
		}

		if i == len(numbers) {
			return true
		}

		n = n2
		n2 = numbers[i]
		i++
	}
}

func fulfils(cmp func(int, int) bool, n, n2 int) bool {
	diff := math.Abs(float64(n - n2))
	return cmp(n, n2) && diff > 0 && diff < 4
}

func gt(a, b int) bool {
	return a > b
}

func lt(a, b int) bool {
	return a < b
}

// FIXME: not working when first position should  be skipped
// i.e.:
//
//	2 1 3 4 5
func sol2(lines [][]int) {
	var acc int

	for _, numbers := range lines {
		n := numbers[0]
		n2 := numbers[1]

		var cmp func(int, int) bool
		if n > n2 {
			cmp = gt
		} else {
			cmp = lt
		}

		firstFail := false
		i := 2
		for {
			log("%d %d ", n, n2)

			if !fulfils(cmp, n, n2) {
				log("fail\n")
				if firstFail {
					log("%v - BREAK\n", numbers)
					break
				}
				firstFail = true
				n2 = n
			} else {
				log("ok\n")
			}

			if i == len(numbers) {
				log("%v - OK\n", numbers)
				acc++
				break
			}

			n = n2
			n2 = numbers[i]

			log("%d\n", i)
			if firstFail && i == 3 {
				// recalculate cmp if first level needs to be removed
				if n > n2 {
					cmp = gt
					log("%d %d are now gt\n", n, n2)
				} else {
					cmp = lt
					log("%d %d are now lt\n", n, n2)
				}
			}

			i++

		}
	}

	fmt.Printf("Total of safe reports (damp) [WRONG]: %d of %d\n", acc, len(lines))
}
