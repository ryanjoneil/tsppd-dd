package ddo

import (
	"sort"
)

// Merger functions map the states in a Layer to a smaller set of states.
type Merger func(states []State, width uint) []State

// MaxCostRestrictionMerger removes the max cost states and returns no more than width states.
func MaxCostRestrictionMerger(states []State, width uint) []State {
	sort.Sort(ByCost(states))
	return states[:width]
}
