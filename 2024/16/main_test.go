package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {

	file, err := os.Open("test-0.txt")
	if err != nil {
		panic(err)
	}

	grid, start, end := parse(file)
	cost := pathCost(grid, start, end)

	assert.Equal(t, 2, cost)

}
