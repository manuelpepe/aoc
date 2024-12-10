package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Eq struct {
	Exp        int64
	Components []int64
}

func main() {
	fh, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	r := bufio.NewScanner(fh)
	eqs := make([]Eq, 0)

	for r.Scan() {
		line := r.Text()
		sepIx := strings.Index(line, ":")

		exp, err := strconv.ParseInt(line[:sepIx], 10, 64)
		if err != nil {
			panic(err)
		}

		compsStr := strings.Split(line[sepIx+2:], " ")
		comps := make([]int64, len(compsStr))
		for ix, v := range compsStr {
			comps[ix], err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				panic(err)
			}
		}

		eqs = append(eqs, Eq{
			Exp:        exp,
			Components: comps,
		})
	}

	fmt.Printf("Calibration is: %d\n", sol1(eqs))
	fmt.Printf("Calibration with join is: %d\n", sol2(eqs))
}

func sol1(eqs []Eq) int {
	var OPS = []byte{
		'*',
		'+',
	}
	return process(eqs, OPS)
}

func sol2(eqs []Eq) int {
	var OPS = []byte{
		'*',
		'+',
		'|', // operator is actually defined as || but this makes it easier
	}

	return process(eqs, OPS)
}

func process(eqs []Eq, ops []byte) int {
	total := 0

EQS:
	for _, eq := range eqs {
		perms := permutations(ops, len(eq.Components)-1)

		for _, perm := range perms {
			acc := eq.Components[0]

			for ix, op := range perm {
				acc = doOp(acc, eq.Components[ix+1], op)
			}

			if acc == eq.Exp {
				total += int(eq.Exp)
				continue EQS
			}
		}
	}

	return total
}

func permutations(chars []byte, reps int) [][]byte {
	return permsRec([][]byte{}, chars, reps-1)
}

func permsRec(acc [][]byte, options []byte, depth int) [][]byte {
	if depth < 0 {
		panic("invalid depth")
	}
	if depth == 0 {
		for _, v := range options {
			acc = append(acc, []byte{v})
		}
		return acc
	}

	base := permsRec(acc, options, depth-1)

	newacc := make([][]byte, 0)
	for _, row := range base {
		for _, opt := range options {
			newRow := append([]byte{}, row...)
			newRow = append(newRow, opt)
			newacc = append(newacc, newRow)
		}
	}

	return newacc
}

func doOp(acc int64, comp int64, op byte) int64 {
	switch op {
	case '*':
		acc *= comp
	case '+':
		acc += comp
	case '|':
		acc = join(acc, comp) // join components
	}

	return acc
}

func join(a, b int64) int64 {
	lenB := math.Floor(math.Log10(float64(b))) + 1
	return a*int64(math.Pow(10, lenB)) + b
}
