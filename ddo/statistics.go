package ddo

// Statistics represent solver execution information.
type Statistics struct {
	ClockSeconds float64
	CPUSeconds   float64
	Optimal      bool
	Fails        uint64
	Nodes        uint64
}
