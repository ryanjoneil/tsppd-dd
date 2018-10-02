package successor

import (
	"fmt"

	"github.com/ryanjoneil/tsppd-dd/ddo"
)

func (s *State) printStates(states []ddo.State) {
	if s.verbosity != 2 {
		return
	}

	for _, state := range states {
		state.(*State).print()
	}
}

func (s *State) print() {
	for i := 0; i < 120; i++ {
		fmt.Print("-")
	}
	fmt.Println()

	strDomain := []string{}
	for _, index := range s.domain {
		strDomain = append(strDomain, s.problem.Nodes[index])
	}

	fmt.Printf("cost=%05d domain=%v partials=%v\n\n", s.Cost(), strDomain, s.partials())

	fmt.Print("      ")
	for _, node := range s.problem.Nodes {
		fmt.Printf("%4s", node)
	}
	fmt.Println()

	fmt.Print("prev: ")
	for _, index := range s.prev {
		if index >= 0 {
			fmt.Printf("%4s", s.problem.Nodes[index])
		} else {
			fmt.Print("    ")
		}
	}
	fmt.Println()

	fmt.Print("next: ")
	for _, index := range s.next {
		if index >= 0 {
			fmt.Printf("%4s", s.problem.Nodes[index])
		} else {
			fmt.Print("    ")
		}
	}
	fmt.Println()
	fmt.Println()

	fmt.Print("part: ")
	for index, p := range s.partial {
		if index > 0 {
			fmt.Print("      ")
		}
		fmt.Printf("%p ", s.partial[index])
		fmt.Printf("[%s] ", s.problem.Nodes[index])
		for i, v := range *p {
			if v {
				fmt.Printf("%s ", s.problem.Nodes[i])
			}
		}
		fmt.Println()
	}
	fmt.Println()

	fmt.Print("pred: ")
	for index, p := range s.pred {
		if index > 0 {
			fmt.Print("      ")
		}
		fmt.Printf("%p ", s.pred[index])
		fmt.Printf("[%s] ", s.problem.Nodes[index])
		for i, v := range *p {
			if v {
				fmt.Printf("%s ", s.problem.Nodes[i])
			}
		}
		fmt.Println()
	}
	fmt.Println()

	fmt.Print("succ: ")
	for index, p := range s.succ {
		if index > 0 {
			fmt.Print("      ")
		}
		fmt.Printf("%p ", s.succ[index])
		fmt.Printf("[%s] ", s.problem.Nodes[index])
		for i, v := range *p {
			if v {
				fmt.Printf("%s ", s.problem.Nodes[i])
			}
		}
		fmt.Println()
	}
	fmt.Println()

	for i := 0; i < 120; i++ {
		fmt.Print("-")
	}
	fmt.Println()
}

func (s *State) partials() [][]string {
	strPartials := [][]string{}

	for index := range s.prev {
		if s.prev[index] < 0 {
			strPartial := []string{}

			current := index
			for current >= 0 {
				strPartial = append(strPartial, s.problem.Nodes[current])
				current = s.next[current]
			}

			if len(strPartial) > 1 {
				strPartials = append(strPartials, strPartial)
			}
		}
	}

	return strPartials
}
