package ddo

// Logger functions give the user or other systems insight into solver execution.
type Logger func(bounds *Bounds, stats Statistics)
