package detect

// Confidence represents how certain a detector is about the format.
type Confidence int

const (
	None Confidence = iota
	Low
	Medium
	High
)
