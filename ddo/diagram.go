package ddo

// Diagram instances
type Diagram struct {
	Layer   *Layer
	Mergers []Merger
	Width   uint
}

// CreateDiagram makes a Decision Diagram that minimizes some objective.
func CreateDiagram(state State, mergers []Merger, width uint) *Diagram {
	return &Diagram{
		Layer:   CreateRootLayer(state, mergers, width),
		Mergers: mergers,
		Width:   width,
	}
}

// Next returns the next layer of a Diagram.
func (d *Diagram) Next(inferenceDual State, incumbent State) *Layer {
	d.Layer = d.Layer.Next(inferenceDual, incumbent)
	return d.Layer
}

// IsDone returns true is the current layer is empty.
func (d *Diagram) IsDone() bool {
	return d.Layer.IsEmpty()
}
