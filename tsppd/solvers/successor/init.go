package successor

func (s *State) initDomain() {
	s.domain = []int{}
	for index, node := range s.problem.Nodes {
		// Domain includes everything that can be assigned to next.
		if !s.problem.IsStart(node) {
			s.domain = append(s.domain, index)
		}
	}
}

func (s *State) initPartial() {
	s.partial = make([]*[]bool, 0, len(s.problem.Nodes))
	for index := range s.problem.Nodes {
		partial := make([]bool, len(s.problem.Nodes))
		partial[index] = true
		s.partial = append(s.partial, &partial)
	}
}

func (s *State) initPrevNext() {
	s.prev = []int{}
	s.next = []int{}
	for range s.problem.Nodes {
		s.prev = append(s.prev, -1)
		s.next = append(s.next, -1)
	}
}

func (s *State) initPredSucc() {
	s.pred = make([]*[]bool, 0, len(s.next))
	s.succ = make([]*[]bool, 0, len(s.next))

	for index1, node1 := range s.problem.Nodes {
		pred := make([]bool, len(s.next))
		succ := make([]bool, len(s.next))

		for index2, node2 := range s.problem.Nodes {
			if index2 == index1 {
				continue
			}

			if s.problem.IsStart(node1) {
				succ[index2] = true

			} else if s.problem.IsPickup(node1) {
				if s.problem.Precedes(node1, node2) || s.problem.IsEnd(node2) {
					succ[index2] = true
				}
				if s.problem.IsStart(node2) {
					pred[index2] = true
				}

			} else if s.problem.IsDelivery(node1) {
				if s.problem.Precedes(node2, node1) || s.problem.IsStart(node2) {
					pred[index2] = true
				}
				if s.problem.IsEnd(node2) {
					succ[index2] = true
				}

			} else if s.problem.IsEnd(node1) {
				pred[index2] = true
			}
		}

		s.pred = append(s.pred, &pred)
		s.succ = append(s.succ, &succ)
	}
}
