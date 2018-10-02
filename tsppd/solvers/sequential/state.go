package sequential

import (
	"github.com/ryanjoneil/tsppd-dd/ddo"
	"github.com/ryanjoneil/tsppd-dd/tsppd"
	"github.com/ryanjoneil/tsppd-dd/tsppd/solvers/apdual"
)

// State represents a current feasible path order.
type State struct {
	cost      int64
	feasible  []string
	node      string
	parent    *State
	problem   *tsppd.Problem
	verbosity uint
	width     uint
	ap        *apdual.State
	relax     bool
}

// CreateRootState makes the initial state for a sequential DD TSPPD solver.
func CreateRootState(problem *tsppd.Problem, infer, relax, ordering string, width, verbosity uint) *State {
	feasible := []string{}
	for _, n := range problem.Nodes {
		if problem.IsPickup(n) {
			feasible = append(feasible, n)
		}
	}

	var ap *apdual.State
	if infer == "ap" {
		ap = apdual.CreateAPDualState(problem)
	}

	state := &State{
		cost:      0,
		feasible:  feasible,
		node:      "+0",
		parent:    nil,
		problem:   problem,
		verbosity: verbosity,
		width:     width,
		ap:        ap,
		relax:     relax == "dd",
	}
	return state
}

// Cost returns the cost of the complete or partial path represented by a State.
func (s *State) Cost() int64 {
	return s.cost
}

// IsSolved returns true if this state is a final solution.
func (s *State) IsSolved() bool {
	return len(s.feasible) == 0
}

// Next creates the next feasible states accessible from a State.
func (s *State) Next(inferenceDual ddo.State, incumbent ddo.State) []ddo.State {
	states := make([]ddo.State, 0, len(s.problem.Nodes)/2)

	for _, next := range s.feasible {
		// Don't generate solutions that are worse than the current incumbent.
		arcCost, _ := s.problem.Cost(s.node, next)
		cost := s.Cost() + arcCost
		if incumbent != nil && cost >= incumbent.Cost() {
			continue
		}

		// AP reduced cost-based domain filtering.
		if inferenceDual != nil && inferenceDual.(*apdual.State).Filter(s.node, next, incumbent) {
			continue
		}

		states = append(states, &State{
			cost:      cost,
			feasible:  s.nextFeasible(next),
			node:      next,
			parent:    s,
			problem:   s.problem,
			verbosity: s.verbosity,
			width:     s.width,
			ap:        s.ap,
			relax:     s.relax,
		})
	}

	s.printStates(states)
	return states
}

// Solution returns the full or partial solution of a sequential TSPPD State.
func (s *State) Solution() *tsppd.Solution {
	rpath := []string{}
	for state := s; state != nil; state = state.parent {
		rpath = append(rpath, state.node)
	}

	path := []string{}
	for i := len(rpath) - 1; i >= 0; i-- {
		path = append(path, rpath[i])
	}

	return &tsppd.Solution{
		Path:    path,
		Problem: s.problem,
	}
}

// Infer creates an inference diagram.
func (s *State) Infer() *ddo.Diagram {
	if s.ap != nil {
		// AP Relaxation
		if s.parent != nil {
			s.ap = s.ap.Set(s.parent.node, s.node)
		}
		return ddo.CreateDiagram(s.ap, []ddo.Merger{}, s.width)
	}
	return nil
}

// Relax creates a relaxation diagram.
func (s *State) Relax() *ddo.Diagram {
	if s.relax {
		return ddo.CreateDiagram(s, []ddo.Merger{MaxCostRelaxationMerger}, s.width)
	}
	return nil
}

// Restrict creates a restriction diagram.
func (s *State) Restrict() *ddo.Diagram {
	return ddo.CreateDiagram(s, []ddo.Merger{ddo.MaxCostRestrictionMerger}, s.width)
}

// Node returns the last node in the route.
func (s *State) Node() string {
	return s.node
}

func (s *State) nextFeasible(next string) []string {
	if len(s.feasible) == 1 && s.problem.IsDelivery(s.feasible[0]) {
		return []string{"-0"}
	}

	// New feasible set = current - next node + delivery.
	feasible := make([]string, 0, len(s.feasible))

	for _, node := range s.feasible {
		if node != next {
			feasible = append(feasible, node)
		}
	}

	// If next node is a pickup, add delivery.
	if s.problem.IsPickup(next) {
		feasible = append(feasible, s.problem.Precedence[next])
	}

	return feasible
}
