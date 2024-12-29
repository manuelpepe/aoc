package parser

import (
	"os"
	"strconv"
	"strings"
)

func Parse(fn string) []int {
	rawdata, err := os.ReadFile(fn)
	if err != nil {
		panic(err)
	}
	data, _ := strings.CutSuffix(string(rawdata), "\n")
	items := strings.Split(data, ",")

	out := make([]int, 0)

	for _, i := range items {
		n, err := strconv.ParseInt(i, 10, 64)
		if err != nil {
			panic(err)
		}
		out = append(out, int(n))
	}

	return out
}
