package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

func main() {
	edges := parse(os.Args[1])
	fmt.Println(sol1(edges))
}

func parse(file string) map[string][]string {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}

	edges := make(map[string][]string)

	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			panic("expected 2 items in each line")
		}

		edges[parts[0]] = append(edges[parts[0]], parts[1])
		edges[parts[1]] = append(edges[parts[1]], parts[0])
	}

	return edges
}

func sol1(edges map[string][]string) int {
	return len(findCycles(edges))
}

type Cycle [3]string

func findCycles(edges map[string][]string) []Cycle {

	// find cycles for all edges
	acc := make([]Cycle, 0)
	for node, _ := range edges {
		// fmt.Printf("-> finding for %s\n", node)
		acc = append(acc, findCyclesRec(node, edges, make([]string, 0))...)
	}

	// dedupe cycles
	ordered := make([]Cycle, 0)
	for _, c := range acc {
		cs := c[:]
		slices.Sort(cs)

		if slices.Contains(ordered, Cycle(cs)) {
			continue
		}

		if !slices.ContainsFunc(cs, func(e string) bool { return e[0] == 't' }) {
			continue
		}

		ordered = append(ordered, Cycle(cs))
		// fmt.Printf("%+v\n", c)
	}

	return ordered
}

func findCyclesRec(node string, edges map[string][]string, currentPath []string) []Cycle {
	// fmt.Printf("  -> %+v - %s\n", currentPath, node)
	neighs := edges[node]

	if len(currentPath) == 2 {
		if slices.Contains(neighs, currentPath[0]) {
			// fmt.Printf("    -> found loop: %+v\n", currentPath)
			return []Cycle{Cycle(append(currentPath, node))}
		}
		return nil
	}

	out := make([]Cycle, 0)

	for _, neigh := range neighs {
		if len(currentPath) > 0 && neigh == currentPath[len(currentPath)-1] {
			continue // skip last node visited
		}
		out = append(out, findCyclesRec(neigh, edges, append(currentPath, node))...)
	}

	return out
}
