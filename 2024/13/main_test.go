package main

import (
	"os"
	"testing"
)

func Test1(t *testing.T) {
	file, err := os.Open("input-0.txt")
	if err != nil {
		t.Fail()
	}

	machines := parse(file)

	if sol2(machines) != 480 {
		t.Fail()
	}
}
