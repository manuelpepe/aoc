package pq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertpop[T any](t *testing.T, q *PriorityQueue[T], exp T) {
	v, ok := q.Pop()
	assert.True(t, ok)
	assert.Equal(t, exp, v)
}

func TestBasicInt(t *testing.T) {
	q := New[int](MIN_HEAP)

	q.Push(3, 3)
	q.Push(1, 1)
	q.Push(2, 2)

	assertpop(t, q, 1)
	assertpop(t, q, 2)
	assertpop(t, q, 3)

	q = New[int](MAX_HEAP)

	q.Push(3, 3)
	q.Push(1, 1)
	q.Push(2, 2)

	assertpop(t, q, 3)
	assertpop(t, q, 2)
	assertpop(t, q, 1)
}

func TestBasicStr(t *testing.T) {
	q := New[string](MIN_HEAP)

	q.Push("a", 3)
	q.Push("b", 1)
	q.Push("c", 2)

	assertpop(t, q, "b")
	assertpop(t, q, "c")
	assertpop(t, q, "a")

	q = New[string](MAX_HEAP)

	q.Push("a", 3)
	q.Push("b", 1)
	q.Push("c", 2)

	assertpop(t, q, "a")
	assertpop(t, q, "c")
	assertpop(t, q, "b")
}

func TestPushPopPush(t *testing.T) {
	q := New[int](MIN_HEAP)

	q.Push(30, 30)
	q.Push(10, 10)
	q.Push(20, 20)

	assertpop(t, q, 10)

	q.Push(5, 5)
	assertpop(t, q, 5)

	q.Push(25, 25)
	assertpop(t, q, 20)
	assertpop(t, q, 25)
	assertpop(t, q, 30)

	for i := 0; i < 3; i++ {
		q.Push(30, 30)
		q.Push(10, 10)
		q.Push(20, 20)

		assertpop(t, q, 10)
		assertpop(t, q, 20)
		assertpop(t, q, 30)
	}

	q.Push(30, 30)
	q.Push(10, 10)
	q.Push(20, 20)

	assertpop(t, q, 10)

	q.Push(5, 5)
	q.Push(35, 35)
	q.Push(25, 25)

	assertpop(t, q, 5)
	assertpop(t, q, 20)
	assertpop(t, q, 25)
	assertpop(t, q, 30)
	assertpop(t, q, 35)
}
