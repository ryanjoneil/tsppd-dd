package tsppd

import (
	"encoding/json"
	"math"
	"strings"
)

// Problem represents a TSPPD instance.
type Problem struct {
	Name       string
	Comment    string
	Nodes      []string
	Precedence map[string]string
	Edges      [][]int64

	index map[string]int
}

// Decode converts a JSON byte array into a TSPPD Problem instance.
func Decode(b []byte) (Problem, error) {
	var p Problem
	if err := json.Unmarshal(b, &p); err != nil {
		return Problem{}, err
	}

	p.init()
	return p, nil
}

// IsEmpty is true if a Problem has no data.
func (p *Problem) IsEmpty() bool {
	return len(p.Nodes) == 0 && len(p.Precedence) == 0 && len(p.Edges) == 0
}

// Index returns the index of a node.
func (p *Problem) Index(node string) (int, bool) {
	i, ok := p.index[node]
	return i, ok
}

// Precedes returns true if node1 precedes node2 in any feasible path.
func (p *Problem) Precedes(node1, node2 string) bool {
	return p.Precedence[node1] == node2
}

// IsStart returns true if a node is +0.
func (p *Problem) IsStart(node string) bool {
	return node == "+0"
}

// IsEnd returns true if a node is -0.
func (p *Problem) IsEnd(node string) bool {
	return node == "-0"
}

// IsPickup returns true if a node is a pickup (i.e. contains +, not +0).
func (p *Problem) IsPickup(node string) bool {
	return strings.Contains(node, "+") && !p.IsStart(node)
}

// IsDelivery returns true if a node is a delivery (i.e. contains -, not -0).
func (p *Problem) IsDelivery(node string) bool {
	return strings.Contains(node, "-") && !p.IsEnd(node)
}

// IsFeasible returns true if a directed edge from node1 to node2 is feasible.
func (p *Problem) IsFeasible(node1, node2 string) bool {
	// +0 has no predecessor. -0 has no successor.
	if p.IsStart(node2) || p.IsEnd(node1) {
		return false
	}

	// Nodes can't connect to themselves.
	if node1 == node2 {
		return false
	}

	// Precedence relations can't be violated.
	if p.Precedes(node2, node1) {
		return false
	}

	// +0 can't connect to a delivery or directly to the end node.
	if p.IsStart(node1) && (p.IsDelivery(node2) || p.IsEnd(node2)) {
		return false
	}

	// Pickups can't connect to -0.
	if p.IsPickup(node1) && p.IsEnd(node2) {
		return false
	}

	return true
}

// FeasibleEdges returns the possible end nodes starting at a given node.
func (p *Problem) FeasibleEdges(node string) []string {
	edges := []string{}
	for _, end := range p.Nodes {
		if p.IsFeasible(node, end) {
			edges = append(edges, end)
		}
	}
	return edges
}

// Cost returns the cost of a directed arc from node1 to node2.
func (p *Problem) Cost(node1, node2 string) (int64, bool) {
	row, okRow := p.Index(node1)
	col, okCol := p.Index(node2)
	if !okRow || !okCol {
		return math.MaxInt16, false
	}
	return p.Edges[row][col], true
}

func (p *Problem) init() {
	p.index = map[string]int{}
	for index, node := range p.Nodes {
		p.index[node] = index
	}
}
