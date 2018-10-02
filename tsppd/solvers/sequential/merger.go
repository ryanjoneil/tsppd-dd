package sequential

import (
	"sort"

	"github.com/ryanjoneil/tsppd-dd/ddo"
)

// MaxCostRelaxationMerger combines the top states by cost.
func MaxCostRelaxationMerger(states []ddo.State, width uint) []ddo.State {
	sort.Sort(ddo.ByCost(states))

	feasible := map[string]bool{}
	for _, state := range states[width-1:] {
		s := state.(*State)
		for _, next := range s.feasible {
			feasible[next] = true
		}
	}

	mergedFeasible := []string{}
	for next := range feasible {
		mergedFeasible = append(mergedFeasible, next)
	}

	mergedStates := []ddo.State{}
	for _, state := range states[:width-1] {
		mergedStates = append(mergedStates, state)
	}

	lastState := states[width-1].(*State)
	mergedStates = append(mergedStates, &State{
		cost:     lastState.cost,
		feasible: mergedFeasible,
		node:     lastState.node,
		parent:   lastState.parent,
		problem:  lastState.problem,
	})

	return mergedStates
}
