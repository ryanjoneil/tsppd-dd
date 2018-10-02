package sequential

import (
	"fmt"

	"github.com/ryanjoneil/tsppd-dd/ddo"
)

func (s *State) printStates(states []ddo.State) {
	if s.verbosity != 2 {
		return
	}

	for _, state := range states {
		fmt.Printf("cost=%05d path=%v\n", state.Cost(), state.(*State).Solution().Path)
	}
}
