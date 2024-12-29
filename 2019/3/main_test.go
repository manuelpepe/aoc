package main

import (
	"fmt"
	"testing"
)

func TestIntersection(t *testing.T) {

	tcs := []struct {
		l  Segment
		l2 Segment
	}{
		{
			l:  Segment{From: Pos{2, 3}, To: Pos{6, 3}},
			l2: Segment{From: Pos{4, 1}, To: Pos{4, 6}},
		},

		{
			l:  Segment{From: Pos{2, 3}, To: Pos{6, 3}},
			l2: Segment{From: Pos{4, 6}, To: Pos{4, 1}},
		},

		{
			l:  Segment{From: Pos{6, 3}, To: Pos{2, 3}},
			l2: Segment{From: Pos{4, 1}, To: Pos{4, 6}},
		},

		{
			l:  Segment{From: Pos{6, 3}, To: Pos{2, 3}},
			l2: Segment{From: Pos{4, 6}, To: Pos{4, 1}},
		},
	}

	for ix, tc := range tcs {
		t.Run(fmt.Sprintf("%d", ix), func(t *testing.T) {
			inters, ok := intersect(tc.l, tc.l2)
			if !ok {
				t.Fatalf("1: expected to find intersection")
			}

			if inters.Pos != (Pos{4, 3}) {
				t.Fatalf("1: expected intersection to be %+v, got %+v", Pos{4, 3}, inters.Pos)
			}
		})
	}

}
