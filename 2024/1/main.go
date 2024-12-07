package main

import (
	"bufio"
	"fmt"
	"math"
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

	var list1 []int
	var list2 []int

	for r.Scan() {
		line := r.Text()

		parts := strings.Split(line, "   ")

		v1, err := strconv.ParseInt(parts[0], 10, 32)
		if err != nil {
			panic(err)
		}

		v2, err := strconv.ParseInt(parts[1], 10, 32)
		if err != nil {
			panic(err)
		}

		list1 = append(list1, int(v1))
		list2 = append(list2, int(v2))
	}

	slices.Sort(list1)
	slices.Sort(list2)

	if len(list1) != len(list2) {
		panic("arrs of different length")
	}

	sol1(list1, list2)
	sol2(list1, list2)
}

func sol1(list1, list2 []int) {
	acc := 0
	for i := range list2 {
		d := int(math.Abs(float64(list1[i] - list2[i])))
		acc += d
		fmt.Printf("distance from %d to %d is %d \n", list1[i], list2[i], d)
	}

	fmt.Printf("Total distance is %d\n", acc)
}

func sol2(list1, list2 []int) {
	repeats := make(map[int]int)
	for _, v := range list2 {
		repeats[v]++
	}

	acc := 0
	for _, v := range list1 {
		acc += v * repeats[v]
	}

	fmt.Printf("Similarity is %d\n", acc)
}
