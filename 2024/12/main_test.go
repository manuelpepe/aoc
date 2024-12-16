package main

import (
	"strings"
	"testing"
)

func Test1(t *testing.T) {
	data := "AAAA\nBBCD\nBBCC\nEEEC"

	table := make([][]*Plot, 0)
	for _, row := range strings.Split(data, "\n") {
		line := []byte(row)

		row := make([]*Plot, len(line))
		for ix, p := range line {
			row[ix] = &Plot{Plant: p, Visited: false}
		}

		table = append(table, row)
	}

	if sol1(table) != 140 {
		t.Fail()
	}
}
