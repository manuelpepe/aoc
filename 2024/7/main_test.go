package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPerms(t *testing.T) {
	out := permutations([]byte{'0', '1'}, 3)

	if len(out) == 0 {
		t.Fatal("len is 0")
	}

	for _, r := range out {
		for _, c := range r {
			fmt.Printf("%c ", c)
		}
		fmt.Printf("\n")
	}

	if len(out) != 8 {
		t.Fail()
	}
}

func TestJoin(t *testing.T) {
	if join(12, 345) != 12345 {
		t.Fail()
	}
}

func TestWrong_Sol2(t *testing.T) {
	// res := sol2([]Eq{
	// 	{
	// 		Exp:        7290,
	// 		Components: []int64{6, 8, 6, 15},
	// 	},
	// })

	// if res != 7290 {
	// 	t.Fatalf("expected %d to equal 7290", res)
	// }

	var acc int64

	acc = 6

	acc = doOp(acc, 8, '*')
	assert.Equal(t, int(acc), 6*8)

	acc = doOp(acc, 6, '|')
	assert.Equal(t, int(acc), 486)

	acc = doOp(acc, 15, '*')
	assert.Equal(t, int(acc), 486*15)

}
