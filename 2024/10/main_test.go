package main

import "testing"

func Test1(t *testing.T) {

	inp := []string{
		"89010123",
		"78121874",
		"87430965",
		"96549874",
		"45678903",
		"32019012",
		"01329801",
		"10456732",
	}

	table := make([][]byte, 0)
	for ix := range inp {
		line := []byte(inp[ix])
		for ix := range line {
			line[ix] -= '0'
		}
		table = append(table, line)
	}

	r1, r2 := sol1and2(table)
	if r1 != 36 {
		t.Fail()
	}

	if r2 != 81 {
		t.Fail()
	}
}
