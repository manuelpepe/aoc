package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	lines := parse(os.Args[1])

	fmt.Printf("solution 1: %d\n", sol1(lines))
	fmt.Printf("solution 2: %d\n", sol2(lines))
}

func sol2(lines [][]Segment) int {
	aLines := lines[0]
	bLines := lines[1]

	aLinesH, aLinesV := splitHV(aLines)
	bLinesH, bLinesV := splitHV(bLines)

	insersections := make([]Intersection, 0)
	insersections = append(insersections, getIntersections(aLinesH, bLinesV)...)
	insersections = append(insersections, getIntersections(aLinesV, bLinesH)...)

	var minDistance int

	for _, i := range insersections {
		distA := distanceTo(aLines, i.IndexA)
		halfSeg := Segment{From: aLines[i.IndexA].From, To: i.Pos}
		distA += halfSeg.Distance()

		distB := distanceTo(bLines, i.IndexB)
		halfSeg = Segment{From: bLines[i.IndexB].From, To: i.Pos}
		distB += halfSeg.Distance()

		totalDist := distA + distB
		if minDistance == 0 || totalDist < minDistance {
			minDistance = totalDist
		}

	}

	return minDistance
}

func distanceTo(lines []Segment, ix int) int {
	acc := 0

	for lix := range lines {
		if lix == ix {
			break
		}

		acc += lines[lix].Distance()
	}

	return acc
}

func sol1(lines [][]Segment) int {
	aLines := lines[0]
	bLines := lines[1]

	aLinesH, aLinesV := splitHV(aLines)
	bLinesH, bLinesV := splitHV(bLines)

	insersections := make([]Intersection, 0)
	insersections = append(insersections, getIntersections(aLinesH, bLinesV)...)
	insersections = append(insersections, getIntersections(aLinesV, bLinesH)...)

	// fmt.Printf("found %d intersections\n", len(insersections))
	// fmt.Printf("intersections: %+v\n", insersections)

	min := 0

	for _, i := range insersections {
		d := manhattanDistance(Pos{0, 0}, i.Pos)
		if min == 0 || d < min {
			min = d
		}
	}

	return min
}

func getIntersections(lines, lines2 []Segment) []Intersection {
	insersections := make([]Intersection, 0)

	for _, line := range lines {
		for _, line2 := range lines2 {
			inters, found := intersect(line, line2)
			if !found {
				continue
			}

			insersections = append(insersections, inters)
		}
	}

	return insersections
}

// https://www.youtube.com/watch?v=bvlIYX9cgls
func intersect(line, line2 Segment) (Intersection, bool) {
	a, b, c, d := line.From, line.To, line2.From, line2.To

	denom := float64((d.X-c.X)*(b.Y-a.Y) - (d.Y-c.Y)*(b.X-a.X))

	if denom == 0 {
		panic("lines are parallel or coincident (unexpected)")
	}

	alpha := float64((d.X-c.X)*(c.Y-a.Y)-(d.Y-c.Y)*(c.X-a.X)) / denom
	beta := float64((b.X-a.X)*(c.Y-a.Y)-(b.Y-a.Y)*(c.X-a.X)) / denom

	if alpha < 0 || alpha > 1 {
		return Intersection{}, false // intersects beyond given segments
	}

	if beta < 0 || beta > 1 {
		return Intersection{}, false // intersects beyond given segments
	}

	intersection := Pos{
		X: a.X + int(alpha*float64(b.X-a.X)),
		Y: a.Y + int(alpha*float64(b.Y-a.Y)),
	}

	// fmt.Printf("INTERSECTS: AB=%+v  CD=%+v  P=%+v\n", line, line2, intersection)

	return Intersection{
		Pos:    intersection,
		IndexA: line.Index,
		IndexB: line2.Index,
	}, true

}

func splitHV(lines []Segment) ([]Segment, []Segment) {
	horizontal, vertical := make([]Segment, 0), make([]Segment, 0)

	for _, line := range lines {
		if line.From.X != line.To.X {
			horizontal = append(horizontal, line)
			continue
		}

		if line.From.Y != line.To.Y {
			vertical = append(vertical, line)
			continue
		}

		panic("found a point")
	}

	return horizontal, vertical
}

func parse(fn string) [][]Segment {
	fh, err := os.Open(fn)
	if err != nil {
		panic(err)
	}

	out := make([][]Segment, 0)

	s := bufio.NewScanner(fh)
	for s.Scan() {
		line := make([]Segment, 0)

		last := Pos{0, 0}

		for ix, move := range strings.Split(s.Text(), ",") {
			next := makeMove(last, move)
			line = append(line, Segment{ix, last, next})
			last = next
		}

		out = append(out, line)
	}

	return out
}

func makeMove(from Pos, move string) Pos {
	delta, err := strconv.ParseInt(move[1:], 10, 64)
	if err != nil {
		panic(err)
	}

	newPos := from

	switch move[0] {
	case 'U':
		newPos.Y += int(delta)
	case 'D':
		newPos.Y -= int(delta)
	case 'R':
		newPos.X += int(delta)
	case 'L':
		newPos.X -= int(delta)
	default:
		panic(fmt.Sprintf("unexpected move: %s", move))
	}

	return newPos
}

type Pos struct {
	X, Y int
}

type Segment struct {
	Index int
	From  Pos
	To    Pos
}

func (s Segment) Distance() int {
	return manhattanDistance(s.From, s.To)
}

type Intersection struct {
	Pos            Pos
	IndexA, IndexB int
}

func manhattanDistance(a, b Pos) int {
	return abs(a.X-b.X) + abs(a.Y-b.Y)
}

func abs[T ~int](x T) T {
	if x < 0 {
		return -x
	}
	return x
}
