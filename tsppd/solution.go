package tsppd

import "math"

// Solution represents a TSPPD path.
type Solution struct {
	Problem *Problem
	Path    []string
}

// Cost computes the cost of a solution.
func (s *Solution) Cost() (int64, bool) {
	var cost int64
	for i := 0; i < len(s.Path)-1; i++ {
		c, ok := s.Problem.Cost(s.Path[i], s.Path[i+1])
		if !ok {
			return math.MaxInt64, false
		}
		cost += c
	}
	return cost, true
}
