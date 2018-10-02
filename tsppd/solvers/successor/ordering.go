package successor

import (
	"math"
	"sort"
)

func (s *State) initOrdering(ordering string) {
	s.ordering = []int{}

	switch ordering {
	case "input":
		s.initOrderingInput()
	case "greedy":
		s.initOrderingGreedy()
	case "regret":
		s.initOrderingRegret()
	}
}

func (s *State) initOrderingInput() {
	for index, node := range s.problem.Nodes {
		if !s.problem.IsEnd(node) {
			s.ordering = append(s.ordering, index)
		}
	}
}

func (s *State) initOrderingGreedy() {
	greedyIndexCosts := []indexCost{}

	for index1, node1 := range s.problem.Nodes {
		if s.problem.IsEnd(node1) {
			continue
		}

		var minCost int64 = math.MaxInt64
		for _, index2 := range s.feasible(index1) {
			cost, _ := s.problem.Cost(node1, s.problem.Nodes[index2])
			if cost < minCost {
				minCost = cost
			}
		}

		greedyIndexCosts = append(greedyIndexCosts, indexCost{index1, minCost})
	}

	sort.Sort(byIndexCost(greedyIndexCosts))
	for _, gic := range greedyIndexCosts {
		s.ordering = append(s.ordering, gic.index)
	}
}

func (s *State) initOrderingRegret() {
	regretIndexCosts := []indexCost{}

	for index1, node1 := range s.problem.Nodes {
		if s.problem.IsEnd(node1) {
			continue
		}

		var minCost1 int64 = math.MaxInt64
		var minCost2 int64 = math.MaxInt64
		for _, index2 := range s.feasible(index1) {
			cost, _ := s.problem.Cost(node1, s.problem.Nodes[index2])
			if cost < minCost1 {
				minCost1, minCost2 = cost, minCost1
			} else if cost < minCost2 {
				minCost2 = cost
			}
		}

		regretIndexCosts = append(regretIndexCosts, indexCost{index1, minCost2 - minCost1})
	}

	sort.Sort(sort.Reverse(byIndexCost(regretIndexCosts)))
	for _, gic := range regretIndexCosts {
		s.ordering = append(s.ordering, gic.index)
	}
}

type indexCost struct {
	index int
	cost  int64
}

type byIndexCost []indexCost

func (b byIndexCost) Len() int {
	return len(b)
}

func (b byIndexCost) Less(i, j int) bool {
	return b[i].cost < b[j].cost
}

func (b byIndexCost) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
