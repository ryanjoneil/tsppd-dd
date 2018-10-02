package successor

func (s *State) feasible(index1 int) []int {
	f := make([]int, 0, len(s.domain))

	for _, index2 := range s.domain {
		// (i j) would create a cycle
		if s.partial[index1] == s.partial[index2] {
			continue
		}

		if intersect(s.partial[index1], s.succ[index2]) || intersect(s.pred[index1], s.partial[index2]) {
			continue
		}

		if intersect(s.pred[index1], s.succ[index2]) || intersect(s.succ[index1], s.pred[index2]) {
			continue
		}

		f = append(f, index2)
	}

	return f
}

func (s *State) nextCost(index1, index2 int) int64 {
	node1 := s.problem.Nodes[index1]
	node2 := s.problem.Nodes[index2]
	cost, _ := s.problem.Cost(node1, node2)
	return s.cost + cost
}

func (s *State) nextDomain(index2 int) []int {
	domain := make([]int, 0, len(s.domain)-1)
	for _, index := range s.domain {
		if index != index2 {
			domain = append(domain, index)
		}
	}
	return domain
}

func (s *State) nextPartial(index1, index2 int) []*[]bool {
	u := union(s.partial[index1], s.partial[index2])

	partial := make([]*[]bool, len(s.partial))
	for index := range s.partial {
		if (*u)[index] {
			partial[index] = u
		} else {
			partial[index] = s.partial[index]
		}
	}

	return partial
}

func (s *State) nextPrev(index1, index2 int) []int {
	prev := make([]int, len(s.prev))
	for index, prevIndex := range s.prev {
		if index == index2 {
			prev[index] = index1
		} else {
			prev[index] = prevIndex
		}
	}
	return prev
}

func (s *State) nextNext(index1, index2 int) []int {
	next := make([]int, len(s.next))
	for index, nextIndex := range s.next {
		if index == index1 {
			next[index] = index2
		} else {
			next[index] = nextIndex
		}
	}
	return next
}

func (s *State) nextPred(index1, index2 int, partial *[]bool) []*[]bool {
	oldPred1, oldPred2 := s.pred[index1], s.pred[index2]
	u := unionMinus(oldPred1, oldPred2, partial)

	pred := make([]*[]bool, len(s.pred))
	for index := range s.pred {
		if (*partial)[index] {
			pred[index] = u
		} else {
			pred[index] = s.pred[index]
		}
	}

	return pred
}

func (s *State) nextSucc(index1, index2 int, partial *[]bool) []*[]bool {
	oldSucc1, oldSucc2 := s.succ[index1], s.succ[index2]
	u := unionMinus(oldSucc1, oldSucc2, partial)

	succ := make([]*[]bool, len(s.succ))
	for index := range s.succ {
		if (*partial)[index] {
			succ[index] = u
		} else {
			succ[index] = s.succ[index]
		}
	}

	return succ
}

func (s *State) inferPred(index int) {
	inferredPred := union(s.pred[index], s.partial[index])

	oldToNew := make([]*[]bool, len(s.succ))
	for successor, v := range *s.succ[index] {
		if !v {
			continue
		}

		if oldToNew[successor] != nil {
			s.pred[successor] = oldToNew[successor]
		} else {
			oldToNew[successor] = unionMinus(inferredPred, s.pred[successor], s.partial[successor])
			s.pred[successor] = oldToNew[successor]
		}
	}
}

func (s *State) inferSucc(index int) {
	inferredSucc := union(s.succ[index], s.partial[index])

	oldToNew := make([]*[]bool, len(s.succ))
	for predecessor, v := range *s.pred[index] {
		if !v {
			continue
		}

		if oldToNew[predecessor] != nil {
			s.succ[predecessor] = oldToNew[predecessor]
		} else {
			oldToNew[predecessor] = unionMinus(inferredSucc, s.succ[predecessor], s.partial[predecessor])
			s.succ[predecessor] = oldToNew[predecessor]
		}
	}
}

func intersect(set1, set2 *[]bool) bool {
	for i := range *set1 {
		if (*set1)[i] && (*set2)[i] {
			return true
		}
	}
	return false
}

func union(set1, set2 *[]bool) *[]bool {
	u := make([]bool, len(*set1))
	for i := range *set1 {
		if (*set1)[i] || (*set2)[i] {
			u[i] = true
		}
	}
	return &u
}

func unionMinus(set1, set2, out *[]bool) *[]bool {
	u := make([]bool, len(*set1))
	for i := range *set1 {
		if !(*out)[i] && ((*set1)[i] || (*set2)[i]) {
			u[i] = true
		}
	}
	return &u
}
