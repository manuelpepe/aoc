package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"
)

func main() {
	edges := parse(os.Args[1])
	fmt.Println(sol1(edges))
	fmt.Println(sol2(edges))

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

func sol2(edges map[string][]string) string {
	allverts := slices.Collect(maps.Keys(edges))
	items := bk(edges, []string{}, allverts, []string{})
	slices.Sort(items)
	return strings.Join(items, ",")
}

// N = neighbors
// R = all vertices contained in clique
// P = potential vertices in clique
// X = vertices not contained in clique
//
// Bron–Kerbosch algorithm
// https://www.dcs.gla.ac.uk/~pat/jchoco/clique/enumeration/tex/report.pdf
func bk(N map[string][]string, R []string, P []string, X []string) []string {

	if len(P) == 0 && len(X) == 0 {
		// fmt.Printf("found max clique: %+v\n", R)
		return R
	}

	max := []string{}

	for len(P) > 0 {
		v := P[0]

		r := bk(N, union(R, v), intersection(P, N[v]), intersection(X, N[v]))
		if r != nil && len(r) > len(max) {
			max = r
		}

		P = complement(P, []string{v})
		X = union(X, v)
	}

	return max
}

func intersection(A, B []string) []string {
	out := make([]string, 0)
	for _, s := range A {
		if slices.Contains(B, s) {
			out = append(out, s)
		}
	}
	return out
}

func complement(A, B []string) []string {
	out := make([]string, 0)
	for _, s := range A {
		if !slices.Contains(B, s) {
			out = append(out, s)
		}
	}
	return out
}

func union(A []string, b string) []string {
	cln := slices.Clone(A)
	if slices.Contains(A, b) {
		return cln
	}
	return append(cln, b)
}
