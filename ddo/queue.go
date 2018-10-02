package ddo

import (
	"container/list"
	"sort"
)

type node struct {
	state  State
	dual   int64
	primal int64
}

type nodevec []node

func (nv nodevec) Len() int {
	return len(nv)
}

func (nv nodevec) Less(i, j int) bool {
	if nv[i].dual == nv[j].dual {
		return nv[i].primal < nv[j].primal
	}
	return nv[i].dual < nv[j].dual
}

func (nv nodevec) Swap(i, j int) {
	nv[i], nv[j] = nv[j], nv[i]
}

type queue struct {
	lnv *list.List
	nn  int
}

func createQueue(root State) *queue {
	lnv := list.New()
	lnv.PushBack(nodevec{node{state: root}})
	return &queue{lnv: lnv, nn: 1}
}

func (q *queue) len() int {
	return q.nn
}

func (q *queue) extend(r nodevec, incumbent State) {
	newR := make(nodevec, 0, r.Len())
	for _, s := range r {
		if incumbent == nil || (s.dual < incumbent.Cost() && s.primal < incumbent.Cost()) {
			newR = append(newR, s)
		}
	}
	sort.Sort(newR)

	if len(newR) > 0 {
		q.lnv.PushFront(newR)
		q.nn += len(newR)
	}
}

func (q *queue) pop() node {
	front := q.lnv.Front()
	nv := front.Value.(nodevec)
	n := nv[0]

	if len(nv) > 1 {
		front.Value = nv[1:]
	} else {
		q.lnv.Remove(front)
	}

	q.nn--
	return n
}
