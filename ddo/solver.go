package ddo

//#include <time.h>
import "C"
import "time"

// Solver implements a basic Branch-and-Bound.
type Solver struct {
	Batch   int
	Workers int

	MaxMillis uint64
	MaxNodes  uint64

	queue     *queue
	root      State
	incumbent State

	logger    Logger
	wallStart time.Time
	cpuStart  _Ctype_long
	fails     uint64
	nodes     uint64
}

// CreateSolver constructs a basic Branch-and-Bound solver.
func CreateSolver(root State, logger Logger) *Solver {
	return &Solver{
		logger:    logger,
		queue:     createQueue(root),
		root:      root,
		wallStart: time.Now(),
		cpuStart:  C.clock(),
	}
}

// Minimize runs a full optimization from the root node.
func (s *Solver) Minimize() State {
	done := false
	for !done && s.queue.len() > 0 {
		boundsChan := make(chan []*Bounds)

		for i := 0; i < s.Workers; i++ {
			states := s.batch()
			go func(s *Solver, states []State) {
				b := []*Bounds{}
				for _, state := range states {
					b = append(b, s.bound(state))
				}
				boundsChan <- b
			}(s, states)
		}

		splitstates := nodevec{}
		for i := 0; i < s.Workers; i++ {
			if s.stop() {
				done = true
				break
			}

			for _, b := range <-boundsChan {
				if b.IsFailed() {
					s.fails++
					continue
				}
				s.nodes++

				if b.betterThan(s.incumbent) {
					s.incumbent = b.Primal
					s.logger(b, Statistics{s.elapsedSeconds(), s.elapsedCPU(), false, s.fails, s.nodes})
				}

				if b.IsRelaxed() {
					for _, next := range b.Root.Next(b.InferenceDual, s.incumbent) {
						if s.better(next) {
							splitstates = append(splitstates, node{
								state:  next,
								dual:   b.DualBound(),
								primal: next.Cost(),
							})
						}
					}
				}
			}
		}
		s.queue.extend(splitstates, s.incumbent)
	}

	// If we proved optimality, then say so.
	if s.incumbent != nil && !s.stop() {
		s.logger(
			&Bounds{
				Root:           s.incumbent,
				InferenceDual:  s.incumbent,
				RelaxationDual: s.incumbent,
				Primal:         s.incumbent,
				label:          exact,
			},
			Statistics{s.elapsedSeconds(), s.elapsedCPU(), true, s.fails, s.nodes},
		)
	}

	return s.incumbent
}

func (s *Solver) batch() []State {
	size := int(s.Batch)
	if size < 1 {
		size = 1
	}
	if s.queue.len() < size {
		size = s.queue.len()
	}

	states := make([]State, 0, size)

	i := 0
	for len(states) < size && i < s.queue.len() {
		state := s.queue.pop().state

		if s.incumbent == nil || state.Cost() < s.incumbent.Cost() {
			states = append(states, state)
		} else {
			s.fails++
		}
	}

	return states
}

// Bound solves a relaxation and then a restriction based on the current Diagram state.
func (s *Solver) bound(state State) *Bounds {
	dualBound := state.Cost()

	var inferenceDual State
	var relaxationDual State
	primal := state

	inferenceDiagram := state.Infer()
	relaxationDiagram := state.Relax()
	restrictionDiagram := state.Restrict()

	// Construct new layers for all diagrams until the restriction is done.
	for !restrictionDiagram.IsDone() {
		if inferenceDiagram != nil {
			if inferenceDiagram.Layer.IsEmpty() {
				return &Bounds{state, inferenceDual, relaxationDual, primal, failed}
			}

			inferenceDual = inferenceDiagram.Layer.Best()
			if s.worse(inferenceDual) {
				return &Bounds{state, inferenceDual, relaxationDual, primal, failed}
			}

			if inferenceDual.Cost() > dualBound {
				dualBound = inferenceDual.Cost()
			}
		}

		if relaxationDiagram != nil {
			if relaxationDiagram.Layer.IsEmpty() {
				return &Bounds{state, inferenceDual, relaxationDual, primal, failed}
			}

			relaxationDual = relaxationDiagram.Layer.Best()
			if s.worse(relaxationDual) {
				return &Bounds{state, inferenceDual, relaxationDual, primal, failed}
			}

			if relaxationDual.Cost() > dualBound {
				dualBound = relaxationDual.Cost()
			}
		}

		primal = restrictionDiagram.Layer.Best()
		if restrictionDiagram.Layer.IsExact && s.worse(primal) {
			return &Bounds{state, inferenceDual, relaxationDual, primal, failed}
		}

		if inferenceDiagram != nil {
			inferenceDiagram.Next(inferenceDual, s.incumbent)
		}
		if relaxationDiagram != nil {
			relaxationDiagram.Next(inferenceDual, s.incumbent)
		}
		restrictionDiagram.Next(inferenceDual, s.incumbent)
	}

	if primal.IsSolved() {
		// Restriction solution should be valid.
		if dualBound < primal.Cost() {
			return &Bounds{state, inferenceDual, relaxationDual, primal, relaxed}
		}
		return &Bounds{state, inferenceDual, relaxationDual, primal, exact}

	} else if restrictionDiagram.Layer.IsExact {
		// If a restriction is infeasible and it is exact, we can fathom it.
		return &Bounds{state, inferenceDual, relaxationDual, primal, failed}
	}

	// Return parent state because that's all we have. Likely
	// the restriction diagram got cut off partway through search
	// and can still generate feasible solutions.
	return &Bounds{state, inferenceDual, relaxationDual, state, relaxed}
}

func (s *Solver) stop() bool {
	if s.MaxMillis > 0 && s.elapsedMilliSeconds() >= float64(s.MaxMillis) {
		return true
	}
	if s.MaxNodes > 0 && uint64(s.nodes+s.fails) >= s.MaxNodes {
		return true
	}
	return false
}

func (s *Solver) better(state State) bool {
	return s.incumbent == nil || state.Cost() < s.incumbent.Cost()
}

func (s *Solver) worse(state State) bool {
	return s.incumbent != nil && state.Cost() >= s.incumbent.Cost()
}

func (s *Solver) elapsedMilliSeconds() float64 {
	return float64(time.Since(s.wallStart)) / float64(time.Millisecond)
}

func (s *Solver) elapsedSeconds() float64 {
	return time.Since(s.wallStart).Seconds()
}

func (s *Solver) elapsedCPU() float64 {
	return float64(C.clock()-s.cpuStart) / float64(C.CLOCKS_PER_SEC)
}
