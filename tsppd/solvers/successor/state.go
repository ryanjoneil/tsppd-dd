package successor

import (
	"fmt"

	"github.com/ryanjoneil/tsppd-dd/ddo"
	"github.com/ryanjoneil/tsppd-dd/tsppd"
	"github.com/ryanjoneil/tsppd-dd/tsppd/solvers/apdual"
)

// State represents a current feasible path order.
type State struct {
	cost int64

	// These are indexed the same way problem.Nodes is.
	domain  []int     // Unused values for next, i.e. {j | there is no next[i] = j}
	partial []*[]bool // partial[i] = nodes contained in partial route of which i is a part
	prev    []int     // prev[i] = j if (j i)
	next    []int     // next[i] = j if (i j)
	pred    []*[]bool // pred[i] = {nodes that must precede i}
	succ    []*[]bool // succ[i] = {nodes that must succeed i}

	// This is the order we assign to next in.
	ordering []int
	orderIdx int

	problem   *tsppd.Problem
	verbosity uint
	width     uint
	ap        *apdual.State
}

// CreateRootState makes the initial state for a successor DD TSPPD solver.
func CreateRootState(problem *tsppd.Problem, infer, relax, ordering string, width, verbosity uint) *State {
	var ap *apdual.State
	if infer == "ap" {
		ap = apdual.CreateAPDualState(problem)
	}

	s := &State{
		cost: 0,

		problem:   problem,
		verbosity: verbosity,
		width:     width,
		ap:        ap,
	}

	s.initDomain()
	s.initPartial()
	s.initPrevNext()
	s.initPredSucc()

	s.initOrdering(ordering)

	if s.verbosity == 2 {
		s.print()
	}

	return s
}

// Cost returns the cost of the complete or partial path represented by a State.
func (s *State) Cost() int64 {
	return s.cost
}

// IsSolved returns true if this state is a final solution.
func (s *State) IsSolved() bool {
	return len(s.domain) == 0
}

// Next creates the next feasible states accessible from a State.
func (s *State) Next(inferenceDual ddo.State, incumbent ddo.State) []ddo.State {
	if s.orderIdx >= len(s.ordering) {
		return []ddo.State{}
	}

	states := make([]ddo.State, 0, len(s.domain))
	index1 := s.ordering[s.orderIdx]

	for _, index2 := range s.feasible(index1) {
		if inferenceDual != nil && inferenceDual.(*apdual.State).FilterIndex(index1, index2, incumbent) {
			continue
		}

		nextPartial := s.nextPartial(index1, index2)

		state := &State{
			cost: s.nextCost(index1, index2),

			domain:  s.nextDomain(index2),
			partial: nextPartial,
			prev:    s.nextPrev(index1, index2),
			next:    s.nextNext(index1, index2),
			pred:    s.nextPred(index1, index2, nextPartial[index1]),
			succ:    s.nextSucc(index1, index2, nextPartial[index1]),

			ordering: s.ordering,
			orderIdx: s.orderIdx + 1,

			problem:   s.problem,
			verbosity: s.verbosity,
			width:     s.width,
			ap:        s.ap,
		}

		state.inferPred(index1)
		state.inferSucc(index1)

		states = append(states, state)

		if s.verbosity == 2 {
			fmt.Println()
			fmt.Printf("(%s %s)\n", s.problem.Nodes[index1], s.problem.Nodes[index2])
			states[len(states)-1].(*State).print()
		}
	}

	return states
}

// Solution returns the full or partial solution of a sequential TSPPD State.
func (s *State) Solution() *tsppd.Solution {
	path := []string{}

	current := 0
	for current >= 0 {
		path = append(path, s.problem.Nodes[current])
		current = s.next[current]
	}

	return &tsppd.Solution{
		Path:    path,
		Problem: s.problem,
	}
}

// Infer creates an inference diagram.
func (s *State) Infer() *ddo.Diagram {
	if s.ap == nil {
		return nil
	}

	if s.orderIdx > 0 {
		s.ap = s.ap.Copy()

		index1 := s.ordering[s.orderIdx-1]
		index2 := s.next[index1]

		// We can't connect anything but index to index2.
		for index3 := 0; index3 < len(s.problem.Nodes); index3++ {
			if index3 != index1 {
				s.ap.Remove(index3, index2)
			}
		}

		// index1's predecessors can't connect to its successors.
		preds := make([]int, 0, len(s.next))
		succs := make([]int, 0, len(s.next))

		for index3, v := range *s.pred[index1] {
			if v {
				preds = append(preds, index3)
			}
		}

		for index3, v := range *s.succ[index1] {
			if v {
				succs = append(succs, index3)
			}
		}

		for _, index3 := range preds {
			for _, index4 := range succs {
				s.ap.Remove(index3, index4)
			}
		}

		s.ap.Solve()
	}
	return ddo.CreateDiagram(s.ap, []ddo.Merger{}, s.width)
}

// Relax creates a relaxation diagram.
func (s *State) Relax() *ddo.Diagram {
	return nil
}

// Restrict creates a restriction diagram.
func (s *State) Restrict() *ddo.Diagram {
	return ddo.CreateDiagram(s, []ddo.Merger{ddo.MaxCostRestrictionMerger}, s.width)
}
