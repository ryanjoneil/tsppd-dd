package apdual

import (
	"github.com/ryanjoneil/tsppd-dd/ddo"
	"github.com/ryanjoneil/tsppd-dd/tsppd"
	"github.com/ryanjoneil/ap"
)

const big = 10 * 1000 * 1000

// State represents a current AP relaxation.
type State struct {
	ap      *ap.AP
	problem *tsppd.Problem
}

// CreateAPDualState creates a State that maintains an AP formulation.
func CreateAPDualState(problem *tsppd.Problem) *State {
	ap := createAP(problem)
	ap.Solve()

	return &State{
		ap:      ap,
		problem: problem,
	}
}

// Set returns a new AP State with an edge forced on.
func (s *State) Set(node1, node2 string) *State {
	index1, _ := s.problem.Index(node1)
	index2, _ := s.problem.Index(node2)
	return s.SetIndex(index1, index2)
}

// SetIndex returns a new AP State with an edge forced on.
func (s *State) SetIndex(index1, index2 int) *State {
	newState := s.Copy()
	for index3 := 0; index3 < len(s.problem.Nodes); index3++ {
		if index3 != index1 {
			newState.Remove(index3, index2)
		}
	}
	newState.Solve()
	return newState
}

// UnsetIndex returns a new AP State with an edge forced off if needed.
func (s *State) UnsetIndex(index1, index2 int) *State {
	newState := s.Copy()
	newState.Remove(index1, index2)
	newState.Solve()
	return newState
}

// Cost returns the objective value of an AP.
func (s *State) Cost() int64 {
	return s.ap.Z
}

// IsSolved will always be true.
func (s *State) IsSolved() bool {
	return true
}

// Next just returns the same state.
func (s *State) Next(inferenceDual ddo.State, incumbent ddo.State) []ddo.State {
	return []ddo.State{s}
}

// Infer doesn't do much for the AP since it is always optimal.
func (s *State) Infer() *ddo.Diagram {
	return ddo.CreateDiagram(s, []ddo.Merger{}, 0)
}

// Relax doesn't do much for the AP since it is always optimal.
func (s *State) Relax() *ddo.Diagram {
	return ddo.CreateDiagram(s, []ddo.Merger{}, 0)
}

// Restrict doesn't do much for the AP either.
func (s *State) Restrict() *ddo.Diagram {
	return ddo.CreateDiagram(s, []ddo.Merger{}, 0)
}

// Filter returns true if the given edge can't be in an optimal solution.
func (s *State) Filter(node1, node2 string, incumbent ddo.State) bool {
	if incumbent == nil {
		return false
	}
	index1, _ := s.problem.Index(node1)
	index2, _ := s.problem.Index(node2)
	return s.FilterIndex(index1, index2, incumbent)
}

// FilterIndex returns true if the given edge can't be in an optimal solution.
func (s *State) FilterIndex(index1, index2 int, incumbent ddo.State) bool {
	if incumbent == nil {
		return false
	}
	return s.ap.Z+s.ap.RC(index1, index2) >= incumbent.Cost()
}

func createAP(problem *tsppd.Problem) *ap.AP {
	ap := ap.Create(len(problem.Nodes))
	for i := range problem.Nodes {
		for j := range problem.Nodes {
			node1 := problem.Nodes[i]
			node2 := problem.Nodes[j]
			if problem.IsFeasible(node1, node2) || node1 == "-0" && node2 == "+0" {
				ap.A[i][j], _ = problem.Cost(node1, node2)
			} else {
				ap.A[i][j] = big
			}
		}
	}
	return ap
}

// Copy makes a complete copy of the state.
func (s *State) Copy() *State {
	return &State{
		ap:      s.ap.Copy(),
		problem: s.problem,
	}
}

// Remove takes an edge out of the feasible set by giving it a large cost.
func (s *State) Remove(index1, index2 int) {
	s.ap.Remove(index1, index2, big)
}

// Solve re-solves an AP relaxation.
func (s *State) Solve() {
	s.ap.Solve()
}
