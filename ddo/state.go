package ddo

// State types are defined by the user based on recursive models.
type State interface {
	Cost() int64
	IsSolved() bool
	Next(inferenceDual State, incumbent State) []State
	Infer() *Diagram
	Relax() *Diagram
	Restrict() *Diagram
}

// ByCost allows []State slices to be sorted by their costs.
type ByCost []State

func (s ByCost) Len() int {
	return len(s)
}

func (s ByCost) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ByCost) Less(i, j int) bool {
	return s[i].Cost() < s[j].Cost()
}
