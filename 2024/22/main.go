package main

import (
	"bufio"
	"fmt"
	"iter"
	"maps"
	"os"
	"slices"
	"strconv"
)

func main() {
	nums := parse(os.Args[1])

	r1 := sol1(nums)
	fmt.Println(r1)

	r2, seq := sol2(nums)
	fmt.Printf("%d -> %+v\n", r2, seq)

}

func parse(fn string) []int {
	file, err := os.Open(fn)
	if err != nil {
		panic("no file found")
	}

	out := make([]int, 0)

	s := bufio.NewScanner(file)
	for s.Scan() {
		v, err := strconv.ParseInt(s.Text(), 10, 64)
		if err != nil {
			panic("no file found")
		}

		out = append(out, int(v))
	}

	return out
}

func sol1(nums []int) int {
	acc := 0
	for _, n := range nums {
		acc += getRandomN(n, 2000)
	}
	return acc
}

type Seq [4]int

// basically bruteforcing, really slow but i had to log off for a while so i stuck with this
func sol2(nums []int) (int, Seq) {
	// find seq permutations that actually appear
	allSeqs := make(map[Seq]int)

	for _, n := range nums {
		seqs := findSeqs(n)
		for s := range seqs {
			allSeqs[s]++
		}
	}

	validSeqs := make([]Seq, 0)
	for s := range maps.Keys(allSeqs) {
		if allSeqs[s] > len(nums)/10 {
			validSeqs = append(validSeqs, s)
		}
	}

	var maxSeq Seq
	maxAcc := 0

	for ix, perm := range validSeqs {
		if ix%100 == 0 {
			fmt.Printf("check perm %d/%d : %+v... \n", ix, len(validSeqs), perm)
		}
		acc := 0
		for _, n := range nums {
			r := findFirstMatchingSeq(n, slices.Clone(perm[:]))
			acc += r
		}

		if acc > maxAcc {
			fmt.Printf("%+v - found acc of %d - updating previous max of %d\n", perm, acc, maxAcc)
			maxAcc = acc
			maxSeq = perm
		}
	}

	return maxAcc, maxSeq
}

func findSeqs(n int) iter.Seq[Seq] {
	out := make(map[Seq]struct{})

	runningPerm := make([]int, 0)

	last := n % 10

	// load first 4 perms
	n = getNextRandom(n)
	runningPerm = append(runningPerm, (n%10)-(last%10))
	last = n

	n = getNextRandom(n)
	runningPerm = append(runningPerm, (n%10)-(last%10))
	last = n

	n = getNextRandom(n)
	runningPerm = append(runningPerm, (n%10)-(last%10))
	last = n

	n = getNextRandom(n)
	runningPerm = append(runningPerm, (n%10)-(last%10))
	last = n

	out[Seq(runningPerm)] = struct{}{}

	for i := 4; i < 2000; i++ {
		n = getNextRandom(n)
		runningPerm = []int{runningPerm[1], runningPerm[2], runningPerm[3], (n % 10) - (last % 10)}
		last = n
		out[Seq(runningPerm)] = struct{}{}
	}

	return maps.Keys(out)
}

func findFirstMatchingSeq(n int, perm []int) int {
	runningPerm := make([]int, 0)

	last := n % 10

	// load first 4 perms
	n = getNextRandom(n)
	runningPerm = append(runningPerm, (n%10)-(last%10))
	last = n

	n = getNextRandom(n)
	runningPerm = append(runningPerm, (n%10)-(last%10))
	last = n

	n = getNextRandom(n)
	runningPerm = append(runningPerm, (n%10)-(last%10))
	last = n

	n = getNextRandom(n)
	runningPerm = append(runningPerm, (n%10)-(last%10))
	last = n

	// fmt.Printf("starting perm: %+v\n", runningPerm)

	if slices.Equal(perm, runningPerm) {
		return last % 10
	}

	for i := 4; i < 2000; i++ {
		n = getNextRandom(n)
		runningPerm = []int{runningPerm[1], runningPerm[2], runningPerm[3], (n % 10) - (last % 10)}
		if slices.Equal(perm, runningPerm) {
			return n % 10
		}
		last = n
	}

	return 0
}

func getRandomN(n int, iters int) int {
	for i := 0; i < iters; i++ {
		n = getNextRandom(n)
	}
	return n
}

func getNextRandom(n int) int {
	n = (n * 64) ^ n
	n = n % 16777216

	n = int(n/32) ^ n
	n = n % 16777216

	n = (n * 2048) ^ n
	n = n % 16777216

	return n
}
