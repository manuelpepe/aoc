package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	locks, keys := parse(os.Args[1])
	fmt.Printf("parsed %d locks and %d keys\n", len(locks), len(keys))
	sol1(locks, keys)
}

func sol1(locks []Lock, keys []Key) {
	acc := 0
	for _, lock := range locks {
		for _, key := range keys {
			if fits(lock, key) {
				acc++
			}
		}
	}

	fmt.Printf("found %d unique lock/key pairs that fit\n", acc)
}

func fits(lock Lock, key Key) bool {
	for i := 0; i < 5; i++ {
		if lock[i]+key[i] > 5 {
			return false
		}
	}
	return true
}

type Lock [5]int
type Key [5]int

func parse(fn string) ([]Lock, []Key) {
	file, err := os.Open(fn)
	if err != nil {
		panic(err)
	}

	s := bufio.NewScanner(file)

	type state int

	var WaitingForFirstRow state = 0
	var ParsingLock state = 1
	var ParsingKey state = 2
	var ExpectingNL state = 3

	currentState := WaitingForFirstRow

	locks := make([]Lock, 0)
	keys := make([]Key, 0)

	last := make([]int, 5)
	parsedRows := 0

	lineNum := 0

	for s.Scan() {
		lineNum++
		line := []byte(s.Text())

		switch currentState {

		case WaitingForFirstRow:
			last = make([]int, 5)
			parsedRows = 0

			if line[0] == '#' {
				currentState = ParsingLock
			} else {
				currentState = ParsingKey
			}

		case ParsingLock:
			assert(len(line) == 5, "expected lock to have 5 columns at line %d, got %d", lineNum, len(line))
			assert(parsedRows < 6, "parsed rows out of bounds at line %d", lineNum)

			if parsedRows == 5 {
				locks = append(locks, Lock(last))
				currentState = ExpectingNL
				continue
			}

			for i := range line {
				if line[i] == '#' {
					last[i]++
				}
			}

			parsedRows++

		case ParsingKey:
			assert(len(line) == 5, "expected key to have 5 columns at line %d, got %d", lineNum, len(line))
			assert(parsedRows < 6, "parsed rows out of bounds at line %d", lineNum)

			if parsedRows == 5 {
				keys = append(keys, Key(last))
				currentState = ExpectingNL
				continue
			}

			for i := range line {
				if line[i] == '#' {
					last[i]++
				}
			}

			parsedRows++

		case ExpectingNL:
			assert(len(line) == 0, "expected newline at line %d, got '%s'", lineNum, line)
			currentState = WaitingForFirstRow

		default:
			panic("unexpected state")
		}
	}

	return locks, keys
}

func assert(cond bool, msg string, args ...any) {
	if !cond {
		formatted := fmt.Sprintf(msg, args...)
		panic(formatted)
	}
}
