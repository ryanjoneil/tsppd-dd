package ddo

const (
	exact   = iota
	failed  = iota
	relaxed = iota
)

// Bounds provide access to dual and primal bounding data.
type Bounds struct {
	Root           State
	InferenceDual  State
	RelaxationDual State
	Primal         State
	label          uint
}

// IsExact returns true if bounds are exact. That means either the
// dual bound == primal bound or the space has been fully searched.
func (b *Bounds) IsExact() bool {
	return b.label == exact
}

// IsFailed returns true if bounds or inference indicate a section
// of the search tree can be fathomed.
func (b *Bounds) IsFailed() bool {
	return b.label == failed
}

// IsRelaxed returns true dual bound < primal bound.
func (b *Bounds) IsRelaxed() bool {
	return b.label == relaxed
}

// DualBound provides dual bounds for a set of state.
func (b *Bounds) DualBound() int64 {
	var dual int64
	if b.InferenceDual != nil && b.InferenceDual.Cost() > dual {
		dual = b.InferenceDual.Cost()
	}
	if b.RelaxationDual != nil && b.RelaxationDual.Cost() > dual {
		dual = b.RelaxationDual.Cost()
	}
	return dual
}

// PrimalBound provides primal bounds for a set of state.
func (b *Bounds) PrimalBound() int64 {
	return b.Primal.Cost()
}

func (b *Bounds) betterThan(incumbent State) bool {
	return !b.IsFailed() && b.Primal.IsSolved() && (incumbent == nil || b.PrimalBound() < incumbent.Cost())
}
