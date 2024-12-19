package pq

import "slices"

// Down left:  2*ix+1
// Down right:  2*ix+2
// Up:  int(ix-1)/2

type container[T any] struct {
	item     T
	priority int
}

type HEAP_TYPE string

const (
	MIN_HEAP HEAP_TYPE = "MIN"
	MAX_HEAP HEAP_TYPE = "MAX"
)

type PriorityQueue[T any] struct {
	arr    []container[T]
	lastIx int
	order  HEAP_TYPE
}

func New[T any](ht HEAP_TYPE) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		arr:    make([]container[T], 0),
		lastIx: -1,
		order:  ht, // TODO: rename
	}
}

func (pq *PriorityQueue[T]) Empty() bool {
	return pq.lastIx >= 0
}

func (pq *PriorityQueue[T]) Push(val T, priority int) {
	pq.lastIx++
	pq.arr = slices.Insert(pq.arr, pq.lastIx, container[T]{val, priority}) // is this leaking memory if push/pop/push/pop/push/pop/... ?
	pq.heapifyUp(pq.lastIx)
}

func (pq *PriorityQueue[T]) Pop() (T, bool) {
	if pq.lastIx < 0 {
		var t T
		return t, false
	}

	out := pq.arr[0]
	pq.arr[0] = pq.arr[pq.lastIx]
	pq.lastIx--
	pq.heapifyDown(0)
	return out.item, true
}

func (pq *PriorityQueue[T]) UnordereredItems() []T {
	out := make([]T, pq.lastIx+1)
	for ix := range pq.arr {
		out = append(out, pq.arr[ix].item)
	}
	return out
}

func (pq *PriorityQueue[T]) heapifyUp(ix int) {
	if ix <= 0 {
		return
	}

	parent := int(ix-1) / 2
	if pq.compare(ix, parent) {
		pq.swap(parent, ix)
		pq.heapifyUp(parent)
	}
}

func (pq *PriorityQueue[T]) heapifyDown(ix int) {
	leftIx := 2*ix + 1
	rightIx := 2*ix + 2

	// no childs
	if pq.lastIx < leftIx {
		return
	}

	// only left child
	if pq.lastIx < rightIx {
		if pq.compare(leftIx, ix) {
			pq.swap(leftIx, ix)
			pq.heapifyDown(leftIx)
		}
	}

	// two childs
	var minChildIx int
	if pq.compare(leftIx, rightIx) {
		minChildIx = leftIx
	} else {
		minChildIx = rightIx
	}

	if pq.compare(minChildIx, ix) {
		pq.swap(minChildIx, ix)
		pq.heapifyDown(minChildIx)
	}
}

func (pq *PriorityQueue[T]) compare(a, b int) bool {
	if pq.order == MIN_HEAP {
		return pq.arr[a].priority < pq.arr[b].priority
	}
	return pq.arr[a].priority > pq.arr[b].priority
}

func (pq *PriorityQueue[T]) swap(a, b int) {
	tmp := pq.arr[a]
	pq.arr[a] = pq.arr[b]
	pq.arr[b] = tmp
}
