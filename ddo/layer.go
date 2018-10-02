package ddo

// Layer instances represent layers within a Diagram.
type Layer struct {
	Depth   uint
	Mergers []Merger
	States  []State
	Width   uint
	IsExact bool
}

// CreateRootLayer builds a new layer with depth 0 and a single state.
func CreateRootLayer(state State, mergers []Merger, width uint) *Layer {
	states := []State{state}

	return &Layer{
		Depth:   0,
		Mergers: mergers,
		States:  states,
		Width:   width,
		IsExact: true,
	}
}

// Best returns the state with min cost.
func (l *Layer) Best() State {
	var best State

	for _, state := range l.States {
		if best == nil || state.Cost() < best.Cost() {
			best = state
		}
	}

	return best
}

// IsEmpty returns true is a Layer has no states.
func (l *Layer) IsEmpty() bool {
	return len(l.States) == 0
}

// Next builds the next Layer in a Diagram.
func (l *Layer) Next(inferenceDual State, incumbent State) *Layer {
	size := 0
	nextStateSlices := make([][]State, 0, len(l.States))
	for _, state := range l.States {
		next := state.Next(inferenceDual, incumbent)
		nextStateSlices = append(nextStateSlices, next)
		size += len(next)
	}

	nextStates := make([]State, 0, size)
	for _, next := range nextStateSlices {
		for _, nextState := range next {
			nextStates = append(nextStates, nextState)
		}
	}

	mergedStates := l.mergeStates(nextStates)

	return &Layer{
		Depth:   l.Depth + 1,
		Mergers: l.Mergers,
		States:  mergedStates,
		Width:   l.Width,
		IsExact: l.IsExact && len(mergedStates) == len(nextStates),
	}
}

func (l *Layer) mergeStates(states []State) []State {
	if l.Width == 0 {
		return states
	}

	for uint(len(states)) > l.Width {
		for _, merger := range l.Mergers {
			states = merger(states, l.Width)
			if uint(len(states)) <= l.Width {
				break
			}
		}
	}
	return states
}
