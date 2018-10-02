package tsppd

// State types for TSPPD return a TSPPD Solution.
type State interface {
	Solution() *Solution
}
